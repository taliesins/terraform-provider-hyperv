package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"strings"
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

func ExpandNetworkAdapters(d *schema.ResourceData) ([]VmNetworkAdapter, error) {
	expandedNetworkAdapters := make([]VmNetworkAdapter, 0)

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

			expandedNetworkAdapter := VmNetworkAdapter{
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

func FlattenNetworkAdapters(networkAdapters *[]VmNetworkAdapter) []interface{} {
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

type VmNetworkAdapterWaitForIp struct {
	Name       string
	WaitForIps bool
}

type VmNetworkAdapter struct {
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

func ExpandVmNetworkAdapterWaitForIps(d *schema.ResourceData) ([]VmNetworkAdapterWaitForIp, uint32, uint32, error) {
	expandVmNetworkAdapterWaitForIps := make([]VmNetworkAdapterWaitForIp, 0)
	waitForIpsTimeout := uint32((d.Get("wait_for_ips_timeout")).(int))
	waitForIpsPollPeriod := uint32((d.Get("wait_for_ips_poll_period")).(int))

	if v, ok := d.GetOk("network_adaptors"); ok {
		networkAdapters := v.([]interface{})

		for _, networkAdapter := range networkAdapters {
			networkAdapter, ok := networkAdapter.(map[string]interface{})
			if !ok {
				return nil, waitForIpsTimeout, waitForIpsPollPeriod, fmt.Errorf("[ERROR][hyperv] network_adaptors should be a Hash - was '%+v'", networkAdapter)
			}

			expandedNetworkAdapterWaitForIp := VmNetworkAdapterWaitForIp{
				Name:       networkAdapter["name"].(string),
				WaitForIps: networkAdapter["wait_for_ips"].(bool),
			}

			expandVmNetworkAdapterWaitForIps = append(expandVmNetworkAdapterWaitForIps, expandedNetworkAdapterWaitForIp)
		}
	}

	return expandVmNetworkAdapterWaitForIps, waitForIpsTimeout, waitForIpsPollPeriod, nil
}

type HypervVmNetworkAdapterClient interface {
	CreateVmNetworkAdapter(
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
	) (err error)
	WaitForVmNetworkAdaptersIps(
		vmName string,
		timeout uint32,
		pollPeriod uint32,
		vmNetworkAdaptersWaitForIps []VmNetworkAdapterWaitForIp,
	) (err error)
	GetVmNetworkAdapters(vmName string, networkAdaptersWaitForIps []VmNetworkAdapterWaitForIp) (result []VmNetworkAdapter, err error)
	UpdateVmNetworkAdapter(
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
	) (err error)
	DeleteVmNetworkAdapter(vmName string, index int) (err error)
	CreateOrUpdateVmNetworkAdapters(vmName string, networkAdapters []VmNetworkAdapter) (err error)
}
