package api

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"text/template"
)

func DefaultVmProcessors() (interface{}, error) {
	result := make([]vmProcessor, 0)
	vmProcessor := vmProcessor{
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

func ExpandVmProcessors(d *schema.ResourceData) ([]vmProcessor, error) {
	expandedVmProcessors := make([]vmProcessor, 0)

	if v, ok := d.GetOk("vm_processor"); ok {
		vmProcessors := v.([]interface{})
		for _, processor := range vmProcessors {
			processor, ok := processor.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] vm_processor should be a Hash - was '%+v'", processor)
			}

			log.Printf("[DEBUG] processor =  [%+v]", processor)

			expandedVmProcessor := vmProcessor{
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

func FlattenVmProcessors(vmProcessors *[]vmProcessor) []interface{} {
	flattenedVmProcessors := make([]interface{}, 0)

	if vmProcessors != nil {
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
	}

	return flattenedVmProcessors
}

type vmProcessor struct {
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

type createOrUpdateVmProcessorArgs struct {
	VmProcessorJson string
}

var createOrUpdateVmProcessorTemplate = template.Must(template.New("CreateOrUpdateVmProcessor").Parse(`
$ErrorActionPreference = 'Stop'
Import-Module Hyper-V
$vmProcessor = '{{.VmProcessorJson}}' | ConvertFrom-Json

$SetVMProcessorArgs = @{}
$SetVMProcessorArgs.VMName=$vmProcessor.VmName
#$SetVMProcessorArgs.Count=$vmProcessor.ProcessorCount
$SetVMProcessorArgs.CompatibilityForMigrationEnabled=$vmProcessor.CompatibilityForMigrationEnabled
$SetVMProcessorArgs.CompatibilityForOlderOperatingSystemsEnabled=$vmProcessor.CompatibilityForOlderOperatingSystemsEnabled
$SetVMProcessorArgs.HwThreadCountPerCore=$vmProcessor.HwThreadCountPerCore
$SetVMProcessorArgs.Maximum=$vmProcessor.Maximum
$SetVMProcessorArgs.Reserve=$vmProcessor.Reserve
$SetVMProcessorArgs.RelativeWeight=$vmProcessor.RelativeWeight
if ($vmProcessor.MaximumCountPerNumaNode -eq 0){
	$vmProcessor.MaximumCountPerNumaNode = (Get-WmiObject -class Win32_ComputerSystem).numberoflogicalprocessors
}
$SetVMProcessorArgs.MaximumCountPerNumaNode=$vmProcessor.MaximumCountPerNumaNode
if ($vmProcessor.MaximumCountPerNumaSocket -eq 0){
	$vmProcessor.MaximumCountPerNumaSocket = (Get-WmiObject -class Win32_ComputerSystem).numberofprocessors
}
$SetVMProcessorArgs.MaximumCountPerNumaSocket=$vmProcessor.MaximumCountPerNumaSocket
$SetVMProcessorArgs.EnableHostResourceProtection=$vmProcessor.EnableHostResourceProtection
$SetVMProcessorArgs.ExposeVirtualizationExtensions=$vmProcessor.ExposeVirtualizationExtensions

Set-VMProcessor @SetVMProcessorArgs
`))

func (c *HypervClient) CreateOrUpdateVmProcessor(
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
) (err error) {
	vmProcessorJson, err := json.Marshal(vmProcessor{
		VmName:                           vmName,
		CompatibilityForMigrationEnabled: compatibilityForMigrationEnabled,
		CompatibilityForOlderOperatingSystemsEnabled: compatibilityForOlderOperatingSystemsEnabled,
		HwThreadCountPerCore:                         hwThreadCountPerCore,
		Maximum:                                      maximum,
		Reserve:                                      reserve,
		RelativeWeight:                               relativeWeight,
		MaximumCountPerNumaNode:                      maximumCountPerNumaNode,
		MaximumCountPerNumaSocket:                    maximumCountPerNumaSocket,
		EnableHostResourceProtection:                 enableHostResourceProtection,
		ExposeVirtualizationExtensions:               exposeVirtualizationExtensions,
	})

	err = c.runFireAndForgetScript(createOrUpdateVmProcessorTemplate, createOrUpdateVmProcessorArgs{
		VmProcessorJson: string(vmProcessorJson),
	})

	return err
}

type getVmProcessorArgs struct {
	VmName string
}

var getVmProcessorTemplate = template.Must(template.New("GetVmProcessor").Parse(`
$ErrorActionPreference = 'Stop'

$vmProcessorObject = Get-VMProcessor -VMName '{{.VmName}}' | %{ @{
	CompatibilityForMigrationEnabled=$_.CompatibilityForMigrationEnabled
	CompatibilityForOlderOperatingSystemsEnabled=$_.CompatibilityForOlderOperatingSystemsEnabled
	HwThreadCountPerCore=$_.HwThreadCountPerCore
	Maximum=$_.Maximum
	Reserve=$_.Reserve
	RelativeWeight=$_.RelativeWeight
	MaximumCountPerNumaNode=$_.MaximumCountPerNumaNode
	MaximumCountPerNumaSocket=$_.MaximumCountPerNumaSocket
	EnableHostResourceProtection=$_.EnableHostResourceProtection
	ExposeVirtualizationExtensions=$_.ExposeVirtualizationExtensions
}}

if ($vmProcessorObject) {
	$vmProcessor = ConvertTo-Json -InputObject $vmProcessorObject
	$vmProcessor
} else {
	"{}"
}
`))

func (c *HypervClient) GetVmProcessor(vmName string) (result vmProcessor, err error) {
	err = c.runScriptWithResult(getVmProcessorTemplate, getVmProcessorArgs{
		VmName: vmName,
	}, &result)

	return result, err
}

func (c *HypervClient) GetVmProcessors(vmName string) (result []vmProcessor, err error) {
	result = make([]vmProcessor, 0)
	vmProcessor, err := c.GetVmProcessor(vmName)
	if err != nil {
		return result, err
	}
	result = append(result, vmProcessor)
	return result, err
}

func (c *HypervClient) CreateOrUpdateVmProcessors(vmName string, vmProcessors []vmProcessor) (err error) {
	if len(vmProcessors) == 0 {
		return nil
	}
	if len(vmProcessors) > 1 {
		return fmt.Errorf("Only 1 vm processor setting allowed per a vm")
	}

	vmProcessor := vmProcessors[0]

	return c.CreateOrUpdateVmProcessor(vmName,
		vmProcessor.CompatibilityForMigrationEnabled,
		vmProcessor.CompatibilityForOlderOperatingSystemsEnabled,
		vmProcessor.HwThreadCountPerCore,
		vmProcessor.Maximum,
		vmProcessor.Reserve,
		vmProcessor.RelativeWeight,
		vmProcessor.MaximumCountPerNumaNode,
		vmProcessor.MaximumCountPerNumaSocket,
		vmProcessor.EnableHostResourceProtection,
		vmProcessor.ExposeVirtualizationExtensions)
}
