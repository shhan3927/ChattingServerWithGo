package network

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

type SessionDelegate func(s *Session)

type Session struct {
	Socket  net.Conn
	recvCh  chan *Session
	recvBuf []byte
	wg      *sync.WaitGroup
}

func NewSession(inSocket net.Conn, recv chan *Session) *Session {
	s := &Session{
		Socket:  inSocket,
		recvCh:  recv,
		recvBuf: make([]byte, 4096),
	}
	go s.process()
	return s
}

// func (s *Session) Start() {
// 	go s.process()
// }

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
