package action

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testCaseInputsParse struct {
	message  string
	input    Inputs
	expected Config
	err      error
}

func TestInputsParse(t *testing.T) {
	testCases := []testCaseInputsParse{
		{
			message:  "empty inputs",
			input:    Inputs{},
			expected: Config{},
		},
		{
			message: "basic environments",
			input: Inputs{
				Environments: `---
- staging
- production`,
			},
			expected: Config{
				Environments: []string{"staging", "production"},
			},
		},
		{
			message: "basic workspace tags",
			input: Inputs{
				EnvironmentsTags: `---
staging:
  - foo:bar
production:
  - baz:woz`,
				Environments: `---
  - staging
  - production`,
			},
			expected: Config{
				Environments: []string{"staging", "production"},
				EnvironmentsTags: map[string][]string{
					"staging":    {"foo:bar"},
					"production": {"baz:woz"},
				},
			},
		},
		{
			message: "tags for one environment",
			input: Inputs{
				Environments: `---
  - staging
  - production`,
				EnvironmentsTags: `---
staging:
  - foo:bar`,
			},
			expected: Config{
				Environments: []string{"staging", "production"},
				EnvironmentsTags: map[string][]string{
					"staging": {"foo:bar"},
				},
			},
		},
		{
			message: "tags non existent environment",
			input: Inputs{
				Environments: `---
- production`,
				EnvironmentsTags: `---
staging:
  - foo:bar`,
			},
			err: ErrEnvironmentNotFound,
		},
		{
			message: "basic variables",
			input: Inputs{
				Environments: `---
- staging
- production`,
				EnvironmentsVariables: `---
staging:
- key: foo
  value: bar
  category: terraform
production:
- key: baz
  value: woz
  category: terraform`,
			},
			expected: Config{
				Environments: []string{"staging", "production"},
				EnvironmentsVariables: map[string][]Variable{
					"staging":    {{Key: "foo", Value: "bar", Category: "terraform"}},
					"production": {{Key: "baz", Value: "woz", Category: "terraform"}},
				},
			},
		},
		{
			message: "variable for one environment",
			input: Inputs{
				Environments: `---
- staging
- production`,
				EnvironmentsVariables: `---
staging:
- key: foo
  value: bar
  category: terraform`,
			},
			expected: Config{
				Environments: []string{"staging", "production"},
				EnvironmentsVariables: map[string][]Variable{
					"staging": {{Key: "foo", Value: "bar", Category: "terraform"}},
				},
			},
		},
		{
			message: "variable for missing environment",
			input: Inputs{
				Environments: `---
- production`,
				EnvironmentsVariables: `---
staging:
- key: foo
  value: bar
  category: terraform`,
			},
			err: ErrEnvironmentNotFound,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.message, func(t *testing.T) {
			t.Parallel()

			i, err := tc.input.Parse()
			if tc.err != nil {
				assert.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, i)
			}
		})
	}
}
