package router

type BaseRouter struct {
}

func (b *BaseRouter) Init() IRouter {
	return b
}
