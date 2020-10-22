package states

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/lecafard/mcproxy/schemas"

	"github.com/lecafard/mcproxy/lib"
	"github.com/lecafard/mcproxy/lib/datatypes"
)

func StateHandshaking(w *bufio.Writer, r *bufio.Reader, c *Client) State {
	var err error

	packet, err := lib.ReadPacket(r)
	if err != nil {
		fmt.Println("error reading packet", err)
		return State{nil}
	}

	rbuf := bytes.NewReader(packet.Payload)
	c.ProtocolVersion, err = datatypes.ReadVarint(rbuf)
	if err != nil {
		fmt.Println("failed to get protocol version")
		return State{nil}
	}
	fmt.Println("protocol version:", c.ProtocolVersion)

	// read server name
	nameBuf, err := datatypes.ReadVarBytearray(rbuf)
	if err != nil {
		fmt.Println("error reading ServerName")
	}
	c.ServerName = string(nameBuf)

	bufPort := make([]byte, 2)
	_, err = rbuf.Read(bufPort)
	if err != nil {
		fmt.Println("error reading ServerPort", err)
		return State{nil}
	}
	c.ServerPort = binary.BigEndian.Uint16(bufPort)

	// read next state
	next, err := datatypes.ReadVarint(rbuf)
	if err != nil {
		fmt.Println("invalid next state")
		return State{nil}
	}

	if next == 1 {
		return State{StateStatus}
	} else if next == 2 {
		return State{StateLogin}
	} else {
		return State{nil}
	}
}

func StateStatus(w *bufio.Writer, r *bufio.Reader, c *Client) State {
	status, err := json.Marshal(schemas.ServerStatus{
		Version: schemas.ServerVersion{
			Name:     "MCProxy 0.1",
			Protocol: 753,
		},
		Description: schemas.ServerDescription{
			Text: "lol",
		},
		Players: schemas.ServerPlayers{
			Maximum: 20,
			Online:  0,
		},
	})

	if err != nil {
		fmt.Println("error generating packet", err)
		return State{nil}
	}
	fmt.Println(status)

	for {
		packet, err := lib.ReadPacket(r)
		if err != nil {
			fmt.Println("error reading packet", err)
			return State{nil}
		}

		rbuf := bytes.NewReader(packet.Payload)
		switch packet.ID {
		// server list info
		case 0x00:
			lib.WritePacket(w, 0x0, datatypes.EncodeVarBytearray(status))
			break
		// ping
		case 0x01:
			bufLong := make([]byte, 8)
			rbuf.Read(bufLong)
			lib.WritePacket(w, 0x1, bufLong)
		}
	}
}

func StateLogin(w *bufio.Writer, r *bufio.Reader, c *Client) State {
	lib.WriteKick(w, "Server starting up...")
	return State{nil}
}
