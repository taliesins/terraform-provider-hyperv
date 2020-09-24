package api

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
	"text/template"
)

type VMSwitchBandwidthMode int

const (
	VMSwitchBandwidthMode_Default  VMSwitchBandwidthMode = 0
	VMSwitchBandwidthMode_Weight   VMSwitchBandwidthMode = 1
	VMSwitchBandwidthMode_Absolute VMSwitchBandwidthMode = 2
	VMSwitchBandwidthMode_None     VMSwitchBandwidthMode = 3
)

var VMSwitchBandwidthMode_name = map[VMSwitchBandwidthMode]string{
	VMSwitchBandwidthMode_Default:  "Default",
	VMSwitchBandwidthMode_Weight:   "Weight",
	VMSwitchBandwidthMode_Absolute: "Absolute",
	VMSwitchBandwidthMode_None:     "None",
}

var VMSwitchBandwidthMode_value = map[string]VMSwitchBandwidthMode{
	"default":  VMSwitchBandwidthMode_Default,
	"weight":   VMSwitchBandwidthMode_Weight,
	"absolute": VMSwitchBandwidthMode_Absolute,
	"none":     VMSwitchBandwidthMode_None,
}

func (x VMSwitchBandwidthMode) String() string {
	return VMSwitchBandwidthMode_name[x]
}

func ToVMSwitchBandwidthMode(x string) VMSwitchBandwidthMode {
	if integerValue, err := strconv.Atoi(x); err == nil {
		return VMSwitchBandwidthMode(integerValue)
	}

	return VMSwitchBandwidthMode_value[strings.ToLower(x)]
}

func (d *VMSwitchBandwidthMode) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(d.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (d *VMSwitchBandwidthMode) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		var i int
		err2 := json.Unmarshal(b, &i)
		if err2 == nil {
			*d = VMSwitchBandwidthMode(i)
			return nil
		}

		return err
	}
	*d = ToVMSwitchBandwidthMode(s)
	return nil
}

type VMSwitchType int

const (
	VMSwitchType_Private  VMSwitchType = 0
	VMSwitchType_Internal VMSwitchType = 1
	VMSwitchType_External VMSwitchType = 2
)

var VMSwitchType_name = map[VMSwitchType]string{
	VMSwitchType_Private:  "Private",
	VMSwitchType_Internal: "Internal",
	VMSwitchType_External: "External",
}

var VMSwitchType_value = map[string]VMSwitchType{
	"private":  VMSwitchType_Private,
	"internal": VMSwitchType_Internal,
	"external": VMSwitchType_External,
}

func (x VMSwitchType) String() string {
	return VMSwitchType_name[x]
}

func ToVMSwitchType(x string) VMSwitchType {
	if integerValue, err := strconv.Atoi(x); err == nil {
		return VMSwitchType(integerValue)
	}

	return VMSwitchType_value[strings.ToLower(x)]
}

func (d *VMSwitchType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(d.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (d *VMSwitchType) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		var i int
		err2 := json.Unmarshal(b, &i)
		if err2 == nil {
			*d = VMSwitchType(i)
			return nil
		}

		return err
	}
	*d = ToVMSwitchType(s)
	return nil
}

type vmSwitch struct {
	Name                                string
	Notes                               string
	AllowManagementOS                   bool
	EmbeddedTeamingEnabled              bool
	IovEnabled                          bool
	PacketDirectEnabled                 bool
	BandwidthReservationMode            VMSwitchBandwidthMode
	SwitchType                          VMSwitchType
	NetAdapterNames                     []string
	DefaultFlowMinimumBandwidthAbsolute int64
	DefaultFlowMinimumBandwidthWeight   int64
	DefaultQueueVmmqEnabled             bool
	DefaultQueueVmmqQueuePairs          int32
	DefaultQueueVrssEnabled             bool
}

type createVMSwitchArgs struct {
	VmSwitchJson string
}

