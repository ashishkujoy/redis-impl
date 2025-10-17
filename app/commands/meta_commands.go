package commands

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

func RegisterMetaCommands(registry *CommandRegistry) {
	registry.Register("ping", func(args [][]byte) (Command, error) {
		return NewPingCommand(args)
	})
	registry.Register("echo", func(args [][]byte) (Command, error) {
		return NewEchoCommand(args)
	})
}
