package main

import (
	_ "fmt"
	"io"
	"log"
	"net"

	"github.com/golang/protobuf/proto"
	"github.com/shhan3927/ChattingServerWithGo/protomessage"
)

type ChattingServer struct {
	ChattingMgr ChattingMgr_Server
}

func (this *ChattingServer) Start(addr string) bool {
	this.ChattingMgr.Init()
	l, err := net.Listen("tcp", addr)

	if err != nil {
		log.Println(err)
		return false
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		defer conn.Close()
		this.ChattingMgr.AddUser(conn)
		go this.HandleMessage(conn)
	}
}

func (this *ChattingServer) HandleMessage(conn net.Conn) {
	recvBuf := make([]byte, MESSAGE_MAX_SIZE)
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
			this.ParseHeader(conn, recvBuf)
		}
	}
}

func (this *ChattingServer) ParseHeader(conn net.Conn, msg []byte) {
	headerBuffer := msg[:HEADER_SIZE]
	var header Header
	e := header.Unmarshal(headerBuffer)

	if e != NoError {
		log.Println(e)
	}

	switch header.messageType {
	case uint32(protomessage.MessageType_value["kCreateNicknameRequest"]):
		this.HandleCreateNickName(conn, msg[HEADER_SIZE:HEADER_SIZE+header.bodyLength])
	default:
	}
}

func (this *ChattingServer) HandleCreateNickName(conn net.Conn, msg []byte) {
	var packet protomessage.CreateNicknameRequest

	e := proto.Unmarshal(msg, &packet)
	if e != nil {
		log.Println(e)
	}

	// response...
	user := this.ChattingMgr.ModifyUserNickname(conn, packet.Name)
	if user != nil {
		var response protomessage.CreateNicknameResponse
		messageType, typeValue := GetPacketType(response)
		response.MessageType = messageType
		response.UserId = user.userId
		response.Name = user.name
		payload, _ := proto.Marshal(&response)
		this.SendToClient(typeValue, uint32(response.XXX_Size()), payload, user.conn)
	}
}

func (this *ChattingServer) SendToClient(packetType uint32, bodySize uint32, payload []byte, conn net.Conn) {
	head := Header{
		messageType: packetType,
		bodyLength:  bodySize,
	}

	headerBuffer := head.Marshal()
	buffer := append(headerBuffer, payload...)

	_, err := conn.Write(buffer)
	if err != nil {
		log.Println(err)
		return
	}
}
