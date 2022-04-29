package action

import "fmt"

func NewDefaults(environments []string) Config {
	defaults := Config{
		Environments:          []string{},
		EnvironmentsVariables: map[string][]Variable{},
		EnvironmentsTags:      map[string][]string{},
	}

	for _, e := range environments {
		defaults.Environments = append(defaults.Environments, e)
		defaults.EnvironmentsTags[e] = []string{fmt.Sprintf("environment:%s", e)}
		defaults.EnvironmentsVariables[e] = []Variable{{Key: "environment", Value: e, Category: "terraform"}}
	}

	return defaults
}
