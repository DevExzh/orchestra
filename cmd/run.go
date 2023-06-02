package cmd

import (
	"flag"
	"fmt"
	"github.com/spf13/cobra"
	"orchestra/global"
	"orchestra/utils"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	flagLang string

	runCommand = &cobra.Command{
		Use:   "run",
		Short: "Run scripts with appropriate runtime (interpreter, VM, etc.)",
		Run: func(cmd *cobra.Command, args []string) {
			autoExecuteByExtension(args[0])
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return global.RuntimeSearcher.FindRuntimes()
		},
	}
)

func autoExecuteByExtension(fileName string) {
	switch filepath.Ext(fileName) {
	case ".py":
		fmt.Println(global.RuntimeSearcher.Runtime("python").DefaultVersion().Path + " " + fileName)
	case ".jar":
		if jar, err := utils.IsRunnableJarFile(fileName,
			global.RuntimeSearcher.Runtime("java").Versions()); err == nil {
			str, er := filepath.Abs(fileName)
			if er != nil {
				return
			}
			global.Prompter.Error(run(jar, "-jar", str, "-Dfile-encoding=UTF-8"))
			global.Prompter.Info("Invoking Java Runtime...")

		} else {
			global.Prompter.Fatal(err.Error())
		}
	case ".js":
		fmt.Println(global.RuntimeSearcher.Runtime("node").DefaultVersion().Path + " " + fileName)
	default:
		global.Prompter.Warning("Unknown type of the executable \"" +
			flag.Arg(1) + "\". Use the operating system's default opening method...")
		args := flag.Args()
		if len(args) > 3 {
			global.Prompter.Error(run(flag.Arg(1), args[2:]...))
		} else {
			global.Prompter.Error(run(flag.Arg(1)))
		}
	}
}

func run(name string, args ...string) error {
	invoker := exec.Command(name, args...)
	invoker.Stdin = os.Stdin
	invoker.Stdout = os.Stdout
	invoker.Stderr = os.Stderr
	return invoker.Run()
}
