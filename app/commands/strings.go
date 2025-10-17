package commands

import (
	"strconv"
	"strings"
)

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
