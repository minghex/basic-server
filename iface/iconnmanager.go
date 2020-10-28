package iface

type IConnManager interface {
	Add(IClient)
	Remove(IClient)
	Stop()
	Len() int
}
