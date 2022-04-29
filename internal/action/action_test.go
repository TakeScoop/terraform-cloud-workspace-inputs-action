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

func TestRun(t *testing.T) {
	testCases := []testCaseRun{
		{
			message: "empty string inputs",
			input: Inputs{
				Environments:          "",
				EnvironmentsTags:      "",
				EnvironmentsVariables: "",
				Name:                  "empty",
				Tags:                  "",
			},
			expected: testCaseRunExpected{
				outputs: map[string]string{
					"workspaces":          `[]`,
					"workspace_tags":      `{}`,
					"workspace_variables": `{}`,
					"tags":                `["service:empty"]`,
					"name":                "empty",
				},
				masked: []string{},
			},
		},
		{
			message: "empty yaml inputs",
			input: Inputs{
				Environments:          `[]`,
				EnvironmentsTags:      `{}`,
				EnvironmentsVariables: `{}`,
				Name:                  "empty",
				Tags:                  `[]`,
			},
			expected: testCaseRunExpected{
				outputs: map[string]string{
					"workspaces":          `[]`,
					"workspace_tags":      `{}`,
					"workspace_variables": `{}`,
					"tags":                `["service:empty"]`,
					"name":                "empty",
				},
				masked: []string{},
			},
		},
		{
			message: "single workspace",
			input: Inputs{
				Name:         "workspace",
				Environments: "",
			},
			expected: testCaseRunExpected{
				outputs: map[string]string{
					"workspaces":          `[]`,
					"workspace_tags":      `{}`,
					"workspace_variables": `{}`,
					"tags":                `["service:workspace"]`,
					"name":                "workspace",
				},
				masked: []string{},
			},
		},
		{
			message: "multi-environment workspace",
			input: Inputs{
				Name: "workspace",
				Environments: `---
- staging
- production
`,
			},
			expected: testCaseRunExpected{
				outputs: map[string]string{
					"workspaces": `["staging","production"]`,
					"workspace_tags": `{
						"staging": ["environment:staging"],
						"production": ["environment:production"]
					}`,
					"workspace_variables": `{
						"staging": [
							{"key": "environment", "value": "staging", "category": "terraform"}
						],
						"production": [
							{"key": "environment", "value": "production", "category": "terraform"}
						]
					}`,
					"tags": `["service:workspace"]`,
					"name": "workspace",
				},
				masked: []string{},
			},
		},
		{
			message: "override variable input",
			input: Inputs{
				Name: "workspace",
				Environments: `---
- staging
- production`,
				EnvironmentsVariables: `---
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
					"tags": `["service:workspace"]`,
					"name": "workspace",
				},
				masked: []string{},
			},
		},
		{
			message: "extend tags",
			input: Inputs{
				Name: "workspace",
				Tags: `---
- foo:bar`,
			},
			expected: testCaseRunExpected{
				outputs: map[string]string{
					"workspaces":          `[]`,
					"workspace_tags":      `{}`,
					"workspace_variables": `{}`,
					"tags":                `["foo:bar","service:workspace"]`,
					"name":                "workspace",
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
				if k == "name" {
					assert.Equal(t, o, out.outputs["name"])
				} else {
					assert.JSONEq(t, o, out.outputs[k], fmt.Sprintf("JSON output value at key %s does not match expected value", k))
				}
			}

			assert.Equal(t, tc.expected.masked, out.masked, "Actual masked values do not match expected masked values")
		})
	}
}
