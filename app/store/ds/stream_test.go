package ds

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterInANewStream(t *testing.T) {
	streams := NewStreams()
	err := streams.Register("Key1", 1526919030473, 0)

	assert.NoError(t, err)
}

func TestRegisterInAExistingStream(t *testing.T) {
	streams := NewStreams()
	_ = streams.Register("Key1", 1526919030473, 0)
	err := streams.Register("Key1", 1526919030473, 1)

	assert.NoError(t, err)
}

func TestRegisterWithSameSequenceAsHeadOfStream(t *testing.T) {
	streams := NewStreams()
	_ = streams.Register("Key1", 1526919030473, 0)
	err := streams.Register("Key1", 1526919030473, 0)

	assert.Error(t, err)
}

func TestRegisterWithHigherTimestampThanHeadOfStream(t *testing.T) {
	streams := NewStreams()
	_ = streams.Register("Key1", 1526919030473, 10)
	err := streams.Register("Key1", 1526919030500, 9)

	assert.NoError(t, err)
}

func TestRegisterWithLowerTimestampThanHeadOfStream(t *testing.T) {
	streams := NewStreams()
	_ = streams.Register("Key1", 1526919030473, 10)
	err := streams.Register("Key1", 1526919030400, 9)

	assert.Error(t, err)
}

func TestRegisterWithZeroZeroKey(t *testing.T) {
	streams := NewStreams()
	err := streams.Register("Key1", 0, 0)

	assert.Error(t, err)
}

func TestIncludes(t *testing.T) {
	streams := NewStreams()
	contains := streams.Contains("Key1")
	assert.False(t, contains)

	_ = streams.Register("Key2", 1526919030400, 10)
	contains = streams.Contains("Key1")
	assert.False(t, contains)

	_ = streams.Register("Key1", 1526919030400, 10)
	contains = streams.Contains("Key1")
	assert.True(t, contains)

}
