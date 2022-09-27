package thirdmodule

import (
	"github.com/cat3306/gameserver/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var (
	MysqlDb *gorm.DB
)

func InitDb() {
	//conf.GameConfig.Mysql.MysqlConn = fmt.Sprintf("%s:%s@tcp(%s)/?charset=utf8mb4&parseTime=true&loc=Local", user, pwd, host)
	db, err := gorm.Open(mysql.Open(conf.GameConfig.Mysql.MysqlConn))
	if err != nil {
		panic(err)
	}
	sqlDb, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDb.SetMaxOpenConns(10)
	sqlDb.SetMaxIdleConns(5)
	err = sqlDb.Ping()
	if err != nil {
		panic(err)
	}
	if true {
		newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second, // Slow SQL threshold
				LogLevel:                  logger.Info, // Log level
				IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
				Colorful:                  true,        // Disable color
			},
		)
		db.Logger = newLogger
	}
	MysqlDb = db

}
