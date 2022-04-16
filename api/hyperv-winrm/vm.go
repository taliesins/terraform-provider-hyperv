package hyperv_winrm

import (
	"encoding/json"
	"github.com/taliesins/terraform-provider-hyperv/api"
	"text/template"
)

type createVmArgs struct {
	VmJson string
}

var createVmTemplate = template.Must(template.New("CreateVm").Parse(`
$ErrorActionPreference = 'Stop'
Import-Module Hyper-V
$vm = '{{.VmJson}}' | ConvertFrom-Json
$automaticCriticalErrorAction = [Microsoft.HyperV.PowerShell.CriticalErrorAction]$vm.AutomaticCriticalErrorAction
$automaticStartAction = [Microsoft.HyperV.PowerShell.StartAction]$vm.AutomaticStartAction
$automaticStopAction = [Microsoft.HyperV.PowerShell.StopAction]$vm.AutomaticStopAction
$checkpointType = [Microsoft.HyperV.PowerShell.CheckpointType]$vm.CheckpointType
$lockOnDisconnect = [Microsoft.HyperV.PowerShell.OnOffState]$vm.LockOnDisconnect
$allowUnverifiedPaths = $true #Not a property set on the vm object, skips validation when changing path

$vmObject = Get-VM -Name "$($vm.Name)*" | ?{$_.Name -eq $vm.Name}

if ($vmObject){
	throw "VM already exists - $($vm.Name)"
}

$NewVmArgs = @{
	Name=$vm.Name
	Generation=$vm.Generation
	MemoryStartupBytes=$vm.MemoryStartupBytes
	NoVHD=$true
}

if ($vm.Path) {
	$NewVmArgs.Path = $vm.Path
}

New-Vm @NewVmArgs

#Delete any auto-generated network adapter
Get-VMNetworkAdapter -VmName $vm.Name | Remove-VMNetworkAdapter

#Delete any auto-generated dvd drive
Get-VMDvdDrive -VmName $vm.Name | Remove-VMDvdDrive

#Set static and dynamic properties can't be set at the same time, but we need the values to match terraforms state
$SetVmArgs = @{}
$SetVmArgs.Name=$vm.Name
$SetVmArgs.StaticMemory=$true
$SetVmArgs.MemoryStartupBytes=$vm.MemoryStartupBytes
Set-Vm @SetVmArgs

$SetVmArgs = @{}
$SetVmArgs.Name=$vm.Name
$SetVmArgs.DynamicMemory=$true
$SetVmArgs.MemoryMinimumBytes=$vm.MemoryMinimumBytes
$SetVmArgs.MemoryMaximumBytes=$vm.MemoryMaximumBytes
Set-Vm @SetVmArgs

$SetVmArgs = @{}
$SetVmArgs.Name=$vm.Name
$SetVmArgs.GuestControlledCacheTypes=$vm.GuestControlledCacheTypes
$SetVmArgs.LowMemoryMappedIoSpace=$vm.LowMemoryMappedIoSpace
$SetVmArgs.HighMemoryMappedIoSpace=$vm.HighMemoryMappedIoSpace
$SetVmArgs.ProcessorCount=$vm.ProcessorCount
$SetVmArgs.AutomaticStartAction=$automaticStartAction
$SetVmArgs.AutomaticStopAction=$automaticStopAction
$SetVmArgs.AutomaticStartDelay=$vm.AutomaticStartDelay
$SetVmArgs.AutomaticCriticalErrorAction=$automaticCriticalErrorAction
$SetVmArgs.AutomaticCriticalErrorActionTimeout=$vm.AutomaticCriticalErrorActionTimeout
$SetVmArgs.LockOnDisconnect=$lockOnDisconnect
$SetVmArgs.Notes=$vm.Notes
$SetVmArgs.SnapshotFileLocation=$vm.SnapshotFileLocation
$SetVmArgs.SmartPagingFilePath=$vm.SmartPagingFilePath
$SetVmArgs.CheckpointType=$checkpointType
$SetVmArgs.AllowUnverifiedPaths=$allowUnverifiedPaths
if ($vm.StaticMemory) {
	$SetVmArgs.StaticMemory = $vm.StaticMemory
} else {
	$SetVmArgs.DynamicMemory = $vm.DynamicMemory
}

Set-Vm @SetVmArgs

`))

