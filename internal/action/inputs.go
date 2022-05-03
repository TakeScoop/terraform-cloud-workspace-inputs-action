package action

import (
	"errors"
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

	wsVars, err := ParseEnvironmentsVariables(i.EnvironmentsVariables, environments)
	if err != nil {
		return Config{}, err
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
