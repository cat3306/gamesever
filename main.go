package main

import (
	"github.com/cat3306/gameserver/conf"
	"github.com/cat3306/gameserver/engine"
	"github.com/cat3306/gameserver/glog"
	"github.com/cat3306/gameserver/router"
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
	)
	e.Run()
}
