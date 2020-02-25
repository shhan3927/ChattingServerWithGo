package common

type UserInfo struct {
	Id   uint32
	Name string
}

type Room struct {
	Id     uint32
	Users  map[uint32]*UserInfo
	Master uint32
	Name   string
}

func (r *Room) GetUserCount() int {
	return len(r.Users)
}

func (r *Room) AddUser(userId uint32) {
	if !r.HasUser(userId) {
		r.Users[userId] = &UserInfo{Id: userId}
	}
}

func (r *Room) HasUser(userId uint32) bool {
	_, exist := r.Users[userId]
	return exist
}

func (r *Room) SetMaster(userId uint32) {
	r.Master = userId
}

func NewRoom(roomId uint32, name string) *Room {
	return &Room{
		Id:    roomId,
		Name:  name,
		Users: make(map[uint32]*UserInfo),
	}
}
