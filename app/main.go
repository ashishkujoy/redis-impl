package main

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/server"
)

func main() {
	args := os.Args[1:]
	port := 6379
	if len(args) >= 1 {
		port, _ = strconv.Atoi(args[0])
	}
	add := fmt.Sprintf("0.0.0.0:%d", port)
	s := server.NewServer()
	l, err := net.Listen("tcp", add)
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}
		go func() {
			_ = s.Serve(conn)
		}()
	}
}
