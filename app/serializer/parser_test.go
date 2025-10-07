package serializer

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadToken(t *testing.T) {
	token, cursor, err := readToken([]byte("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"), 1)
	assert.NoError(t, err)
	assert.Equal(t, 4, cursor)
	assert.Equal(t, []byte("2"), token)

	token, cursor, err = readToken([]byte("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"), cursor)
	assert.NoError(t, err)
	assert.Equal(t, 8, cursor)
	assert.Equal(t, []byte("$5"), token)
}

func TestReadNextElement(t *testing.T) {
	element, cursor, err := readNextElement([]byte("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"), 4)
	assert.NoError(t, err)
	assert.Equal(t, 15, cursor)
	assert.Equal(t, "hello", element)
}

func TestReadNextElementForEnd(t *testing.T) {
	_, _, err := readNextElement([]byte("*2\r\n$5\r\nhello\r\n"), 15)

	assert.Error(t, err)
	assert.Equal(t, "EOF", err.Error())
}

func TestEmptyArrayParsing(t *testing.T) {
	reader := strings.NewReader("*0\r\n")
	array, err := ParseArray(reader)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(array))
}

func TestArrayParsingSingleElementArray(t *testing.T) {
	reader := strings.NewReader("*0\r\n$5\r\nhello\r\n")
	array, err := ParseArray(reader)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(array))
	assert.Equal(t, "hello", array[0])
}

func TestArrayParsingMultiElementArray(t *testing.T) {
	reader := strings.NewReader("*3\r\n$5\r\nhello\r\n$5\r\nworld\r\n$9\r\nsomething\r\n")
	array, err := ParseArray(reader)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(array))
	assert.Equal(t, "hello", array[0])
	assert.Equal(t, "world", array[1])
	assert.Equal(t, "something", array[2])
}
