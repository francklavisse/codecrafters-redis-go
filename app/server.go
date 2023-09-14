package main

import (
	"fmt"
	"net"
	"os"
	s "strings"
)

func readCmd(conn net.Conn) ([]string, error) {
	var b = make([]byte, 1000)
	n, err := conn.Read(b)

	if err != nil {
		return nil, err
	}

	receivedData := b[:n]
	cmd := s.Split(string(receivedData), "\r\n")

	return cmd, nil
}

func getResponse(cmd []string) string {
	db := make(map[string]string)

	fmt.Println(cmd)

	switch cmd[2] {
	case "ECHO", "echo":
		return "+" + cmd[4] + "\r\n"
	case "SET", "set":
		db[cmd[4]] = cmd[6]
		return "+OK\r\n"
	case "GET", "get":
		return "+" + db[cmd[4]] + "\r\n"
	default:
		return "+PONG\r\n"
	}
}

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
				cmd, err := readCmd(conn)
				if err != nil {
					fmt.Println(err.Error())
					break
				}

				resp := getResponse(cmd)

				_, err = conn.Write([]byte(resp))
				if err != nil {
					fmt.Println(err.Error())
					break
				}
			}
		}(conn)
	}
}
