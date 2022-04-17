package hyperv_winrm

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/taliesins/terraform-provider-hyperv/api"
	"text/template"
)

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

func (c *ClientConfig) CreateOrUpdateVmProcessor(
	ctx context.Context,
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
	vmProcessorJson, err := json.Marshal(api.VmProcessor{
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

	if err != nil {
		return err
	}

	err = c.WinRmClient.RunFireAndForgetScript(ctx, createOrUpdateVmProcessorTemplate, createOrUpdateVmProcessorArgs{
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

func (c *ClientConfig) GetVmProcessor(ctx context.Context, vmName string) (result api.VmProcessor, err error) {
	err = c.WinRmClient.RunScriptWithResult(ctx, getVmProcessorTemplate, getVmProcessorArgs{
		VmName: vmName,
	}, &result)

	return result, err
}

func (c *ClientConfig) GetVmProcessors(ctx context.Context, vmName string) (result []api.VmProcessor, err error) {
	result = make([]api.VmProcessor, 0)
	vmProcessor, err := c.GetVmProcessor(ctx, vmName)
	if err != nil {
		return result, err
	}
	result = append(result, vmProcessor)
	return result, err
}

func (c *ClientConfig) CreateOrUpdateVmProcessors(ctx context.Context, vmName string, vmProcessors []api.VmProcessor) (err error) {
	if len(vmProcessors) == 0 {
		return nil
	}
	if len(vmProcessors) > 1 {
		return fmt.Errorf("Only 1 vm processor setting allowed per a vm")
	}

	vmProcessor := vmProcessors[0]

	return c.CreateOrUpdateVmProcessor(ctx, vmName,
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
