package main

import (
	"cess-httpservice/internal/chain"
	"cess-httpservice/internal/encryption"
	"cess-httpservice/internal/handler"
	"cess-httpservice/internal/logger"
)

// Program entry
func main() {
	logger.Init()
	encryption.Init()
	chain.Init()
	handler.Main()
}
