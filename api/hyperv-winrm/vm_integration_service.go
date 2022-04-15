package hyperv_winrm

import (
	"github.com/taliesins/terraform-provider-hyperv/api"
	"text/template"
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

func (c *ClientConfig) GetVmIntegrationServices(vmName string) (result []api.VmIntegrationService, err error) {
	err = c.WinRmClient.RunScriptWithResult(getVmIntegrationServicesTemplate, getVmIntegrationServicesArgs{
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

func (c *ClientConfig) EnableVmIntegrationService(vmName string, name string) (err error) {
	err = c.WinRmClient.RunFireAndForgetScript(enableVmIntegrationServiceTemplate, enableVmIntegrationServiceArgs{
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

func (c *ClientConfig) DisableVmIntegrationService(vmName string, name string) (err error) {
	err = c.WinRmClient.RunFireAndForgetScript(disableVmIntegrationServiceTemplate, disableVmIntegrationServiceArgs{
		VmName: vmName,
		Name:   name,
	})

	return err
}

func (c *ClientConfig) CreateOrUpdateVmIntegrationServices(vmName string, integrationServices []api.VmIntegrationService) (err error) {
	for _, integrationService := range integrationServices {
		if integrationService.Enabled {
			err = c.EnableVmIntegrationService(vmName, integrationService.Name)
		} else {
			err = c.DisableVmIntegrationService(vmName, integrationService.Name)
		}
		if err != nil {
			return err
		}
	}

	return nil
}
