provider "hyperv" {
}

resource "hyperv_network_switch" "dmz_network_switch" {
  name = "dmz"
}

resource "hyperv_machine_instance" "web_Server" {
  name = "web_server"

  network_adaptors {
      name = "wan"
      switch_name = hyperv_network_switch.dmz_network_switch.name
  }
}