var createVMSwitchTemplate = template.Must(template.New("CreateVMSwitch").Parse(`
$ErrorActionPreference = 'Stop'
Import-Module Hyper-V
$vmSwitch = '{{.VmSwitchJson}}' | ConvertFrom-Json
$minimumBandwidthMode = [Microsoft.HyperV.PowerShell.VMSwitchBandwidthMode]$vmSwitch.BandwidthReservationMode
$switchType = [Microsoft.HyperV.PowerShell.VMSwitchType]$vmSwitch.SwitchType
$NetAdapterNames = @($vmSwitch.NetAdapterNames)
#when EnablePacketDirect=true it seems to throw an exception if EnableIov=true or EnableEmbeddedTeaming=true

$switchObject = Get-VMSwitch -Name "$($vmSwitch.Name)*" | ?{$_.Name -eq $vmSwitch.Name}

if ($switchObject){
	throw "Switch already exists - $($vmSwitch.Name)"
}

$NewVmSwitchArgs = @{}
$NewVmSwitchArgs.Name=$vmSwitch.Name
$NewVmSwitchArgs.MinimumBandwidthMode=$minimumBandwidthMode
$NewVmSwitchArgs.EnableEmbeddedTeaming=$vmSwitch.EmbeddedTeamingEnabled
$NewVmSwitchArgs.EnableIov=$vmSwitch.IovEnabled
$NewVmSwitchArgs.EnablePacketDirect=$vmSwitch.PacketDirectEnabled

if ($NetAdapterNames) {
	$NewVmSwitchArgs.AllowManagementOS=$vmSwitch.AllowManagementOS
	$NewVmSwitchArgs.NetAdapterName=$NetAdapterNames
} else {
	$NewVmSwitchArgs.SwitchType=$switchType
	#not used unless interface is specified
	#-AllowManagementOS $vmSwitch.AllowManagementOS
}
New-VMSwitch @NewVmSwitchArgs

$switchObject = Get-VMSwitch -Name "$($vmSwitch.Name)" | ?{$_.Name -eq $vmSwitch.Name}

if (!$switchObject){
	throw "Switch does not exist - $($vmSwitch.Name)"
}

$SetVmSwitchArgs = @{}
$SetVmSwitchArgs.Name=$vmSwitch.Name
$SetVmSwitchArgs.Notes=$vmSwitch.Notes
if (($minimumBandwidthMode -eq [Microsoft.HyperV.PowerShell.VMSwitchBandwidthMode]::Absolute) -and $switchObject.DefaultFlowMinimumBandwidthAbsolute -ne $vmSwitch.DefaultFlowMinimumBandwidthAbsolute) {
	$SetVmSwitchArgs.DefaultFlowMinimumBandwidthAbsolute=$vmSwitch.DefaultFlowMinimumBandwidthAbsolute
}
if ((($minimumBandwidthMode -eq [Microsoft.HyperV.PowerShell.VMSwitchBandwidthMode]::Weight) -or (($minimumBandwidthMode -eq [Microsoft.HyperV.PowerShell.VMSwitchBandwidthMode]::Default) -and (-not ($vmSwitch.IovEnabled)))) -and $switchObject.DefaultFlowMinimumBandwidthWeight -ne $vmSwitch.DefaultFlowMinimumBandwidthWeight) {
	$SetVmSwitchArgs.DefaultFlowMinimumBandwidthWeight=$vmSwitch.DefaultFlowMinimumBandwidthWeight
}
$SetVmSwitchArgs.DefaultQueueVmmqEnabled=$vmSwitch.DefaultQueueVmmqEnabled
$SetVmSwitchArgs.DefaultQueueVmmqQueuePairs=$vmSwitch.DefaultQueueVmmqQueuePairs
$SetVmSwitchArgs.DefaultQueueVrssEnabled=$vmSwitch.DefaultQueueVrssEnabled

Set-VMSwitch @SetVmSwitchArgs

`))

func (c *HypervClient) CreateVMSwitch(
	name string,
	notes string,
	allowManagementOS bool,
	embeddedTeamingEnabled bool,
	iovEnabled bool,
	packetDirectEnabled bool,
	bandwidthReservationMode VMSwitchBandwidthMode,
	switchType VMSwitchType,
	netAdapterNames []string,
	defaultFlowMinimumBandwidthAbsolute int64,
	defaultFlowMinimumBandwidthWeight int64,
	defaultQueueVmmqEnabled bool,
	defaultQueueVmmqQueuePairs int32,
	defaultQueueVrssEnabled bool,
) (err error) {

	vmSwitchJson, err := json.Marshal(vmSwitch{
		Name:                                name,
		Notes:                               notes,
		AllowManagementOS:                   allowManagementOS,
		EmbeddedTeamingEnabled:              embeddedTeamingEnabled,
		IovEnabled:                          iovEnabled,
		PacketDirectEnabled:                 packetDirectEnabled,
		BandwidthReservationMode:            bandwidthReservationMode,
		SwitchType:                          switchType,
		NetAdapterNames:                     netAdapterNames,
		DefaultFlowMinimumBandwidthAbsolute: defaultFlowMinimumBandwidthAbsolute,
		DefaultFlowMinimumBandwidthWeight:   defaultFlowMinimumBandwidthWeight,
		DefaultQueueVmmqEnabled:             defaultQueueVmmqEnabled,
		DefaultQueueVmmqQueuePairs:          defaultQueueVmmqQueuePairs,
		DefaultQueueVrssEnabled:             defaultQueueVrssEnabled,
	})

	if err != nil {
		return err
	}

	err = c.runFireAndForgetScript(createVMSwitchTemplate, createVMSwitchArgs{
		VmSwitchJson: string(vmSwitchJson),
	})

	return err
}

