package xnet

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"xserve/iface"
)

type Client struct {
	s           iface.IServer
	conn        *net.TCPConn
	connID      uint32
	msgBuffChan chan []byte
	msgHandler  iface.IMsgHandle
	ctx         context.Context
	cancel      context.CancelFunc
	isclosed    bool
	mu          sync.RWMutex
}

func NewClient(s iface.IServer, conn *net.TCPConn, msgHandler iface.IMsgHandle) iface.IClient {
	c := &Client{
		s:           s,
		conn:        conn,
		msgBuffChan: make(chan []byte, 1024),
		msgHandler:  msgHandler,
		isclosed:    false,
	}

	c.s.GetConnMgr().Add(c)
	return c
}

func (c *Client) Start() {
	c.ctx, c.cancel = context.WithCancel(context.Background())

	go c.StartReader()
	go c.StartWriter()

	c.s.CallOnConnStart(c)
}

func (c *Client) StartReader() {
	fmt.Println("Client Start Reader")
	defer fmt.Println("Client Stop Reader")
	defer c.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			dp := NewDataPack()

			hd := make([]byte, dp.GetDataLen())
			if _, err := io.ReadFull(c.conn, hd); err != nil {
				fmt.Println("read msg head error =", err)
				break
			}

			msg, err := dp.UnPack(hd)
			if err != nil {
				fmt.Println("UnPack error =", err)
				break
			}

			var data []byte
			if msg.GetDataLen() > 0 {
				data = make([]byte, msg.GetDataLen())
				if _, err := io.ReadFull(c.conn, data); err != nil {
					fmt.Println("read msg data error =", err)
					break
				}
			}
			msg.SetData(data)

			req := &Request{
				msg:  msg,
				conn: c,
			}

			// c.msgHandler.SendMsgToTaskQueue(req)
			go c.msgHandler.DoMsgHandler(req)
		}
	}
}

func (c *Client) StartWriter() {
	fmt.Println("Client Start Write")
	defer fmt.Println("Client Stop Write")

	for {
		select {
		case data, ok := <-c.msgBuffChan:
			if ok {
				if _, err := c.conn.Write(data); err != nil {
					fmt.Println("Send Data error:", err)
					return
				}
			} else {
				break
			}
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *Client) Stop() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isclosed == true {
		return
	}
	c.isclosed = true

	//关闭连接
	c.conn.Close()
	//关闭writer
	c.cancel()
	//将链接删除
	c.s.GetConnMgr().Remove(c)
	//关闭管道
	close(c.msgBuffChan)
}

func (c *Client) GetConnID() uint32 {
	return c.connID
}

func (c *Client) SendBuffMsg(msgId uint32, data []byte) error {
	c.mu.RLock()
	if c.isclosed == true {
		c.mu.RUnlock()
		return errors.New("connection closed when send buff msg")
	}
	c.mu.RUnlock()

	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("msg pack error")
		return err
	}

	c.msgBuffChan <- msg
	return nil
}
