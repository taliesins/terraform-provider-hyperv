package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strconv"
	"strings"
)

type ConsoleModeType int

const (
	ConsoleModeType_Default ConsoleModeType = 0
	ConsoleModeType_Com1    ConsoleModeType = 1
	ConsoleModeType_Com2    ConsoleModeType = 2
	ConsoleModeType_None    ConsoleModeType = 3
)

var ConsoleModeType_name = map[ConsoleModeType]string{
	ConsoleModeType_Default: "Default",
	ConsoleModeType_Com1:    "COM1",
	ConsoleModeType_Com2:    "COM2",
	ConsoleModeType_None:    "None",
}

var ConsoleModeType_value = map[string]ConsoleModeType{
	"default": ConsoleModeType_Default,
	"com1":    ConsoleModeType_Com1,
	"com2":    ConsoleModeType_Com2,
	"none":    ConsoleModeType_None,
}

func (x ConsoleModeType) String() string {
	return ConsoleModeType_name[x]
}

func ToConsoleModeType(x string) ConsoleModeType {
	if integerValue, err := strconv.Atoi(x); err == nil {
		return ConsoleModeType(integerValue)
	}
	return ConsoleModeType_value[strings.ToLower(x)]
}

func (d *ConsoleModeType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(d.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (d *ConsoleModeType) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		var i int
		err2 := json.Unmarshal(b, &i)
		if err2 == nil {
			*d = ConsoleModeType(i)
			return nil
		}

		return err
	}
	*d = ToConsoleModeType(s)
	return nil
}

type IPProtocolPreference int

const (
	IPProtocolPreference_IPv4 IPProtocolPreference = 0
	IPProtocolPreference_IPv6 IPProtocolPreference = 1
)

var IPProtocolPreference_name = map[IPProtocolPreference]string{
	IPProtocolPreference_IPv4: "IPv4",
	IPProtocolPreference_IPv6: "IPv6",
}

var IPProtocolPreference_value = map[string]IPProtocolPreference{
	"ipv4": IPProtocolPreference_IPv4,
	"ipv6": IPProtocolPreference_IPv6,
}

func (x IPProtocolPreference) String() string {
	return IPProtocolPreference_name[x]
}

func ToIPProtocolPreference(x string) IPProtocolPreference {
	if integerValue, err := strconv.Atoi(x); err == nil {
		return IPProtocolPreference(integerValue)
	}
	return IPProtocolPreference_value[strings.ToLower(x)]
}

func (d *IPProtocolPreference) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(d.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (d *IPProtocolPreference) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		var i int
		err2 := json.Unmarshal(b, &i)
		if err2 == nil {
			*d = IPProtocolPreference(i)
			return nil
		}

		return err
	}
	*d = ToIPProtocolPreference(s)
	return nil
}

type VmFirmware struct {
	VmName                       string
	EnableSecureBoot             OnOffState
	SecureBootTemplate           string
	PreferredNetworkBootProtocol IPProtocolPreference
	ConsoleMode                  ConsoleModeType
	PauseAfterBootFailure        OnOffState
}

func DefaultVmFirmwares() (interface{}, error) {
	result := make([]VmFirmware, 0)
	vmFirmware := VmFirmware{
		EnableSecureBoot:             OnOffState_On,
		SecureBootTemplate:           "MicrosoftWindows",
		PreferredNetworkBootProtocol: IPProtocolPreference_IPv4,
		ConsoleMode:                  ConsoleModeType_Default,
		PauseAfterBootFailure:        OnOffState_Off,
	}

	result = append(result, vmFirmware)
	return result, nil
}

func ExpandVmFirmwares(d *schema.ResourceData) ([]VmFirmware, error) {
	expandedVmFirmwares := make([]VmFirmware, 0)

	if v, ok := d.GetOk("vm_firmware"); ok {
		vmFirmwares := v.([]interface{})
		for _, firmware := range vmFirmwares {
			firmware, ok := firmware.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] vm_firmware should be a Hash - was '%+v'", firmware)
			}

			log.Printf("[DEBUG] firmware = [%+v]", firmware)

			expandedVmFirmware := VmFirmware{
				EnableSecureBoot:             ToOnOffState(firmware["enable_secure_boot"].(string)),
				SecureBootTemplate:           firmware["secure_boot_template"].(string),
				PreferredNetworkBootProtocol: ToIPProtocolPreference(firmware["preferred_network_boot_protocol"].(string)),
				ConsoleMode:                  ToConsoleModeType(firmware["console_mode"].(string)),
				PauseAfterBootFailure:        ToOnOffState(firmware["pause_after_boot_failure"].(string)),
			}

			expandedVmFirmwares = append(expandedVmFirmwares, expandedVmFirmware)
		}
	} else {
		vmFirmware := VmFirmware{
			EnableSecureBoot:             OnOffState_On,
			SecureBootTemplate:           "MicrosoftWindows",
			PreferredNetworkBootProtocol: IPProtocolPreference_IPv4,
			ConsoleMode:                  ConsoleModeType_Default,
			PauseAfterBootFailure:        OnOffState_Off,
		}
		expandedVmFirmwares = append(expandedVmFirmwares, vmFirmware)
	}

	return expandedVmFirmwares, nil
}

func FlattenVmFirmwares(vmFirmwares *[]VmFirmware) []interface{} {
	flattenedVmFirmwares := make([]interface{}, 0)

	if vmFirmwares != nil {
		for _, vmFirmware := range *vmFirmwares {
			flattenedVmFirmware := make(map[string]interface{})
			flattenedVmFirmware["enable_secure_boot"] = vmFirmware.EnableSecureBoot.String()
			flattenedVmFirmware["secure_boot_template"] = vmFirmware.SecureBootTemplate
			flattenedVmFirmware["preferred_network_boot_protocol"] = vmFirmware.PreferredNetworkBootProtocol.String()
			flattenedVmFirmware["console_mode"] = vmFirmware.ConsoleMode.String()
			flattenedVmFirmware["pause_after_boot_failure"] = vmFirmware.PauseAfterBootFailure.String()
			flattenedVmFirmwares = append(flattenedVmFirmwares, flattenedVmFirmware)
		}
	}

	return flattenedVmFirmwares
}

type HypervVmFirmwareClient interface {
	CreateOrUpdateVmFirmware(
		vmName string,
		enableSecureBoot OnOffState,
		secureBootTemplate string,
		preferredNetworkBootProtocol IPProtocolPreference,
		consoleMode ConsoleModeType,
		pauseAfterBootFailure OnOffState,
	) (err error)
	GetVmFirmware(vmName string) (result VmFirmware, err error)
	GetNoVmFirmwares() (result []VmFirmware)
	GetVmFirmwares(vmName string) (result []VmFirmware, err error)
	CreateOrUpdateVmFirmwares(vmName string, vmFirmwares []VmFirmware) (err error)
}
