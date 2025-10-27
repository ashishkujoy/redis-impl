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
	EncodeError(string) []byte
	EncodeXRange(entries []*ds.StreamEntryView) ([]byte, error)
	EncodeXRead(entries []*ds.StreamView) ([]byte, error)
}

type ExecutionContext struct {
	Kv                   *store.KVStore
	Lists                *ds.Lists
	Streams              *ds.Streams
	Serializer           Serializer
	BlockingQueueManager *ds.BlockingQueueManager
}

func NewExecutionContext(
	kv *store.KVStore,
	lists *ds.Lists,
	streams *ds.Streams,
	serializer Serializer,
	blockingQueueManager *ds.BlockingQueueManager,
) *ExecutionContext {
	return &ExecutionContext{
		Kv:                   kv,
		Lists:                lists,
		Streams:              streams,
		Serializer:           serializer,
		BlockingQueueManager: blockingQueueManager,
	}
}

type Command interface {
	Execute(ctx *ExecutionContext) ([]byte, error)
}
