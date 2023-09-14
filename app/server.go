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

	done := make(chan bool, 1)
	for i := 0; i < 10; i++ {
		go func(done chan bool) {
			conn, err := l.Accept()

			defer func() {
				conn.Close()
				done <- true
			}()

			for {
				if err != nil {
					fmt.Println("Error accepting connection: ", err.Error())
					os.Exit(1)
				}
				var b = make([]byte, 1000)
				_, err = conn.Read(b)
				if err != nil {
					fmt.Println("Failed to read")
				}

				_, err = conn.Write([]byte("+PONG\r\n"))
				if err != nil {
					fmt.Println("Failed to write")
				}
			}
		}(done)
	}

	<-done
}
