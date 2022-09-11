package router

import (
	"github.com/cat3306/gameserver/glog"
	"github.com/cat3306/gameserver/protocol"
	"github.com/panjf2000/gnet/v2"
	"sync"
)

type ConnManager struct {
	connections map[string]gnet.Conn
	locker      sync.RWMutex
}

func newConnManager() *ConnManager {
	return &ConnManager{
		connections: map[string]gnet.Conn{},
	}
}
func (c *ConnManager) Len() int {
	c.locker.Lock()
	defer c.locker.Unlock()
	return len(c.connections)
}
func (c *ConnManager) Add(conn gnet.Conn) {
	c.locker.Lock()
	defer c.locker.Unlock()
	c.connections[conn.ID()] = conn
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

type Room struct {
	maxNum    int    //人数
	pwd       string //密码
	joinState bool   //是否能加入
	gameState bool   //游戏状态
	scene     int    //游戏场景
	Id        string
	connMgr   *ConnManager
}

func (r *Room) Broadcast(v interface{}, ctx *protocol.Context) {
	raw := protocol.Encode(v, ctx.CodeType, ctx.Proto)
	r.connMgr.Broadcast(raw)
}
