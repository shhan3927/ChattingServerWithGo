package main

import "encoding/binary"

const HeaderSize uint32 = 12

type ErrorCode uint32

const (
	Error1 ErrorCode = 0
	Error2 ErrorCode = 1
)

type Header struct {
	messageType uint32
	errorCode   ErrorCode
	bodyLength  uint32
}

func (h *Header) Unmarshal(buf []byte) ErrorCode {
	if len(buf) != int(HeaderSize) {
		return Error1
	}
	h.messageType = binary.BigEndian.Uint32(buf[:4])
	h.errorCode = ErrorCode(binary.BigEndian.Uint32(buf[4:8]))
	h.bodyLength = binary.BigEndian.Uint32(buf[8:])
	return nil
}
func (h *Header) Marshal() []byte {
	buf := make([]byte, HeaderSize)
	binary.BigEndian.PutUint32(buf[:4], uint32(h.messageType))
	binary.BigEndian.PutUint32(buf[4:8], uint32(h.errorCode))
	binary.BigEndian.PutUint32(buf[8:], h.bodyLength)
	return buf
}
