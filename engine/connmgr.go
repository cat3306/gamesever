package engine

import (
	"github.com/cat3306/gameserver/glog"
	"github.com/cat3306/gameserver/util"
	"github.com/panjf2000/gnet/v2"
	"sync"
)

type ConnManager struct {
	connections map[string]gnet.Conn
	locker      sync.RWMutex
}

func newConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[string]gnet.Conn),
	}
}
func (c *ConnManager) Add(conn gnet.Conn) {
	c.locker.Lock()
	defer c.locker.Unlock()
	cId := util.GenId(9)
	conn.SetId(cId)
	c.connections[cId] = conn
}
func (c *ConnManager) Remove(id string) {
	c.locker.Lock()
	defer c.locker.Unlock()
	delete(c.connections, id)
}
func (c *ConnManager) Broadcast(raw []byte) {
	c.locker.RLock()
	defer c.locker.RUnlock()
	for _, v := range c.connections {
		err := v.AsyncWrite(raw, nil)
		if err != nil {
			glog.Logger.Sugar().Errorf("AsyncWrite err:%s", err.Error())
		}
	}
}
