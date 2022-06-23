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
	FileMap_FileMetaInfo      = "File"
	FileMap_SchedulerInfo     = "SchedulerMap"
	FileBank_UserSpaceList    = "UserSpaceList"
	FileBank_UserSpaceInfo    = "UserHoldSpaceDetails"
	FileBank_UserFilelistInfo = "UserHoldFileList"
	Sminer_PurchasedSpace     = "PurchasedSpace"
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
	File_Name   types.Bytes
	FileSize    types.U64
	FileHash    types.Bytes
	Public      types.Bool
	UserAddr    types.AccountID
	FileState   types.Bytes
	Backups     types.U8
	Downloadfee types.U128
	FileDupl    []FileDuplicateInfo
}

type FileDuplicateInfo struct {
	MinerId   types.U64
	BlockNum  types.U32
	ScanSize  types.U32
	Acc       types.AccountID
	MinerIp   types.Bytes
	DuplId    types.Bytes
	RandKey   types.Bytes
	BlockInfo []BlockInfo
}
type BlockInfo struct {
	BlockIndex types.Bytes
	BlockSize  types.U32
}

//---UserInfo
type UserSpaceListInfo struct {
	Size     types.U128 `json:"size"`
	Deadline types.U32  `json:"deadline"`
}

type UserStorageSpace struct {
	Purchased_space types.U128 `json:"purchased_space"`
	Used_space      types.U128 `json:"used_space"`
	Remaining_space types.U128 `json:"remaining_space"`
}
