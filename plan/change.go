package plan

import (
	"fmt"
	"log"

	"github.com/bhoriuchi/optimus/types"
	"github.com/fatih/color"
)

// converts an object that uses a count to a map
func countToMap(in *Input, op types.ChangeOperation) {
	for _, resource := range in.state.Resources {
		if resource.Addr() != op.Address {
			continue
		}

		resource.EachMode = "map"

		for _, instance := range resource.Instances {
			indexKey := fmt.Sprintf("%v", instance.IndexKey)
			for _, change := range op.Count {
				if fmt.Sprintf("%d", change.Index) == indexKey {
					instance.IndexKey = change.Key

					if change.Validate != "" {
						attrs, err := instance.Attrs()
						if err != nil {
							log.Panic(err)
						}

						if valid := Validate(change.Validate, map[string]interface{}{
							"Instance":   instance,
							"Count":      change,
							"Attributes": attrs,
						}); !valid {
							color.Yellow("! %s failed validation", resource.Addr(instance.IndexKey))
						}
					}
				}
			}
		}
	}
}

// adds a new resource
func addResource(in *Input, op types.ChangeOperation) {
	in.state.Resources = append(in.state.Resources, op.Resource)
}

// removes a resource
func removeResource(in *Input, op types.ChangeOperation) {
	resources := []*types.ResourceStateV4{}
	for _, resource := range in.state.Resources {
		if resource.Addr() != op.Address {
			resources = append(resources, resource)
		}
	}
	in.state.Resources = resources
}

// moves an instance to a resource and updates it index
func moveInstance(in *Input, op types.ChangeOperation) {
	var target *types.ResourceStateV4

	sourceAddr, _ := splitAddress(op.Address)

	targetAddr, indexKey := splitAddress(op.NewAddress)
	if indexKey == nil {
		log.Panicf("invalid new address: %s", op.NewAddress)
	}

	// find the target resource
	for _, resource := range in.state.Resources {
		if resource.Addr() == targetAddr {
			target = resource
			break
		}
	}

	if target == nil {
		log.Panicf("target resource %q not found", targetAddr)
	}

	if len(target.Instances) == 0 {
		target.Instances = []*types.InstanceObjectStateV4{}
	}

	for _, resource := range in.state.Resources {
		if resource.Addr() != sourceAddr {
			continue
		}

		sourceInstances := []*types.InstanceObjectStateV4{}

		for _, instance := range resource.Instances {
			addr := resource.Addr(instance.IndexKey)
			if addr == op.Address {
				instance.IndexKey = indexKey
				target.Instances = append(target.Instances, instance)
			} else {
				sourceInstances = append(sourceInstances, instance)
			}
		}

		resource.Instances = sourceInstances
	}
}
