package hyperv_winrm

import (
	"context"
	"encoding/json"
	"github.com/taliesins/terraform-provider-hyperv/api"
	"text/template"
)

type createVmNetworkAdapterArgs struct {
	VmNetworkAdapterJson string
}

var createVmNetworkAdapterTemplate = template.Must(template.New("CreateVmNetworkAdapter").Parse(`
$ErrorActionPreference = 'Stop'
Import-Module Hyper-V
$vmNetworkAdapter = '{{.VmNetworkAdapterJson}}' | ConvertFrom-Json

$dhcpGuard = [Microsoft.HyperV.PowerShell.OnOffState]$vmNetworkAdapter.DhcpGuard
$routerGuard = [Microsoft.HyperV.PowerShell.OnOffState]$vmNetworkAdapter.RouterGuard
$portMirroring = [Microsoft.HyperV.PowerShell.VMNetworkAdapterPortMirroringMode]$vmNetworkAdapter.PortMirroring
$ieeePriorityTag = [Microsoft.HyperV.PowerShell.OnOffState]$vmNetworkAdapter.IeeePriorityTag
$iovInterruptModeration = [Microsoft.HyperV.PowerShell.IovInterruptModerationValue]$vmNetworkAdapter.IovInterruptModeration
$allowTeaming = [Microsoft.HyperV.PowerShell.OnOffState]$vmNetworkAdapter.AllowTeaming
$deviceNaming = [Microsoft.HyperV.PowerShell.OnOffState]$vmNetworkAdapter.DeviceNaming
$fixSpeed10G = [Microsoft.HyperV.PowerShell.OnOffState]$vmNetworkAdapter.FixSpeed10G

$NewVmNetworkAdapterArgs = @{
	VmName=$vmNetworkAdapter.VmName
	Name=$vmNetworkAdapter.Name
	IsLegacy=$vmNetworkAdapter.IsLegacy
	SwitchName=$vmNetworkAdapter.SwitchName
}

Add-VmNetworkAdapter @NewVmNetworkAdapterArgs

$minimumBandwidthMode = [Microsoft.HyperV.PowerShell.VMSwitchBandwidthMode]::None

if ($vmNetworkAdapter.SwitchName) {
	$vmSwitch = Get-VMSwitch -Name $vmNetworkAdapter.SwitchName
	if ($vmSwitch) {
		$minimumBandwidthMode = $vmSwitch.BandwidthReservationMode
	}
}

$SetVmNetworkAdapterArgs = @{}
$SetVmNetworkAdapterArgs.VmName=$vmNetworkAdapter.VmName
$SetVmNetworkAdapterArgs.Name=$vmNetworkAdapter.Name
if ($vmNetworkAdapter.DynamicMacAddress) {
	$SetVmNetworkAdapterArgs.DynamicMacAddress=$vmNetworkAdapter.DynamicMacAddress
} elseif ($vmNetworkAdapter.StaticMacAddress) {
	$SetVmNetworkAdapterArgs.StaticMacAddress=$vmNetworkAdapter.StaticMacAddress
}
if ($vmNetworkAdapter.MacAddressSpoofing) {
	$SetVmNetworkAdapterArgs.MacAddressSpoofing=$vmNetworkAdapter.MacAddressSpoofing
}
$SetVmNetworkAdapterArgs.DhcpGuard=$dhcpGuard
$SetVmNetworkAdapterArgs.RouterGuard=$routerGuard
$SetVmNetworkAdapterArgs.PortMirroring=$portMirroring
$SetVmNetworkAdapterArgs.IeeePriorityTag=$ieeePriorityTag
$SetVmNetworkAdapterArgs.VmqWeight=$vmNetworkAdapter.VmqWeight
$SetVmNetworkAdapterArgs.IovQueuePairsRequested=$vmNetworkAdapter.IovQueuePairsRequested
$SetVmNetworkAdapterArgs.IovInterruptModeration=$iovInterruptModeration
$SetVmNetworkAdapterArgs.IovWeight=$vmNetworkAdapter.IovWeight
$SetVmNetworkAdapterArgs.IPsecOffloadMaximumSecurityAssociation=$vmNetworkAdapter.IPsecOffloadMaximumSecurityAssociation
$SetVmNetworkAdapterArgs.MaximumBandwidth=$vmNetworkAdapter.MaximumBandwidth
if ($minimumBandwidthMode -eq [Microsoft.HyperV.PowerShell.VMSwitchBandwidthMode]::Absolute){
	$SetVmNetworkAdapterArgs.MinimumBandwidthAbsolute=$vmNetworkAdapter.MinimumBandwidthAbsolute
}
if ($minimumBandwidthMode -eq [Microsoft.HyperV.PowerShell.VMSwitchBandwidthMode]::Weight -or $minimumBandwidthMode -eq [Microsoft.HyperV.PowerShell.VMSwitchBandwidthMode]::Default){
	$SetVmNetworkAdapterArgs.MinimumBandwidthWeight=$vmNetworkAdapter.MinimumBandwidthWeight
}
$SetVmNetworkAdapterArgs.MandatoryFeatureId=$vmNetworkAdapter.MandatoryFeatureId
if ($vmNetworkAdapter.ResourcePoolName) {
	$SetVmNetworkAdapterArgs.ResourcePoolName=$vmNetworkAdapter.ResourcePoolName
}
$SetVmNetworkAdapterArgs.TestReplicaPoolName=$vmNetworkAdapter.TestReplicaPoolName
$SetVmNetworkAdapterArgs.TestReplicaSwitchName=$vmNetworkAdapter.TestReplicaSwitchName
$SetVmNetworkAdapterArgs.VirtualSubnetId=$vmNetworkAdapter.VirtualSubnetId
$SetVmNetworkAdapterArgs.AllowTeaming=$allowTeaming
$SetVmNetworkAdapterArgs.NotMonitoredInCluster=$vmNetworkAdapter.NotMonitoredInCluster
$SetVmNetworkAdapterArgs.StormLimit=$vmNetworkAdapter.StormLimit
$SetVmNetworkAdapterArgs.DynamicIPAddressLimit=$vmNetworkAdapter.DynamicIPAddressLimit
$SetVmNetworkAdapterArgs.DeviceNaming=$deviceNaming
$SetVmNetworkAdapterArgs.FixSpeed10G=$fixSpeed10G
$SetVmNetworkAdapterArgs.PacketDirectNumProcs=$vmNetworkAdapter.PacketDirectNumProcs
$SetVmNetworkAdapterArgs.PacketDirectModerationCount=$vmNetworkAdapter.PacketDirectModerationCount
$SetVmNetworkAdapterArgs.PacketDirectModerationInterval=$vmNetworkAdapter.PacketDirectModerationInterval
$SetVmNetworkAdapterArgs.VrssEnabled=$vmNetworkAdapter.VrssEnabled
$SetVmNetworkAdapterArgs.VmmqEnabled=$vmNetworkAdapter.VmmqEnabled
$SetVmNetworkAdapterArgs.VmmqQueuePairs=$vmNetworkAdapter.VmmqQueuePairs

Set-VmNetworkAdapter @SetVmNetworkAdapterArgs

if ($vmNetworkAdapter.VlanAccess -and $vmNetworkAdapter.VlanId) {
	$SetVmNetworkAdapterVlanArgs = @{}

	$SetVmNetworkAdapterVlanArgs.VMName  = $vmNetworkAdapter.VmName
	$SetVmNetworkAdapterVlanArgs.VMNetworkAdapterName   = $vmNetworkAdapter.Name
	$SetVmNetworkAdapterVlanArgs.Access = $true
	$SetVmNetworkAdapterVlanArgs.VlanId = $vmNetworkAdapter.VlanId

	Set-VmNetworkAdapterVlan @SetVmNetworkAdapterVlanArgs
}

`))

