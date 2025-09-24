package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const inputFilePath = "messages.txt"

func main() {
	messagesFile, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatal(err)
	}

	linesChannel := getLinesChannel(messagesFile)

	for line := range linesChannel {
		fmt.Printf("read: %v\n", line)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		defer f.Close()
		defer close(ch)

		currentLine := ""

		for {
			fileBytes := make([]byte, 8, 8)
			numBytes, err := f.Read(fileBytes)
			if err != nil {
				if currentLine != "" {
					ch <- currentLine
				}

				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error: %s\n", err.Error())
				return
			}

			bytesString := string(fileBytes[:numBytes])
			stringArr := strings.Split(bytesString, "\n")

			for i := 0; i < len(stringArr)-1; i++ {
				ch <- fmt.Sprintf("%s%s", currentLine, stringArr[i])
				currentLine = ""
			}

			currentLine += stringArr[len(stringArr)-1]
		}
	}()

	return ch
}
