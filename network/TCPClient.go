package network

import (
	"net"
)

type TCPClient struct {
	Socket net.Conn
	Data   chan []byte
}

// func (client *TCPClient) Receive() {
// 	for {
// 		//message := make([]byte, MESSAGE_MAX_SIZE)
// 		length, err := *client.socket.Read(client.data)
// 		if err != nil {
// 			*client.socket.Close()
// 			break
// 		}
// 		if length > 0 {
// 			//client.ParseHeader()
// 		}
// 	}
// }
