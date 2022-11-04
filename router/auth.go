package router

import (
	"github.com/cat3306/gameserver/conf"
	"github.com/cat3306/gameserver/glog"
	"github.com/cat3306/gameserver/protocol"
	"github.com/cat3306/gameserver/router/protoc"
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

func (c *ClientAuth) ClientAuth(ctx *protocol.Context, v interface{}) {
	req := protoc.AuthRequest{}
	err := ctx.Bind(&req)
	if err != nil {
		glog.Logger.Sugar().Errorf("param err:%s", err.Error())
		return
	}
	//glog.Logger.Sugar().Errorf("req.CipherText:%s,c.rawPrivateKey:%s", req.CipherText, req.CipherText)
	text := cryptoutil.RsaDecrypt(req.CipherText, c.rawPrivateKey)
	if string(text) != req.Text {
		glog.Logger.Sugar().Errorf("认证失败!")
		ctx.SendWithCodeType(JsonRspErr("auth failed"), protocol.Json)
		time.Sleep(time.Second)
	} else {
		ctx.Conn.SetProperty(protocol.Auth, "ok")
		ctx.SendWithCodeType(JsonRspOK("auth ok"), protocol.Json)
	}
}
