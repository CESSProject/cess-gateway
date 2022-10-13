package tcp

type NetConn interface {
	// HandlerLoop 不能阻塞
	HandlerLoop()
	GetMsg() (*Message, bool)
	SendMsg(m *Message)
	Close() error
	IsClose() bool
}
