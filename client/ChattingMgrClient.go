package client

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/shhan3927/ChattingServerWithGo/common"
	"github.com/shhan3927/ChattingServerWithGo/network"
	"github.com/shhan3927/ChattingServerWithGo/protomessage"
)

var instance *ChattingMgrClient
var once sync.Once

type OnResponseDelegate func()

func GetChattingMgr() *ChattingMgrClient {
	once.Do(func() {
		instance = &ChattingMgrClient{
			roomList: make(map[uint32]*common.Room),
		}
	})
	return instance
}

type ClientNetwork struct {
	socket  net.Conn
	recvBuf []byte
}

type User struct {
	name   string
	userId uint32
}

type ChattingMgrClient struct {
	master     User
	roomList   map[uint32]*common.Room
	joinedRoom *common.Room
	//networkInfo      *ClientNetwork
	OnCreateNickname OnResponseDelegate
	OnCreateRoom     OnResponseDelegate
	networkMgr       *NetworkMgr

	chattingUI *ChattingUI
}

func (c *ChattingMgrClient) Start() {
	c.networkMgr = NewNetworkMgr()
	c.networkMgr.OnRecvMsg = c.OnRecvMsg
	c.chattingUI = NewChattingUI()
	c.networkMgr.start()
	c.chattingUI.Start()

	for {
		time.Sleep(time.Millisecond)
	}
}

func (c *ChattingMgrClient) OnRecvMsg(msg *common.Message) {
	switch msg.CmdType {
	case uint32(protomessage.MessageType_kCreateNicknameResponse):
		var response protomessage.CreateNicknameResponse
		e := proto.Unmarshal(msg.Body, &response)
		if e != nil {
			log.Println(e)
		}
		c.master.userId = response.UserId
		c.master.name = response.Name
		c.OnCreateNickname()
		fmt.Println("Response : CreateNickname")
	case uint32(protomessage.MessageType_kCreateRoomResponse):
		var response protomessage.CreateRoomResponse
		e := proto.Unmarshal(msg.Body, &response)
		if e != nil {
			log.Println(e)
		}

		// room 처리
		fmt.Println("room : ", response.RoomId, response.Name)
		c.OnCreateRoom()
		fmt.Println("Response : CreateRoom")
	}
}

func (c *ChattingMgrClient) ReqCreateRoom(name string) {
	fmt.Println("ReqCreateRoom")
	var req protomessage.CreateRoomRequest
	messageType, typeValue := network.GetPacketType(req)
	req.MessageType = messageType
	req.Name = name

	head := network.Header{
		MessageType: typeValue,
		BodyLength:  uint32(req.XXX_Size()),
	}

	headerBuffer := head.Marshal()
	payloadBuffer, _ := proto.Marshal(&req)
	buffer := append(headerBuffer, payloadBuffer...)
	c.networkMgr.socket.Write(buffer)
}

func (c *ChattingMgrClient) ReqCreateNickname(name string) {
	var nicknameReq protomessage.CreateNicknameRequest
	messageType, typeValue := network.GetPacketType(nicknameReq)
	nicknameReq.MessageType = messageType
	nicknameReq.Name = name

	head := network.Header{
		MessageType: typeValue,
		BodyLength:  uint32(nicknameReq.XXX_Size()),
	}

	headerBuffer := head.Marshal()
	payloadBuffer, _ := proto.Marshal(&nicknameReq)
	buffer := append(headerBuffer, payloadBuffer...)
	c.networkMgr.socket.Write(buffer)
}
