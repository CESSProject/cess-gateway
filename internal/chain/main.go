package chain

import (
	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
)

func NewRpcClient(rpcaddr string) (*gsrpc.SubstrateAPI, error) {
	return gsrpc.NewSubstrateAPI(rpcaddr)
}
