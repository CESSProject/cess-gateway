package main

import (
	"cess-httpservice/internal/chain"
	"cess-httpservice/internal/logger"
)

func main() {
	logger.Init()
	chain.Init()
}
