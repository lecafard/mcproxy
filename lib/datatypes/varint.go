package datatypes

import (
	"bytes"
	"errors"
	"io"
)

var InvalidVarint = errors.New("Invalid Varint")

func ReadVarint(r io.ByteReader) (int32, error) {
	numRead := 0
	var result int32 = 0
	var b byte = 128
	var err error
	for (b & 0b10000000) != 0 {
		b, err = r.ReadByte()
		if err != nil {
			return 0, err
		}

		value := int32(b & 0b01111111)
		result |= (value << (7 * numRead))

		numRead++
		if numRead > 5 {
			return 0, InvalidVarint
		}
	}

	return result, nil
}

func WriteVarint(w io.ByteWriter, v int32) error {
	if v == 0 {
		w.WriteByte(0)
		return nil
	}

	for v != 0 {
		var temp byte
		temp = byte(v & 0b01111111)

		// Note: >>> means that the sign bit is shifted with the rest of the number rather than being left alone
		v >>= 7
		if v != 0 {
			temp |= 0b10000000
		}
		err := w.WriteByte(temp)
		if err != nil {
			return err
		}
	}

	return nil
}

func EncodeVarint(v int32) []byte {
	buf := bytes.NewBuffer(nil)
	WriteVarint(buf, v)
	return buf.Bytes()
}
