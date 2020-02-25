package client

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/shhan3927/ChattingServerWithGo/common"
	"github.com/shhan3927/ChattingServerWithGo/network"
)

type NetworkMgr struct {
	socket    net.Conn
	recvBuf   []byte
	sendBufCh chan []byte
	OnRecvMsg MessageDelegate
	wg        sync.WaitGroup
}

type MessageDelegate func(*common.Message)

func NewNetworkMgr() *NetworkMgr {
	return &NetworkMgr{
		recvBuf: make([]byte, 4096),
	}
}

func (n *NetworkMgr) start() {
	socket, error := net.Dial("tcp", network.CONNECT_PORT)
	if error != nil {
		fmt.Println(error)
	}

	//defer socket.Close()

	n.socket = socket
	go n.process()
}

func (n *NetworkMgr) process() {
	for {
		length, err := n.socket.Read(n.recvBuf)
		if err != nil {
			n.socket.Close()
			break
		}
		if length > 0 {
			fmt.Println("Recv response!!")
			head, payload, err := n.parseMessage(n.recvBuf)
			if err == nil {
				fmt.Println("Send response to client manager")
				n.OnRecvMsg(&common.Message{
					CmdType: head.MessageType,
					Body:    payload,
				})
			} else {
				fmt.Println("Recv response!")
			}
		}

		//n.recvBuf = n.recvBuf[:0]
	}
}

func (n *NetworkMgr) parseMessage(message []byte) (*network.Header, []byte, error) {
	head := &network.Header{}
	err := head.Unmarshal(message[:network.HEADER_SIZE])
	if err != nil {
		return head, nil, err
	}
	if head.BodyLength == 0 {
		return head, nil, nil
	}

	return head, message[network.HEADER_SIZE : network.HEADER_SIZE+head.BodyLength], nil
}

func (n *NetworkMgr) SendMessage(msg *common.Message, bodySize uint32) {
	head := network.Header{
		MessageType: msg.CmdType,
		BodyLength:  bodySize,
	}

	headerBuffer := head.Marshal()
	buffer := append(headerBuffer, msg.Body...)

	_, err := n.socket.Write(buffer)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("Send request")
}
