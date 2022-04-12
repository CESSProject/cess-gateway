package main

import (
	"cess-httpservice/configs"
	"cess-httpservice/internal/handler"

	"fmt"
	"os"
)

// init
func init() {
	if _, err := os.Stat(configs.BaseDir); err != nil {
		if err = os.MkdirAll(configs.BaseDir, os.ModeDir); err != nil {
			fmt.Printf("\x1b[%dm[err]\x1b[0m %v\n", 41, err)
			os.Exit(1)
		}
	}
	if _, err := os.Stat(configs.LogfileDir); err != nil {
		if err = os.MkdirAll(configs.LogfileDir, os.ModeDir); err != nil {
			fmt.Printf("\x1b[%dm[err]\x1b[0m %v\n", 41, err)
			os.Exit(1)
		}
	}
	if _, err := os.Stat(configs.DbDir); err != nil {
		if err = os.MkdirAll(configs.DbDir, os.ModeDir); err != nil {
			fmt.Printf("\x1b[%dm[err]\x1b[0m %v\n", 41, err)
			os.Exit(1)
		}
	}
	if _, err := os.Stat(configs.FileCacheDir); err != nil {
		if err = os.MkdirAll(configs.FileCacheDir, os.ModeDir); err != nil {
			fmt.Printf("\x1b[%dm[err]\x1b[0m %v\n", 41, err)
			os.Exit(1)
		}
	}
	if err := configs.ParseConfile(); err != nil {
		fmt.Printf("\x1b[%dm[err]\x1b[0m %v\n", 41, err)
		os.Exit(1)
	}
}

// Program entry
func main() {
	//init
	//logger.LogInit()
	// encryption.Init()
	// chain.Init()
	//start-up
	handler.Main()
}
