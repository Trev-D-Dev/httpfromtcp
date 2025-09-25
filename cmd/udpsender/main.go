package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const port = "localhost:42069"

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", port)
	if err != nil {
		log.Fatalf("error resolving udp address: %s", err.Error())
	}

	conn, err := net.DialUDP("udp", udpAddr, udpAddr)
	if err != nil {
		log.Fatalf("error creating udp connection: %s", err.Error())
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")

		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("error reading input: %s", err.Error())
		}

		_, err = conn.Write([]byte(input))
		if err != nil {
			log.Fatalf("error writing input to udp connection: %s", err.Error())
		}
	}
}
