package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type PortMirroring int

const (
	PortMirroring_None        PortMirroring = 0
	PortMirroring_Destination PortMirroring = 1
	PortMirroring_Source      PortMirroring = 2
)

var PortMirroring_name = map[PortMirroring]string{
	PortMirroring_None:        "None",
	PortMirroring_Destination: "Destination",
	PortMirroring_Source:      "Source",
}

var PortMirroring_value = map[string]PortMirroring{
	"none":        PortMirroring_None,
	"destination": PortMirroring_Destination,
	"source":      PortMirroring_Source,
}

func (x PortMirroring) String() string {
	return PortMirroring_name[x]
}

func ToPortMirroring(x string) PortMirroring {
	if integerValue, err := strconv.Atoi(x); err == nil {
		return PortMirroring(integerValue)
	}
	return PortMirroring_value[strings.ToLower(x)]
}

func (d *PortMirroring) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(d.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (d *PortMirroring) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		var i int
		err2 := json.Unmarshal(b, &i)
		if err2 == nil {
			*d = PortMirroring(i)
			return nil
		}

		return err
	}
	*d = ToPortMirroring(s)
	return nil
}

type IovInterruptModerationValue int

const (
	IovInterruptModerationValue_Default  IovInterruptModerationValue = 0
	IovInterruptModerationValue_Adaptive IovInterruptModerationValue = 1
	IovInterruptModerationValue_Off      IovInterruptModerationValue = 2
	IovInterruptModerationValue_Low      IovInterruptModerationValue = 100
	IovInterruptModerationValue_Medium   IovInterruptModerationValue = 200
	IovInterruptModerationValue_High     IovInterruptModerationValue = 300
)

var IovInterruptModerationValue_name = map[IovInterruptModerationValue]string{
	IovInterruptModerationValue_Default:  "Default",
	IovInterruptModerationValue_Adaptive: "Adaptive",
	IovInterruptModerationValue_Off:      "Off",
	IovInterruptModerationValue_Low:      "Low",
	IovInterruptModerationValue_Medium:   "Medium",
	IovInterruptModerationValue_High:     "High",
}

var IovInterruptModerationValue_value = map[string]IovInterruptModerationValue{
	"default":  IovInterruptModerationValue_Default,
	"adaptive": IovInterruptModerationValue_Adaptive,
	"off":      IovInterruptModerationValue_Off,
	"low":      IovInterruptModerationValue_Low,
	"medium":   IovInterruptModerationValue_Medium,
	"high":     IovInterruptModerationValue_High,
}

func (x IovInterruptModerationValue) String() string {
	return IovInterruptModerationValue_name[x]
}

func ToIovInterruptModerationValue(x string) IovInterruptModerationValue {
	if integerValue, err := strconv.Atoi(x); err == nil {
		return IovInterruptModerationValue(integerValue)
	}
	return IovInterruptModerationValue_value[strings.ToLower(x)]
}

func (d *IovInterruptModerationValue) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(d.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (d *IovInterruptModerationValue) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		var i int
		err2 := json.Unmarshal(b, &i)
		if err2 == nil {
			*d = IovInterruptModerationValue(i)
			return nil
		}

		return err
	}
	*d = ToIovInterruptModerationValue(s)
	return nil
}

func DiffSuppressVmStaticMacAddress(key, old, new string, d *schema.ResourceData) bool {
	//Static Mac Address has not been set, so we don't mind what ever value is automatically generated
	if new == "" {
		return true
	}

	return new == old
}

