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
			message: "name required",
			input:   Inputs{},
			err:     ErrNameNotSet,
		},
		{
			message:  "empty inputs",
			input:    Inputs{Name: "empty"},
			expected: Config{Name: "empty"},
		},
		{
			message: "basic workspace tags",
			input: Inputs{
				Name: "workspace",
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
				Name:         "workspace",
				Environments: Environments{"staging", "production"},
				EnvironmentsTags: map[string][]string{
					"staging":    {"foo:bar"},
					"production": {"baz:woz"},
				},
			},
		},
		{
			message: "tags for one environment",
			input: Inputs{
				Name: "workspace",
				Environments: `---
  - staging
  - production`,
				EnvironmentsTags: `---
staging:
  - foo:bar`,
			},
			expected: Config{
				Name:         "workspace",
				Environments: Environments{"staging", "production"},
				EnvironmentsTags: map[string][]string{
					"staging": {"foo:bar"},
				},
			},
		},
		{
			message: "tags non existent environment",
			input: Inputs{
				Name: "workspace",
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
				Name: "workspace",
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
				Name:         "workspace",
				Environments: Environments{"staging", "production"},
				EnvironmentsVariables: map[string][]Variable{
					"staging":    {{Key: "foo", Value: "bar", Category: "terraform"}},
					"production": {{Key: "baz", Value: "woz", Category: "terraform"}},
				},
			},
		},
		{
			message: "variable for one environment",
			input: Inputs{
				Name: "workspace",
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
				Name:         "workspace",
				Environments: Environments{"staging", "production"},
				EnvironmentsVariables: map[string][]Variable{
					"staging": {{Key: "foo", Value: "bar", Category: "terraform"}},
				},
			},
		},
		{
			message: "variable for missing environment",
			input: Inputs{
				Name: "workspace",
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
		{
			message:  "workspace name",
			input:    Inputs{Name: "foo"},
			expected: Config{Name: "foo"},
		},
		{
			message: "tags",
			input: Inputs{
				Name: "workspace",
				Tags: `---
- department:engineering
- division:platform`,
			},
			expected: Config{
				Name: "workspace",
				Tags: []string{"department:engineering", "division:platform"},
			},
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
