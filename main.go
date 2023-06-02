package main

import (
	"bufio"
	"github.com/fatih/color"
	"orchestra/cmd"
	"orchestra/global"
	"orchestra/utils"
	"os"
	"strings"
)

func main() {
	// TODO: Implement Graphic User Interface
	/*if len(os.Args) == 1 { // GUI
		gui := "./.orchestra/virtuoso/virtuoso"
		if runtime.GOOS == "windows" {
			gui += ".exe"
		}
		_, e := os.Stat(gui)
		switch {
		case e == nil:
			global.Prompter.Info("Calling virtuoso services, please wait...")
			return
		case os.IsNotExist(e):
			// Not found
			global.Prompter.Fatal("Cannot find the virtuoso GUI executable.")
		default:
			return
		}
	} else { // CUI*/
	if len(os.Args) == 1 || os.Args[1] == "repl" {
		reader, fg := bufio.NewReader(os.Stdin), color.New(color.FgWhite)
		for {
			_, _ = fg.Print("orchestra> ")
			input, err := reader.ReadString('\n')
			input = input[:len(input)-global.LengthOfNewLine]
			if err == nil {
				switch input {
				case "quit", "exit":
					global.Prompter.Info("Good bye ~")
					return
				case "":
					continue
				}
				os.Args = utils.PrependString(strings.Split(input, " "), "")
				cmd.ParseCommands()
			}
		}
	} else {
		cmd.ParseCommands()
	}
	global.RuntimeSearcher.Close()
	//}
}
