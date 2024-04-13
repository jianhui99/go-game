package net

type Connection interface {
	Close()
	SendMessage(buf []byte) error
	GetSession() *Session
}

type MsgPack struct {
	Body []byte
	Cid  string
}
