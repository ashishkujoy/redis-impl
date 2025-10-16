package commands

type PingCommand struct{}

type EchoCommand struct {
	Message string
}

func NewEchoCommand(elements [][]byte) (*EchoCommand, error) {
	return &EchoCommand{
		Message: string(elements[1]),
	}, nil
}
