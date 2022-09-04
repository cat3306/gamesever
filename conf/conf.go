package conf

import (
	"encoding/json"
	"flag"
	"io/ioutil"
)

type GlobalConf struct {
	Ip              string `json:"host"`     //ip
	Port            int    `json:"tcp_port"` //port
	MaxConn         int    `json:"max_conn"` //最大连接数
	ConnWriteBuffer int    `json:"conn_write_buffer"`
	ConnReadBuffer  int    `json:"conn_read_buffer"`
}

var GameConfig *GlobalConf

func Init() {
	var file string
	flag.StringVar(&file, "c", "", "use -c to bind conf file")
	flag.Parse()
	gameConfig := new(GlobalConf)
	err := LoadJsonConfigLocal(file, gameConfig)
	if err != nil {
		panic(err)
	}
	GameConfig=gameConfig
}
func LoadJsonConfigLocal(file string, v interface{}) error {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}
