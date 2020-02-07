package jclass

import (
	"bytes"
	"encoding/binary"
	"unicode/utf16"
)

func decodeUtf16(s []byte) (string, error) {
	// Not very efficient, but correct. Can speed it up later.
	pairs := make([]uint16, len(s)/2)
	if err := binary.Read(bytes.NewReader(s), binary.BigEndian, &pairs); err != nil {
		return "", err
	}
	return string(utf16.Decode(pairs)), nil
}