type getVMSwitchArgs struct {
	Name string
}

var getVMSwitchTemplate = template.Must(template.New("GetVMSwitch").Parse(`
$ErrorActionPreference = 'Stop'
$vmSwitchObject = Get-VMSwitch -Name '{{.Name}}*' | ?{$_.Name -eq '{{.Name}}' } | %{ @{
	Name=$_.Name;
	Notes=$_.Notes;
	AllowManagementOS=$_.AllowManagementOS;
	EmbeddedTeamingEnabled=$_.EmbeddedTeamingEnabled;
	IovEnabled=$_.IovEnabled;
	PacketDirectEnabled=$_.PacketDirectEnabled;
	BandwidthReservationMode=$_.BandwidthReservationMode;
	SwitchType=$_.SwitchType;
	NetAdapterNames=@(if($_.NetAdapterInterfaceDescriptions){@(Get-NetAdapter -InterfaceDescription $_.NetAdapterInterfaceDescriptions | %{$_.Name})});
	DefaultFlowMinimumBandwidthAbsolute=$_.DefaultFlowMinimumBandwidthAbsolute;
	DefaultFlowMinimumBandwidthWeight=$_.DefaultFlowMinimumBandwidthWeight;
	DefaultQueueVmmqEnabled=$_.DefaultQueueVmmqEnabledRequested;
	DefaultQueueVmmqQueuePairs=$_.DefaultQueueVmmqQueuePairsRequested;
	DefaultQueueVrssEnabled=$_.DefaultQueueVrssEnabledRequested;
}}

if ($vmSwitchObject){
	$vmSwitch = ConvertTo-Json -InputObject $vmSwitchObject
	$vmSwitch
} else {
	"{}"
}
`))

func (c *HypervClient) GetVMSwitch(name string) (result vmSwitch, err error) {
	err = c.runScriptWithResult(getVMSwitchTemplate, getVMSwitchArgs{
		Name: name,
	}, &result)

	return result, err
}

type updateVMSwitchArgs struct {
	VmSwitchJson string
}

