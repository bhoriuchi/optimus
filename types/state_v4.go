package types

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type StateV4 struct {
	Version          StateVersionV4           `json:"version"`
	TerraformVersion string                   `json:"terraform_version"`
	Serial           uint64                   `json:"serial"`
	Lineage          string                   `json:"lineage"`
	RootOutputs      map[string]OutputStateV4 `json:"outputs"`
	Resources        []*ResourceStateV4       `json:"resources"`
}

type OutputStateV4 struct {
	ValueRaw     json.RawMessage `json:"value"`
	ValueTypeRaw json.RawMessage `json:"type"`
	Sensitive    bool            `json:"sensitive,omitempty"`
}

type ResourceStateV4 struct {
	Module         string                   `json:"module,omitempty" yaml:"module,omitempty"`
	Mode           string                   `json:"mode" yaml:"mode"`
	Type           string                   `json:"type" yaml:"type"`
	Name           string                   `json:"name" yaml:"name"`
	EachMode       string                   `json:"each,omitempty" yaml:"each,omitempty"`
	ProviderConfig string                   `json:"provider" yaml:"provider"`
	Instances      []*InstanceObjectStateV4 `json:"instances" yaml:"instances"`
}

func (r ResourceStateV4) Addr(indexKey ...interface{}) string {
	path := []string{r.Module}

	if r.Mode == "data" {
		path = append(path, "data")
	}

	path = append(path, r.Type, r.Name)
	out := strings.Join(path, ".")

	if r.Mode == "data" {
		return out
	}

	if len(indexKey) == 1 {
		switch key := indexKey[0].(type) {
		case int, int32, int64, float32, float64:
			return fmt.Sprintf("%s[%v]", out, key)
		case string:
			return fmt.Sprintf(`%s["%v"]`, out, key)
		}
	}

	return out
}

type InstanceObjectStateV4 struct {
	IndexKey interface{} `json:"index_key,omitempty"`
	Status   string      `json:"status,omitempty"`
	Deposed  string      `json:"deposed,omitempty"`

	SchemaVersion           uint64            `json:"schema_version"`
	AttributesRaw           json.RawMessage   `json:"attributes,omitempty"`
	AttributesFlat          map[string]string `json:"attributes_flat,omitempty"`
	AttributeSensitivePaths json.RawMessage   `json:"sensitive_attributes,omitempty,"`

	PrivateRaw []byte `json:"private,omitempty"`

	Dependencies []string `json:"dependencies,omitempty"`

	CreateBeforeDestroy bool `json:"create_before_destroy,omitempty"`
}

// convert attributes to a map
func (i *InstanceObjectStateV4) Attrs() (map[string]interface{}, error) {
	attrs := map[string]interface{}{}
	if i.AttributesRaw != nil {
		if err := json.Unmarshal(i.AttributesRaw, &attrs); err != nil {
			return nil, err
		}
	}

	return attrs, nil
}

// stateVersionV4 is a weird special type we use to produce our hard-coded
// "version": 4 in the JSON serialization.
type StateVersionV4 struct{}

func (sv StateVersionV4) MarshalJSON() ([]byte, error) {
	return []byte{'4'}, nil
}

func (sv StateVersionV4) UnmarshalJSON([]byte) error {
	// Nothing to do: we already know we're version 4
	return nil
}

// normalize makes some in-place changes to normalize the way items are
// stored to ensure that two functionally-equivalent states will be stored
// identically.
func (s *StateV4) Normalize() {
	sort.Stable(SortResourcesV4(s.Resources))
	for _, rs := range s.Resources {
		sort.Stable(SortInstancesV4(rs.Instances))
	}
}

type SortResourcesV4 []*ResourceStateV4

func (sr SortResourcesV4) Len() int      { return len(sr) }
func (sr SortResourcesV4) Swap(i, j int) { sr[i], sr[j] = sr[j], sr[i] }
func (sr SortResourcesV4) Less(i, j int) bool {
	switch {
	case sr[i].Module != sr[j].Module:
		return sr[i].Module < sr[j].Module
	case sr[i].Mode != sr[j].Mode:
		return sr[i].Mode < sr[j].Mode
	case sr[i].Type != sr[j].Type:
		return sr[i].Type < sr[j].Type
	case sr[i].Name != sr[j].Name:
		return sr[i].Name < sr[j].Name
	default:
		return false
	}
}

type SortInstancesV4 []*InstanceObjectStateV4

func (si SortInstancesV4) Len() int      { return len(si) }
func (si SortInstancesV4) Swap(i, j int) { si[i], si[j] = si[j], si[i] }
func (si SortInstancesV4) Less(i, j int) bool {
	ki := si[i].IndexKey
	kj := si[j].IndexKey
	if ki != kj {
		if (ki == nil) != (kj == nil) {
			return ki == nil
		}
		if kii, isInt := ki.(int); isInt {
			if kji, isInt := kj.(int); isInt {
				return kii < kji
			}
			return true
		}
		if kis, isStr := ki.(string); isStr {
			if kjs, isStr := kj.(string); isStr {
				return kis < kjs
			}
			return true
		}
	}
	if si[i].Deposed != si[j].Deposed {
		return si[i].Deposed < si[j].Deposed
	}
	return false
}
