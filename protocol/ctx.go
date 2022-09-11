package protocol

import (
	"github.com/cat3306/gameserver/glog"
	"github.com/cat3306/gameserver/util"
	"github.com/panjf2000/gnet/v2"
)

const (
	userId = "user_id"
	roomId = "room_id"
)

type Context struct {
	Payload  []byte
	CodeType CodeType
	Proto    uint32
	Conn     gnet.Conn
	rawChan  chan []byte
}

func (c *Context) Bind(v interface{}) error {
	return GameCoder(c.CodeType).Unmarshal(c.Payload, v)
}
func (c *Context) Send(v interface{}) {
	raw := Encode(v, c.CodeType, c.Proto)
	err := c.Conn.AsyncWrite(raw, nil)
	if err != nil {
		glog.Logger.Sugar().Errorf("AsyncWrite err:%s", err.Error())
	}
}
func (c *Context) SendWithCodeType(v interface{}, codeType CodeType) {
	raw := Encode(v, codeType, c.Proto)
	err := c.Conn.AsyncWrite(raw, nil)
	if err != nil {
		glog.Logger.Sugar().Errorf("AsyncWrite err:%s", err.Error())
	}
}
func (c *Context) SendWithParams(v interface{}, codeType CodeType, method string) {
	raw := Encode(v, codeType, util.MethodHash(method))
	err := c.Conn.AsyncWrite(raw, nil)
	if err != nil {
		glog.Logger.Sugar().Errorf("AsyncWrite err:%s", err.Error())
	}
}

func (c *Context) SetRawChan(rawChan chan []byte) {
	c.rawChan = rawChan
}
func (c *Context) GlobalBroadcast(v interface{}, ) {
	raw := Encode(v, c.CodeType, c.Proto)
	c.rawChan <- raw
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
