package hyperv

import (
	"text/template"
	"encoding/json"
	"strings"
)

type VMSwitchBandwidthMode int

const (
	VMSwitchBandwidthMode_Default VMSwitchBandwidthMode = 0
	VMSwitchBandwidthMode_Weight VMSwitchBandwidthMode = 1
	VMSwitchBandwidthMode_Absolute VMSwitchBandwidthMode = 2
	VMSwitchBandwidthMode_None VMSwitchBandwidthMode = 3
)

var VMSwitchBandwidthMode_name = map[VMSwitchBandwidthMode]string{
	VMSwitchBandwidthMode_Default: "Default",
	VMSwitchBandwidthMode_Weight: "Weight",
	VMSwitchBandwidthMode_Absolute: "Absolute",
	VMSwitchBandwidthMode_None: "None",
}

var VMSwitchBandwidthMode_value = map[string]VMSwitchBandwidthMode{
	"default": VMSwitchBandwidthMode_Default,
	"weight": VMSwitchBandwidthMode_Weight,
	"absolute": VMSwitchBandwidthMode_Absolute,
	"none": VMSwitchBandwidthMode_None,
}

func (x VMSwitchBandwidthMode) String() string {
	return VMSwitchBandwidthMode_name[x]
}

func ToVMSwitchBandwidthMode(x string) VMSwitchBandwidthMode {
	return VMSwitchBandwidthMode_value[strings.ToLower(x)]
}

type VMSwitchType int

const (
	VMSwitchType_Private VMSwitchType = 0
	VMSwitchType_Internal VMSwitchType = 1
	VMSwitchType_External VMSwitchType = 2
)

var VMSwitchType_name = map[VMSwitchType]string{
	VMSwitchType_Private: "Private",
	VMSwitchType_Internal: "Internal",
	VMSwitchType_External: "External",
}

var VMSwitchType_value = map[string]VMSwitchType{
	"private": VMSwitchType_Private,
	"internal": VMSwitchType_Internal,
	"external": VMSwitchType_External,
}

func (x VMSwitchType) String() string {
	return VMSwitchType_name[x]
}

func ToVMSwitchType(x string) VMSwitchType {
	return VMSwitchType_value[strings.ToLower(x)]
}

type vmSwitch struct {
	Name		string
	Notes		string
	AllowManagementOS bool
	EmbeddedTeamingEnabled bool
	IovEnabled bool
	PacketDirectEnabled bool
	BandwidthReservationMode VMSwitchBandwidthMode
	SwitchType	VMSwitchType
	NetAdapterInterfaceDescriptions []string
	NetAdapterNames []string
	DefaultFlowMinimumBandwidthAbsolute int64
	DefaultFlowMinimumBandwidthWeight int64
	DefaultQueueVmmqEnabled bool
	DefaultQueueVmmqQueuePairs int32
	DefaultQueueVrssEnabled bool
}

type createVMSwitchArgs struct {
	VmSwitchJson		string
}

