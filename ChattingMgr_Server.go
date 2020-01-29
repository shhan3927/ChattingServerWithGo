package main

import "net"

///////////////////////////////////
type User struct {
	userId uint32
	name   string
	conn   net.Conn
}

///////////////////////////////////
type Room struct {
	roomId uint32
	name   string
	users  []User
}

///////////////////////////////////
type ChattingMgr_Server struct {
	users      map[uint32]*User
	rooms      map[uint32]*Room
	connMap    map[net.Conn]uint32
	roomSeqNum uint32
	userSeqNum uint32
}

func (this *ChattingMgr_Server) Init() {
	this.users = make(map[uint32]*User)
	this.rooms = make(map[uint32]*Room)
	this.connMap = make(map[net.Conn]uint32)
	this.roomSeqNum = 0
	this.userSeqNum = 0
}

func (this *ChattingMgr_Server) AddUser(conn net.Conn) {
	var user User
	user.userId = uint32(this.userSeqNum + 1)
	user.conn = conn
	this.users[user.userId] = &User{userId: uint32(this.userSeqNum + 1), conn: conn}
	this.connMap[conn] = user.userId
	this.userSeqNum++
}

func (this *ChattingMgr_Server) GetUser(userId uint32) (*User, bool) {
	if val, ok := this.users[userId]; ok {
		return val, true
	}

	return nil, false
}

func (this *ChattingMgr_Server) AddRoom(name string) {
	var room Room
	room.roomId = uint32(this.roomSeqNum + 1)
	room.name = name
	this.rooms[room.roomId] = &Room{roomId: uint32(this.roomSeqNum + 1), name: name}
	this.roomSeqNum++
}

func (this *ChattingMgr_Server) GetRoom(roomId uint32) (*Room, bool) {
	if val, ok := this.rooms[roomId]; ok {
		return val, true
	}

	return nil, false
}

func (this *ChattingMgr_Server) ModifyUserNickname(conn net.Conn, nickname string) *User {
	user, success := this.GetUser(this.connMap[conn])
	if !success {
		return nil
	}

	user.name = nickname
	return user
}
