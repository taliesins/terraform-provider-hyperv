package hyperv_winrm

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/taliesins/terraform-provider-hyperv/api"
	"text/template"
)

type createOrUpdateVmFirmwareArgs struct {
	VmFirmwareJson string
}

var createOrUpdateVmFirmwareTemplate = template.Must(template.New("CreateOrUpdateVmFirmware").Parse(`
$ErrorActionPreference = 'Stop'
Import-Module Hyper-V
$vmFirmware = '{{.VmFirmwareJson}}' | ConvertFrom-Json

$bootOrders = @($vmFirmware.BootOrders | %{
	$bootOrder = $_
	if ($bootOrder.Type -eq 'NetworkAdapter') {
		$networkAdapter = Get-VM -Name "$($vmFirmware.VmName)*" | ?{$_.Name -eq $vmFirmware.VmName } | Get-VMNetworkAdapter
		if ($bootOrder.NetworkAdapterName) {
			$networkAdapter = $networkAdapter | ?{$_.Name -eq $bootOrder.NetworkAdapterName}
		}

		if ($bootOrder.SwitchName) {
			$networkAdapter = $networkAdapter | ?{$_.SwitchName -eq $bootOrder.SwitchName}
		}

		if ($bootOrder.MacAddress) {
			$networkAdapter = $networkAdapter | ?{$_.MacAddress -ieq $bootOrder.MacAddress}
		}

		$networkAdapter
	} elseif ($bootOrder.Type -eq 'HardDiskDrive') {
		$hardDiskDrive = Get-VM -Name "$($vmFirmware.VmName)*" | ?{$_.Name -eq $vmFirmware.VmName } | Get-VMHardDiskDrive

		if ($bootOrder.Path) {
			$hardDiskDrive = $hardDiskDrive | ?{$_.Path -ieq $bootOrder.Path}
		}

		if ($bootOrder.ControllerNumber -gt -1) {
			$hardDiskDrive = $hardDiskDrive | ?{$_.ControllerNumber -eq $bootOrder.ControllerNumber}
		}

		if ($bootOrder.ControllerLocation -gt -1) {
			$hardDiskDrive = $hardDiskDrive | ?{$_.ControllerLocation -eq $bootOrder.ControllerLocation}
		}

		$hardDiskDrive

	} elseif ($bootOrder.Type -eq 'DvdDrive') {
		$dvdDrive = Get-VM -Name "$($vmFirmware.VmName)*" | ?{$_.Name -eq $vmFirmware.VmName } | Get-VMDvdDrive

		if ($bootOrder.Path) {
			$dvdDrive = $dvdDrive | ?{$_.Path -ieq $bootOrder.Path}
		}

		if ($bootOrder.ControllerNumber -gt -1) {
			$dvdDrive = $dvdDrive | ?{$_.ControllerNumber -eq $bootOrder.ControllerNumber}
		}

		if ($bootOrder.ControllerLocation -gt -1) {
			$dvdDrive = $dvdDrive | ?{$_.ControllerLocation -eq $bootOrder.ControllerLocation}
		}

		$dvdDrive
	}
})

$SetVMFirmwareArgs = @{}
$SetVMFirmwareArgs.VMName=$vmFirmware.VmName
$SetVMFirmwareArgs.BootOrder=$bootOrders
$SetVMFirmwareArgs.EnableSecureBoot=$vmFirmware.EnableSecureBoot
$SetVMFirmwareArgs.SecureBootTemplate=$vmFirmware.SecureBootTemplate
$SetVMFirmwareArgs.PreferredNetworkBootProtocol=$vmFirmware.PreferredNetworkBootProtocol
$SetVMFirmwareArgs.ConsoleMode=$vmFirmware.ConsoleMode
$SetVMFirmwareArgs.PauseAfterBootFailure=$vmFirmware.PauseAfterBootFailure

Set-VMFirmware @SetVMFirmwareArgs
`))

