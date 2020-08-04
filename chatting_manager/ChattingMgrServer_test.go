package chatting_manager

import (
	"fmt"
	"testing"
)

func init() {
	c := NewChattingMgr()
	c.Init()
	// chattingMgr.networkMgr = NewTCPServer()
	// chattingMgr.networkMgr.OnConnect = chattingMgr.RegisterUser
	// chattingMgr.networkMgr.OnRecvMessage = chattingMgr.dispatchMessage
}

func TestCreateNickname(t *testing.T) {
	fmt.Println("Starting client...")
	//client.GetChattingMgr().Start()

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
