package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseConfig(t *testing.T) {
	config, err := Parse("testdata/valid-config.yml")
	require.NoError(t, err)

	assert.Len(t, config.Databases, 2)
	assert.Equal(t, "dba", config.Databases[0].Name)
	assert.Equal(t, "postgres://localhost:5432/dba?sslmode=disable", config.Databases[0].URL)
}

func TestParseBadConfigs(t *testing.T) {
	tests := []struct {
		ConfigFile    string
		ExpectedError string
	}{
		{
			ConfigFile:    "testdata/unknown-config.yml",
			ExpectedError: "failed to read config file 'testdata/unknown-config.yml': open testdata/unknown-config.yml: no such file or directory",
		},
		{
			ConfigFile:    "testdata/invalid-duplicated-config.yml",
			ExpectedError: "failed to validate configuration. A database named 'dba' has already been declared",
		},
		{
			ConfigFile:    "testdata/invalid-empty-name-config.yml",
			ExpectedError: "failed to validate configuration. Database name cannot be empty",
		},
		{
			ConfigFile:    "testdata/invalid-empty-url-config.yml",
			ExpectedError: "failed to validate configuration. URL for database 'dba' cannot be empty",
		},
		{
			ConfigFile:    "testdata/invalid-yaml-config.yml",
			ExpectedError: "failed to unmarshall config file 'testdata/invalid-yaml-config.yml': yaml: unmarshal errors:\n  line 1: cannot unmarshal !!seq into config.Config",
		},
	}

	for _, test := range tests {
		_, err := Parse(test.ConfigFile)
		if err == nil {
			t.Errorf("In case %v:\nExpected: %v\nGot: nil", test.ConfigFile, test.ExpectedError)
			continue
		}
		if err.Error() != test.ExpectedError {
			t.Errorf("In case %v:\nExpected: %v\nGot: %v", test.ConfigFile, test.ExpectedError, err.Error())
		}
	}
}
