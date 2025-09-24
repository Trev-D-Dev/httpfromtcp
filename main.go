package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	messagesFile, err := os.Open("./messages.txt")
	if err != nil {
		log.Fatal(err)
	}

	fileBytes := make([]byte, 8)

	offset := int64(0)
	for {
		_, err := messagesFile.ReadAt(fileBytes, offset)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}

		fmt.Printf("read: %s\n", fileBytes)

		if err == io.EOF {
			os.Exit(0)
		}

		offset += 8
	}
}
