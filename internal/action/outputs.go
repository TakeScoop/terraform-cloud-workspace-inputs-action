package action

import (
	"encoding/json"
	"fmt"
)

func (c Config) SetOutputs(o Outputter) error {
	o.SetOutput("name", c.Name)

	if err := c.Environments.setOutputs(o); err != nil {
		return fmt.Errorf("failed to set environments output: %w", err)
	}

	if err := c.EnvironmentsTags.setOutputs(o); err != nil {
		return fmt.Errorf("failed to set environment tags output: %w", err)
	}

	if err := c.EnvironmentsVariables.setOutputs(o); err != nil {
		return fmt.Errorf("failed to set environment variables output: %w", err)
	}

	if err := c.Tags.setOutputs(o); err != nil {
		return fmt.Errorf("failed to set tags output: %w", err)
	}

	return nil
}

func setJSONOutput(o Outputter, key string, item any) error {
	b, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal %s: %w", key, err)
	}

	o.SetOutput(key, string(b))

	return nil
}
