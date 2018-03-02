package api

import (
	"encoding/json"
	"strings"
	"text/template"
)

type ControllerType int

const (
	ControllerType_Ide  ControllerType = 0
	ControllerType_Scsi ControllerType = 1
)

var ControllerType_name = map[ControllerType]string{
	ControllerType_Ide:  "Ide",
	ControllerType_Scsi: "Scsi",
}

var ControllerType_value = map[string]ControllerType{
	"ide":  ControllerType_Ide,
	"scsi": ControllerType_Scsi,
}

func (x ControllerType) String() string {
	return ControllerType_name[x]
}

func ToControllerType(x string) ControllerType {
	return ControllerType_value[strings.ToLower(x)]
}

type CacheAttributes int

const (
	CacheAttributes_Default                 CacheAttributes = 0
	CacheAttributes_WriteCacheEnabled       CacheAttributes = 1
	CacheAttributes_WriteCacheAndFUAEnabled CacheAttributes = 2
	CacheAttributes_WriteCacheDisabled      CacheAttributes = 3
)

var CacheAttributes_name = map[CacheAttributes]string{
	CacheAttributes_Default:                 "Default",
	CacheAttributes_WriteCacheEnabled:       "WriteCacheEnabled",
	CacheAttributes_WriteCacheAndFUAEnabled: "WriteCacheAndFUAEnabled",
	CacheAttributes_WriteCacheDisabled:      "WriteCacheDisabled",
}

var CacheAttributes_value = map[string]CacheAttributes{
	"default":                 CacheAttributes_Default,
	"writecacheenabled":       CacheAttributes_WriteCacheEnabled,
	"writecacheandfuaenabled": CacheAttributes_WriteCacheAndFUAEnabled,
	"writecachedisabled":      CacheAttributes_WriteCacheDisabled,
}

func (x CacheAttributes) String() string {
	return CacheAttributes_name[x]
}

func ToCacheAttributes(x string) CacheAttributes {
	return CacheAttributes_value[strings.ToLower(x)]
}

func ExpandHardDiskDrives(hardDiskDrives *[]map[string]interface{}) []vmHardDiskDrive {
	expandedHardDiskDrives := make([]vmHardDiskDrive, 0)

	for _, hardDiskDrive := range *hardDiskDrives {
		expandedHardDiskDrive := vmHardDiskDrive{
			ControllerType:                ToControllerType(hardDiskDrive["controller_type"].(string)),
			ControllerNumber:              hardDiskDrive["controller_number"].(int),
			ControllerLocation:            hardDiskDrive["controller_location"].(int),
			Path:                          hardDiskDrive["path"].(string),
			DiskNumber:                    hardDiskDrive["disk_number"].(int),
			ResourcePoolName:              hardDiskDrive["resource_pool_name"].(string),
			SupportPersistentReservations: hardDiskDrive["support_persistent_reservations"].(bool),
			MaximumIops:                   hardDiskDrive["maximum_iops"].(int),
			MinimumIops:                   hardDiskDrive["minimum_iops"].(int),
			QosPolicyId:                   hardDiskDrive["qos_policy_id"].(string),
			OverrideCacheAttributes:       ToCacheAttributes(hardDiskDrive["override_cache_attributes"].(string)),
		}

		expandedHardDiskDrives = append(expandedHardDiskDrives, expandedHardDiskDrive)
	}

	if len(expandedHardDiskDrives) > 0 {
		return expandedHardDiskDrives
	}

	return nil
}

func FlattenHardDiskDrives(hardDiskDrives *[]vmHardDiskDrive) []map[string]interface{} {
	flattenedHardDiskDrives := make([]map[string]interface{}, 0)

	for _, hardDiskDrive := range *hardDiskDrives {
		flattenedHardDiskDrive := make(map[string]interface{})
		flattenedHardDiskDrive["controller_type"] = hardDiskDrive.ControllerType
		flattenedHardDiskDrive["controller_number"] = hardDiskDrive.ControllerNumber
		flattenedHardDiskDrive["controller_location"] = hardDiskDrive.ControllerLocation
		flattenedHardDiskDrive["path"] = hardDiskDrive.Path
		flattenedHardDiskDrive["disk_number"] = hardDiskDrive.DiskNumber
		flattenedHardDiskDrive["resource_pool_name"] = hardDiskDrive.ResourcePoolName
		flattenedHardDiskDrive["support_persistent_reservations"] = hardDiskDrive.SupportPersistentReservations
		flattenedHardDiskDrive["maximum_iops"] = hardDiskDrive.MaximumIops
		flattenedHardDiskDrive["minimum_iops"] = hardDiskDrive.MinimumIops
		flattenedHardDiskDrive["qos_policy_id"] = hardDiskDrive.QosPolicyId
		flattenedHardDiskDrive["override_cache_attributes"] = hardDiskDrive.OverrideCacheAttributes
		flattenedHardDiskDrives = append(flattenedHardDiskDrives, flattenedHardDiskDrive)
	}

	if len(flattenedHardDiskDrives) > 0 {
		return flattenedHardDiskDrives
	}

	return nil
}

