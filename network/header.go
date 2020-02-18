package network

import (
	"encoding/binary"
	"errors"
)

const HEADER_SIZE uint32 = 12

var HeaderELength = errors.New("Invalid header length")

type Header struct {
	MessageType uint32
	Err         uint32
	BodyLength  uint32
}

func (h *Header) Unmarshal(buf []byte) error {
	if len(buf) != int(HEADER_SIZE) {
		return HeaderELength
	}
	h.MessageType = binary.BigEndian.Uint32(buf[:4])
	h.Err = binary.BigEndian.Uint32(buf[4:8])
	h.BodyLength = binary.BigEndian.Uint32(buf[8:])
	return nil
}

func (h *Header) Marshal() []byte {
	buf := make([]byte, HEADER_SIZE)
	binary.BigEndian.PutUint32(buf[:], uint32(h.MessageType))
	binary.BigEndian.PutUint32(buf[4:8], uint32(h.Err))
	binary.BigEndian.PutUint32(buf[8:], h.BodyLength)
	return buf
}
