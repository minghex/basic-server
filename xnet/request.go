package xnet

import "xserve/iface"

type Request struct {
	msg  iface.IMessage
	conn iface.IClient
}

func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgId()
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetClient() iface.IClient {
	return r.conn
}
