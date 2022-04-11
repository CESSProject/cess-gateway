package handler

// http response message
type RespMsg struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
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
	Blocknumber int64  `json:"blocknumber"`
	Random2     int    `json:"random2"`
}

// Request structure when user get randomkey
type ReqRandomkeyMsg struct {
	Walletaddr string `json:"walletaddr"`
}
