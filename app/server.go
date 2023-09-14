package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go func(conn net.Conn) {
			defer conn.Close()

			for {
				var b = make([]byte, 1000)
				_, err = conn.Read(b)
				if err != nil {
					fmt.Println(err.Error())
					break
				}

				s := strings.Split(string(b), "\r\n")

				if s[2] == "ECHO" {
					_, err = conn.Write([]byte("+PONG\r\n"))
					if err != nil {
						fmt.Println(err.Error())
						break
					}

					b = make([]byte, 1000)
					_, err = conn.Read(b)

					s = strings.Split(string(b), "\r\n")
					resp := strings.Join([]string{"+", s[2], "\r\n"}, "")
					_, err = conn.Write([]byte(resp))
					if err != nil {
						fmt.Println(err.Error())
						break
					}
				} else {
					_, err = conn.Write([]byte("+PONG\r\n"))
					if err != nil {
						fmt.Println(err.Error())
						break
					}
				}

			}
		}(conn)
	}
}
