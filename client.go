package main

import (
	"log"
	"net"

	"github.com/golang/protobuf/proto"
	"github.com/shhan3927/ChattingServerWithGo/protomessage"
)

type Client struct {
	userId uint32
	conn   net.Conn
}

func (this *Client) ReqCreateNickname(name string) {
	var nicknameReq protomessage.CreateNicknameRequest
	messageType, typeValue := protomessage.GetPacketType(nicknameReq)
	nicknameReq.MessageType = messageType
	nicknameReq.Name = name

	head := Header{
		messageType: typeValue,
		bodyLength:  uint32(nicknameReq.XXX_Size()),
	}

	headerBuffer := head.Marshal()
	payloadBuffer, _ := proto.Marshal(&nicknameReq)
	buffer := append(headerBuffer, payloadBuffer...)
	this.conn.Write(buffer)
}

func HandleSendMessage(conn net.Conn) {
	for {
		//fmt.Scan(&name)
	}
}

func HandleRecvMessage(conn net.Conn) {
	data := make([]byte, protomessage.MESSAGE_MAX_SIZE)

	for {
		_, err := conn.Read(data)
		if err != nil {
			log.Println(err)
			return
		}

		// parsing header...
	}
}

// func main() {
// 	conn, err := net.Dial("tcp", ":4321")
// 	if nil != err {
// 		log.Println(err)
// 	}

// 	defer conn.Close()

// 	var client Client
// 	client.conn = conn

// 	var name string
// 	fmt.Println("Input your name")
// 	fmt.Scan(&name)
// 	client.ReqCreateNickname(name)

// 	go HandleRecvMessage(client.conn)
// }
