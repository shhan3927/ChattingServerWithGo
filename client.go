package main

import (
	"log"
	"net"
	"fmt"
)

func main() {
	conn, err := net.Dial("tcp", ":4321")
	if nil != err {
		log.Println(err)
	}

	var name string
	fmt.Println("Input your name")
	fmt.Scan(&name)
}
