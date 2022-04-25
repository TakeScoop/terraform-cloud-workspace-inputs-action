package action

type Variable struct {
	Key         string `yaml:"key" json:"key"`
	Value       string `yaml:"value" json:"value"`
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
	Category    string `yaml:"category" json:"category"`
	Sensitive   bool   `yaml:"sensitive,omitempty" json:"sensitive,omitempty"`
}
