package main

import (
	"fmt"
	"net"
	"os"
	"strconv"

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
	for _, arg := range args {
		fmt.Println(arg)
	}
	port := 6379
	if len(args) >= 1 {
		port, _ = strconv.Atoi(args[0])
	}
	kvStore := store.NewKVStore()
	lists := ds.NewLists()

	add := fmt.Sprintf("0.0.0.0:%d", port)
	fmt.Println(add)
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
				switch c := command.(type) {
				case *serializer.PingCommand:
					{
						fmt.Println("Received PING command")
						_, _ = conn.Write([]byte("+PONG\r\n"))
						fmt.Println("Replied for ping command")
					}
				case *serializer.EchoCommand:
					{
						fmt.Println("Received Echo command")
						bytes, err := serializer.EncodeBulkString(c.Message)
						if err != nil {
							fmt.Println("Error encoding command: ", err.Error())
						}
						_, err = conn.Write(bytes)
					}
				case *serializer.SetCommand:
					{
						fmt.Println("Received Set command")
						kvStore.Set(c.Key, c.Value, c.PX)
						_, _ = conn.Write([]byte("+OK\r\n"))
					}
				case *serializer.GetCommand:
					{
						fmt.Println("Received Get command")
						value, found := kvStore.Get(c.Key)
						res := []byte("$-1\r\n")
						if found {
							res, err = serializer.EncodeBulkString(value)
							if err != nil {
								fmt.Println("Error encoding value: ", err.Error())
								return
							}
						}
						_, _ = conn.Write(res)
					}
				case *serializer.RPushCommand:
					{
						fmt.Println("Received RPush command")
						length := lists.RPush(c.Key, c.Value)
						res, err := serializer.EncodeNumber(length)
						if err != nil {
							fmt.Println("Error encoding value: ", err.Error())
							return
						}
						_, _ = conn.Write(res)
					}
				case *serializer.LRangeCommand:
					{
						fmt.Println("Received LRange command")
						elements := lists.LRange(c.Key, c.Start, c.End)
						res, _ := serializer.EncodeAsBulkArray(elements)
						_, _ = conn.Write(res)
					}
				}
			}
		}(conn)
	}
}
