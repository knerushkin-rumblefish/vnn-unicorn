package main

import (
	"bytes"
	"encoding/binary"
	"os"
)

func create_binary_blob() {
	size := 4 * 1024 * 1024 // 16GB

	file, err := os.Create("blob.bin")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	for i := 0; i < 1024; i++ {

		buf := make([]byte, 4)

		binary.LittleEndian.PutUint32(buf, uint32(1))
		data := bytes.Repeat(buf, size/4)
		_, err = file.Write(data)

		if err != nil {
			panic(err)
		}
	}
}
