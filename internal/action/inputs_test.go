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
				Names: []string{"staging", "production"},
			},
		},
		{
			message: "basic workspace tags",
			input: Inputs{
				Tags: `---
staging:
  - foo:bar
production:
  - baz:woz`,
				Environments: `---
  - staging
  - production`,
			},
			expected: Config{
				Names: []string{"staging", "production"},
				Tags: map[string][]string{
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
				Tags: `---
staging:
  - foo:bar`,
			},
			expected: Config{
				Names: []string{"staging", "production"},
				Tags: map[string][]string{
					"staging": {"foo:bar"},
				},
			},
		},
		{
			message: "tags non existent environment",
			input: Inputs{
				Environments: `---
- production`,
				Tags: `---
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
				Variables: `---
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
				Names: []string{"staging", "production"},
				Variables: map[string][]Variable{
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
				Variables: `---
staging:
- key: foo
  value: bar
  category: terraform`,
			},
			expected: Config{
				Names: []string{"staging", "production"},
				Variables: map[string][]Variable{
					"staging": {{Key: "foo", Value: "bar", Category: "terraform"}},
				},
			},
		},
		{
			message: "variable for missing environment",
			input: Inputs{
				Environments: `---
- production`,
				Variables: `---
staging:
- key: foo
  value: bar
  category: terraform`,
			},
			err: ErrEnvironmentNotFound,
		},
	}

	for _, tc := range testCases {
		func(tc testCaseInputsParse) {
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
		}(tc)
	}
}
