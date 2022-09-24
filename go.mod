module github.com/cat3306/gameserver

go 1.17

require (
	github.com/panjf2000/gnet/v2 v2.1.1
	github.com/valyala/bytebufferpool v1.0.0
	go.uber.org/zap v1.23.0
)

require (
	github.com/panjf2000/ants/v2 v2.4.8 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/sys v0.0.0-20220224120231-95c6836cb0e7 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
)

replace github.com/panjf2000/gnet/v2 => github.com/cat3306/gnet/v2 v2.1.3
