package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"text/template"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

type vmFirmware struct {
	VmName                       string
	EnableSecureBoot             OnOffState
	SecureBootTemplate           string
	PreferredNetworkBootProtocol IPProtocolPreference
	ConsoleMode                  ConsoleModeType
	PauseAfterBootFailure        OnOffState
}

func DefaultVmFirmwares() (interface{}, error) {
	result := make([]vmFirmware, 0)
	vmFirmware := vmFirmware{
		EnableSecureBoot:             OnOffState_On,
		SecureBootTemplate:           "MicrosoftWindows",
		PreferredNetworkBootProtocol: IPProtocolPreference_IPv4,
		ConsoleMode:                  ConsoleModeType_Default,
		PauseAfterBootFailure:        OnOffState_Off,
	}

	result = append(result, vmFirmware)
	return result, nil
}

func ExpandVmFirmwares(d *schema.ResourceData) ([]vmFirmware, error) {
	expandedVmFirmwares := make([]vmFirmware, 0)

	if v, ok := d.GetOk("vm_firmware"); ok {
		vmFirmwares := v.([]interface{})
		for _, firmware := range vmFirmwares {
			firmware, ok := firmware.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] vm_firmware should be a Hash - was '%+v'", firmware)
			}

			log.Printf("[DEBUG] firmware = [%+v]", firmware)

			expandedVmFirmware := vmFirmware{
				EnableSecureBoot:             ToOnOffState(firmware["enable_secure_boot"].(string)),
				SecureBootTemplate:           firmware["secure_boot_template"].(string),
				PreferredNetworkBootProtocol: ToIPProtocolPreference(firmware["preferred_network_boot_protocol"].(string)),
				ConsoleMode:                  ToConsoleModeType(firmware["console_mode"].(string)),
				PauseAfterBootFailure:        ToOnOffState(firmware["pause_after_boot_failure"].(string)),
			}

			expandedVmFirmwares = append(expandedVmFirmwares, expandedVmFirmware)
		}
	}

	return expandedVmFirmwares, nil
}

func FlattenVmFirmwares(vmFirmwares *[]vmFirmware) []interface{} {
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

func (c *HypervClient) CreateOrUpdateVmFirmware(
	vmName string,
	enableSecureBoot OnOffState,
	secureBootTemplate string,
	preferredNetworkBootProtocol IPProtocolPreference,
	consoleMode ConsoleModeType,
	pauseAfterBootFailure OnOffState,
) (err error) {
	vmFirmwareJson, err := json.Marshal(vmFirmware{
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

	err = c.runFireAndForgetScript(createOrUpdateVmFirmwareTemplate, createOrUpdateVmFirmwareArgs{
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

func (c *HypervClient) GetVmFirmware(vmName string) (result vmFirmware, err error) {
	err = c.runScriptWithResult(getVmFirmwareTemplate, getVmFirmwareArgs{
		VmName: vmName,
	}, &result)

	return result, err
}

func (c *HypervClient) GetNoVmFirmwares() (result []vmFirmware) {
	result = make([]vmFirmware, 0)
	return result
}

func (c *HypervClient) GetVmFirmwares(vmName string) (result []vmFirmware, err error) {
	result = make([]vmFirmware, 0)
	vmFirmware, err := c.GetVmFirmware(vmName)
	if err != nil {
		return result, err
	}
	result = append(result, vmFirmware)
	return result, err
}

func (c *HypervClient) CreateOrUpdateVmFirmwares(vmName string, vmFirmwares []vmFirmware) (err error) {
	if len(vmFirmwares) == 0 {
		return nil
	}
	if len(vmFirmwares) > 1 {
		return fmt.Errorf("Only 1 vm firmware setting allowed per a vm")
	}

	vmFirmware := vmFirmwares[0]

	return c.CreateOrUpdateVmFirmware(vmName,
		vmFirmware.EnableSecureBoot,
		vmFirmware.SecureBootTemplate,
		vmFirmware.PreferredNetworkBootProtocol,
		vmFirmware.ConsoleMode,
		vmFirmware.PauseAfterBootFailure,
	)
}
