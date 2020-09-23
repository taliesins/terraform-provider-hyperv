---
subcategory: "VHD"
layout: "hyperv"
page_title: "HyperV: hyperv_vhd"
description: |-
  Creates and manages vhd/vhdx/vhds.
---

# hyperv\_vhd

The ``hyperv_vhd`` resource creates and manages vhd/vhdx/vhds on a HyperV environment.

## Example Usage

```hcl
resource "hyperv_vhd" "web_server_vhd" {
  path = "C:\\data\\VirtualMachines\\web_server\\Virtual Hard Disks\\MobyLinuxVM.vhdx"
  source = ""
  source_vm = ""
  source_disk = 0
  vhd_type = "Dynamic"
  parent_path = ""
  size = 21474836480
  block_size = 0
  logical_sector_size = 0
  physical_sector_size = 0
}
```

## Argument Reference

The following arguments are supported:

* `path` - (Required). Path to the new virtual hard disk file(s) that is being created or being copied to. If a filename or relative path is specified, the new virtual hard disk path is calculated relative to the current working directory. Depending on the source selected, the path will be used to determine where to copy source vhd/vhdx/vhds file to.

* `source` - (Optional) empty (default). This field is mutually exclusive with the fields "source_vm", "parent_path", "source_disk". This value can be a url or a path (including wildcards). Box, Zip and 7z files will automatically be expanded. The destination folder will be the directory portion of the path. If expanded files have a folder called "Virtual Machines", then the "Virtual Machines" folder will be used instead of the entire archive contents. 

* `source_vm` - (Optional) empty (default). This field is mutually exclusive with the fields "source", "parent_path", "source_disk". This value is the name of the vm to copy the vhds from.

* `source_disk` - (Optional) `0` (default). This field is mutually exclusive with the fields "source", "source_vm", "parent_path". Specifies the physical disk to be used as the source for the virtual hard disk to be created.

* `vhd_type` - (Optional) `Dynamic` (default). Valid values to use are `Unknown`, `Fixed`, `Dynamic`, `Differencing`. This field is mutually exclusive with the fields "source", "source_vm", "parent_path".

* `parent_path` - (Optional) empty (default). This field is mutually exclusive with the fields "source", "source_vm", "source_disk", "size". Specifies the path to the parent of the differencing disk to be created (this parameter may be specified only for the creation of a differencing disk).

* `size` - (Optional) `0` (default). This field is mutually exclusive with the field "parent_path". The maximum size, in bytes, of the virtual hard disk to be created.

* `block_size` - (Optional) `0` (default). This field is mutually exclusive with the fields "source", "source_vm", "parent_path". Specifies the block size, in bytes, of the virtual hard disk to be created.

* `logical_sector_size` - (Optional) `0` (default). Valid values to use are `0`, `512`, `4096`. This field is mutually exclusive with the fields "source", "source_vm", "parent_path". Specifies the logical sector size, in bytes, of the virtual hard disk to be created. 

* `physical_sector_size` - (Optional) `0` (default). Valid values to use are `0`, `512`, `4096`. This field is mutually exclusive with the fields	"source",	"source_vm", "parent_path". Specifies the physical sector size, in bytes.
