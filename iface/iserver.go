package iface

type IServer interface {
	Start()

	Stop()

	Serve()

	GetConnMgr() IConnManager

	AddRouter(uint32, IRouter)

	CallOnConnStart(IClient)

	SetOnConnStart(func(IClient))
}
