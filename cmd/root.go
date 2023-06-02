package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"orchestra/global"
	"os"
	"path/filepath"
)

var (
	rootCommand = &cobra.Command{
		Use:     "orchestra",
		Aliases: []string{"orc"},
		Short:   "A powerful tool to do everything.",
		Long:    "Orchestra is a tool for dispatching scripts, managing applications.",
	}

	versionCommand = &cobra.Command{
		Use:   "version",
		Short: "Print the version of this tool.",
		Run: func(cmd *cobra.Command, args []string) {
			hc, w := color.New(color.FgHiCyan), color.New(color.FgYellow)
			_, _ = hc.Print("Orchestra: ")
			_, _ = w.Println("1.3.1")
			//_, _ = hc.Print("Core Library (Virtuoso): ")
			//_, _ = w.Println(global.CoreLibraryVersion())
			fmt.Println("Copyright (C) DevExzh (Ryker Zhu), licensed under GPLv3.")
		},
	}

	isInitialized = false
)

func initialize() {
	if isInitialized {
		return
	}

	/*
		Prepare database
	*/
	global.RuntimeSearcher.SetLoadFromDatabase(true)
	global.RuntimeSearcher.SetSaveToDatabase(true)
	d, _ := os.Getwd()
	d = filepath.Join(d, ".orchestra")
	_ = global.CreateDirectoryIfNotExist(d)
	global.Prompter.Error(global.RuntimeSearcher.SetDatabasePath(d))

	/*
		Set commands and flags
	*/
	flags := runCommand.Flags()
	flags.StringVarP(
		&flagLang,
		"language",
		"L",
		"auto",
		"The language that a executable employs.")

	initGetCommand()
	initSumCommand()

	rootCommand.AddCommand(
		runCommand,
		envCommand,
		versionCommand,
		getCommand,
		sumCommand,
		equalCommand,
	)
	isInitialized = true
}

func ParseCommands() {
	initialize()
	if err := runCommand.Execute(); err != nil {
		global.Prompter.Error(err)
	}
}
