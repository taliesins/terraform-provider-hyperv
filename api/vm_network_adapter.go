package api

import (
	"encoding/json"
	"strings"
	"text/template"
	"github.com/hashicorp/terraform/helper/schema"
	"fmt"
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
	return PortMirroring_value[strings.ToLower(x)]
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
	return IovInterruptModerationValue_value[strings.ToLower(x)]
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

			expandedNetworkAdapter := vmNetworkAdapter{
				Name:                                   networkAdapter["name"].(string),
				SwitchName:                             networkAdapter["switch_name"].(string),
				ManagementOs:                           networkAdapter["management_os"].(bool),
				IsLegacy:                               networkAdapter["is_legacy"].(bool),
				DynamicMacAddress:                      networkAdapter["dynamic_mac_address"].(bool),
				StaticMacAddress:                       networkAdapter["static_mac_address"].(string),
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
				MandatoryFeatureId:                     networkAdapter["mandatory_feature_id"].(string),
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
			}

			expandedNetworkAdapters = append(expandedNetworkAdapters, expandedNetworkAdapter)
		}
	}

	return expandedNetworkAdapters, nil
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
			flattenedNetworkAdapter["dhcp_guard"] = networkAdapter.DhcpGuard
			flattenedNetworkAdapter["router_guard"] = networkAdapter.RouterGuard
			flattenedNetworkAdapter["port_mirroring"] = networkAdapter.PortMirroring
			flattenedNetworkAdapter["ieee_priority_tag"] = networkAdapter.IeeePriorityTag
			flattenedNetworkAdapter["vmq_weight"] = networkAdapter.VmqWeight
			flattenedNetworkAdapter["iov_queue_pairs_requested"] = networkAdapter.IovQueuePairsRequested
			flattenedNetworkAdapter["iov_interrupt_moderation"] = networkAdapter.IovInterruptModeration
			flattenedNetworkAdapter["iov_weight"] = networkAdapter.IovWeight
			flattenedNetworkAdapter["ipsec_offload_maximum_security_association"] = networkAdapter.IpsecOffloadMaximumSecurityAssociation
			flattenedNetworkAdapter["maximum_bandwidth"] = networkAdapter.MaximumBandwidth
			flattenedNetworkAdapter["minimum_bandwidth_absolute"] = networkAdapter.MinimumBandwidthAbsolute
			flattenedNetworkAdapter["minimum_bandwidth_weight"] = networkAdapter.MinimumBandwidthWeight
			flattenedNetworkAdapter["mandatory_feature_id"] = networkAdapter.MandatoryFeatureId
			flattenedNetworkAdapter["resource_pool_name"] = networkAdapter.ResourcePoolName
			flattenedNetworkAdapter["test_replica_pool_name"] = networkAdapter.TestReplicaPoolName
			flattenedNetworkAdapter["test_replica_switch_name"] = networkAdapter.TestReplicaSwitchName
			flattenedNetworkAdapter["virtual_subnet_id"] = networkAdapter.VirtualSubnetId
			flattenedNetworkAdapter["allow_teaming"] = networkAdapter.AllowTeaming
			flattenedNetworkAdapter["not_monitored_in_cluster"] = networkAdapter.NotMonitoredInCluster
			flattenedNetworkAdapter["storm_limit"] = networkAdapter.StormLimit
			flattenedNetworkAdapter["dynamic_ip_address_limit"] = networkAdapter.DynamicIpAddressLimit
			flattenedNetworkAdapter["device_naming"] = networkAdapter.DeviceNaming
			flattenedNetworkAdapter["fix_speed_10g"] = networkAdapter.FixSpeed10G
			flattenedNetworkAdapter["packet_direct_num_procs"] = networkAdapter.PacketDirectNumProcs
			flattenedNetworkAdapter["packet_direct_moderation_count"] = networkAdapter.PacketDirectModerationCount
			flattenedNetworkAdapter["packet_direct_moderation_interval"] = networkAdapter.PacketDirectModerationInterval
			flattenedNetworkAdapter["vrss_enabled"] = networkAdapter.VrssEnabled
			flattenedNetworkAdapter["vmmq_enabled"] = networkAdapter.VmmqEnabled
			flattenedNetworkAdapter["vmmq_queue_pairs"] = networkAdapter.VmmqQueuePairs

			flattenedNetworkAdapters = append(flattenedNetworkAdapters, flattenedNetworkAdapter)
		}
	}

	return flattenedNetworkAdapters
}