func (c *ClientConfig) CreateVmNetworkAdapter(
	ctx context.Context,
	vmName string,
	name string,
	switchName string,
	managementOs bool,
	isLegacy bool,
	dynamicMacAddress bool,
	staticMacAddress string,
	macAddressSpoofing api.OnOffState,
	dhcpGuard api.OnOffState,
	routerGuard api.OnOffState,
	portMirroring api.PortMirroring,
	ieeePriorityTag api.OnOffState,
	vmqWeight int,
	iovQueuePairsRequested int,
	iovInterruptModeration api.IovInterruptModerationValue,
	iovWeight int,
	ipsecOffloadMaximumSecurityAssociation int,
	maximumBandwidth int,
	minimumBandwidthAbsolute int,
	minimumBandwidthWeight int,
	mandatoryFeatureId []string,
	resourcePoolName string,
	testReplicaPoolName string,
	testReplicaSwitchName string,
	virtualSubnetId int,
	allowTeaming api.OnOffState,
	notMonitoredInCluster bool,
	stormLimit int,
	dynamicIpAddressLimit int,
	deviceNaming api.OnOffState,
	fixSpeed10G api.OnOffState,
	packetDirectNumProcs int,
	packetDirectModerationCount int,
	packetDirectModerationInterval int,
	vrssEnabled bool,
	vmmqEnabled bool,
	vmmqQueuePairs int,
	vlanAccess bool,
	vlanId int,
) (err error) {

	vmNetworkAdapterJson, err := json.Marshal(api.VmNetworkAdapter{
		VmName:                                 vmName,
		Name:                                   name,
		SwitchName:                             switchName,
		ManagementOs:                           managementOs,
		IsLegacy:                               isLegacy,
		DynamicMacAddress:                      dynamicMacAddress,
		StaticMacAddress:                       staticMacAddress,
		MacAddressSpoofing:                     macAddressSpoofing,
		DhcpGuard:                              dhcpGuard,
		RouterGuard:                            routerGuard,
		PortMirroring:                          portMirroring,
		IeeePriorityTag:                        ieeePriorityTag,
		VmqWeight:                              vmqWeight,
		IovQueuePairsRequested:                 iovQueuePairsRequested,
		IovInterruptModeration:                 iovInterruptModeration,
		IovWeight:                              iovWeight,
		IpsecOffloadMaximumSecurityAssociation: ipsecOffloadMaximumSecurityAssociation,
		MaximumBandwidth:                       maximumBandwidth,
		MinimumBandwidthAbsolute:               minimumBandwidthAbsolute,
		MinimumBandwidthWeight:                 minimumBandwidthWeight,
		MandatoryFeatureId:                     mandatoryFeatureId,
		ResourcePoolName:                       resourcePoolName,
		TestReplicaPoolName:                    testReplicaPoolName,
		TestReplicaSwitchName:                  testReplicaSwitchName,
		VirtualSubnetId:                        virtualSubnetId,
		AllowTeaming:                           allowTeaming,
		NotMonitoredInCluster:                  notMonitoredInCluster,
		StormLimit:                             stormLimit,
		DynamicIpAddressLimit:                  dynamicIpAddressLimit,
		DeviceNaming:                           deviceNaming,
		FixSpeed10G:                            fixSpeed10G,
		PacketDirectNumProcs:                   packetDirectNumProcs,
		PacketDirectModerationCount:            packetDirectModerationCount,
		PacketDirectModerationInterval:         packetDirectModerationInterval,
		VrssEnabled:                            vrssEnabled,
		VmmqEnabled:                            vmmqEnabled,
		VmmqQueuePairs:                         vmmqQueuePairs,
		VlanAccess:                             vlanAccess,
		VlanId:                                 vlanId,
	})

	if err != nil {
		return err
	}

	err = c.WinRmClient.RunFireAndForgetScript(ctx, createVmNetworkAdapterTemplate, createVmNetworkAdapterArgs{
		VmNetworkAdapterJson: string(vmNetworkAdapterJson),
	})

	return err
}

