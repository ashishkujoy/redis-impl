package serializer

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
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

func readNextElement(buf []byte, cursor int) (string, int, error) {
	token, cursor, err := readToken(buf, cursor+1)
	if err != nil {
		return "", cursor, err
	}
	elementSize, err := toNumber(token)
	if err != nil {
		return "", cursor, err
	}
	elem := string(buf[cursor : cursor+elementSize])
	return elem, cursor + elementSize + 2, nil
}

func ParseArray(c io.Reader) ([]string, error) {
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
	elements := make([]string, 0, elemCount)
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

func ParseCommand(c net.Conn) (*Command, error) {
	elements, err := ParseArray(c)
	for _, element := range elements {
		fmt.Println(element)
	}
	if err != nil {
		return nil, err
	}

	return &Command{
		Name: elements[0],
		Args: elements[1:],
	}, errors.New("invalid command")
}
