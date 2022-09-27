package thirdmodule

import (
	"github.com/cat3306/gameserver/util"
	"time"
)

func Init() {
	util.PanicRepeatRun(InitDb, "InitDb", time.Second)
}
