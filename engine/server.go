package engine

import (
	"fmt"
	"github.com/cat3306/gameserver/conf"
	"github.com/cat3306/gameserver/glog"
	"github.com/cat3306/gameserver/protocol"
	"github.com/cat3306/gameserver/router"
	"github.com/cat3306/gameserver/util"
	"github.com/panjf2000/gnet/v2"
	"time"
)

type Server struct {
	connMgr *protocol.ConnManager
	gnet.BuiltinEventEngine
	eng        gnet.Engine
	handlerMgr *HandlerManager
}

func NewEngine() *Server {
	return &Server{
		connMgr:    protocol.NewConnManager(),
		handlerMgr: NewHandlerManager(),
	}
}
func (s *Server) OnBoot(e gnet.Engine) (action gnet.Action) {
	s.eng = e
	glog.Logger.Sugar().Infof("game Server is listening on:%d", conf.GameConfig.Port)
	return
}
func (s *Server) OnTraffic(c gnet.Conn) gnet.Action {
	defer func() {
		err := recover()
		if err != nil {
			glog.Logger.Sugar().Errorf("OnTraffic panic %v", err)
		}
	}()
	s.eng.CountConnections()
	context, err := protocol.Decode(c)
	if err != nil {
		glog.Logger.Sugar().Warnf("OnTraffic err:%s", err.Error())
		return gnet.None
	}
	if context == nil {
		panic("context nil")
	}
	context.SetConnMgr(s.connMgr)
	s.handlerMgr.ExeHandler(conf.GameConfig.AuthConfig.IsAuth, context)
	return gnet.None
}
func (s *Server) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	reason := ""
	if err != nil {
		reason = err.Error()
	}
	glog.Logger.Sugar().Infof("cid:%s close,reason:%s", c.ID(), reason)
	s.connMgr.Remove(c.ID())
	ctx := protocol.Context{
		Conn: c,
	}
	roomId := ctx.GetRoomId()
	if roomId != "" && router.RoomMgr != nil {
		router.RoomMgr.LeaveRoomByConnClose(roomId, c.ID())
	}

	return gnet.None
}
func (s *Server) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	cId := util.GenId(9)
	c.SetId(cId)
	s.connMgr.Add(c)
	glog.Logger.Sugar().Infof("cid:%s connect", c.ID())
	return nil, gnet.None
}
func (s *Server) OnShutdown(e gnet.Engine) {

}
func (s *Server) Run() {
	addr := fmt.Sprintf("tcp://:%d", conf.GameConfig.Port)
	f := func() {
		err := gnet.Run(s, addr,
			gnet.WithMulticore(true),
			gnet.WithSocketSendBuffer(conf.GameConfig.ConnWriteBuffer),
			gnet.WithSocketRecvBuffer(conf.GameConfig.ConnWriteBuffer),
		)
		panic(err)
	}
	defer func() {
		s.handlerMgr.gPool.Release()
	}()
	util.PanicRepeatRun(f, util.PanicRepeatRunArgs{
		Sleep: time.Second,
		Try:   20,
	})
}

func (s *Server) AddRouter(routers ...router.IRouter) {
	for _, v := range routers {
		s.handlerMgr.RegisterRouter(v.Init())
	}
}
func (s *Server) AddHandler(method string, f func(c *protocol.Context)) {
	s.handlerMgr.Register(util.MethodHash(method), f)
}