type getVmNetworkAdaptersArgs struct {
	VmName string
}

var getVmNetworkAdaptersTemplate = template.Must(template.New("GetVmNetworkAdapters").Parse(`
$ErrorActionPreference = 'Stop'
#First 3 requests fails to get ip address
Get-VMNetworkAdapter -VmName '{{.VmName}}' | Out-Null
Get-VMNetworkAdapter -VmName '{{.VmName}}' | Out-Null
Get-VMNetworkAdapter -VmName '{{.VmName}}' | Out-Null

$vmNetworkAdaptersObject = @(Get-VM -Name '{{.VmName}}*' | ?{$_.Name -eq '{{.VmName}}' } | Get-VMNetworkAdapter | %{ @{
     Name=$_.Name;
     SwitchName=$_.SwitchName;
     ManagementOs=$_.IsManagementOs;
     IsLegacy=$_.IsLegacy;
     DynamicMacAddress=$_.DynamicMacAddressEnabled;
     StaticMacAddress=if ($_.MacAddress -eq '000000000000') { '' } else { $_.MacAddress };
     MacAddressSpoofing=$_.MacAddressSpoofing;
     DhcpGuard=$_.DhcpGuard;
     RouterGuard=$_.RouterGuard;
     PortMirroring=$_.PortMirroringMode;
     IeeePriorityTag=$_.IeeePriorityTag;
     VmqWeight=$_.VmqWeight;
     IovQueuePairsRequested=$_.IovQueuePairsRequested;
     IovInterruptModeration=$_.IovInterruptModeration;
     IovWeight=$_.IovWeight;
     IpsecOffloadMaximumSecurityAssociation=$_.IPsecOffloadMaxSA;
     MaximumBandwidth=$_.BandwidthSetting.MaximumBandwidth;
     MinimumBandwidthAbsolute=$_.BandwidthSetting.MinimumBandwidthAbsolute;
     MinimumBandwidthWeight=$_.BandwidthSetting.MinimumBandwidthWeight;
     MandatoryFeatureId=$_.MandatoryFeatureId;
     ResourcePoolName=$_.PoolName;
     TestReplicaPoolName=$_.TestReplicaPoolName;
     TestReplicaSwitchName=$_.TestReplicaSwitchName;
     VirtualSubnetId=$_.VirtualSubnetId;
     AllowTeaming=$_.AllowTeaming;
     NotMonitoredInCluster=!$_.ClusterMonitored;
     StormLimit=$_.StormLimit;
     DynamicIpAddressLimit=$_.DynamicIpAddressLimit;
     DeviceNaming=$_.DeviceNaming;
     FixSpeed10G=$_.FixSpeed10G;
     PacketDirectNumProcs=$_.PacketDirectNumProcs;
     PacketDirectModerationCount=$_.PacketDirectModerationCount;
     PacketDirectModerationInterval=$_.PacketDirectModerationInterval;
     VrssEnabled=$_.VrssEnabledRequested;
     VmmqEnabled=$_.VmmqEnabledRequested;
     VmmqQueuePairs=$_.VmmqQueuePairsRequested;
	 IpAddresses=@($_.IpAddresses);
	 VlanAccess=if ($_.VLanSetting.OperationMode -eq 'Access') {$true} else {$false};
	 VlanId=$_.VLanSetting.AccessVlanId;
}})

if ($vmNetworkAdaptersObject) {
	$vmNetworkAdapters = ConvertTo-Json -InputObject $vmNetworkAdaptersObject
	$vmNetworkAdapters
} else {
	"[]"
}
`))

