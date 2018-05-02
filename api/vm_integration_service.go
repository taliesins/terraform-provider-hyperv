package api

import (
	"text/template"
)

func ExpandIntegrationServices(integrationServices *[]map[string]interface{}) []vmIntegrationService {
	expandedIntegrationServices := make([]vmIntegrationService, 0)

	for _, integrationService := range *integrationServices {
		expandedIntegrationService := vmIntegrationService{
			Name:    integrationService["name"].(string),
			Enabled: integrationService["enabled"].(bool),
		}

		expandedIntegrationServices = append(expandedIntegrationServices, expandedIntegrationService)
	}

	if len(expandedIntegrationServices) > 0 {
		return expandedIntegrationServices
	}

	return nil
}

func FlattenIntegrationServices(integrationServices *[]vmIntegrationService) []map[string]interface{} {
	flattenedIntegrationServices := make([]map[string]interface{}, 0)

	for _, integrationService := range *integrationServices {
		flattenedIntegrationService := make(map[string]interface{})
		flattenedIntegrationService["name"] = integrationService.Name
		flattenedIntegrationService["enabled"] = integrationService.Enabled
		flattenedIntegrationServices = append(flattenedIntegrationServices, flattenedIntegrationService)
	}

	if len(flattenedIntegrationServices) > 0 {
		return flattenedIntegrationServices
	}

	return nil
}

type vmIntegrationService struct {
	Name    string
	Enabled bool
}

type getVMIntegrationServicesArgs struct {
	VMName string
}

var getVMIntegrationServicesTemplate = template.Must(template.New("GetVMIntegrationServices").Parse(`
$ErrorActionPreference = 'Stop'
$vmIntegrationServicesObject = Get-VMIntegrationService -VMName '{{.VMName}}' | %{ @{
	Name=$_.Name;
	Enabled=$_.Enabled;
}}

if ($vmIntegrationServicesObject) {
	$vmIntegrationServices = ConvertTo-Json -InputObject $vmIntegrationServicesObject
	$vmIntegrationServices
} else {
	"{}"
}
`))

func (c *HypervClient) GetVMIntegrationServices(vmName string) (result []vmIntegrationService, err error) {
	err = c.runScriptWithResult(getVMIntegrationServicesTemplate, getVMIntegrationServicesArgs{
		VMName: vmName,
	}, &result)

	return result, err
}

type enableVMIntegrationServiceArgs struct {
	VMName string
	Name   string
}

var enableVMIntegrationServiceTemplate = template.Must(template.New("EnableVMIntegrationService").Parse(`
$ErrorActionPreference = 'Stop'

Enable-VMIntegrationService -VMName '{{.VMName}}' -Name '{{.Name}}'
`))

func (c *HypervClient) EnableVMIntegrationService(vmname string, name string) (err error) {
	err = c.runFireAndForgetScript(enableVMIntegrationServiceTemplate, enableVMIntegrationServiceArgs{
		VMName: vmname,
		Name:   name,
	})

	return err
}

type disableVMIntegrationServiceArgs struct {
	VMName string
	Name   string
}

var disableVMIntegrationServiceTemplate = template.Must(template.New("DisableVMIntegrationService").Parse(`
$ErrorActionPreference = 'Stop'

Disable-VMIntegrationService -VMName '{{.VMName}}' -Name '{{.Name}}'
`))

func (c *HypervClient) DisableVMIntegrationService(vmname string, name string) (err error) {
	err = c.runFireAndForgetScript(disableVMIntegrationServiceTemplate, disableVMIntegrationServiceArgs{
		VMName: vmname,
		Name:   name,
	})

	return err
}

func (c *HypervClient) CreateOrUpdateVMIntegrationServices(vmName string, integrationServices []vmIntegrationService) (err error) {
	for _, integrationService := range integrationServices {
		if integrationService.Enabled {
			err = c.EnableVMIntegrationService(vmName, integrationService.Name)
		} else {
			err = c.DisableVMIntegrationService(vmName, integrationService.Name)
		}
		if err != nil {
			return err
		}
	}

	return nil
}
