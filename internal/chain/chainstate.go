package chain

import (
	. "cess-httpservice/internal/logger"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/pkg/errors"
)

type Chain_RegisterMsg struct {
	Acc    types.Bytes `json:"acc"`
	Random types.U32   `json:"random"`
}

// Get miner information on the cess chain
func GetUserRegisterMsg(blocknumber uint64) (Chain_RegisterMsg, error) {
	var (
		err error
		msg Chain_RegisterMsg
	)
	api := getSubstrateAPI()
	defer func() {
		releaseSubstrateAPI()
		err := recover()
		if err != nil {
			Err.Sugar().Errorf("[panic] [%v] [%v]", blocknumber, err)
		}
	}()
	blockHash, err := api.RPC.Chain.GetBlockHash(blocknumber)
	if err != nil {
		return msg, errors.Wrap(err, "GetBlockHash err")
	}

	events := MyEventRecords{}

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return msg, errors.Errorf("GetMetadataLatest [%v] [%v]", blocknumber, err)
	}

	keye, err := types.CreateStorageKey(meta, "System", "Events", nil)
	if err != nil {
		return msg, errors.Errorf("CreateStorageKey [%v] [%v]", blocknumber, err)
	}

	h, err := api.RPC.State.GetStorageRaw(keye, blockHash)
	if err != nil {
		return msg, errors.Errorf("GetStorageRaw [%v] [%v]", blocknumber, err)
	}
	err = types.EventRecordsRaw(*h).DecodeEventRecords(meta, &events)
	if err != nil {
		Out.Sugar().Infof("[%v]Decode event err:%v", blocknumber, err)
	}
	// TODO:Waiting for the chain to define the interface
	// if events.FileMap_RegistrationUser != nil {
	// 	for i := 0; i < len(events.FileMap_RegistrationUser); i++ {

	// 	}
	// 	return msg, errors.Errorf("[%v]events.FileMap_RegistrationUser data err", blocknumber)
	// }
	return msg, errors.Errorf("[%v]events.FileMap_RegistrationUser not found", blocknumber)
}
