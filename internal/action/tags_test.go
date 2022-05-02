package action

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

			assert.ElementsMatch(t, tc.expected, MergeTags(tc.input[0], tc.input[1]))
		})
	}
}

type testCaseMergeEnvironmentsTags struct {
	message  string
	expected EnvironmentsTags
	input    [2]EnvironmentsTags
}

func TestMergeEnvironmentsTags(t *testing.T) {
	testCases := []testCaseMergeEnvironmentsTags{
		{
			message:  "empty",
			input:    [2]EnvironmentsTags{},
			expected: EnvironmentsTags{},
		},
		{
			message: "dedupe",
			input: [2]EnvironmentsTags{
				{"staging": Tags{"foo", "bar"}},
				{"staging": Tags{"foo"}},
			},
			expected: EnvironmentsTags{"staging": Tags{"foo", "bar"}},
		},
		{
			message: "add",
			input: [2]EnvironmentsTags{
				{
					"staging":    Tags{"foo"},
					"production": Tags{"bar"},
				},
				{
					"staging":    Tags{"environment:staging"},
					"production": Tags{"environment:production"},
				},
			},
			expected: EnvironmentsTags{
				"staging":    Tags{"foo", "environment:staging"},
				"production": Tags{"bar", "environment:production"},
			},
		},
		{
			message: "add empty",
			input: [2]EnvironmentsTags{
				{"staging": Tags{"foo"}},
				{"staging": Tags{}},
			},
			expected: EnvironmentsTags{"staging": Tags{"foo"}},
		},
		{
			message: "add from empty",
			input: [2]EnvironmentsTags{
				{"staging": Tags{}},
				{"staging": Tags{"foo"}},
			},
			expected: EnvironmentsTags{"staging": Tags{"foo"}},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.message, func(t *testing.T) {
			t.Parallel()

			actual := MergeEnvironmentsTags(tc.input[0], tc.input[1])

			for e, tags := range tc.expected {
				assert.ElementsMatch(t, tags, actual[e])
			}
		})
	}
}

type testCaseParseTags struct {
	message  string
	input    string
	expected Tags
}

func TestParseTags(t *testing.T) {
	testCases := []testCaseParseTags{
		{
			message:  "empty",
			input:    "",
			expected: Tags{},
		},
		{
			message: "list of tags",
			input: `---
- foo:bar
- baz:woz
`,
			expected: Tags{"foo:bar", "baz:woz"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.message, func(t *testing.T) {
			t.Parallel()

			tags, err := ParseTags(tc.input)
			require.NoError(t, err)

			assert.ElementsMatch(t, tc.expected, tags)
		})
	}
}

type testCaseParseEnvironmentsTags struct {
	message  string
	input    testCaseParseEnvironmentsTagsInput
	expected EnvironmentsTags
	err      error
}

type testCaseParseEnvironmentsTagsInput struct {
	value string
	envs  Environments
}

func TestParseEnvironmentsTags(t *testing.T) {
	testCases := []testCaseParseEnvironmentsTags{
		{
			message: "empty",
			input: testCaseParseEnvironmentsTagsInput{
				value: "",
				envs:  Environments{},
			},
			expected: EnvironmentsTags{},
		},
		{
			message: "map of tags",
			input: testCaseParseEnvironmentsTagsInput{
				value: `---
staging:
- foo:bar
- baz:woz
production:
- environment:production
`,
				envs: Environments{"staging"},
			},
			expected: EnvironmentsTags{
				"staging":    {"foo:bar", "baz:woz"},
				"production": {"environment:production"},
			},
		},
		{
			message: "should error if tag parsed for non existent environment",
			input: testCaseParseEnvironmentsTagsInput{
				value: `---
staging:
- foo:bar
production:
- baz:woz
`,
				envs: Environments{"staging"},
			},
			err: ErrEnvironmentNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.message, func(t *testing.T) {
			t.Parallel()

			tags, err := ParseEnvironmentsTags(tc.input.value, tc.input.envs)
			if tc.err != nil {
				assert.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)

				for env, expected := range tc.expected {
					assert.ElementsMatch(t, expected, tags[env])
				}
			}
		})
	}
}

type testCaseTagsSetOutput struct {
	message  string
	input    Tags
	expected string
}

func TestTagsSetOutput(t *testing.T) {
	testCases := []testCaseTagsSetOutput{
		{
			message:  "empty",
			input:    Tags{},
			expected: `[]`,
		},
		{
			message:  "with tags",
			input:    Tags{"service:name"},
			expected: `["service:name"]`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.message, func(t *testing.T) {
			t.Parallel()
			out := newTestOutputter()
			require.NoError(t, tc.input.SetOutputs(&out))
			assert.JSONEq(t, tc.expected, out.outputs["tags"])
		})
	}
}

type testCaseEnvironmentsTagsSetOutput struct {
	message  string
	input    EnvironmentsTags
	expected string
}

func TestEnvironmentsTagsSetOutput(t *testing.T) {
	testCases := []testCaseEnvironmentsTagsSetOutput{
		{
			message:  "empty",
			input:    EnvironmentsTags{},
			expected: `{}`,
		},
		{
			message: "with tags",
			input: EnvironmentsTags{
				"staging":    Tags{"environment:staging"},
				"production": Tags{"environment:production"},
			},
			expected: `{
				"staging": ["environment:staging"],
				"production": ["environment:production"]
			}`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.message, func(t *testing.T) {
			t.Parallel()
			out := newTestOutputter()
			require.NoError(t, tc.input.SetOutputs(&out))
			assert.JSONEq(t, tc.expected, out.outputs["workspace_tags"])
		})
	}
}
