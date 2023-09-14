package main

import (
	"fmt"
	"net"
	"os"
)

func main() {

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	conn, err := l.Accept()
	defer conn.Close()

	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	for {
		var b = make([]byte, 1000)
		_, err = conn.Read(b)
		if err != nil {
			fmt.Println("Failed to read")
			os.Exit(1)
		}

		_, err = conn.Write([]byte("+PONG\r\n"))
		if err != nil {
			fmt.Println("Failed to write")
			os.Exit(1)
		}
	}
}
