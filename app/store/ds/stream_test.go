package ds

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterInANewStream(t *testing.T) {
	streams := NewStreams()
	_, err := streams.Register("Key1", "1526919030473-0")

	assert.NoError(t, err)
}

func TestRegisterInAExistingStream(t *testing.T) {
	streams := NewStreams()
	_, _ = streams.Register("Key1", "1526919030473-0")
	_, err := streams.Register("Key1", "1526919030473-1")

	assert.NoError(t, err)
}

func TestRegisterWithSameSequenceAsHeadOfStream(t *testing.T) {
	streams := NewStreams()
	_, _ = streams.Register("Key1", "1526919030473-0")
	_, err := streams.Register("Key1", "1526919030473-0")

	assert.Error(t, err)
}

func TestRegisterWithHigherTimestampThanHeadOfStream(t *testing.T) {
	streams := NewStreams()
	_, _ = streams.Register("Key1", "1526919030473-10")
	_, err := streams.Register("Key1", "1526919030500-9")

	assert.NoError(t, err)
}

func TestRegisterWithLowerTimestampThanHeadOfStream(t *testing.T) {
	streams := NewStreams()
	_, _ = streams.Register("Key1", "1526919030473-10")
	_, err := streams.Register("Key1", "1526919030400-9")

	assert.Error(t, err)
}

func TestRegisterWithZeroZeroKey(t *testing.T) {
	streams := NewStreams()
	_, err := streams.Register("Key1", "0-0")

	assert.Error(t, err)
}

func TestIncludes(t *testing.T) {
	streams := NewStreams()
	contains := streams.Contains("Key1")
	assert.False(t, contains)

	_, _ = streams.Register("Key2", "1526919030400-10")
	contains = streams.Contains("Key1")
	assert.False(t, contains)

	_, _ = streams.Register("Key1", "1526919030400-10")
	contains = streams.Contains("Key1")
	assert.True(t, contains)
}

func TestRegisterWithAutogenerateSequenceInNonExistingStream(t *testing.T) {
	streams := NewStreams()
	id, err := streams.Register("Key1", "1526919030473-*")

	assert.NoError(t, err)
	assert.Equal(t, "1526919030473-0", id)
}

func TestRegisterWithAutogenerateSequenceInExistingStream(t *testing.T) {
	streams := NewStreams()
	id, err := streams.Register("Key1", "1526919030473-10")

	assert.NoError(t, err)
	assert.Equal(t, "1526919030473-10", id)

	id, err = streams.Register("Key1", "1526919030473-*")

	assert.NoError(t, err)
	assert.Equal(t, "1526919030473-11", id)
}

func TestRegisterWithAutogenerateSequenceInExistingStreamWithMultipleTimestampEnteries(t *testing.T) {
	streams := NewStreams()
	id, err := streams.Register("Key1", "1526919030473-10")
	assert.NoError(t, err)
	assert.Equal(t, "1526919030473-10", id)

	id, err = streams.Register("Key1", "1526919030474-10")
	assert.NoError(t, err)
	assert.Equal(t, "1526919030474-10", id)

	id, err = streams.Register("Key1", "1526919030474-*")
	assert.NoError(t, err)
	assert.Equal(t, "1526919030474-11", id)
}

func TestRegisterWithAutogenerateSequenceWithTimestampZero(t *testing.T) {
	streams := NewStreams()
	id, err := streams.Register("Key1", "0-*")

	assert.NoError(t, err)
	assert.Equal(t, "0-1", id)
}
