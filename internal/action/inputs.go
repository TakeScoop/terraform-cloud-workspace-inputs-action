package action

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Inputs struct {
	Environments          string
	EnvironmentsTags      string
	EnvironmentsVariables string
}

func (i Inputs) Parse() (Config, error) {
	var environments []string
	if err := yaml.Unmarshal([]byte(i.Environments), &environments); err != nil {
		return Config{}, fmt.Errorf("failed to parse Names: %w", err)
	}

	var wsVars map[string][]Variable
	if err := yaml.Unmarshal([]byte(i.EnvironmentsVariables), &wsVars); err != nil {
		return Config{}, fmt.Errorf("failed to parse workspace variables: %w", err)
	}

	for env := range wsVars {
		found := false
		for _, e := range environments {
			if e == env {
				found = true
			}
		}
		if !found {
			return Config{}, fmt.Errorf("environment %s in passed tags not found in environments %v: %w", env, environments, ErrEnvironmentNotFound)
		}
	}

	var wsTags map[string][]string
	if err := yaml.Unmarshal([]byte(i.EnvironmentsTags), &wsTags); err != nil {
		return Config{}, fmt.Errorf("failed to parse workspace tags: %w", err)
	}

	for env := range wsTags {
		found := false
		for _, e := range environments {
			if e == env {
				found = true
			}
		}
		if !found {
			return Config{}, fmt.Errorf("environment %s in passed variables not found in environments %v: %w", env, environments, ErrEnvironmentNotFound)
		}
	}

	return Config{
		Environments:          environments,
		EnvironmentsVariables: wsVars,
		EnvironmentsTags:      wsTags,
	}, nil
}
