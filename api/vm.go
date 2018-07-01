package api

import (
	"encoding/json"
	"strings"
	"text/template"
	"bytes"
	"strconv"
)

type CriticalErrorAction int

const (
	CriticalErrorAction_None  CriticalErrorAction = 0
	CriticalErrorAction_Pause CriticalErrorAction = 1
)

var CriticalErrorAction_name = map[CriticalErrorAction]string{
	CriticalErrorAction_None:  "None",
	CriticalErrorAction_Pause: "Pause",
}

var CriticalErrorAction_value = map[string]CriticalErrorAction{
	"none":  CriticalErrorAction_None,
	"pause": CriticalErrorAction_Pause,
}

func (x CriticalErrorAction) String() string {
	return CriticalErrorAction_name[x]
}

func ToCriticalErrorAction(x string) CriticalErrorAction {
	if integerValue, err := strconv.Atoi(x); err == nil {
		return CriticalErrorAction(integerValue)
	}

	return CriticalErrorAction_value[strings.ToLower(x)]
}

func (d *CriticalErrorAction) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(d.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (d *CriticalErrorAction) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		var i int
		err2 := json.Unmarshal(b, &i)
		if err2 == nil {
			*d = CriticalErrorAction(i)
			return nil
		}

		return err
	}
	*d = ToCriticalErrorAction(s)
	return nil
}

type StartAction int

const (
	StartAction_Nothing        StartAction = 2
	StartAction_StartIfRunning StartAction = 3
	StartAction_Start          StartAction = 4
)

var StartAction_name = map[StartAction]string{
	StartAction_Nothing:        "Nothing",
	StartAction_StartIfRunning: "StartIfRunning",
	StartAction_Start:          "Start",
}

var StartAction_value = map[string]StartAction{
	"nothing":        StartAction_Nothing,
	"startifrunning": StartAction_StartIfRunning,
	"start":          StartAction_Start,
}

func (x StartAction) String() string {
	return StartAction_name[x]
}

func ToStartAction(x string) StartAction {
	if integerValue, err := strconv.Atoi(x); err == nil {
		return StartAction(integerValue)
	}

	return StartAction_value[strings.ToLower(x)]
}

func (d *StartAction) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(d.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (d *StartAction) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		var i int
		err2 := json.Unmarshal(b, &i)
		if err2 == nil {
			*d = StartAction(i)
			return nil
		}

		return err
	}
	*d = ToStartAction(s)
	return nil
}

type StopAction int

const (
	StopAction_TurnOff  StopAction = 2
	StopAction_Save     StopAction = 3
	StopAction_ShutDown StopAction = 4
)

var StopAction_name = map[StopAction]string{
	StopAction_TurnOff:  "TurnOff",
	StopAction_Save:     "Save",
	StopAction_ShutDown: "ShutDown",
}

var StopAction_value = map[string]StopAction{
	"turnoff":  StopAction_TurnOff,
	"save":     StopAction_Save,
	"shutdown": StopAction_ShutDown,
}

func (x StopAction) String() string {
	return StopAction_name[x]
}

func ToStopAction(x string) StopAction {
	if integerValue, err := strconv.Atoi(x); err == nil {
		return StopAction(integerValue)
	}
	return StopAction_value[strings.ToLower(x)]
}

func (d *StopAction) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(d.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (d *StopAction) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		var i int
		err2 := json.Unmarshal(b, &i)
		if err2 == nil {
			*d = StopAction(i)
			return nil
		}

		return err
	}
	*d = ToStopAction(s)
	return nil
}

type CheckpointType int

const (
	CheckpointType_Disabled       CheckpointType = 2
	CheckpointType_Production     CheckpointType = 3
	CheckpointType_ProductionOnly CheckpointType = 4
	CheckpointType_Standard       CheckpointType = 5
)

var CheckpointType_name = map[CheckpointType]string{
	CheckpointType_Disabled:       "Disabled",
	CheckpointType_Production:     "Production",
	CheckpointType_ProductionOnly: "ProductionOnly",
	CheckpointType_Standard:       "Standard",
}

