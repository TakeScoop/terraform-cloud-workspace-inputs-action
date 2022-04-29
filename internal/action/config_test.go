package action

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCaseMergeConfigs struct {
	message  string
	input    [2]Config
	expected Config
}

func TestMergeConfig(t *testing.T) {
	testCases := []testCaseMergeConfigs{
		{
			message: "no tags",
			input: [2]Config{
				{Environments: []string{"staging"}},
				{},
			},
			expected: Config{
				Environments:          []string{"staging"},
				EnvironmentsTags:      map[string][]string{"staging": {}},
				EnvironmentsVariables: map[string][]Variable{"staging": {}},
			},
		},
		{
			message: "dedupe workspace tags",
			input: [2]Config{
				{
					Environments:     []string{"staging"},
					EnvironmentsTags: map[string][]string{"staging": {"environment:staging"}},
				},
				{EnvironmentsTags: map[string][]string{"staging": {"environment:staging"}}},
			},
			expected: Config{
				Environments:          []string{"staging"},
				EnvironmentsTags:      map[string][]string{"staging": {"environment:staging"}},
				EnvironmentsVariables: map[string][]Variable{"staging": {}},
			},
		},
		{
			message: "dedupe variables",
			input: [2]Config{
				{
					Environments: []string{"staging"},
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
				Environments:     []string{"staging"},
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
					Environments: []string{"staging", "production"},
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
					Environments: []string{"staging", "production"},
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
				Environments: []string{"staging", "production"},
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
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.message, func(t *testing.T) {
			t.Parallel()

			actual := MergeConfigs(tc.input[0], tc.input[1])
			assert.ElementsMatch(t, tc.expected.Environments, actual.Environments)

			for _, e := range actual.Environments {
				assert.ElementsMatch(t, tc.expected.EnvironmentsTags[e], actual.EnvironmentsTags[e])
				assert.ElementsMatch(t, tc.expected.EnvironmentsVariables[e], actual.EnvironmentsVariables[e])
			}
		})
	}
}
