package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

func StartServerMode() {
	fmt.Println("Starting server...")
	listener, error := net.Listen("tcp", ":4321")
	if error != nil {
		fmt.Println(error)
	}
	var server ChattingServer
	server.Init()

	go server.Start()
	for {
		connection, _ := listener.Accept()
		if error != nil {
			fmt.Println(error)
		}
		client := &Client{socket: connection, data: make([]byte, MESSAGE_MAX_SIZE)}
		server.register <- client
		go server.HandleMessage(client)
	}
}

func StartClientMode() {
	fmt.Println("Starting client...")
	connection, error := net.Dial("tcp", ":4321")
	if error != nil {
		fmt.Println(error)
	}
	client := &Client{socket: connection}
	go client.Receive()

	client.ReqCreateNickname("Test")

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
		StartServerMode()
	} else {
		StartClientMode()
	}
}
