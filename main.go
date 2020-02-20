package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/shhan3927/ChattingServerWithGo/client"
	"github.com/shhan3927/ChattingServerWithGo/network"
)

func StartClientMode() {
	fmt.Println("Starting client...")
	connection, error := net.Dial("tcp", network.CONNECT_PORT)
	if error != nil {
		fmt.Println(error)
	}
	client := &network.Client{Socket: connection, Data: make([]byte, network.MESSAGE_MAX_SIZE)}
	go client.Receive()

	var name string
	fmt.Println("Input nickname : ")
	fmt.Scanf("%s", &name)
	client.ReqCreateNickname(name)

	// fmt.Println("Input Room Name : ")
	// fmt.Scanf("%s", &name)
	// client.ReqCreateRoom(name)

	for {
		reader := bufio.NewReader(os.Stdin)
		message, _ := reader.ReadString('\n')
		fmt.Println(message)
		client.ReqCreateRoom("Test")

		//message, _ := reader.ReadString('\n')
		//connection.Write([]byte(strings.TrimRight(message, "\n")))
	}
}

func main() {
	flagMode := flag.String("mode", "client", "start in client or server mode")
	flag.Parse()

	if strings.ToLower(*flagMode) == "server" {
		chatMgr := NewChattingMgr()
		chatMgr.Init()
	} else {
		client.GetChattingMgr().Start()
		//StartClientMode()
	}
}
