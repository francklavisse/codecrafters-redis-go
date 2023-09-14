package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	s "strings"
	"time"
)

func readCmd(conn net.Conn) ([]string, error) {
	var b = make([]byte, 1000)
	n, err := conn.Read(b)

	if err != nil {
		return nil, err
	}

	receivedData := b[:n]
	cmd := s.Split(string(receivedData), "\r\n")
	fmt.Println(cmd)
	return cmd, nil
}

type Data struct {
	PX        int
	CreatedAt time.Time
	Value     string
}

var db = make(map[string]Data)

func getResponse(cmd []string) string {
	switch cmd[2] {
	case "ECHO", "echo":
		return "+" + cmd[4] + "\r\n"
	case "SET", "set":
		fmt.Println(cmd)
		d := Data{Value: cmd[6], CreatedAt: time.Now(), PX: -1}
		if len(cmd) > 6 && (cmd[7] == "px" || cmd[7] == "PX") {
			px, err := strconv.Atoi(cmd[9])
			if err == nil {
				d.PX = px
			}
		}
		db[cmd[4]] = d
		return "+OK\r\n"
	case "GET", "get":
		d := db[cmd[4]]
		if d.PX != -1 && d.CreatedAt.Add(time.Duration(d.PX)*time.Millisecond).UTC().Before(time.Now()) {
			return "+NULL\r\n"
		}
		return "+" + db[cmd[4]].Value + "\r\n"
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
