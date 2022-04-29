package action

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
			actual := mergeEnvironmentsVariables(tc.input[0], tc.input[1])
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

			tc.input.setOutputs(&o)

			assert.ElementsMatch(t, tc.expected, o.masked)
		})
	}
}
