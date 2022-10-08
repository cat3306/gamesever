package router

import (
	"github.com/cat3306/gameserver/protocol"
)

type ClientAuth struct {
	BaseRouter
}

func (c *ClientAuth) Init() IRouter {
	return c
}
func (c *ClientAuth) ClientAuth(ctx *protocol.Context, v interface{}) {
	ctx.Conn.SetProperty(protocol.Auth, "")
}
