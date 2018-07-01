package api

import (
	"text/template"
	"github.com/hashicorp/terraform/helper/schema"
	"strconv"
	"log"
	"strings"
)

func DefaultVmIntegrationServices() (interface{}, error) {
	flattenedIntegrationServices := make(map[string]interface{}, 0)

	flattenedIntegrationServices["VSS"] = true
	flattenedIntegrationServices["Shutdown"] = true
	flattenedIntegrationServices["Time Synchronization"] = true
	flattenedIntegrationServices["Heartbeat"] = true
	flattenedIntegrationServices["Key-Value Pair Exchange"] = true
	flattenedIntegrationServices["Guest Service Interface"] = false

	return flattenedIntegrationServices, nil
}

func getDefaultValueForVmIntegrationService(integrationServiceKey string, d *schema.ResourceData) bool{
	v, _ := DefaultVmIntegrationServices()
	integrationServices := v.(map[string]interface{})
	if integrationServiceValueInterface, found := integrationServices[integrationServiceKey]; found {
		if integrationServiceValue, ok := integrationServiceValueInterface.(bool); ok {
			return integrationServiceValue
		}
		//its not a bool something went wrong
	}

	return false
}

func DiffSuppressVmIntegrationServices (key, old, new string, d *schema.ResourceData) bool {
	integrationServiceKey := strings.TrimPrefix(key, "integration_services.")

	if integrationServiceKey == "%" {
		//We do not care about the number of elements as we only tack things we have specified
		return true
	}

	if new == "" {
		//We have not explicitly set a value, so allow any value as we are not tracking it
		return true
	}

	newValue, newValueError := strconv.ParseBool(new)
	oldValue, oldValueError := strconv.ParseBool(old)

	if newValueError != nil  {
		newValue = getDefaultValueForVmIntegrationService(integrationServiceKey, d)
		log.Printf("'[%s]' New value '[%s]' defaulted to '[%v]' ", integrationServiceKey, new, newValue)
	}

	if oldValueError != nil  {
		oldValue = getDefaultValueForVmIntegrationService(integrationServiceKey, d)
		log.Printf("'[%s]' Old value '[%s]' defaulted to '[%v]' ", integrationServiceKey, old, oldValue)
	}

	log.Printf("'[%s]' Comparing old value '[%v]' with new value '[%v]' ", integrationServiceKey, oldValue, newValue)
	return newValue == oldValue
}

func GetChangedIntegrationServices(vmIntegrationServices []vmIntegrationService, d *schema.ResourceData) []vmIntegrationService {
	changedIntegrationServices := make([]vmIntegrationService, 0)

	for _, integrationServiceValue := range vmIntegrationServices {
		key := "integration_services."+integrationServiceValue.Name

		if d.HasChange(key) {
			log.Printf("integration service '[%s]' changed", key)
			changedIntegrationServices = append(changedIntegrationServices, integrationServiceValue)
		} else {
			log.Printf("integration service '[%s]' not changed", key)
		}
	}

	return changedIntegrationServices
}

func ExpandIntegrationServices(d *schema.ResourceData) ([]vmIntegrationService, error) {
	expandedIntegrationServices := make([]vmIntegrationService, 0)

	if v, ok := d.GetOk("integration_services"); ok {
		integrationServices := v.(map[string]interface{})

		for integrationServiceKey, integrationServiceValue := range integrationServices {
			integrationService := vmIntegrationService{
				Name:    integrationServiceKey,
				Enabled: integrationServiceValue.(bool),
			}

			expandedIntegrationServices = append(expandedIntegrationServices, integrationService)
		}
	}

	return expandedIntegrationServices, nil
}

func FlattenIntegrationServices(integrationServices *[]vmIntegrationService) map[string]interface{} {
	flattenedIntegrationServices := make(map[string]interface{}, 0)

	if integrationServices != nil {
		for _, integrationService := range *integrationServices {
			flattenedIntegrationServices[integrationService.Name] = integrationService.Enabled
		}
	}

	return flattenedIntegrationServices
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
$vmIntegrationServicesObject = @(Get-VMIntegrationService -VMName '{{.VMName}}' | %{ @{
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
