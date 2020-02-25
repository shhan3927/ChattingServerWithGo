package common

import "github.com/shhan3927/ChattingServerWithGo/protomessage"

func GetUserInfoWithProto(userInfo *protomessage.UserInfo) *UserInfo {
	return &UserInfo{
		Id:   userInfo.Id,
		Name: userInfo.Name,
	}
}

func GetRoomWithProto(roomInfo *protomessage.RoomInfo) *Room {
	room := &Room{
		Id:     roomInfo.Id,
		Name:   roomInfo.Name,
		Master: roomInfo.MasterUserId,
	}

	room.Users = make(map[uint32]*UserInfo)

	for _, userInfo := range roomInfo.Users {
		room.Users[userInfo.Id] = GetUserInfoWithProto(userInfo)
	}

	return room
}
