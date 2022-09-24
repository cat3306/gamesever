package router

import (
	"github.com/cat3306/gameserver/glog"
	"github.com/cat3306/gameserver/protocol"
)

type GameManager struct {
	BaseRouter
}

func (g *GameManager) Init() IRouter {
	return g
}
func (g *GameManager) CreatePlayer(ctx *protocol.Context) {

}
func (g *GameManager) ObjectMove(ctx *protocol.Context) {
	glog.Logger.Sugar().Infof(string(ctx.Payload))
}
