package util

import (
	"github.com/cat3306/gameserver/glog"
	"math"
	"reflect"
	"runtime"
	"time"
)

type PanicRepeatRunArgs struct {
	Sleep time.Duration
	Try   int
}

func runPanicLess(f func()) (panicLess bool) {
	defer func() {
		err := recover()
		panicLess = err == nil
		if err != nil {
			name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
			glog.Logger.Sugar().Errorf("%s err:%v", name, err)
		}
	}()
	f()
	return
}

func PanicRepeatRun(f func(), args ...PanicRepeatRunArgs) {
	param := PanicRepeatRunArgs{
		Sleep: 0,
		Try:   math.MaxInt16,
	}
	if len(args) != 0 {
		param = args[0]
	}
	if param.Try == 0 {
		param.Try = math.MaxInt8
	}
	total := param.Try
	for !runPanicLess(f) && param.Try >= 1 {
		if param.Sleep != 0 {
			time.Sleep(param.Sleep)
		}
		param.Try--
	}
	if param.Try == 0 {
		name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
		glog.Logger.Sugar().Errorf("%s:finally failed,total:%d", name, total)
	}
}
