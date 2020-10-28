package iface

type IMsgHandle interface {
	DoMsgHandler(request IRequest)
	AddRouter(uint32, IRouter)
	StartWorkerPool()
	SendMsgToTaskQueue(IRequest)
}
