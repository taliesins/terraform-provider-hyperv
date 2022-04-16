package api

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

func DefaultVmProcessors() (interface{}, error) {
	result := make([]VmProcessor, 0)
	vmProcessor := VmProcessor{
		CompatibilityForMigrationEnabled:             false,
		CompatibilityForOlderOperatingSystemsEnabled: false,
		HwThreadCountPerCore:                         0,
		Maximum:                                      100,
		Reserve:                                      0,
		RelativeWeight:                               100,
		MaximumCountPerNumaNode:                      0,
		MaximumCountPerNumaSocket:                    0,
		EnableHostResourceProtection:                 false,
		ExposeVirtualizationExtensions:               false,
	}

	result = append(result, vmProcessor)
	return result, nil
}

func DiffSuppressVmProcessorMaximumCountPerNumaNode(key, old, new string, d *schema.ResourceData) bool {
	log.Printf("[DEBUG] '[%s]' Comparing old value '[%v]' with new value '[%v]' ", key, old, new)
	if new == "0" {
		//We have not explicitly set a value, so allow any value as we are not tracking it
		return true
	}

	return new == old
}

func DiffSuppressVmProcessorMaximumCountPerNumaSocket(key, old, new string, d *schema.ResourceData) bool {
	log.Printf("[DEBUG] '[%s]' Comparing old value '[%v]' with new value '[%v]' ", key, old, new)
	if new == "0" {
		//We have not explicitly set a value, so allow any value as we are not tracking it
		return true
	}

	return new == old
}

func ExpandVmProcessors(d *schema.ResourceData) ([]VmProcessor, error) {
	expandedVmProcessors := make([]VmProcessor, 0)

	if v, ok := d.GetOk("vm_processor"); ok {
		vmProcessors := v.([]interface{})
		for _, processor := range vmProcessors {
			processor, ok := processor.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] vm_processor should be a Hash - was '%+v'", processor)
			}

			log.Printf("[DEBUG] processor =  [%+v]", processor)

			expandedVmProcessor := VmProcessor{
				CompatibilityForMigrationEnabled:             processor["compatibility_for_migration_enabled"].(bool),
				CompatibilityForOlderOperatingSystemsEnabled: processor["compatibility_for_older_operating_systems_enabled"].(bool),
				HwThreadCountPerCore:                         int64(processor["hw_thread_count_per_core"].(int)),
				Maximum:                                      int64(processor["maximum"].(int)),
				Reserve:                                      int64(processor["reserve"].(int)),
				RelativeWeight:                               int32(processor["relative_weight"].(int)),
				MaximumCountPerNumaNode:                      int32(processor["maximum_count_per_numa_node"].(int)),
				MaximumCountPerNumaSocket:                    int32(processor["maximum_count_per_numa_socket"].(int)),
				EnableHostResourceProtection:                 processor["enable_host_resource_protection"].(bool),
				ExposeVirtualizationExtensions:               processor["expose_virtualization_extensions"].(bool),
			}

			expandedVmProcessors = append(expandedVmProcessors, expandedVmProcessor)
		}
	}

	return expandedVmProcessors, nil
}

func FlattenVmProcessors(vmProcessors *[]VmProcessor) []interface{} {
	if vmProcessors == nil || len(*vmProcessors) < 1  {
		return nil
	}

	flattenedVmProcessors := make([]interface{}, 0)

	for _, vmProcessor := range *vmProcessors {
		flattenedVmProcessor := make(map[string]interface{})
		flattenedVmProcessor["compatibility_for_migration_enabled"] = vmProcessor.CompatibilityForMigrationEnabled
		flattenedVmProcessor["compatibility_for_older_operating_systems_enabled"] = vmProcessor.CompatibilityForOlderOperatingSystemsEnabled
		flattenedVmProcessor["hw_thread_count_per_core"] = vmProcessor.HwThreadCountPerCore
		flattenedVmProcessor["maximum"] = vmProcessor.Maximum
		flattenedVmProcessor["reserve"] = vmProcessor.Reserve
		flattenedVmProcessor["relative_weight"] = vmProcessor.RelativeWeight
		flattenedVmProcessor["maximum_count_per_numa_node"] = vmProcessor.MaximumCountPerNumaNode
		flattenedVmProcessor["maximum_count_per_numa_socket"] = vmProcessor.MaximumCountPerNumaSocket
		flattenedVmProcessor["enable_host_resource_protection"] = vmProcessor.EnableHostResourceProtection
		flattenedVmProcessor["expose_virtualization_extensions"] = vmProcessor.ExposeVirtualizationExtensions
		flattenedVmProcessors = append(flattenedVmProcessors, flattenedVmProcessor)
	}

	return flattenedVmProcessors
}

type VmProcessor struct {
	VmName                                       string
	CompatibilityForMigrationEnabled             bool
	CompatibilityForOlderOperatingSystemsEnabled bool
	HwThreadCountPerCore                         int64
	Maximum                                      int64
	Reserve                                      int64
	RelativeWeight                               int32
	MaximumCountPerNumaNode                      int32
	MaximumCountPerNumaSocket                    int32
	EnableHostResourceProtection                 bool
	ExposeVirtualizationExtensions               bool
}

type HypervVmProcessorClient interface {
	CreateOrUpdateVmProcessor(
		vmName string,
		compatibilityForMigrationEnabled bool,
		compatibilityForOlderOperatingSystemsEnabled bool,
		hwThreadCountPerCore int64,
		maximum int64,
		reserve int64,
		relativeWeight int32,
		maximumCountPerNumaNode int32,
		maximumCountPerNumaSocket int32,
		enableHostResourceProtection bool,
		exposeVirtualizationExtensions bool,
	) (err error)
	GetVmProcessors(vmName string) (result []VmProcessor, err error)
	CreateOrUpdateVmProcessors(vmName string, vmProcessors []VmProcessor) (err error)
}
