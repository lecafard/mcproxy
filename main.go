package main

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/lecafard/mcproxy/states"
)

const MaxPacketLength = 2097152

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide port number")
		return
	}

	port := ":" + arguments[1]
	l, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		fmt.Println("accepted new connection")
		if err != nil {
			fmt.Println("error accepting connection")
		}
		go handler(c)
	}
}

func handler(c net.Conn) {
	defer c.Close()

	w := bufio.NewWriter(c)
	r := bufio.NewReader(c)

	// State machine logic
	state := states.State{Handler: states.StateHandshaking}
	client := states.Client{}
	for {
		state = state.Handler(w, r, &client)
		if state.Handler == nil {
			return
		}
	}
}