type vmNetworkAdapter struct {
	VMName                                 string
	Index                                  int
	Name                                   string
	SwitchName                             string
	ManagementOs                           bool
	IsLegacy                               bool
	DynamicMacAddress                      bool
	StaticMacAddress                       string
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
	MandatoryFeatureId                     string
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
}

type createVMNetworkAdapterArgs struct {
	VmNetworkAdapterJson string
}

var createVMNetworkAdapterTemplate = template.Must(template.New("CreateVMNetworkAdapter").Parse(`
$ErrorActionPreference = 'Stop'
Get-Vm | Out-Null
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
	VMName=$vmNetworkAdapter.VmName
	Name=$vmNetworkAdapter.Name
	IsLegacy=$vmNetworkAdapter.IsLegacy
	SwitchName=$vmNetworkAdapter.SwitchName
}

Add-VmNetworkAdapter @NewVmNetworkAdapterArgs

$SetVmNetworkAdapterArgs = @{}
$SetVmNetworkAdapterArgs.VMName=$vmNetworkAdapter.VMName
$SetVmNetworkAdapterArgs.Name=$vmNetworkAdapter.Name
if ($vmNetworkAdapter.DynamicMacAddress) {
	$SetVmNetworkAdapterArgs.DynamicMacAddress=$vmNetworkAdapter.DynamicMacAddress
} else {
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
$SetVmNetworkAdapterArgs.MinimumBandwidthAbsolute=$vmNetworkAdapter.MinimumBandwidthAbsolute
$SetVmNetworkAdapterArgs.MinimumBandwidthWeight=$vmNetworkAdapter.MinimumBandwidthWeight
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

`))

func (c *HypervClient) CreateVMNetworkAdapter(
	vmName string,
	name string,
	switchName string,
	managementOs bool,
	isLegacy bool,
	dynamicMacAddress bool,
	staticMacAddress string,
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
	mandatoryFeatureId string,
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
) (err error) {

	vmNetworkAdapterJson, err := json.Marshal(vmNetworkAdapter{
		VMName:                                 vmName,
		Name:                                   name,
		SwitchName:                             switchName,
		ManagementOs:                           managementOs,
		IsLegacy:                               isLegacy,
		DynamicMacAddress:                      dynamicMacAddress,
		StaticMacAddress:                       staticMacAddress,
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
	})

	err = c.runFireAndForgetScript(createVMNetworkAdapterTemplate, createVMNetworkAdapterArgs{
		VmNetworkAdapterJson: string(vmNetworkAdapterJson),
	})

	return err
}

type getVMNetworkAdaptersArgs struct {
	VMName string
}

