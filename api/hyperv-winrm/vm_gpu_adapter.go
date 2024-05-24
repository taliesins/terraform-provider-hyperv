package hyperv_winrm

import (
	"context"
	"encoding/json"
	"log"
	"text/template"

	"github.com/taliesins/terraform-provider-hyperv/api"
)

type createVmGpuAdapterArgs struct {
	VMName           string
	VmGpuAdapterJson string
}

var createVmGpuAdapterTemplate = template.Must(template.New("CreateVmGpuAdapter").Parse(`
$ErrorActionPreference = 'Stop'
$vmGpuAdapter = '{{.VmGpuAdapterJson}}' | ConvertFrom-Json


$NewVmGpuAdapterArgs = @{
    InstancePath=$vmGpuAdapter.InstancePath
}

@(Get-VM -Name '{{.VMName}}') | Add-VMGpuPartitionAdapter @NewVmGpuAdapterArgs

$SetVmGpuAdapterArgs = @{
    MinPartitionVRAM=$vmGpuAdapter.MinPartitionVRAM
    MaxPartitionVRAM=$vmGpuAdapter.MaxPartitionVRAM
    OptimalPartitionVRAM=$vmGpuAdapter.OptimalPartitionVRAM

    MinPartitionEncode=$vmGpuAdapter.MinPartitionEncode
    MaxPartitionEncode=$vmGpuAdapter.MaxPartitionEncode
    OptimalPartitionEncode=$vmGpuAdapter.OptimalPartitionEncode

    MinPartitionDecode=$vmGpuAdapter.MinPartitionDecode
    MaxPartitionDecode=$vmGpuAdapter.MaxPartitionDecode
    OptimalPartitionDecode=$vmGpuAdapter.OptimalPartitionDecode

    MinPartitionCompute=$vmGpuAdapter.MinPartitionCompute
    MaxPartitionCompute=$vmGpuAdapter.MaxPartitionCompute
    OptimalPartitionCompute=$vmGpuAdapter.OptimalPartitionCompute

}

@(Get-VM -Name '{{.VMName}}' | Get-VMGpuPartitionAdapter | ?{$_.InstancePath -eq $vmGpuAdapter.InstancePath }) | Set-VMGpuPartitionAdapter @SetVmGpuAdapterArgs

`))

func (c *ClientConfig) CreateVmGpuAdapter(
	ctx context.Context,
	vmName string,
	instancePath string,
	minPartitionVram int32,
	maxPartitionVram int32,
	optimalPartitionVram int32,
	minPartitionEncode int64,
	maxPartitionEncode int64,
	optimalPartitionEncode int64,
	minPartitionDecode int32,
	maxPartitionDecode int32,
	optimalPartitionDecode int32,
	minPartitionCompute int32,
	maxPartitionCompute int32,
	optimalPartitionCompute int32,
) (err error) {
	vmGpuAdapterJson, err := json.Marshal(api.VmGpuAdapter{
		VmName:                  vmName,
		InstancePath:            instancePath,
		MinPartitionVram:        minPartitionVram,
		MaxPartitionVram:        maxPartitionVram,
		OptimalPartitionVram:    optimalPartitionVram,
		MinPartitionEncode:      minPartitionEncode,
		MaxPartitionEncode:      maxPartitionEncode,
		OptimalPartitionEncode:  optimalPartitionEncode,
		MinPartitionDecode:      minPartitionDecode,
		MaxPartitionDecode:      maxPartitionDecode,
		OptimalPartitionDecode:  optimalPartitionDecode,
		MinPartitionCompute:     minPartitionCompute,
		MaxPartitionCompute:     maxPartitionCompute,
		OptimalPartitionCompute: optimalPartitionCompute,
	})

	if err != nil {
		return err
	}

	err = c.WinRmClient.RunFireAndForgetScript(ctx, createVmGpuAdapterTemplate, createVmGpuAdapterArgs{
		VMName:           vmName,
		VmGpuAdapterJson: string(vmGpuAdapterJson),
	})

	return err
}

type getVmGpuAdaptersArgs struct {
	VmName string
}

