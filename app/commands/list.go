package commands

import (
	"errors"
	"strconv"
)

type RPushCommand struct {
	Key   string
	Value []string
}

func (R RPushCommand) Execute(ctx *ExecutionContext) ([]byte, error) {
	length := ctx.Lists.RPush(R.Key, R.Value)
	return ctx.Serializer.Encode(length)
}

type LPushCommand struct {
	Key   string
	Value []string
}

func (L LPushCommand) Execute(ctx *ExecutionContext) ([]byte, error) {
	length := ctx.Lists.LPush(L.Key, L.Value)
	return ctx.Serializer.Encode(length)
}

func toStrings(elements [][]byte) []string {
	strs := make([]string, len(elements))
	for i, e := range elements {
		strs[i] = string(e)
	}
	return strs
}

func NewRPushCommand(elements [][]byte) (*RPushCommand, error) {
	command := &RPushCommand{}
	command.Key = string(elements[0])
	command.Value = toStrings(elements[1:])
	return command, nil
}

func NewLPushCommand(elements [][]byte) (*LPushCommand, error) {
	command := &LPushCommand{}
	command.Key = string(elements[0])
	command.Value = toStrings(elements[1:])
	return command, nil
}

type LRangeCommand struct {
	Key   string
	Start int
	End   int
}

func (L LRangeCommand) Execute(ctx *ExecutionContext) ([]byte, error) {
	values := ctx.Lists.LRange(L.Key, L.Start, L.End)
	return ctx.Serializer.Encode(values)
}

func NewLRangeCommand(elements [][]byte) (*LRangeCommand, error) {
	command := &LRangeCommand{}
	command.Key = string(elements[0])
	start, err := strconv.Atoi(string(elements[1]))
	if err != nil {
		return nil, err
	}
	command.Start = start
	end, err := strconv.Atoi(string(elements[2]))
	if err != nil {
		return nil, err
	}
	command.End = end
	return command, nil
}

type LLENCommand struct {
	Key string
}

func (L LLENCommand) Execute(ctx *ExecutionContext) ([]byte, error) {
	length := ctx.Lists.LLen(L.Key)
	return ctx.Serializer.Encode(length)
}

func NewLLENCommand(elements [][]byte) (*LLENCommand, error) {
	command := &LLENCommand{}
	command.Key = string(elements[0])
	return command, nil
}

type LPopCommand struct {
	Key string
}

func (L *LPopCommand) Execute(ctx *ExecutionContext) ([]byte, error) {
	element, err := ctx.Lists.LPop(L.Key)
	res := ctx.Serializer.NullBulkByte()
	if err == nil {
		res, _ = ctx.Serializer.Encode(element)
	}

	return res, nil
}

func NewLPopCommand(element [][]byte) (*LPopCommand, error) {
	if len(element) != 1 {
		return nil, errors.New("not enough arguments for LPop command")
	}
	command := &LPopCommand{}
	command.Key = string(element[0])
	return command, nil
}

func RegisterListCommands(registry *CommandRegistry) *CommandRegistry {
	registry.Register("rpush", func(elements [][]byte) (Command, error) {
		return NewRPushCommand(elements)
	})
	registry.Register("lpush", func(elements [][]byte) (Command, error) {
		return NewLPushCommand(elements)
	})
	registry.Register("lrange", func(elements [][]byte) (Command, error) {
		return NewLRangeCommand(elements)
	})
	registry.Register("llen", func(elements [][]byte) (Command, error) {
		return NewLLENCommand(elements)
	})
	registry.Register("lpop", func(elements [][]byte) (Command, error) {
		return NewLPopCommand(elements)
	})
	return registry
}
