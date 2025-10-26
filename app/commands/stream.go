package commands

import (
	"errors"
)

type XADDCommand struct {
	Key  string
	Id   string
	Data [][]byte
}

func (x *XADDCommand) Execute(ctx *ExecutionContext) ([]byte, error) {
	id, err := ctx.Streams.Add(x.Key, x.Id, x.Data)
	if err != nil {
		return ctx.Serializer.EncodeError(err.Error()), nil
	}
	return ctx.Serializer.Encode(id)
}

func NewXADDCommand(elements [][]byte) (*XADDCommand, error) {
	if len(elements) < 2 {
		return nil, errors.New("not enough arguments")
	}

	return &XADDCommand{
		Key:  string(elements[0]),
		Id:   string(elements[1]),
		Data: elements[2:],
	}, nil
}

type XRANGECommand struct {
	Key   string
	Start string
	End   string
}

func (x *XRANGECommand) Execute(ctx *ExecutionContext) ([]byte, error) {
	entries := ctx.Streams.List(x.Key, x.Start, x.End)
	return ctx.Serializer.EncodeXRange(entries)
}

func NewXRANGECommand(elements [][]byte) (*XRANGECommand, error) {
	if len(elements) < 2 {
		return nil, errors.New("not enough arguments")
	}
	return &XRANGECommand{
		Key:   string(elements[0]),
		Start: string(elements[1]),
		End:   string(elements[2]),
	}, nil
}

func RegisterStreamCommands(registry *CommandRegistry) {
	registry.Register("xadd", func(i [][]byte) (Command, error) {
		return NewXADDCommand(i)
	})
	registry.Register("xrange", func(i [][]byte) (Command, error) {
		return NewXRANGECommand(i)
	})
}
