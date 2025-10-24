package commands

import "errors"

type XADDCommand struct {
	Key string
	Id  string
}

func (x *XADDCommand) Execute(ctx *ExecutionContext) ([]byte, error) {
	ctx.Streams.Register(x.Key)
	return ctx.Serializer.Encode(x.Id)
}

func NewXADDCommand(elements [][]byte) (*XADDCommand, error) {
	if len(elements) < 2 {
		return nil, errors.New("not enough arguments")
	}
	return &XADDCommand{
		Key: string(elements[0]),
		Id:  string(elements[1]),
	}, nil
}

func RegisterStreamCommands(registry *CommandRegistry) {
	registry.Register("xadd", func(i [][]byte) (Command, error) {
		return NewXADDCommand(i)
	})
}
