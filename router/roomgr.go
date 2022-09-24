package router

import (
	"errors"
	"fmt"
	"github.com/cat3306/gameserver/glog"
	"github.com/cat3306/gameserver/protocol"
	"github.com/cat3306/gameserver/util"
	"sync"
)

type RoomManager struct {
	BaseRouter
	rooms  map[string]*Room
	locker sync.RWMutex
}

var RoomMgr *RoomManager

func (r *RoomManager) Init() IRouter {
	r.rooms = make(map[string]*Room)
	RoomMgr = r
	return r
}

type CreateRoomReq struct {
	Pwd       string `json:"Pwd"`
	MaxNum    int    `json:"MaxNum"`    //最大人数
	JoinState bool   `json:"JoinState"` //是否能加入
}
type CreateRoomRsp struct {
	Id string `json:"id"`
}

func (r *RoomManager) CreateRoom(ctx *protocol.Context) {
	req := &CreateRoomReq{}
	err := ctx.Bind(req)
	if err != nil {
		glog.Logger.Sugar().Errorf("ctx.Bind err:%s", err.Error())
		ctx.Send(JsonRspErr(err.Error()))
		return
	}
	if req.MaxNum == 0 {
		req.MaxNum = 1
	}
	room := &Room{
		maxNum:    req.MaxNum,
		pwd:       req.Pwd,
		joinState: req.JoinState,
		gameState: false,
		scene:     0,
		Id:        util.GenId(7),
		connMgr:   protocol.NewConnManager(),
	}
	room.connMgr.Add(ctx.Conn)
	r.AddRoom(room)
	ctx.SetRoomId(room.Id)
	ctx.Send(JsonRspOK(CreateRoomRsp{Id: room.Id}))
}
func (r *RoomManager) AddRoom(room *Room) {
	r.locker.Lock()
	defer r.locker.Unlock()
	r.rooms[room.Id] = room
}
func (r *RoomManager) DelRoom(id string) {
	r.locker.Lock()
	defer r.locker.Unlock()
	delete(r.rooms, id)
}
func (r *RoomManager) GetRoom(id string) (*Room, bool) {
	r.locker.RLock()
	defer r.locker.RUnlock()
	room, ok := r.rooms[id]
	return room, ok
}
func (r *RoomManager) LeaveRoomByConnClose(roomId string, connId string) {
	if roomId == "" {
		return
	}
	room, _ := r.GetRoom(roomId)
	if room == nil {
		return
	}
	room.connMgr.Remove(connId)
	if room.connMgr.Len() == 0 {
		r.DelRoom(roomId)
	}
}
func (r *RoomManager) LeaveRoom(ctx *protocol.Context) {
	roomId := ctx.GetRoomId()
	room, _ := r.GetRoom(roomId)
	if room == nil {
		ctx.SendWithCodeType(JsonRspErr(fmt.Sprintf("room not found,id:%s", roomId)), protocol.Json)
		glog.Logger.Sugar().Errorf("room not found roomId:%s", roomId)
		return
	}

	room.connMgr.Remove(ctx.Conn.ID())
	if room.connMgr.Len() == 0 {
		r.DelRoom(roomId)
	}
	ctx.DelRoomId()
	glog.Logger.Sugar().Infof("ok,roomId:%s,clientId:%s", roomId, ctx.Conn.ID())

	ctx.SendWithCodeType(JsonRspOK(""), protocol.Json)
}

type JoinRoomReq struct {
	RoomId string `json:"RoomId"`
	Pwd    string `json:"Pwd"`
}

func (r *RoomManager) JoinRoom(ctx *protocol.Context) {
	req := &JoinRoomReq{}
	if err := ctx.Bind(req); err != nil {
		glog.Logger.Sugar().Errorf("ctx.Bind err:%s", err.Error())
		ctx.Send(JsonRspErr(err.Error()))
		return
	}
	room, err := r.joinRoom(req, ctx)
	if err != nil {
		glog.Logger.Sugar().Errorf("JoinRoom err:%s,req:%+v", err.Error(), req)
		ctx.Send(JsonRspErr(err.Error()))
		return
	}

	room.Broadcast(JsonRspOK(""), ctx)
	//ctx.Send()
}
func (r *RoomManager) joinRoom(req *JoinRoomReq, ctx *protocol.Context) (*Room, error) {
	room, _ := r.GetRoom(req.RoomId)
	if room == nil {
		return nil, errors.New("room not found")
	}
	if !room.joinState {
		return nil, errors.New("room not allow join")
	}
	if room.pwd != "" {
		if req.Pwd != room.pwd {
			return nil, errors.New("room pwd not correct")
		}
	}
	if room.connMgr.Len() >= room.maxNum {
		return nil, errors.New("room full")
	}
	room.connMgr.Add(ctx.Conn)
	return room, nil
}
