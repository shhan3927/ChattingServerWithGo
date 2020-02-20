package client

import (
	"bufio"
	"fmt"
	"os"
)

type UIState uint

const (
	None       UIState = 0
	Main       UIState = 1
	NickName   UIState = 2
	CreateRoom UIState = 3
	ShowRoom   UIState = 4
	InRoom     UIState = 5
)

type UIInterface interface {
	Execute()
	//NextState() UIState
}

//////////////////////////
type InputNicknameUI struct {
}

func (u InputNicknameUI) Execute() {
	//var name string
	fmt.Println("=============================")
	fmt.Println("=== SH's Chatting Program ===")
	fmt.Println()
	fmt.Println("Input your nickname : ")
	//fmt.Scanf("%s", &name)

	// Request create nickname
	GetChattingMgr().ReqCreateNickname("TEst")
}

//////////////////////////
type ShowRoomUI struct {
}

func (u ShowRoomUI) Execute() {
	fmt.Println("=============================")
	fmt.Println("=======   Room List   =======")
	fmt.Println("=============================")
	fmt.Println("Input room number to join(If you want to make room, input \"c\"): ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')

	// join room or create room
	if input == "c" {

	} else {

	}
}

//////////////////////////
type MainUI struct {
}

func (u MainUI) Execute() {

}

//////////////////////////
type CreateRoomUI struct {
}

func (u CreateRoomUI) Execute() {
	var name string
	fmt.Println("=============================")
	fmt.Println("=======  Create Room  =======")
	fmt.Println("=============================")
	fmt.Println("Input room name : ")
	fmt.Scanf("%s", &name)

	// request create room
	GetChattingMgr().ReqCreateRoom(name)
}

//////////////////////////
type InRoomUI struct {
}

func (u InRoomUI) Execute() {
	fmt.Println("=============================")
	fmt.Println("=======  Room : ", "Room1")
	fmt.Println("=============================")

	for {
		// reader := bufio.NewReader(os.Stdin)
		// input, _ := reader.ReadString('\n')
		// send message
	}
}

////////////////////////////////////
type ChattingUI struct {
	beforeState  UIState
	currentState UIState
	stateMap     map[UIState]UIInterface
}

func (c *ChattingUI) Start() {
	chattingMgr := GetChattingMgr()
	chattingMgr.OnCreateNickname = c.EventCreateNickname
	chattingMgr.OnCreateRoom = c.EventCreateRoom
	c.stateMap[Main] = MainUI{}
	c.stateMap[NickName] = InputNicknameUI{}
	c.stateMap[CreateRoom] = CreateRoomUI{}
	c.stateMap[InRoom] = InRoomUI{}
	c.SetState(NickName)
}

func (c *ChattingUI) SetState(state UIState) {
	if c.currentState == state {
		return
	}

	c.beforeState = c.currentState
	c.currentState = state
	c.stateMap[c.currentState].Execute()
}

func NewChattingUI() *ChattingUI {
	return &ChattingUI{
		beforeState:  None,
		currentState: None,
		stateMap:     make(map[UIState]UIInterface),
	}
}

func (c *ChattingUI) EventCreateNickname() {
	c.SetState(CreateRoom)
}

func (c *ChattingUI) EventCreateRoom() {
	c.SetState(InRoom)
}
