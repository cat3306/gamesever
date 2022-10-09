package conf

import (
	"encoding/json"
	"flag"
	"io/ioutil"
)

type GlobalConf struct {
	Ip              string      `json:"host"`     //ip
	Port            int         `json:"tcp_port"` //port
	MaxConn         int         `json:"max_conn"` //最大连接数
	ConnWriteBuffer int         `json:"conn_write_buffer"`
	ConnReadBuffer  int         `json:"conn_read_buffer"`
	Mysql           MysqlConfig `json:"mysql"`
	AuthConfig      AuthConfig  `json:"auth_config"`
}
type AuthConfig struct {
	IsAuth         bool   `json:"is_auth"` //是否客户端验签
	PublicKeyPath  string `json:"public_key_path"`
	PrivateKeyPath string `json:"private_key_path"`
}
type MysqlConfig struct {
	Host                 string `json:"host"`
	Port                 int    `json:"port"`
	User                 string `json:"user"`
	Pwd                  string `json:"pwd"`
	MysqlConn            string `json:"mysql_conn"`
	MysqlConnectPoolSize int    `json:"mysql_connect_pool_size"`
	SetLog               bool   `json:"set_log"`
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
	GameConfig = gameConfig
}
func LoadJsonConfigLocal(file string, v interface{}) error {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}
