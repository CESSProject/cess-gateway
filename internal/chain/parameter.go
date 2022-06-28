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
	FileMap_FileMetaInfo   = "File"
	FileMap_SchedulerInfo  = "SchedulerMap"
	FileBank_UserSpaceList = "UserSpaceList"
	FileBank_UserSpaceInfo = "UserHoldSpaceDetails"
	FileBank_UserFilelist  = "UserHoldFileList"
	Sminer_PurchasedSpace  = "PurchasedSpace"
)

// cess chain Transaction name
const (
	ChainTx_FileBank_Upload            = "FileBank.upload"
	ChainTx_FileBank_DeleteFile        = "FileBank.delete_file"
	ChainTx_FileBank_UploadDeclaration = "FileBank.upload_declaration"
)

//---RegisterMsg
type RegisterMsg struct {
	Acc      types.Bytes `json:"acc"`
	Collrate types.U128  `json:"collrate"`
	Random   types.U32   `json:"random"`
}

//---SchedulerInfo
type SchedulerInfo struct {
	Ip              types.Bytes
	Stash_user      types.AccountID
	Controller_user types.AccountID
}

//---FileMetaInfo
type FileMetaInfo struct {
	MinerId     types.U64         `json:"Miner_id"`
	FileSize    types.U64         `json:"File_size"`
	BlockNum    types.U32         `json:"Block_num"`
	ScanSize    types.U32         `json:"Scan_size"`
	SegmentSize types.U32         `json:"Segment_size"`
	MinerAcc    types.AccountID   `json:"Miner_acc"`
	MinerIp     types.Bytes       `json:"Miner_ip"`
	Users       []types.AccountID `json:"Users"`
	Names       []types.Bytes     `json:"Names"`
	FileState   types.Bytes       `json:"File_state"`
}

//---UserInfo
type UserSpaceListInfo struct {
	Size     types.U128 `json:"size"`
	Deadline types.U32  `json:"deadline"`
}

type UserStorageSpace struct {
	Purchased_space types.U128
	Used_space      types.U128
	Remaining_space types.U128
}

type UserFileList struct {
	File_hash types.Bytes
	File_size types.U64
}
