package protomessage

import (
	"reflect"
)

const MessageBodySizeMax uint32 = 4096

func GetPacketType(i interface{}) MessageType int32 {
	typeValue := MessageType_value["k"+reflect.TypeOf(i).String()]
	return MessageType(typeValue), typeValue 
}


