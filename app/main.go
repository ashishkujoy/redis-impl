package main

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/commands"
	"github.com/codecrafters-io/redis-starter-go/app/parser"
	"github.com/codecrafters-io/redis-starter-go/app/serializer"
	"github.com/codecrafters-io/redis-starter-go/app/store"
	"github.com/codecrafters-io/redis-starter-go/app/store/ds"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")
	args := os.Args[1:]
	port := 6379
	if len(args) >= 1 {
		port, _ = strconv.Atoi(args[0])
	}
	kvStore := store.NewKVStore()
	lists := ds.NewLists()
	var respSerializer = serializer.NewRESPSerializer()
	executionContext := commands.NewExecutionContext(kvStore, lists, respSerializer)
	registry := commands.SetupCommandRegistry()

	add := fmt.Sprintf("0.0.0.0:%d", port)
	fmt.Println(add)
	l, err := net.Listen("tcp", add)
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
		go func(c net.Conn) {
			defer func(conn net.Conn) {
				_ = conn.Close()
			}(conn)
			for {
				bytes, err := parser.ParseArray(c)
				if err != nil {
					fmt.Println("Error parsing command: ", err.Error())
					return
				}
				command, err := registry.Create(string(bytes[0]), bytes[1:])
				if err != nil {
					fmt.Println("Error creating command: ", err.Error())
					return
				}

				res, err := command.Execute(executionContext)
				if err != nil {
					fmt.Println("Error executing command: ", err.Error())
					return
				}
				_, err = conn.Write(res)
				if err != nil {
					fmt.Println("Error writing response: ", err.Error())
					return
				}
			}
		}(conn)
	}
}
