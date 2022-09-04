package engine

import (
	"github.com/panjf2000/gnet/v2"
	"sync"
)

type connManager struct {
	connections map[string]gnet.Conn
	locker      sync.Mutex
}

func newConnManager() *connManager {
	return &connManager{
		connections: make(map[string]gnet.Conn),
	}
}
