package main

import (
	"github.com/sethvargo/go-githubactions"
	"github.com/takescoop/compose-inputs/internal/action"
)

func main() {
	if err := action.Run(action.Inputs{
		Environments:          githubactions.GetInput("environments"),
		EnvironmentsVariables: githubactions.GetInput("environments_variables"),
		EnvironmentsTags:      githubactions.GetInput("environments_tags"),
	}, githubactions.New()); err != nil {
		githubactions.Fatalf("Error: %s", err)
	}
}
