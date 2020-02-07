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

type ChattingMgrServer struct {
	users      map[uint32]*ChattingUser
	userIdMap  map[*ChattingUser]uint32
	userSeqNum uint32
	networkMgr network.TCPServer
}

func (c *ChattingMgrServer) Init() {
	c.users = make(map[uint32]*ChattingUser)
	c.userIdMap = make(map[*ChattingUser]uint32)

	go c.Start()
	c.networkMgr.Init()
	c.networkMgr.Start(":4321")
}

func (c *ChattingMgrServer) Start() {
	for {
		select {
		case client := <-c.networkMgr.Connect:
			c.RegisterUser(&ChattingUser{client: client})
		}
	}
}

func (c *ChattingMgrServer) RegisterUser(user *ChattingUser) {
	c.userSeqNum++
	c.users[c.userSeqNum] = user
	c.userIdMap[user] = c.userSeqNum
	go func(u *ChattingUser) {
		for {
			select {
			case message := <-u.client.Data:
				fmt.Println("Client Recv Message...")
				c.ParseMessage(u, message)
			}
		}
	}(user)
}

func (c *ChattingMgrServer) ParseMessage(user *ChattingUser, message []byte) {
	header, payload, err := GetHeadAndPayload(message)
	if err != nil {
		return
	}

	switch header.messageType {
	case uint32(protomessage.MessageType_value["kCreateNicknameRequest"]):
		c.HandleCreateNickName(user, payload)
	default:
	}
}

func (c *ChattingMgrServer) HandleCreateNickName(user *ChattingUser, msg []byte) {
	var packet protomessage.CreateNicknameRequest

	e := proto.Unmarshal(msg, &packet)
	if e != nil {
		log.Println(e)
	}

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
func (c *ChattingMgrServer) SendToClient(packetType uint32, bodySize uint32, payload []byte, user *ChattingUser) {
	head := Header{
		messageType: packetType,
		bodyLength:  bodySize,
	}

	headerBuffer := head.Marshal()
	buffer := append(headerBuffer, payload...)

	_, err := user.client.Socket.Write(buffer)
	if err != nil {
		log.Println(err)
		return
	}
}

func (c *ChattingMgrServer) ModifyUserNickname(user *ChattingUser, nickname string) *ChattingUser {
	userId := c.userIdMap[user]
	u, isExist := c.users[userId]

	if !isExist {
		return nil
	}

	u.nickname = nickname
	return u
}

/////////////////////////////////////////
func GetHeadAndPayload(message []byte) (Header, []byte, error) {
	head := Header{}
	// n, err := client.socket.Read(client.data[:HEADER_SIZE])
	// if err != nil {
	// 	if n == 0 || err == io.EOF {
	// 		//	syscall.WSAECONNRESET on windows
	// 		return head, nil, SessionEDisconnect
	// 	}
	// 	return head, nil, err
	// }
	err := head.Unmarshal(message[:HEADER_SIZE])
	if err != nil {
		return head, nil, err
	}
	if head.bodyLength == 0 {
		return head, nil, nil
	}
	// read body
	// _, err = client.socket.Read(client.data[:head.bodyLength])
	// if err != nil {
	// 	return head, nil, err
	// }
	return head, message[HEADER_SIZE : HEADER_SIZE+head.bodyLength], nil
}

/////////////////////////////////////
type ChattingUser struct {
	userId   uint32
	nickname string
	client   *network.TCPClient
}

func (c *ChattingUser) Recv(data []byte) {

}

func (c *ChattingUser) Send(data []byte) {

}
