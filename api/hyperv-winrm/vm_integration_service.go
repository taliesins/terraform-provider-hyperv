package hyperv_winrm

import (
	"context"
	"text/template"

	"github.com/taliesins/terraform-provider-hyperv/api"
)

type getVmIntegrationServicesArgs struct {
	VmName string
}

var getVmIntegrationServicesTemplate = template.Must(template.New("GetVmIntegrationServices").Parse(`
$ErrorActionPreference = 'Stop'
$vmIntegrationServicesObject = @(Get-VMIntegrationService -VmName '{{.VmName}}' | %{ @{
	Name=$_.Name;
	Enabled=$_.Enabled;
}})

if ($vmIntegrationServicesObject) {
	$vmIntegrationServices = ConvertTo-Json -InputObject $vmIntegrationServicesObject
	$vmIntegrationServices
} else {
	"[]"
}
`))

func (c *ClientConfig) GetVmIntegrationServices(ctx context.Context, vmName string) (result []api.VmIntegrationService, err error) {
	err = c.WinRmClient.RunScriptWithResult(ctx, getVmIntegrationServicesTemplate, getVmIntegrationServicesArgs{
		VmName: vmName,
	}, &result)

	return result, err
}

type enableVmIntegrationServiceArgs struct {
	VmName string
	Name   string
}

var enableVmIntegrationServiceTemplate = template.Must(template.New("EnableVmIntegrationService").Parse(`
$ErrorActionPreference = 'Stop'

Enable-VMIntegrationService -VmName '{{.VmName}}' -Name '{{.Name}}'
`))

func (c *ClientConfig) EnableVmIntegrationService(ctx context.Context, vmName string, name string) (err error) {
	err = c.WinRmClient.RunFireAndForgetScript(ctx, enableVmIntegrationServiceTemplate, enableVmIntegrationServiceArgs{
		VmName: vmName,
		Name:   name,
	})

	return err
}

type disableVmIntegrationServiceArgs struct {
	VmName string
	Name   string
}

var disableVmIntegrationServiceTemplate = template.Must(template.New("DisableVmIntegrationService").Parse(`
$ErrorActionPreference = 'Stop'

Disable-VMIntegrationService -VmName '{{.VmName}}' -Name '{{.Name}}'
`))

func (c *ClientConfig) DisableVmIntegrationService(ctx context.Context, vmName string, name string) (err error) {
	err = c.WinRmClient.RunFireAndForgetScript(ctx, disableVmIntegrationServiceTemplate, disableVmIntegrationServiceArgs{
		VmName: vmName,
		Name:   name,
	})

	return err
}

func (c *ClientConfig) CreateOrUpdateVmIntegrationServices(ctx context.Context, vmName string, integrationServices []api.VmIntegrationService) (err error) {
	for _, integrationService := range integrationServices {
		if integrationService.Enabled {
			err = c.EnableVmIntegrationService(ctx, vmName, integrationService.Name)
		} else {
			err = c.DisableVmIntegrationService(ctx, vmName, integrationService.Name)
		}
		if err != nil {
			return err
		}
	}

	return nil
}
