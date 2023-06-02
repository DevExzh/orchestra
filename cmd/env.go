package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"orchestra/global"
)

var (
	envCommand = &cobra.Command{
		Use:     "runtimes",
		Aliases: []string{"env"},
		Short:   "Print all runtimes detected.",
		Run: func(cmd *cobra.Command, args []string) {
			for _, rt := range global.RuntimeSearcher.Runtimes() {
				if rt.VersionCount() == 0 {
					continue
				}
				_, _ = color.New(color.FgCyan).Println(rt.Name())
				for _, v := range rt.Versions() {
					fmt.Print("    ", "Version: ")
					_, _ = color.New(color.FgYellow).Print(v.Name)
					fmt.Print(", Location: ")
					_, _ = color.New(color.FgGreen).Print(v.Path)
					fmt.Print("\n")
				}
			}
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			global.RuntimeSearcher.SetLoadFromDatabase(false)
			return global.RuntimeSearcher.FindRuntimes()
		},
	}
)
