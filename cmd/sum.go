package cmd

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"hash"
	"io"
	"orchestra/global"
	"os"
	"strings"
)

var (
	sumCommand = &cobra.Command{
		Use:   "sum",
		Short: "Calculate hash summary of files.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				global.Prompter.Fatal("Files must be provided to calculate summary.")
			}
			var data [][]string
			switch strings.ToLower(flagAlgorithm) {
			case "md5":
				data = calcHash(md5.New(), args)
			case "sha1":
				data = calcHash(sha1.New(), args)
			case "sha256":
				data = calcHash(sha256.New(), args)
			case "sha512":
				data = calcHash(sha512.New(), args)
			default:
				global.Prompter.Fatal("Unknown algorithm \"", args[0], "\".")
				return
			}
			_ = pterm.DefaultTable.
				WithHasHeader().
				WithData(data).
				WithBoxed(true).
				WithLeftAlignment().
				Render()
		},
	}

	equalCommand = &cobra.Command{
		Use:   "eql",
		Short: "Checks if multiple files are equal.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				global.Prompter.Fatal("At least two files must be provided.")
				return
			}
			s := make(chan int64, len(args))
			defer close(s)
			for _, arg := range args {
				stat, err := os.Stat(arg)
				if os.IsNotExist(err) {
					global.Prompter.Fatal("File \"", arg, "\" cannot be found")
					return
				}
				s <- stat.Size()
			}
			var c = <-s
			for d := range s {
				if d != c {
					global.Prompter.Info("File sizes do not match.")
					return
				}
			}
			h := calcHash(sha256.New(), args)
			h1 := h[1][1]
			for i := 2; i < len(h); i++ {
				if h[i][1] != h1 {
					global.Prompter.Info("These files are not equal.")
					return
				}
			}
			global.Prompter.Info("Files are equivalent.")
		},
	}

	flagAlgorithm string
)

func initSumCommand() {
	flag := sumCommand.Flags()
	flag.StringVarP(
		&flagAlgorithm,
		"algorithm",
		"a",
		"sha256",
		"Algorithm employed in the calculation.",
	)
}

func calcHash(hash hash.Hash, files []string) [][]string {
	d := [][]string{{"File Name", "Hash"}}
	var j = 0
	for _, f := range files {
		file, err := os.Open(f)
		if err != nil {
			continue
		} else {
			if b, er := io.ReadAll(file); er == nil {
				hash.Write(b)
				d = append(d, []string{f, hex.EncodeToString(hash.Sum(nil))})
			} else {
				continue
			}
			j++
			_ = file.Close()
		}
	}
	return d
}
