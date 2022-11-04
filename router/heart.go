package router

import (
	"github.com/cat3306/gameserver/glog"
	"github.com/cat3306/gameserver/protocol"
	"time"
)

type HeartBeat struct {
	BaseRouter
}

func (h *HeartBeat) Init() IRouter {
	return h
}
func (h *HeartBeat) HeartBeat(ctx *protocol.Context) {
	str := ""
	err := ctx.Bind(&str)
	glog.Logger.Sugar().Infof("HeartBeat:%s", str)
	if err != nil {
		glog.Logger.Sugar().Errorf("Bind err:%s", err)
	}
	ctx.Send("*")
}
func (h *HeartBeat) GoHeartBeat(ctx *protocol.Context, n struct{}) {
	str := ""
	err := ctx.Bind(&str)
	glog.Logger.Sugar().Infof("HeartBeat:%s", str)
	if err != nil {
		glog.Logger.Sugar().Errorf("Bind err:%s", err)
	}
	time.Sleep(time.Second * 10)
	ctx.Send("*")
}