func (c *ClientConfig) CreateVm(
	name string,
	path string,
	generation int,
	automaticCriticalErrorAction api.CriticalErrorAction,
	automaticCriticalErrorActionTimeout int32,
	automaticStartAction api.StartAction,
	automaticStartDelay int32,
	automaticStopAction api.StopAction,
	checkpointType api.CheckpointType,
	dynamicMemory bool,
	guestControlledCacheTypes bool,
	highMemoryMappedIoSpace int64,
	lockOnDisconnect api.OnOffState,
	lowMemoryMappedIoSpace int32,
	memoryMaximumBytes int64,
	memoryMinimumBytes int64,
	memoryStartupBytes int64,
	notes string,
	processorCount int64,
	smartPagingFilePath string,
	snapshotFileLocation string,
	staticMemory bool,
) (err error) {
	vmJson, err := json.Marshal(api.Vm{
		Name:                                name,
		Path:                                path,
		Generation:                          generation,
		AutomaticCriticalErrorAction:        automaticCriticalErrorAction,
		AutomaticCriticalErrorActionTimeout: automaticCriticalErrorActionTimeout,
		AutomaticStartAction:                automaticStartAction,
		AutomaticStartDelay:                 automaticStartDelay,
		AutomaticStopAction:                 automaticStopAction,
		CheckpointType:                      checkpointType,
		DynamicMemory:                       dynamicMemory,
		GuestControlledCacheTypes:           guestControlledCacheTypes,
		HighMemoryMappedIoSpace:             highMemoryMappedIoSpace,
		LockOnDisconnect:                    lockOnDisconnect,
		LowMemoryMappedIoSpace:              lowMemoryMappedIoSpace,
		MemoryMaximumBytes:                  memoryMaximumBytes,
		MemoryMinimumBytes:                  memoryMinimumBytes,
		MemoryStartupBytes:                  memoryStartupBytes,
		Notes:                               notes,
		ProcessorCount:                      processorCount,
		SmartPagingFilePath:                 smartPagingFilePath,
		SnapshotFileLocation:                snapshotFileLocation,
		StaticMemory:                        staticMemory,
	})

	if err != nil {
		return err
	}

	err = c.WinRmClient.RunFireAndForgetScript(createVmTemplate, createVmArgs{
		VmJson: string(vmJson),
	})

	return err
}

type getVmArgs struct {
	Name string
}

var getVmTemplate = template.Must(template.New("GetVm").Parse(`
$ErrorActionPreference = 'Stop'
$vmObject = Get-VM -Name '{{.Name}}*' | ?{$_.Name -eq '{{.Name}}' } | %{ @{
	Name=$_.Name;
	Path=$_.Path;
	Generation=$_.Generation;
	AutomaticCriticalErrorAction=$_.AutomaticCriticalErrorAction;
	AutomaticCriticalErrorActionTimeout=$_.AutomaticCriticalErrorActionTimeout;
	AutomaticStartAction=$_.AutomaticStartAction;
	AutomaticStartDelay=$_.AutomaticStartDelay;
	AutomaticStopAction=$_.AutomaticStopAction;
	CheckpointType=$_.CheckpointType;
	DynamicMemory=$_.DynamicMemoryEnabled;
	GuestControlledCacheTypes=$_.GuestControlledCacheTypes;
	HighMemoryMappedIoSpace=$_.HighMemoryMappedIoSpace;
	LockOnDisconnect=$_.LockOnDisconnect;
	LowMemoryMappedIoSpace=$_.LowMemoryMappedIoSpace;
	MemoryMaximumBytes=$_.MemoryMaximum;
	MemoryMinimumBytes=$_.MemoryMinimum;
	MemoryStartupBytes=$_.MemoryStartup;
	Notes=$_.Notes;
	ProcessorCount=$_.ProcessorCount;
	SmartPagingFilePath=$_.SmartPagingFilePath;
	SnapshotFileLocation=$_.SnapshotFileLocation;
	StaticMemory=!$_.DynamicMemoryEnabled;
}}

if ($vmObject) {
	$vm = ConvertTo-Json -InputObject $vmObject
	$vm
} else {
	"{}"
}
`))

func (c *ClientConfig) GetVm(name string) (result api.Vm, err error) {
	err = c.WinRmClient.RunScriptWithResult(getVmTemplate, getVmArgs{
		Name: name,
	}, &result)

	return result, err
}

type updateVmArgs struct {
	VmJson string
}

