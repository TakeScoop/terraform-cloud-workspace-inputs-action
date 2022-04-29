package action

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCaseMergeTags struct {
	message  string
	expected Tags
	input    [2]Tags
}

func TestMergeTags(t *testing.T) {
	testCases := []testCaseMergeTags{
		{
			message:  "empty",
			input:    [2]Tags{},
			expected: Tags{},
		},
		{
			message: "dedupe",
			input: [2]Tags{
				{"foo:bar", "environment:staging"},
				{"environment:staging"},
			},
			expected: Tags{"foo:bar", "environment:staging"},
		},
		{
			message: "add",
			input: [2]Tags{
				{"environment:staging"},
				{"environment:staging", "foo:bar"},
			},
			expected: Tags{"foo:bar", "environment:staging"},
		},
		{
			message: "add empty",
			input: [2]Tags{
				{"environment:staging"},
				{},
			},
			expected: Tags{"environment:staging"},
		},
		{
			message: "add to empty",
			input: [2]Tags{
				{},
				{"environment:staging"},
			},
			expected: Tags{"environment:staging"},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.message, func(t *testing.T) {
			t.Parallel()

			assert.ElementsMatch(t, tc.expected, mergeTags(tc.input[0], tc.input[1]))
		})
	}
}

type testCaseMergeEnvironmnetsTags struct {
	message  string
	expected EnvironmentsTags
	input    [2]EnvironmentsTags
}

func TestMergeEnvironmnetsTags(t *testing.T) {
	testCases := []testCaseMergeEnvironmnetsTags{
		{
			message:  "empty",
			input:    [2]EnvironmentsTags{},
			expected: EnvironmentsTags{},
		},
		{
			message: "dedupe",
			input: [2]EnvironmentsTags{
				{"staging": []string{"foo", "bar"}},
				{"staging": []string{"foo"}},
			},
			expected: EnvironmentsTags{"staging": []string{"foo", "bar"}},
		},
		{
			message: "add",
			input: [2]EnvironmentsTags{
				{
					"staging":    []string{"foo"},
					"production": []string{"bar"},
				},
				{
					"staging":    []string{"environment:staging"},
					"production": []string{"environment:production"},
				},
			},
			expected: EnvironmentsTags{
				"staging":    []string{"foo", "environment:staging"},
				"production": []string{"bar", "environment:production"},
			},
		},
		{
			message: "add empty",
			input: [2]EnvironmentsTags{
				{"staging": []string{"foo"}},
				{"staging": []string{}},
			},
			expected: EnvironmentsTags{"staging": []string{"foo"}},
		},
		{
			message: "add from empty",
			input: [2]EnvironmentsTags{
				{"staging": []string{}},
				{"staging": []string{"foo"}},
			},
			expected: EnvironmentsTags{"staging": []string{"foo"}},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.message, func(t *testing.T) {
			t.Parallel()

			actual := mergeEnvironmentsTags(tc.input[0], tc.input[1])

			for e, tags := range tc.expected {
				assert.ElementsMatch(t, tags, actual[e])
			}
		})
	}
}
