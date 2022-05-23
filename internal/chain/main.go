package chain

import (
	"cess-gateway/configs"
	. "cess-gateway/internal/logger"
	"fmt"
	"os"
	"sync"
	"time"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
)

var (
	l *sync.Mutex
	r *gsrpc.SubstrateAPI
)

// init
func Init() {
	var err error
	r, err = gsrpc.NewSubstrateAPI(configs.Confile.RpcAddr)
	if err != nil {
		fmt.Printf("\x1b[%dm[err]\x1b[0m %v\n", 41, err)
		os.Exit(1)
	}
	l = new(sync.Mutex)
	go substrateAPIKeepAlive()
}

// KeepAlive
func substrateAPIKeepAlive() {
	var (
		err   error
		count uint8  = 0
		peer  uint64 = 0
		tk           = time.Tick(time.Second * 25)
	)

	for range tk {
		if count <= 1 {
			peer, err = healthcheck(r)
			if err != nil || peer == 0 {
				count++
			}
		}
		if count > 1 {
			count = 2
			r, err = gsrpc.NewSubstrateAPI(configs.ChainAddr)
			if err == nil {
				Err.Sugar().Errorf("%v", err)
			} else {
				count = 0
			}
		}
	}
}

// health check
func healthcheck(a *gsrpc.SubstrateAPI) (uint64, error) {
	defer func() {
		err := recover()
		if err != nil {
			Err.Sugar().Errorf("[panic]: %v", err)
		}
	}()
	h, err := a.RPC.System.Health()
	return uint64(h.Peers), err
}

// get SubstrateAPI and lock
func getSubstrateAPI() *gsrpc.SubstrateAPI {
	l.Lock()
	return r
}

// unlock
func releaseSubstrateAPI() {
	l.Unlock()
}
