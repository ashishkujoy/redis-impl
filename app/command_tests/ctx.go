package command_tests

import (
	"github.com/codecrafters-io/redis-starter-go/app/commands"
	"github.com/codecrafters-io/redis-starter-go/app/serializer"
	"github.com/codecrafters-io/redis-starter-go/app/store"
	"github.com/codecrafters-io/redis-starter-go/app/store/ds"
)

func CreateExecutionContext() *commands.ExecutionContext {
	blockingQueueManager := ds.NewBlockingQueueManager()
	return &commands.ExecutionContext{
		Kv:                   store.NewKVStore(),
		Lists:                ds.NewLists(blockingQueueManager),
		Streams:              ds.NewStreams(),
		Serializer:           serializer.NewRESPSerializer(),
		BlockingQueueManager: blockingQueueManager,
	}
}
