package action

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testCaseMergeEnvironmnets struct {
	message  string
	input    [2]Environments
	expected Environments
}

func TestMergeEnvironments(t *testing.T) {
	testCases := []testCaseMergeEnvironmnets{
		{
			message:  "empty environments",
			input:    [2]Environments{{}, {}},
			expected: Environments{},
		},
		{
			message:  "dedupe",
			input:    [2]Environments{{"staging", "production"}, {"staging"}},
			expected: Environments{"staging", "production"},
		},
		{
			message:  "add",
			input:    [2]Environments{{"staging", "production"}, {"staging", "playground"}},
			expected: Environments{"staging", "production", "playground"},
		},
		{
			message:  "add to empty",
			input:    [2]Environments{{}, {"staging", "production"}},
			expected: Environments{"staging", "production"},
		},
		{
			message:  "add empty",
			input:    [2]Environments{{"staging", "production"}, {}},
			expected: Environments{"staging", "production"},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.message, func(t *testing.T) {
			t.Parallel()
			actual := MergeEnvironments(tc.input[0], tc.input[1])
			assert.ElementsMatch(t, tc.expected, actual)
		})
	}
}

type testCaseParseEnvironments struct {
	message  string
	input    string
	expected Environments
}

func TestParseEnvironments(t *testing.T) {
	testCases := []testCaseParseEnvironments{
		{
			message:  "empty environments",
			input:    "",
			expected: Environments{},
		},
		{
			message: "list of environments",
			input: `---
- staging			
- production
`,
			expected: Environments{"staging", "production"},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.message, func(t *testing.T) {
			envs, err := ParseEnvironments(tc.input)
			require.NoError(t, err)

			assert.ElementsMatch(t, tc.expected, envs)
		})
	}
}
