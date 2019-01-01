## 0.5.0.23 (Beta Release)

FEATURES:
* **Resource:** `hyperv_machine_instance`
    * **Sub Resource:** `network_adaptors`
        * exposes list of ip addresses for each network adaptor
        * can specify to wait for network adaptor to get ip addresses

NOTES:

- Can now specify timeouts and poll periods for waiting for vm to be in state and for waiting for ip addresses on network adaptor

## 0.5.0.18 (Beta Release)

FEATURES:

* **New Resource:** `hyperv_network_switch`
* **New Resource:** `hyperv_vhd` 
* **New Resource:** `hyperv_machine_instance`
    * **New Sub Resource:** `network_adaptors`
    * **New Sub Resource:** `dvd_drives`
    * **New Sub Resource:** `hard_disk_drives`

NOTES:

- Remote scheduled task powershell runner does not run into issues with escaping variables or escaping between the different scripting layers.
- Changed Winrmcp to use Powershell commands directly rather then use base64 encoded strings as we want to prevent Powershell progress leaking.
- Changed Winrmcp to return path of files on remote box as the location of $env:temp can change in Powershell depending on the session instance.
- Runs all HyperV commands remotely i.e. so the provider can run on a linux machine and connect remotely to a windows machine running HyperV.
- Almost all functionality of Powershell HyperV commandlets for the resources is exposed via Terraform resources.
- Support for downloading zip or 7z format for VHDs
- Support for downloading Packer box format for VHDs
