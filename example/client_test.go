package example

import (
	"fmt"
	"github.com/cat3306/gameserver/glog"
	"github.com/cat3306/gameserver/protocol"
	"github.com/cat3306/gameserver/util"
	"github.com/cat3306/gocommon/cryptoutil"
	"io/ioutil"
	"net"
	"testing"
	"time"
)

func Conn() net.Conn {
	//conf.DefaultConf()
	conn, err := net.Dial("tcp", "127.0.0.1:8840")
	if err != nil {
		fmt.Println(err)
		//os.Exit(0)
	}
	return conn
}
func init() {
	glog.Init()
}
func receive(conn net.Conn) {
	go func() {
		for {
			payload, _, _, err := protocol.ReadFull(conn)
			if err != nil {
				panic(err)
			}
			glog.Logger.Sugar().Infof(string(payload))
		}
	}()
}
func TestHeartBeat(t *testing.T) {
	conn := Conn()
	auth(conn)
	receive(conn)
	heartBeat(conn, false)
}
func TestGoHeartBeat(t *testing.T) {
	conn := Conn()
	receive(conn)
	heartBeat(conn, true)
}
func TestHearBeatMore(t *testing.T) {
	for i := 0; i < 100; i++ {
		go func() {
			conn := Conn()
			if conn == nil {
				return
			}
			auth(conn)
			receive(conn)
			heartBeat(conn, false)
		}()
		time.Sleep(time.Millisecond * 20)
	}
	select {}
}
func heartBeat(conn net.Conn, isGo bool) {
	m := "HeartBeat"
	if isGo {
		m = "GoHeartBeat"
	}
	raw, msgLen := protocol.Encode("💓", protocol.String, util.MethodHash(m))

	for {
		_, err := conn.Write(raw[:msgLen])
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}
		time.Sleep(1000 * time.Millisecond)
	}
}

func TestCreateRoom(t *testing.T) {
	conn := Conn()
	receive(conn)
	createRoom(conn)
}

func createRoom(conn net.Conn) {

	type CreateRoomReq struct {
		Pwd       string `json:"Pwd"`
		MaxNum    int    `json:"MaxNum"`    //最大人数
		JoinState bool   `json:"JoinState"` //是否能加入
	}
	req := CreateRoomReq{
		Pwd:       "123456",
		MaxNum:    10,
		JoinState: true,
	}
	raw, msgLen := protocol.Encode(req, protocol.Json, util.MethodHash("CreateRoom"))
	_, err := conn.Write(raw[:msgLen])
	if err != nil {
		fmt.Println("write error err ", err)
		return
	}
	select {}
}

func TestJoinRoom(t *testing.T) {
	conn := Conn()
	receive(conn)
	joinRoom(conn)
}
func joinRoom(conn net.Conn) {
	type CreateRoomReq struct {
		Pwd    string `json:"Pwd"`
		RoomId string `json:"RoomId"`
	}
	req := CreateRoomReq{
		Pwd:    "123456",
		RoomId: "kInXQNE",
	}
	raw, msgLen := protocol.Encode(req, protocol.Json, util.MethodHash("JoinRoom"))
	_, err := conn.Write(raw[:msgLen])
	if err != nil {
		fmt.Println("write error err ", err)
		return
	}
	select {}
}

func TestBinFile(t *testing.T) {
	conn := Conn()
	receive(conn)
	binFile(conn)
}
func binFile(conn net.Conn) {
	bin, err := ioutil.ReadFile("/Users/joker/Downloads/GitHubDesktop-x64.zip")
	if err != nil {
		fmt.Println(err.Error())
	}
	raw, msgLen := protocol.EncodeBin(bin, protocol.String, util.MethodHash("TestBinFile"))
	_, err = conn.Write(raw[:msgLen])
	fmt.Println(err)
}

func TestAuth(t *testing.T) {
	conn := Conn()
	receive(conn)
	auth(conn)
}
func auth(conn net.Conn) {
	var req struct {
		CipherText []byte `json:"CipherText"`
		Text       string `json:"Text"`
	}
	req.Text = "life is short"
	pubKey, err := cryptoutil.RawRSAKey("./public_key.pem")
	if err != nil {
		fmt.Println(err)
		return
	}
	req.CipherText = cryptoutil.RsaEncrypt([]byte(req.Text), pubKey)
	raw, msgLen := protocol.Encode(req, protocol.Json, util.MethodHash("ClientAuth"))
	_, err = conn.Write(raw[:msgLen])
	fmt.Println(err)
}
