terraform {
  required_providers {
    hyperv = {
      source  = "taliesins/hyperv"
      version = ">= 1.0.3"
    }
  }
}

provider "hyperv" {
}

resource "hyperv_vhd" "web_server_vhd" {
  path = "c:\\web_server\\web_server_g2.vhdx"
  #source               = ""
  #source_vm            = ""
  #source_disk          = 0
  vhd_type = "Dynamic"
  #parent_path          = ""
  size = 10737418240 #10GB
  #block_size           = 0
  #logical_sector_size  = 0
  #physical_sector_size = 0
}

data "hyperv_vhd" "web_server_vhd" {
  path = hyperv_vhd.web_server_vhd.path
}

output "hyperv_vhd" {
  value = data.hyperv_vhd.web_server_vhd
}