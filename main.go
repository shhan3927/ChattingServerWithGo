package main

import (
	"log"
	"net"
)

///////////////////////////////////
type User struct {
	userId uint32
	name   string
	conn   net.Listener
}

///////////////////////////////////
type Room struct {
	roomId uint32
	name   string
	users  []User
}

///////////////////////////////////
type ChattingManager struct {
	users map[uint32]User
	rooms map[uint32]Room
}

func (this *ChattingManager) AddUser(name string, conn net.Listener) {
	var user User
	user.userId = uint32(len(this.users) + 1)
	user.name = name
	user.conn = conn
	this.users[user.userId] = user
}

func (this *ChattingManager) AddRoom(name string) {
	var room Room
	room.roomId = uint32(len(this.rooms) + 1)
	room.name = name
	this.rooms[room.roomId] = room
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
