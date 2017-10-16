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

* `name` - (Required) The name of the instance.

* `generation` - (Optional) `1`(default). Specifies a note to be associated with the switch to be created.

* `allow_unverified_paths` - (Optional) `false` (default). Specifies if the HyperV cluster will not throw an error if the specified path is not verified by the cluster.

* `automatic_critical_error_action` - (Optional) `0`(default). Specifies the action to take when the VM encounters a critical error, and exceeds the timeout duration specified by the AutomaticCriticalErrorActionTimeout cmdlet. The acceptable values for this parameter are: `Pause` and `None`.

* `automatic_critical_error_action_timeout` - (Optional) `0`(default). Specifies the amount of time, in minutes, to wait in critical pause before powering off the virtual machine.

* `automatic_start_action` - (Optional) `0`(default). Specifies the action the virtual machine is to take upon start. Allowed values are `Nothing`, `StartIfRunning`, and `Start`.

* `automatic_start_delay` - (Optional) `0`(default). Specifies the number of seconds by which the virtual machine's start should be delayed.

* `automatic_stop_action` - (Optional) `0`(default). Specifies the action the virtual machine is to take when the virtual machine host shuts down. Allowed values are `TurnOff`, `Save`, and `ShutDown`.

* `checkpoint_type` - (Optional) `0`(default). Allows you to configure the type of checkpoints created by Hyper-V. Allowed values are `Disabled`, `Standard`, `Production`, and `ProductionOnly`. If `Disabled` is specified, block creation of checkpoints. If `Standard` is specified, create standard checkpoints. If `Production` is specified, create production checkpoints if supported by guest operating system. Otherwise, create standard checkpoints. If `ProductionOnly` is specified, create production checkpoints if supported by guest operating system. Otherwise, the operation fails.

* `dynamic_memory` - (Optional) `false` (default). Specifies if machine instance will have dynamic memory enabled.

* `guest_controlled_cache_types` - (Optional) `false` (default). Specifies if the machine instance will use guest controlled cache types.

* `high_memory_mapped_io_space` - (Optional) `0`(default). 

* `lock_on_disconnect` - (Optional) `false` (default). Specifies whether virtual machine connection in basic mode locks the console after a user disconnects.

* `low_memory_mapped_io_space` - (Optional) `0`(default).

* `memory_maximum_bytes` - (Optional) `0`(default). Specifies the maximum amount of memory that the virtual machine is to be allocated. (Applies only to virtual machines using dynamic memory.)

* `memory_minimum_bytes` - (Optional) `0`(default). Specifies the minimum amount of memory that the virtual machine is to be allocated. (Applies only to virtual machines using dynamic memory.)

* `memory_startup_bytes` - (Optional) `0`(default). Specifies the amount of memory that the virtual machine is to be allocated upon startup. (If the virtual machine does not use dynamic memory, then this is the static amount of memory to be allocated.)

* `notes` - (Optional) `` (Default). Specifies a note to be associated with the machine to be created.

* `processor_count` - (Optional) `1`(default). Specifies the number of virtual processors for the virtual machine.

* `smart_paging_file_path` - (Optional) `` (Default). Specifies the folder in which the Smart Paging file is to be stored.

* `snapshot_file_location` - (Optional) `` (Default). Specifies the folder in which the virtual machine is to store its snapshot files.

* `static_memory` - (Optional) `false` (default). Specifies if the machine instance will use static memory.











