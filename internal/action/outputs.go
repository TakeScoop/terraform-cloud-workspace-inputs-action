package action

import (
	"encoding/json"
	"fmt"
)

type Outputs struct {
	Workspaces string
	Tags       string
	Variables  string
}

func (c Config) SetOutputs(o Outputter) error {
	if err := setOutput(o, "workspaces", c.Names); err != nil {
		return err
	}

	if err := setOutput(o, "workspace_tags", c.Tags); err != nil {
		return err
	}

	for _, vars := range c.Variables {
		for _, v := range vars {
			fmt.Println(v.Value, v.Sensitive)

			if v.Sensitive {
				o.AddMask(v.Value)
			}
		}
	}

	if err := setOutput(o, "workspace_variables", c.Variables); err != nil {
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
