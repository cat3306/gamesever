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
	connMgr *ConnManager
	gnet.BuiltinEventEngine
	eng        gnet.Engine
	handlerMgr *HandlerManager
	rawChan    chan []byte
}

func NewEngine() *Server {

	return &Server{
		connMgr:    newConnManager(),
		gPool:      goroutine.Default(),
		handlerMgr: NewHandlerManager(),
		rawChan:    make(chan []byte, 1000),
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
	context.SetRawChan(s.rawChan)
	f := s.handlerMgr.GetHandler(context.Proto)
	if f != nil {
		err = s.gPool.Submit(func() {
			f(context)
		})
		if err != nil {
			glog.Logger.Sugar().Errorf("OnTraffic err:%s", err.Error())
		}
		//f(context)
	} else {
		glog.Logger.Sugar().Errorf("OnTraffic not found hander,proto:%d", context.Proto)
	}
	return gnet.None
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
	s.connMgr.Add(c)
	glog.Logger.Sugar().Infof("cid:%s connect", c.ID())
	return nil, gnet.None
}
func (s *Server) OnShutdown(e gnet.Engine) {

}
func (s *Server) OnTick() (delay time.Duration, action gnet.Action) {
	delay = 200 * time.Millisecond
	return
}
func (s *Server) Run() {
	addr := fmt.Sprintf("tcp://:%d", conf.GameConfig.Port)
	s.loop()
	err := gnet.Run(s, addr, gnet.WithMulticore(true))
	panic(err)
}
func (s *Server) loop() {
	go func() {
		for {
			select {
			case buf := <-s.rawChan:
				s.connMgr.Broadcast(buf)
			}
		}
	}()
}

func (s *Server) AddRouter(routers ...router.IGameObject) {
	for _, v := range routers {
		s.handlerMgr.RegisterRouter(v.Init())
	}
}
func (s *Server) AddHandler(method string, f func(c *protocol.Context)) {
	s.handlerMgr.Register(util.MethodHash(method), f)
}
