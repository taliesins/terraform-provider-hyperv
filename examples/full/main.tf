provider "hyperv" {
}

resource "hyperv_vhd" "web_server_vhd" {
  path = "C:\\data\\VirtualMachines\\web_server\\Virtual Hard Disks\\MobyLinuxVM.vhdx"
  source_vm = "MobyLinuxVM"
}

resource "hyperv_network_switch" "dmz_network_switch" {
  name = "dmz"
}

resource "hyperv_machine_instance" "web_server" {
  name = "web_server"
  
  integration_services {
    "VSS" = true
  }

  network_adaptors {
    name = "wan"
    switch_name = "${hyperv_network_switch.dmz_network_switch.name}"
  }

  hard_disk_drives {
    path = "${hyperv_vhd.web_server_vhd.path}"
    controller_number = "0"
    controller_location = "0"
  }
}