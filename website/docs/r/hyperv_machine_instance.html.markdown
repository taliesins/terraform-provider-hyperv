---
layout: "hyperv"
page_title: "HyperV: hyperv_machine_instance"
sidebar_current: "docs-hyperv-resource-machine-instance"
description: |-
  Creates and manages a machine instance.
---

# hyperv\_machine\_instance

The ``hyperv_machine_instance`` resource creates and manages an instance on a HyperV environment.

## Example Usage

```hcl
resource "hyperv_machine_instance" "default" {
  name = "WebServer"
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required). Specifies the name of the new virtual machine.

* `generation` - (Optional) `1` (default). Valid values to use are `1`, `2`. Specifies the generation, as an integer, for the virtual machine.

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

* `integration_services` - (optional) default integration services (default). A map of all the integration services and if the integration service should be enabled/disabled. Integration services that are not specified will not be enforced.

```hcl
resource "hyperv_machine_instance" "default" {
  name = "WebServer"
  integration_services {
    "Guest Service Interface" = false
    "Heartbeat"               = true
    "Key-Value Pair Exchange" = true
    "Shutdown"                = true
    "Time Synchronization"    = true
    "VSS"                     = true
  }
}
```

* `network_adaptors` - (Optional) empty array (default). An array of all the network adaptors connected to vm.

```hcl
resource "hyperv_machine_instance" "default" {
  name = "WebServer"
  network_adaptors {
    name = "wan"
    switch_name = "${hyperv_network_switch.dmz_network_switch.name}"
  }
}

* `dvd_drives` - (Optional) empty array (default). An array of all the dvd drives connected to vm.

```hcl
resource "hyperv_machine_instance" "default" {
  name = "WebServer"
  dvd_drives {
    path = "${hyperv_vhd.web_server_vhd.path}"
    controller_number = "0"
    controller_location = "0"
  }
}

* `hard_disk_drives` - (Optional) empty array (default). An array of all the hard disk drives connected to vm.

```hcl
resource "hyperv_machine_instance" "default" {
  name = "WebServer"
  hard_disk_drives {
    path = "${hyperv_vhd.web_server_vhd.path}"
    controller_number = "0"
    controller_location = "0"
  }
}

### Integration Service
### Network Adaptors
### Dvd drives
### Hard Disk Drives