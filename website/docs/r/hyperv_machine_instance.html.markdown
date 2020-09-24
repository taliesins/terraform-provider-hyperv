---
subcategory: "Machine Instance"
layout: "hyperv"
page_title: "HyperV: hyperv_machine_instance"
description: |-
  Creates and manages a machine instance.
---

# hyperv\_machine\_instance

The ``hyperv_machine_instance`` resource creates and manages an instance on a HyperV environment.

## Example Usage

```hcl
resource "hyperv_machine_instance" "default" {
  name = "WebServer"
  generation = 1
  automatic_critical_error_action = "Pause"
  automatic_critical_error_action_timeout = 30
  automatic_start_action = "StartIfRunning"
  automatic_start_delay = 0
  automatic_stop_action = "Save"
  checkpoint_type = "Production"
  dynamic_memory = false
  guest_controlled_cache_types = false
  high_memory_mapped_io_space = 536870912
  lock_on_disconnect = "Off"
  low_memory_mapped_io_space = 134217728
  memory_maximum_bytes = 1099511627776
  memory_minimum_bytes = 536870912
  memory_startup_bytes = 536870912
  notes = ""
  processor_count = 1
  smart_paging_file_path = "C:\ProgramData\Microsoft\Windows\Hyper-V"
  snapshot_file_location = "C:\ProgramData\Microsoft\Windows\Hyper-V"
  static_memory = true
  state = "Running"

  # Configure integration services
  integration_services {
  }

  # Create a network adaptor
  network_adaptors {
  }

  # Create dvd drive
  dvd_drives {
  }

  # Create a hard disk drive
  hard_disk_drives {
  }
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required). Specifies the name of the new virtual machine.

* `generation` - (Optional) `2` (default). Valid values to use are `1`, `2`. Specifies the generation, as an integer, for the virtual machine.

* `automatic_critical_error_action` - (Optional) `Pause` (default). Valid values to use are `Pause`, `None`. Specifies the action to take when the VM encounters a critical error, and exceeds the timeout duration specified by the AutomaticCriticalErrorActionTimeout cmdlet. 

* `automatic_critical_error_action_timeout` - (Optional) `30` (default). Specifies the amount of time, in minutes, to wait in critical pause before powering off the virtual machine.

* `automatic_start_action` - (Optional) `StartIfRunning` (default). Valid values to use are `Nothing`, `StartIfRunning`, `Start`. Specifies the action the virtual machine is to take upon start. 

* `automatic_start_delay` - (Optional) `0` (default). Specifies the number of seconds by which the virtual machine's start should be delayed.

* `automatic_stop_action` - (Optional) `Save` (default). Valid values to use are `TurnOff`, `Save`, `ShutDown`. Specifies the action the virtual machine is to take when the virtual machine host shuts down. 

* `checkpoint_type` - (Optional) `Production` (default). Valid values to use are `Disabled`, `Standard`, `Production`, `ProductionOnly`. Allows you to configure the type of checkpoints created by Hyper-V. If `Disabled` is specified, block creation of checkpoints. If `Standard` is specified, create standard checkpoints. If `Production` is specified, create production checkpoints if supported by guest operating system. Otherwise, create standard checkpoints. If `ProductionOnly` is specified, create production checkpoints if supported by guest operating system. Otherwise, the operation fails.

* `dynamic_memory` - (Optional) `false` (default). Specifies if machine instance will have dynamic memory enabled.

* `guest_controlled_cache_types` - (Optional) `false` (default). Specifies if the machine instance will use guest controlled cache types.

* `high_memory_mapped_io_space` - (Optional) `536870912` (default). 

* `lock_on_disconnect` - (Optional) `Off` (default). Valid values to use are `On`, `Off`. Specifies whether virtual machine connection in basic mode locks the console after a user disconnects.

* `low_memory_mapped_io_space` - (Optional) `134217728` (default).

* `memory_maximum_bytes` - (Optional) `1099511627776` (default). Specifies the maximum amount of memory that the virtual machine is to be allocated. (Applies only to virtual machines using dynamic memory.)

* `memory_minimum_bytes` - (Optional) `536870912` (default). Specifies the minimum amount of memory that the virtual machine is to be allocated. (Applies only to virtual machines using dynamic memory.)

* `memory_startup_bytes` - (Optional) `536870912` (default). Specifies the amount of memory that the virtual machine is to be allocated upon startup. (If the virtual machine does not use dynamic memory, then this is the static amount of memory to be allocated.)

* `notes` - (Optional) empty (Default). Specifies a note to be associated with the machine to be created.

* `processor_count` - (Optional) `1` (default). Specifies the number of virtual processors for the virtual machine.

* `smart_paging_file_path` - (Optional) `C:\ProgramData\Microsoft\Windows\Hyper-V` (Default). Specifies the folder in which the Smart Paging file is to be stored.

* `snapshot_file_location` - (Optional) `C:\ProgramData\Microsoft\Windows\Hyper-V` (Default). Specifies the folder in which the virtual machine is to store its snapshot files.

* `static_memory` - (Optional) `true` (default). Specifies if the machine instance will use static memory.

* `state` - (Optional) `Running` (default).  Valid values to use are `Running`, `Off`. Specifies if the machine instance will be running or off.

* `wait_for_state_timeout` - (Optional) `120` (default). The amount of time in seconds to wait before throwing an exception when trying to change for the virtual machine to the desired state.

* `wait_for_state_poll_period` - (Optional) `2` (default). The amount of time in seconds to wait between trying to change for the virtual machine to the desired state.

* `wait_for_ips_timeout` - (Optional) `300` (default). The amount of time in seconds to wait before throwing an exception when trying to get ip addresses for network cards on the virtual machine.

* `wait_for_ips_poll_period` - (Optional) `5` (default). The amount of time in seconds to wait between trying to get ip addresses for network cards on the virtual machine.

* `network_adaptors` - (Optional) empty array (default). An array of all the network adaptors connected to vm.

* `integration_services` - (optional) default integration services (default). A map of all the integration services and if the integration service should be enabled/disabled. Integration services that are not specified will not be enforced.

* `vm_processor` - (optional) default vm processor (default). All the vm processor settings connected to vm.

* `dvd_drives` - (Optional) empty array (default). An array of all the dvd drives connected to vm.

* `hard_disk_drives` - (Optional) empty array (default). An array of all the hard disk drives connected to vm.

### VM Firmware

Note: this terraform resource will be skipped if the VM generation is 1. Terraform Schema does not provide a way to valid schema against other properties hence this approach.

```hcl
resource "hyperv_machine_instance" "default" {
  name = "WebServer"
  vm_firmware {
    enable_secure_boot = ""
    secure_boot_template = ""
    secure_boot_template_id = ""
    preferred_network_boot_protocol = ""
    console_mode = ""
    pause_after_boot_failure = ""
  }
}
```

* `enable_secure_boot` - (Optional) `On` (default). Valid values to use are `On`, `Off`. Specifies whether to enable secure boot.
* `secure_boot_template` - (Optional) `MicrosoftWindows` (default). Example values to use are `MicrosoftWindows`,`MicrosoftUEFICertificateAuthority`, `OpenSourceShieldedVM`. Specifies the name of the secure boot template. If secure boot is enabled, you must have a valid secure boot template for the guest operating system to start.
* `preferred_network_boot_protocol` - (Optional) `IPv4` (default). Valid values to use are `IPv4`, `IPv6`. Specifies the IP protocol version to use during a network boot.
* `console_mode` - (Optional) `Default` (default). Valid values to use are `Default`, `COM1`, `COM2`, `None`. Specifies the console mode type for the virtual machine. This parameter allows a virtual machine to run without graphical user interface.
* `pause_after_boot_failure` - (Optional) `Off` (default). Valid values to use are `On`, `Off`. Specifies the behavior of the virtual machine after a start failure. For a value of On, if the virtual machine fails to start correctly from a device, the virtual machine is paused.

### VM Processor

```hcl
resource "hyperv_machine_instance" "default" {
  name = "WebServer"
  vm_processor {
    compatibility_for_migration_enabled = false
    compatibility_for_older_operating_systems_enabled = false
    hw_thread_count_per_core = 0
    maximum = 100
    reserve = 0
    relative_weight = 100
    maximum_count_per_numa_node = 0
    maximum_count_per_numa_socket = 0
    enable_host_resource_protection = false
    expose_virtualization_extensions = false
  }
}
```

* `compatibility_for_migration_enabled` - (Optional) `false` (default). Specifies whether the virtual processor's features are to be limited for compatibility when migrating the virtual machine to another host.
* `compatibility_for_older_operating_systems_enabled` - (Optional) `false` (default). Specifies whether the virtual processor's features are to be limited for compatibility with older operating systems.
* `hw_thread_count_per_core` - (Optional) `0` (default). Specifies the number of virtual SMT threads exposed to the virtual machine. Setting this value to 0 indicates the virtual machine will inherit the host's number of threads per core. This setting may not exceed the host's number of threads per core.

Note: Windows Server 2016 does not support setting HwThreadCountPerCore to 0. For more details, see Configuring VM SMT settings using PowerShell.
* `maximum` - (Optional) `100` (default). Specifies the maximum percentage of resources available to the virtual machine processor to be configured. Allowed values range from 0 to 100.
* `reserve` - (Optional) `0` (default). Specifies the percentage of processor resources to be reserved for this virtual machine. Allowed values range from 0 to 100.
* `relative_weight` - (Optional) `100` (default). Specifies the priority for allocating the physical computer's processing power to this virtual machine relative to others. Allowed values range from 1 to 10000.
* `maximum_count_per_numa_node` - (Optional) `0` (default). Specifies the maximum number of processors per NUMA node to be configured for the virtual machine.
* `maximum_count_per_numa_socket` - (Optional) `0` (default). Specifies the maximum number of sockets per NUMA node to be configured for the virtual machine.
* `enable_host_resource_protection` - (Optional) `false` (default). Specifies whether to enable host resource protection on the virtual machine. When enabled, the host will enforce limits on some aspects of the virtual machine's activity, preventing excessive consumption of host compute resources. VM activities controlled by this setting include the VMbus pipe messages associated with a subset of the VM's virtual devices, and intercepts generated by the VM. The virtual devices affected include the video, keyboard, mouse, and dynamic memory VDEVs.
* `expose_virtualization_extensions` - (Optional) `false` (default). Specifies whether the hypervisor should expose the presence of virtualization extensions to the virtual machine, which enables support for nested virtualization.

### Integration Service

```hcl
resource "hyperv_machine_instance" "default" {
  name = "WebServer"
  vm_proce {
    "Guest Service Interface" = false
    "Heartbeat"               = true
    "Key-Value Pair Exchange" = true
    "Shutdown"                = true
    "Time Synchronization"    = true
    "VSS"                     = true
  }
}
```

* `Guest Service Interface` - (Optional) `false` (default). Provides an interface for the Hyper-V host to copy files to or from the virtual machine.

* `Heartbeat` - (Optional) `true` (default). Reports that the virtual machine is running correctly.

* `Key-Value Pair Exchange` - (Optional) `true` (default). Provides a way to exchange basic metadata b etween the virtual machine and the host.

* `Shutdown` - (Optional) `true` (default). Allows the host to trigger virtual machines shutdown.

* `Time Synchronization` - (Optional) `true` (default). Synchronizes the virtual machine's clock with the host computer's clock.

* `VSS` - (Optional) `true` (default). Allows Volume Shadow Copy Service to back up the virtual machine with out shutting it down.

### Network Adaptors

```hcl
resource "hyperv_machine_instance" "default" {
  name = "WebServer"
  network_adaptors {
    name = "wan"
    switch_name = "${hyperv_network_switch.dmz_network_switch.name}"
    management_os = false
    is_legacy = false
    dynamic_mac_address = true
    static_mac_address = ""
    mac_address_spoofing = "Off"
    dhcp_guard = "Off"
    router_guard = "Off"
    port_mirroring = "None"
    ieee_priority_tag = "Off"
    vmq_weight=100
    iov_queue_pairs_requested=1
    iov_interrupt_moderation="Off"
    iov_weight=100
    ipsec_offload_maximum_security_association=512
    maximum_bandwidth=0
    minimum_bandwidth_absolute=0
    minimum_bandwidth_weigh=0
    mandatory_feature_id=[]
    resource_pool_name=""
    test_replica_pool_name=""
    test_replica_switch_name=""
    virtual_subnet_id=0
    allow_teaming="On"
    not_monitored_in_cluster=false
    storm_limit=0
    dynamic_ip_address_limit=0
    device_naming="Off"
    fix_speed_10g="Off"
    packet_direct_num_procs=0
    packet_direct_moderation_count=0
    packet_direct_moderation_interval=0
    vrss_enabled=true
    vmmq_enabled=false
    vmmq_queue_pairs=16
  }
}
```

* `name` - (Required). Specifies the name for the virtual network adapter.

* `switch_name` - (Optional) empty (default). Specifies the name of the virtual switch to connect to the new network adapter. If the switch name is not unique, then the operation fails.

* `management_os` - (Optional) `false` (default). Specifies the virtual network adapter in the management operating system to be configured.

* `is_legacy` - (Optional) `false` (default). Specifies whether the virtual network adapter is the legacy type.

* `dynamic_mac_address` - (Optional) `true` (default). Assigns a dynamically generated MAC address to the virtual network adapter.

* `static_mac_address` - (Optional) empty (default). Assigns a specific a MAC addresss to the virtual network adapter.

* `mac_address_spoofing` - (Optional) `Off` (default). Valid values to use are `On`, `Off`. Specifies whether virtual machines may change the source MAC address in outgoing packets to one not assigned to them. On allows the virtual machine to use a different MAC address. Off only allows the virtual machine to use the MAC address assigned to it. 

* `dhcp_guard` - (Optional) `Off` (default). Valid values to use are `On`, `Off`. Specifies whether to drop DHCP messages from a virtual machine claiming to be a DHCP server. 

* `router_guard` - (Optional) `Off` (default). Valid values to use are `On`, `Off`. Specifies whether to drop Router Advertisement and Redirection messages from unauthorized virtual machines. If On is specified, such messages are dropped. If Off is specified, such messages are sent.

* `port_mirroring` - (Optional) `None` (default). Valid values to use are `None`, `Source`, `Destination`. Specifies the port mirroring mode for the network adapter to be configured. If a virtual network adapter is configured as Source, every packet it sends or receives is copied and forwarded to a virtual network adapter configured to receive the packets. If a virtual network adapter is configured as Destination, it receives copied packets from the source virtual network adapter. The source and destination virtual network adapters must be connected to the same virtual switch. Specify None to disable the feature.

* `ieee_priority_tag` - (Optional) `Off` (default). Valid values to use are `On`, `Off`. Specifies whether IEEE 802.1p tagged packets from the virtual machine should be trusted. If it is on, the IEEE 802.1p tagged packets will be let go as is. If it is off, the priority value is reset to 0.

* `vmq_weight` - (Optional) `100` (default). Valid values to use are between `1` to `100`. Specifies whether virtual machine queue (VMQ) is to be enabled on the virtual network adapter. The relative weight describes the affinity of the virtual network adapter to use VMQ. Specify 0 to disable VMQ on the virtual network adapter.

* `iov_queue_pairs_requested` - (Optional) `1` (default). Valid values to use are between `1` to `4294967295`. Specifies the number of hardware queue pairs to be allocated to an SR-IOV virtual function. If receive-side scaling (RSS) is required, and if the physical network adapter that binds to the virtual switch supports RSS on SR-IOV virtual functions, then more than one queue pair is required.

* `iov_interrupt_moderation` - (Optional) `Off` (default). Valid values to use are `Default`, `Adaptive`, `Off`, `Low `, `Medium`, `High`. Specifies the interrupt moderation value for a single-root I/O virtualization (SR-IOV) virtual function assigned to a virtual network adapter. If Default is chosen, the value is determined by the physical network adapter vendor's setting. If Adaptive is chosen, the interrupt moderation rate will be based on the runtime traffic pattern.

* `iov_weight` - (Optional) `100` (default). Valid values to use are between `0` to `100`. Specifies whether single-root I/O virtualization (SR-IOV) is to be enabled on this virtual network adapter. The relative weight sets the affinity of the virtual network adapter to the assigned SR-IOV virtual function. Specify 0 to disable SR-IOV on the virtual network adapter.

* `ipsec_offload_maximum_security_association` - (Optional) `512` (default). Specifies the maximum number of security associations that can be offloaded to the physical network adapter that is bound to the virtual switch and that supports IPSec Task Offload. Specify zero to disable the feature.

* `maximum_bandwidth` - (Optional) `0` (default). Specifies the maximum bandwidth, in bits per second, for the virtual network adapter. The specified value is rounded to the nearest multiple of eight. Specify zero to disable the feature.

* `minimum_bandwidth_absolute` - (Optional) `0` (default). Specifies the minimum bandwidth, in bits per second, for the virtual network adapter. The specified value is rounded to the nearest multiple of eight. A value larger than 100 Mbps is recommended.

* `minimum_bandwidth_weight` - (Optional) `0` (default). Valid values to use are between `0` to `100`. Specifies the minimum bandwidth, in terms of relative weight, for the virtual network adapter. The weight describes how much bandwidth to provide to the virtual network adapter relative to other virtual network adapters connected to the same virtual switch. Specify 0 to disable the feature.

* `mandatory_feature_id` - (Optional) array of strings (default). Specifies the unique identifiers of the virtual switch extension features that are required for this virtual network adapter to operate.

* `resource_pool_name` - (Optional) empty (default). Specifies the name of the resource pool.

* `test_replica_pool_name` - (Optional) empty (default). This parameter applies only to virtual machines that are enabled for replication. It specifies the name of the network resource pool that will be used by this virtual network adapter when its virtual machine is created during a test failover.

* `test_replica_switch_name` - (Optional) empty (default). This parameter applies only to virtual machines that are enabled for replication. It specifies the name of the virtual switch to which the virtual network adapter should be connected when its virtual machine is created during a test failover.

* `virtual_subnet_id` - (Optional) `0` (default). Valid values to use are `0` or between `4096` to `16777215` (2^24 - 1). Specifies the virtual subnet ID to use with Hyper-V Network Virtualization. Use 0 to clear this parameter.

* `allow_teaming` - (Optional) `On` (default). Valid values to use are `On`, `Off`. Specifies whether the virtual network adapter can be teamed with other network adapters connected to the same virtual switch. 

* `not_monitored_in_cluster` - (Optional) `false` (default). Indicates whether to not monitor the network adapter if the virtual machine that it belongs to is part of a cluster. By default, network adapters for clustered virtual machines are monitored.

* `storm_limit` - (Optional) `0` (default). Specifies the number of broadcast, multicast, and unknown unicast packets per second a virtual machine is allowed to send through the specified virtual network adapter. Broadcast, multicast, and unknown unicast packets beyond the limit during that one second interval are dropped. A value of zero (0) means there is no limit.

* `dynamic_ip_address_limit` - (Optional) `0` (default). Specifies the dynamic IP address limit.

* `device_naming` - (Optional) `Off` (default). Valid values to use are `On`, `Off`. Specifies whether this adapter uses device naming.

* `fix_speed_10g` - (Optional) `Off` (default). Valid values to use are `On`, `Off`. Specifies whether the adapter uses fix speed of 10G.

* `packet_direct_num_procs` - (Optional) `0` (default). Specifies the number of processors to use for virtual switch processing inside of the host.

* `packet_direct_moderation_count` - (Optional) `0` (default). Specifies the number of packets to wait for before signaling an interrupt.

* `packet_direct_moderation_interval` - (Optional) `0` (default). Specifies the amount of time, in milliseconds, to wait before signaling an interrupt after a packet arrives.

* `vrss_enabled` - (Optional) `true` (default). Should Virtual Receive Side Scaling be enabled. This configuration allows the load from a virtual network adapter to be distributed across multiple virtual processors in a virtual machine (VM), allowing the VM to process more network traffic more rapidly than it can with a single logical processor.

* `vmmq_enabled` - (Optional) `false` (default). Should Virtual Machine Multi-Queue be enabled. With set to true multiple queues are allocated to a single VM with each queue affinitized to a core in the VM.

* `vmmq_queue_pairs` - (Optional) `16` (default). The number of Virtual Machine Multi-Queues to create for this VM.

* `wait_for_ips` - (Optional) `true` (default). Wait for the network card to be assigned an ip address. 

* `ip_addresses` - (Computed).  The current list of IP addresses on this machine. If HyperV integration tools is not running on the virtual machine, or if the VM is powered off, or has not been assigned an ip address, this list will be empty. 

### Dvd drives

```hcl
resource "hyperv_machine_instance" "default" {
  name = "WebServer"
  dvd_drives {
    controller_number = "0"
    controller_location = "1"
    path = "c:/iso/windows-server-2016.iso"
    resource_pool_name = ""
  }
}
```

* `controller_number` - (Required).  Specifies the number of the controller to which the DVD drive is to be added. 

* `controller_location` - (Required). Specifies the number of the location on the controller at which the DVD drive is to be added. 

* `path` - (Optional) empty (default). Specifies the full path to the virtual hard disk file or physical hard disk volume for the added DVD drive.

* `resource_pool_name` - (Optional) empty (default). Specifies the friendly name of the ISO resource pool to which this DVD drive is to be associated.

### Hard Disk Drives

```hcl
resource "hyperv_machine_instance" "default" {
  name = "WebServer"
  hard_disk_drives {
    controller_type = "Ide"
    controller_number = "0"
    controller_location = "0"
    path = "c:/virtual machines/WebServer/windows-server-2016.vhd"
    disk_number = 4294967295
    resource_pool_name = "Primordial"
    support_persistent_reservations = false
    maximum_iops = 0
    minimum_iops = 0
    qos_policy_id = "00000000-0000-0000-0000-000000000000"
    override_cache_attributes = "Default"
  }
}
```
* `controller_type` - (Optional) `Scsi` (default). Valid values to use are `Ide`, `Scsi`. Specifies the type of the controller to which the hard disk drive is to be added. 

* `controller_number` - (Required).  Specifies the number of the controller to which the hard disk drive is to be added. 

* `controller_location` - (Required).  Specifies the number of the location on the controller at which the hard disk drive is to be added. 

* `path` - (Optional) empty (default). Specifies the full path of the hard disk drive file to be added.

* `disk_number` - (Optional) `4294967295` (default). If value is 4294967295 then disk number is ignored. Specifies the disk number of the offline physical hard drive to be connected as a passthrough disk.

* `resource_pool_name` - (Optional) `Primordial` (default). Specifies the friendly name of the resource pool to which this virtual hard disk is to be associated.

* `support_persistent_reservations` - (Optional) `false` (default). Indicates that the hard disk supports SCSI persistent reservation semantics. Specify this parameter when the hard disk is a shared disk that is used by multiple virtual machines.

* `maximum_iops` - (Optional) `0` (default). If value is 0 then iops is ignored. Specifies the maximum normalized I/O operations per second (IOPS) for the hard disk. Hyper-V calculates normalized IOPS as the total size of I/O per second divided by 8 KB.

* `minimum_iops` - (Optional) `0` (default). If maximum iops value is 0 then iops is ignored. Specifies the minimum normalized I/O operations per second (IOPS) for the hard disk. Hyper-V calculates normalized IOPS as the total size of I/O per second divided by 8 KB.

* `qos_policy_id` - (Optional) `00000000-0000-0000-0000-000000000000` (default). Specifies the unique ID for a storage QoS policy that this cmdlet associates with the hard disk drive. If value is 00000000-0000-0000-0000-000000000000 then qos policy id is ignored.

* `override_cache_attributes` - (Optional) `Default` (default). Valid values to use are `Default`, `WriteCacheEnabled`, `WriteCacheAndFUAEnabled`, `WriteCacheDisabled`.

  With Default it is equivalent of WriteCacheDisabled.
  
  With WriteCacheEnabled write I/O is acknowledged as written before it is committed to stable media. If your internal disks, DAS, SAN, or NAS has a battery backup system that can guarantee clean cache flushes on a power outage, write caching is generally safe. Internal batteries that report their status and/or automatically disable caching are best. UPS-backed systems are sometimes OK, but they are not foolproof.

  With WriteCacheAndFUAEnabled write I/O is committed to stable media BEFORE the I/O is acknowledged as written.

  With WriteCacheDisabled when I/O is written it is acknowledged as written as there is no cache in between.