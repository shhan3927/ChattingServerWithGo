package main

import (
	"log"
	"net"
)

///////////////////////////////////
type User struct {
	name string
	conn net.Listener
}

///////////////////////////////////
type Room struct {
	name  string
	users []User
}

///////////////////////////////////
type ChattingManager struct {
	users map[string]User
	rooms map[string]Room
}

func (this *ChattingManager) AddUser(name string, conn net.Listener){
	var user User;
	user.name = name;
	user.conn = conn;
	this.users[name] = user
}

///////////////////////////////////
func main() {
	conn, err := net.Listen("tcp", ":4321")
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()

	var room Room
	room.name = "ARoom"
	for {
		user, err := conn.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		defer user.Close()
	}
}
