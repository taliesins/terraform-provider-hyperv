## 0.1.0 (Unreleased)

FEATURES:

* **New Resource:** `hyperv_machine` 
* **New Resource:** `hyperv_network_switch`

NOTES:

* Rewritten Powershell remote execution and elevated remote execution
* Changed Winrmcp to use Powershell commands directly rather then use base64 encoded strings as we want to prevent Powershell progress leaking.
* Changed Winrmcp to return path of files on remote box as the location of $env:temp can change in Powershell depending on the session instance.
