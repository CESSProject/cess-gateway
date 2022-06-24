package main

import (
	"cess-gateway/configs"
	"cess-gateway/internal/chain"
	"cess-gateway/internal/handler"
	"cess-gateway/tools"
	"flag"

	"fmt"
	"os"
)

// version
const VERSION = "CESS-Gateway v0.2.0 pre-release"

var printVersion bool

// init
func init() {
	flag.BoolVar(&printVersion, "v", false, "Print version number")
	flag.Parse()
	if printVersion {
		fmt.Println(VERSION)
		os.Exit(1)
	}

	if err := configs.ParseConfile(); err != nil {
		fmt.Printf("\x1b[%dm[err]\x1b[0m %v\n", 41, err)
		os.Exit(1)
	}

	if err := tools.CreatDirIfNotExist(configs.BaseDir); err != nil {
		fmt.Printf("\x1b[%dm[err]\x1b[0m %v\n", 41, err)
		os.Exit(1)
	}

	if err := tools.CreatDirIfNotExist(configs.LogfileDir); err != nil {
		fmt.Printf("\x1b[%dm[err]\x1b[0m %v\n", 41, err)
		os.Exit(1)
	}

	if err := tools.CreatDirIfNotExist(configs.DbDir); err != nil {
		fmt.Printf("\x1b[%dm[err]\x1b[0m %v\n", 41, err)
		os.Exit(1)
	}

	if err := tools.CreatDirIfNotExist(configs.FileCacheDir); err != nil {
		fmt.Printf("\x1b[%dm[err]\x1b[0m %v\n", 41, err)
		os.Exit(1)
	}
}

// Program entry
func main() {
	chain.Init()
	//start-up
	handler.Main()
}
