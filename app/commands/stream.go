package commands

import (
	"errors"
	"strconv"
	"strings"
)

type XADDCommand struct {
	Key       string
	Id        string
	Timestamp int
	Sequence  int
}

func extractIdParts(id string) (int, int, error) {
	tokens := strings.Split(id, "-")
	if len(tokens) != 2 {
		return 0, 0, errors.New("invalid id")
	}
	timestampToken := tokens[0]
	sequenceToken := tokens[1]
	timestamp, err := strconv.Atoi(timestampToken)
	if err != nil {
		return 0, 0, errors.New("invalid id")
	}
	sequence, err := strconv.Atoi(sequenceToken)
	if err != nil {
		return 0, 0, errors.New("invalid id")
	}
	return timestamp, sequence, nil
}

func (x *XADDCommand) Execute(ctx *ExecutionContext) ([]byte, error) {
	err := ctx.Streams.Register(x.Key, x.Timestamp, x.Sequence)
	if err != nil {
		return ctx.Serializer.EncodeError(err.Error()), nil
	}
	return ctx.Serializer.Encode(x.Id)
}

func NewXADDCommand(elements [][]byte) (*XADDCommand, error) {
	if len(elements) < 2 {
		return nil, errors.New("not enough arguments")
	}
	timestamp, sequence, err := extractIdParts(string(elements[1]))
	if err != nil {
		return nil, err
	}

	return &XADDCommand{
		Key:       string(elements[0]),
		Id:        string(elements[1]),
		Timestamp: timestamp,
		Sequence:  sequence,
	}, nil
}

func RegisterStreamCommands(registry *CommandRegistry) {
	registry.Register("xadd", func(i [][]byte) (Command, error) {
		return NewXADDCommand(i)
	})
}
