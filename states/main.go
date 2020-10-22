package states

import (
	"bufio"

	"github.com/lecafard/mcproxy/schemas"
)

type Client struct {
	ProtocolVersion int32
	ServerName      string
	ServerPort      uint16
	ServerConfig    *schemas.ConfigServer
	RawHandshake    []byte
}

type State struct {
	Handler func(*bufio.Writer, *bufio.Reader, *schemas.Config, *Client) State
}