var getVmGpuAdaptersTemplate = template.Must(template.New("GetVmGpuAdapters").Parse(`
$vmGpuAdaptersObject = @(Get-VM -Name '{{.VmName}}*' | ?{$_.Name -eq '{{.VmName}}' } | Get-VMGpuPartitionAdapter | %{ @{
     InstancePath=$_.InstancePath;

     CurrentPartitionVRAM=$_.CurrentPartitionVRAM;
     MinPartitionVRAM=$_.MinPartitionVRAM;
     MaxPartitionVRAM=$_.MaxPartitionVRAM;
     OptimalPartitionVRAM=$_.OptimalPartitionVRAM;

     CurrentPartitionEncode=$_.CurrentPartitionEncode;
     MinPartitionEncode=$_.MinPartitionEncode;
     MaxPartitionEncode=$_.MaxPartitionEncode;
     OptimalPartitionEncode=$_.OptimalPartitionEncode;

     CurrentPartitionDecode=$_.CurrentPartitionDecode;
     MinPartitionDecode=$_.MinPartitionDecode;
     MaxPartitionDecode=$_.MaxPartitionDecode;
     OptimalPartitionDecode=$_.OptimalPartitionDecode;

     CurrentPartitionCompute=$_.CurrentPartitionCompute;
     MinPartitionCompute=$_.MinPartitionCompute;
     MaxPartitionCompute=$_.MaxPartitionCompute;
     OptimalPartitionCompute=$_.OptimalPartitionCompute;
	 
}})

if ($vmGpuAdaptersObject) {
    # unexpectedly, powershell replaces & wit the unicode representation, so it has to be replaced back
	$vmGpuAdapters = (ConvertTo-Json -InputObject $vmGpuAdaptersObject) -replace '\\u0026', '&'
	$vmGpuAdapters
} else {
	"[]"
}
`))

func (c *ClientConfig) GetVmGpuAdapters(ctx context.Context, vmName string) (result []api.VmGpuAdapter, err error) {
	result = make([]api.VmGpuAdapter, 0)

	err = c.WinRmClient.RunScriptWithResult(ctx, getVmGpuAdaptersTemplate, getVmGpuAdaptersArgs{
		VmName: vmName,
	}, &result)

	return result, err
}

type updateVmGpuAdapterArgs struct {
	VmName           string
	InstancePath     string
	VmGpuAdapterJson string
}

var updateVmGpuAdapterTemplate = template.Must(template.New("UpdateVmGpuAdapter").Parse(`
$ErrorActionPreference = 'Stop'
Import-Module Hyper-V
$vmGpuAdapter = '{{.VmGpuAdapterJson}}' | ConvertFrom-Json

$vmGpuAdaptersObject = @(Get-VM -Name '{{.VmName}}*' | ?{$_.Name -eq '{{.VmName}}' } | Get-VMGpuPartitionAdapter | ?{$_.InstancePath -eq '{{.InstancePath}}' })

if (!$vmGpuAdaptersObject){
	throw "VM gpu adapter does not exist - {{.InstancePath}}"
}


$SetVmGpuAdapterArgs = @{
    MinPartitionVRAM=$vmGpuAdapter.MinPartitionVRAM
    MaxPartitionVRAM=$vmGpuAdapter.MaxPartitionVRAM
    OptimalPartitionVRAM=$vmGpuAdapter.OptimalPartitionVRAM

    MinPartitionEncode=$vmGpuAdapter.MinPartitionEncode
    MaxPartitionEncode=$vmGpuAdapter.MaxPartitionEncode
    OptimalPartitionEncode=$vmGpuAdapter.OptimalPartitionEncode

    MinPartitionDecode=$vmGpuAdapter.MinPartitionDecode
    MaxPartitionDecode=$vmGpuAdapter.MaxPartitionDecode
    OptimalPartitionDecode=$vmGpuAdapter.OptimalPartitionDecode

    MinPartitionCompute=$vmGpuAdapter.MinPartitionCompute
    MaxPartitionCompute=$vmGpuAdapter.MaxPartitionCompute
    OptimalPartitionCompute=$vmGpuAdapter.OptimalPartitionCompute

}

@(Get-VM -Name '{{.VmName}}' | Get-VMGpuPartitionAdapter | ?{$_.InstancePath -eq $vmGpuAdapter.InstancePath }) | Set-VMGpuPartitionAdapter @SetVmGpuAdapterArgs
`))

func (c *ClientConfig) UpdateVmGpuAdapter(
	ctx context.Context,
	vmName string,
	instancePath string,
	minPartitionVram int32,
	maxPartitionVram int32,
	optimalPartitionVram int32,
	minPartitionEncode int64,
	maxPartitionEncode int64,
	optimalPartitionEncode int64,
	minPartitionDecode int32,
	maxPartitionDecode int32,
	optimalPartitionDecode int32,
	minPartitionCompute int32,
	maxPartitionCompute int32,
	optimalPartitionCompute int32,
) (err error) {
	vmGpuAdapterJson, err := json.Marshal(api.VmGpuAdapter{
		VmName:                  vmName,
		InstancePath:            instancePath,
		MinPartitionVram:        minPartitionVram,
		MaxPartitionVram:        maxPartitionVram,
		OptimalPartitionVram:    optimalPartitionVram,
		MinPartitionEncode:      minPartitionEncode,
		MaxPartitionEncode:      maxPartitionEncode,
		OptimalPartitionEncode:  optimalPartitionEncode,
		MinPartitionDecode:      minPartitionDecode,
		MaxPartitionDecode:      maxPartitionDecode,
		OptimalPartitionDecode:  optimalPartitionDecode,
		MinPartitionCompute:     minPartitionCompute,
		MaxPartitionCompute:     maxPartitionCompute,
		OptimalPartitionCompute: optimalPartitionCompute,
	})

	if err != nil {
		return err
	}

	err = c.WinRmClient.RunFireAndForgetScript(ctx, updateVmGpuAdapterTemplate, updateVmGpuAdapterArgs{
		VmName:           vmName,
		InstancePath:     instancePath,
		VmGpuAdapterJson: string(vmGpuAdapterJson),
	})

	return err
}

