package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/golang/protobuf/proto"
	"github.com/shhan3927/ChattingServerWithGo/protomessage"
)

///////////////////////////////////
type User struct {
	userId uint32
	name   string
	conn   net.Conn
}

///////////////////////////////////
type Room struct {
	roomId uint32
	name   string
	users  []User
}

///////////////////////////////////
type ChattingManager struct {
	users   map[uint32]User
	rooms   map[uint32]Room
	connMap map[net.Conn]uint32
}

func (this *ChattingManager) Init() {
	this.users = make(map[uint32]User)
	this.rooms = make(map[uint32]Room)
	this.connMap = make(map[net.Conn]uint32)
}

func (this *ChattingManager) AddUser(conn net.Conn) {
	var user User
	user.userId = uint32(len(this.users) + 1)
	user.conn = conn
	this.users[user.userId] = user
	this.connMap[conn] = user.userId
}

func (this *ChattingManager) AddRoom(name string) {
	var room Room
	room.roomId = uint32(len(this.rooms) + 1)
	room.name = name
	this.rooms[room.roomId] = room
}

///////////////////////////////////
func HandleMessage(conn net.Conn) {
	recvBuf := make([]byte, protomessage.MESSAGE_MAX_SIZE)
	for {
		n, err := conn.Read(recvBuf)
		if nil != err {
			if io.EOF == err {
				log.Println(err)
				return
			}
			log.Println(err)
			return
		}

		if n > 0 {
			ParseHeader(recvBuf)
		}
	}
}

func ParseHeader(msg []byte) {
	headerBuffer := msg[:HEADER_SIZE]
	var header Header
	e := header.Unmarshal(headerBuffer)

	if e != NoError {
		log.Println(e)
	}

	switch header.messageType {
	case uint32(protomessage.MessageType_value["kCreateNicknameRequest"]):
		HandleCreateNickName(msg[HEADER_SIZE : HEADER_SIZE+header.bodyLength])
	default:
	}
}

func HandleCreateNickName(msg []byte) {
	var packet protomessage.CreateNicknameRequest

	e := proto.Unmarshal(msg, &packet)
	if e != nil {
		log.Println(e)
	}

	// response...
	fmt.Println(packet.Name)
}

func main() {
	var chattingMgr ChattingManager
	chattingMgr.Init()

	l, err := net.Listen("tcp", ":4321")
	if err != nil {
		log.Println(err)
	}
	defer l.Close()

	var room Room
	room.name = "ARoom"
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		defer conn.Close()

		chattingMgr.AddUser(conn)
		go HandleMessage(conn)
	}
}
