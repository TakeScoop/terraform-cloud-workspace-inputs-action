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

type EnvironmentsVariables map[string][]Variable

func ParseEnvironmentsVariables(s string, envs Environments) (EnvironmentsVariables, error) {
	var wsVars EnvironmentsVariables

	if err := yaml.Unmarshal([]byte(s), &wsVars); err != nil {
		return EnvironmentsVariables{}, fmt.Errorf("failed to parse workspace variables: %w", err)
	}

	eMap := map[string]bool{}
	for _, e := range envs {
		eMap[e] = true
	}

	for env := range wsVars {
		if _, ok := eMap[env]; !ok {
			return nil, fmt.Errorf("environment %s in passed tags not found in environments %v: %w", env, envs, ErrEnvironmentNotFound)
		}
	}

	return wsVars, nil
}

func mergeEnvironmentsVariables(a EnvironmentsVariables, b EnvironmentsVariables, envs Environments) EnvironmentsVariables {
	out := EnvironmentsVariables{}

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

func (ev EnvironmentsVariables) setOutputs(o Outputter) error {
	for _, vars := range ev {
		for _, v := range vars {
			if v.Sensitive {
				o.AddMask(v.Value)
			}
		}
	}

	return setJSONOutput(o, "workspace_variables", ev)
}
