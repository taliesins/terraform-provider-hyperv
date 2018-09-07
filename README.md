HyperV Terraform Provider
=========================

#This is beta code. Working on adding acceptance tests so that I can mark it as release quality.

- Website: https://github.com/taliesins/terraform-provider-hyperv
- Documentation: https://github.com/taliesins/terraform-provider-hyperv/blob/master/website/docs/index.html.markdown
- Issues: https://github.com/taliesins/terraform-provider-hyperv/issues

![Hashi Logo](/website/logo-hashicorp.svg?raw=true "Hashi Logo")
![Windows Server Logo](/website/windows-server-2016-logo.svg?raw=true "Windows Server Logo")

Features
------------
- Remote scheduled task powershell runner does not run into issues with escaping variables or escaping between the different scripting layers.
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

-	[Terraform](https://www.terraform.io/downloads.html) 0.10.x
-	[Go](https://golang.org/doc/install) 1.8 (to build the provider plugin)
-   Connectivity and credentials to a server running HyperV

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

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.8+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make bin
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
