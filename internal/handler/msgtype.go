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
	UserId     int64  `json:"userId"`
	Deposit    string `json:"deposit"`
	TotalSpace string `json:"totalSpace"`
	UsedSpace  string `json:"usedSpace"`
	FreeSpace  string `json:"freeSpace"`
	Walletaddr string `json:"walletaddr"`
}
