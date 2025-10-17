package commands

import "fmt"

type CommandFactory func([][]byte) (Command, error)

type CommandRegistry struct {
	factories map[string]CommandFactory
}

func (r *CommandRegistry) Register(name string, factory CommandFactory) {
	r.factories[name] = factory
}

func (r *CommandRegistry) Create(name string, args [][]byte) (Command, error) {
	factory, ok := r.factories[name]
	if !ok {
		return nil, fmt.Errorf("no such command: %s", name)
	}
	return factory(args)
}

func SetupCommandRegistry() *CommandRegistry {
	factories := make(map[string]CommandFactory)
	r := &CommandRegistry{factories: factories}

	r.Register("echo", func(i [][]byte) (Command, error) {
		return NewEchoCommand(i)
	})
	r.Register("ping", func(i [][]byte) (Command, error) {
		return NewPingCommand(i)
	})

	return r
}
