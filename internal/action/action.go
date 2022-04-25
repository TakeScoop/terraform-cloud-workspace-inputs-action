package action

type Outputter interface {
	SetOutput(k string, v string)
	AddMask(p string)
}

func Run(inputs Inputs, out Outputter) error {
	config, err := inputs.Parse()
	if err != nil {
		return err
	}

	merged := MergeConfigs(config, NewDefaults(config.Names))

	if err := merged.SetOutputs(out); err != nil {
		return err
	}

	return nil
}
