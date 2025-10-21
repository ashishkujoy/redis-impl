package serializer

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/commands"
)

type RESPSerializer struct {
}

func (r RESPSerializer) NullArray() []byte {
	return []byte("*-1\r\n")
}

func (r RESPSerializer) EncodeBulkString(msg string) ([]byte, error) {
	return EncodeBytesAsBulkString([]byte(msg))
}

func NewRESPSerializer() RESPSerializer {
	return RESPSerializer{}
}

func (r RESPSerializer) Encode(i interface{}) ([]byte, error) {
	switch c := i.(type) {
	case int:
		{
			return EncodeNumber(c)
		}
	case []byte:
		{
			return EncodeBytesAsBulkString(c)
		}
	case string:
		{
			return EncodeBulkString(c)
		}
	case []string:
		{
			return EncodeAsBulkArray(c)
		}
	}
	return nil, nil
}

func (r RESPSerializer) Decode(bytes []byte) (commands.Command, error) {
	panic("implement me")
}

func (r RESPSerializer) NullBulkByte() []byte {
	return []byte("$-1\r\n")
}

func EncodeBulkString(message string) ([]byte, error) {
	return EncodeBytesAsBulkString([]byte(message))
}

func EncodeBytesAsBulkString(message []byte) ([]byte, error) {
	length := len(message)
	lenStr := strconv.Itoa(length)
	encoded := make([]byte, 0, 1+len(lenStr)+2+len(message)+2)
	encoded = append(encoded, '$')
	encoded = append(encoded, lenStr...)
	encoded = append(encoded, '\r', '\n')
	encoded = append(encoded, message...)
	encoded = append(encoded, '\r', '\n')

	return encoded, nil
}

func EncodeNumber(number int) ([]byte, error) {
	str := strconv.Itoa(number)
	encoded := make([]byte, 0, len(str)+3)
	encoded = append(encoded, ':')
	encoded = append(encoded, str...)
	encoded = append(encoded, '\r', '\n')

	return encoded, nil
}

func EncodeAsBulkArray(message []string) ([]byte, error) {
	str := strconv.Itoa(len(message))
	encoded := make([]byte, 0, len(message)+len(str))
	encoded = append(encoded, '*')
	encoded = append(encoded, str...)
	encoded = append(encoded, '\r', '\n')

	for _, element := range message {
		length := strconv.Itoa(len(element))
		encoded = append(encoded, '$')
		encoded = append(encoded, length...)
		encoded = append(encoded, '\r', '\n')
		encoded = append(encoded, element...)
		encoded = append(encoded, '\r', '\n')
	}

	return encoded, nil
}
