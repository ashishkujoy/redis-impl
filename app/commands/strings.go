package commands

import (
	"strconv"
	"strings"
)

type GetCommand struct {
	Key string
}

func (g *GetCommand) Execute(ctx *ExecutionContext) ([]byte, error) {
	value, found := ctx.Kv.Get(g.Key)
	if !found {
		return []byte("$-1\r\n"), nil
	}
	return ctx.Serializer.Encode(value)
}

func NewGetCommand(elements [][]byte) (*GetCommand, error) {
	return &GetCommand{
		Key: string(elements[0]),
	}, nil
}

type SetCommand struct {
	Key   string
	Value string
	PX    int
}

func (s *SetCommand) Execute(ctx *ExecutionContext) ([]byte, error) {
	ctx.Kv.Set(s.Key, s.Value, s.PX)
	return []byte("+OK\r\n"), nil
}

func NewSetCommand(elements [][]byte) (*SetCommand, error) {
	command := &SetCommand{}
	command.Key = string(elements[0])
	command.Value = string(elements[1])
	command.PX = -1
	if len(elements) == 4 {
		timeOptionName := string(elements[2])
		PX, err := strconv.Atoi(string(elements[3]))
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

func RegisterKVCommands(registry *CommandRegistry) *CommandRegistry {
	registry.Register("get", func(args [][]byte) (Command, error) {
		return NewGetCommand(args)
	})
	registry.Register("set", func(args [][]byte) (Command, error) {
		return NewSetCommand(args)
	})

	return registry
}
