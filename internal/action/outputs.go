package action

import (
	"encoding/json"
	"fmt"
)

func (c Config) SetOutputs(o Outputter) error {
	if err := setJSONOutput(o, "workspaces", c.Environments); err != nil {
		return err
	}

	if err := setJSONOutput(o, "workspace_tags", c.EnvironmentsTags); err != nil {
		return err
	}

	for _, vars := range c.EnvironmentsVariables {
		for _, v := range vars {
			if v.Sensitive {
				o.AddMask(v.Value)
			}
		}
	}

	if err := setJSONOutput(o, "workspace_variables", c.EnvironmentsVariables); err != nil {
		return err
	}

	o.SetOutput("name", c.Name)

	if err := setJSONOutput(o, "tags", c.Tags); err != nil {
		return err
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
