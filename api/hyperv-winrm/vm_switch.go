package hyperv_winrm

import (
	"context"
	"encoding/json"
	"text/template"

	"github.com/taliesins/terraform-provider-hyperv/api"
)

type existsVMSwitchArgs struct {
	Name string
}

var existsVMSwitchTemplate = template.Must(template.New("ExistsVMSwitch").Parse(`
$ErrorActionPreference = 'Stop'
$vmSwitchObject = Get-VMSwitch -Name '{{.Name}}*' | ?{$_.Name -eq '{{.Name}}' }

if ($vmSwitchObject){
	$exists = ConvertTo-Json -InputObject @{Exists=$true}
	$exists
} else {
	$exists = ConvertTo-Json -InputObject @{Exists=$false}
	$exists
}
`))

func (c *ClientConfig) VMSwitchExists(ctx context.Context, name string) (result api.VmSwitchExists, err error) {
	err = c.WinRmClient.RunScriptWithResult(ctx, existsVMSwitchTemplate, existsVMSwitchArgs{
		Name: name,
	}, &result)

	return result, err
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

if ($vmSwitch.OperationMode -eq 1) {
	$SetVmNetworkAdapterVlanArgs = @{}
	$SetVmNetworkAdapterVlanArgs.VMNetworkAdapterName=$vmSwitch.Name
	$SetVmNetworkAdapterVlanArgs.ManagementOS=$true
	$SetVmNetworkAdapterVlanArgs.Access=$true
	$SetVmNetworkAdapterVlanArgs.VlanID=$vmSwitch.VlanId

	Set-VMNetworkAdapterVlan @SetVmNetworkAdapterVlanArgs
}

`))

func (c *ClientConfig) CreateVMSwitch(
	ctx context.Context,
	name string,
	notes string,
	allowManagementOS bool,
	embeddedTeamingEnabled bool,
	iovEnabled bool,
	packetDirectEnabled bool,
	bandwidthReservationMode api.VMSwitchBandwidthMode,
	switchType api.VMSwitchType,
	netAdapterNames []string,
	defaultFlowMinimumBandwidthAbsolute int64,
	defaultFlowMinimumBandwidthWeight int64,
	defaultQueueVmmqEnabled bool,
	defaultQueueVmmqQueuePairs int32,
	defaultQueueVrssEnabled bool,
	operationMode api.VMSwitchOperationMode,
	vlanId int,
) (err error) {
	vmSwitchJson, err := json.Marshal(api.VmSwitch{
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
		OperationMode:                       operationMode,
		VlanID:                              vlanId,
	})

	if err != nil {
		return err
	}

	err = c.WinRmClient.RunFireAndForgetScript(ctx, createVMSwitchTemplate, createVMSwitchArgs{
		VmSwitchJson: string(vmSwitchJson),
	})

	return err
}

type getVMSwitchArgs struct {
	Name string
}

var getVMSwitchTemplate = template.Must(template.New("GetVMSwitch").Parse(`
$ErrorActionPreference = 'Stop'
$vmAdapterVlanObject = Get-VMNetworkAdapterVlan -ManagementOS -VMNetworkAdapterName '{{.Name}}*' | %{ @{
	OperationMode=$_.OperationMode
	AccessVlanId=$_.AccessVlanId
}}

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
	OperationMode=$vmAdapterVlanObject.OperationMode
	VlanID=$vmAdapterVlanObject.AccessVlanId
}}


if ($vmSwitchObject){
	$vmSwitch = ConvertTo-Json -InputObject $vmSwitchObject
	$vmSwitch
} else {
	"{}"
}
`))

func (c *ClientConfig) GetVMSwitch(ctx context.Context, name string) (result api.VmSwitch, err error) {
	err = c.WinRmClient.RunScriptWithResult(ctx, getVMSwitchTemplate, getVMSwitchArgs{
		Name: name,
	}, &result)

	return result, err
}

type updateVMSwitchArgs struct {
	OldName      string
	VmSwitchJson string
}

var updateVMSwitchTemplate = template.Must(template.New("UpdateVMSwitch").Parse(`
$ErrorActionPreference = 'Stop'
Import-Module Hyper-V
$oldName = '{{.OldName}}'
$vmSwitch = '{{.VmSwitchJson}}' | ConvertFrom-Json
$minimumBandwidthMode = [Microsoft.HyperV.PowerShell.VMSwitchBandwidthMode]$vmSwitch.BandwidthReservationMode
$switchType = [Microsoft.HyperV.PowerShell.VMSwitchType]$vmSwitch.SwitchType
$NetAdapterNames = @($vmSwitch.NetAdapterNames)

#when EnablePacketDirect=true it seems to throw an exception if EnableIov=true or EnableEmbeddedTeaming=true

$switchObject = Get-VMSwitch -Name "$($oldName)*" | ?{$_.Name -eq $oldName}

if (!$switchObject){
	throw "Switch does not exist - $($oldName)"
}

if ($oldName -ne $vmSwitch.Name) {
	Rename-VMSwitch -Name $oldName -NewName $vmSwitch.Name
}

$SetVmSwitchArgs = @{}
$SetVmSwitchArgs.Name=$vmSwitch.Name
$SetVmSwitchArgs.Notes=$vmSwitch.Notes
if ($NetAdapterNames) {
	$SetVmSwitchArgs.AllowManagementOS=$vmSwitch.AllowManagementOS
	# Converts the incoming Object[] to a String as expected by the command
	$SetVmSwitchArgs.NetAdapterName=[system.String]::Join(",", $NetAdapterNames)
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

if ($vmSwitch.OperationMode -eq 1) {
	$SetVmNetworkAdapterVlanArgs = @{}
	$SetVmNetworkAdapterVlanArgs.VMNetworkAdapterName=$vmSwitch.Name
	$SetVmNetworkAdapterVlanArgs.ManagementOS=$true
	$SetVmNetworkAdapterVlanArgs.Access=$true
	$SetVmNetworkAdapterVlanArgs.VlanID=$vmSwitch.VlanId

	Set-VMNetworkAdapterVlan @SetVmNetworkAdapterVlanArgs
}
`))

func (c *ClientConfig) UpdateVMSwitch(
	ctx context.Context,
	oldName string,
	name string,
	notes string,
	allowManagementOS bool,
	// embeddedTeamingEnabled bool,
	// iovEnabled bool,
	// packetDirectEnabled bool,
	// bandwidthReservationMode api.VMSwitchBandwidthMode,
	switchType api.VMSwitchType,
	netAdapterNames []string,
	defaultFlowMinimumBandwidthAbsolute int64,
	defaultFlowMinimumBandwidthWeight int64,
	defaultQueueVmmqEnabled bool,
	defaultQueueVmmqQueuePairs int32,
	defaultQueueVrssEnabled bool,
	operationMode api.VMSwitchOperationMode,
	vlanId int,
) (err error) {
	vmSwitchJson, err := json.Marshal(api.VmSwitch{
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
		OperationMode:                       operationMode,
		VlanID:                              vlanId,
	})

	if err != nil {
		return err
	}

	err = c.WinRmClient.RunFireAndForgetScript(ctx, updateVMSwitchTemplate, updateVMSwitchArgs{
		OldName:      oldName,
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

func (c *ClientConfig) DeleteVMSwitch(ctx context.Context, name string) (err error) {
	err = c.WinRmClient.RunFireAndForgetScript(ctx, deleteVMSwitchTemplate, deleteVMSwitchArgs{
		Name: name,
	})

	return err
}