var updateVMSwitchTemplate = template.Must(template.New("UpdateVMSwitch").Parse(`
$ErrorActionPreference = 'Stop'
Import-Module Hyper-V
$vmSwitch = '{{.VmSwitchJson}}' | ConvertFrom-Json
$minimumBandwidthMode = [Microsoft.HyperV.PowerShell.VMSwitchBandwidthMode]$vmSwitch.BandwidthReservationMode
$switchType = [Microsoft.HyperV.PowerShell.VMSwitchType]$vmSwitch.SwitchType
$NetAdapterNames = @($vmSwitch.NetAdapterNames)
#when EnablePacketDirect=true it seems to throw an exception if EnableIov=true or EnableEmbeddedTeaming=true

$switchObject = Get-VMSwitch -Name "$($vmSwitch.Name)*" | ?{$_.Name -eq $vmSwitch.Name}

if (!$switchObject){
	throw "Switch does not exist - $($vmSwitch.Name)"
}

$SetVmSwitchArgs = @{}
$SetVmSwitchArgs.Name=$vmSwitch.Name
$SetVmSwitchArgs.Notes=$vmSwitch.Notes
if ($NetAdapterNames) {
	$SetVmSwitchArgs.AllowManagementOS=$vmSwitch.AllowManagementOS
	$SetVmSwitchArgs.NetAdapterName=$NetAdapterNames
	#Updates not supported on:
	#-EnableEmbeddedTeaming $vmSwitch.EmbeddedTeamingEnabled
	#-EnableIov $vmSwitch.IovEnabled
	#-EnablePacketDirect $vmSwitch.PacketDirectEnabled
	#-MinimumBandwidthMode $minimumBandwidthMode
} else {
	$SetVmSwitchArgs.SwitchType=$switchType
	#Updates not supported on:
	#-EnableEmbeddedTeaming $vmSwitch.EmbeddedTeamingEnabled
	#-EnableIov $vmSwitch.IovEnabled
	#-EnablePacketDirect $vmSwitch.PacketDirectEnabled
	#-MinimumBandwidthMode $minimumBandwidthMode

	#not used unless interface is specified
	#-AllowManagementOS $vmSwitch.AllowManagementOS
}

if (($minimumBandwidthMode -eq [Microsoft.HyperV.PowerShell.VMSwitchBandwidthMode]::Absolute) -and $switchObject.DefaultFlowMinimumBandwidthAbsolute -ne $vmSwitch.DefaultFlowMinimumBandwidthAbsolute) {
	$SetVmSwitchArgs.DefaultFlowMinimumBandwidthAbsolute=$vmSwitch.DefaultFlowMinimumBandwidthAbsolute
}
if ((($minimumBandwidthMode -eq [Microsoft.HyperV.PowerShell.VMSwitchBandwidthMode]::Weight) -or (($minimumBandwidthMode -eq [Microsoft.HyperV.PowerShell.VMSwitchBandwidthMode]::Default) -and (-not ($vmSwitch.IovEnabled)))) -and $switchObject.DefaultFlowMinimumBandwidthWeight -ne $vmSwitch.DefaultFlowMinimumBandwidthWeight) {
	$SetVmSwitchArgs.DefaultFlowMinimumBandwidthWeight=$vmSwitch.DefaultFlowMinimumBandwidthWeight
}
$SetVmSwitchArgs.DefaultQueueVmmqEnabled=$vmSwitch.DefaultQueueVmmqEnabled
$SetVmSwitchArgs.DefaultQueueVmmqQueuePairs=$vmSwitch.DefaultQueueVmmqQueuePairs
$SetVmSwitchArgs.DefaultQueueVrssEnabled=$vmSwitch.DefaultQueueVrssEnabled

Set-VMSwitch @SetVmSwitchArgs
`))

func (c *HypervClient) UpdateVMSwitch(
	name string,
	notes string,
	allowManagementOS bool,
	//embeddedTeamingEnabled bool,
	//iovEnabled bool,
	//packetDirectEnabled bool,
	//bandwidthReservationMode VMSwitchBandwidthMode,
	switchType VMSwitchType,
	netAdapterNames []string,
	defaultFlowMinimumBandwidthAbsolute int64,
	defaultFlowMinimumBandwidthWeight int64,
	defaultQueueVmmqEnabled bool,
	defaultQueueVmmqQueuePairs int32,
	defaultQueueVrssEnabled bool,
) (err error) {

	vmSwitchJson, err := json.Marshal(vmSwitch{
		Name:              name,
		Notes:             notes,
		AllowManagementOS: allowManagementOS,
		//EmbeddedTeamingEnabled:embeddedTeamingEnabled,
		//IovEnabled:iovEnabled,
		//PacketDirectEnabled:packetDirectEnabled,
		//BandwidthReservationMode:bandwidthReservationMode,
		SwitchType:                          switchType,
		NetAdapterNames:                     netAdapterNames,
		DefaultFlowMinimumBandwidthAbsolute: defaultFlowMinimumBandwidthAbsolute,
		DefaultFlowMinimumBandwidthWeight:   defaultFlowMinimumBandwidthWeight,
		DefaultQueueVmmqEnabled:             defaultQueueVmmqEnabled,
		DefaultQueueVmmqQueuePairs:          defaultQueueVmmqQueuePairs,
		DefaultQueueVrssEnabled:             defaultQueueVrssEnabled,
	})

	if err != nil {
		return err
	}

	err = c.runFireAndForgetScript(updateVMSwitchTemplate, updateVMSwitchArgs{
		VmSwitchJson: string(vmSwitchJson),
	})

	return err
}

type deleteVMSwitchArgs struct {
	Name string
}

var deleteVMSwitchTemplate = template.Must(template.New("DeleteVMSwitch").Parse(`
$ErrorActionPreference = 'Stop'
Get-VMSwitch -Name '{{.Name}}*' | ?{$_.Name -eq '{{.Name}}'} | Remove-VMSwitch -Force
`))

func (c *HypervClient) DeleteVMSwitch(name string) (err error) {
	err = c.runFireAndForgetScript(deleteVMSwitchTemplate, deleteVMSwitchArgs{
		Name: name,
	})

	return err
}