var updateVmTemplate = template.Must(template.New("UpdateVm").Parse(`
$ErrorActionPreference = 'Stop'
Import-Module Hyper-V
$vm = '{{.VmJson}}' | ConvertFrom-Json
$automaticCriticalErrorAction = [Microsoft.HyperV.PowerShell.CriticalErrorAction]$vm.AutomaticCriticalErrorAction
$automaticStartAction = [Microsoft.HyperV.PowerShell.StartAction]$vm.AutomaticStartAction
$automaticStopAction = [Microsoft.HyperV.PowerShell.StopAction]$vm.AutomaticStopAction
$checkpointType = [Microsoft.HyperV.PowerShell.CheckpointType]$vm.CheckpointType
$lockOnDisconnect = [Microsoft.HyperV.PowerShell.OnOffState]$vm.LockOnDisconnect
$allowUnverifiedPaths = $true #Not a property set on the vm object, skips validation when changing path
$vmObject = Get-VM -Name "$($vm.Name)*" | ?{$_.Name -eq $vm.Name}

if (!$vmObject){
	throw "VM does not exist - $($vm.Name)"
}

#Set static and dynamic properties can't be set at the same time, but we need the values to match terraforms state
$SetVmArgs = @{}
$SetVmArgs.Name=$vm.Name
$SetVmArgs.StaticMemory=$true
$SetVmArgs.MemoryStartupBytes=$vm.MemoryStartupBytes
Set-Vm @SetVmArgs

$SetVmArgs = @{}
$SetVmArgs.Name=$vm.Name
$SetVmArgs.DynamicMemory=$true
$SetVmArgs.MemoryMinimumBytes=$vm.MemoryMinimumBytes
$SetVmArgs.MemoryMaximumBytes=$vm.MemoryMaximumBytes
Set-Vm @SetVmArgs

$SetVmArgs = @{}
$SetVmArgs.Name=$vm.Name
$SetVmArgs.GuestControlledCacheTypes=$vm.GuestControlledCacheTypes
$SetVmArgs.LowMemoryMappedIoSpace=$vm.LowMemoryMappedIoSpace
$SetVmArgs.HighMemoryMappedIoSpace=$vm.HighMemoryMappedIoSpace
$SetVmArgs.ProcessorCount=$vm.ProcessorCount
$SetVmArgs.AutomaticStartAction=$automaticStartAction
$SetVmArgs.AutomaticStopAction=$automaticStopAction
$SetVmArgs.AutomaticStartDelay=$vm.AutomaticStartDelay
$SetVmArgs.AutomaticCriticalErrorAction=$automaticCriticalErrorAction
$SetVmArgs.AutomaticCriticalErrorActionTimeout=$vm.AutomaticCriticalErrorActionTimeout
$SetVmArgs.LockOnDisconnect=$lockOnDisconnect
$SetVmArgs.Notes=$vm.Notes
$SetVmArgs.SnapshotFileLocation=$vm.SnapshotFileLocation
$SetVmArgs.SmartPagingFilePath=$vm.SmartPagingFilePath
$SetVmArgs.CheckpointType=$checkpointType
$SetVmArgs.AllowUnverifiedPaths=$allowUnverifiedPaths
if ($vm.StaticMemory) {
	$SetVmArgs.StaticMemory = $vm.StaticMemory
} else {
	$SetVmArgs.DynamicMemory = $vm.DynamicMemory
}

Set-Vm @SetVmArgs
`))

func (c *ClientConfig) UpdateVm(
	name string,
	//	generation int,
	automaticCriticalErrorAction api.CriticalErrorAction,
	automaticCriticalErrorActionTimeout int32,
	automaticStartAction api.StartAction,
	automaticStartDelay int32,
	automaticStopAction api.StopAction,
	checkpointType api.CheckpointType,
	dynamicMemory bool,
	guestControlledCacheTypes bool,
	highMemoryMappedIoSpace int64,
	lockOnDisconnect api.OnOffState,
	lowMemoryMappedIoSpace int32,
	memoryMaximumBytes int64,
	memoryMinimumBytes int64,
	memoryStartupBytes int64,
	notes string,
	processorCount int64,
	smartPagingFilePath string,
	snapshotFileLocation string,
	staticMemory bool,
) (err error) {
	vmJson, err := json.Marshal(api.Vm{
		Name: name,
		//Generation:generation,
		AutomaticCriticalErrorAction:        automaticCriticalErrorAction,
		AutomaticCriticalErrorActionTimeout: automaticCriticalErrorActionTimeout,
		AutomaticStartAction:                automaticStartAction,
		AutomaticStartDelay:                 automaticStartDelay,
		AutomaticStopAction:                 automaticStopAction,
		CheckpointType:                      checkpointType,
		DynamicMemory:                       dynamicMemory,
		GuestControlledCacheTypes:           guestControlledCacheTypes,
		HighMemoryMappedIoSpace:             highMemoryMappedIoSpace,
		LockOnDisconnect:                    lockOnDisconnect,
		LowMemoryMappedIoSpace:              lowMemoryMappedIoSpace,
		MemoryMaximumBytes:                  memoryMaximumBytes,
		MemoryMinimumBytes:                  memoryMinimumBytes,
		MemoryStartupBytes:                  memoryStartupBytes,
		Notes:                               notes,
		ProcessorCount:                      processorCount,
		SmartPagingFilePath:                 smartPagingFilePath,
		SnapshotFileLocation:                snapshotFileLocation,
		StaticMemory:                        staticMemory,
	})

	if err != nil {
		return err
	}

	err = c.WinRmClient.RunFireAndForgetScript(updateVmTemplate, updateVmArgs{
		VmJson: string(vmJson),
	})

	return err
}

type deleteVmArgs struct {
	Name string
}

var deleteVmTemplate = template.Must(template.New("DeleteVm").Parse(`
$ErrorActionPreference = 'Stop'
Get-VM -Name '{{.Name}}*' | ?{$_.Name -eq '{{.Name}}'} | Remove-VM -force
`))

func (c *ClientConfig) DeleteVm(name string) (err error) {
	err = c.WinRmClient.RunFireAndForgetScript(deleteVmTemplate, deleteVmArgs{
		Name: name,
	})

	return err
}
