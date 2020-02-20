package network

import (
	"sync"
	"testing"
)

var sg sync.WaitGroup

func init() {
	// networkMgr := NewTCPServer()
	// networkMgr.Start(":4323")
	// networkMgr.OnRecvMessage = func(header *Header, payload []byte) {
	// 	switch header.MessageType {
	// 	case uint32(protomessage.MessageType_kCreateNicknameRequest):
	// 		var packet protomessage.CreateNicknameRequest
	// 		e := proto.Unmarshal(payload, &packet)
	// 		if e != nil {
	// 			fmt.Println(e)
	// 		}

	// 		fmt.Println(packet.Name)
	// 		sg.Done()
	// 	}
	// }
}

func TestCreateNickname(t *testing.T) {
	// fmt.Println("Starting client...")
	// connection, error := net.Dial("tcp", ":4323")
	// if error != nil {
	// 	fmt.Println(error)
	// }
	// client := &Client{socket: connection, data: make([]byte, MESSAGE_MAX_SIZE)}
	// go client.Receive()

	// sg.Add(1)
	// client.ReqCreateNickname("nickname")
	// sg.Wait()

	// for {
	// 	reader := bufio.NewReader(os.Stdin)
	// 	message, _ := reader.ReadString('\n')
	// 	connection.Write([]byte(strings.TrimRight(message, "\n")))
	// }
}
