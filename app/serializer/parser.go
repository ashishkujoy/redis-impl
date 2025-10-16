package serializer

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

type CommandName = string

const (
	PING CommandName = "ping"
	ECHO CommandName = "echo"
	GET  CommandName = "get"
	SET  CommandName = "set"
)

type Command interface{}

type PingCommand struct{}

type EchoCommand struct {
	Message string
}

func NewEchoCommand(elements [][]byte) (*EchoCommand, error) {
	return &EchoCommand{
		Message: string(elements[1]),
	}, nil
}

type GetCommand struct {
	Key string
}

func NewGetCommand(elements [][]byte) (*GetCommand, error) {
	return &GetCommand{
		Key: string(elements[1]),
	}, nil
}

type SetCommand struct {
	Key   string
	Value string
	PX    int
}

func NewSetCommand(elements [][]byte) (*SetCommand, error) {
	command := &SetCommand{}
	command.Key = string(elements[1])
	command.Value = string(elements[2])
	command.PX = -1
	if len(elements) == 5 {
		timeOptionName := string(elements[3])
		PX, err := strconv.Atoi(string(elements[4]))
		if err != nil {
			return nil, err
		}
		if strings.ToLower(timeOptionName) == "ex" {
			PX = 1000 * PX
		}
		command.PX = PX
	}
	return command, nil
}

type RPushCommand struct {
	Key   string
	Value string
}

func NewRPushCommand(elements [][]byte) (*RPushCommand, error) {
	command := &RPushCommand{}
	command.Key = string(elements[1])
	command.Value = string(elements[2])
	return command, nil
}

var EofError = errors.New("EOF")

func readToken(buf []byte, cursor int) ([]byte, int, error) {
	currentPosition := cursor
	endFound := false
	for ; currentPosition < len(buf)-1; currentPosition++ {
		if buf[currentPosition] == byte('\r') && buf[currentPosition+1] == byte('\n') {
			endFound = true
			break
		}
	}

	if !endFound {
		return nil, len(buf), EofError
	}

	return buf[cursor:currentPosition], currentPosition + 2, nil
}

func toNumber(n []byte) (int, error) {
	return strconv.Atoi(string(n))
}

func readNextElement(buf []byte, cursor int) ([]byte, int, error) {
	token, cursor, err := readToken(buf, cursor+1)
	if err != nil {
		return nil, cursor, err
	}
	elementSize, err := toNumber(token)
	if err != nil {
		return nil, cursor, err
	}
	elem := buf[cursor : cursor+elementSize]
	return elem, cursor + elementSize + 2, nil
}

func ParseArray(c io.Reader) ([][]byte, error) {
	buf := make([]byte, 1024)
	length, err := c.Read(buf)
	if err != nil {
		return nil, err
	}
	buf = buf[:length]
	cursor := 1
	elemCountByte, cursor, err := readToken(buf, cursor)
	if err != nil {
		return nil, err
	}
	elemCount, err := toNumber(elemCountByte)
	elements := make([][]byte, 0, elemCount)
	for cursor < length {
		element, nextCursor, err := readNextElement(buf, cursor)
		if errors.Is(err, EofError) {
			break
		}
		cursor = nextCursor
		elements = append(elements, element)
	}
	return elements, err
}

func ParseCommand(c net.Conn) (Command, error) {
	elements, err := ParseArray(c)
	if err != nil {
		return nil, err
	}

	switch strings.ToLower(string(elements[0])) {
	case "ping":
		return &PingCommand{}, nil
	case "echo":
		return NewEchoCommand(elements)
	case "get":
		return NewGetCommand(elements)
	case "set":
		return NewSetCommand(elements)
	case "rpush":
		return NewRPushCommand(elements)
	}
	return nil, fmt.Errorf("unknown command: %s", string(elements[0]))
}
