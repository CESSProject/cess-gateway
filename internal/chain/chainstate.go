package chain

import (
	. "cess-httpservice/internal/logger"
	"cess-httpservice/tools"
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/pkg/errors"
)

// Get miner information on the cess chain
func GetUserRegisterMsg(blocknumber uint64, walletadddr string) (RegisterMsg, error) {
	var (
		err error
		msg RegisterMsg
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

	bytes, err := tools.DecodeToPub(walletadddr)
	if err != nil {
		return msg, errors.Errorf("DecodeToPub [%v] [%v] %v", blocknumber, walletadddr, err)
	}

	if events.FileBank_UserAuth != nil {
		for i := 0; i < len(events.FileBank_UserAuth); i++ {
			if events.FileBank_UserAuth[i].User == types.NewAccountID(bytes) {
				msg.Acc = bytes
				msg.Random = events.FileBank_UserAuth[i].Random
				msg.Collrate = events.FileBank_UserAuth[i].Collrate
				return msg, nil
			}
		}
		return msg, errors.Errorf("[%v]events.FileBank_UserAuth data err", blocknumber)
	}
	return msg, errors.Errorf("[%v]events.FileBank_UserAuth not found", blocknumber)
}

// Get scheduler information on the cess chain
func GetSchedulerInfo() ([]SchedulerInfo, error) {
	var (
		err  error
		data []SchedulerInfo
	)
	api := getSubstrateAPI()
	defer func() {
		releaseSubstrateAPI()
		err := recover()
		if err != nil {
			Err.Sugar().Errorf("[panic] %v", err)
		}
	}()
	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return nil, errors.Wrapf(err, "[%v.%v:GetMetadataLatest]", State_FileMap, FileMap_SchedulerInfo)
	}

	key, err := types.CreateStorageKey(meta, State_FileMap, FileMap_SchedulerInfo)
	if err != nil {
		return nil, errors.Wrapf(err, "[%v.%v:CreateStorageKey]", State_FileMap, FileMap_SchedulerInfo)
	}

	ok, err := api.RPC.State.GetStorageLatest(key, &data)
	if err != nil {
		return nil, errors.Wrapf(err, "[%v.%v:GetStorageLatest]", State_FileMap, FileMap_SchedulerInfo)
	}
	if !ok {
		return data, errors.Errorf("[%v.%v:GetStorageLatest value is nil]", State_FileMap, FileMap_SchedulerInfo)
	}
	return data, nil
}

// Get file meta information on the cess chain
func GetFileMetaInfo(fileid int64) (FileMetaInfo, error) {
	var (
		err  error
		data FileMetaInfo
	)

	api := getSubstrateAPI()
	defer func() {
		releaseSubstrateAPI()
		err := recover()
		if err != nil {
			Err.Sugar().Errorf("[panic] %v", err)
		}
	}()

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:GetMetadataLatest]", State_FileBank, FileMap_FileMetaInfo)
	}

	id, err := types.EncodeToBytes(fmt.Sprintf("%v", fileid))
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:EncodeToBytes]", State_FileBank, FileMap_FileMetaInfo)
	}

	key, err := types.CreateStorageKey(meta, State_FileBank, FileMap_FileMetaInfo, id)
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:CreateStorageKey]", State_FileBank, FileMap_FileMetaInfo)
	}

	ok, err := api.RPC.State.GetStorageLatest(key, &data)
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:GetStorageLatest]", State_FileBank, FileMap_FileMetaInfo)
	}
	if !ok {
		return data, errors.Errorf("[%v.%v:GetStorageLatest value is nil]", State_FileBank, FileMap_FileMetaInfo)
	}
	return data, nil
}

// Get user information on the cess chain
func GetSpaceDetailsInfo(wallet string) ([]UserSpaceListInfo, error) {
	var (
		err  error
		data []UserSpaceListInfo
	)

	api := getSubstrateAPI()
	defer func() {
		releaseSubstrateAPI()
		err := recover()
		if err != nil {
			Err.Sugar().Errorf("[panic] %v", err)
		}
	}()

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:GetMetadataLatest]", State_FileBank, FileBank_UserSpaceList)
	}

	bytes, err := tools.DecodeToPub(wallet)
	if err != nil {
		return data, err
	}

	key, err := types.CreateStorageKey(meta, State_FileBank, FileBank_UserSpaceList, bytes)
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:CreateStorageKey]", State_FileBank, FileBank_UserSpaceList)
	}

	ok, err := api.RPC.State.GetStorageLatest(key, &data)
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:GetStorageLatest]", State_FileBank, FileBank_UserSpaceList)
	}
	if !ok {
		return data, errors.Errorf("[%v.%v:GetStorageLatest value is nil]", State_FileBank, FileBank_UserSpaceList)
	}
	return data, nil
}

