package main

import (
	"fmt"
	"log"
	"time"

	"go.bug.st/serial"
)

func GetCoolantTemp() {
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

	// Wrap the command sending and response reading in a loop
	for {
		// Send the "01 05" command to request the coolant temperature
		if _, err = port.Write([]byte("01 05\r")); err != nil {
			log.Fatal(err)
		}

		// Read and log the response from the device
		response := ReadAndLog(port)
		if len(response) < 7 {
			log.Fatal("Invalid response length")
		}

		// Calculate the coolant temperature from the response
		tempByte := response[6]
		tempCelsius := int(tempByte) - 40
		// Print the coolant temperature
		fmt.Printf("Coolant Temperature: %d°C\n", tempCelsius)

		// Delay before sending the next request (e.g., 2 seconds)
		time.Sleep(2 * time.Second)
	}
}
