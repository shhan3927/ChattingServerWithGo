package network

import (
	"fmt"
	"log"
	"net"
)

type CmdType uint32
type ErrorCode uint32
type MessageDelegate func(SessionInfo, *Message)
type MessageDelegate2 func(SessionInfo, *Message, uint32)

type Message struct {
	CmdType uint32
	ErrCode ErrorCode
	Body    []byte
}

func NewTCPServer() *TCPServer {
	server := &TCPServer{
		sessions:      make(map[uint64]*Session),
		connectCh:     make(chan *Session, 10000),
		disconnectCh:  make(chan *Session, 10000),
		recvCh:        make(chan *Session, 10000),
		sendCh:        make(chan *Session, 10000),
		MessageCh:     make(chan *Message, 10000),
		sessionSeqNum: 0,
	}
	return server
}

type TCPServer struct {
	clients  map[net.Conn]*TCPClient
	sessions map[uint64]*Session
	Connect  chan *TCPClient

	connectCh    chan *Session
	disconnectCh chan *Session
	recvCh       chan *Session
	sendCh       chan *Session
	MessageCh    chan *Message

	listener net.Listener

	OnConnect     SessionInfoDelegate
	OnRecvMessage MessageDelegate
	OnSendMessage MessageDelegate2

	sessionSeqNum uint64
}

func (s *TCPServer) Start(address string) (err error) {
	fmt.Println("start tcp server...")
	s.listener, err = net.Listen("tcp", address)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer s.listener.Close()

	go s.accept()
	s.process()

	return
}

func (s *TCPServer) accept() {
	for {
		fmt.Println("tcp server accept...")
		socket, err := s.listener.Accept()
		if err != nil {
			fmt.Println(err)
		}

		s.sessionSeqNum++
		s.connectCh <- NewSession(socket, s.sessionSeqNum, s.recvCh)
	}
}

func (s *TCPServer) process() {
	for {
		select {
		case session := <-s.connectCh:
			s.registerSession(session)
			s.OnConnect(session.info)
		case session := <-s.disconnectCh:
			s.unregisterSession(session)
		case session := <-s.recvCh:
			head, payload, err := s.parseMessage(session.recvBuf)
			if err != nil {
				fmt.Println(err)
			} else {
				s.OnRecvMessage(session.info, &Message{
					CmdType: head.MessageType,
					ErrCode: 0,
					Body:    payload,
				})
			}
		}
	}
}

func (s *TCPServer) registerSession(session *Session) {
	session.isActive = true
	s.sessions[session.info.SessionId] = session
}

func (s *TCPServer) unregisterSession(session *Session) {
	session.isActive = false
	delete(s.sessions, session.info.SessionId)
}

func (s *TCPServer) parseMessage(message []byte) (*Header, []byte, error) {
	head := &Header{}
	err := head.Unmarshal(message[:HEADER_SIZE])
	if err != nil {
		return head, nil, err
	}
	if head.BodyLength == 0 {
		return head, nil, nil
	}

	return head, message[HEADER_SIZE : HEADER_SIZE+head.BodyLength], nil
}

func (s *TCPServer) SendMessage(sessionInfo SessionInfo, msg *Message, bodySize uint32) {
	if session, exist := s.sessions[sessionInfo.SessionId]; exist {
		head := Header{
			MessageType: msg.CmdType,
			BodyLength:  bodySize,
		}

		headerBuffer := head.Marshal()
		buffer := append(headerBuffer, msg.Body...)

		_, err := session.Socket.Write(buffer)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
