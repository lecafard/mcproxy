package states

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os/exec"
	"time"

	"github.com/lecafard/mcproxy/lib"
	"github.com/lecafard/mcproxy/lib/datatypes"
	"github.com/lecafard/mcproxy/schemas"
)

func StateHandshaking(w *bufio.Writer, r *bufio.Reader, s *schemas.Config, c *Client) State {
	var err error

	packet, err := lib.ReadPacket(r)
	if err != nil {
		fmt.Println("error reading packet", err)
		return State{nil}
	}
	c.RawHandshake = packet.Payload

	rbuf := bytes.NewReader(packet.Payload)
	c.ProtocolVersion, err = datatypes.ReadVarint(rbuf)
	if err != nil {
		fmt.Println("failed to get protocol version")
		return State{nil}
	}

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

	// check if server is up
	if s.Servers[c.ServerName] != nil {
		c.ServerConfig = s.Servers[c.ServerName]
	} else {
		c.ServerConfig = s.Servers["default"]
	}

	if c.ServerConfig == nil {
		lib.WriteKick(w, "Server endpoint doesn't exist...")
		return State{nil}
	}

	dialler := net.Dialer{Timeout: time.Duration(s.Timeout) * time.Second}
	_, err = dialler.Dial("tcp", c.ServerConfig.Proxy)
	if err == nil {
		return State{StateProxy}
	}

	// read next state
	next, err := datatypes.ReadVarint(rbuf)
	if err != nil {
		fmt.Println("invalid next state")
		return State{nil}
	}

	if next == 1 {
		return State{StateStatus}
	} else if next == 2 {
		// run startup script
		cmd := exec.Command("sh", "-c", c.ServerConfig.Startup)
		cmd.Run()
		lib.WriteKick(w, "Submitted startup request! Try again in a bit.")
		return State{nil}
	} else {
		return State{nil}
	}
}

func StateStatus(w *bufio.Writer, r *bufio.Reader, s *schemas.Config, c *Client) State {
	if c.ServerConfig == nil {
		lib.WriteKick(w, "Server endpoint doesn't exist...")
		return State{nil}
	}
	status, err := json.Marshal(schemas.ServerStatus{
		Version: schemas.ServerVersion{
			Name:     "MCProxy 0.1",
			Protocol: 753,
		},
		Description: schemas.ServerDescription{
			Text: c.ServerConfig.Description,
		},
		Players: schemas.ServerPlayers{
			Maximum: 0,
			Online:  0,
		},
	})
	if err != nil {
		fmt.Println("error generating packet", err)
		return State{nil}
	}

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

func StateProxy(w *bufio.Writer, r *bufio.Reader, s *schemas.Config, c *Client) State {
	stop := make(chan bool)

	if c.ServerConfig == nil {
		lib.WriteKick(w, "Server endpoint doesn't exist...")
		return State{nil}
	}

	dialler := net.Dialer{Timeout: time.Duration(s.Timeout) * time.Second}
	p, err := dialler.Dial("tcp", c.ServerConfig.Proxy)
	if err != nil {
		lib.WriteKick(w, "Server starting up...")
		return State{nil}
	}

	rp := bufio.NewReader(p)
	wp := bufio.NewWriter(p)

	lib.WritePacket(wp, 0x0, c.RawHandshake)

	// do the proxy
	go pipe(r, wp, stop)
	pipe(rp, w, stop)
	<-stop
	return State{nil}
}

func pipe(src *bufio.Reader, dst *bufio.Writer, complete chan bool) {
	var err error = nil
	var bytes []byte = make([]byte, 512)
	var read int = 0

	for {
		read, err = src.Read(bytes)
		if err != nil {
			complete <- true
			break
		}
		_, err = dst.Write(bytes[:read])
		dst.Flush()
		if err != nil {
			complete <- true
			break
		}
	}
}
