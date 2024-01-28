package main

import (
	"fmt"
	"log"
	"time"

	"go.bug.st/serial"
)

func GetRPMs() {
	mode := &serial.Mode{BaudRate: 9600}
	port, err := serial.Open("COM4", mode)
	if err != nil {
		log.Fatal(err)
	}
	defer port.Close()

	err = port.SetReadTimeout(10 * time.Second)
	if err != nil {
		log.Fatal(err)
	}

	if _, err = port.Write([]byte("ATZ\r")); err != nil {
		log.Fatal(err)
	}
	ReadAndLog(port)

	for {
		// Send the "01 0C" command to request the engine RPM
		if _, err = port.Write([]byte("01 0C\r")); err != nil {
			log.Fatal(err)
		}

		// Read and log the response from the device
		response := ReadAndLog(port)
		if len(response) < 8 { // Expecting at least 8 bytes for RPM response
			log.Fatal("Invalid response length")
		}

		// Calculate the engine RPM from the response
		rpm := ((int(response[6]) * 256) + int(response[7])) / 4
		// Print the engine RPM
		fmt.Printf("Engine RPM: %d\n", rpm)

		// Delay before sending the next request (e.g., 2 seconds)
		time.Sleep(2 * time.Second)
	}
}
