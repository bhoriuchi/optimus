package plan

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/bhoriuchi/optimus/types"
	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
)

type Input struct {
	Plan        string
	plan        *types.Plan
	State       string
	state       *types.StateV4
	Change      string
	change      *types.ChangeFile
	Out         string
	HideUpdates bool
}

func (in *Input) absPaths() error {
	var err error
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	if in.Plan != "" && !filepath.IsAbs(in.Plan) {
		in.Plan, err = filepath.Abs(filepath.Join(wd, in.Plan))
		if err != nil {
			return err
		}
	}

	if in.State != "" && !filepath.IsAbs(in.State) {
		in.State, err = filepath.Abs(filepath.Join(wd, in.State))
		if err != nil {
			return err
		}
	}

	if in.Change != "" && !filepath.IsAbs(in.Change) {
		in.Change, err = filepath.Abs(filepath.Join(wd, in.Change))
		if err != nil {
			return err
		}
	}

	if in.Out != "" && !filepath.IsAbs(in.Out) {
		in.Out, err = filepath.Abs(filepath.Join(wd, in.Out))
		if err != nil {
			return err
		}
	}

	return nil
}

func Plan(in *Input) error {
	in.plan = &types.Plan{}
	in.state = &types.StateV4{}

	if err := in.absPaths(); err != nil {
		return err
	}

	if in.Plan == "" {
		return fmt.Errorf("no plan file specified")
	}
	if in.State == "" {
		return fmt.Errorf("no state file specified")
	}

	plan, err := ioutil.ReadFile(in.Plan)
	if err != nil {
		return fmt.Errorf("failed to read plan file: %s", err)
	}

	if err := json.Unmarshal(plan, in.plan); err != nil {
		return fmt.Errorf("failed to unmarshal plan file: %s", err)
	}

	state, err := ioutil.ReadFile(in.State)
	if err != nil {
		return fmt.Errorf("failed to read state file: %s", err)
	}

	if err := json.Unmarshal(state, in.state); err != nil {
		return fmt.Errorf("failed to unmarshal state file: %s", err)
	}

	if in.Change != "" {
		in.change = &types.ChangeFile{}
		change, err := ioutil.ReadFile(in.Change)
		if err != nil {
			return fmt.Errorf("failed to read change file: %s", err)
		}
		if err := yaml.Unmarshal(change, in.change); err != nil {
			return fmt.Errorf("failed to unmarshal change file: %s", err)
		}
	}

	if in.change != nil {
		applyChanges(in)
		in.state.Serial++
	}

	add := []string{}
	remove := []string{}

	for _, change := range in.plan.ResourceChanges {

		for _, action := range change.Change.Actions {
			switch action {
			case "create":
				add = append(add, change.Address)
			case "delete":
				remove = append(remove, change.Address)
			}
		}
	}

	sort.Strings(add)
	sort.Strings(remove)

	for _, addr := range remove {
		found, _, _ := findStateResourceByAddress(in.state, addr)
		if found {
			color.Red("- %s", addr)
		} else if !in.HideUpdates {
			color.Blue("* %s", addr)
		}
	}

	for _, addr := range add {
		found, _, _ := findStateResourceByAddress(in.state, addr)
		if !found {
			color.Green("+ %s", addr)
		} else if !in.HideUpdates {
			color.Blue("* %s", addr)
		}
	}

	if in.Out != "" {
		b, err := json.MarshalIndent(in.state, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to serialize updated state file: %s", err)
		}
		if err := ioutil.WriteFile(in.Out, b, 0755); err != nil {
			return fmt.Errorf("faild to write updated state file to disk: %s", err)
		}
	}

	return nil
}

func findStateResourceByAddress(state *types.StateV4, address string) (bool, *types.ResourceStateV4, *types.InstanceObjectStateV4) {
	for _, resource := range state.Resources {
		for _, instance := range resource.Instances {
			raddr := resource.Addr()
			iaddr := resource.Addr(instance.IndexKey)

			switch idx := instance.IndexKey.(type) {
			case int, int32, int64, float32, float64:
				if (idx == 0.0 && address == raddr) || (address == iaddr) {
					return true, resource, instance
				}
			case string:
				if address == iaddr {
					return true, resource, instance
				}
			}
		}
	}

	return false, nil, nil
}

func applyChanges(in *Input) {
	for _, op := range in.change.Operations {
		switch op.Operation {
		case "count_to_map":
			countToMap(in, op)
		case "add_resource":
			addResource(in, op)
		case "remove_resource":
			removeResource(in, op)
		case "move_instance":
			moveInstance(in, op)
		}
	}
}
