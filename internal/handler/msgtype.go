package handler

// http response msg
type RespMsg struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Request structure when user registers
type RegistrationReq struct {
	Blocknumber int    `json:"blocknumber"`
	Walletaddr  string `json:"walletaddr"`
}
