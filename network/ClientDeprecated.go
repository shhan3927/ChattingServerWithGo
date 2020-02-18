package network

import (
	"fmt"
	_ "fmt"
	"log"
	"net"

	"github.com/golang/protobuf/proto"
	"github.com/shhan3927/ChattingServerWithGo/protomessage"
)

///////////////////////////////////
type Client struct {
	userId   uint32
	Data     []byte
	nickname string
	Socket   net.Conn
	roomId   uint32
}

func (client *Client) Receive() {
	for {
		//message := make([]byte, MESSAGE_MAX_SIZE)
		length, err := client.Socket.Read(client.Data)
		if err != nil {
			client.Socket.Close()
			break
		}
		if length > 0 {
			client.ParseHeader()
		}
	}
}

func (this *Client) ReqCreateNickname(name string) {
	var nicknameReq protomessage.CreateNicknameRequest
	messageType, typeValue := GetPacketType(nicknameReq)
	nicknameReq.MessageType = messageType
	nicknameReq.Name = name

	head := Header{
		MessageType: typeValue,
		BodyLength:  uint32(nicknameReq.XXX_Size()),
	}

	headerBuffer := head.Marshal()
	payloadBuffer, _ := proto.Marshal(&nicknameReq)
	buffer := append(headerBuffer, payloadBuffer...)
	this.Socket.Write(buffer)
}

func (client *Client) ParseHeader() {
	header, payload, err := GetHeadAndPayload2(client)
	if err != nil {
		return
	}

	switch header.MessageType {
	case uint32(protomessage.MessageType_value["kCreateNicknameResponse"]):
		var response protomessage.CreateNicknameResponse
		e := proto.Unmarshal(payload, &response)
		if e != nil {
			log.Println(e)
		}
		client.userId = response.UserId
		client.nickname = response.Name
		fmt.Println(client.nickname)
	default:
		fmt.Println("dddd")
	}
}

// ///////////////////////////////////
// type Room struct {
// 	roomId  uint32
// 	name    string
// 	clients []*Client
// }

// ///////////////////////////////////
// type ChattingServer struct {
// 	userIds    map[uint32]*Client
// 	clients    map[*Client]uint32
// 	broadcast  chan []byte
// 	register   chan *Client
// 	unregister chan *Client
// 	rooms      map[uint32]*Room
// 	roomSeqNum uint32
// 	userSeqNum uint32
// }

// func (this *ChattingServer) Init() {
// 	this.userIds = make(map[uint32]*Client)
// 	this.clients = make(map[*Client]uint32)
// 	this.rooms = make(map[uint32]*Room)
// 	this.register = make(chan *Client)
// 	this.unregister = make(chan *Client)
// 	this.roomSeqNum = 0
// 	this.userSeqNum = 0
// }

// func (this *ChattingServer) Start() {
// 	for {
// 		select {
// 		case client := <-this.register:
// 			this.RegisterClient(client)
// 		case client := <-this.unregister:
// 			this.UnregisterClient(client)
// 			// case message := <-this.broadcast:
// 			// 	for connection := range this.clients {
// 			// 		select {
// 			// 		case connection.data <- message:
// 			// 		default:
// 			// 			close(connection.data)
// 			// 			delete(this.clients, connection)
// 			// 		}
// 			// 	}
// 		}
// 	}
// }

// func (this *ChattingServer) RegisterClient(client *Client) {
// 	if _, found := this.clients[client]; !found {
// 		this.userSeqNum++
// 		client.userId = this.userSeqNum
// 		this.clients[client] = this.userSeqNum
// 		this.userIds[this.userSeqNum] = client
// 		fmt.Printf("Client %d is connected..\n", this.userSeqNum)
// 	}
// }

// func (this *ChattingServer) UnregisterClient(client *Client) {
// 	if userId, found := this.clients[client]; found {
// 		delete(this.clients, client)
// 		delete(this.userIds, userId)
// 	}
// }

