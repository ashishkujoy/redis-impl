package ds

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Stream struct {
	Id        string
	Timestamp int
	Sequence  int
}

func NewStream(id string, timestamp int, sequence int) *Stream {
	return &Stream{Id: id, Timestamp: timestamp, Sequence: sequence}
}

type Streams struct {
	streams map[string][]*Stream
}

func NewStreams() *Streams {
	return &Streams{
		streams: make(map[string][]*Stream),
	}
}

func (s *Streams) validateKey(key string, timestamp, sequence int) error {
	if timestamp < 1 && sequence < 1 {
		return errors.New("ERR The ID specified in XADD must be greater than 0-0")
	}
	existingEntry, ok := s.streams[key]
	if !ok || len(existingEntry) == 0 {
		return nil
	}
	lastEntry := existingEntry[len(existingEntry)-1]
	if lastEntry.Timestamp < timestamp {
		return nil
	}
	if lastEntry.Timestamp == timestamp && lastEntry.Sequence < sequence {
		return nil
	}
	return errors.New("ERR The ID specified in XADD is equal or smaller than the target stream top item")
}

func (s *Streams) generateSequence(key string, timestamp int) int {
	existingStream, ok := s.streams[key]
	if !ok {
		if timestamp == 0 {
			return 1
		}
		return 0
	}
	var timestampHead *Stream
	for _, stream := range existingStream {
		if stream.Timestamp == timestamp {
			timestampHead = stream
		}
	}
	if timestampHead == nil {
		return 0
	}

	return timestampHead.Sequence + 1
}

func (s *Streams) getTimestampAndSequence(key string, id string) (int, int, error) {
	tokens := strings.Split(id, "-")
	if len(tokens) != 2 {
		return 0, 0, errors.New("ERR The ID specified in XADD is invalid")
	}
	timestamp, err := strconv.Atoi(tokens[0])
	if err != nil {
		return 0, 0, errors.New("ERR The ID specified in XADD is invalid")
	}
	sequenceToken := tokens[1]
	if sequenceToken == "*" {
		sequence := s.generateSequence(key, timestamp)
		return timestamp, sequence, nil
	}
	sequence, err := strconv.Atoi(sequenceToken)
	if err != nil {
		return 0, 0, errors.New("ERR The ID specified in XADD is invalid")
	}
	return timestamp, sequence, nil
}

func (s *Streams) Register(key string, id string) (string, error) {
	timestamp, sequence, err := s.getTimestampAndSequence(key, id)
	if err != nil {
		return "", err
	}
	err = s.validateKey(key, timestamp, sequence)
	if err != nil {
		return "", err
	}
	stream := NewStream(key, timestamp, sequence)
	s.streams[key] = []*Stream{stream}
	return fmt.Sprintf("%d-%d", timestamp, sequence), nil
}

func (s *Streams) Contains(key string) bool {
	_, ok := s.streams[key]
	return ok
}
