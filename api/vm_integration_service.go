package api

import (
	"log"
	"strconv"
	"strings"
	"text/template"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DefaultVmIntegrationServices() (interface{}, error) {
	flattenedIntegrationServices := make(map[string]interface{})

	flattenedIntegrationServices["VSS"] = true
	flattenedIntegrationServices["Shutdown"] = true
	flattenedIntegrationServices["Time Synchronization"] = true
	flattenedIntegrationServices["Heartbeat"] = true
	flattenedIntegrationServices["Key-Value Pair Exchange"] = true
	flattenedIntegrationServices["Guest Service Interface"] = false

	return flattenedIntegrationServices, nil
}

func getDefaultValueForVmIntegrationService(integrationServiceKey string, _ *schema.ResourceData) bool {
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

func DiffSuppressVmIntegrationServices(key, old, new string, d *schema.ResourceData) bool {
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

	if newValueError != nil {
		newValue = getDefaultValueForVmIntegrationService(integrationServiceKey, d)
		log.Printf("[DEBUG] '[%s]' New value '[%s]' defaulted to '[%v]' ", integrationServiceKey, new, newValue)
	}

	if oldValueError != nil {
		oldValue = getDefaultValueForVmIntegrationService(integrationServiceKey, d)
		log.Printf("[DEBUG] '[%s]' Old value '[%s]' defaulted to '[%v]' ", integrationServiceKey, old, oldValue)
	}

	log.Printf("[DEBUG] '[%s]' Comparing old value '[%v]' with new value '[%v]' ", integrationServiceKey, oldValue, newValue)
	return newValue == oldValue
}

func GetChangedIntegrationServices(vmIntegrationServices []vmIntegrationService, d *schema.ResourceData) []vmIntegrationService {
	changedIntegrationServices := make([]vmIntegrationService, 0)

	for _, integrationServiceValue := range vmIntegrationServices {
		key := "integration_services." + integrationServiceValue.Name

		if d.HasChange(key) {
			log.Printf("[DEBUG] integration service '[%s]' changed", key)
			changedIntegrationServices = append(changedIntegrationServices, integrationServiceValue)
		} else {
			log.Printf("[DEBUG] integration service '[%s]' not changed", key)
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
	flattenedIntegrationServices := make(map[string]interface{})

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

func (c *HypervClient) GetVmIntegrationServices(vmName string) (result []vmIntegrationService, err error) {
	err = c.runScriptWithResult(getVmIntegrationServicesTemplate, getVmIntegrationServicesArgs{
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

func (c *HypervClient) EnableVmIntegrationService(vmName string, name string) (err error) {
	err = c.runFireAndForgetScript(enableVmIntegrationServiceTemplate, enableVmIntegrationServiceArgs{
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

func (c *HypervClient) DisableVmIntegrationService(vmName string, name string) (err error) {
	err = c.runFireAndForgetScript(disableVmIntegrationServiceTemplate, disableVmIntegrationServiceArgs{
		VmName: vmName,
		Name:   name,
	})

	return err
}

func (c *HypervClient) CreateOrUpdateVmIntegrationServices(vmName string, integrationServices []vmIntegrationService) (err error) {
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