func ExpandNetworkAdapters(d *schema.ResourceData) ([]vmNetworkAdapter, error) {
	expandedNetworkAdapters := make([]vmNetworkAdapter, 0)

	if v, ok := d.GetOk("network_adaptors"); ok {
		networkAdapters := v.([]interface{})

		for _, networkAdapter := range networkAdapters {
			networkAdapter, ok := networkAdapter.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] network_adaptors should be a Hash - was '%+v'", networkAdapter)
			}

			mandatoryFeatureIdSet := networkAdapter["mandatory_feature_id"].(*schema.Set).List()
			mandatoryFeatureIds := make([]string, 0)
			for _, mandatoryFeatureId := range mandatoryFeatureIdSet {
				mandatoryFeatureIds = append(mandatoryFeatureIds, mandatoryFeatureId.(string))
			}

			ipAddressesSet := networkAdapter["ip_addresses"].([]interface{})
			ipAddresses := make([]string, 0)
			for _, ipAddress := range ipAddressesSet {
				ipAddresses = append(ipAddresses, ipAddress.(string))
			}

			expandedNetworkAdapter := vmNetworkAdapter{
				Name:                                   networkAdapter["name"].(string),
				SwitchName:                             networkAdapter["switch_name"].(string),
				ManagementOs:                           networkAdapter["management_os"].(bool),
				IsLegacy:                               networkAdapter["is_legacy"].(bool),
				DynamicMacAddress:                      networkAdapter["dynamic_mac_address"].(bool),
				StaticMacAddress:                       networkAdapter["static_mac_address"].(string),
				MacAddressSpoofing:                     ToOnOffState(networkAdapter["mac_address_spoofing"].(string)),
				DhcpGuard:                              ToOnOffState(networkAdapter["dhcp_guard"].(string)),
				RouterGuard:                            ToOnOffState(networkAdapter["router_guard"].(string)),
				PortMirroring:                          ToPortMirroring(networkAdapter["port_mirroring"].(string)),
				IeeePriorityTag:                        ToOnOffState(networkAdapter["ieee_priority_tag"].(string)),
				VmqWeight:                              networkAdapter["vmq_weight"].(int),
				IovQueuePairsRequested:                 networkAdapter["iov_queue_pairs_requested"].(int),
				IovInterruptModeration:                 ToIovInterruptModerationValue(networkAdapter["iov_interrupt_moderation"].(string)),
				IovWeight:                              networkAdapter["iov_weight"].(int),
				IpsecOffloadMaximumSecurityAssociation: networkAdapter["ipsec_offload_maximum_security_association"].(int),
				MaximumBandwidth:                       networkAdapter["maximum_bandwidth"].(int),
				MinimumBandwidthAbsolute:               networkAdapter["minimum_bandwidth_absolute"].(int),
				MinimumBandwidthWeight:                 networkAdapter["minimum_bandwidth_weight"].(int),
				MandatoryFeatureId:                     mandatoryFeatureIds,
				ResourcePoolName:                       networkAdapter["resource_pool_name"].(string),
				TestReplicaPoolName:                    networkAdapter["test_replica_pool_name"].(string),
				TestReplicaSwitchName:                  networkAdapter["test_replica_switch_name"].(string),
				VirtualSubnetId:                        networkAdapter["virtual_subnet_id"].(int),
				AllowTeaming:                           ToOnOffState(networkAdapter["allow_teaming"].(string)),
				NotMonitoredInCluster:                  networkAdapter["not_monitored_in_cluster"].(bool),
				StormLimit:                             networkAdapter["storm_limit"].(int),
				DynamicIpAddressLimit:                  networkAdapter["dynamic_ip_address_limit"].(int),
				DeviceNaming:                           ToOnOffState(networkAdapter["device_naming"].(string)),
				FixSpeed10G:                            ToOnOffState(networkAdapter["fix_speed_10g"].(string)),
				PacketDirectNumProcs:                   networkAdapter["packet_direct_num_procs"].(int),
				PacketDirectModerationCount:            networkAdapter["packet_direct_moderation_count"].(int),
				PacketDirectModerationInterval:         networkAdapter["packet_direct_moderation_interval"].(int),
				VrssEnabled:                            networkAdapter["vrss_enabled"].(bool),
				VmmqEnabled:                            networkAdapter["vmmq_enabled"].(bool),
				VmmqQueuePairs:                         networkAdapter["vmmq_queue_pairs"].(int),
				VlanAccess:                             networkAdapter["vlan_access"].(bool),
				VlanId:                                 networkAdapter["vlan_id"].(int),
				WaitForIps:                             networkAdapter["wait_for_ips"].(bool),
				IpAddresses:                            ipAddresses,
			}

			expandedNetworkAdapters = append(expandedNetworkAdapters, expandedNetworkAdapter)
		}
	}

	return expandedNetworkAdapters, nil
}

func FlattenMandatoryFeatureIds(mandatoryFeatureIdStrings []string) *schema.Set {
	var mandatoryFeatureIds []interface{}

	for _, mandatoryFeatureId := range mandatoryFeatureIdStrings {
		mandatoryFeatureIds = append(mandatoryFeatureIds, mandatoryFeatureId)
	}

	return schema.NewSet(schema.HashString, mandatoryFeatureIds)
}

