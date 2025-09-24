package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	messagesFile, err := os.Open("./messages.txt")
	if err != nil {
		log.Fatal(err)
	}

	// bytes will be stored in this
	fileBytes := make([]byte, 8)

	// keeps track of the current line
	currentLine := ""

	// offset to keep track of position in file
	offset := int64(0)

	// loops until io.EOF exits program
	for {

		// reads from file, if err occurs and isn't io.EOF, logs a fatal error
		_, err := messagesFile.ReadAt(fileBytes, offset)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		// converts bytes to a string
		bytesString := string(fileBytes)
		// splits string on new line
		stringArr := strings.Split(bytesString, "\n")
		// adds first element of array to current line
		currentLine += stringArr[0]

		if err == io.EOF {
			fmt.Printf("read: %s\n", currentLine)
			os.Exit(0)
		}

		// if array's length is greater than 1, meaning a line ended
		if len(stringArr) > 1 {
			// prints line and then resets it to the second array element
			fmt.Printf("read: %s\n", currentLine)
			currentLine = stringArr[1]
		}

		// adds to offset and clears byte array
		offset += 8
		fileBytes = make([]byte, 8)
	}
}
