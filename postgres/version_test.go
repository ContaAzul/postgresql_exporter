package postgres

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFunctionVersions(t *testing.T) {
	is96 := Version.is96
	is10 := Version.is10
	is96Or10 := Version.Is96Or10

	tt := []struct {
		str      string
		fn       func(Version) bool
		expected bool
	}{
		{"9.6.6", is96, true},
		{"9.5.4", is96, false},
		{"10.3", is10, true},
		{"9.6.6", is10, false},
		{"9.6.6", is96Or10, true},
		{"9.5.4", is96Or10, false},
		{"10.4", is96Or10, true},
		{"11.0", is96Or10, false},
	}

	for _, tc := range tt {
		testName := fmt.Sprintf("expecting %v(\"%v\") to be %v", tc.str, tc.fn, tc.expected)
		t.Run(testName, func(t *testing.T) {
			assert.Equal(t, tc.fn(Version(tc.str)), tc.expected)
		})
	}
}
