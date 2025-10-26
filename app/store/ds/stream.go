package ds

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/store"
)

const (
	AutoGenerateIDToken             = "*"
	ErrInvalidStreamIDGreaterThanZero = "ERR The ID specified in XADD must be greater than 0-0"
	ErrInvalidStreamIDEqualOrSmaller  = "ERR The ID specified in XADD is equal or smaller than the target stream top item"
	ErrInvalidStreamID                = "ERR The ID specified in XADD is invalid"
)

type StreamEntry struct {
	Timestamp int
	Sequence  int
}

type StreamID struct {
	Timestamp int
	Sequence  int
}

func NewStreamEntry(timestamp int, sequence int) *StreamEntry {
	return &StreamEntry{Timestamp: timestamp, Sequence: sequence}
}

type Stream struct {
	entries []*StreamEntry
}

func (s *Stream) lastEntry() *StreamEntry {
	if len(s.entries) == 0 {
		return nil
	}
	return s.entries[len(s.entries)-1]
}

type Streams struct {
	streams map[string]*Stream
	clock   store.Clock
}

func NewStreams() *Streams {
	return &Streams{
		streams: make(map[string]*Stream),
	}
}

func NewStreamsWithClock(clock store.Clock) *Streams {
	return &Streams{
		streams: make(map[string]*Stream),
		clock:   clock,
	}
}

func (s *Streams) validateEntryID(key string, timestamp, sequence int) error {
	if timestamp < 1 && sequence < 1 {
		return errors.New(ErrInvalidStreamIDGreaterThanZero)
	}
	existingStream, ok := s.streams[key]
	if !ok || len(existingStream.entries) == 0 {
		return nil
	}
	lastEntry := existingStream.lastEntry()
	if lastEntry.Timestamp < timestamp {
		return nil
	}
	if lastEntry.Timestamp == timestamp && lastEntry.Sequence < sequence {
		return nil
	}
	return errors.New(ErrInvalidStreamIDEqualOrSmaller)
}

func (s *Streams) generateSequence(key string, timestamp int) int {
	existingStream, ok := s.streams[key]
	if !ok || len(existingStream.entries) == 0 {
		if timestamp == 0 {
			return 1
		}
		return 0
	}
	lastEntry := existingStream.lastEntry()
	if lastEntry.Timestamp == timestamp {
		return lastEntry.Sequence + 1
	}
	if lastEntry.Timestamp < timestamp {
		return 0
	}
	if timestamp == 0 {
		return 1
	}
	return 0
}

func (s *Streams) generateTimestamp(key string, timestampToken string) int {
	timestamp, err := strconv.Atoi(timestampToken)
	if err == nil {
		return timestamp
	}
	stream, ok := s.streams[key]
	defaultTimestamp := s.clock.CurrentMillis()
	if !ok {
		return defaultTimestamp
	}
	if len(stream.entries) == 0 {
		return defaultTimestamp
	}
	lastEntry := stream.lastEntry()
	if lastEntry.Timestamp >= defaultTimestamp {
		return lastEntry.Timestamp + 1
	}
	return defaultTimestamp
}

func parseStreamID(id string) (string, string, error) {
	if id == AutoGenerateIDToken {
		return AutoGenerateIDToken, AutoGenerateIDToken, nil
	}
	tokens := strings.Split(id, "-")
	if len(tokens) == 1 {
		return tokens[0], AutoGenerateIDToken, nil
	}
	if len(tokens) == 2 {
		return tokens[0], tokens[1], nil
	}
	return "", "", errors.New(ErrInvalidStreamID)
}

func (s *Streams) generateStreamID(key string, id string) (*StreamID, error) {
	timestampToken, sequenceToken, err := parseStreamID(id)
	if err != nil {
		return nil, err
	}

	timestamp := s.generateTimestamp(key, timestampToken)

	var sequence int
	if sequenceToken == AutoGenerateIDToken {
		sequence = s.generateSequence(key, timestamp)
	} else {
		sequence, err = strconv.Atoi(sequenceToken)
		if err != nil {
			return nil, errors.New(ErrInvalidStreamID)
		}
	}

	return &StreamID{Timestamp: timestamp, Sequence: sequence}, nil
}

func (s *Streams) Add(key string, id string) (string, error) {
	streamID, err := s.generateStreamID(key, id)
	if err != nil {
		return "", err
	}
	err = s.validateEntryID(key, streamID.Timestamp, streamID.Sequence)
	if err != nil {
		return "", err
	}

	if _, ok := s.streams[key]; !ok {
		s.streams[key] = &Stream{entries: make([]*StreamEntry, 0)}
	}

	entry := NewStreamEntry(streamID.Timestamp, streamID.Sequence)
	s.streams[key].entries = append(s.streams[key].entries, entry)
	return fmt.Sprintf("%d-%d", streamID.Timestamp, streamID.Sequence), nil
}

func (s *Streams) Contains(key string) bool {
	_, ok := s.streams[key]
	return ok
}
