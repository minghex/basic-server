package xnet

import (
	"fmt"
	"sync"
	"xserve/iface"
)

type ConnManager struct {
	connections map[uint32]iface.IClient
	mu          sync.RWMutex
}

func NewConnManager() iface.IConnManager {
	return &ConnManager{
		connections: make(map[uint32]iface.IClient),
	}
}

func (cm *ConnManager) Add(c iface.IClient) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.connections[c.GetConnID()] = c
}

func (cm *ConnManager) Remove(c iface.IClient) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	delete(cm.connections, c.GetConnID())
}

func (cm *ConnManager) Stop() {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	for connID, conn := range cm.connections {
		conn.Stop()
		delete(cm.connections, connID)
	}

	fmt.Println("clear All conncetions successfully: conn num = ", cm.Len())
}

func (cm *ConnManager) Len() int {
	return len(cm.connections)
}