var getVMNetworkAdaptersTemplate = template.Must(template.New("GetVMNetworkAdapters").Parse(`
$ErrorActionPreference = 'Stop'
$vmNetworkAdaptersObject = Get-VMNetworkAdapter -VMName '{{.VMName}}' | %{ @{
	Name=$_.Name;
	SwitchName=$_.SwitchName;
	ManagementOs=$_.ManagementOs;
	IsLegacy=$_.IsLegacy;
	DynamicMacAddress=$_.DynamicMacAddress;
	StaticMacAddress=$_.StaticMacAddress;
	DhcpGuard=$_.DhcpGuard;
	RouterGuard=$_.RouterGuard;
	PortMirroring=$_.PortMirroring;
	IeeePriorityTag=$_.IeeePriorityTag;
	VmqWeight=$_.VmqWeight;
	IovQueuePairsRequested=$_.IovQueuePairsRequested;
	IovInterruptModeration=$_.IovInterruptModeration;
	IovWeight=$_.IovWeight;
	IpsecOffloadMaximumSecurityAssociation=$_.IpsecOffloadMaximumSecurityAssociation;
	MaximumBandwidth=$_.MaximumBandwidth;
	MinimumBandwidthAbsolute=$_.MinimumBandwidthAbsolute;
	MinimumBandwidthWeight=$_.MinimumBandwidthWeight;
	MandatoryFeatureId=$_.MandatoryFeatureId;
	ResourcePoolName=$_.ResourcePoolName;
	TestReplicaPoolName=$_.TestReplicaPoolName;
	TestReplicaSwitchName=$_.TestReplicaSwitchName;
	VirtualSubnetId=$_.VirtualSubnetId;
	AllowTeaming=$_.AllowTeaming;
	NotMonitoredInCluster=$_.NotMonitoredInCluster;
	StormLimit=$_.StormLimit;
	DynamicIpAddressLimit=$_.DynamicIpAddressLimit;
	DeviceNaming=$_.DeviceNaming;
	FixSpeed10G=$_.FixSpeed10G;
	PacketDirectNumProcs=$_.PacketDirectNumProcs;
	PacketDirectModerationCount=$_.PacketDirectModerationCount;
	PacketDirectModerationInterval=$_.PacketDirectModerationInterval;
	VrssEnabled=$_.VrssEnabled;
	VmmqEnabled=$_.VmmqEnabled;
	VmmqQueuePairs=$_.VmmqQueuePairs;
}}

if ($vmNetworkAdaptersObject) {
	$vmNetworkAdapters = ConvertTo-Json -InputObject $vmNetworkAdaptersObject
	$vmNetworkAdapters
} else {
	"[]"
}
`))

func (c *HypervClient) GetVMNetworkAdapters(vmname string) (result []vmNetworkAdapter, err error) {
	result = make([]vmNetworkAdapter, 0)

	err = c.runScriptWithResult(getVMNetworkAdaptersTemplate, getVMNetworkAdaptersArgs{
		VMName: vmname,
	}, result)

	return result, err
}

type updateVMNetworkAdapterArgs struct {
	VMName               string
	Index                int
	VmNetworkAdapterJson string
}

var updateVMNetworkAdapterTemplate = template.Must(template.New("UpdateVMNetworkAdapter").Parse(`
$ErrorActionPreference = 'Stop'
Get-Vm | Out-Null
$vmNetworkAdapter = '{{.VmNetworkAdapterJson}}' | ConvertFrom-Json

$dhcpGuard = [Microsoft.HyperV.PowerShell.OnOffState]$vmNetworkAdapter.DhcpGuard
$routerGuard = [Microsoft.HyperV.PowerShell.OnOffState]$vmNetworkAdapter.RouterGuard
$portMirroring = [Microsoft.HyperV.PowerShell.VMNetworkAdapterPortMirroringMode]$vmNetworkAdapter.PortMirroring
$ieeePriorityTag = [Microsoft.HyperV.PowerShell.OnOffState]$vmNetworkAdapter.IeeePriorityTag
$iovInterruptModeration = [Microsoft.HyperV.PowerShell.IovInterruptModerationValue]$vmNetworkAdapter.IovInterruptModeration
$allowTeaming = [Microsoft.HyperV.PowerShell.OnOffState]$vmNetworkAdapter.AllowTeaming
$deviceNaming = [Microsoft.HyperV.PowerShell.OnOffState]$vmNetworkAdapter.DeviceNaming
$fixSpeed10G = [Microsoft.HyperV.PowerShell.OnOffState]$vmNetworkAdapter.FixSpeed10G

$vmNetworkAdaptersObject = @(Get-VMNetworkAdapter -VMName '{{.VMName}}')[{{.Index}}]

if (!$vmNetworkAdaptersObject){
	throw "VM network adapter does not exist - {{.Index}}"
}

$SetVmNetworkAdapterArgs = @{}
$SetVmNetworkAdapterArgs.VMName=$vmNetworkAdapter.VMName
$SetVmNetworkAdapterArgs.Name=$vmNetworkAdapter.Name
if ($vmNetworkAdapter.DynamicMacAddress) {
	$SetVmNetworkAdapterArgs.DynamicMacAddress=$vmNetworkAdapter.DynamicMacAddress
} else {
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
$SetVmNetworkAdapterArgs.MinimumBandwidthAbsolute=$vmNetworkAdapter.MinimumBandwidthAbsolute
$SetVmNetworkAdapterArgs.MinimumBandwidthWeight=$vmNetworkAdapter.MinimumBandwidthWeight
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

`))

