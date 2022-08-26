package hyperv_winrm

import (
	"context"
	"encoding/json"
	"fmt"
	"text/template"

	"github.com/taliesins/terraform-provider-hyperv/api"
)

type createOrUpdateVmFirmwareArgs struct {
	VmFirmwareJson string
}

var createOrUpdateVmFirmwareTemplate = template.Must(template.New("CreateOrUpdateVmFirmware").Parse(`
$ErrorActionPreference = 'Stop'
Import-Module Hyper-V
$vmFirmware = '{{.VmFirmwareJson}}' | ConvertFrom-Json

$SetVMFirmwareArgs = @{}
$SetVMFirmwareArgs.VMName=$vmFirmware.VmName

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
	enableSecureBoot api.OnOffState,
	secureBootTemplate string,
	preferredNetworkBootProtocol api.IPProtocolPreference,
	consoleMode api.ConsoleModeType,
	pauseAfterBootFailure api.OnOffState,
) (err error) {
	vmFirmwareJson, err := json.Marshal(api.VmFirmware{
		VmName:                       vmName,
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

$vmFirmwareObject = Get-VMFirmware -VMName '{{.VmName}}' | %{ @{
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
		vmFirmware.EnableSecureBoot,
		vmFirmware.SecureBootTemplate,
		vmFirmware.PreferredNetworkBootProtocol,
		vmFirmware.ConsoleMode,
		vmFirmware.PauseAfterBootFailure,
	)
}
