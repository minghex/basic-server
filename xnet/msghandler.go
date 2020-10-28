package xnet

import (
	"fmt"
	"xserver/iface"
)

type MsgHandle struct {
	Apis           map[uint32]iface.IRouter
	TaskQueue      []chan iface.IRequest
	WorkerPoolSize uint32
}

func NewMsgHandler() iface.IMsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]iface.IRouter),
		WorkerPoolSize: 8,
		TaskQueue:      make([]chan iface.IRequest, 8),
	}
}

func (mh *MsgHandle) DoMsgHandler(req iface.IRequest) {
	msgId := req.GetMsgID()

	if handler, ok := mh.Apis[msgId]; ok {
		handler.Handle(req)
	} else {
		fmt.Println("not found handler")
	}
}

func (mh *MsgHandle) AddRouter(msgID uint32, handler iface.IRouter) {
	if _, ok := mh.Apis[msgID]; !ok {
		mh.Apis[msgID] = handler
	}
}

func (mh *MsgHandle) SendMsgToTaskQueue(req iface.IRequest) {
	workID := req.GetClient().GetConnID() % 8
	mh.TaskQueue[workID] <- req
}

func (mh *MsgHandle) StartOneWork(workID int, taskQueue chan iface.IRequest) {

	for {
		select {
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

func (mh *MsgHandle) StartWorkerPool() {
	for i := 0; i < 8; i++ {
		mh.TaskQueue[i] = make(chan iface.IRequest, 1024)
		go mh.StartOneWork(i, mh.TaskQueue[i])
	}
}
