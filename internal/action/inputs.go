package action

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

type Inputs struct {
	Environments          string
	EnvironmentsTags      string
	EnvironmentsVariables string
	Name                  string
	Tags                  string
}

var (
	ErrNameNotSet = errors.New("The input 'name' is required and cannot be an empty string")
)

// Parse unmarshals the raw inputs stored in an Inputs struct and returns a parsed Config struct
func (i Inputs) Parse() (Config, error) {
	if i.Name == "" {
		return Config{}, ErrNameNotSet
	}

	environments, err := ParseEnvironments(i.Environments)
	if err != nil {
		return Config{}, err
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

	wsTags, err := ParseEnvironmentsTags(i.EnvironmentsTags, environments)
	if err != nil {
		return Config{}, err
	}

	tags, err := ParseTags(i.Tags)
	if err != nil {
		return Config{}, err
	}

	return Config{
		Name:                  i.Name,
		Tags:                  tags,
		Environments:          environments,
		EnvironmentsVariables: wsVars,
		EnvironmentsTags:      wsTags,
	}, nil
}
