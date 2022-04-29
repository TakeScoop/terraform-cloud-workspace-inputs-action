package action

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Environments []string

func ParseEnvironments(s string) (Environments, error) {
	var environments []string

	if err := yaml.Unmarshal([]byte(s), &environments); err != nil {
		return Environments{}, fmt.Errorf("failed to parse environments: %w", err)
	}

	return environments, nil
}

func mergeEnvironments(a Environments, b Environments) Environments {
	out := Environments{}

	eMap := map[string]bool{}

	for _, v := range append(a, b...) {
		if _, ok := eMap[v]; !ok {
			eMap[v] = true
			out = append(out, v)
		}
	}

	return out
}

func (e Environments) setOutputs(o Outputter) error {
	return setJSONOutput(o, "workspaces", e)
}
