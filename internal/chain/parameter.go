package chain

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
