package main

import (
	"encoding/binary"
	"errors"
)

const HEADER_SIZE uint32 = 12

var HeaderELength = errors.New("Invalid header length")

type Header struct {
	messageType uint32
	err         uint32
	bodyLength  uint32
}

func (h *Header) Unmarshal(buf []byte) error {
	if len(buf) != int(HEADER_SIZE) {
		return HeaderELength
	}
	h.messageType = binary.BigEndian.Uint32(buf[:4])
	h.err = binary.BigEndian.Uint32(buf[4:8])
	h.bodyLength = binary.BigEndian.Uint32(buf[8:])
	return nil
}

func (h *Header) Marshal() []byte {
	buf := make([]byte, HEADER_SIZE)
	binary.BigEndian.PutUint32(buf[:], uint32(h.messageType))
	binary.BigEndian.PutUint32(buf[4:8], uint32(h.err))
	binary.BigEndian.PutUint32(buf[8:], h.bodyLength)
	return buf
}
