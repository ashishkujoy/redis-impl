package server

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/commands"
	"github.com/codecrafters-io/redis-starter-go/app/parser"
	"github.com/codecrafters-io/redis-starter-go/app/serializer"
	"github.com/codecrafters-io/redis-starter-go/app/store"
	"github.com/codecrafters-io/redis-starter-go/app/store/ds"
)

type Server struct {
	KV               *store.KVStore
	Lists            *ds.Lists
	serializer       *commands.Serializer
	executionContext *commands.ExecutionContext
	registry         *commands.CommandRegistry
}

func NewServer() *Server {
	kvStore := store.NewKVStore()
	lists := ds.NewLists()
	var respSerializer commands.Serializer = serializer.NewRESPSerializer()
	executionContext := commands.NewExecutionContext(kvStore, lists, respSerializer)
	registry := commands.SetupCommandRegistry()

	return &Server{
		KV:               kvStore,
		Lists:            lists,
		registry:         registry,
		serializer:       &respSerializer,
		executionContext: executionContext,
	}
}

func (s *Server) Serve(c net.Conn) error {
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(c)

	for {
		bytes, err := parser.ParseArray(c)
		if err != nil {
			return err
		}
		command, err := s.registry.Create(string(bytes[0]), bytes[1:])
		if err != nil {
			return err
		}

		res, err := command.Execute(s.executionContext)
		if err != nil {
			return err
		}
		_, err = c.Write(res)
		if err != nil {
			return err
		}
	}
}
