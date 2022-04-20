package handler

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
type ReqRegistrationMsg struct {
	Walletaddr  string `json:"walletaddr"`
	Blocknumber uint64 `json:"blocknumber"`
	Random2     int    `json:"random2"`
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
