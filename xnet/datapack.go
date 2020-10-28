package xnet

import (
	"bytes"
	"encoding/binary"
	"xserve/iface"
)

type DataPack struct{}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (dp *DataPack) GetDataLen() uint32 {
	return 8
}

func (dp *DataPack) UnPack(binaryData []byte) (iface.IMessage, error) {
	dbuff := bytes.NewReader(binaryData)

	msg := &Message{}

	if err := binary.Read(dbuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	if err := binary.Read(dbuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	return msg, nil
}

func (dp *DataPack) Pack(m iface.IMessage) ([]byte, error) {
	dbuff := bytes.NewBuffer([]byte{})

	if err := binary.Write(dbuff, binary.LittleEndian, m.GetDataLen()); err != nil {
		return nil, err
	}

	if err := binary.Write(dbuff, binary.LittleEndian, m.GetMsgId()); err != nil {
		return nil, err
	}

	if err := binary.Write(dbuff, binary.LittleEndian, m.GetData()); err != nil {
		return nil, err
	}

	return dbuff.Bytes(), nil
}
