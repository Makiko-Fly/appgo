package main

import (
	"github.com/vitess/go/cgzip"
	"fmt"
	"bytes"
)

func main() {
	buf := new(bytes.Buffer)
	level := cgzip.Z_DEFAULT_COMPRESSION
	bufSize := 10240
	writer, err := cgzip.NewWriterLevelBuffer(buf, level, bufSize)
	if err != nil {
		fmt.Println(err)
	}
	defer writer.Close()



	writer.Write([]byte("abcdefghijklmnopqrstuvwxyz"))
	writer.Flush()

	buf.Reset()
	writer


	fmt.Println(buf.String())
}
