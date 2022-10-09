package protocol

import (
	"github.com/cat3306/gameserver/glog"
	"github.com/cat3306/gameserver/util"
	"github.com/panjf2000/gnet/v2"
)

const (
	userId = "user_id"
	roomId = "room_id"
	Auth   = "auth"
)

type Context struct {
	Payload  []byte
	CodeType CodeType
	Proto    uint32
	Conn     gnet.Conn
	connMgr  *ConnManager
}

func (c *Context) Bind(v interface{}) error {
	return GameCoder(c.CodeType).Unmarshal(c.Payload, v)
}
func (c *Context) Send(v interface{}) {
	err := c.AsyncWrite(Encode(v, c.CodeType, c.Proto))
	if err != nil {
		glog.Logger.Sugar().Errorf("AsyncWrite err:%s", err.Error())
	}
}
func (c *Context) SendWithCodeType(v interface{}, codeType CodeType) {
	err := c.AsyncWrite(Encode(v, c.CodeType, c.Proto))
	if err != nil {
		glog.Logger.Sugar().Errorf("AsyncWrite err:%s", err.Error())
	}
}
func (c *Context) SendWithParams(v interface{}, codeType CodeType, method string) {
	err := c.AsyncWrite(Encode(v, codeType, util.MethodHash(method)))
	if err != nil {
		glog.Logger.Sugar().Errorf("AsyncWrite err:%s", err.Error())
	}
}

func (c *Context) AsyncWrite(raw []byte, msgLen int) error {
	return c.Conn.AsyncWrite(raw[:msgLen], func(c gnet.Conn) error {
		BUFFERPOOL.Put(raw)
		return nil
	})
}
func (c *Context) SetConnMgr(connMgr *ConnManager) {
	c.connMgr = connMgr
}
func (c *Context) GlobalBroadcast(v interface{}, ) {
	if c.connMgr != nil {
		raw, msgLen := Encode(v, c.CodeType, c.Proto)
		c.connMgr.Broadcast(raw[:msgLen])
	}
}
func (c *Context) SetUserId(uid string) {
	c.Conn.SetProperty(userId, uid)
}
func (c *Context) GetUserId() string {
	return c.Conn.GetProperty(userId)
}

func (c *Context) SetRoomId(id string) {
	c.Conn.SetProperty(roomId, id)
}
func (c *Context) GetRoomId() string {
	return c.Conn.GetProperty(roomId)
}
func (c *Context) DelRoomId() {
	c.Conn.SetProperty(roomId, "")
}
