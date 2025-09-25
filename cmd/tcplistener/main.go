package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const port = ":42069"

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("error listening for TCP traffic: %s\n", err.Error())
	}
	defer listener.Close()

	fmt.Println("Listening for TCP traffic on", port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("error: %s\n", err.Error())
		}
		fmt.Println("Accepted connection from", conn.RemoteAddr())

		linesChannel := getLinesChannel(conn)

		for line := range linesChannel {
			fmt.Println(line)
		}

		fmt.Println("Connection to ", conn.RemoteAddr(), "closed")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		defer f.Close()
		defer close(ch)

		currentLine := ""

		for {
			fileBytes := make([]byte, 8)
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
