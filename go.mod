module github.com/cat3306/gameserver

go 1.17

require (
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/panjf2000/gnet/v2 v2.1.1
	github.com/valyala/bytebufferpool v1.0.0
	go.uber.org/zap v1.23.0
	gorm.io/driver/mysql v1.3.6
	gorm.io/gorm v1.23.10
)

require (
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.20.2 // indirect
	github.com/panjf2000/ants/v2 v2.4.8 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
)

replace github.com/panjf2000/gnet/v2 => github.com/cat3306/gnet/v2 v2.1.3