var CheckpointType_value = map[string]CheckpointType{
	"disabled":       CheckpointType_Disabled,
	"production":     CheckpointType_Production,
	"productiononly": CheckpointType_ProductionOnly,
	"standard":       CheckpointType_Standard,
}

func (x CheckpointType) String() string {
	return CheckpointType_name[x]
}

func ToCheckpointType(x string) CheckpointType {
	if integerValue, err := strconv.Atoi(x); err == nil {
		return CheckpointType(integerValue)
	}
	return CheckpointType_value[strings.ToLower(x)]
}

func (d *CheckpointType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(d.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (d *CheckpointType) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		var i int
		err2 := json.Unmarshal(b, &i)
		if err2 == nil {
			*d = CheckpointType(i)
			return nil
		}

		return err
	}
	*d = ToCheckpointType(s)
	return nil
}

type OnOffState int

const (
	OnOffState_On  OnOffState = 0
	OnOffState_Off OnOffState = 1
)

var OnOffState_name = map[OnOffState]string{
	OnOffState_On:  "On",
	OnOffState_Off: "Off",
}

var OnOffState_value = map[string]OnOffState{
	"on":  OnOffState_On,
	"off": OnOffState_Off,
}

func (x OnOffState) String() string {
	return OnOffState_name[x]
}

func ToOnOffState(x string) OnOffState {
	if integerValue, err := strconv.Atoi(x); err == nil {
		return OnOffState(integerValue)
	}
	return OnOffState_value[strings.ToLower(x)]
}

func (d *OnOffState) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(d.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (d *OnOffState) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		var i int
		err2 := json.Unmarshal(b, &i)
		if err2 == nil {
			*d = OnOffState(i)
			return nil
		}

		return err
	}
	*d = ToOnOffState(s)
	return nil
}

type vm struct {
	Name                                string
	Generation                          int
	AutomaticCriticalErrorAction        CriticalErrorAction
	AutomaticCriticalErrorActionTimeout int32
	AutomaticStartAction                StartAction
	AutomaticStartDelay                 int32
	AutomaticStopAction                 StopAction
	CheckpointType                      CheckpointType
	DynamicMemory                       bool
	GuestControlledCacheTypes           bool
	HighMemoryMappedIoSpace             int64
	LockOnDisconnect                    OnOffState
	LowMemoryMappedIoSpace              int32
	MemoryMaximumBytes                  int64
	MemoryMinimumBytes                  int64
	MemoryStartupBytes                  int64
	Notes                               string
	ProcessorCount                      int64
	SmartPagingFilePath                 string
	SnapshotFileLocation                string
	StaticMemory                        bool
	//ParentCheckpointName				string  this will allow us to set the checkpoint to use
}

type createVMArgs struct {
	VmJson string
}

