package states

import "bufio"

type Client struct {
	ProtocolVersion int32
	ServerName      string
	ServerPort      uint16
}

type State struct {
	Handler func(w *bufio.Writer, r *bufio.Reader, c *Client) State
}