// func (this *ChattingServer) GetClient(userId uint32) (*Client, bool) {
// 	if val, ok := this.userIds[userId]; ok {
// 		return val, true
// 	}

// 	return nil, false
// }

// func (this *ChattingServer) ModifyUserNickname(inClient *Client, nickname string) *Client {
// 	userId := this.clients[inClient]
// 	client, isExist := this.GetClient(userId)

// 	if !isExist {
// 		return nil
// 	}

// 	client.nickname = nickname
// 	return client
// }

// func (this *ChattingServer) AddRoom(name string) {
// 	var room Room
// 	room.roomId = uint32(this.roomSeqNum + 1)
// 	room.name = name
// 	this.rooms[room.roomId] = &Room{roomId: uint32(this.roomSeqNum + 1), name: name}
// 	this.roomSeqNum++
// }

// func (this *ChattingServer) GetRoom(roomId uint32) (*Room, bool) {
// 	if val, ok := this.rooms[roomId]; ok {
// 		return val, true
// 	}

// 	return nil, false
// }

// func (this *ChattingServer) HandleMessage(client *Client) {
// 	for {
// 		n, err := client.Socket.Read(client.Data)
// 		if nil != err {
// 			if io.EOF == err {
// 				log.Println(err)
// 				return
// 			}
// 			log.Println(err)
// 			return
// 		}

// 		if n > 0 {
// 			this.ParseMessage(client)
// 		}
// 	}
// }

// func (this *ChattingServer) ParseMessage(client *Client) {
// 	header, payload, err := GetHeadAndPayload2(client)
// 	if err != nil {
// 		return
// 	}

// 	switch header.messageType {
// 	case uint32(protomessage.MessageType_value["kCreateNicknameRequest"]):
// 		this.HandleCreateNickName(client, payload)
// 	default:
// 	}
// }

// func (this *ChattingServer) HandleCreateNickName(inClient *Client, msg []byte) {
// 	var packet protomessage.CreateNicknameRequest

// 	e := proto.Unmarshal(msg, &packet)
// 	if e != nil {
// 		log.Println(e)
// 	}

// 	// response...
// 	client := this.ModifyUserNickname(inClient, packet.Name)
// 	if client != nil {
// 		var response protomessage.CreateNicknameResponse
// 		messageType, typeValue := GetPacketType(response)
// 		response.MessageType = messageType
// 		response.UserId = client.userId
// 		response.Name = client.nickname
// 		payload, _ := proto.Marshal(&response)
// 		this.SendToClient(typeValue, uint32(response.XXX_Size()), payload, client)
// 	}
// }

// func (this *ChattingServer) SendToClient(packetType uint32, bodySize uint32, payload []byte, client *Client) {
// 	head := Header{
// 		messageType: packetType,
// 		bodyLength:  bodySize,
// 	}

// 	headerBuffer := head.Marshal()
// 	buffer := append(headerBuffer, payload...)

// 	_, err := client.Socket.Write(buffer)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// }

// /////////////////////////////////////////
func GetHeadAndPayload2(client *Client) (Header, []byte, error) {
	head := Header{}
	// n, err := client.Socket.Read(client.Data[:HEADER_SIZE])
	// if err != nil {
	// 	if n == 0 || err == io.EOF {
	// 		//	syscall.WSAECONNRESET on windows
	// 		return head, nil, SessionEDisconnect
	// 	}
	// 	return head, nil, err
	// }
	err := head.Unmarshal(client.Data[:HEADER_SIZE])
	if err != nil {
		return head, nil, err
	}
	if head.BodyLength == 0 {
		return head, nil, nil
	}
	// read body
	// _, err = client.Socket.Read(client.Data[:head.bodyLength])
	// if err != nil {
	// 	return head, nil, err
	// }
	return head, client.Data[HEADER_SIZE : HEADER_SIZE+head.BodyLength], nil
}
