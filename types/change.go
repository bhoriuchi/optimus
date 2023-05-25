package types

type ReplaceType string

const (
	ReplaceTypePrefix ReplaceType = "prefix"
	ReplaceTypeExact  ReplaceType = "exact"
)

type ChangeFile struct {
	Operations []ChangeOperation `yaml:"operations" json:"operations"`
}

type ChangeOperation struct {
	Operation        string           `yaml:"operation" json:"operation"`
	Address          string           `yaml:"address" json:"address"`
	NewAddress       string           `yaml:"new_address" json:"new_address"`
	NewName          string           `yaml:"new_name" json:"new_name"`
	Count            []ChangeCount    `yaml:"count" json:"count"`
	Resource         *ResourceStateV4 `yaml:"resource" json:"resource"`
	Validate         string           `yaml:"validate" json:"validate"`
	TerraformVersion string           `yaml:"terraform_version" json:"terraform_version"`
	ReplaceType      ReplaceType      `yaml:"replace_type" json:"replace_type"`
	Replace          string           `yaml:"replace" json:"replace"`
	With             string           `yaml:"with" json:"with"`
}

type ChangeCount struct {
	Index    int    `yaml:"index" json:"index"`
	Key      string `yaml:"key" json:"key"`
	Validate string `yaml:"validate" json:"validate"`
}
