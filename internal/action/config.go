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
	//EnvironmentsVariables hold variables for specific environments
	EnvironmentsVariables map[string][]Variable
	// EnvironmentsTags hold tags for specific environments
	EnvironmentsTags map[string][]string
	// Environments entires indicate that one workspace should be created per environment
	Environments Environments
	// Name represents the workspace name, or a workspace prefix in a multi environment configuration
	Name string
	// Tags are universally applied to all managed workspaces
	Tags []string
}

// NewConfig returns an empty Config struct
func NewConfig(name string) Config {
	return Config{
		Environments:          Environments{},
		EnvironmentsVariables: map[string][]Variable{},
		EnvironmentsTags:      map[string][]string{},
		Tags:                  []string{},
		Name:                  name,
	}
}

// ExtendConfig takes two Config structs and extends the values in a with values from b, removing duplicates
func ExtendConfig(a Config, b Config) Config {
	merged := NewConfig(a.Name)

	envs := MergeEnvironments(a.Environments, b.Environments)

	merged.Environments = envs

	for _, e := range envs {
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

	for _, e := range envs {
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

	tMap := map[string]bool{}
	for _, t := range append(a.Tags, b.Tags...) {
		if _, ok := tMap[t]; !ok {
			tMap[t] = true
			merged.Tags = append(merged.Tags, t)
		}
	}

	return merged
}
