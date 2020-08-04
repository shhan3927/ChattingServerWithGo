package main

import (
	"flag"
	"strings"

	"github.com/shhan3927/ChattingServerWithGo/chatting_manager"
	"github.com/shhan3927/ChattingServerWithGo/client"
)

func main() {
	flagMode := flag.String("mode", "client", "start in client or server mode")
	flag.Parse()

	if strings.ToLower(*flagMode) == "server" {
		chatMgr := chatting_manager.NewChattingMgr()
		chatMgr.Init()
	} else {
		client.GetChattingMgr().Start()
	}
}
