package router

import (
	"fmt"
	"github.com/cat3306/gameserver/protocol"
)

type HeartBeat struct {
	BaseRouter
}

func (h *HeartBeat) Init() IGameObject {
	return h
}
func (h *HeartBeat) HeartBeat(ctx *protocol.Context) {
	fmt.Println(string(ctx.Payload))
	err := ctx.Conn.AsyncWrite(protocol.Encode("haha", protocol.String, ctx.Proto), nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}
