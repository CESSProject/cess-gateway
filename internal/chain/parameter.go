package chain

import "github.com/centrifuge/go-substrate-rpc-client/v4/types"

// cess chain state
const (
	State_Sminer      = "Sminer"
	State_SegmentBook = "SegmentBook"
	State_FileBank    = "FileBank"
	State_FileMap     = "FileMap"
)

// cess chain module method
const (
	Sminer_AllMinerItems      = "AllMiner"
	Sminer_MinerItems         = "MinerItems"
	Sminer_SegInfo            = "SegInfo"
	SegmentBook_ParamSet      = "ParamSet"
	SegmentBook_ConProofInfoA = "ConProofInfoA"
	SegmentBook_UnVerifiedA   = "UnVerifiedA"
	SegmentBook_UnVerifiedB   = "UnVerifiedB"
	SegmentBook_UnVerifiedC   = "UnVerifiedC"
	SegmentBook_UnVerifiedD   = "UnVerifiedD"
	FileMap_FileMetaInfo      = "File"
	FileMap_SchedulerInfo     = "SchedulerMap"
	FileBank_UserInfoMap      = "UserInfoMap"
	FileBank_UserSpaceInfo    = "UserHoldSpaceDetails"
	FileBank_UserFilelistInfo = "UserHoldFileList"
)

// cess chain Transaction name
const (
	ChainTx_SegmentBook_VerifyInVpa  = "SegmentBook.verify_in_vpa"
	ChainTx_SegmentBook_VerifyInVpb  = "SegmentBook.verify_in_vpb"
	ChainTx_SegmentBook_VerifyInVpc  = "SegmentBook.verify_in_vpc"
	ChainTx_SegmentBook_VerifyInVpd  = "SegmentBook.verify_in_vpd"
	ChainTx_SegmentBook_IntentSubmit = "SegmentBook.intent_submit"
	ChainTx_FileBank_Update          = "FileBank.update"
	ChainTx_FileMap_Add_schedule     = "FileMap.registration_scheduler"
	ChainTx_FileBank_PutMetaInfo     = "FileBank.update_dupl"
	ChainTx_FileBank_HttpUpload      = "FileBank.http_upload"
)

//---RegisterMsg
type RegisterMsg struct {
	Acc      types.Bytes `json:"acc"`
	Collrate types.U128  `json:"collrate"`
	Random   types.U32   `json:"random"`
}

//---SchedulerInfo
type SchedulerInfo struct {
	Ip    types.Bytes     `json:"ip"`
	Owner types.AccountID `json:"acc"`
}

//---FileMetaInfo
type FileMetaInfo struct {
	File_Name   types.Bytes         `json:"file_name"`   //File name
	FileSize    types.U64           `json:"file_size"`   //File size
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

//---UserInfo
type UserInfo struct {
	Collaterals types.U128 `json:"collaterals"`
}

type UserStorageSpace struct {
	Purchased_space types.U128 `json:"purchased_space"`
	Used_space      types.U128 `json:"used_space"`
	Remaining_space types.U128 `json:"remaining_space"`
}
