package action

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Variable struct {
	Key         string `yaml:"key" json:"key"`
	Value       string `yaml:"value" json:"value"`
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
	Category    string `yaml:"category" json:"category"`
	Sensitive   bool   `yaml:"sensitive,omitempty" json:"sensitive,omitempty"`
}

// EnvironmentsVariables maps environment names to list of Terraform workspace variables
type EnvironmentsVariables map[string][]Variable

// ParseEnvironmentsVariables returns an EnvironmentsVariables struct from a string
func ParseEnvironmentsVariables(s string, envs Environments) (EnvironmentsVariables, error) {
	wsVars := EnvironmentsVariables{}

	if err := yaml.Unmarshal([]byte(s), &wsVars); err != nil {
		return EnvironmentsVariables{}, fmt.Errorf("failed to parse workspace variables: %w", err)
	}

	eMap := map[string]bool{}
	for _, e := range envs {
		eMap[e] = true
	}

	for env := range wsVars {
		if _, ok := eMap[env]; !ok {
			return nil, fmt.Errorf("environment %s in passed variables not found in environments %v: %w", env, envs, ErrEnvironmentNotFound)
		}
	}

	return wsVars, nil
}

func (ev EnvironmentsVariables) environments() Environments {
	envs := Environments{}

	for e := range ev {
		envs = append(envs, e)
	}

	return envs
}

// MergeEnvironmentsVariables merges two EnvironmentsVariables structs, removing duplicates
func MergeEnvironmentsVariables(a EnvironmentsVariables, b EnvironmentsVariables) EnvironmentsVariables {
	out := EnvironmentsVariables{}

	envs := MergeEnvironments(a.environments(), b.environments())

	for _, e := range envs {
		out[e] = []Variable{}

		varMap := map[string]Variable{}

		if aVars, ok := a[e]; ok {
			for _, v := range aVars {
				if _, ok := varMap[fmt.Sprintf("%s-%s", v.Key, v.Category)]; !ok {
					out[e] = append(out[e], v)
					varMap[fmt.Sprintf("%s-%s", v.Key, v.Category)] = v
				}
			}
		}

		if bVars, ok := b[e]; ok {
			for _, v := range bVars {
				if _, ok := varMap[fmt.Sprintf("%s-%s", v.Key, v.Category)]; !ok {
					out[e] = append(out[e], v)
					varMap[fmt.Sprintf("%s-%s", v.Key, v.Category)] = v
				}
			}
		}
	}

	return out
}

// SetOutputs uses the passed outputter to output environment variable information
func (ev EnvironmentsVariables) SetOutputs(o Outputter) error {
	for _, vars := range ev {
		for _, v := range vars {
			if v.Sensitive {
				o.AddMask(v.Value)
			}
		}
	}

	return setJSONOutput(o, "workspace_variables", ev)
}
