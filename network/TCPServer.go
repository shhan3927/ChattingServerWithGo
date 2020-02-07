package network

import (
	"fmt"
	"io"
	"log"
	"net"
)

type TCPServer struct {
	clients map[net.Conn]*TCPClient
	Connect chan *TCPClient
}

func (s *TCPServer) Init() {
	s.clients = make(map[net.Conn]*TCPClient)
	s.Connect = make(chan *TCPClient)
}

func (s *TCPServer) Start(address string) {
	fmt.Println("start tcp server...")
	listener, error := net.Listen("tcp", address)
	if error != nil {
		fmt.Println(error)
	}

	defer listener.Close()

	for {
		fmt.Println("tcp server accept...")
		socket, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}

		s.RegisterSocket(socket, nil)
		client := &TCPClient{Socket: socket, Data: make(chan []byte)}
		s.Connect <- client
		go s.HandleMessage(client)
	}
}

func (s *TCPServer) RegisterSocket(socket net.Conn, client *TCPClient) {
	if ret, isExist := s.clients[socket]; !isExist || ret == nil {
		s.clients[socket] = client
	}
}

func (s *TCPServer) UnregisterSocket(socket net.Conn) {
	delete(s.clients, socket)
}

func (s *TCPServer) HandleMessage(client *TCPClient) {
	buffer := make([]byte, MESSAGE_MAX_SIZE)

	for {
		n, err := client.Socket.Read(buffer)
		fmt.Println("Read Socket")
		if nil != err {
			defer func() {
				s.UnregisterSocket(client.Socket)
				client.Socket.Close()
			}()

			if io.EOF == err {
				log.Println(err)
				return
			}
			log.Println(err)
			return
		}

		if n > 0 {
			fmt.Println("Read Socket Success..")
			client.Data <- buffer
		}
	}
}