func (c *ClientConfig) GetVmNetworkAdapters(ctx context.Context, vmName string, networkAdaptersWaitForIps []api.VmNetworkAdapterWaitForIp) (result []api.VmNetworkAdapter, err error) {
	result = make([]api.VmNetworkAdapter, 0)

	err = c.WinRmClient.RunScriptWithResult(ctx, getVmNetworkAdaptersTemplate, getVmNetworkAdaptersArgs{
		VmName: vmName,
	}, &result)

	//Enrich network adapter with config settings that are not stored in hyperv
	for _, networkAdapterWaitForIps := range networkAdaptersWaitForIps {
		for networkAdapterIndex, networkAdapter := range result {
			if networkAdapterWaitForIps.Name == networkAdapter.Name {
				result[networkAdapterIndex].WaitForIps = networkAdapterWaitForIps.WaitForIps
				break
			}
		}
	}

	return result, err
}

type waitForVmNetworkAdaptersIpsArgs struct {
	VmName                          string
	Timeout                         uint32
	PollPeriod                      uint32
	VmNetworkAdaptersWaitForIpsJson string
}

var waitForVmNetworkAdaptersIpsTemplate = template.Must(template.New("WaitForVmNetworkAdaptersIps").Parse(`
$ErrorActionPreference = 'Stop'

function Test-CanGetIpsForState($State){
	$states = @([Microsoft.HyperV.PowerShell.VMState]::Running,
			[Microsoft.HyperV.PowerShell.VMState]::RunningCritical
        )
    return $states -contains $state 
}

function Test-CanNotGetIpsForState($State){
    $states = @([Microsoft.HyperV.PowerShell.VMState]::Stopping,
			[Microsoft.HyperV.PowerShell.VMState]::StoppingCritical,
			[Microsoft.HyperV.PowerShell.VMState]::ForceShutdown,
			[Microsoft.HyperV.PowerShell.VMState]::Off,
			[Microsoft.HyperV.PowerShell.VMState]::OffCritical,
			[Microsoft.HyperV.PowerShell.VMState]::Paused,
			[Microsoft.HyperV.PowerShell.VMState]::PausedCritical
        )
    return $states -contains $state 
}

function Test-IsNotInFinalTransitionState($State){
    $states = @([Microsoft.HyperV.PowerShell.VMState]::Other,
		[Microsoft.HyperV.PowerShell.VMState]::Stopping,
		[Microsoft.HyperV.PowerShell.VMState]::Saved,
		[Microsoft.HyperV.PowerShell.VMState]::Starting,
		[Microsoft.HyperV.PowerShell.VMState]::Reset,
		[Microsoft.HyperV.PowerShell.VMState]::Saving,
		[Microsoft.HyperV.PowerShell.VMState]::Pausing,
		[Microsoft.HyperV.PowerShell.VMState]::Resuming,
		[Microsoft.HyperV.PowerShell.VMState]::FastSaved,
		[Microsoft.HyperV.PowerShell.VMState]::FastSaving,
		[Microsoft.HyperV.PowerShell.VMState]::ForceShutdown,
		[Microsoft.HyperV.PowerShell.VMState]::ForceReboot,
        [Microsoft.HyperV.PowerShell.VMState]::StoppingCritical,
        [Microsoft.HyperV.PowerShell.VMState]::SavedCritical,
        [Microsoft.HyperV.PowerShell.VMState]::StartingCritical,
        [Microsoft.HyperV.PowerShell.VMState]::ResetCritical,
        [Microsoft.HyperV.PowerShell.VMState]::SavingCritical,
        [Microsoft.HyperV.PowerShell.VMState]::PausingCritical,
        [Microsoft.HyperV.PowerShell.VMState]::ResumingCritical,
        [Microsoft.HyperV.PowerShell.VMState]::FastSavedCritical,
        [Microsoft.HyperV.PowerShell.VMState]::FastSavingCritical
        )
	   
    return $states -contains $State 
}

function Wait-ForNetworkAdapterIps($Name, $Timeout, $PollPeriod, $VmNetworkAdaptersToWaitForIps){
	$timer = [Diagnostics.Stopwatch]::StartNew()
	while ($timer.Elapsed.TotalSeconds -lt $Timeout) {
        $vmObject = Get-VM -Name "$($vmName)*" | ?{$_.Name -eq $vmName}

        if (!(Test-IsNotInFinalTransitionState $vmObject.state)){
            if (Test-CanGetIpsForState $vmObject.state) {
                $waitForIp = $false

                $VmNetworkAdaptersToWaitForIps | ?{$_.WaitForIps} | %{
                    $name = $_.Name
                    $ipAddresses = @($vmObject.NetworkAdapters | ?{$_.Name -eq $name} | %{$_.IPAddresses} |?{$_})

                    if ((!($ipAddresses)) -or ($ipAddresses -contains '0.0.0.0')){
                        $waitForIp = $true
                    } 
                }

                if (!$waitForIp){
                    break
                }
           	} elseif (Test-CanNotGetIpsForState $vmObject.state) {
               	break
           	}
       	}

        Start-Sleep -Seconds $PollPeriod
	}
	$timer.Stop()

	if ($timer.Elapsed.TotalSeconds -gt $Timeout) {
		throw 'Timeout while waiting for vm $($Name) to read network adapter ips'
	} 
}

Import-Module Hyper-V
$vmNetworkAdaptersToWaitForIps = '{{.VmNetworkAdaptersWaitForIpsJson}}' | ConvertFrom-Json
$vmName = '{{.VmName}}'
$vmObject = Get-VM -Name "$($vmName)*" | ?{$_.Name -eq $vmName}
$timeout = {{.Timeout}}
$pollPeriod = {{.PollPeriod}}

if (!$vmObject){
	throw "VM does not exist - $($vmName)"
}

Wait-ForNetworkAdapterIps -Name $vmName -Timeout $timeout -PollPeriod $pollPeriod -VmNetworkAdaptersToWaitForIps $vmNetworkAdaptersToWaitForIps

`))

