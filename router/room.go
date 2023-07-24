package router

import (
	"github.com/cat3306/gameserver/protocol"
)

type Room struct {
	maxNum    int    //人数
	pwd       string //密码
	joinState bool   //是否能加入
	gameState bool   //游戏状态
	scene     int    //游戏场景
	Id        string
	connMgr   *protocol.ConnManager
}

func (r *Room) Broadcast(v interface{}, ctx *protocol.Context) {
	raw := protocol.Encode(v, ctx.CodeType, ctx.Proto)
	r.connMgr.Broadcast(raw)
}
