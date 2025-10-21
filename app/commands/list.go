package commands

import (
	"context"
	"errors"
	"strconv"
	"time"
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
	Key   string
	Count int
}

func (L *LPopCommand) Execute(ctx *ExecutionContext) ([]byte, error) {
	elements := ctx.Lists.LPop(L.Key, L.Count)

	if len(elements) == 0 {
		return ctx.Serializer.NullBulkByte(), nil
	}
	if len(elements) == 1 {
		return ctx.Serializer.EncodeBulkString(elements[0])
	}

	return ctx.Serializer.Encode(elements)
}

func NewLPopCommand(element [][]byte) (*LPopCommand, error) {
	if len(element) == 0 {
		return nil, errors.New("not enough arguments for LPop command")
	}
	command := &LPopCommand{}
	command.Key = string(element[0])
	command.Count = 1
	if len(element) > 1 {
		count, err := strconv.Atoi(string(element[1]))
		if err != nil {
			return nil, err
		}
		command.Count = count
	}
	return command, nil
}

type BLPopCommand struct {
	Key     string
	Timeout int
}

func NewBLPopCommand(element [][]byte) (*BLPopCommand, error) {
	if len(element) == 0 {
		return nil, errors.New("not enough arguments for BLPop command")
	}
	timeout := 0
	if len(element) > 1 {
		count, err := strconv.Atoi(string(element[1]))
		if err != nil {
			return nil, err
		}
		timeout = count
	}
	return &BLPopCommand{
		Key:     string(element[0]),
		Timeout: timeout,
	}, nil
}

func (c *BLPopCommand) Execute(ctx *ExecutionContext) ([]byte, error) {
	values := ctx.Lists.LPop(c.Key, 1)

	if len(values) == 0 {
		value := ""
		if c.Timeout == 0 {
			value = c.blockIndefinitely(ctx)
		} else {
			v, err := c.blockWithTimeout(ctx)
			if err != nil {
				return ctx.Serializer.NullArray(), nil
			}
			value = v
		}
		values = append(values, c.Key)
		values = append(values, value)

	}
	return ctx.Serializer.Encode(values)
}

func (c *BLPopCommand) blockIndefinitely(ctx *ExecutionContext) string {
	ct := context.Background()
	blockedClient := ctx.BlockingQueueManager.BlockOn(c.Key, ct)
	return <-blockedClient.WakeChan
}

func (c *BLPopCommand) blockWithTimeout(ctx *ExecutionContext) (string, error) {
	ct, cancel := context.WithTimeout(context.Background(), time.Duration(c.Timeout)*time.Second)
	defer cancel()
	blockedClient := ctx.BlockingQueueManager.BlockOn(c.Key, ct)
	for {
		select {
		case value := <-blockedClient.WakeChan:
			return value, nil
		case <-blockedClient.TimeoutChan:
			return "", errors.New("timeout")
		}
	}
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
	registry.Register("blpop", func(elements [][]byte) (Command, error) {
		return NewBLPopCommand(elements)
	})
	return registry
}
