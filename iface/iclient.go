package iface

type IClient interface {
	Start()
	Stop()
	GetConnID() uint32
	SendBuffMsg(uint32, []byte) error
}
