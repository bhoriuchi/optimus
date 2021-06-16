package types

type ChangeFile struct {
	Operations []ChangeOperation `yaml:"operations" json:"operations"`
}

type ChangeOperation struct {
	Operation  string           `yaml:"operation" json:"operation"`
	Address    string           `yaml:"address" json:"address"`
	NewAddress string           `yaml:"new_address" json:"new_address"`
	Count      []ChangeCount    `yaml:"count" json:"count"`
	Resource   *ResourceStateV4 `yaml:"resource" json:"resource"`
	Validate   string           `yaml:"validate" json:"validate"`
}

type ChangeCount struct {
	Index    int    `yaml:"index" json:"index"`
	Key      string `yaml:"key" json:"key"`
	Validate string `yaml:"validate" json:"validate"`
}
