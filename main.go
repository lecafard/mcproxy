package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"github.com/lecafard/mcproxy/schemas"
	"github.com/lecafard/mcproxy/states"
)

const MaxPacketLength = 2097152

var config schemas.Config

func main() {
	var configFile string
	arguments := os.Args
	if len(arguments) == 1 {
		configFile = "config.json"
	} else {
		configFile = arguments[1]
	}

	configBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Println("unable to read config file:", configFile)
		return
	}
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		fmt.Println("unable to unmarshal config file:", configFile)
		return
	}
	fmt.Println("read config file:", configFile)

	l, err := net.Listen("tcp", config.Listen)
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
		state = state.Handler(w, r, &config, &client)
		if state.Handler == nil {
			return
		}
	}
}
