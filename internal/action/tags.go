package action

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// Tags are workspace tags applied to all workspaces
type Tags []string

// EnvironmentsTags is a map of environment names to a list of tags applied to that environment's workspace
type EnvironmentsTags map[string][]string

// ParseTags parses the passed string into universal tags applied to all workspaces
func ParseTags(s string) (Tags, error) {
	var tags Tags

	if err := yaml.Unmarshal([]byte(s), &tags); err != nil {
		return Tags{}, fmt.Errorf("failed to parse tags: %w", err)
	}

	return tags, nil
}

// ParseEnvironmentsTags parses a string into a map of environments to workspace tags for that environment
func ParseEnvironmentsTags(s string, envs Environments) (EnvironmentsTags, error) {
	var wsTags EnvironmentsTags

	if err := yaml.Unmarshal([]byte(s), &wsTags); err != nil {
		return EnvironmentsTags{}, fmt.Errorf("failed to parse workspace tags: %w", err)
	}

	eMap := map[string]bool{}
	for _, e := range envs {
		eMap[e] = true
	}

	for env := range wsTags {
		if _, ok := eMap[env]; !ok {
			return EnvironmentsTags{}, fmt.Errorf("environment %s in passed variables not found in environments %v: %w", env, envs, ErrEnvironmentNotFound)
		}
	}

	return wsTags, nil
}

// MergeTags takes two Tags structs and merges them, removing duplicates
func MergeTags(a Tags, b Tags) Tags {
	out := Tags{}

	tMap := map[string]bool{}

	for _, t := range append(a, b...) {
		if _, ok := tMap[t]; !ok {
			tMap[t] = true
			out = append(out, t)
		}
	}

	return out
}

func (et EnvironmentsTags) environments() Environments {
	envs := Environments{}

	for e := range et {
		envs = append(envs, e)
	}

	return envs
}

// MergeEnvironmentsTags takes two EnvironmentsTags structs and merges them, removing duplicates
func MergeEnvironmentsTags(a EnvironmentsTags, b EnvironmentsTags) EnvironmentsTags {
	out := EnvironmentsTags{}

	envs := MergeEnvironments(a.environments(), b.environments())

	for _, e := range envs {
		out[e] = Tags{}

		tagMap := map[string]bool{}

		if aTags, ok := a[e]; ok {
			for _, t := range aTags {
				if _, ok := tagMap[t]; !ok {
					out[e] = append(out[e], t)
					tagMap[t] = true
				}
			}
		}

		if bTags, ok := b[e]; ok {
			for _, t := range bTags {
				if _, ok := tagMap[t]; !ok {
					out[e] = append(out[e], t)
					tagMap[t] = true
				}
			}
		}
	}

	return out
}

func (t Tags) SetOutputs(o Outputter) error {
	return setJSONOutput(o, "tags", t)
}

func (et EnvironmentsTags) SetOutputs(o Outputter) error {
	return setJSONOutput(o, "workspace_tags", et)
}
