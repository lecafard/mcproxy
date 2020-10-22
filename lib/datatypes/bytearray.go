package datatypes

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

var InvalidBytearray = errors.New("invalid byte array")
var InvalidByteCount = errors.New("invalid amount of bytes")

func ReadVarBytearray(r io.Reader) ([]byte, error) {
	length, err := ReadVarint(r.(io.ByteReader))
	if err != nil {
		fmt.Println("failed to get server name length")
		return nil, InvalidBytearray
	}

	buf := make([]byte, length)
	n, err := r.Read(buf)
	if n != int(length) {
		return nil, InvalidByteCount
	}
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func WriteVarBytearray(w io.Writer, b []byte) error {
	WriteVarint(w.(io.ByteWriter), int32(len(b)))
	n, err := w.Write(b)
	if n != len(b) {
		return InvalidByteCount
	}
	if err != nil {
		return err
	}
	return nil
}

func EncodeVarBytearray(b []byte) []byte {
	buf := bytes.NewBuffer(nil)
	WriteVarBytearray(buf, b)
	return buf.Bytes()
}
