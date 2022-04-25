package action

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testCaseRun struct {
	message  string
	input    Inputs
	expected testCaseRunExpected
	actual   testOutputter
}

type testCaseRunExpected struct {
	outputs map[string]string
	masked  []string
}

type testOutputter struct {
	outputs map[string]string
	masked  []string
}

func (o *testOutputter) SetOutput(k string, v string) {
	o.outputs[k] = v
}

func (o *testOutputter) AddMask(p string) {
	o.masked = append(o.masked, p)
}

func newTestOutputter() testOutputter {
	return testOutputter{
		outputs: map[string]string{},
		masked:  []string{},
	}
}

func TestRun(t *testing.T) {
	testCases := []testCaseRun{
		{
			message: "empty string inputs",
			input: Inputs{
				Environments: "",
				Tags:         "",
				Variables:    "",
			},
			expected: testCaseRunExpected{
				outputs: map[string]string{
					"workspaces":          `[]`,
					"workspace_tags":      `{}`,
					"workspace_variables": `{}`,
				},
				masked: []string{},
			},
		},
		{
			message: "empty yaml inputs",
			input: Inputs{
				Environments: "",
				Tags:         "",
				Variables:    "",
			},
			expected: testCaseRunExpected{
				outputs: map[string]string{
					"workspaces":          `[]`,
					"workspace_tags":      `{}`,
					"workspace_variables": `{}`,
				},
				masked: []string{},
			},
		},
		{
			message: "environments",
			input: Inputs{
				Environments: `---
- staging
- production`,
				Variables: `---
staging:
- key: secret
  value: masked				
  category: terraform
  sensitive: true`,
			},
			expected: testCaseRunExpected{
				outputs: map[string]string{
					"workspaces": `["staging", "production"]`,
					"workspace_tags": `{
						"staging": ["environment:staging"],
						"production": ["environment:production"]
					}`,
					"workspace_variables": `{
						"staging": [
							{"key": "secret", "value": "masked", "category": "terraform", "sensitive": true},
							{"key": "environment", "value": "staging", "category": "terraform"}
						],
						"production": [
							{"key": "environment", "value": "production", "category": "terraform"}
						]
					}`,
				},
				masked: []string{"masked"},
			},
		},
		{
			message: "override variable input",
			input: Inputs{
				Environments: `---
- staging
- production`,
				Variables: `---
staging:
- key: environment
  value: bar
  category: terraform`,
			},
			expected: testCaseRunExpected{
				outputs: map[string]string{
					"workspaces": `["staging", "production"]`,
					"workspace_tags": `{
						"staging": ["environment:staging"],
						"production": ["environment:production"]
					}`,
					"workspace_variables": `{
						"staging": [
							{"key": "environment", "value": "bar", "category": "terraform"}
						],
						"production": [
							{"key": "environment", "value": "production", "category": "terraform"}
						]
					}`,
				},
				masked: []string{},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.message, func(t *testing.T) {
			t.Parallel()
			out := newTestOutputter()

			require.NoError(t, Run(tc.input, &out))

			for k, o := range tc.expected.outputs {
				assert.JSONEq(t, o, out.outputs[k], fmt.Sprintf("JSON output value at key %s does not match expected value", k))
			}

			assert.Equal(t, tc.expected.masked, out.masked, "Actual masked values do not match expected masked values")
		})
	}
}
