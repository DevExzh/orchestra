package cmd

import (
	"fmt"
	"github.com/pterm/pterm"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"net/url"
	"orchestra/global"
	"orchestra/network"
	"orchestra/utils"
)

var (
	getCommand = &cobra.Command{
		Use:   "get",
		Short: "Get online resources rapidly.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				global.Prompter.Fatal("URL must be provided to request the file.")
				return
			}
			getResource(args[0])
		},
	}

	// Flags for "get" command
	flagOutputPath  string
	flagThreadCount uint
	flagProxy       string
	flagMemoryCache uint
	flagSummaryOnly bool
)

func initGetCommand() {
	flags := getCommand.Flags()
	flags.StringVarP(
		&flagOutputPath,
		"output",
		"o",
		"./",
		"The output path where the file will be download.",
	)
	flags.UintVarP(
		&flagThreadCount,
		"thread",
		"t",
		0,
		"How many threads should be employed in the procedure.",
	)
	flags.StringVarP(
		&flagProxy,
		"proxy",
		"x",
		"",
		"URL of the proxy server.",
	)
	flags.UintVarP(
		&flagMemoryCache,
		"memory-cache",
		"M",
		0,
		"Count of bytes that every single thread will store in memory for cache.",
	)
	flags.BoolVarP(
		&flagSummaryOnly,
		"summary",
		"A",
		false,
		"Print summary but do not download the file.",
	)
}

func getResource(u string) {
	if ur, err := url.Parse(u); err != nil {
		global.Prompter.Error(err)
	} else {
		worker := network.Downloader{
			URL:            ur,
			ThreadCount:    flagThreadCount,
			TargetFilePath: flagOutputPath,
			EnableProxy:    flagProxy != "",
			ProxyURL:       flagProxy,
			AutoRename:     true,
		}
		global.Prompter.Error(worker.Init())
		pterm.DefaultTable.WithBoxed(true).WithLeftAlignment().WithData(pterm.TableData{
			{"Length", utils.AppendUnit(worker.TotalLength(), false, true, 2, true)},
			{"Host", ur.Host},
			{"File Name", worker.Suggested},
		})
		if !flagSummaryOnly {
			pb := progressbar.DefaultBytes(int64(worker.TotalLength()), "")
			go func(downloader *network.Downloader, bar *progressbar.ProgressBar) {
				for {
					select {
					case t := <-downloader.TaskResponses:
						switch t.Status {
						case network.TaskStarted:
							fmt.Println("[Thread", t.Id, "] Started.")
						case network.TaskFinished:
							fmt.Println("[Thread", t.Id, "] Finished.")
						case network.TaskUpdated:
							global.Prompter.Error(bar.Set64(int64(downloader.BytesDownloaded())))
						}
					case d := <-downloader.DownloaderResponses:
						switch d.Status {
						case network.DownloaderFinished:
							return
						default:
							global.Prompter.HandleError(downloader)
						}
					}
				}
			}(&worker, pb)
			global.Prompter.Error(worker.Start())
			worker.Wait()
		}
	}
}
