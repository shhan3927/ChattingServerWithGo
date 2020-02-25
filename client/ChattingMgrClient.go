package client

import (
	"fmt"
	"log"
	"net"
	"sync"

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
			isClose:  false,
			message:  make(chan *common.Message, 1000),
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
	message          chan *common.Message
	chattingUI       *ChattingUI
	isClose          bool
}

func (c *ChattingMgrClient) Start() {
	c.networkMgr = NewNetworkMgr()
	c.networkMgr.OnRecvMsg = c.OnRecvMsg
	c.chattingUI = NewChattingUI()
	c.networkMgr.start()
	c.chattingUI.Start()

	c.Process()
}

func (c *ChattingMgrClient) Process() {
	for {
		select {
		case message := <-c.message:
			fmt.Println("OnRecvMsg")
			c.dispatch(message)
		}
	}
}

func (c *ChattingMgrClient) OnRecvMsg(msg *common.Message) {
	c.message <- msg
}

func (c *ChattingMgrClient) dispatch(msg *common.Message) {
	switch msg.CmdType {
	case uint32(protomessage.MessageType_kCreateNicknameResponse):
		var response protomessage.CreateNicknameResponse
		e := proto.Unmarshal(msg.Body, &response)
		if e != nil {
			log.Println(e)
		}
		c.master.userId = response.UserId
		c.master.name = response.Name
		fmt.Println("Response : CreateNickname")
		c.OnCreateNickname()
	case uint32(protomessage.MessageType_kCreateRoomResponse):
		var response protomessage.CreateRoomResponse
		e := proto.Unmarshal(msg.Body, &response)
		if e != nil {
			log.Println(e)
		}

		// room 처리
		fmt.Println("room : ", response.RoomId, response.Name)
		fmt.Println("Response : CreateRoom")
		c.OnCreateRoom()
	}
}

func (c *ChattingMgrClient) ReqCreateRoom(name string) {
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
	fmt.Println("Write Req packet to socket")
	c.networkMgr.socket.Write(buffer)
}
