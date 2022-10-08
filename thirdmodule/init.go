package thirdmodule

import (
	"github.com/cat3306/gameserver/util"
	"time"
)

func Init() {
	util.PanicRepeatRun(InitDb, util.PanicRepeatRunArgs{
		Sleep: time.Second,
		Try:   3,
	})
}
