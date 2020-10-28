package xnet

import (
	"fmt"
	"net"
	"xserver/iface"
)

type Server struct {
	Name        string
	IP          string
	Port        int
	ConnMgr     iface.IConnManager
	MsgHandle   iface.IMsgHandle
	OnConnStart func(conn iface.IClient)
}

func NewServer() iface.IServer {
	return &Server{
		Name:      "node1",
		IP:        "localhost",
		Port:      4040,
		ConnMgr:   NewConnManager(),
		MsgHandle: NewMsgHandler(),
	}
}

func (s *Server) Serve() {
	s.Start()

	select {}
}

func (s *Server) Start() {
	fmt.Printf("[START] Server ,listenner at IP: %s, Port %d is starting\n", s.IP, s.Port)
	go func() {
		s.MsgHandle.StartWorkerPool()

		addr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("ResolveTCPAddr error = ", err)
			return
		}

		ls, err := net.ListenTCP("tcp", addr)
		if err != nil {
			fmt.Println("Listen error = ", err)
			return
		}

		for {
			conn, err := ls.AcceptTCP()
			if err != nil {
				fmt.Println("Accept error = ", err)
				continue
			}

			c := NewClient(s, conn, s.MsgHandle)
			go c.Start()
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("[Stop] server, name ", s.Name)
	s.ConnMgr.Stop()
}

func (s *Server) GetConnMgr() iface.IConnManager {
	return s.ConnMgr
}

func (s *Server) AddRouter(msgID uint32, handler iface.IRouter) {
	s.MsgHandle.AddRouter(msgID, handler)
}

func (s *Server) CallOnConnStart(c iface.IClient) {
	if s.OnConnStart != nil {
		s.OnConnStart(c)
	}
}

func (s *Server) SetOnConnStart(fn func(iface.IClient)) {
	s.OnConnStart = fn
}
