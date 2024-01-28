package main

import (
	"bytes"
	"encoding/hex"
	"log"

	"go.bug.st/serial"
)

func ReadAndLog(port serial.Port) []byte {
	var response bytes.Buffer
	buf := make([]byte, 128)
	for {
		n, err := port.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		if n == 0 {
			break
		}
		response.Write(buf[:n])
		if bytes.Contains(buf[:n], []byte(">")) {
			break
		}
	}
	log.Printf("Received %d bytes: %s (hex: %s)", response.Len(), response.String(), hex.EncodeToString(response.Bytes()))
	return response.Bytes()
}