func FlattenNetworkAdapters(networkAdapters *[]vmNetworkAdapter) []interface{} {
	flattenedNetworkAdapters := make([]interface{}, 0)

	if networkAdapters != nil {
		for _, networkAdapter := range *networkAdapters {
			flattenedNetworkAdapter := make(map[string]interface{})
			flattenedNetworkAdapter["name"] = networkAdapter.Name
			flattenedNetworkAdapter["switch_name"] = networkAdapter.SwitchName
			flattenedNetworkAdapter["management_os"] = networkAdapter.ManagementOs
			flattenedNetworkAdapter["is_legacy"] = networkAdapter.IsLegacy
			flattenedNetworkAdapter["dynamic_mac_address"] = networkAdapter.DynamicMacAddress
			flattenedNetworkAdapter["static_mac_address"] = networkAdapter.StaticMacAddress
			flattenedNetworkAdapter["mac_address_spoofing"] = networkAdapter.MacAddressSpoofing.String()
			flattenedNetworkAdapter["dhcp_guard"] = networkAdapter.DhcpGuard.String()
			flattenedNetworkAdapter["router_guard"] = networkAdapter.RouterGuard.String()
			flattenedNetworkAdapter["port_mirroring"] = networkAdapter.PortMirroring.String()
			flattenedNetworkAdapter["ieee_priority_tag"] = networkAdapter.IeeePriorityTag.String()
			flattenedNetworkAdapter["vmq_weight"] = networkAdapter.VmqWeight
			flattenedNetworkAdapter["iov_queue_pairs_requested"] = networkAdapter.IovQueuePairsRequested
			flattenedNetworkAdapter["iov_interrupt_moderation"] = networkAdapter.IovInterruptModeration.String()
			flattenedNetworkAdapter["iov_weight"] = networkAdapter.IovWeight
			flattenedNetworkAdapter["ipsec_offload_maximum_security_association"] = networkAdapter.IpsecOffloadMaximumSecurityAssociation
			flattenedNetworkAdapter["maximum_bandwidth"] = networkAdapter.MaximumBandwidth
			flattenedNetworkAdapter["minimum_bandwidth_absolute"] = networkAdapter.MinimumBandwidthAbsolute
			flattenedNetworkAdapter["minimum_bandwidth_weight"] = networkAdapter.MinimumBandwidthWeight
			flattenedNetworkAdapter["mandatory_feature_id"] = FlattenMandatoryFeatureIds(networkAdapter.MandatoryFeatureId)
			flattenedNetworkAdapter["resource_pool_name"] = networkAdapter.ResourcePoolName
			flattenedNetworkAdapter["test_replica_pool_name"] = networkAdapter.TestReplicaPoolName
			flattenedNetworkAdapter["test_replica_switch_name"] = networkAdapter.TestReplicaSwitchName
			flattenedNetworkAdapter["virtual_subnet_id"] = networkAdapter.VirtualSubnetId
			flattenedNetworkAdapter["allow_teaming"] = networkAdapter.AllowTeaming.String()
			flattenedNetworkAdapter["not_monitored_in_cluster"] = networkAdapter.NotMonitoredInCluster
			flattenedNetworkAdapter["storm_limit"] = networkAdapter.StormLimit
			flattenedNetworkAdapter["dynamic_ip_address_limit"] = networkAdapter.DynamicIpAddressLimit
			flattenedNetworkAdapter["device_naming"] = networkAdapter.DeviceNaming.String()
			flattenedNetworkAdapter["fix_speed_10g"] = networkAdapter.FixSpeed10G.String()
			flattenedNetworkAdapter["packet_direct_num_procs"] = networkAdapter.PacketDirectNumProcs
			flattenedNetworkAdapter["packet_direct_moderation_count"] = networkAdapter.PacketDirectModerationCount
			flattenedNetworkAdapter["packet_direct_moderation_interval"] = networkAdapter.PacketDirectModerationInterval
			flattenedNetworkAdapter["vrss_enabled"] = networkAdapter.VrssEnabled
			flattenedNetworkAdapter["vmmq_enabled"] = networkAdapter.VmmqEnabled
			flattenedNetworkAdapter["vmmq_queue_pairs"] = networkAdapter.VmmqQueuePairs
			flattenedNetworkAdapter["vlan_access"] = networkAdapter.VlanAccess
			flattenedNetworkAdapter["vlan_id"] = networkAdapter.VlanId
			flattenedNetworkAdapter["wait_for_ips"] = networkAdapter.WaitForIps
			flattenedNetworkAdapter["ip_addresses"] = networkAdapter.IpAddresses

			flattenedNetworkAdapters = append(flattenedNetworkAdapters, flattenedNetworkAdapter)
		}
	}

	return flattenedNetworkAdapters
}

type vmNetworkAdapterWaitForIp struct {
	Name       string
	WaitForIps bool
}

