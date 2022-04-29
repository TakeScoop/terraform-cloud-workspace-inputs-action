package action

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCaseNewDefaults struct {
	message  string
	input    Config
	expected Config
	err      error
}

func TestNewDefaults(t *testing.T) {
	testCases := []testCaseNewDefaults{
		{
			message: "empty environments",
			input:   NewDefaults([]string{}, "empty"),
			expected: Config{
				Environments:          []string{},
				EnvironmentsTags:      map[string][]string{},
				EnvironmentsVariables: map[string][]Variable{},
				Tags:                  []string{"service:empty"},
				Name:                  "empty",
			},
		},
		{
			message: "environments",
			input:   NewDefaults([]string{"staging", "production"}, "name"),
			expected: Config{
				Environments: []string{"staging", "production"},
				EnvironmentsTags: map[string][]string{
					"staging":    {"environment:staging"},
					"production": {"environment:production"},
				},
				EnvironmentsVariables: map[string][]Variable{
					"staging":    {{Key: "environment", Value: "staging", Category: "terraform"}},
					"production": {{Key: "environment", Value: "production", Category: "terraform"}},
				},
				Tags: []string{"service:name"},
				Name: "name",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.message, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.expected, tc.input)
		})
	}
}
