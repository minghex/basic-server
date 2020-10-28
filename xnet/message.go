package xnet

type Message struct {
	DataLen uint32
	Id      uint32
	Data    []byte
}

func NewMsgPackage(id uint32, data []byte) *Message {
	return &Message{
		Id:      id,
		Data:    data,
		DataLen: uint32(len(data)),
	}
}

func (m *Message) GetDataLen() uint32 {
	return m.DataLen
}

func (m *Message) GetMsgId() uint32 {
	return m.Id
}

func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) SetData(data []byte) {
	m.Data = data
}

func (m *Message) SetMsgId(id uint32) {
	m.Id = id
}

func (m *Message) SetDataLen(len uint32) {
	m.DataLen = len
}