type vmHardDiskDrive struct {
	VMName                        string
	ControllerType                ControllerType
	ControllerNumber              int
	ControllerLocation            int
	Path                          string
	DiskNumber                    int
	ResourcePoolName              string
	SupportPersistentReservations bool
	MaximumIops                   int
	MinimumIops                   int
	QosPolicyId                   string
	OverrideCacheAttributes       CacheAttributes
	//AllowUnverifiedPaths          bool no way of checking if its turned on so always turn on
}

type createVMHardDiskDriveArgs struct {
	VmHardDiskDriveJson string
}

var createVMHardDiskDriveTemplate = template.Must(template.New("CreateVMHardDiskDrive").Parse(`
$ErrorActionPreference = 'Stop'
Get-Vm | Out-Null
$vmHardDiskDrive = '{{.VmHardDiskDriveJson}}' | ConvertFrom-Json

$NewVmHardDiskDriveArgs = @{
	VMName=$vmHardDiskDrive.VmName
	ControllerType=$vmHardDiskDrive.ControllerType
	ControllerNumber=$vmHardDiskDrive.ControllerNumber
	ControllerLocation=$vmHardDiskDrive.ControllerLocation
	Path=$vmHardDiskDrive.Path
	DiskNumber=$vmHardDiskDrive.DiskNumber
	ResourcePoolName=$vmHardDiskDrive.ResourcePoolName
	SupportPersistentReservations=$vmHardDiskDrive.SupportPersistentReservations
	MaximumIops=$vmHardDiskDrive.MaximumIops
	MinimumIops=$vmHardDiskDrive.MinimumIops
	QosPolicyId=$vmHardDiskDrive.QosPolicyId
	OverrideCacheAttributes=$vmHardDiskDrive.OverrideCacheAttributes
	AllowUnverifiedPaths=$true
}

Add-VmHardDiskDrive @NewVmHardDiskDriveArgs
`))

func (c *HypervClient) CreateVMHardDiskDrive(
	vmName string,
	controllerType ControllerType,
	controllerNumber int,
	controllerLocation int,
	path string,
	diskNumber int,
	resourcePoolName string,
	supportPersistentReservations bool,
	maximumIops int,
	minimumIops int,
	qosPolicyId string,
	overrideCacheAttributes CacheAttributes,

) (err error) {

	vmHardDiskDriveJson, err := json.Marshal(vmHardDiskDrive{
		VMName:                        vmName,
		ControllerType:                controllerType,
		ControllerNumber:              controllerNumber,
		ControllerLocation:            controllerLocation,
		Path:                          path,
		DiskNumber:                    diskNumber,
		ResourcePoolName:              resourcePoolName,
		SupportPersistentReservations: supportPersistentReservations,
		MaximumIops:                   maximumIops,
		MinimumIops:                   minimumIops,
		QosPolicyId:                   qosPolicyId,
		OverrideCacheAttributes:       overrideCacheAttributes,
	})

	err = c.runFireAndForgetScript(createVMHardDiskDriveTemplate, createVMHardDiskDriveArgs{
		VmHardDiskDriveJson: string(vmHardDiskDriveJson),
	})

	return err
}

type getVMHardDiskDrivesArgs struct {
	VMName string
}

