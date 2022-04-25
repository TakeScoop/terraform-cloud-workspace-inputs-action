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
				{Names: []string{"staging"}},
				{},
			},
			expected: Config{
				Names:     []string{"staging"},
				Tags:      map[string][]string{"staging": {}},
				Variables: map[string][]Variable{"staging": {}},
			},
		},
		{
			message: "dedupe workspace tags",
			input: [2]Config{
				{
					Names: []string{"staging"},
					Tags:  map[string][]string{"staging": {"environment:staging"}},
				},
				{Tags: map[string][]string{"staging": {"environment:staging"}}},
			},
			expected: Config{
				Names:     []string{"staging"},
				Tags:      map[string][]string{"staging": {"environment:staging"}},
				Variables: map[string][]Variable{"staging": {}},
			},
		},
		{
			message: "dedupe variables",
			input: [2]Config{
				{
					Names: []string{"staging"},
					Variables: map[string][]Variable{
						"staging": {{Key: "environment", Value: "staging", Category: "terraform"}},
					},
				},
				{
					Variables: map[string][]Variable{
						"staging": {{Key: "environment", Value: "staging", Category: "terraform"}},
					},
				},
			},
			expected: Config{
				Names: []string{"staging"},
				Tags:  map[string][]string{"staging": {}},
				Variables: map[string][]Variable{
					"staging": {{Key: "environment", Value: "staging", Category: "terraform"}},
				},
			},
		},
		{
			message: "extend default",
			input: [2]Config{
				{
					Names: []string{"staging", "production"},
					Variables: map[string][]Variable{
						"staging": {
							{Key: "environment", Value: "staging", Category: "terraform"},
						},
						"production": {
							{Key: "environment", Value: "production", Category: "terraform"},
						},
					},
					Tags: map[string][]string{
						"staging":    {"environment:staging"},
						"production": {"environment:production"},
					},
				},
				{
					Names: []string{"staging", "production"},
					Variables: map[string][]Variable{
						"staging": {
							{Key: "environment", Value: "staging", Category: "terraform"},
							{Key: "foo", Value: "bar", Category: "env"},
						},
						"production": {
							{Key: "baz", Value: "woz", Category: "terraform"},
						},
					},
					Tags: map[string][]string{
						"staging":    {"foo:bar"},
						"production": {"environment:production", "baz:woz"},
					},
				},
			},
			expected: Config{
				Names: []string{"staging", "production"},
				Variables: map[string][]Variable{
					"staging": {
						{Key: "environment", Value: "staging", Category: "terraform"},
						{Key: "foo", Value: "bar", Category: "env"},
					},
					"production": {
						{Key: "environment", Value: "production", Category: "terraform"},
						{Key: "baz", Value: "woz", Category: "terraform"},
					},
				},
				Tags: map[string][]string{
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
			assert.ElementsMatch(t, tc.expected.Names, actual.Names)

			for _, e := range actual.Names {
				assert.ElementsMatch(t, tc.expected.Tags[e], actual.Tags[e])
				assert.ElementsMatch(t, tc.expected.Variables[e], actual.Variables[e])
			}
		})
	}
}
