package action

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testCaseMergeEnvironmentsVariables struct {
	message  string
	input    [2]EnvironmentsVariables
	expected EnvironmentsVariables
}

func TestMergeEnvironmentsVariables(t *testing.T) {
	testCases := []testCaseMergeEnvironmentsVariables{
		{
			message:  "empty",
			input:    [2]EnvironmentsVariables{{}, {}},
			expected: EnvironmentsVariables{},
		},
		{
			message: "dedupe",
			input: [2]EnvironmentsVariables{
				{
					"staging": []Variable{
						{Key: "environment", Value: "staging", Category: "terraform"},
					},
					"production": []Variable{
						{Key: "environment", Value: "staging", Category: "terraform"},
					},
				},
				{
					"staging": []Variable{
						{Key: "environment", Value: "staging", Category: "terraform"},
					},
				},
			},
			expected: EnvironmentsVariables{
				"staging": []Variable{
					{Key: "environment", Value: "staging", Category: "terraform"},
				},
				"production": []Variable{
					{Key: "environment", Value: "staging", Category: "terraform"},
				},
			},
		},
		{
			message: "add",
			input: [2]EnvironmentsVariables{
				{
					"staging": []Variable{
						{Key: "foo", Value: "bar", Category: "terraform"},
					},
				},
				{
					"staging": []Variable{
						{Key: "environment", Value: "staging", Category: "terraform"},
					},
					"production": []Variable{
						{Key: "environment", Value: "staging", Category: "terraform"},
					},
				},
			},
			expected: EnvironmentsVariables{
				"staging": []Variable{
					{Key: "foo", Value: "bar", Category: "terraform"},
					{Key: "environment", Value: "staging", Category: "terraform"},
				},
				"production": []Variable{
					{Key: "environment", Value: "staging", Category: "terraform"},
				},
			},
		},
		{
			message: "add to empty",
			input: [2]EnvironmentsVariables{
				{},
				{
					"staging": []Variable{
						{Key: "environment", Value: "staging", Category: "terraform"},
					},
					"production": []Variable{
						{Key: "environment", Value: "staging", Category: "terraform"},
					},
				},
			},
			expected: EnvironmentsVariables{
				"staging": []Variable{
					{Key: "environment", Value: "staging", Category: "terraform"},
				},
				"production": []Variable{
					{Key: "environment", Value: "staging", Category: "terraform"},
				},
			},
		},
		{
			message: "add empty",
			input: [2]EnvironmentsVariables{
				{
					"staging": []Variable{
						{Key: "environment", Value: "staging", Category: "terraform"},
					},
					"production": []Variable{
						{Key: "environment", Value: "staging", Category: "terraform"},
					},
				},
				{},
			},
			expected: EnvironmentsVariables{
				"staging": []Variable{
					{Key: "environment", Value: "staging", Category: "terraform"},
				},
				"production": []Variable{
					{Key: "environment", Value: "staging", Category: "terraform"},
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.message, func(t *testing.T) {
			t.Parallel()
			actual := MergeEnvironmentsVariables(tc.input[0], tc.input[1])
			assert.Equal(t, tc.expected, actual)
		})
	}
}

type testCaseEnvironmentsVariablesSetOutput struct {
	message  string
	input    EnvironmentsVariables
	expected []string
}

func TestEnvironmentsVariablesSetOutput(t *testing.T) {
	testCases := []testCaseEnvironmentsVariablesSetOutput{
		{
			message:  "empty",
			input:    EnvironmentsVariables{},
			expected: []string{},
		},
		{
			message: "no sensitive",
			input: EnvironmentsVariables{
				"staging": []Variable{
					{Category: "env", Key: "foo", Value: "bar"},
				},
			},
			expected: []string{},
		},
		{
			message: "sensitive",
			input: EnvironmentsVariables{
				"staging": []Variable{
					{Category: "env", Key: "foo", Value: "bar", Sensitive: true},
					{Category: "env", Key: "ok", Value: "not sensitive"},
				},
				"production": []Variable{
					{Category: "env", Key: "bar", Value: "baz", Sensitive: true},
				},
			},
			expected: []string{"bar", "baz"},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.message, func(t *testing.T) {
			t.Parallel()

			o := newTestOutputter()

			require.NoError(t, tc.input.SetOutputs(&o))

			assert.ElementsMatch(t, tc.expected, o.masked)
		})
	}
}

type testCaseParseEnvironmentsVariables struct {
	message  string
	input    testCaseParseEnvironmentsVariablesInput
	expected EnvironmentsVariables
	err      error
}

type testCaseParseEnvironmentsVariablesInput struct {
	value string
	envs  Environments
}

func TestParseEnvironmentsVariables(t *testing.T) {
	testCases := []testCaseParseEnvironmentsVariables{
		{
			message: "empty",
			input: testCaseParseEnvironmentsVariablesInput{
				value: "",
				envs:  Environments{},
			},
		},
		{
			message: "with values",
			input: testCaseParseEnvironmentsVariablesInput{
				value: `---
staging:
- key: environment
  value: staging 
  category: terraform
production:
- key: environment
  value: production
  category: terraform`,
				envs: Environments{"staging", "production"},
			},
			expected: EnvironmentsVariables{
				"staging": []Variable{
					{Key: "environment", Value: "staging", Category: "terraform"},
				},
				"production": []Variable{
					{Key: "environment", Value: "production", Category: "terraform"},
				},
			},
		},
		{
			message: "error when environment not passed",
			input: testCaseParseEnvironmentsVariablesInput{
				value: `---
staging:
- key: environment
  value: staging 
  category: terraform
production:
- key: environment
  value: production
  category: terraform`,
				envs: Environments{"staging"},
			},
			err: ErrEnvironmentNotFound,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.message, func(t *testing.T) {
			vars, err := ParseEnvironmentsVariables(tc.input.value, tc.input.envs)

			if tc.err != nil {
				assert.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, vars)
			}
		})
	}
}
