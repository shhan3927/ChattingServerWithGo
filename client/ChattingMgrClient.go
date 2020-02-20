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
			networkInfo: &ClientNetwork{
				recvBuf: make([]byte, 4096),
			},
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
	master           User
	roomList         map[uint32]*common.Room
	joinedRoom       *common.Room
	networkInfo      *ClientNetwork
	OnCreateNickname OnResponseDelegate
	OnCreateRoom     OnResponseDelegate
}

func (c *ChattingMgrClient) Start() {
	socket, error := net.Dial("tcp", network.CONNECT_PORT)
	if error != nil {
		fmt.Println(error)
	}

	defer socket.Close()

	c.networkInfo.socket = socket
	go c.process()

	chattingUI := NewChattingUI()
	chattingUI.Start()
}

func (c *ChattingMgrClient) process() {
	for {
		length, err := c.networkInfo.socket.Read(c.networkInfo.recvBuf)
		if err != nil {
			c.networkInfo.socket.Close()
			break
		}
		if length > 0 {
			c.ParseHeader()
		}
	}
}

func (c *ChattingMgrClient) ParseHeader() {
	header, payload, err := c.GetHeadAndPayload()
	if err != nil {
		return
	}

	switch header.MessageType {
	case uint32(protomessage.MessageType_kCreateNicknameResponse):
		var response protomessage.CreateNicknameResponse
		e := proto.Unmarshal(payload, &response)
		if e != nil {
			log.Println(e)
		}
		c.master.userId = response.UserId
		c.master.name = response.Name
		c.OnCreateNickname()

	case uint32(protomessage.MessageType_kCreateRoomResponse):
		var response protomessage.CreateRoomResponse
		e := proto.Unmarshal(payload, &response)
		if e != nil {
			log.Println(e)
		}

		// room 처리
		c.OnCreateRoom()
		fmt.Println("room : ", response.RoomId, response.Name)
	default:
		fmt.Println("dddd")
	}
}

func (c *ChattingMgrClient) GetHeadAndPayload() (network.Header, []byte, error) {
	head := network.Header{}
	err := head.Unmarshal(c.networkInfo.recvBuf[:network.HEADER_SIZE])
	if err != nil {
		return head, nil, err
	}
	if head.BodyLength == 0 {
		return head, nil, nil
	}

	return head, c.networkInfo.recvBuf[network.HEADER_SIZE : network.HEADER_SIZE+head.BodyLength], nil
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
	c.networkInfo.socket.Write(buffer)
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
	c.networkInfo.socket.Write(buffer)
}