func (c *ClientConfig) CreateOrUpdateVmFirmware(
	ctx context.Context,
	vmName string,
	bootOrders []api.Gen2BootOrder,
	enableSecureBoot api.OnOffState,
	secureBootTemplate string,
	preferredNetworkBootProtocol api.IPProtocolPreference,
	consoleMode api.ConsoleModeType,
	pauseAfterBootFailure api.OnOffState,
) (err error) {
	vmFirmwareJson, err := json.Marshal(api.VmFirmware{
		VmName:                       vmName,
		BootOrders:                   bootOrders,
		EnableSecureBoot:             enableSecureBoot,
		SecureBootTemplate:           secureBootTemplate,
		PreferredNetworkBootProtocol: preferredNetworkBootProtocol,
		ConsoleMode:                  consoleMode,
		PauseAfterBootFailure:        pauseAfterBootFailure,
	})

	if err != nil {
		return err
	}

	err = c.WinRmClient.RunFireAndForgetScript(ctx, createOrUpdateVmFirmwareTemplate, createOrUpdateVmFirmwareArgs{
		VmFirmwareJson: string(vmFirmwareJson),
	})

	return err
}

type getVmFirmwareArgs struct {
	VmName string
}

var getVmFirmwareTemplate = template.Must(template.New("GetVmFirmware").Parse(`
$ErrorActionPreference = 'Stop'

$vmFirmwareObject = Get-VM -Name '{{.VmName}}*' | ?{$_.Name -eq '{{.VmName}}' } | Get-VMFirmware | %{ @{
	BootOrders= @($_.BootOrder | %{
		if ($_.BootType -eq 'Network') {
			@{Type='NetworkAdapter';NetworkAdapterName=$_.Device.Name;SwitchName=$_.Device.SwitchName;MacAddress=$_.Device.MacAddress;Path='';ControllerNumber=-1;ControllerLocation=-1;}
		} elseif ($_.BootType -eq 'Drive') {
			@{Type=@(if ($_.Device.Name.StartsWith('Hard Drive')) { 'HardDiskDrive' } else {'DvdDrive'});NetworkAdapterName='';SwitchName='';MacAddress='';Path=$_.Device.Path;ControllerNumber=$_.Device.ControllerNumber;ControllerLocation=$_.Device.ControllerLocation;}
		}
	})
	EnableSecureBoot=             $_.SecureBoot
	SecureBootTemplate=           $_.SecureBootTemplate
	PreferredNetworkBootProtocol= $_.PreferredNetworkBootProtocol
	ConsoleMode=                  $_.ConsoleMode
	PauseAfterBootFailure=        $_.PauseAfterBootFailure
}}

if ($vmFirmwareObject) {
	$vmFirmware = ConvertTo-Json -InputObject $vmFirmwareObject
	$vmFirmware
} else {
	"{}"
}
`))

func (c *ClientConfig) GetVmFirmware(ctx context.Context, vmName string) (result api.VmFirmware, err error) {
	err = c.WinRmClient.RunScriptWithResult(ctx, getVmFirmwareTemplate, getVmFirmwareArgs{
		VmName: vmName,
	}, &result)

	return result, err
}

func (c *ClientConfig) GetNoVmFirmwares(ctx context.Context) (result []api.VmFirmware) {
	result = make([]api.VmFirmware, 0)
	return result
}

func (c *ClientConfig) GetVmFirmwares(ctx context.Context, vmName string) (result []api.VmFirmware, err error) {
	result = make([]api.VmFirmware, 0)
	vmFirmware, err := c.GetVmFirmware(ctx, vmName)
	if err != nil {
		return result, err
	}
	result = append(result, vmFirmware)
	return result, err
}

func (c *ClientConfig) CreateOrUpdateVmFirmwares(ctx context.Context, vmName string, vmFirmwares []api.VmFirmware) (err error) {
	if len(vmFirmwares) == 0 {
		return nil
	}
	if len(vmFirmwares) > 1 {
		return fmt.Errorf("Only 1 vm firmware setting allowed per a vm")
	}

	vmFirmware := vmFirmwares[0]

	return c.CreateOrUpdateVmFirmware(ctx, vmName,
		vmFirmware.BootOrders,
		vmFirmware.EnableSecureBoot,
		vmFirmware.SecureBootTemplate,
		vmFirmware.PreferredNetworkBootProtocol,
		vmFirmware.ConsoleMode,
		vmFirmware.PauseAfterBootFailure,
	)
}
