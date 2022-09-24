package router

import (
	"github.com/cat3306/gameserver/aoi"
	"github.com/panjf2000/gnet/v2"
)

type GameObject struct {
	Id       string
	Position *Vector3
	Rotation *Vector3
	aoi      aoi.AOI
	conn     gnet.Conn
}

func (g *GameObject) Init() IRouter {
	return g
}
