package engine

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/cat3306/gameserver/glog"
	"github.com/cat3306/gameserver/protocol"
	"github.com/cat3306/gameserver/router"
	"github.com/cat3306/gameserver/util"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/pool/goroutine"
)

var (
	ErrHandlerNotFound   = errors.New("handler not found")
	ErrHandlerAuthFailed = errors.New("handler auth failed")
)

type Handler func(c *protocol.Context)
type GoHandler func(c *protocol.Context, none struct{})
type SHandler func(c *protocol.Context, v interface{})

func NewHandlerManager() *HandlerManager {
	return &HandlerManager{
		handlers:        make(map[uint32]Handler),
		goHandler:       make(map[uint32]GoHandler),
		specialHandlers: make(map[uint32]SHandler),
		gPool:           goroutine.Default(),
	}
}

type HandlerManager struct {
	handlers        map[uint32]Handler
	goHandler       map[uint32]GoHandler
	specialHandlers map[uint32]SHandler
	gPool           *goroutine.Pool
}

func (h *HandlerManager) Register(hashCode uint32, handler Handler) {
	if _, ok := h.handlers[hashCode]; ok {
		panic(fmt.Sprintf("Register repeated method:%d", hashCode))
	}
	h.handlers[hashCode] = handler
}
func (h *HandlerManager) GoRegister(hashCode uint32, handler GoHandler) {
	if _, ok := h.goHandler[hashCode]; ok {
		panic(fmt.Sprintf("Register repeated method:%d", hashCode))
	}
	h.goHandler[hashCode] = handler
}
func (h *HandlerManager) SpecialRegister(hashCode uint32, handler SHandler) {
	if _, ok := h.specialHandlers[hashCode]; ok {
		panic(fmt.Sprintf("Register repeated method:%d", hashCode))
	}
	h.specialHandlers[hashCode] = handler
}
func (h *HandlerManager) RegisterRouter(iG router.IRouter) {
	t := reflect.TypeOf(iG)
	tName := t.String()
	vl := reflect.ValueOf(iG)
	for i := 0; i < t.NumMethod(); i++ {
		name := t.Method(i).Name
		v, ok := vl.Method(i).Interface().(func(ctx *protocol.Context))
		if ok {
			if checkoutMethod(name) {
				hashId := util.MethodHash(name)
				h.Register(hashId, v)
				glog.Logger.Sugar().Infof("[%s.%s] hashId:%d", tName, name, hashId)
			}
		}
		v1, ok1 := vl.Method(i).Interface().(func(c *protocol.Context, none struct{}))
		if ok1 {
			if checkoutMethod(name) {
				hashId := util.MethodHash(name)
				h.GoRegister(hashId, v1)
				glog.Logger.Sugar().Infof("[%s.go_%s] hashId:%d", tName, name, hashId)
			}
		}
		v2, ok2 := vl.Method(i).Interface().(func(c *protocol.Context, v interface{}))
		if ok2 {
			if checkoutMethod(name) {
				hashId := util.MethodHash(name)
				h.SpecialRegister(hashId, v2)
				glog.Logger.Sugar().Infof("[%s.special_%s] hashId:%d", tName, name, hashId)
			}
		}
	}
}

//函数签名首字母大写才会被注入
func checkoutMethod(m string) bool {
	if len(m) == 0 {
		return false
	}
	if m[0] >= 'A' && m[0] <= 'W' {
		return true
	}
	return false
}
func (h *HandlerManager) GetHandler(proto uint32) Handler {
	f := h.handlers[proto]
	return f
}
func (h *HandlerManager) GetGoHandler(proto uint32) GoHandler {
	f := h.goHandler[proto]
	return f
}
func (h *HandlerManager) GetSHandler(proto uint32) SHandler {
	f := h.specialHandlers[proto]
	return f
}
func (h *HandlerManager) checkConnAuth(c gnet.Conn) bool {
	ok := c.GetProperty(protocol.Auth)
	if ok == "" {
		ip := c.RemoteAddr()
		f := func(raw []byte) error {
			return c.AsyncWrite(raw, func(c gnet.Conn) error {
				protocol.BUFFERPOOL.Put(raw)
				_ = c.Close()
				return nil
			})
		}
		err := f(protocol.Encode("auth error ", protocol.String, 0))
		glog.Logger.Sugar().Errorf("checkConnAuth auth falied err:%v,ip:%s", err, ip.String())
		return false
	}
	return true
}

//同步handler
func (h *HandlerManager) exeSyncHandler(auth bool, ctx *protocol.Context) error {
	f := h.GetHandler(ctx.Proto)
	if f != nil {
		if auth {
			if !h.checkConnAuth(ctx.Conn) {
				return ErrHandlerAuthFailed
			}
		}
		f(ctx)
		return nil
	}
	return ErrHandlerNotFound
}

//异步handler
func (h *HandlerManager) exeAsyncHandler(auth bool, ctx *protocol.Context) error {
	f := h.GetGoHandler(ctx.Proto)
	if f != nil {
		if auth {
			if !h.checkConnAuth(ctx.Conn) {
				return ErrHandlerAuthFailed
			}
		}

		newBuffer := protocol.BUFFERPOOL.Get(uint32(len(ctx.Payload)))
		copy(*newBuffer, ctx.Payload)
		ctx.Payload = *newBuffer
		err := h.gPool.Submit(func() {
			f(ctx, struct{}{})
		})
		if err != nil {
			glog.Logger.Sugar().Errorf("exeGoHandler err:%s", err.Error())
			return err
		}
		return nil
	}
	return ErrHandlerNotFound
}

func (h *HandlerManager) exeNotNeedAuthHandler(ctx *protocol.Context) error {
	f := h.GetSHandler(ctx.Proto)
	if f != nil {
		f(ctx, nil)
		return nil
	}
	return ErrHandlerNotFound
}
func (h *HandlerManager) ExeHandler(auth bool, ctx *protocol.Context) {
	err := h.exeSyncHandler(auth, ctx)
	if !errors.Is(err, ErrHandlerNotFound) {
		return
	}
	err = h.exeAsyncHandler(auth, ctx)
	if !errors.Is(err, ErrHandlerNotFound) {
		return
	}
	err = h.exeNotNeedAuthHandler(ctx)
	if err != nil {
		glog.Logger.Sugar().Errorf("ExeHandler err:%s,hash:%d", err, ctx.Proto)
	}
}