func (c *ClientConfig) WaitForVmNetworkAdaptersIps(
	ctx context.Context,
	vmName string,
	timeout uint32,
	pollPeriod uint32,
	vmNetworkAdaptersWaitForIps []api.VmNetworkAdapterWaitForIp,
) (err error) {

	vmNetworkAdaptersWaitForIpsJson, err := json.Marshal(vmNetworkAdaptersWaitForIps)

	if err != nil {
		return err
	}

	err = c.WinRmClient.RunFireAndForgetScript(ctx, waitForVmNetworkAdaptersIpsTemplate, waitForVmNetworkAdaptersIpsArgs{
		VmName:                          vmName,
		Timeout:                         timeout,
		PollPeriod:                      pollPeriod,
		VmNetworkAdaptersWaitForIpsJson: string(vmNetworkAdaptersWaitForIpsJson),
	})

	return err
}

type updateVmNetworkAdapterArgs struct {
	VmName               string
	Index                int
	VmNetworkAdapterJson string
}

var updateVmNetworkAdapterTemplate = template.Must(template.New("UpdateVmNetworkAdapter").Parse(`
$ErrorActionPreference = 'Stop'
#First 3 requests fails to get ip address
Get-VMNetworkAdapter -VmName '{{.VmName}}' | Out-Null
Get-VMNetworkAdapter -VmName '{{.VmName}}' | Out-Null
Get-VMNetworkAdapter -VmName '{{.VmName}}' | Out-Null

$vmNetworkAdapter = '{{.VmNetworkAdapterJson}}' | ConvertFrom-Json

$dhcpGuard = [Microsoft.HyperV.PowerShell.OnOffState]$vmNetworkAdapter.DhcpGuard
$routerGuard = [Microsoft.HyperV.PowerShell.OnOffState]$vmNetworkAdapter.RouterGuard
$portMirroring = [Microsoft.HyperV.PowerShell.VMNetworkAdapterPortMirroringMode]$vmNetworkAdapter.PortMirroring
$ieeePriorityTag = [Microsoft.HyperV.PowerShell.OnOffState]$vmNetworkAdapter.IeeePriorityTag
$iovInterruptModeration = [Microsoft.HyperV.PowerShell.IovInterruptModerationValue]$vmNetworkAdapter.IovInterruptModeration
$allowTeaming = [Microsoft.HyperV.PowerShell.OnOffState]$vmNetworkAdapter.AllowTeaming
$deviceNaming = [Microsoft.HyperV.PowerShell.OnOffState]$vmNetworkAdapter.DeviceNaming
$fixSpeed10G = [Microsoft.HyperV.PowerShell.OnOffState]$vmNetworkAdapter.FixSpeed10G

$vmNetworkAdaptersObject = @(Get-VM -Name '{{.VmName}}*' | ?{$_.Name -eq '{{.VmName}}' } | Get-VMNetworkAdapter)[{{.Index}}]

if (!$vmNetworkAdaptersObject){
	throw "VM network adapter does not exist - {{.Index}}"
}

if ($vmNetworkAdapter.SwitchName) {
	$vmSwitch = Get-VMSwitch -Name $vmNetworkAdapter.SwitchName
	if ($vmSwitch) {
		$minimumBandwidthMode = $vmSwitch.BandwidthReservationMode
	}
}

if ($vmNetworkAdaptersObject.SwitchName -ne $vmNetworkAdapter.SwitchName) {
	if ($vmNetworkAdapter.SwitchName) {
		$null = $vmNetworkAdaptersObject | Connect-VMNetworkAdapter -SwitchName $vmNetworkAdapter.SwitchName
	} else {
		$null = $vmNetworkAdaptersObject | Disconnect-VMNetworkAdapter
	}
}

if ($vmNetworkAdaptersObject.Name -ne $vmNetworkAdapter.Name) {
	$null = $vmNetworkAdaptersObject | Rename-VMNetworkAdapter -NewName $vmNetworkAdapter.Name
}

$SetVmNetworkAdapterArgs = @{}
$SetVmNetworkAdapterArgs.VmName=$vmNetworkAdapter.VmName
$SetVmNetworkAdapterArgs.Name=$vmNetworkAdapter.Name
if ($vmNetworkAdapter.DynamicMacAddress) {
	$SetVmNetworkAdapterArgs.DynamicMacAddress=$vmNetworkAdapter.DynamicMacAddress
} elseif ($vmNetworkAdapter.StaticMacAddress) {
	$SetVmNetworkAdapterArgs.StaticMacAddress=$vmNetworkAdapter.StaticMacAddress
}
if ($vmNetworkAdapter.MacAddressSpoofing) {
	$SetVmNetworkAdapterArgs.MacAddressSpoofing=$vmNetworkAdapter.MacAddressSpoofing
}
$SetVmNetworkAdapterArgs.DhcpGuard=$dhcpGuard
$SetVmNetworkAdapterArgs.RouterGuard=$routerGuard
$SetVmNetworkAdapterArgs.PortMirroring=$portMirroring
$SetVmNetworkAdapterArgs.IeeePriorityTag=$ieeePriorityTag
$SetVmNetworkAdapterArgs.VmqWeight=$vmNetworkAdapter.VmqWeight
$SetVmNetworkAdapterArgs.IovQueuePairsRequested=$vmNetworkAdapter.IovQueuePairsRequested
$SetVmNetworkAdapterArgs.IovInterruptModeration=$iovInterruptModeration
$SetVmNetworkAdapterArgs.IovWeight=$vmNetworkAdapter.IovWeight
$SetVmNetworkAdapterArgs.IPsecOffloadMaximumSecurityAssociation=$vmNetworkAdapter.IPsecOffloadMaximumSecurityAssociation
$SetVmNetworkAdapterArgs.MaximumBandwidth=$vmNetworkAdapter.MaximumBandwidth
if ($minimumBandwidthMode -eq [Microsoft.HyperV.PowerShell.VMSwitchBandwidthMode]::Absolute){
	$SetVmNetworkAdapterArgs.MinimumBandwidthAbsolute=$vmNetworkAdapter.MinimumBandwidthAbsolute
}
if ($minimumBandwidthMode -eq [Microsoft.HyperV.PowerShell.VMSwitchBandwidthMode]::Weight -or $minimumBandwidthMode -eq [Microsoft.HyperV.PowerShell.VMSwitchBandwidthMode]::Default){
	$SetVmNetworkAdapterArgs.MinimumBandwidthWeight=$vmNetworkAdapter.MinimumBandwidthWeight
}
$SetVmNetworkAdapterArgs.MandatoryFeatureId=$vmNetworkAdapter.MandatoryFeatureId
if ($vmNetworkAdaptersObject.ResourcePoolName -ne $vmNetworkAdapter.ResourcePoolName) {
	if ($vmNetworkAdapter.ResourcePoolName) {
		$SetVmNetworkAdapterArgs.ResourcePoolName=$vmNetworkAdapter.ResourcePoolName
	} else {
		$null = $vmNetworkAdaptersObject | Disconnect-VMNetworkAdapter
	}
}

$SetVmNetworkAdapterArgs.TestReplicaPoolName=$vmNetworkAdapter.TestReplicaPoolName
$SetVmNetworkAdapterArgs.TestReplicaSwitchName=$vmNetworkAdapter.TestReplicaSwitchName
$SetVmNetworkAdapterArgs.VirtualSubnetId=$vmNetworkAdapter.VirtualSubnetId
$SetVmNetworkAdapterArgs.AllowTeaming=$allowTeaming
$SetVmNetworkAdapterArgs.NotMonitoredInCluster=$vmNetworkAdapter.NotMonitoredInCluster
$SetVmNetworkAdapterArgs.StormLimit=$vmNetworkAdapter.StormLimit
$SetVmNetworkAdapterArgs.DynamicIPAddressLimit=$vmNetworkAdapter.DynamicIPAddressLimit
$SetVmNetworkAdapterArgs.DeviceNaming=$deviceNaming
$SetVmNetworkAdapterArgs.FixSpeed10G=$fixSpeed10G
$SetVmNetworkAdapterArgs.PacketDirectNumProcs=$vmNetworkAdapter.PacketDirectNumProcs
$SetVmNetworkAdapterArgs.PacketDirectModerationCount=$vmNetworkAdapter.PacketDirectModerationCount
$SetVmNetworkAdapterArgs.PacketDirectModerationInterval=$vmNetworkAdapter.PacketDirectModerationInterval
$SetVmNetworkAdapterArgs.VrssEnabled=$vmNetworkAdapter.VrssEnabled
$SetVmNetworkAdapterArgs.VmmqEnabled=$vmNetworkAdapter.VmmqEnabled
$SetVmNetworkAdapterArgs.VmmqQueuePairs=$vmNetworkAdapter.VmmqQueuePairs

Set-VmNetworkAdapter @SetVmNetworkAdapterArgs

if ($vmNetworkAdapter.VlanAccess -and $vmNetworkAdapter.VlanId) {
	$SetVmNetworkAdapterVlanArgs = @{}

	$SetVmNetworkAdapterVlanArgs.VMName = $vmNetworkAdapter.VmName
	$SetVmNetworkAdapterVlanArgs.VMNetworkAdapterName = $vmNetworkAdapter.Name
	$SetVmNetworkAdapterVlanArgs.Access = $true
	$SetVmNetworkAdapterVlanArgs.VlanId = $vmNetworkAdapter.VlanId

	Set-VmNetworkAdapterVlan @SetVmNetworkAdapterVlanArgs
}

`))

