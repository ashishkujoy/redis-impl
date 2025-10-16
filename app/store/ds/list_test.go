package ds

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeStartAndEnd(t *testing.T) {
	tests := []struct {
		name          string
		length        int
		start         int
		end           int
		expectedStart int
		expectedEnd   int
	}{
		{name: "Start as 0", length: 6, start: 0, end: 2, expectedStart: 0, expectedEnd: 2},
		{name: "Start as non zero", length: 6, start: 3, end: 5, expectedStart: 3, expectedEnd: 5},
		{name: "Start greater than length", length: 6, start: 7, end: 5, expectedStart: 7, expectedEnd: 5},
		{name: "Start is negative", length: 6, start: -2, end: 5, expectedStart: 4, expectedEnd: 5},
		{name: "End as 0", length: 6, start: 2, end: 0, expectedStart: 2, expectedEnd: 0},
		{name: "End as last element index", length: 6, start: 3, end: 5, expectedStart: 3, expectedEnd: 5},
		{name: "End greater than length", length: 6, start: 2, end: 15, expectedStart: 2, expectedEnd: 5},
		{name: "End is negative", length: 6, start: 2, end: -2, expectedStart: 2, expectedEnd: 4},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualStart, actualEnd := makeStartAndEnd(test.length, test.start, test.end)
			assert.Equal(t, test.expectedStart, actualStart, test.name)
			assert.Equal(t, test.expectedEnd, actualEnd, test.name)
		})
	}
}
