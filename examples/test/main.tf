terraform {
  required_providers {
    hyperv = {
      version = "1.0.3"
      source  = "registry.terraform.io/taliesins/hyperv"
    }
  }
}

provider "hyperv" {

}

data "hyperv_network_switch" "dmz_network_switch" {
  name = "dmz"
}
