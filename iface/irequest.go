package iface

type IRequest interface {
	GetMsgID() uint32
	GetData() []byte
	GetClient() IClient
}
