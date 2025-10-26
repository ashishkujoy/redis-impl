package ds

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/store"
)

const (
	AutoGenerateIDToken               = "*"
	ErrInvalidStreamIDGreaterThanZero = "ERR The ID specified in XADD must be greater than 0-0"
	ErrInvalidStreamIDEqualOrSmaller  = "ERR The ID specified in XADD is equal or smaller than the target stream top item"
	ErrInvalidStreamID                = "ERR The ID specified in XADD is invalid"
)

type StreamEntry struct {
	Timestamp int
	Sequence  int
	Data      [][]byte
}

type StreamID struct {
	Timestamp int
	Sequence  int
}

func NewStreamID(timestamp, sequence int) *StreamID {
	return &StreamID{
		Timestamp: timestamp,
		Sequence:  sequence,
	}
}

func (id *StreamID) isInRange(start, end *StreamID) bool {
	isTimestampInRange := start.Timestamp <= id.Timestamp && id.Timestamp <= end.Timestamp
	if !isTimestampInRange {
		return false
	}
	if id.Timestamp == start.Timestamp {
		return id.Sequence >= start.Sequence
	}
	if id.Timestamp == end.Timestamp {
		return id.Sequence <= end.Sequence
	}
	return true
}

func NewStreamEntry(timestamp int, sequence int, data [][]byte) *StreamEntry {
	return &StreamEntry{Timestamp: timestamp, Sequence: sequence, Data: data}
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
	if lastEntry == nil {
		return nil
	}
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
	if lastEntry == nil {
		if timestamp == 0 {
			return 1
		}
		return 0
	}
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

func (s *Streams) resolveTimestamp(key string, timestampToken string) int {
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
	if lastEntry == nil || lastEntry.Timestamp <= defaultTimestamp {
		return defaultTimestamp
	}
	return defaultTimestamp + 1
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

func (s *Streams) resolveSequence(key string, timestamp int, sequenceToken string) (int, error) {
	if sequenceToken == AutoGenerateIDToken {
		return s.generateSequence(key, timestamp), nil
	}
	sequence, err := strconv.Atoi(sequenceToken)
	if err != nil {
		return 0, errors.New(ErrInvalidStreamID)
	}
	return sequence, nil
}

func (s *Streams) generateStreamID(key string, id string) (*StreamID, error) {
	timestampToken, sequenceToken, err := parseStreamID(id)
	if err != nil {
		return nil, err
	}

	timestamp := s.resolveTimestamp(key, timestampToken)
	sequence, err := s.resolveSequence(key, timestamp, sequenceToken)
	if err != nil {
		return nil, err
	}

	return &StreamID{Timestamp: timestamp, Sequence: sequence}, nil
}

func (s *Streams) Add(key string, id string, data [][]byte) (string, error) {
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

	entry := NewStreamEntry(streamID.Timestamp, streamID.Sequence, data)
	s.streams[key].entries = append(s.streams[key].entries, entry)
	return fmt.Sprintf("%d-%d", streamID.Timestamp, streamID.Sequence), nil
}

func generateStreamId(idStr string, lastSequence int, isEnd bool) *StreamID {
	tokens := strings.Split(idStr, "-")
	timestamp, _ := strconv.Atoi(tokens[0])
	sequence := lastSequence
	if len(tokens) > 1 {
		sequence, _ = strconv.Atoi(tokens[1])
	}
	return &StreamID{Timestamp: timestamp, Sequence: sequence}
}

func (s *Streams) List(key string, startStr string, endStr string) []*StreamEntry {
	stream, ok := s.streams[key]
	if !ok {
		return make([]*StreamEntry, 0)
	}
	if len(stream.entries) == 0 {
		return make([]*StreamEntry, 0)
	}
	entry := stream.lastEntry()
	if entry == nil {
		return make([]*StreamEntry, 0)
	}
	start := generateStreamId(startStr, entry.Sequence, false)
	end := generateStreamId(endStr, entry.Sequence, true)
	var entries []*StreamEntry
	for _, entry := range stream.entries {
		streamID := NewStreamID(entry.Timestamp, entry.Sequence)
		if streamID.isInRange(start, end) {
			entries = append(entries, entry)
		}
	}
	return entries
}

func (s *Streams) Contains(key string) bool {
	_, ok := s.streams[key]
	return ok
}