type vmNetworkAdapter struct {
	VmName                                 string
	Index                                  int
	Name                                   string
	SwitchName                             string
	ManagementOs                           bool
	IsLegacy                               bool
	DynamicMacAddress                      bool
	StaticMacAddress                       string
	MacAddressSpoofing                     OnOffState
	DhcpGuard                              OnOffState
	RouterGuard                            OnOffState
	PortMirroring                          PortMirroring
	IeeePriorityTag                        OnOffState
	VmqWeight                              int
	IovQueuePairsRequested                 int
	IovInterruptModeration                 IovInterruptModerationValue
	IovWeight                              int
	IpsecOffloadMaximumSecurityAssociation int
	MaximumBandwidth                       int
	MinimumBandwidthAbsolute               int
	MinimumBandwidthWeight                 int
	MandatoryFeatureId                     []string
	ResourcePoolName                       string
	TestReplicaPoolName                    string
	TestReplicaSwitchName                  string
	VirtualSubnetId                        int
	AllowTeaming                           OnOffState
	NotMonitoredInCluster                  bool
	StormLimit                             int
	DynamicIpAddressLimit                  int
	DeviceNaming                           OnOffState
	FixSpeed10G                            OnOffState
	PacketDirectNumProcs                   int
	PacketDirectModerationCount            int
	PacketDirectModerationInterval         int
	VrssEnabled                            bool
	VmmqEnabled                            bool
	VmmqQueuePairs                         int
	VlanAccess                             bool
	VlanId                                 int
	WaitForIps                             bool
	IpAddresses                            []string
}

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

