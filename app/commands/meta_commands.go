package commands

import "fmt"

type PingCommand struct{}

func (c *PingCommand) Execute(_ *ExecutionContext) ([]byte, error) {
	return []byte("+PONG\r\n"), nil
}

func NewPingCommand(_ [][]byte) (*PingCommand, error) {
	return &PingCommand{}, nil
}

type EchoCommand struct {
	Message string
}

func (e EchoCommand) Execute(ctx *ExecutionContext) ([]byte, error) {
	return ctx.Serializer.Encode(e.Message)
}

func NewEchoCommand(elements [][]byte) (*EchoCommand, error) {
	return &EchoCommand{
		Message: string(elements[0]),
	}, nil
}

type TypeCommand struct {
	Key string
}

func NewTypeCommand(elements [][]byte) (*TypeCommand, error) {
	if len(elements) != 1 {
		return nil, fmt.Errorf("expected 1 element, got %d", len(elements))
	}
	return &TypeCommand{
		Key: string(elements[0]),
	}, nil
}

func (t *TypeCommand) Execute(ctx *ExecutionContext) ([]byte, error) {
	isString := ctx.Kv.Contains(t.Key)
	isArray := ctx.Lists.Contains(t.Key)
	valueType := "none"
	if isString {
		valueType = "string"
	}
	if isArray {
		valueType = "array"
	}
	return ctx.Serializer.EncodeSimpleString(valueType)
}

func RegisterMetaCommands(registry *CommandRegistry) {
	registry.Register("ping", func(args [][]byte) (Command, error) {
		return NewPingCommand(args)
	})
	registry.Register("echo", func(args [][]byte) (Command, error) {
		return NewEchoCommand(args)
	})
	registry.Register("type", func(args [][]byte) (Command, error) {
		return NewTypeCommand(args)
	})
}
