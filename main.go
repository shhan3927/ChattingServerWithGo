package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/shhan3927/ChattingServerWithGo/network"
)

func StartClientMode() {
	fmt.Println("Starting client...")
	connection, error := net.Dial("tcp", ":4321")
	if error != nil {
		fmt.Println(error)
	}
	client := &Client{socket: connection, data: make([]byte, network.MESSAGE_MAX_SIZE)}
	go client.Receive()

	var nickname string
	fmt.Println("Input nickname : ")
	fmt.Scanf("%s", &nickname)
	client.ReqCreateNickname(nickname)

	for {
		reader := bufio.NewReader(os.Stdin)
		message, _ := reader.ReadString('\n')
		connection.Write([]byte(strings.TrimRight(message, "\n")))
	}
}

func main() {
	flagMode := flag.String("mode", "server", "start in client or server mode")
	flag.Parse()

	if strings.ToLower(*flagMode) == "server" {
		var chatMgr ChattingMgrServer
		chatMgr.Init()
	} else {
		StartClientMode()
	}
}