func (c *HypervClient) CreateVmNetworkAdapter(
	vmName string,
	name string,
	switchName string,
	managementOs bool,
	isLegacy bool,
	dynamicMacAddress bool,
	staticMacAddress string,
	macAddressSpoofing OnOffState,
	dhcpGuard OnOffState,
	routerGuard OnOffState,
	portMirroring PortMirroring,
	ieeePriorityTag OnOffState,
	vmqWeight int,
	iovQueuePairsRequested int,
	iovInterruptModeration IovInterruptModerationValue,
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
	allowTeaming OnOffState,
	notMonitoredInCluster bool,
	stormLimit int,
	dynamicIpAddressLimit int,
	deviceNaming OnOffState,
	fixSpeed10G OnOffState,
	packetDirectNumProcs int,
	packetDirectModerationCount int,
	packetDirectModerationInterval int,
	vrssEnabled bool,
	vmmqEnabled bool,
	vmmqQueuePairs int,
	vlanAccess bool,
	vlanId int,
) (err error) {

	vmNetworkAdapterJson, err := json.Marshal(vmNetworkAdapter{
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

	err = c.runFireAndForgetScript(createVmNetworkAdapterTemplate, createVmNetworkAdapterArgs{
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

$vmNetworkAdaptersObject = @(Get-VMNetworkAdapter -VmName '{{.VmName}}' | %{ @{
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

func (c *HypervClient) GetVmNetworkAdapters(vmName string, networkAdaptersWaitForIps []vmNetworkAdapterWaitForIp) (result []vmNetworkAdapter, err error) {
	result = make([]vmNetworkAdapter, 0)

	err = c.runScriptWithResult(getVmNetworkAdaptersTemplate, getVmNetworkAdaptersArgs{
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

func (c *HypervClient) WaitForVmNetworkAdaptersIps(
	vmName string,
	timeout uint32,
	pollPeriod uint32,
	vmNetworkAdaptersWaitForIps []vmNetworkAdapterWaitForIp,
) (err error) {

	vmNetworkAdaptersWaitForIpsJson, err := json.Marshal(vmNetworkAdaptersWaitForIps)

	if err != nil {
		return err
	}

	err = c.runFireAndForgetScript(waitForVmNetworkAdaptersIpsTemplate, waitForVmNetworkAdaptersIpsArgs{
		VmName:                          vmName,
		Timeout:                         timeout,
		PollPeriod:                      pollPeriod,
		VmNetworkAdaptersWaitForIpsJson: string(vmNetworkAdaptersWaitForIpsJson),
	})

	return err
}

func ExpandVmNetworkAdapterWaitForIps(d *schema.ResourceData) ([]vmNetworkAdapterWaitForIp, uint32, uint32, error) {
	expandVmNetworkAdapterWaitForIps := make([]vmNetworkAdapterWaitForIp, 0)
	waitForIpsTimeout := uint32((d.Get("wait_for_ips_timeout")).(int))
	waitForIpsPollPeriod := uint32((d.Get("wait_for_ips_poll_period")).(int))

	if v, ok := d.GetOk("network_adaptors"); ok {
		networkAdapters := v.([]interface{})

		for _, networkAdapter := range networkAdapters {
			networkAdapter, ok := networkAdapter.(map[string]interface{})
			if !ok {
				return nil, waitForIpsTimeout, waitForIpsPollPeriod, fmt.Errorf("[ERROR][hyperv] network_adaptors should be a Hash - was '%+v'", networkAdapter)
			}

			expandedNetworkAdapterWaitForIp := vmNetworkAdapterWaitForIp{
				Name:       networkAdapter["name"].(string),
				WaitForIps: networkAdapter["wait_for_ips"].(bool),
			}

			expandVmNetworkAdapterWaitForIps = append(expandVmNetworkAdapterWaitForIps, expandedNetworkAdapterWaitForIp)
		}
	}

	return expandVmNetworkAdapterWaitForIps, waitForIpsTimeout, waitForIpsPollPeriod, nil
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

$vmNetworkAdaptersObject = @(Get-VMNetworkAdapter -VmName '{{.VmName}}')[{{.Index}}]

if (!$vmNetworkAdaptersObject){
	throw "VM network adapter does not exist - {{.Index}}"
}

if ($vmNetworkAdapter.SwitchName) {
	$vmSwitch = Get-VMSwitch -Name $vmNetworkAdapter.SwitchName
	if ($vmSwitch) {
		$minimumBandwidthMode = $vmSwitch.BandwidthReservationMode
	}
	$null = Get-VMNetworkAdapter -VmName '{{.VmName}}' -Name "$($vmNetworkAdapter.Name)" | Connect-VMNetworkAdapter -SwitchName "$($vmNetworkAdapter.SwitchName)"
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

	$SetVmNetworkAdapterVlanArgs.VMName = $vmNetworkAdapter.VmName
	$SetVmNetworkAdapterVlanArgs.VMNetworkAdapterName = $vmNetworkAdapter.Name
	$SetVmNetworkAdapterVlanArgs.Access = $true
	$SetVmNetworkAdapterVlanArgs.VlanId = $vmNetworkAdapter.VlanId

	Set-VmNetworkAdapterVlan @SetVmNetworkAdapterVlanArgs
}

`))

func (c *HypervClient) UpdateVmNetworkAdapter(
	vmName string,
	index int,
	name string,
	switchName string,
	managementOs bool,
	isLegacy bool,
	dynamicMacAddress bool,
	staticMacAddress string,
	macAddressSpoofing OnOffState,
	dhcpGuard OnOffState,
	routerGuard OnOffState,
	portMirroring PortMirroring,
	ieeePriorityTag OnOffState,
	vmqWeight int,
	iovQueuePairsRequested int,
	iovInterruptModeration IovInterruptModerationValue,
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
	allowTeaming OnOffState,
	notMonitoredInCluster bool,
	stormLimit int,
	dynamicIpAddressLimit int,
	deviceNaming OnOffState,
	fixSpeed10G OnOffState,
	packetDirectNumProcs int,
	packetDirectModerationCount int,
	packetDirectModerationInterval int,
	vrssEnabled bool,
	vmmqEnabled bool,
	vmmqQueuePairs int,
	vlanAccess bool,
	vlanId int,
) (err error) {

	vmNetworkAdapterJson, err := json.Marshal(vmNetworkAdapter{
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

	err = c.runFireAndForgetScript(updateVmNetworkAdapterTemplate, updateVmNetworkAdapterArgs{
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

@(Get-VMNetworkAdapter -VmName '{{.VmName}}')[{{.Index}}] | Remove-VMNetworkAdapter -Force
`))

func (c *HypervClient) DeleteVmNetworkAdapter(vmName string, index int) (err error) {
	err = c.runFireAndForgetScript(deleteVmNetworkAdapterTemplate, deleteVmNetworkAdapterArgs{
		VmName: vmName,
		Index:  index,
	})

	return err
}

func (c *HypervClient) CreateOrUpdateVmNetworkAdapters(vmName string, networkAdapters []vmNetworkAdapter) (err error) {
	networkAdaptersWaitForIps := make([]vmNetworkAdapterWaitForIp, 0)

	//Empty networkAdaptersWaitForIps is ok as we aren't using the results anywhere
	currentNetworkAdapters, err := c.GetVmNetworkAdapters(vmName, networkAdaptersWaitForIps)
	if err != nil {
		return err
	}

	currentNetworkAdaptersLength := len(currentNetworkAdapters)
	desiredNetworkAdaptersLength := len(networkAdapters)

	for i := currentNetworkAdaptersLength - 1; i > desiredNetworkAdaptersLength-1; i-- {
		currentNetworkAdapter := currentNetworkAdapters[i]
		err = c.DeleteVmNetworkAdapter(vmName, currentNetworkAdapter.Index)
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
