package common

type Room struct {
	Id     uint32
	Users  []uint32
	Master uint32
	Name   string
}

func (r *Room) GetUserCount() int {
	return len(r.Users)
}

func (r *Room) AddUser(userId uint32) {
	if !r.HasUser(userId) {
		r.Users = append(r.Users, userId)
	}
}

func (r *Room) HasUser(userId uint32) bool {
	for _, n := range r.Users {
		if userId == n {
			return true
		}
	}

	return false
}

func (r *Room) SetMaster(userId uint32) {
	r.Master = userId
}

func NewRoom(roomId uint32, name string) *Room {
	return &Room{
		Id:   roomId,
		Name: name,
	}
}
