package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/golang/protobuf/proto"
	"github.com/shhan3927/ChattingServerWithGo/protomessage"
)

// 일단 전역으로...
var client Client

type Client struct {
	userId uint32
	name   string
	//conn   net.Conn
}

func (this *Client) ReqCreateNickname(conn net.Conn, name string) {
	var nicknameReq protomessage.CreateNicknameRequest
	messageType, typeValue := GetPacketType(nicknameReq)
	nicknameReq.MessageType = messageType
	nicknameReq.Name = name

	head := Header{
		messageType: typeValue,
		bodyLength:  uint32(nicknameReq.XXX_Size()),
	}

	headerBuffer := head.Marshal()
	payloadBuffer, _ := proto.Marshal(&nicknameReq)
	buffer := append(headerBuffer, payloadBuffer...)
	conn.Write(buffer)
}

func HandleSendMessage(conn net.Conn) {
	for {
		//fmt.Scan(&name)
	}
}

func HandleRecvMessage(conn net.Conn) {
	data := make([]byte, 1024)

	for {
		_, err := conn.Read(data)
		if err != nil {
			if io.EOF == err {
				log.Println(err)
				return
			}
			log.Println(err)
			return
		}

		// parsing header...
		ParseHeader(conn, data)
	}
}

func ParseHeader(conn net.Conn, msg []byte) {
	headerBuffer := msg[:HEADER_SIZE]
	var header Header
	e := header.Unmarshal(headerBuffer)

	if e != NoError {
		log.Println(e)
	}

	switch header.messageType {
	case uint32(protomessage.MessageType_value["kCreateNicknameResponse"]):
		payload := msg[HEADER_SIZE : HEADER_SIZE+header.bodyLength]
		var response protomessage.CreateNicknameResponse
		e := proto.Unmarshal(payload, &response)
		if e != nil {
			log.Println(e)
		}
		client.userId = response.UserId
		client.name = response.Name
	default:
		fmt.Println("dddd")
	}
}

func main() {
	conn, err := net.Dial("tcp", ":4321")
	if nil != err {
		log.Println(err)
	}

	defer conn.Close()
	go HandleRecvMessage(conn)
	//var name string

	fmt.Println("Input your name")
	//fmt.Scan(&name)
	client.ReqCreateNickname(conn, "Test")
}
