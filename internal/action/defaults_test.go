package action

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCaseNewDefaults struct {
	message  string
	input    Config
	expected Config
}

func TestNewDefaults(t *testing.T) {
	testCases := []testCaseNewDefaults{
		{
			message: "empty environments",
			input:   NewDefaults(Environments{}, "empty"),
			expected: Config{
				Environments:          Environments{},
				EnvironmentsTags:      EnvironmentsTags{},
				EnvironmentsVariables: map[string][]Variable{},
				Tags:                  Tags{"service:empty"},
				Name:                  "empty",
			},
		},
		{
			message: "environments",
			input:   NewDefaults(Environments{"staging", "production"}, "name"),
			expected: Config{
				Environments: Environments{"staging", "production"},
				EnvironmentsTags: EnvironmentsTags{
					"staging":    {"environment:staging"},
					"production": {"environment:production"},
				},
				EnvironmentsVariables: map[string][]Variable{
					"staging":    {{Key: "environment", Value: "staging", Category: "terraform"}},
					"production": {{Key: "environment", Value: "production", Category: "terraform"}},
				},
				Tags: Tags{"service:name"},
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
