terraform {
  required_providers {
    hyperv = {
      version = "1.0.3"
      source = "registry.terraform.io/taliesins/hyperv"
    }
  }
}

provider "hyperv" {

}

data  "hyperv_network_switch" "dmz_network_switch" {
  name = "dmz"
}

data  "hyperv_machine_instance" "web_server_g1" {
  name = "web_server_g1"
}

resource "hyperv_vhd" "web_server_g3_vhd" {
  path = "c:\\vhdx\\web_server_g3.vhdx"
  source_vm = data.hyperv_machine_instance.web_server_g1.name
}

resource "hyperv_machine_instance" "web_server_g3" {
  name = "web_server_g3"
  static_memory = true
  
  network_adaptors {
    name = "wan"
    switch_name = data.hyperv_network_switch.dmz_network_switch.name
  }

  hard_disk_drives {
    path = hyperv_vhd.web_server_g3_vhd.path
    controller_number = "0"
    controller_location = "0"
  }

  integration_services = {
    VSS = true
  }
}