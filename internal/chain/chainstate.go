package chain

import (
	. "cess-httpservice/internal/logger"
	"cess-httpservice/tools"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/pkg/errors"
)

type Chain_RegisterMsg struct {
	Acc      types.Bytes `json:"acc"`
	Collrate types.U128  `json:"collrate"`
	Random   types.U32   `json:"random"`
}

type SchedulerInfo struct {
	Ip    types.Bytes     `json:"ip"`
	Owner types.AccountID `json:"acc"`
}

type FileMetaInfo struct {
	//FileId      types.Bytes         `json:"acc"`         //File id
	File_Name   types.Bytes         `json:"file_name"`   //File name
	FileSize    types.U128          `json:"file_size"`   //File size
	FileHash    types.Bytes         `json:"file_hash"`   //File hash
	Public      types.Bool          `json:"public"`      //Public or not
	UserAddr    types.AccountID     `json:"user_addr"`   //Upload user's address
	FileState   types.Bytes         `json:"file_state"`  //File state
	Backups     types.U8            `json:"backups"`     //Number of backups
	Downloadfee types.U128          `json:"downloadfee"` //Download fee
	FileDupl    []FileDuplicateInfo `json:"file_dupl"`   //File backup information list
}

type FileDuplicateInfo struct {
	DuplId    types.Bytes     `json:"dupl_id"`    //Backup id
	RandKey   types.Bytes     `json:"rand_key"`   //Random key
	SliceNum  types.U16       `json:"slice_num"`  //Number of slices
	FileSlice []FileSliceInfo `json:"file_slice"` //Slice information list
}

type FileSliceInfo struct {
	SliceId   types.Bytes   `json:"slice_id"`   //Slice id
	SliceSize types.U32     `json:"slice_size"` //Slice size
	SliceHash types.Bytes   `json:"slice_hash"` //Slice hash
	FileShard FileShardInfo `json:"file_shard"` //Shard information
}

type FileShardInfo struct {
	DataShardNum  types.U8      `json:"data_shard_num"`  //Number of data shard
	RedunShardNum types.U8      `json:"redun_shard_num"` //Number of redundant shard
	ShardHash     []types.Bytes `json:"shard_hash"`      //Shard hash list
	ShardAddr     []types.Bytes `json:"shard_addr"`      //Store miner service addr list
	Peerid        []types.U64   `json:"wallet_addr"`     //Store miner wallet addr list
}

// Get miner information on the cess chain
func GetUserRegisterMsg(blocknumber uint64, walletadddr string) (Chain_RegisterMsg, error) {
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

	_, err = api.RPC.State.GetStorageLatest(key, &data)
	if err != nil {
		return nil, errors.Wrapf(err, "[%v.%v:GetStorageLatest]", State_FileMap, FileMap_SchedulerInfo)
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
		return data, errors.Wrapf(err, "[%v.%v:GetMetadataLatest]", State_FileMap, FileMap_FileMetaInfo)
	}

	id, err := types.EncodeToBytes(fileid)
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:EncodeToBytes]", State_FileMap, FileMap_FileMetaInfo)
	}

	key, err := types.CreateStorageKey(meta, State_FileMap, FileMap_FileMetaInfo, types.Bytes(id))
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:CreateStorageKey]", State_FileMap, FileMap_FileMetaInfo)
	}

	_, err = api.RPC.State.GetStorageLatest(key, &data)
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:GetStorageLatest]", State_FileMap, FileMap_FileMetaInfo)
	}
	return data, nil
}
