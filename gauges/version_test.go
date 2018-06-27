package gauges

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	var assert = assert.New(t)
	_, gauges, close := prepare(t)
	defer close()
	assert.NotEmpty(gauges.version())
	assertNoErrs(t, gauges)
}

func TestVersionIsPG96(t *testing.T) {
	var assert = assert.New(t)
	assert.True(isPG96("9.6.6"))
}

func TestVersionIsNotPG96(t *testing.T) {
	var assert = assert.New(t)
	assert.False(isPG96("9.5.4"))
}

func TestVersionIsPG10(t *testing.T) {
	var assert = assert.New(t)
	assert.True(isPG10("10.3"))
}

func TestVersionIsNotPG10(t *testing.T) {
	var assert = assert.New(t)
	assert.False(isPG10("9.6.6"))
}