var createVMTemplate = template.Must(template.New("CreateVM").Parse(`
$ErrorActionPreference = 'Stop'
Get-Vm | Out-Null
$vm = '{{.VmJson}}' | ConvertFrom-Json
$automaticCriticalErrorAction = [Microsoft.HyperV.PowerShell.CriticalErrorAction]$vm.AutomaticCriticalErrorAction
$automaticStartAction = [Microsoft.HyperV.PowerShell.StartAction]$vm.AutomaticStartAction
$automaticStopAction = [Microsoft.HyperV.PowerShell.StopAction]$vm.AutomaticStopAction
$checkpointType = [Microsoft.HyperV.PowerShell.CheckpointType]$vm.CheckpointType
$lockOnDisconnect = [Microsoft.HyperV.PowerShell.OnOffState]$vm.LockOnDisconnect
$allowUnverifiedPaths = $true #Not a property set on the vm object, skips validation when changing path

$vmObject = Get-VM | ?{$_.Name -eq $vm.Name}

if ($vmObject){
	throw "VM already exists - $($vm.Name)"
}

$NewVmArgs = @{
	Name=$vm.Name
	Generation=$vm.Generation
	MemoryStartupBytes=$vm.MemoryStartupBytes
	NoVHD=$true
}

New-Vm @NewVmArgs

#Delete any auto-generated network adapter
Get-VMNetworkAdapter -VMName web_server | Remove-VMNetworkAdapter

#Delete any auto-generated dvd drive
Get-VMDvdDrive -VMName web_server | Remove-VMDvdDrive

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

func (c *HypervClient) CreateVM(
	name string,
	generation int,
	automaticCriticalErrorAction CriticalErrorAction,
	automaticCriticalErrorActionTimeout int32,
	automaticStartAction StartAction,
	automaticStartDelay int32,
	automaticStopAction StopAction,
	checkpointType CheckpointType,
	dynamicMemory bool,
	guestControlledCacheTypes bool,
	highMemoryMappedIoSpace int64,
	lockOnDisconnect OnOffState,
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

	vmJson, err := json.Marshal(vm{
		Name:                                name,
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

	err = c.runFireAndForgetScript(createVMTemplate, createVMArgs{
		VmJson: string(vmJson),
	})

	return err
}

type getVMArgs struct {
	Name string
}

var getVMTemplate = template.Must(template.New("GetVM").Parse(`
$ErrorActionPreference = 'Stop'
$vmObject = Get-VM | ?{$_.Name -eq '{{.Name}}' } | %{ @{
	Name=$_.Name;
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

func (c *HypervClient) GetVM(name string) (result vm, err error) {
	err = c.runScriptWithResult(getVMTemplate, getVMArgs{
		Name: name,
	}, &result)

	return result, err
}

type updateVMArgs struct {
	VmJson string
}

var updateVMTemplate = template.Must(template.New("UpdateVM").Parse(`
$ErrorActionPreference = 'Stop'
Get-Vm | Out-Null
$vm = '{{.VmJson}}' | ConvertFrom-Json
$automaticCriticalErrorAction = [Microsoft.HyperV.PowerShell.CriticalErrorAction]$vm.AutomaticCriticalErrorAction
$automaticStartAction = [Microsoft.HyperV.PowerShell.StartAction]$vm.AutomaticStartAction
$automaticStopAction = [Microsoft.HyperV.PowerShell.StopAction]$vm.AutomaticStopAction
$checkpointType = [Microsoft.HyperV.PowerShell.CheckpointType]$vm.CheckpointType
$lockOnDisconnect = [Microsoft.HyperV.PowerShell.OnOffState]$vm.LockOnDisconnect
$allowUnverifiedPaths = $true #Not a property set on the vm object, skips validation when changing path
$vmObject = Get-VM | ?{$_.Name -eq $vm.Name}

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

func (c *HypervClient) UpdateVM(
	name string,
	//	generation int,
	automaticCriticalErrorAction CriticalErrorAction,
	automaticCriticalErrorActionTimeout int32,
	automaticStartAction StartAction,
	automaticStartDelay int32,
	automaticStopAction StopAction,
	checkpointType CheckpointType,
	dynamicMemory bool,
	guestControlledCacheTypes bool,
	highMemoryMappedIoSpace int64,
	lockOnDisconnect OnOffState,
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

	vmJson, err := json.Marshal(vm{
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

	err = c.runFireAndForgetScript(updateVMTemplate, updateVMArgs{
		VmJson: string(vmJson),
	})

	return err
}

type deleteVMArgs struct {
	Name string
}

var deleteVMTemplate = template.Must(template.New("DeleteVM").Parse(`
$ErrorActionPreference = 'Stop'
Get-VM | ?{$_.Name -eq '{{.Name}}'} | Remove-VM -force
`))

func (c *HypervClient) DeleteVM(name string) (err error) {
	err = c.runFireAndForgetScript(deleteVMTemplate, deleteVMArgs{
		Name: name,
	})

	return err
}
