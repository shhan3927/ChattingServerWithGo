package common

type CmdType uint32
type ErrorCode uint32

type Message struct {
	CmdType uint32
	ErrCode ErrorCode
	Body    []byte
}
