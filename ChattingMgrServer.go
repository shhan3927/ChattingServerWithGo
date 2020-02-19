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
}

func (c *ChattingMgr) RegisterUser(session *network.Session) {
	user := NewChattingUser(c.userSeqNum, session)
	c.userSeqNum++
	c.users[c.userSeqNum] = user
	c.userIdMap[user] = c.userSeqNum
	c.userSessionMap[session] = user
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

		m := &network.Message{
			CmdType: typeValue,
			Body:    payload,
		}

		fmt.Println("Client Recv Message : ", u.nickname)
		c.networkMgr.SendMessage(user.session, m, uint32(response.XXX_Size()))
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
