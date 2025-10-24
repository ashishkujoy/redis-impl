package ds

import "errors"

type Stream struct {
	Id        string
	Timestamp int
	Sequence  int
}

func NewStream(id string, timestamp int, sequence int) *Stream {
	return &Stream{Id: id, Timestamp: timestamp, Sequence: sequence}
}

type Streams struct {
	streams map[string]*Stream
}

func NewStreams() *Streams {
	return &Streams{
		streams: make(map[string]*Stream),
	}
}

func (s *Streams) validateKey(key string, timestamp, sequence int) error {
	if timestamp < 1 && sequence < 1 {
		return errors.New("ERR The ID specified in XADD must be greater than 0-0")
	}
	existingEntry, ok := s.streams[key]
	if !ok {
		return nil
	}
	if existingEntry.Timestamp <= timestamp && existingEntry.Sequence < sequence {
		return nil
	}
	return errors.New("ERR The ID specified in XADD is equal or smaller than the target stream top item")
}

func (s *Streams) Register(key string, timestamp int, sequence int) error {
	err := s.validateKey(key, timestamp, sequence)
	if err != nil {
		return err
	}
	stream := NewStream(key, timestamp, sequence)
	s.streams[key] = stream
	return nil
}

func (s *Streams) Contains(key string) bool {
	_, ok := s.streams[key]
	return ok
}
