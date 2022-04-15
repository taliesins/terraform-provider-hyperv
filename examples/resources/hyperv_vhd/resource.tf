resource "hyperv_vhd" "web_server_vhd" {
  path                 = "C:\\data\\VirtualMachines\\web_server\\Virtual Hard Disks\\MobyLinuxVM.vhdx"
  source               = ""
  source_vm            = ""
  source_disk          = 0
  vhd_type             = "Dynamic"
  parent_path          = ""
  size                 = 21474836480
  block_size           = 0
  logical_sector_size  = 0
  physical_sector_size = 0
}