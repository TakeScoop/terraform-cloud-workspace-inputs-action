package action

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Environments []string

func ParseEnvironments(s string) (Environments, error) {
	var environments Environments

	if err := yaml.Unmarshal([]byte(s), &environments); err != nil {
		return Environments{}, fmt.Errorf("failed to parse environments: %w", err)
	}

	return environments, nil
}

func MergeEnvironments(a Environments, b Environments) Environments {
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

func EnvironmentsFromKeys(a map[string]any, b map[string]any) Environments {
	e1 := Environments{}
	for k := range a {
		e1 = append(e1, k)
	}

	e2 := Environments{}
	for k := range b {
		e2 = append(e2, k)
	}

	return MergeEnvironments(e1, e2)
}
