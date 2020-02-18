package main

import (
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/shhan3927/ChattingServerWithGo/network"
	"github.com/shhan3927/ChattingServerWithGo/protomessage"
)

type ITCPSocket interface {
	Recv([]byte)
	Send([]byte)
}

func NewChattingMgr() *ChattingMgr {
	chattingMgr := &ChattingMgr{
		users:          make(map[uint32]*ChattingUser),
		userIdMap:      make(map[*ChattingUser]uint32),
		userSessionMap: make(map[*network.Session]*ChattingUser),
	}
	return chattingMgr
}

type ChattingMgr struct {
	users          map[uint32]*ChattingUser
	userIdMap      map[*ChattingUser]uint32
	userSessionMap map[*network.Session]*ChattingUser
	userSeqNum     uint32
	networkMgr     *network.TCPServer
}

func (c *ChattingMgr) Init() {
	c.networkMgr = network.NewTCPServer()
	c.networkMgr.OnConnect = c.RegisterUser
	c.networkMgr.OnRecvMessage = c.dispatchMessage
	c.networkMgr.Start(":4300")
	//.Start()
}

func (c *ChattingMgr) Start() {
	// for {
	// 	select {
	// 	case session := <-c.networkMgr.ConnectCh:
	// 		c.RegisterUser(session)
	// 	case message := <-c.networkMgr.MessageCh:
	// 		c.dispatchMessage(message)
	// 	}
	// }
}

func (c *ChattingMgr) RegisterUser(session *network.Session) {
	user := NewChattingUser(c.userSeqNum, session)
	c.userSeqNum++
	c.users[c.userSeqNum] = user
	c.userIdMap[user] = c.userSeqNum
	c.userSessionMap[session] = user

	// go func(u *ChattingUser) {
	// 	for {
	// 		select {
	// 		case message := <-u.client.Data:
	// 			fmt.Println("Client Recv Message...")
	// 			c.ParseMessage(user, message)
	// 		}
	// 	}
	// }(user)
}

func (c *ChattingMgr) ParseMessage(session *network.Session, message []byte) {
	// header, payload, err := GetHeadAndPayload(message)
	// if err != nil {
	// 	return
	// }

	// switch header.messageType {
	// case uint32(protomessage.MessageType_value["kCreateNicknameRequest"]):
	// 	c.HandleCreateNickName(nil, payload)
	// default:
	// }
}

func (c *ChattingMgr) dispatchMessage(session *network.Session, msg *network.Message) {
	switch msg.CmdType {
	case uint32(protomessage.MessageType_kCreateNicknameRequest):
		c.HandleCreateNickName(c.userSessionMap[session], msg.Body)
	}
}

func (c *ChattingMgr) HandleCreateNickName(user *ChattingUser, msg []byte) {
	var packet protomessage.CreateNicknameRequest

	e := proto.Unmarshal(msg, &packet)
	if e != nil {
		log.Println(e)
	}

	fmt.Println(packet.Name)

	// 나중에 수정
	// response...
	u := c.ModifyUserNickname(user, packet.Name)
	if u != nil {
		var response protomessage.CreateNicknameResponse
		messageType, typeValue := network.GetPacketType(response)
		response.MessageType = messageType
		response.UserId = u.userId
		response.Name = u.nickname
		payload, _ := proto.Marshal(&response)
		fmt.Println("Client Recv Message : ", u.nickname)
		c.SendToClient(typeValue, uint32(response.XXX_Size()), payload, u)
	}
}

// 나중에 수정
func (c *ChattingMgr) SendToClient(packetType uint32, bodySize uint32, payload []byte, user *ChattingUser) {
	head := network.Header{
		MessageType: packetType,
		BodyLength:  bodySize,
	}

	headerBuffer := head.Marshal()
	buffer := append(headerBuffer, payload...)

	_, err := user.session.Socket.Write(buffer)
	if err != nil {
		log.Println(err)
		return
	}
}

func (c *ChattingMgr) ModifyUserNickname(user *ChattingUser, nickname string) *ChattingUser {
	userId := c.userIdMap[user]
	u, isExist := c.users[userId]

	if !isExist {
		return nil
	}

	u.nickname = nickname
	return u
}

/////////////////////////////////////////

/////////////////////////////////////

func NewChattingUser(userId uint32, session *network.Session) *ChattingUser {
	c := &ChattingUser{
		userId:  userId,
		session: session,
	}

	return c
}

type ChattingUser struct {
	userId   uint32
	nickname string
	session  *network.Session
}

func (c *ChattingUser) Recv(data []byte) {

}

func (c *ChattingUser) Send(data []byte) {

}
