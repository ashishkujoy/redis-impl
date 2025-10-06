package serializer

import (
	"bytes"
	"errors"
	"net"
)

type CommandName = string

const (
	PING CommandName = "ping"
	ECHO CommandName = "echo"
)

type Command struct {
	Name CommandName
	Args []string
}

func ParseCommand(c net.Conn) (*Command, error) {
	commandBuf := make([]byte, 4)
	read, err := c.Read(commandBuf)
	if err != nil {
		return nil, err
	}
	if read != 4 {
		return nil, errors.New("invalid command")
	}
	if bytes.Equal(commandBuf, []byte("ping")) {
		return &Command{Name: "ping"}, nil
	}
	if bytes.Equal(commandBuf, []byte("echo")) {
		argsBuf := make([]byte, 1024)
		read, err := c.Read(argsBuf)
		if err != nil {
			return nil, err
		}
		arg := string(argsBuf[:read])
		args := make([]string, 1)
		args[0] = arg
		return &Command{Name: "echo", Args: args}, nil
	}
	return nil, errors.New("invalid command")
}
