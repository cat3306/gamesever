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
	s.exeHandler(context)
	return gnet.None
}
func (s *Server) onConnAuth(c gnet.Conn) bool {
	ok := c.GetProperty(protocol.Auth)
	if ok == "" {
		ip := c.RemoteAddr()
		f := func(raw []byte, msgLen int) error {
			return c.AsyncWrite(raw[:msgLen], func(c gnet.Conn) error {
				protocol.BUFFERPOOL.Put(raw)
				c.Close()
				return nil
			})
		}
		err := f(protocol.Encode("auth error ", protocol.String, 0))
		glog.Logger.Sugar().Errorf("onConnAuth auth falied err:%v,ip:%s", err, ip.String())
		return false
	}
	return true
}
func (s *Server) executeHandler(ctx *protocol.Context) error {
	// 同步

	f := s.handlerMgr.GetHandler(ctx.Proto)

	if f != nil {
		f(ctx)
		return nil
	}
	// 异步
	gf := s.handlerMgr.GetGoHandler(ctx.Proto)
	if gf != nil {

		err := s.gPool.Submit(func() {
			gf(ctx, struct{}{})
		})
		if err != nil {
			glog.Logger.Sugar().Errorf("executeHandler err:%s", err.Error())
		}
		return nil
	}
	return fmt.Errorf("not found hander,proto:%d", ctx.Proto)
}
func (s *Server) exeHandler(ctx *protocol.Context) {
	if conf.GameConfig.AuthConfig.IsAuth {
		sf := s.handlerMgr.GetSHandler(ctx.Proto)
		if sf != nil {
			sf(ctx, nil)
			return
		}
	}
	if conf.GameConfig.AuthConfig.IsAuth {
		if !s.onConnAuth(ctx.Conn) {
			return
		}
	}
	err := s.executeHandler(ctx)
	if err != nil {
		glog.Logger.Sugar().Errorf("executeHandler err:%s", err.Error())
	}
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
