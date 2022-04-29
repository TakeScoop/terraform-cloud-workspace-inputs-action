package action

import "fmt"

func NewDefaults(environments []string, name string) Config {
	defaults := NewConfig(name)

	for _, e := range environments {
		defaults.Environments = append(defaults.Environments, e)
		defaults.EnvironmentsTags[e] = []string{fmt.Sprintf("environment:%s", e)}
		defaults.EnvironmentsVariables[e] = []Variable{{Key: "environment", Value: e, Category: "terraform"}}
	}

	defaults.Tags = append(defaults.Tags, fmt.Sprintf("source:%s", name))

	return defaults
}
