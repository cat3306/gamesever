package util

import (
	"github.com/cat3306/gameserver/glog"
	"time"
)

func runPanicLess(f func(), method string) (panicLess bool) {
	defer func() {
		err := recover()
		panicLess = err == nil
		if err != nil {
			glog.Logger.Sugar().Errorf("%s err:%v", method, err)
		}
	}()
	f()
	return
}

func PanicRepeatRun(f func(), method string, timeSleep time.Duration) {
	for !runPanicLess(f, method) {
		if timeSleep != 0 {
			time.Sleep(timeSleep)
		}
	}
}
