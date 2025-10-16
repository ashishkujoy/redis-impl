package commands

type Command interface{}

type CommandFactory = func([][]byte) (Command, error)

type CommandFactories = map[string]CommandFactory
