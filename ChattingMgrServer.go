package main

import (
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/shhan3927/ChattingServerWithGo/common"
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
		userSessionMap: make(map[uint64]*ChattingUser),
		rooms:          make(map[uint32]*common.Room),
	}
	return chattingMgr
}

type ChattingMgr struct {
	users          map[uint32]*ChattingUser
	userIdMap      map[*ChattingUser]uint32
	userSessionMap map[uint64]*ChattingUser
	userSeqNum     uint32
	networkMgr     *network.TCPServer
	rooms          map[uint32]*common.Room
	roomSeqNum     uint32
}

func (c *ChattingMgr) Init() {
	c.networkMgr = network.NewTCPServer()
	c.networkMgr.OnConnect = c.RegisterUser
	c.networkMgr.OnRecvMessage = c.dispatchMessage
	c.networkMgr.Start(network.CONNECT_PORT)
}

func (c *ChattingMgr) RegisterUser(sessionInfo network.SessionInfo) {
	c.userSeqNum++
	user := NewChattingUser(c.userSeqNum, sessionInfo)
	c.users[c.userSeqNum] = user
	c.userIdMap[user] = c.userSeqNum
	c.userSessionMap[sessionInfo.SessionId] = user
}

func (c *ChattingMgr) dispatchMessage(sessionInfo network.SessionInfo, msg *common.Message) {
	switch msg.CmdType {
	case uint32(protomessage.MessageType_kCreateNicknameRequest):
		c.HandleCreateNickName(c.userSessionMap[sessionInfo.SessionId], msg.Body)
	case uint32(protomessage.MessageType_kCreateRoomRequest):
		c.HandleCreateRoom(c.userSessionMap[sessionInfo.SessionId].userId, msg.Body)
	}
}

func (c *ChattingMgr) HandleCreateNickName(user *ChattingUser, msg []byte) {
	var packet protomessage.CreateNicknameRequest

	e := proto.Unmarshal(msg, &packet)
	if e != nil {
		log.Println(e)
	}

	fmt.Println(packet.Name)

	u := c.ModifyUserNickname(user, packet.Name)
	if u != nil {
		var response protomessage.CreateNicknameResponse
		messageType, typeValue := network.GetPacketType(response)
		response.MessageType = messageType
		response.UserId = u.userId
		response.Name = u.nickname
		payload, _ := proto.Marshal(&response)

		m := &common.Message{
			CmdType: typeValue,
			Body:    payload,
		}

		fmt.Println("Client Recv Message : ", u.nickname)
		c.networkMgr.SendMessage(user.sessionInfo, m, uint32(response.XXX_Size()))
	}
}

func (c *ChattingMgr) HandleCreateRoom(userId uint32, msg []byte) {
	var packet protomessage.CreateRoomRequest
	e := proto.Unmarshal(msg, &packet)
	if e != nil {
		log.Println(e)
	}

	// Create Room
	c.roomSeqNum++
	room := common.NewRoom(c.roomSeqNum, packet.Name)
	room.AddUser(c.roomSeqNum)
	room.SetMaster(userId)

	// Response
	var response protomessage.CreateRoomResponse
	messageType, typeValue := network.GetPacketType(response)
	response.MessageType = messageType
	response.RoomId = room.Id
	response.Name = room.Name
	payload, _ := proto.Marshal(&response)

	m := &common.Message{
		CmdType: typeValue,
		Body:    payload,
	}

	c.networkMgr.SendMessage(c.users[userId].sessionInfo, m, uint32(response.XXX_Size()))
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

func NewChattingUser(userId uint32, info network.SessionInfo) *ChattingUser {
	c := &ChattingUser{
		userId:      userId,
		sessionInfo: info,
	}

	return c
}

type ChattingUser struct {
	userId      uint32
	nickname    string
	sessionInfo network.SessionInfo
}
