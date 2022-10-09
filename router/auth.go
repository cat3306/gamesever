package router

import (
	"github.com/cat3306/gameserver/conf"
	"github.com/cat3306/gameserver/glog"
	"github.com/cat3306/gameserver/protocol"
	"github.com/cat3306/gocommon/cryptoutil"
	"time"
)

type ClientAuth struct {
	PublicKeyPath  string
	PrivateKeyPath string
	rawPrivateKey  []byte
	BaseRouter
}

func (c *ClientAuth) Init() IRouter {
	c.PrivateKeyPath = conf.GameConfig.AuthConfig.PrivateKeyPath
	c.initCert()
	return c
}
func (c *ClientAuth) initCert() {

	privateKeyRaw, err := cryptoutil.RawRSAKey(c.PrivateKeyPath)
	if err != nil {
		panic(err)
	}
	c.rawPrivateKey = privateKeyRaw
}

var req struct {
	CipherText []byte `json:"CipherText"`
	Text       string `json:"Text"`
}

func (c *ClientAuth) ClientAuth(ctx *protocol.Context, v interface{}) {
	err := ctx.Bind(&req)
	if err != nil {
		glog.Logger.Sugar().Errorf("param err:%s", err.Error())
	}
	//glog.Logger.Sugar().Errorf("req.CipherText:%s,c.rawPrivateKey:%s", req.CipherText, c.rawPrivateKey)
	text := cryptoutil.RsaDecrypt(req.CipherText, c.rawPrivateKey)
	if string(text) != req.Text {
		glog.Logger.Sugar().Errorf("认证失败!")
		ctx.Send(JsonRspErr("认证失败"))
		time.Sleep(time.Second)
	} else {
		ctx.Conn.SetProperty(protocol.Auth, "ok")
		ctx.Send(JsonRspOK("认证成功"))
	}
}
