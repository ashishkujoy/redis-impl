package main

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/serializer"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")
	args := os.Args[1:]
	port := 6379
	if len(args) > 1 {
		port, _ = strconv.Atoi(args[1])
	}

	// Uncomment this block to pass the first stage
	//
	add := fmt.Sprintf("0.0.0.0:%d", port)
	l, err := net.Listen("tcp", add)
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	for {
		conn, err := l.Accept()
		fmt.Println("Accepting connection")
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go func(c net.Conn) {
			defer func(conn net.Conn) {
				_ = conn.Close()
			}(conn)
			fmt.Println("Connection accepted")
			for {
				command, err := serializer.ParseCommand(c)
				if err != nil {
					fmt.Println("Error parsing command: ", err.Error())
					return
				}
				fmt.Printf("Command: %s, args count%d\n", command.Name, len(command.Args))
				if command.Name == serializer.PING {
					fmt.Println("Received command")
					_, _ = conn.Write([]byte("+PONG\r\n"))
					fmt.Println("Replied for ping command")
					continue
				}
				fmt.Println("Received command: ", command.Name)
				fmt.Println(len(command.Args))
				bulkString, err := serializer.EncodeBulkString(command.Args[0])
				fmt.Println("Encoded BulkString: ", bulkString)
				if err != nil {
					return
				}
				_, err = conn.Write(bulkString)
				if err != nil {
					return
				}
			}
		}(conn)
	}
}
