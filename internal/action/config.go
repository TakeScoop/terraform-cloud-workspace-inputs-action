package action

import (
	"errors"
	"fmt"
)

var (
	// Error returned when a workspaced input is set for an environment not passed in the environments input
	ErrEnvironmentNotFound = errors.New("Environment not found")
)

// Config holds the parsed workspace values
type Config struct {
	EnvironmentsVariables map[string][]Variable
	EnvironmentsTags      map[string][]string
	Environments          []string
}

// NewConfig returns an empty Config struct
func NewConfig() Config {
	return Config{
		Environments:          []string{},
		EnvironmentsVariables: map[string][]Variable{},
		EnvironmentsTags:      map[string][]string{},
	}
}

// MergeConfigs takes two Config structs and merges the values, removing duplicates
func MergeConfigs(a Config, b Config) Config {
	merged := NewConfig()

	envMap := map[string]bool{}
	for _, e := range append(a.Environments, b.Environments...) {
		if _, ok := envMap[e]; !ok {
			merged.Environments = append(merged.Environments, e)
			envMap[e] = true
		}
	}

	for _, e := range merged.Environments {
		merged.EnvironmentsTags[e] = []string{}

		tagMap := map[string]bool{}

		if aTags, ok := a.EnvironmentsTags[e]; ok {
			for _, t := range aTags {
				if _, ok := tagMap[t]; !ok {
					merged.EnvironmentsTags[e] = append(merged.EnvironmentsTags[e], t)
					tagMap[t] = true
				}
			}
		}

		if bTags, ok := b.EnvironmentsTags[e]; ok {
			for _, t := range bTags {
				if _, ok := tagMap[t]; !ok {
					merged.EnvironmentsTags[e] = append(merged.EnvironmentsTags[e], t)
					tagMap[t] = true
				}
			}
		}
	}

	for _, e := range merged.Environments {
		merged.EnvironmentsVariables[e] = []Variable{}

		varMap := map[string]Variable{}

		if aVars, ok := a.EnvironmentsVariables[e]; ok {
			for _, v := range aVars {
				if _, ok := varMap[fmt.Sprintf("%s-%s", v.Key, v.Category)]; !ok {
					merged.EnvironmentsVariables[e] = append(merged.EnvironmentsVariables[e], v)
					varMap[fmt.Sprintf("%s-%s", v.Key, v.Category)] = v
				}
			}
		}

		if bVars, ok := b.EnvironmentsVariables[e]; ok {
			for _, v := range bVars {
				if _, ok := varMap[fmt.Sprintf("%s-%s", v.Key, v.Category)]; !ok {
					merged.EnvironmentsVariables[e] = append(merged.EnvironmentsVariables[e], v)
					varMap[fmt.Sprintf("%s-%s", v.Key, v.Category)] = v
				}
			}
		}
	}

	return merged
}
