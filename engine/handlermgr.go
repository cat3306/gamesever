package engine

import (
	"context"
	"fmt"
	"github.com/cat3306/gameserver/glog"
	"github.com/cat3306/gameserver/protocol"
	"github.com/cat3306/gameserver/router"
	"github.com/cat3306/gameserver/util"
	"reflect"
)

type Handler func(c *protocol.Context)

func NewHandlerManager() *HandlerManager {
	ctx, f := context.WithCancel(context.Background())
	return &HandlerManager{
		handlers: make(map[uint32]Handler),
		ctx:      ctx,
		cancel:   f,
	}
}

type HandlerCtx struct {
	Ctx *protocol.Context
	f   Handler
}
type HandlerManager struct {
	handlers  map[uint32]Handler
	taskQueue chan *HandlerCtx
	ctx       context.Context
	cancel    context.CancelFunc
}

func (h *HandlerManager) Register(hashCode uint32, handler Handler) {
	if _, ok := h.handlers[hashCode]; ok {
		panic(fmt.Sprintf("Register repeated method:%d", hashCode))
	}
	h.handlers[hashCode] = handler
}
func (h *HandlerManager) RegisterRouter(iG router.IGameObject) {
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
	}
}
func (h *HandlerManager) Cancel() {
	h.cancel()
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
	f, _ := h.handlers[proto]
	return f
}
