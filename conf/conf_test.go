package conf

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"testing"
)

func TestGen(t *testing.T) {
	g := GlobalConf{
		Ip:              "0.0.0.0",
		Port:            8848,
		MaxConn:         1000,
		ConnWriteBuffer: 1048576,
		ConnReadBuffer:  1048576,
	}
	raw, err := json.Marshal(g)
	if err != nil {
		t.Logf(err.Error())
	}
	t.Log(ioutil.WriteFile("conf.json", raw, fs.ModePerm))
}
