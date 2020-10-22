package lib

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/lecafard/mcproxy/lib/datatypes"
)

// Packet minecraft
type Packet struct {
	Length  int
	ID      int32
	Payload []byte
}

var InvalidPacketSize = errors.New("invalid packet size")

// ReadPacket
func ReadPacket(r io.Reader) (*Packet, error) {
	packet := Packet{}
	length, err := datatypes.ReadVarint(r.(io.ByteReader))
	if err != nil {
		return nil, err
	}

	data := make([]byte, length)
	n, err := r.Read(data)
	dr := bytes.NewReader(data)

	if err != nil {
		return nil, err
	} else if n != int(length) {
		return nil, InvalidPacketSize
	}

	packet.ID, err = datatypes.ReadVarint(dr)
	if err != nil {
		return nil, err
	}

	packet.Payload = make([]byte, dr.Len())
	packet.Length = dr.Len()
	dr.Read(packet.Payload)

	return &packet, nil
}

// WritePacket writes a minecraft packet
func WritePacket(w *bufio.Writer, id int32, payload []byte) {
	l := int32(len(payload) + len(datatypes.EncodeVarint(id)))
	datatypes.WriteVarint(w, l)
	datatypes.WriteVarint(w, id)
	w.Write(payload)
	w.Flush()
}

func WriteKick(w *bufio.Writer, message string) {
	fmt.Println("kicking client")
	buf, err := json.Marshal(message)
	if err != nil {
		fmt.Println("error writing kick")
		return
	}

	WritePacket(w, 0x0, datatypes.EncodeVarBytearray(buf))
	return
}
