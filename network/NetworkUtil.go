package network

import (
	"reflect"
	"strings"

	"github.com/shhan3927/ChattingServerWithGo/protomessage"
)

const MESSAGE_MAX_SIZE uint32 = 4096

func GetPacketType(i interface{}) (protomessage.MessageType, uint32) {
	aaa := reflect.TypeOf(i).String()
	buffer := strings.Split(aaa, ".")
	typeValue := protomessage.MessageType_value["k"+buffer[len(buffer)-1]]
	return protomessage.MessageType(typeValue), uint32(typeValue)
}