func (c *ClientConfig) UpdateVmNetworkAdapter(
	ctx context.Context,
	vmName string,
	index int,
	name string,
	switchName string,
	managementOs bool,
	isLegacy bool,
	dynamicMacAddress bool,
	staticMacAddress string,
	macAddressSpoofing api.OnOffState,
	dhcpGuard api.OnOffState,
	routerGuard api.OnOffState,
	portMirroring api.PortMirroring,
	ieeePriorityTag api.OnOffState,
	vmqWeight int,
	iovQueuePairsRequested int,
	iovInterruptModeration api.IovInterruptModerationValue,
	iovWeight int,
	ipsecOffloadMaximumSecurityAssociation int,
	maximumBandwidth int,
	minimumBandwidthAbsolute int,
	minimumBandwidthWeight int,
	mandatoryFeatureId []string,
	resourcePoolName string,
	testReplicaPoolName string,
	testReplicaSwitchName string,
	virtualSubnetId int,
	allowTeaming api.OnOffState,
	notMonitoredInCluster bool,
	stormLimit int,
	dynamicIpAddressLimit int,
	deviceNaming api.OnOffState,
	fixSpeed10G api.OnOffState,
	packetDirectNumProcs int,
	packetDirectModerationCount int,
	packetDirectModerationInterval int,
	vrssEnabled bool,
	vmmqEnabled bool,
	vmmqQueuePairs int,
	vlanAccess bool,
	vlanId int,
) (err error) {

	vmNetworkAdapterJson, err := json.Marshal(api.VmNetworkAdapter{
		VmName:                                 vmName,
		Index:                                  index,
		Name:                                   name,
		SwitchName:                             switchName,
		ManagementOs:                           managementOs,
		IsLegacy:                               isLegacy,
		DynamicMacAddress:                      dynamicMacAddress,
		StaticMacAddress:                       staticMacAddress,
		MacAddressSpoofing:                     macAddressSpoofing,
		DhcpGuard:                              dhcpGuard,
		RouterGuard:                            routerGuard,
		PortMirroring:                          portMirroring,
		IeeePriorityTag:                        ieeePriorityTag,
		VmqWeight:                              vmqWeight,
		IovQueuePairsRequested:                 iovQueuePairsRequested,
		IovInterruptModeration:                 iovInterruptModeration,
		IovWeight:                              iovWeight,
		IpsecOffloadMaximumSecurityAssociation: ipsecOffloadMaximumSecurityAssociation,
		MaximumBandwidth:                       maximumBandwidth,
		MinimumBandwidthAbsolute:               minimumBandwidthAbsolute,
		MinimumBandwidthWeight:                 minimumBandwidthWeight,
		MandatoryFeatureId:                     mandatoryFeatureId,
		ResourcePoolName:                       resourcePoolName,
		TestReplicaPoolName:                    testReplicaPoolName,
		TestReplicaSwitchName:                  testReplicaSwitchName,
		VirtualSubnetId:                        virtualSubnetId,
		AllowTeaming:                           allowTeaming,
		NotMonitoredInCluster:                  notMonitoredInCluster,
		StormLimit:                             stormLimit,
		DynamicIpAddressLimit:                  dynamicIpAddressLimit,
		DeviceNaming:                           deviceNaming,
		FixSpeed10G:                            fixSpeed10G,
		PacketDirectNumProcs:                   packetDirectNumProcs,
		PacketDirectModerationCount:            packetDirectModerationCount,
		PacketDirectModerationInterval:         packetDirectModerationInterval,
		VrssEnabled:                            vrssEnabled,
		VmmqEnabled:                            vmmqEnabled,
		VmmqQueuePairs:                         vmmqQueuePairs,
		VlanAccess:                             vlanAccess,
		VlanId:                                 vlanId,
	})

	if err != nil {
		return err
	}

	err = c.WinRmClient.RunFireAndForgetScript(ctx, updateVmNetworkAdapterTemplate, updateVmNetworkAdapterArgs{
		VmName:               vmName,
		Index:                index,
		VmNetworkAdapterJson: string(vmNetworkAdapterJson),
	})

	return err
}

