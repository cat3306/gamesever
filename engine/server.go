package engine

import (
	"fmt"
	"github.com/cat3306/gameserver/conf"
	"github.com/cat3306/gameserver/glog"
	"github.com/cat3306/gameserver/protocol"
	"github.com/cat3306/gameserver/router"
	"github.com/cat3306/gameserver/util"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/pool/goroutine"
	"time"
)

type Server struct {
	gPool   *goroutine.Pool
	connMgr *protocol.ConnManager
	gnet.BuiltinEventEngine
	eng        gnet.Engine
	handlerMgr *HandlerManager
}

func NewEngine() *Server {

	return &Server{
		connMgr:    protocol.NewConnManager(),
		gPool:      goroutine.Default(),
		handlerMgr: NewHandlerManager(),
	}
}
func (s *Server) OnBoot(e gnet.Engine) (action gnet.Action) {
	s.eng = e
	glog.Logger.Sugar().Infof("game Server is listening on:%d", conf.GameConfig.Port)
	return
}
func (s *Server) OnTraffic(c gnet.Conn) gnet.Action {
	context, err := protocol.Decode(c)
	if err != nil {
		glog.Logger.Sugar().Errorf("OnTraffic err:%s", err.Error())
		return gnet.None
	}
	if context == nil {
		panic("context nil")
	}
	context.SetConnMgr(s.connMgr)
	s.exeHandler(context)
	return gnet.None
}
func (s *Server) exeHandler(ctx *protocol.Context) {
	f := s.handlerMgr.GetHandler(ctx.Proto)
	if f != nil {
		f(ctx)
		return
	}
	gf := s.handlerMgr.GetGoHandler(ctx.Proto)
	if gf != nil {
		err := s.gPool.Submit(func() {
			gf(ctx, struct{}{})
		})
		if err != nil {
			glog.Logger.Sugar().Errorf("OnTraffic err:%s", err.Error())
		}
	} else {
		glog.Logger.Sugar().Errorf("OnTraffic not found hander,proto:%d", ctx.Proto)
	}
}

func (s *Server) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	glog.Logger.Sugar().Infof("cid:%s close,reason:%s", c.ID(), err.Error())
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
	err := c.SetReadBuffer(conf.GameConfig.ConnReadBuffer)
	if err != nil {
		panic(err)
	}
	err = c.SetWriteBuffer(conf.GameConfig.ConnWriteBuffer)
	if err != nil {
		panic(err)
	}
	glog.Logger.Sugar().Infof("cid:%s connect", c.ID())
	return nil, gnet.None
}
func (s *Server) OnShutdown(e gnet.Engine) {

}
func (s *Server) OnTick() (delay time.Duration, action gnet.Action) {
	delay = 2 * time.Millisecond
	glog.Logger.Info("a")
	return
}
func (s *Server) Run() {
	addr := fmt.Sprintf("tcp://:%d", conf.GameConfig.Port)
	//protocol.InitBufferPool()
	err := gnet.Run(s, addr, gnet.WithMulticore(true))
	panic(err)
}

func (s *Server) AddRouter(routers ...router.IRouter) {
	for _, v := range routers {
		s.handlerMgr.RegisterRouter(v.Init())
	}
}
func (s *Server) AddHandler(method string, f func(c *protocol.Context)) {
	s.handlerMgr.Register(util.MethodHash(method), f)
}
