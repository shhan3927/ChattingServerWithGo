package main

import (
	"ChattingServerWithGo/protomessage"
	"fmt"
	"log"
	"net"

	"github.com/golang/protobuf/proto"
)

type Client struct {
	uint32 userId
	conn   net.Conn
}

func (this *Client) ReqCreateNickName(name String) {
	var nicknameReq protomessage.CreateNicknameRequest
	messageType, typeValue := protomessage.GetPacketType(nicknameReq)
	nicknameReq.MessageType = messageType
	nicknameReq.Name = "name"

	head := Header{
		messageType: typeValue,
	}

	headerBuffer := head.Marshal()
	payloadBuffer, err := proto.Marshal()
	buffer := append(headerBuffer, payloadBuffer...)
	this.conn.Write(buffer)
}

func main() {
	conn, err := net.Dial("tcp", ":4321")
	if nil != err {
		log.Println(err)
	}

	defer conn.Close()

	var client Client
	client.conn = conn

	var name string
	fmt.Println("Input your name")
	fmt.Scan(&name)
	client.ReqCreateNickname(name)

	go HandleRecvMessage(client.conn)
}

func HandleSendMessage(conn net.Conn) {
	for {
		//fmt.Scan(&name)
	}
}

func HandleRecvMessage(conn net.Conn) {
	data := make([]byte, MessageBodySizeMax)

	for {
		n, err := conn.Read(data)
		if err != nil {
			log.Println(err)
			return
		}

		// parsing header...
	}
}
