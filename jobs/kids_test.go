package jobs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessKidAgeValue(t *testing.T) {
	tValues := []struct{
		input int
		expected int
	}{
		{0, 0}, {1, 0},
		{2, 2}, {3, 2}, {4, 2}, {5, 2}, {6, 2},
		{7, 7}, {8, 7},
		{9, 9}, {10, 9}, {11, 9}, {12, 9},
		{13, 13}, {14, 13}, {42, 13},
	}

	for _, tValue := range tValues {
		result := processKidAgeValue(tValue.input)
		assert.Equal(t, tValue.expected, result)
	}
}

