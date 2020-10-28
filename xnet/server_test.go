package xnet

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"
	"xserver/iface"
)

type PingRouter struct {
	BaseRouter
}

func (pr *PingRouter) Handle(req iface.IRequest) {
	fmt.Println("C2S msgID =", req.GetMsgID(), ", data =", string(req.GetData()))

	err := req.GetClient().SendBuffMsg(1, []byte("pong...pong...pong\n"))
	if err != nil {
		fmt.Println("Handle SendMsg err: ", err)
	}
}

func ClientTest() {

	time.Sleep(3 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:4040")
	if err != nil {
		fmt.Println("Dial error = ", err)
		return
	}

	for {
		dp := NewDataPack()

		msg, _ := dp.Pack(NewMsgPackage(1, []byte("c2s ping ping ping")))
		_, err := conn.Write(msg)
		if err != nil {
			fmt.Println("write error=", err)
			return
		}

		headData := make([]byte, dp.GetDataLen())
		_, err = io.ReadFull(conn, headData)
		if err != nil {
			fmt.Println("client read head err :", err)
			return
		}

		msgHead, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("unpack err :", err)
			return
		}

		if msgHead.GetDataLen() > 0 {
			msg := msgHead.(*Message)
			msg.Data = make([]byte, msg.GetDataLen())

			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("client read data err :", err)
				return
			}

			fmt.Printf("client recv :id = %d,len = %d, data = %s\n", msg.Id, msg.DataLen, msg.Data)
		}
	}
}

func TestServer(t *testing.T) {
	s := NewServer()

	s.AddRouter(1, &PingRouter{})

	go ClientTest()

	go s.Serve()

	select {
	case <-time.After(time.Second * 10):
		return
	}
}
