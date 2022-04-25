package main

import (
	"github.com/sethvargo/go-githubactions"
	"github.com/takescoop/compose-inputs/internal/action"
)

func main() {
	if err := action.Run(action.Inputs{
		Environments: githubactions.GetInput("environments"),
		Variables:    githubactions.GetInput("environments_variables"),
		Tags:         githubactions.GetInput("environments_tags"),
	}, githubactions.New()); err != nil {
		githubactions.Fatalf("Error: %s", err)
	}
}
