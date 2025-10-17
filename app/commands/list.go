package commands

import "strconv"

type RPushCommand struct {
	Key   string
	Value []string
}

type LPushCommand struct {
	Key   string
	Value []string
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
	command.Key = string(elements[1])
	command.Value = toStrings(elements[2:])
	return command, nil
}

func NewLPushCommand(elements [][]byte) (*LPushCommand, error) {
	command := &LPushCommand{}
	command.Key = string(elements[1])
	command.Value = toStrings(elements[2:])
	return command, nil
}

type LRangeCommand struct {
	Key   string
	Start int
	End   int
}

func NewLRangeCommand(elements [][]byte) (*LRangeCommand, error) {
	command := &LRangeCommand{}
	command.Key = string(elements[1])
	start, err := strconv.Atoi(string(elements[2]))
	if err != nil {
		return nil, err
	}
	command.Start = start
	end, err := strconv.Atoi(string(elements[3]))
	if err != nil {
		return nil, err
	}
	command.End = end
	return command, nil
}

type LLENCommand struct {
	Key string
}

func NewLLENCommand(elements [][]byte) (*LLENCommand, error) {
	command := &LLENCommand{}
	command.Key = string(elements[1])
	return command, nil
}
