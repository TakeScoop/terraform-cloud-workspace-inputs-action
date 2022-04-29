package action

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testCaseOutputsFromInputs struct {
	message  string
	input    Config
	expected testSetOutputsExpected
}

type testSetOutputsExpected struct {
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

func TestOutputsFromInputs(t *testing.T) {
	testCases := []testCaseOutputsFromInputs{
		{
			message: "empty inputs",
			input:   NewConfig("empty"),
			expected: testSetOutputsExpected{
				outputs: map[string]string{
					"workspaces":          `[]`,
					"workspace_variables": `{}`,
					"workspace_tags":      `{}`,
					"tags":                `[]`,
					"name":                "empty",
				},
				masked: []string{},
			},
		},
		{
			message: "default inputs",
			input: Config{
				Environments: []string{"staging", "production"},
				EnvironmentsTags: map[string][]string{
					"staging":    {"environment:staging"},
					"production": {"environment:production"},
				},
				EnvironmentsVariables: map[string][]Variable{
					"staging":    {{Key: "environment", Value: "staging", Category: "terraform"}},
					"production": {{Key: "environment", Value: "production", Category: "terraform"}},
				},
				Tags: []string{"service:workspace"},
				Name: "workspace",
			},
			expected: testSetOutputsExpected{
				outputs: map[string]string{
					"workspaces": `[
						"staging",
						"production"
					]`,
					"workspace_variables": `{
						"staging": [{
							"key": "environment",
							"value": "staging",
							"category": "terraform"
						}],
						"production": [{
							"key": "environment",
							"value": "production",
							"category": "terraform"
						}]
					}`,
					"workspace_tags": `{
						"staging": ["environment:staging"],
						"production": ["environment:production"]
					}`,
					"name": "workspace",
					"tags": `["service:workspace"]`,
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

			require.NoError(t, tc.input.SetOutputs(&out))

			assert.ElementsMatch(t, tc.expected.masked, out.masked)

			for k, v := range tc.expected.outputs {
				if k == "name" {
					assert.Equal(t, v, out.outputs["name"])
				} else {
					assert.JSONEq(t, v, out.outputs[k])
				}
			}
		})
	}
}
