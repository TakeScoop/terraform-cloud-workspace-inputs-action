package action

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCaseExtendConfig struct {
	message  string
	input    [2]Config
	expected Config
}

func TestExtendConfig(t *testing.T) {
	testCases := []testCaseExtendConfig{
		{
			message: "no tags",
			input: [2]Config{
				{Environments: Environments{"staging"}},
				{},
			},
			expected: Config{
				Environments:          Environments{"staging"},
				EnvironmentsTags:      map[string][]string{"staging": {}},
				EnvironmentsVariables: map[string][]Variable{"staging": {}},
			},
		},
		{
			message: "dedupe workspace tags",
			input: [2]Config{
				{
					Environments:     Environments{"staging"},
					EnvironmentsTags: map[string][]string{"staging": {"environment:staging"}},
				},
				{EnvironmentsTags: map[string][]string{"staging": {"environment:staging"}}},
			},
			expected: Config{
				Environments:          Environments{"staging"},
				EnvironmentsTags:      map[string][]string{"staging": {"environment:staging"}},
				EnvironmentsVariables: map[string][]Variable{"staging": {}},
			},
		},
		{
			message: "dedupe variables",
			input: [2]Config{
				{
					Environments: Environments{"staging"},
					EnvironmentsVariables: map[string][]Variable{
						"staging": {{Key: "environment", Value: "staging", Category: "terraform"}},
					},
				},
				{
					EnvironmentsVariables: map[string][]Variable{
						"staging": {{Key: "environment", Value: "staging", Category: "terraform"}},
					},
				},
			},
			expected: Config{
				Environments:     Environments{"staging"},
				EnvironmentsTags: map[string][]string{"staging": {}},
				EnvironmentsVariables: map[string][]Variable{
					"staging": {{Key: "environment", Value: "staging", Category: "terraform"}},
				},
			},
		},
		{
			message: "extend default",
			input: [2]Config{
				{
					Environments: Environments{"staging", "production"},
					EnvironmentsVariables: map[string][]Variable{
						"staging": {
							{Key: "environment", Value: "staging", Category: "terraform"},
						},
						"production": {
							{Key: "environment", Value: "production", Category: "terraform"},
						},
					},
					EnvironmentsTags: map[string][]string{
						"staging":    {"environment:staging"},
						"production": {"environment:production"},
					},
				},
				{
					Environments: Environments{"staging", "production"},
					EnvironmentsVariables: map[string][]Variable{
						"staging": {
							{Key: "environment", Value: "staging", Category: "terraform"},
							{Key: "foo", Value: "bar", Category: "env"},
						},
						"production": {
							{Key: "baz", Value: "woz", Category: "terraform"},
						},
					},
					EnvironmentsTags: map[string][]string{
						"staging":    {"foo:bar"},
						"production": {"environment:production", "baz:woz"},
					},
				},
			},
			expected: Config{
				Environments: Environments{"staging", "production"},
				EnvironmentsVariables: map[string][]Variable{
					"staging": {
						{Key: "environment", Value: "staging", Category: "terraform"},
						{Key: "foo", Value: "bar", Category: "env"},
					},
					"production": {
						{Key: "environment", Value: "production", Category: "terraform"},
						{Key: "baz", Value: "woz", Category: "terraform"},
					},
				},
				EnvironmentsTags: map[string][]string{
					"staging":    {"environment:staging", "foo:bar"},
					"production": {"environment:production", "baz:woz"},
				},
			},
		},
		{
			message: "use name from first config",
			input: [2]Config{
				{Name: "a"},
				{Name: "b"},
			},
			expected: Config{Name: "a"},
		},
		{
			message: "add tags to empty config",
			input: [2]Config{
				{},
				{Tags: []string{"foo:bar"}},
			},
			expected: Config{
				Tags: []string{"foo:bar"},
			},
		},
		{
			message: "dedupe tags from same config",
			input: [2]Config{
				{},
				{Tags: []string{"foo:bar", "foo:bar"}},
			},
			expected: Config{
				Tags: []string{"foo:bar"},
			},
		},
		{
			message: "dedupe tags from different configs",
			input: [2]Config{
				{Tags: []string{"foo:bar"}},
				{Tags: []string{"foo:bar"}},
			},
			expected: Config{
				Tags: []string{"foo:bar"},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.message, func(t *testing.T) {
			t.Parallel()

			actual := ExtendConfig(tc.input[0], tc.input[1])

			assert.Equal(t, tc.expected.Name, actual.Name)

			assert.ElementsMatch(t, tc.expected.Environments, actual.Environments)
			assert.ElementsMatch(t, tc.expected.Tags, actual.Tags)

			for _, e := range actual.Environments {
				assert.ElementsMatch(t, tc.expected.EnvironmentsTags[e], actual.EnvironmentsTags[e])
				assert.ElementsMatch(t, tc.expected.EnvironmentsVariables[e], actual.EnvironmentsVariables[e])
			}
		})
	}
}
