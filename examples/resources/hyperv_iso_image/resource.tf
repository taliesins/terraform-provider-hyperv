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

data "archive_file" "bootstrap" {
  type        = "zip"
  source_dir  = "bootstrap"
  output_path = "bootstrap.zip"
}

resource "hyperv_iso_image" "bootstrap" {
  volume_name               = "BOOTSTRAP"
  source_zip_file_path      = data.archive_file.bootstrap.output_path
  source_zip_file_path_hash = data.archive_file.bootstrap.output_sha
  destination_iso_file_path = "$env:TEMP\\bootstrap.iso"
  iso_media_type            = "dvdplusrw_duallayer"
  iso_file_system_type      = "unknown"
}