HyperV Terraform Provider
=========================

#This is beta code. Working on adding acceptance tests so that I can mark it as release quality.

- [Website](https://github.com/taliesins/terraform-provider-hyperv)
- [Releases](https://github.com/taliesins/terraform-provider-hyperv/releases)
- [Documentation](https://registry.terraform.io/providers/taliesins/hyperv/latest/docs)
- [Issues](https://github.com/taliesins/terraform-provider-hyperv/issues)

![Hashi Logo](https://cdn.rawgit.com/taliesins/terraform-provider-hyperv/master/website/logo-hashicorp.svg "Hashi Logo")
![Windows Server Logo](https://cdn.rawgit.com/taliesins/terraform-provider-hyperv/master/website/windows-server-2016-logo.svg "Windows Server Logo")

Features
------------
- Remote scheduled task powershell runner does not run into issues with escaping variables or escaping between the different scripting layers.
- Changed Winrmcp to use Powershell commands directly rather than using base64 encoded strings as we want to prevent Powershell progress leaking.
- Changed Winrmcp to return path of files on remote box as the location of $env:temp can change in Powershell depending on the session instance.
- Runs all HyperV commands remotely i.e. so the provider can run on a linux machine and connect remotely to a windows machine running HyperV.
- Almost all functionality of Powershell HyperV commandlets for the resources is exposed via Terraform resources.
- Resource - Network Switch
- Resource - VHD
- Resource - Virtual Machine Instance
  - Network adaptors
  - Hard drives
  - Dvd drives

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 1.29.x
-	[Go](https://golang.org/doc/install) 1.19 (to build the provider plugin)
-   Connectivity and credentials to a server running HyperV on Windows 10, Windows Server 2016 or newer

Setting up server for Provider usage
------------------------------------

- Install HyperV `Enable-WindowsOptionalFeature -Online -FeatureName:Microsoft-Hyper-V -All`
- Enable WinRM with negotiate authentication support
```
Enable-PSRemoting -SkipNetworkProfileCheck -Force

Set-WSManInstance WinRM/Config/WinRS -ValueSet @{MaxMemoryPerShellMB = 1024}
Set-WSManInstance WinRM/Config -ValueSet @{MaxTimeoutms=1800000}
Set-WSManInstance WinRM/Config/Client -ValueSet @{TrustedHosts="*"}
Set-WSManInstance WinRM/Config/Service/Auth -ValueSet @{Negotiate = $true}
```
- WinRM allow HTTPS
```
#Create CA certificate
$rootCaName = "DevRootCA"
$rootCaPassword = ConvertTo-SecureString "P@ssw0rd" -asplaintext -force 
$rootCaCertificate = Get-ChildItem cert:\LocalMachine\Root |?{$_.subject -eq "CN=$rootCaName"}
if (!$rootCaCertificate){
  Get-ChildItem cert:\LocalMachine\My |?{$_.subject -eq "CN=$rootCaName"} | remove-item -force
  if (Test-Path .\$rootCaName.cer) {
    remove-item .\$rootCaName.cer -force
  }
  if (Test-Path .\$rootCaName.pfx) {
    remove-item .\$rootCaName.pfx -force
  }
  $params = @{
    Type = 'Custom'
    DnsName = $rootCaName
    Subject = "CN=$rootCaName"
    KeyExportPolicy = 'Exportable'
    CertStoreLocation = 'Cert:\LocalMachine\My'
    KeyUsageProperty = 'All'
    KeyUsage = 'None'
    Provider = 'Microsoft Strong Cryptographic Provider'
    KeySpec = 'KeyExchange'
    KeyLength = 4096
    HashAlgorithm = 'SHA256'
    KeyAlgorithm = 'RSA'
    NotAfter = (Get-Date).AddYears(5)
  }
  $rootCaCertificate = New-SelfSignedCertificate @params

  Export-Certificate -Cert $rootCaCertificate -FilePath .\$rootCaName.cer -Verbose
  Export-PfxCertificate -Cert $rootCaCertificate -FilePath .\$rootCaName.pfx -Password $rootCaPassword -Verbose
  Get-ChildItem cert:\LocalMachine\My |?{$_.subject -eq "CN=$rootCaName"} | remove-item -force
  Import-PfxCertificate -FilePath .\$rootCaName.pfx -CertStoreLocation Cert:\LocalMachine\Root -password $rootCaPassword -Exportable -Verbose
  Import-PfxCertificate -FilePath .\$rootCaName.pfx -CertStoreLocation Cert:\LocalMachine\My -password $rootCaPassword -Exportable -Verbose
  $rootCaCertificate = Get-ChildItem cert:\LocalMachine\My |?{$_.subject -eq "CN=$rootCaName"}
}

#Create host certificate using CA
$hostName = [System.Net.Dns]::GetHostName()
$hostPassword = ConvertTo-SecureString "P@ssw0rd" -asplaintext -force
$hostCertificate = Get-ChildItem cert:\LocalMachine\My |?{$_.subject -eq "CN=$hostName"}
if (!$hostCertificate){
  if (Test-Path .\$hostName.cer) {
    remove-item .\$hostName.cer -force
  }
  if (Test-Path .\$hostName.pfx) {
    remove-item .\$hostName.pfx -force
  }
  $dnsNames = @($hostName, "localhost", "127.0.0.1") + [System.Net.Dns]::GetHostByName($env:computerName).AddressList.IpAddressToString
  
  $params = @{
    Type = 'Custom'
    DnsName = $dnsNames
    Subject = "CN=$hostName"
    KeyExportPolicy = 'Exportable'
    CertStoreLocation = 'Cert:\LocalMachine\My'
    KeyUsageProperty = 'All'
    KeyUsage = @('KeyEncipherment','DigitalSignature','NonRepudiation')
    TextExtension = @("2.5.29.37={text}1.3.6.1.5.5.7.3.1,1.3.6.1.5.5.7.3.2")
    Signer = $rootCaCertificate
    Provider = 'Microsoft Strong Cryptographic Provider'
    KeySpec = 'KeyExchange'
    KeyLength = 2048
    HashAlgorithm = 'SHA256'
    KeyAlgorithm = 'RSA'
    NotAfter = (Get-date).AddYears(2)
  }
  $hostCertificate = New-SelfSignedCertificate @params
  Export-Certificate -Cert $hostCertificate -FilePath .\$hostName.cer -Verbose
  Export-PfxCertificate -Cert $hostCertificate -FilePath .\$hostName.pfx -Password $hostPassword -Verbose
  Get-ChildItem cert:\LocalMachine\My |?{$_.subject -eq "CN=$hostName"} | remove-item -force
  Import-PfxCertificate -FilePath .\$hostName.pfx -CertStoreLocation Cert:\LocalMachine\My -password $hostPassword -Exportable -Verbose
  $hostCertificate = Get-ChildItem cert:\LocalMachine\My |?{$_.subject -eq "CN=$hostName"}
}

Get-ChildItem wsman:\localhost\Listener\ | Where-Object -Property Keys -eq 'Transport=HTTPS' | Remove-Item -Recurse
New-Item -Path WSMan:\localhost\Listener -Transport HTTPS -Address * -CertificateThumbPrint $($hostCertificate.Thumbprint) -Force -Verbose

Restart-Service WinRM -Verbose

New-NetFirewallRule -DisplayName "Windows Remote Management (HTTPS-In)" -Name "WinRMHTTPSIn" -Profile Any -LocalPort 5986 -Protocol TCP -Verbose
```
- WinRM allow HTTP
```
# Get the public networks
$PubNets = Get-NetConnectionProfile -NetworkCategory Public -ErrorAction SilentlyContinue 

# Set the profile to private
foreach ($PubNet in $PubNets) {
    Set-NetConnectionProfile -InterfaceIndex $PubNet.InterfaceIndex -NetworkCategory Private
}

# Configure winrm
Set-WSManInstance WinRM/Config/Service -ValueSet @{AllowUnencrypted = $true}

# Restore network categories
foreach ($PubNet in $PubNets) {
    Set-NetConnectionProfile -InterfaceIndex $PubNet.InterfaceIndex -NetworkCategory Public
}

Get-ChildItem wsman:\localhost\Listener\ | Where-Object -Property Keys -eq 'Transport=HTTP' | Remove-Item -Recurse
New-Item -Path WSMan:\localhost\Listener -Transport HTTP -Address * -Force -Verbose

Restart-Service WinRM -Verbose

New-NetFirewallRule -DisplayName "Windows Remote Management (HTTP-In)" -Name "WinRMHTTPIn" -Profile Any -LocalPort 5985 -Protocol TCP -Verbose
```

[Enable Ssl for WinRM using Powershell](https://frontier.town/2011/12/overthere-control-windows-from-java/)

To debug WinRm issues enable debugging by setting environment variable `WINRMCP_DEBUG=1` and `TF_LOG=DEBUG`

Check if you can connect to WinRM
```
$hostName=[System.Net.Dns]::GetHostName()
$winrmPort = "5986"

# Get the credentials of the machine
$cred = Get-Credential

# Connect to the machine
$soptions = New-PSSessionOption -SkipCACheck -SkipCNCheck
Enter-PSSession -ComputerName $hostName -Port $winrmPort -Credential $cred -SessionOption $soptions -UseSSL
```

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/taliesins/terraform-provider-hyperv`

```sh
$ mkdir -p $GOPATH/src/github.com/taliesins; cd $GOPATH/src/github.com/taliesins
$ git clone https://github.com/taliesins/terraform-provider-hyperv.git
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/taliesins/terraform-provider-hyperv
$ make build
```

Using the provider
----------------------
## Fill in for each provider

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.17+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

You should also use the terraform [documentation](https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-provider#prepare-terraform-for-local-provider-install) to setup the terraform environment correctly so that you can use your locally compiled version.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-hyperv
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```

Debugging the Provider
----------------------

To set Terraform log level:
```
set TF_LOG=TRACE
```

To view powershell commands that are sent:
```
set WINRMCP_DEBUG=TRUE
```
