package serializer

import (
	"testing"

	"github.com/codecrafters-io/redis-starter-go/app/store/ds"
	"github.com/stretchr/testify/assert"
)

func TestEncodeStreamList(t *testing.T) {
	streamEntries := []*ds.StreamEntryView{
		ds.NewStreamEntryView("1526985054069-0", []string{"temperature", "36", "humidity", "95"}),
		ds.NewStreamEntryView("1526985054079-0", []string{"temperature", "37", "humidity", "94"}),
	}
	serializer := NewRESPSerializer()

	bytes, err := serializer.EncodeXRange(streamEntries)
	assert.NoError(t, err)
	expectedResp := "*2\r\n" +
		"*2\r\n" +
		"$15\r\n1526985054069-0\r\n" +
		"*4\r\n" +
		"$11\r\ntemperature\r\n" +
		"$2\r\n36\r\n" +
		"$8\r\nhumidity\r\n" +
		"$2\r\n95\r\n" +
		"*2\r\n" +
		"$15\r\n1526985054079-0\r\n" +
		"*4\r\n" +
		"$11\r\ntemperature\r\n" +
		"$2\r\n37\r\n" +
		"$8\r\nhumidity\r\n" +
		"$2\r\n94\r\n"
	assert.Equal(t, expectedResp, string(bytes))
}

func TestEncodeXRead(t *testing.T) {
	streamEntries := []*ds.StreamEntryView{
		ds.NewStreamEntryView("1526985054069-0", []string{"temperature", "36", "humidity", "95"}),
		ds.NewStreamEntryView("1526985054079-0", []string{"temperature", "37", "humidity", "94"}),
	}
	serializer := NewRESPSerializer()

	bytes, err := serializer.EncodeXRead("some_key", streamEntries)
	assert.NoError(t, err)

	expectedResp := "*1\r\n" +
		"*2\r\n" +
		"$8\r\nsome_key\r\n" +
		"*2\r\n" +
		"*2\r\n" +
		"$15\r\n1526985054069-0\r\n" +
		"*4\r\n" +
		"$11\r\ntemperature\r\n" +
		"$2\r\n36\r\n" +
		"$8\r\nhumidity\r\n" +
		"$2\r\n95\r\n" +
		"*2\r\n" +
		"$15\r\n1526985054079-0\r\n" +
		"*4\r\n" +
		"$11\r\ntemperature\r\n" +
		"$2\r\n37\r\n" +
		"$8\r\nhumidity\r\n" +
		"$2\r\n94\r\n"

	assert.Equal(t, expectedResp, string(bytes))
}
