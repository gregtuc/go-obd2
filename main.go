package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"go.bug.st/serial"
)

func main() {
	// Open the COM4 port with the default settings (baud rate = 9600)
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

	// Send the "ATZ" command to reset the device and prepare it for communication
	if _, err = port.Write([]byte("ATZ\r")); err != nil {
		log.Fatal(err)
	}
	// Read and log the response from the device
	readAndLog(port)

	// Send the "01 05" command to request the coolant temperature
	if _, err = port.Write([]byte("01 05\r")); err != nil {
		log.Fatal(err)
	}

	// Sanity check - response should normally be at least 7 bytes long
	response := readAndLog(port)
	if len(response) < 7 {
		log.Fatal("Invalid response length")
	}

	// Calculate the coolant temperature from the response
	// Assuming the temperature is the 7th byte of the response and converting it to Celsius
	tempByte := response[6]
	tempCelsius := int(tempByte) - 40
	fmt.Printf("Coolant Temperature: %d°C\n", tempCelsius)
}

// readAndLog reads from the serial port until the ">" character (prompt) is detected,
// indicating the end of the device's response.
func readAndLog(port serial.Port) []byte {
	var response bytes.Buffer // Buffer to hold the response
	buf := make([]byte, 128)  // Temporary buffer to hold data read from the port
	for {
		// Read data into the buffer
		n, err := port.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		// If no data was read, exit the loop
		if n == 0 {
			break
		}
		// Write the read data into the response buffer
		response.Write(buf[:n])
		// If the read data contains the ">" character, exit the loop
		if bytes.Contains(buf[:n], []byte(">")) {
			break
		}
	}
	// Log the number of bytes received, the response as a string, and as hex
	log.Printf("Received %d bytes: %s (hex: %s)", response.Len(), response.String(), hex.EncodeToString(response.Bytes()))
	// Return the response as a byte slice
	return response.Bytes()
}
