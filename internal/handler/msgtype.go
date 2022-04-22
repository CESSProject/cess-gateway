package handler

const (
	//200
	Status_200_default = "success"
	Status_200_expired = "captcha has expired and a new captcha has been sent to your mailbox"

	//400
	Status_400_default     = "HTTP error"
	Status_400_EmailFormat = "Email format error"
	Status_400_captcha     = "captcha error"
	Status_400_EmailSmpt   = "Please check your email address and whether to enable SMTP service"

	//401
	Status_401_token   = "Unauthorized"
	Status_401_expired = "token expired"

	//403
	Status_403_expired    = "not enough space"
	Status_403_dufilename = "duplicate filename"

	//500
	Status_500_db         = "Server internal data error"
	Status_500_chain      = "Server internal chain data error"
	Status_500_unexpected = "Server unexpected error"
)

// http response message
type RespMsg struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// http response random number message
type RespRandomMsg struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Random1 int    `json:"random1"`
	Random2 int    `json:"random2"`
}

// Request structure when user registers
type ReqGrantMsg struct {
	Mailbox string `json:"mailbox"`
	Captcha int64  `json:"captcha"`
}

// Request structure when user get randomkey
type ReqRandomkeyMsg struct {
	Walletaddr string `json:"walletaddr"`
}

// user state structure
type UserStateMsg struct {
	TotalSpace   string `json:"totalSpace"`
	UsedSpace    string `json:"usedSpace"`
	FreeSpace    string `json:"freeSpace"`
	SpaceDetails []SpaceDetailsMsg
}

// user space details structure
type SpaceDetailsMsg struct {
	Size     uint64 `json:"size"`
	Deadline uint32 `json:"deadline"`
}

// Request structure when user registers
type ReqDeleteFileMsg struct {
	Token    string `json:"token"`
	Filename string `json:"filename"`
}