var createVMSwitchTemplate = template.Must(template.New("CreateVMSwitch").Parse(`
$ErrorActionPreference = 'Stop'
Get-Vm | Out-Null
$vmSwitch = '{{.VmSwitchJson}}' | ConvertFrom-Json
$minimumBandwidthMode = [Microsoft.HyperV.PowerShell.VMSwitchBandwidthMode]$vmSwitch.BandwidthReservationMode
$switchType = [Microsoft.HyperV.PowerShell.VMSwitchType]$vmSwitch.SwitchType
$NetAdapterInterfaceDescriptions = @($vmSwitch.NetAdapterInterfaceDescriptions)
$NetAdapterNames = @($vmSwitch.$NetAdapterNames)
#when EnablePacketDirect=true it seems to throw an exception if EnableIov=true or EnableEmbeddedTeaming=true

$switchObject = Get-VMSwitch | ?{$_.Name -eq $vmSwitch.Name}

if ($switchObject){
	throw "Switch already exists - $($vmSwitch.Name)"
}

if ($NetAdapterInterfaceDescriptions -or $NetAdapterNames) {
	New-VMSwitch -Name $vmSwitch.Name -Notes $vmSwitch.Notes -AllowManagementOS $vmSwitch.AllowManagementOS -EnableEmbeddedTeaming $vmSwitch.EmbeddedTeamingEnabled -EnableIov $vmSwitch.IovEnabled -EnablePacketDirect $vmSwitch.PacketDirectEnabled -MinimumBandwidthMode $minimumBandwidthMode -NetAdapterInterfaceDescription $NetAdapterInterfaceDescriptions -NetAdapterName $NetAdapterNames
} else {
	New-VMSwitch -Name $vmSwitch.Name -Notes $vmSwitch.Notes -EnableEmbeddedTeaming $vmSwitch.EmbeddedTeamingEnabled -EnableIov $vmSwitch.IovEnabled -EnablePacketDirect $vmSwitch.PacketDirectEnabled -MinimumBandwidthMode $minimumBandwidthMode -SwitchType $switchType

	#not used unless interface is specified
	#-AllowManagementOS $vmSwitch.AllowManagementOS
}

$switchObject = Get-VMSwitch -Name $vmSwitch.Name

if ($switchObject.DefaultFlowMinimumBandwidthAbsolute -ne $vmSwitch.DefaultFlowMinimumBandwidthAbsolute) {
	Set-VMSwitch -Name $vmSwitch.Name -DefaultFlowMinimumBandwidthAbsolute $vmSwitch.DefaultFlowMinimumBandwidthAbsolute
}

if ($switchObject.DefaultFlowMinimumBandwidthWeight -ne $vmSwitch.DefaultFlowMinimumBandwidthWeight) {
	Set-VMSwitch -Name $vmSwitch.Name -DefaultFlowMinimumBandwidthWeight $vmSwitch.DefaultFlowMinimumBandwidthWeight
}

if ($switchObject.DefaultQueueVmmqEnabled -ne $vmSwitch.DefaultQueueVmmqEnabled) {
	Set-VMSwitch -Name $vmSwitch.Name -DefaultQueueVmmqEnabled $vmSwitch.DefaultQueueVmmqEnabled
}

if ($switchObject.DefaultQueueVmmqQueuePairs -ne $vmSwitch.DefaultQueueVmmqQueuePairs) {
	Set-VMSwitch -Name $vmSwitch.Name -DefaultQueueVmmqQueuePairs $vmSwitch.DefaultQueueVmmqQueuePairs
}

if ($switchObject.DefaultQueueVrssEnabled -ne $vmSwitch.DefaultQueueVrssEnabled) {
	Set-VMSwitch -Name $vmSwitch.Name -DefaultQueueVrssEnabled $vmSwitch.DefaultQueueVrssEnabled
}
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
	netAdapterInterfaceDescriptions []string,
	netAdapterNames []string,
	defaultFlowMinimumBandwidthAbsolute int64,
	defaultFlowMinimumBandwidthWeight int64,
	defaultQueueVmmqEnabled bool,
	defaultQueueVmmqQueuePairs int32,
	defaultQueueVrssEnabled bool,
) (err error) {

	vmSwitchJson, err := json.Marshal(vmSwitch{
		Name:name,
		Notes:notes,
		AllowManagementOS:allowManagementOS,
		EmbeddedTeamingEnabled:embeddedTeamingEnabled,
		IovEnabled:iovEnabled,
		PacketDirectEnabled:packetDirectEnabled,
		BandwidthReservationMode:bandwidthReservationMode,
		SwitchType:switchType,
		NetAdapterInterfaceDescriptions:netAdapterInterfaceDescriptions,
		NetAdapterNames:netAdapterNames,
		DefaultFlowMinimumBandwidthAbsolute:defaultFlowMinimumBandwidthAbsolute,
		DefaultFlowMinimumBandwidthWeight:defaultFlowMinimumBandwidthWeight,
		DefaultQueueVmmqEnabled:defaultQueueVmmqEnabled,
		DefaultQueueVmmqQueuePairs:defaultQueueVmmqQueuePairs,
		DefaultQueueVrssEnabled:defaultQueueVrssEnabled,
	})

	err = c.runFireAndForgetScript(createVMSwitchTemplate, createVMSwitchArgs{
		VmSwitchJson:string(vmSwitchJson),
	});

	return err
}

type getVMSwitchArgs struct {
	Name		string
}

var getVMSwitchTemplate = template.Must(template.New("GetVMSwitch").Parse(`
$ErrorActionPreference = 'Stop'
$vmSwitch = Get-VMSwitch | ?{$_.Name -eq '{{.Name}}' } | %{ @{
	Name=$_.Name;
	Notes=$_.Notes;
	AllowManagementOS=$_.AllowManagementOS;
	EmbeddedTeamingEnabled=$_.EmbeddedTeamingEnabled;
	IovEnabled=$_.IovEnabled;
	PacketDirectEnabled=$_.PacketDirectEnabled;
	BandwidthReservationMode=$_.BandwidthReservationMode;
	SwitchType=$_.SwitchType;
	NetAdapterInterfaceDescriptions=$_.NetAdapterInterfaceDescriptions;
	NetAdapterNames=if($_.NetAdapterInterfaceDescriptions){@(Get-NetAdapter -InterfaceDescription $_.NetAdapterInterfaceDescriptions | %{$_.Name})}else{@()};
	DefaultFlowMinimumBandwidthAbsolute=$_.DefaultFlowMinimumBandwidthAbsolute;
	DefaultFlowMinimumBandwidthWeight=$_.DefaultFlowMinimumBandwidthWeight;
	DefaultQueueVmmqEnabled=$_.DefaultQueueVmmqEnabled;
	DefaultQueueVmmqQueuePairs=$_.DefaultQueueVmmqQueuePairs;
	DefaultQueueVrssEnabled=$_.DefaultQueueVrssEnabled;
}} | ConvertTo-Json

if (!$vmSwitch) {
	$vmSwitch = '{"NetAdapterNames":[]}'
}

$vmSwitch
`))

func (c *HypervClient) GetVMSwitch(name string) (result vmSwitch, err error) {
	err = c.runScriptWithResult(getVMSwitchTemplate, getVMSwitchArgs{
		Name:name,
	}, &result);

	return result, err
}

type updateVMSwitchArgs struct {
	VmSwitchJson		string
}

var updateVMSwitchTemplate = template.Must(template.New("UpdateVMSwitch").Parse(`
$ErrorActionPreference = 'Stop'
Get-Vm | Out-Null
$vmSwitch = '{{.VmSwitchJson}}' | ConvertFrom-Json
$minimumBandwidthMode = [Microsoft.HyperV.PowerShell.VMSwitchBandwidthMode]$vmSwitch.BandwidthReservationMode
$switchType = [Microsoft.HyperV.PowerShell.VMSwitchType]$vmSwitch.SwitchType
$NetAdapterInterfaceDescriptions = @($vmSwitch.NetAdapterInterfaceDescriptions)
$NetAdapterNames = @($vmSwitch.$NetAdapterNames)
#when EnablePacketDirect=true it seems to throw an exception if EnableIov=true or EnableEmbeddedTeaming=true

$switchObject = Get-VMSwitch | ?{$_.Name -eq $vmSwitch.Name}

if (!$switchObject){
	throw "Switch does not exist - $($vmSwitch.Name)"
}

if ($NetAdapterInterfaceDescriptions -or $NetAdapterNames) {
	Set-VMSwitch -Name $vmSwitch.Name -AllowManagementOS $vmSwitch.AllowManagementOS -NetAdapterInterfaceDescription $vmSwitch.NetAdapterInterfaceDescriptions -NetAdapterName $NetAdapterNames

	#Updates not supported on:
	#-EnableEmbeddedTeaming $vmSwitch.EmbeddedTeamingEnabled
	#-EnableIov $vmSwitch.IovEnabled
	#-EnablePacketDirect $vmSwitch.PacketDirectEnabled
	#-MinimumBandwidthMode $minimumBandwidthMode
} else {
	Set-VMSwitch -Name $vmSwitch.Name -SwitchType $switchType

	#Updates not supported on:
	#-EnableEmbeddedTeaming $vmSwitch.EmbeddedTeamingEnabled
	#-EnableIov $vmSwitch.IovEnabled
	#-EnablePacketDirect $vmSwitch.PacketDirectEnabled
	#-MinimumBandwidthMode $minimumBandwidthMode

	#not used unless interface is specified
	#-AllowManagementOS $vmSwitch.AllowManagementOS
}

if ($switchObject.Notes -ne $vmSwitch.Notes) {
	Set-VMSwitch -Name $vmSwitch.Name -Notes $vmSwitch.Notes
}

if ($switchObject.DefaultFlowMinimumBandwidthAbsolute -ne $vmSwitch.DefaultFlowMinimumBandwidthAbsolute) {
	Set-VMSwitch -Name $vmSwitch.Name -DefaultFlowMinimumBandwidthAbsolute $vmSwitch.DefaultFlowMinimumBandwidthAbsolute
}

if ($switchObject.DefaultFlowMinimumBandwidthWeight -ne $vmSwitch.DefaultFlowMinimumBandwidthWeight) {
	Set-VMSwitch -Name $vmSwitch.Name -DefaultFlowMinimumBandwidthWeight $vmSwitch.DefaultFlowMinimumBandwidthWeight
}

if ($switchObject.DefaultQueueVmmqEnabled -ne $vmSwitch.DefaultQueueVmmqEnabled) {
	Set-VMSwitch -Name $vmSwitch.Name -DefaultQueueVmmqEnabled $vmSwitch.DefaultQueueVmmqEnabled
}

if ($switchObject.DefaultQueueVmmqQueuePairs -ne $vmSwitch.DefaultQueueVmmqQueuePairs) {
	Set-VMSwitch -Name $vmSwitch.Name -DefaultQueueVmmqQueuePairs $vmSwitch.DefaultQueueVmmqQueuePairs
}

if ($switchObject.DefaultQueueVrssEnabled -ne $vmSwitch.DefaultQueueVrssEnabled) {
	Set-VMSwitch -Name $vmSwitch.Name -DefaultQueueVrssEnabled $vmSwitch.DefaultQueueVrssEnabled
}
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
	netAdapterInterfaceDescriptions []string,
	netAdapterNames []string,
	defaultFlowMinimumBandwidthAbsolute int64,
	defaultFlowMinimumBandwidthWeight int64,
	defaultQueueVmmqEnabled bool,
	defaultQueueVmmqQueuePairs int32,
	defaultQueueVrssEnabled bool,
) (err error) {

	vmSwitchJson, err := json.Marshal(vmSwitch{
		Name:name,
		Notes:notes,
		AllowManagementOS:allowManagementOS,
		//EmbeddedTeamingEnabled:embeddedTeamingEnabled,
		//IovEnabled:iovEnabled,
		//PacketDirectEnabled:packetDirectEnabled,
		//BandwidthReservationMode:bandwidthReservationMode,
		SwitchType:switchType,
		NetAdapterInterfaceDescriptions:netAdapterInterfaceDescriptions,
		NetAdapterNames:netAdapterNames,
		DefaultFlowMinimumBandwidthAbsolute:defaultFlowMinimumBandwidthAbsolute,
		DefaultFlowMinimumBandwidthWeight:defaultFlowMinimumBandwidthWeight,
		DefaultQueueVmmqEnabled:defaultQueueVmmqEnabled,
		DefaultQueueVmmqQueuePairs:defaultQueueVmmqQueuePairs,
		DefaultQueueVrssEnabled:defaultQueueVrssEnabled,
	})

	err = c.runFireAndForgetScript(updateVMSwitchTemplate, updateVMSwitchArgs{
		VmSwitchJson:string(vmSwitchJson),
	});

	return err
}

type deleteVMSwitchArgs struct {
	Name		string
}

var deleteVMSwitchTemplate = template.Must(template.New("DeleteVMSwitch").Parse(`
$ErrorActionPreference = 'Stop'
Get-VMSwitch | ?{$_.Name -eq '{{.Name}}'} | Remove-VMSwitch
`))

func (c *HypervClient) DeleteVMSwitch(name string) (err error) {
	err = c.runFireAndForgetScript(deleteVMSwitchTemplate, deleteVMSwitchArgs{
		Name:name,
	});

	return err
}
