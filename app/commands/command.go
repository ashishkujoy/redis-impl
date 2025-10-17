package commands

import (
	"github.com/codecrafters-io/redis-starter-go/app/store"
	"github.com/codecrafters-io/redis-starter-go/app/store/ds"
)

type Serializer interface {
	Encode(interface{}) ([]byte, error)
	Decode([]byte) (Command, error)
}

type ExecutionContext struct {
	Kv         *store.KVStore
	Lists      *ds.Lists
	Serializer Serializer
}

func NewExecutionContext(kv *store.KVStore, lists *ds.Lists, serializer Serializer) *ExecutionContext {
	return &ExecutionContext{
		Kv:         kv,
		Lists:      lists,
		Serializer: serializer,
	}
}

type Command interface {
	Execute(ctx *ExecutionContext) ([]byte, error)
}