type deleteVmNetworkAdapterArgs struct {
	VmName string
	Index  int
}

var deleteVmNetworkAdapterTemplate = template.Must(template.New("DeleteVmNetworkAdapter").Parse(`
$ErrorActionPreference = 'Stop'

@(Get-VM -Name '{{.VmName}}*' | ?{$_.Name -eq '{{.VmName}}' } | Get-VMNetworkAdapter)[{{.Index}}] | Remove-VMNetworkAdapter
`))

func (c *ClientConfig) DeleteVmNetworkAdapter(ctx context.Context, vmName string, index int) (err error) {
	err = c.WinRmClient.RunFireAndForgetScript(ctx, deleteVmNetworkAdapterTemplate, deleteVmNetworkAdapterArgs{
		VmName: vmName,
		Index:  index,
	})

	return err
}

func (c *ClientConfig) CreateOrUpdateVmNetworkAdapters(ctx context.Context, vmName string, networkAdapters []api.VmNetworkAdapter) (err error) {
	networkAdaptersWaitForIps := make([]api.VmNetworkAdapterWaitForIp, 0)

	//Empty networkAdaptersWaitForIps is ok as we aren't using the results anywhere
	currentNetworkAdapters, err := c.GetVmNetworkAdapters(ctx, vmName, networkAdaptersWaitForIps)
	if err != nil {
		return err
	}

	currentNetworkAdaptersLength := len(currentNetworkAdapters)
	desiredNetworkAdaptersLength := len(networkAdapters)

	for i := currentNetworkAdaptersLength - 1; i > desiredNetworkAdaptersLength-1; i-- {
		currentNetworkAdapter := currentNetworkAdapters[i]
		err = c.DeleteVmNetworkAdapter(ctx, vmName, currentNetworkAdapter.Index)
		if err != nil {
			return err
		}
	}

	if currentNetworkAdaptersLength > desiredNetworkAdaptersLength {
		currentNetworkAdaptersLength = desiredNetworkAdaptersLength
	}

	for i := 0; i <= currentNetworkAdaptersLength-1; i++ {
		currentNetworkAdapter := currentNetworkAdapters[i]
		networkAdapter := networkAdapters[i]
		err = c.UpdateVmNetworkAdapter(
			ctx,
			vmName,
			currentNetworkAdapter.Index,
			networkAdapter.Name,
			networkAdapter.SwitchName,
			networkAdapter.ManagementOs,
			networkAdapter.IsLegacy,
			networkAdapter.DynamicMacAddress,
			networkAdapter.StaticMacAddress,
			networkAdapter.MacAddressSpoofing,
			networkAdapter.DhcpGuard,
			networkAdapter.RouterGuard,
			networkAdapter.PortMirroring,
			networkAdapter.IeeePriorityTag,
			networkAdapter.VmqWeight,
			networkAdapter.IovQueuePairsRequested,
			networkAdapter.IovInterruptModeration,
			networkAdapter.IovWeight,
			networkAdapter.IpsecOffloadMaximumSecurityAssociation,
			networkAdapter.MaximumBandwidth,
			networkAdapter.MinimumBandwidthAbsolute,
			networkAdapter.MinimumBandwidthWeight,
			networkAdapter.MandatoryFeatureId,
			networkAdapter.ResourcePoolName,
			networkAdapter.TestReplicaPoolName,
			networkAdapter.TestReplicaSwitchName,
			networkAdapter.VirtualSubnetId,
			networkAdapter.AllowTeaming,
			networkAdapter.NotMonitoredInCluster,
			networkAdapter.StormLimit,
			networkAdapter.DynamicIpAddressLimit,
			networkAdapter.DeviceNaming,
			networkAdapter.FixSpeed10G,
			networkAdapter.PacketDirectNumProcs,
			networkAdapter.PacketDirectModerationCount,
			networkAdapter.PacketDirectModerationInterval,
			networkAdapter.VrssEnabled,
			networkAdapter.VmmqEnabled,
			networkAdapter.VmmqQueuePairs,
			networkAdapter.VlanAccess,
			networkAdapter.VlanId,
		)
		if err != nil {
			return err
		}
	}

	for i := currentNetworkAdaptersLength - 1 + 1; i <= desiredNetworkAdaptersLength-1; i++ {
		networkAdapter := networkAdapters[i]
		err = c.CreateVmNetworkAdapter(
			ctx,
			vmName,
			networkAdapter.Name,
			networkAdapter.SwitchName,
			networkAdapter.ManagementOs,
			networkAdapter.IsLegacy,
			networkAdapter.DynamicMacAddress,
			networkAdapter.StaticMacAddress,
			networkAdapter.MacAddressSpoofing,
			networkAdapter.DhcpGuard,
			networkAdapter.RouterGuard,
			networkAdapter.PortMirroring,
			networkAdapter.IeeePriorityTag,
			networkAdapter.VmqWeight,
			networkAdapter.IovQueuePairsRequested,
			networkAdapter.IovInterruptModeration,
			networkAdapter.IovWeight,
			networkAdapter.IpsecOffloadMaximumSecurityAssociation,
			networkAdapter.MaximumBandwidth,
			networkAdapter.MinimumBandwidthAbsolute,
			networkAdapter.MinimumBandwidthWeight,
			networkAdapter.MandatoryFeatureId,
			networkAdapter.ResourcePoolName,
			networkAdapter.TestReplicaPoolName,
			networkAdapter.TestReplicaSwitchName,
			networkAdapter.VirtualSubnetId,
			networkAdapter.AllowTeaming,
			networkAdapter.NotMonitoredInCluster,
			networkAdapter.StormLimit,
			networkAdapter.DynamicIpAddressLimit,
			networkAdapter.DeviceNaming,
			networkAdapter.FixSpeed10G,
			networkAdapter.PacketDirectNumProcs,
			networkAdapter.PacketDirectModerationCount,
			networkAdapter.PacketDirectModerationInterval,
			networkAdapter.VrssEnabled,
			networkAdapter.VmmqEnabled,
			networkAdapter.VmmqQueuePairs,
			networkAdapter.VlanAccess,
			networkAdapter.VlanId,
		)

		if err != nil {
			return err
		}
	}

	return nil
}
