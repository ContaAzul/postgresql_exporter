package postgres

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEqualOrGreaterThan96(t *testing.T) {
	tt := []struct {
		version  int
		expected bool
	}{
		{90407, false},
		{90600, true},
		{90606, true},
		{100000, true},
		{100004, true},
		{110000, true},
		{110004, true},
	}

	for _, tc := range tt {
		testName := fmt.Sprintf("expecting IsEqualOrGreaterThan96(\"%v\") to be %v", tc.version, tc.expected)
		t.Run(testName, func(t *testing.T) {
			assert.Equal(t, Version(tc.version).IsEqualOrGreaterThan96(), tc.expected)
		})
	}
}

func TestIsEqualOrGreaterThan10(t *testing.T) {
	tt := []struct {
		version  int
		expected bool
	}{
		{90407, false},
		{90600, false},
		{90606, false},
		{100000, true},
		{100004, true},
		{110000, true},
		{110004, true},
	}

	for _, tc := range tt {
		testName := fmt.Sprintf("expecting IsEqualOrGreaterThan10(\"%v\") to be %v", tc.version, tc.expected)
		t.Run(testName, func(t *testing.T) {
			assert.Equal(t, Version(tc.version).IsEqualOrGreaterThan10(), tc.expected)
		})
	}
}
