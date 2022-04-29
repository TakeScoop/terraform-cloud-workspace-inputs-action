package action

import (
	"encoding/json"
	"fmt"
)

func (c Config) SetOutputs(o Outputter) error {
	if err := setOutput(o, "workspaces", c.Environments); err != nil {
		return err
	}

	if err := setOutput(o, "workspace_tags", c.EnvironmentsTags); err != nil {
		return err
	}

	for _, vars := range c.EnvironmentsVariables {
		for _, v := range vars {
			if v.Sensitive {
				o.AddMask(v.Value)
			}
		}
	}

	if err := setOutput(o, "workspace_variables", c.EnvironmentsVariables); err != nil {
		return err
	}

	return nil
}

func setOutput(o Outputter, key string, item any) error {
	b, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal %s: %w", key, err)
	}

	o.SetOutput(key, string(b))

	return nil
}
