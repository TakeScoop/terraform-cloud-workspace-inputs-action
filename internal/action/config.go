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
	Variables map[string][]Variable
	Tags      map[string][]string
	Names     []string
}

// NewConfig returns an empty Config struct
func NewConfig() Config {
	return Config{
		Names:     []string{},
		Variables: map[string][]Variable{},
		Tags:      map[string][]string{},
	}
}

// MergeConfigs takes two Config structs and merges the values, removing duplicates
func MergeConfigs(a Config, b Config) Config {
	merged := NewConfig()

	envMap := map[string]bool{}
	for _, e := range append(a.Names, b.Names...) {
		if _, ok := envMap[e]; !ok {
			merged.Names = append(merged.Names, e)
			envMap[e] = true
		}
	}

	for _, e := range merged.Names {
		merged.Tags[e] = []string{}

		tagMap := map[string]bool{}

		if aTags, ok := a.Tags[e]; ok {
			for _, t := range aTags {
				if _, ok := tagMap[t]; !ok {
					merged.Tags[e] = append(merged.Tags[e], t)
					tagMap[t] = true
				}
			}
		}

		if bTags, ok := b.Tags[e]; ok {
			for _, t := range bTags {
				if _, ok := tagMap[t]; !ok {
					merged.Tags[e] = append(merged.Tags[e], t)
					tagMap[t] = true
				}
			}
		}
	}

	for _, e := range merged.Names {
		merged.Variables[e] = []Variable{}

		varMap := map[string]Variable{}

		if aVars, ok := a.Variables[e]; ok {
			for _, v := range aVars {
				if _, ok := varMap[fmt.Sprintf("%s-%s", v.Key, v.Category)]; !ok {
					merged.Variables[e] = append(merged.Variables[e], v)
					varMap[fmt.Sprintf("%s-%s", v.Key, v.Category)] = v
				}
			}
		}

		if bVars, ok := b.Variables[e]; ok {
			for _, v := range bVars {
				if _, ok := varMap[fmt.Sprintf("%s-%s", v.Key, v.Category)]; !ok {
					merged.Variables[e] = append(merged.Variables[e], v)
					varMap[fmt.Sprintf("%s-%s", v.Key, v.Category)] = v
				}
			}
		}
	}

	return merged
}
