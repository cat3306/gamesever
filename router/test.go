package router

import (
	"github.com/cat3306/gameserver/glog"
	"github.com/cat3306/gameserver/protocol"
	"github.com/cat3306/gameserver/router/protoc"
	"io/ioutil"
	"os/user"
)

type TestRouter struct {
	BaseRouter
}

func (t *TestRouter) Init() IRouter {
	return t
}

func (t *TestRouter) TestProtoBuffer(ctx *protocol.Context) {
	p := &protoc.Position{}
	err := ctx.Bind(p)
	if err != nil {
		glog.Logger.Sugar().Errorf("err:%s", err.Error())
		return
	}
	glog.Logger.Sugar().Infof("p:%+v", p)
}
func (t *TestRouter) TestBinFile(ctx *protocol.Context, n struct{}) {
	u, err := user.Current()
	if err != nil {
		glog.Logger.Sugar().Errorf(" Current err:%s", err.Error())
		return
	}
	err = ioutil.WriteFile(u.HomeDir+"/test.bin", ctx.Payload, 7777)
	if err != nil {
		glog.Logger.Sugar().Errorf(" ioutil.WriteFile err:%s", err.Error())
	}
}
