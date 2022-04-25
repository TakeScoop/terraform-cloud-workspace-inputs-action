package action

import "fmt"

func NewDefaults(environments []string) Config {
	defaults := Config{
		Names:     []string{},
		Variables: map[string][]Variable{},
		Tags:      map[string][]string{},
	}

	for _, e := range environments {
		defaults.Names = append(defaults.Names, e)
		defaults.Tags[e] = []string{fmt.Sprintf("environment:%s", e)}
		defaults.Variables[e] = []Variable{{Key: "environment", Value: e, Category: "terraform"}}
	}

	return defaults
}