var getVMHardDiskDrivesTemplate = template.Must(template.New("GetVMHardDiskDrives").Parse(`
$ErrorActionPreference = 'Stop'
$vmHardDiskDrivesObject = Get-VMHardDiskDrive -VMName '{{.VMName}}' | %{ @{
	ControllerType=$_.ControllerType;
	ControllerNumber=$_.ControllerNumber;
	ControllerLocation=$_.ControllerLocation;
	Path=$_.Path;
	DiskNumber=$_.DiskNumber;
	ResourcePoolName=$_.PoolName;
	SupportPersistentReservations=$_.SupportPersistentReservations;
	MaximumIops=$_.MaximumIops;
	MinimumIops=$_.MinimumIops;
	QosPolicyId=$_.QosPolicyId;	
	OverrideCacheAttributes=$_.WriteHardeningMethod;
}}

if ($vmHardDiskDrivesObject) {
	$vmHardDiskDrives = ConvertTo-Json -InputObject $vmHardDiskDrivesObject
	$vmHardDiskDrives
} else {
	"{}"
}
`))

func (c *HypervClient) GetVMHardDiskDrives(vmName string) (result []vmHardDiskDrive, err error) {
	err = c.runScriptWithResult(getVMHardDiskDrivesTemplate, getVMHardDiskDrivesArgs{
		VMName: vmName,
	}, &result)

	return result, err
}

type updateVMHardDiskDriveArgs struct {
	VMName              string
	ControllerNumber    int
	ControllerLocation  int
	VmHardDiskDriveJson string
}

var updateVMHardDiskDriveTemplate = template.Must(template.New("UpdateVMHardDiskDrive").Parse(`
$ErrorActionPreference = 'Stop'
Get-Vm | Out-Null
$vmHardDiskDrive = '{{.VmHardDiskDriveJson}}' | ConvertFrom-Json

$vmHardDiskDrivesObject = @(Get-VMHardDiskDrive -VMName '{{.VMName}}' -ControllerLocation {{.ControllerLocation}} -ControllerNumber {{.ControllerNumber}} )

if (!$vmHardDiskDrivesObject){
	throw "VM hard disk drive does not exist - {{.ControllerLocation}} {{.ControllerNumber}}"
}

$SetVmHardDiskDriveArgs = @{}
$SetVmHardDiskDriveArgs.VMName=$vmHardDiskDrivesObject.VMName
$SetVmHardDiskDriveArgs.ControllerType=$vmHardDiskDrivesObject.ControllerType
$SetVmHardDiskDriveArgs.ControllerLocation=$vmHardDiskDrivesObject.ControllerLocation
$SetVmHardDiskDriveArgs.ControllerNumber=$vmHardDiskDrivesObject.ControllerNumber
$SetVmHardDiskDriveArgs.ToControllerLocation=$vmHardDiskDrive.ControllerLocation
$SetVmHardDiskDriveArgs.ToControllerNumber=$vmHardDiskDrive.ControllerNumber
$SetVmHardDiskDriveArgs.Path=$vmHardDiskDrive.Path
$SetVmHardDiskDriveArgs.DiskNumber=$vmHardDiskDrive.DiskNumber
$SetVmHardDiskDriveArgs.ResourcePoolName=$vmHardDiskDrive.ResourcePoolName
$SetVmHardDiskDriveArgs.SupportPersistentReservations=$vmHardDiskDrive.SupportPersistentReservations
$SetVmHardDiskDriveArgs.MaximumIops=$vmHardDiskDrive.MaximumIops
$SetVmHardDiskDriveArgs.MinimumIops=$vmHardDiskDrive.MinimumIops
$SetVmHardDiskDriveArgs.QosPolicyId=$vmHardDiskDrive.QosPolicyId
$SetVmHardDiskDriveArgs.OverrideCacheAttributes=$vmHardDiskDrive.OverrideCacheAttributes	
$SetVmHardDiskDriveArgs.AllowUnverifiedPaths=$true

Set-VMHardDiskDrive @SetVmHardDiskDriveArgs

`))

