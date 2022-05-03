package action

import (
	"errors"
)

var (
	// Error returned when a workspaced input is set for an environment not passed in the environments input
	ErrEnvironmentNotFound = errors.New("Environment not found")
)

// Config holds the parsed workspace values
type Config struct {
	//EnvironmentsVariables hold variables for specific environments
	EnvironmentsVariables EnvironmentsVariables
	// EnvironmentsTags hold tags for specific environments
	EnvironmentsTags EnvironmentsTags
	// Environments entires indicate that one workspace should be created per environment
	Environments Environments
	// Name represents the workspace name, or a workspace prefix in a multi environment configuration
	Name string
	// Tags are universally applied to all managed workspaces
	Tags Tags
}

// NewConfig returns an empty Config struct
func NewConfig(name string) Config {
	return Config{
		Environments:          Environments{},
		EnvironmentsVariables: map[string][]Variable{},
		EnvironmentsTags:      EnvironmentsTags{},
		Tags:                  Tags{},
		Name:                  name,
	}
}

// ExtendConfig takes two Config structs and extends the values in a with values from b, removing duplicates
func ExtendConfig(a Config, b Config) Config {
	envs := MergeEnvironments(a.Environments, b.Environments)

	return Config{
		Name:                  a.Name,
		Environments:          envs,
		Tags:                  MergeTags(a.Tags, b.Tags),
		EnvironmentsTags:      MergeEnvironmentsTags(a.EnvironmentsTags, b.EnvironmentsTags),
		EnvironmentsVariables: MergeEnvironmentsVariables(a.EnvironmentsVariables, b.EnvironmentsVariables),
	}
}
