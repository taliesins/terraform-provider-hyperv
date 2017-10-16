---
layout: "hyperv"
page_title: "HyperV: hyperv_network_switch"
sidebar_current: "docs-hyperv-resource-network-switch"
description: |-
  Creates and manages network switch.
---

# hyperv\_network\_switch

The ``hyperv_network_switch`` resource creates and manages a network switch on a HyperV environment.

## Example Usage

```hcl
resource "hyperv_network_switch" "default" {
  name = "DMZ"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the switch to be created.

* `notes` - (Optional) Specifies a note to be associated with the switch to be created.

* `allow_management_os` - (Optional) `false` (default). Specifies if the HyperV host machine will have access to network switch when created. It provides this access via a virtual adaptor, so you will need to either configure static ips on the virtual adaptor or configure a dhcp on a machine connected to the network switch.

* `enable_embedded_teaming` - (Optional) `false` (default). Specifies if the HyperV host machine will enable teaming for network switch when created. It allows NIC teaming so that you could support scenrios such as redundent links. 

* `enable_iov` - (Optional) `false` (default). Specifies if the HyperV host machine will enable IO virtualization for network switch when created. If your hardware supports it, it enables the virtual machine to talk directly to the NIC.

* `enable_packet_direct` - (Optional) `false` (default). Specifies if the HyperV host machine will enable packet direct path for network switch when created. Increases packet throughoutput and reduces the network latency between vms on the switch.

* `minimum_bandwidth_mode` - (Optional) `0`(default). Specifies how minimum bandwidth is to be configured on the virtual switch. Allowed values are `Absolute`, `Default`, `None`, or `Weight`. If `Absolute` is specified, minimum bandwidth is bits per second. If `Weight` is specified, minimum bandwidth is a value ranging from `1` to `100`. If `None` is specified, minimum bandwidth is disabled on the switch – that is, users cannot configure it on any network adapter connected to the switch. If `Default` is specified, the system will set the mode to Weight, if the switch is not IOV-enabled, or `None` if the switch is IOV-enabled.

* `switch_type` - (Optional) `0`(default). Specifies the type of the switch to be created. Allowed values are `Internal` and `Private`. To create an `External` virtual switch, specify either the NetAdapterInterfaceDescription or the NetAdapterName parameter, which implicitly set the type of the virtual switch to `External`.

* `net_adapter_interface_descriptions` - (Optional) `string[]` (default). Specifies the interface description of the network adapter to be bound to the switch to be created. You can use the Get-NetAdapter cmdlet to get the interface description of a network adapter.

* `net_adapter_names` - (Optional) `string[]` (default). Specifies the name of the network adapter to be bound to the switch to be created. You can use the Get-NetAdapter cmdlet to get the interface description of a network adapter.

* `default_flow_minimum_bandwidth_absolute` - (Optional) `0` (default). Specifies the minimum bandwidth, in bits per second, that is allocated to a special category called "default flow." Any traffic sent by a virtual network adapter that is connected to this virtual switch and does not have minimum bandwidth allocated is filtered into this category. Specify a value for this parameter only if the minimum bandwidth mode on this virtual switch is absolute (See the New-VMSwitch cmdlet). By default, the virtual switch allocates 10% of the total bandwidth, which depends on the physical network adapter it binds to, to this category. For example, if a virtual switch binds to a 1 GbE network adapter, this special category can use at least 100 Mbps. If the value is not a multiple of 8, the value is rounded down to the nearest number that is a multiple of 8. For example, a value input as 1234567 is converted to 1234560.

* `default_flow_minimum_bandwidth_weight` - (Optional) `0` (default). Specifies the minimum bandwidth, in relative weight, that is allocated to a special category called "default flow". Any traffic sent by a virtual network adapter that is connected to this virtual switch and doesn’t have minimum bandwidth allocated is filtered into this category. Specify a value for this parameter only if the minimum bandwidth mode on this virtual switch is weight (See the New-VMSwitch cmdlet). By default, this special category has a weight of 1.

* `default_queue_vmmq_enabled` - (Optional) `false` (default). Specifies if the HyperV host machine will enable virtual machine multi queue for network switch when created.

* `default_queue_vmmq_queue_pairs` - (Optional) `false` (default). Specifies if the HyperV host machine will enable virtual machine multi queue pairs for network switch when created. 

* `default_queue_vrss_enabled` - (Optional) `false` (default). Specifies if the HyperV host machine will enable virtual receive side scaling for network switch when created. 