func (c *HypervClient) UpdateVMHardDiskDrive(
	vmName string,
	controllerNumber int,
	controllerLocation int,
	controllerType ControllerType,
	toControllerNumber int,
	toControllerLocation int,
	path string,
	diskNumber int,
	resourcePoolName string,
	supportPersistentReservations bool,
	maximumIops int,
	minimumIops int,
	qosPolicyId string,
	overrideCacheAttributes CacheAttributes,
) (err error) {

	vmHardDiskDriveJson, err := json.Marshal(vmHardDiskDrive{
		VMName:                        vmName,
		ControllerType:                controllerType,
		ControllerNumber:              toControllerNumber,
		ControllerLocation:            toControllerLocation,
		Path:                          path,
		DiskNumber:                    diskNumber,
		ResourcePoolName:              resourcePoolName,
		SupportPersistentReservations: supportPersistentReservations,
		MaximumIops:                   maximumIops,
		MinimumIops:                   minimumIops,
		QosPolicyId:                   qosPolicyId,
		OverrideCacheAttributes:       overrideCacheAttributes,
	})

	err = c.runFireAndForgetScript(updateVMHardDiskDriveTemplate, updateVMHardDiskDriveArgs{
		VMName:              vmName,
		ControllerNumber:    controllerNumber,
		ControllerLocation:  controllerLocation,
		VmHardDiskDriveJson: string(vmHardDiskDriveJson),
	})

	return err
}

type deleteVMHardDiskDriveArgs struct {
	VMName             string
	ControllerNumber   int
	ControllerLocation int
}

var deleteVMHardDiskDriveTemplate = template.Must(template.New("DeleteVMHardDiskDrive").Parse(`
$ErrorActionPreference = 'Stop'

@(Get-VMHardDiskDrive -VMName '{{.VMName}}' -ControllerNumber {{.ControllerNumber}} -ControllerLocation {{.ControllerLocation}}) | Remove-VMHardDiskDrive -Force
`))

func (c *HypervClient) DeleteVMHardDiskDrive(vmname string, controllerNumber int, controllerLocation int) (err error) {
	err = c.runFireAndForgetScript(deleteVMHardDiskDriveTemplate, deleteVMHardDiskDriveArgs{
		VMName:             vmname,
		ControllerNumber:   controllerNumber,
		ControllerLocation: controllerLocation,
	})

	return err
}

func (c *HypervClient) CreateOrUpdateVMHardDiskDrives(vmName string, hardDiskDrives []vmHardDiskDrive) (err error) {
	currentHardDiskDrives, err := c.GetVMHardDiskDrives(vmName)
	if err != nil {
		return err
	}

	for i := len(currentHardDiskDrives) - 1; i > len(hardDiskDrives)-1; i-- {
		currentHardDiskDrive := currentHardDiskDrives[i]
		err = c.DeleteVMHardDiskDrive(currentHardDiskDrive.VMName, currentHardDiskDrive.ControllerNumber, currentHardDiskDrive.ControllerLocation)
		if err != nil {
			return err
		}
	}

	for i := 0; i < len(currentHardDiskDrives)-1; i++ {
		currentHardDiskDrive := currentHardDiskDrives[i]
		hardDiskDrive := hardDiskDrives[i]

		err = c.UpdateVMHardDiskDrive(
			currentHardDiskDrive.VMName,
			currentHardDiskDrive.ControllerNumber,
			currentHardDiskDrive.ControllerLocation,
			hardDiskDrive.ControllerType,
			hardDiskDrive.ControllerNumber,
			hardDiskDrive.ControllerLocation,
			hardDiskDrive.Path,
			hardDiskDrive.DiskNumber,
			hardDiskDrive.ResourcePoolName,
			hardDiskDrive.SupportPersistentReservations,
			hardDiskDrive.MaximumIops,
			hardDiskDrive.MinimumIops,
			hardDiskDrive.QosPolicyId,
			hardDiskDrive.OverrideCacheAttributes,
		)
		if err != nil {
			return err
		}
	}

	for i := len(currentHardDiskDrives) - 1; i < len(hardDiskDrives)-1; i++ {
		hardDiskDrive := hardDiskDrives[i]
		err = c.CreateVMHardDiskDrive(
			hardDiskDrive.VMName,
			hardDiskDrive.ControllerType,
			hardDiskDrive.ControllerNumber,
			hardDiskDrive.ControllerLocation,
			hardDiskDrive.Path,
			hardDiskDrive.DiskNumber,
			hardDiskDrive.ResourcePoolName,
			hardDiskDrive.SupportPersistentReservations,
			hardDiskDrive.MaximumIops,
			hardDiskDrive.MinimumIops,
			hardDiskDrive.QosPolicyId,
			hardDiskDrive.OverrideCacheAttributes,
		)

		if err != nil {
			return err
		}
	}

	return nil
}
