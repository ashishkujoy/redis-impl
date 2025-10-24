package commands

import (
	"github.com/codecrafters-io/redis-starter-go/app/store"
	"github.com/codecrafters-io/redis-starter-go/app/store/ds"
)

type Serializer interface {
	Encode(interface{}) ([]byte, error)
	EncodeBulkString(string) ([]byte, error)
	EncodeSimpleString(string) ([]byte, error)
	Decode([]byte) (Command, error)
	NullBulkByte() []byte
	NullArray() []byte
}

type ExecutionContext struct {
	Kv                   *store.KVStore
	Lists                *ds.Lists
	Serializer           Serializer
	BlockingQueueManager *ds.BlockingQueueManager
}

func NewExecutionContext(
	kv *store.KVStore,
	lists *ds.Lists,
	serializer Serializer,
	blockingQueueManager *ds.BlockingQueueManager,
) *ExecutionContext {
	return &ExecutionContext{
		Kv:                   kv,
		Lists:                lists,
		Serializer:           serializer,
		BlockingQueueManager: blockingQueueManager,
	}
}

type Command interface {
	Execute(ctx *ExecutionContext) ([]byte, error)
}