type deleteVmGpuAdapterArgs struct {
	VmName       string
	InstancePath string
}

var deleteVmGpuAdapterTemplate = template.Must(template.New("DeleteVmGpuAdapter").Parse(`
$ErrorActionPreference = 'Stop'

@(Get-VM -Name '{{.VmName}}*' | ?{$_.Name -eq '{{.VmName}}' } | Get-VMGpuPartitionAdapter | ?{$_.InstancePath -eq '{{.InstancePath}}' }) | Remove-VMGpuPartitionAdapter
`))

func (c *ClientConfig) DeleteVmGpuAdapter(ctx context.Context, vmName string, instancePath string) (err error) {
	err = c.WinRmClient.RunFireAndForgetScript(ctx, deleteVmGpuAdapterTemplate, deleteVmGpuAdapterArgs{
		VmName:       vmName,
		InstancePath: instancePath,
	})

	return err
}

func (c *ClientConfig) CreateOrUpdateVmGpuAdapters(ctx context.Context, vmName string, gpuAdapters []api.VmGpuAdapter) (err error) {
	currentGpuAdapters, err := c.GetVmGpuAdapters(ctx, vmName)
	if err != nil {
		return err
	}

	currentGpuAdaptersLength := len(currentGpuAdapters)
	desiredGpuAdaptersLength := len(gpuAdapters)

	log.Printf("currentGpuAdaptersLength: %d", currentGpuAdaptersLength)
	log.Printf("desiredGpuAdaptersLength: %d", desiredGpuAdaptersLength)

	for i := currentGpuAdaptersLength - 1; i > desiredGpuAdaptersLength-1; i-- {
		currentGpuAdapter := currentGpuAdapters[i]
		err = c.DeleteVmGpuAdapter(ctx, vmName, currentGpuAdapter.InstancePath)
		if err != nil {
			return err
		}
	}

	if currentGpuAdaptersLength > desiredGpuAdaptersLength {
		currentGpuAdaptersLength = desiredGpuAdaptersLength
	}

	for i := 0; i <= currentGpuAdaptersLength-1; i++ {
		currentGpuAdapter := currentGpuAdapters[i]
		gpuAdapter := gpuAdapters[i]

		err = c.UpdateVmGpuAdapter(
			ctx,
			vmName,
			currentGpuAdapter.InstancePath,
			gpuAdapter.MinPartitionVram,
			gpuAdapter.MaxPartitionVram,
			gpuAdapter.OptimalPartitionVram,
			gpuAdapter.MinPartitionEncode,
			gpuAdapter.MaxPartitionEncode,
			gpuAdapter.OptimalPartitionEncode,
			gpuAdapter.MinPartitionDecode,
			gpuAdapter.MaxPartitionDecode,
			gpuAdapter.OptimalPartitionDecode,
			gpuAdapter.MinPartitionCompute,
			gpuAdapter.MaxPartitionCompute,
			gpuAdapter.OptimalPartitionCompute,
		)
		if err != nil {
			return err
		}
	}

	for i := currentGpuAdaptersLength - 1 + 1; i <= desiredGpuAdaptersLength-1; i++ {
		gpuAdapter := gpuAdapters[i]
		err = c.CreateVmGpuAdapter(
			ctx,
			vmName,
			gpuAdapter.InstancePath,
			gpuAdapter.MinPartitionVram,
			gpuAdapter.MaxPartitionVram,
			gpuAdapter.OptimalPartitionVram,
			gpuAdapter.MinPartitionEncode,
			gpuAdapter.MaxPartitionEncode,
			gpuAdapter.OptimalPartitionEncode,
			gpuAdapter.MinPartitionDecode,
			gpuAdapter.MaxPartitionDecode,
			gpuAdapter.OptimalPartitionDecode,
			gpuAdapter.MinPartitionCompute,
			gpuAdapter.MaxPartitionCompute,
			gpuAdapter.OptimalPartitionCompute,
		)

		if err != nil {
			return err
		}
	}

	return nil
}
