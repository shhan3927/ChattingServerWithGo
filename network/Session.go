package network

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

type SessionInfoDelegate func(SessionInfo)

type SessionInfo struct {
	SessionId uint64
}

type Session struct {
	Socket   net.Conn
	recvCh   chan *Session
	recvBuf  []byte
	wg       *sync.WaitGroup
	isActive bool
	info     SessionInfo
}

func NewSession(inSocket net.Conn, sessionId uint64, recv chan *Session) *Session {
	s := &Session{
		Socket:   inSocket,
		recvCh:   recv,
		recvBuf:  make([]byte, 4096),
		isActive: false,
		info: SessionInfo{
			SessionId: sessionId,
		},
	}
	go s.process()
	return s
}

func (s *Session) process() {
	for {
		n, err := s.Socket.Read(s.recvBuf)
		fmt.Println("Read Socket")
		if nil != err {
			defer func() {
				s.Socket.Close()
			}()

			if io.EOF == err {
				log.Println(err)
				return
			}
			log.Println(err)
			return
		}

		if n > 0 {
			fmt.Println("Read Socket Success..")
			s.recvCh <- s
		}
	}
}

func (s *Session) GetInfo() SessionInfo {
	return s.info
}