// Get user space information on the cess chain
func GetUserSpaceInfo(wallet string) (UserStorageSpace, error) {
	var (
		err  error
		data UserStorageSpace
	)

	api := getSubstrateAPI()
	defer func() {
		releaseSubstrateAPI()
		err := recover()
		if err != nil {
			Err.Sugar().Errorf("[panic] %v", err)
		}
	}()

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:GetMetadataLatest]", State_FileBank, FileBank_UserSpaceInfo)
	}

	bytes, err := tools.DecodeToPub(wallet)
	if err != nil {
		return data, err
	}

	key, err := types.CreateStorageKey(meta, State_FileBank, FileBank_UserSpaceInfo, bytes)
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:CreateStorageKey]", State_FileBank, FileBank_UserSpaceInfo)
	}

	ok, err := api.RPC.State.GetStorageLatest(key, &data)
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:GetStorageLatest]", State_FileBank, FileBank_UserSpaceInfo)
	}
	if !ok {
		return data, errors.Errorf("[%v.%v:GetStorageLatest value is nil]", State_FileBank, FileBank_UserSpaceInfo)
	}
	return data, nil
}

// Get file meta information on the cess chain
func GetFilelistInfo(wallet string) ([]types.Bytes, error) {
	var (
		err  error
		data []types.Bytes
	)

	api := getSubstrateAPI()
	defer func() {
		releaseSubstrateAPI()
		err := recover()
		if err != nil {
			Err.Sugar().Errorf("[panic] %v", err)
		}
	}()

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:GetMetadataLatest]", State_FileBank, FileBank_UserFilelistInfo)
	}

	bytes, err := tools.DecodeToPub(wallet)
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:DecodeToPub]", State_FileBank, FileBank_UserFilelistInfo)
	}

	key, err := types.CreateStorageKey(meta, State_FileBank, FileBank_UserFilelistInfo, bytes)
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:CreateStorageKey]", State_FileBank, FileBank_UserFilelistInfo)
	}

	ok, err := api.RPC.State.GetStorageLatest(key, &data)
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:GetStorageLatest]", State_FileBank, FileBank_UserFilelistInfo)
	}
	if !ok {
		return data, errors.Errorf("[%v.%v:GetStorageLatest value is nil]", State_FileBank, FileBank_UserFilelistInfo)
	}
	return data, nil
}

//Query sold space information on the cess chain
func QuerySoldSpace() (uint64, error) {
	var (
		err  error
		data types.U128
	)
	api := getSubstrateAPI()
	defer func() {
		releaseSubstrateAPI()
		err := recover()
		if err != nil {
			Err.Sugar().Errorf("[panic] %v", err)
		}
	}()
	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return 0, errors.Wrapf(err, "[%v.%v:GetMetadataLatest]", State_Sminer, Sminer_PurchasedSpace)
	}

	key, err := types.CreateStorageKey(meta, State_Sminer, Sminer_PurchasedSpace)
	if err != nil {
		return 0, errors.Wrapf(err, "[%v.%v:CreateStorageKey]", State_Sminer, Sminer_PurchasedSpace)
	}

	ok, err := api.RPC.State.GetStorageLatest(key, &data)
	if err != nil {
		return 0, errors.Wrapf(err, "[%v.%v:GetStorageLatest]", State_Sminer, Sminer_PurchasedSpace)
	}
	if !ok {
		return 0, nil
	}
	return data.Uint64(), nil
}

//Query total space information on the cess chain
func QueryTotalSpace() (uint64, error) {
	var (
		err  error
		data types.U128
	)
	api := getSubstrateAPI()
	defer func() {
		releaseSubstrateAPI()
		err := recover()
		if err != nil {
			Err.Sugar().Errorf("[panic] %v", err)
		}
	}()
	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return 0, errors.Wrapf(err, "[%v.%v:GetMetadataLatest]", State_Sminer, Sminer_TotalSpace)
	}

	key, err := types.CreateStorageKey(meta, State_Sminer, Sminer_TotalSpace)
	if err != nil {
		return 0, errors.Wrapf(err, "[%v.%v:CreateStorageKey]", State_Sminer, Sminer_TotalSpace)
	}

	ok, err := api.RPC.State.GetStorageLatest(key, &data)
	if err != nil {
		return 0, errors.Wrapf(err, "[%v.%v:GetStorageLatest]", State_Sminer, Sminer_TotalSpace)
	}
	if !ok {
		return 0, nil
	}
	return data.Uint64(), nil
}

// Get lastest block height
func GetLastestBlockHeight() (uint32, error) {
	api := getSubstrateAPI()
	defer func() {
		releaseSubstrateAPI()
		err := recover()
		if err != nil {
			Err.Sugar().Errorf("[panic] %v", err)
		}
	}()
	head, err := api.RPC.Chain.GetHeaderLatest()
	if err != nil {
		return 0, errors.Wrapf(err, "[GetHeaderLatest]")
	}
	return uint32(head.Number), nil
}
