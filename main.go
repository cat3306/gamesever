package main

import (
	"github.com/cat3306/gameserver/conf"
	"github.com/cat3306/gameserver/engine"
	"github.com/cat3306/gameserver/router"
)

func main() {
	conf.Init()
	e := engine.NewEngine()
	e.AddRouter(new(router.HeartBeat))
	e.Run()
}