func (c *HypervClient) UpdateVMNetworkAdapter(
	vmName string,
	index int,
	name string,
	switchName string,
	managementOs bool,
	isLegacy bool,
	dynamicMacAddress bool,
	staticMacAddress string,
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
	mandatoryFeatureId string,
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
) (err error) {

	vmNetworkAdapterJson, err := json.Marshal(vmNetworkAdapter{
		VMName:                                 vmName,
		Index:                                  index,
		Name:                                   name,
		SwitchName:                             switchName,
		ManagementOs:                           managementOs,
		IsLegacy:                               isLegacy,
		DynamicMacAddress:                      dynamicMacAddress,
		StaticMacAddress:                       staticMacAddress,
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
	})

	err = c.runFireAndForgetScript(updateVMNetworkAdapterTemplate, updateVMNetworkAdapterArgs{
		VMName:               vmName,
		Index:                index,
		VmNetworkAdapterJson: string(vmNetworkAdapterJson),
	})

	return err
}

type deleteVMNetworkAdapterArgs struct {
	VMName string
	Index  int
}

var deleteVMNetworkAdapterTemplate = template.Must(template.New("DeleteVMNetworkAdapter").Parse(`
$ErrorActionPreference = 'Stop'

@(Get-VMNetworkAdapter -VMName '{{.VMName}}')[{{.Index}}] | Remove-VMNetworkAdapter -Force
`))

func (c *HypervClient) DeleteVMNetworkAdapter(vmName string, index int) (err error) {
	err = c.runFireAndForgetScript(deleteVMNetworkAdapterTemplate, deleteVMNetworkAdapterArgs{
		VMName: vmName,
		Index:  index,
	})

	return err
}

func (c *HypervClient) CreateOrUpdateVMNetworkAdapters(vmName string, networkAdapters []vmNetworkAdapter) (err error) {
	currentNetworkAdapters, err := c.GetVMNetworkAdapters(vmName)
	if err != nil {
		return err
	}

	currentNetworkAdaptersLength := len(currentNetworkAdapters)
	desiredNetworkAdaptersLength := len(networkAdapters)

	for i := currentNetworkAdaptersLength - 1; i > desiredNetworkAdaptersLength-1; i-- {
		currentNetworkAdapter := currentNetworkAdapters[i]
		err = c.DeleteVMNetworkAdapter(vmName, currentNetworkAdapter.Index)
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
		err = c.UpdateVMNetworkAdapter(
			vmName,
			currentNetworkAdapter.Index,
			networkAdapter.Name,
			networkAdapter.SwitchName,
			networkAdapter.ManagementOs,
			networkAdapter.IsLegacy,
			networkAdapter.DynamicMacAddress,
			networkAdapter.StaticMacAddress,
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
		)
		if err != nil {
			return err
		}
	}

	for i := currentNetworkAdaptersLength - 1 + 1; i <= desiredNetworkAdaptersLength-1; i++ {
		networkAdapter := networkAdapters[i]
		err = c.CreateVMNetworkAdapter(
			vmName,
			networkAdapter.Name,
			networkAdapter.SwitchName,
			networkAdapter.ManagementOs,
			networkAdapter.IsLegacy,
			networkAdapter.DynamicMacAddress,
			networkAdapter.StaticMacAddress,
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
		)

		if err != nil {
			return err
		}
	}

	return nil
}
