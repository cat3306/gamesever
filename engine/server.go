package engine

import (
	"fmt"
	"github.com/cat3306/gameserver/conf"
	"github.com/cat3306/gameserver/protocol"
	"github.com/cat3306/gameserver/router"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/pool/goroutine"
	"log"
	"time"
)

type Server struct {
	gPool   *goroutine.Pool
	connMgr *connManager
	gnet.BuiltinEventEngine
	eng        gnet.Engine
	handlerMgr *HandlerManager
}

func NewEngine() *Server {

	return &Server{
		connMgr:    newConnManager(),
		gPool:      goroutine.Default(),
		handlerMgr: NewHandlerManager(),
	}
}
func (s *Server) OnBoot(e gnet.Engine) (action gnet.Action) {

	s.eng = e
	log.Printf("game Server is listening on:%d", conf.GameConfig.Port)
	return
}
func (s *Server) OnTraffic(c gnet.Conn) gnet.Action {
	context, err := protocol.Decode(c)
	if err != nil {
		log.Printf("OnTraffic err:%s", err.Error())
		return gnet.None
	}
	f := s.handlerMgr.GetHandler(context.Proto)
	if f != nil {
		err = s.gPool.Submit(func() {
			f(context)
		})
		if err != nil {
			log.Printf("OnTraffic err:%s", err.Error())
		}
	} else {
		log.Printf("OnTraffic not found hander,proto:%d", context.Proto)
	}
	return gnet.None
}

func (s *Server) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	fmt.Println(err.Error())
	return gnet.None
}
func (s *Server) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
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
	err := gnet.Run(s, addr, gnet.WithMulticore(true))
	panic(err)
}

func (s *Server) AddRouter(routers ...router.IGameObject) {
	for _, v := range routers {
		s.handlerMgr.RegisterRouter(v)
	}
}
