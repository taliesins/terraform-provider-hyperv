package api

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strconv"
	"strings"
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

func GetChangedIntegrationServices(vmIntegrationServices []VmIntegrationService, d *schema.ResourceData) []VmIntegrationService {
	changedIntegrationServices := make([]VmIntegrationService, 0)

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

func ExpandIntegrationServices(d *schema.ResourceData) ([]VmIntegrationService, error) {
	expandedIntegrationServices := make([]VmIntegrationService, 0)

	if v, ok := d.GetOk("integration_services"); ok {
		integrationServices := v.(map[string]interface{})

		for integrationServiceKey, integrationServiceValue := range integrationServices {
			integrationService := VmIntegrationService{
				Name:    integrationServiceKey,
				Enabled: integrationServiceValue.(bool),
			}

			expandedIntegrationServices = append(expandedIntegrationServices, integrationService)
		}
	}

	return expandedIntegrationServices, nil
}

func FlattenIntegrationServices(integrationServices *[]VmIntegrationService) map[string]interface{} {
	if integrationServices == nil || len(*integrationServices) < 1 {
		return nil
	}

	flattenedIntegrationServices := make(map[string]interface{})

	for _, integrationService := range *integrationServices {
		flattenedIntegrationServices[integrationService.Name] = integrationService.Enabled
	}

	return flattenedIntegrationServices
}

type VmIntegrationService struct {
	Name    string
	Enabled bool
}

type HypervVmIntegrationServiceClient interface {
	GetVmIntegrationServices(vmName string) (result []VmIntegrationService, err error)
	EnableVmIntegrationService(vmName string, name string) (err error)
	DisableVmIntegrationService(vmName string, name string) (err error)
	CreateOrUpdateVmIntegrationServices(vmName string, integrationServices []VmIntegrationService) (err error)
}
