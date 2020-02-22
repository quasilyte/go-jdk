package cmdutil

import (
	"os"

	"github.com/quasilyte/go-jdk/jclass"
)

func DecodeClassFile(filename string) (*jclass.File, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var dec jclass.Decoder
	return dec.Decode(f)
}
