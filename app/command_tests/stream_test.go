package command_tests

import (
	"testing"

	"github.com/codecrafters-io/redis-starter-go/app/commands"
	"github.com/stretchr/testify/assert"
)

func TestXADD_CommandParsing(t *testing.T) {
	_, err := commands.NewXADDCommand([][]byte{[]byte("Key1"), []byte("1526919030473-0")})
	assert.NoError(t, err)
}

func TestXADD_CommandForNonExistingKey(t *testing.T) {
	command, err := commands.NewXADDCommand([][]byte{[]byte("Key1"), []byte("1526919030473-0")})
	assert.NoError(t, err)
	ctx := CreateExecutionContext()
	res, err := command.Execute(ctx)

	assert.NoError(t, err)
	assert.Equal(t, "$15\r\n1526919030473-0\r\n", string(res))
}

func TestXADD_CommandWithAutoGenerateSequence(t *testing.T) {
	command, err := commands.NewXADDCommand([][]byte{[]byte("Key1"), []byte("1526919030473-*")})
	assert.NoError(t, err)
	ctx := CreateExecutionContext()
	res, err := command.Execute(ctx)

	assert.NoError(t, err)
	assert.Equal(t, "$15\r\n1526919030473-0\r\n", string(res))
}
