package main

import (
	"github.com/cat3306/gameserver/conf"
	"github.com/cat3306/gameserver/engine"
	"github.com/cat3306/gameserver/glog"
	"github.com/cat3306/gameserver/router"
	"github.com/cat3306/gameserver/thirdmodule"
	"math/rand"
	"time"
)

func main() {
	conf.Init()
	glog.Init()
	e := engine.NewEngine()
	rand.Seed(time.Now().UnixNano())
	e.AddRouter(
		new(router.HeartBeat),
		new(router.RoomManager),
		new(router.GameManager),
		new(router.ClientAuth),
		new(router.TestRouter),
	)
	thirdmodule.Init()
	e.Run()
}
