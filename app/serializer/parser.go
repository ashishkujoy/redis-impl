package serializer

import (
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
	Args []byte
}

func ParseCommand(c net.Conn) (*Command, error) {
	commandBuf := make([]byte, 1024)
	read, err := c.Read(commandBuf)
	if err != nil {
		return nil, err
	}
	commandName := string(commandBuf[:4])
	return &Command{
		Name: commandName,
		Args: commandBuf[5:read],
	}, errors.New("invalid command")
}
