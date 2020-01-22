package protomessage

import (
	"reflect"
)

const MESSAGE_MAX_SIZE uint32 = 4096

func GetPacketType(i interface{}) (MessageType, uint32) {
	typeValue := MessageType_value["k"+reflect.TypeOf(i).String()]
	return MessageType(typeValue), uint32(typeValue)
}
