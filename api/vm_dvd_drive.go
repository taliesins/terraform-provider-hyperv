package api

import (
	"encoding/json"
	"text/template"
	"github.com/hashicorp/terraform/helper/schema"
	"fmt"
)

func ExpandDvdDrives(d *schema.ResourceData) ([]vmDvdDrive, error) {
	expandedDvdDrives := make([]vmDvdDrive, 0)

	if v, ok := d.GetOk("dvd_drives"); ok {
		dvdDrives := v.([]interface{})
		for _, dvdDrive := range dvdDrives {
			dvdDrive, ok := dvdDrive.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] network_adaptors should be a Hash - was '%+v'", dvdDrive)
			}

			expandedDvdDrive := vmDvdDrive{
				ControllerNumber:   dvdDrive["controller_number"].(int),
				ControllerLocation: dvdDrive["controller_location"].(int),
				Path:               dvdDrive["path"].(string),
				ResourcePoolName:   dvdDrive["resource_pool_name"].(string),
			}

			expandedDvdDrives = append(expandedDvdDrives, expandedDvdDrive)
		}
	}

	if len(expandedDvdDrives) > 0 {
		return expandedDvdDrives, nil
	}

	return nil, nil
}

func FlattenDvdDrives(dvdDrives *[]vmDvdDrive) []interface{} {
	flattenedDvdDrives := make([]interface{}, 0)

	if dvdDrives != nil {
		for _, dvdDrive := range *dvdDrives {
			flattenedDvdDrive := make(map[string]interface{})
			flattenedDvdDrive["controller_number"] = dvdDrive.ControllerNumber
			flattenedDvdDrive["controller_location"] = dvdDrive.ControllerLocation
			flattenedDvdDrive["path"] = dvdDrive.Path
			flattenedDvdDrive["resource_pool_name"] = dvdDrive.ResourcePoolName
			flattenedDvdDrives = append(flattenedDvdDrives, flattenedDvdDrive)
		}
	}

	return flattenedDvdDrives
}

type vmDvdDrive struct {
	VMName             string
	ControllerNumber   int
	ControllerLocation int
	Path               string
	//AllowUnverifiedPaths bool no way of checking if its turned on so always turn on
	ResourcePoolName string
}

type createVMDvdDriveArgs struct {
	VmDvdDriveJson string
}

var createVMDvdDriveTemplate = template.Must(template.New("CreateVMDvdDrive").Parse(`
$ErrorActionPreference = 'Stop'
Get-Vm | Out-Null
$vmDvdDrive = '{{.VmDvdDriveJson}}' | ConvertFrom-Json

$NewVmDvdDriveArgs = @{
	VMName=$vmDvdDrive.VmName
	ControllerNumber=$vmDvdDrive.ControllerNumber
	ControllerLocation=$vmDvdDrive.ControllerLocation
	Path=$vmDvdDrive.Path
	ResourcePoolName=$vmDvdDrive.ResourcePoolName
	AllowUnverifiedPaths=$true
}

Add-VmDvdDrive @NewVmDvdDriveArgs
`))

func (c *HypervClient) CreateVMDvdDrive(
	vmName string,
	controllerNumber int,
	controllerLocation int,
	path string,
	resourcePoolName string,
) (err error) {

	vmDvdDriveJson, err := json.Marshal(vmDvdDrive{
		VMName:             vmName,
		ControllerNumber:   controllerNumber,
		ControllerLocation: controllerLocation,
		Path:               path,
		ResourcePoolName:   resourcePoolName,
	})

	err = c.runFireAndForgetScript(createVMDvdDriveTemplate, createVMDvdDriveArgs{
		VmDvdDriveJson: string(vmDvdDriveJson),
	})

	return err
}

type getVMDvdDrivesArgs struct {
	VMName string
}

var getVMDvdDrivesTemplate = template.Must(template.New("GetVMDvdDrives").Parse(`
$ErrorActionPreference = 'Stop'
$vmDvdDrivesObject = Get-VMDvdDrive -VMName '{{.VMName}}' | %{ @{
	ControllerNumber=$_.ControllerNumber;
	ControllerLocation=$_.ControllerLocation;
	Path=$_.Path;
	#ControllerType=$_.ControllerType; not able to set it
	#DvdMediaType=$_.DvdMediaType; not able to set it
	ResourcePoolName=$_.PoolName;
}}

if ($vmDvdDrivesObject) {
	$vmDvdDrives = ConvertTo-Json -InputObject $vmDvdDrivesObject
	$vmDvdDrives
} else {
	"[]"
}
`))

func (c *HypervClient) GetVMDvdDrives(vmName string) (result []vmDvdDrive, err error) {
	result = make([]vmDvdDrive, 0)

	err = c.runScriptWithResult(getVMDvdDrivesTemplate, getVMDvdDrivesArgs{
		VMName: vmName,
	}, result)

	return result, err
}

type updateVMDvdDriveArgs struct {
	VMName             string
	ControllerNumber   int
	ControllerLocation int
	VmDvdDriveJson     string
}

var updateVMDvdDriveTemplate = template.Must(template.New("UpdateVMDvdDrive").Parse(`
$ErrorActionPreference = 'Stop'
Get-Vm | Out-Null
$vmDvdDrive = '{{.VmDvdDriveJson}}' | ConvertFrom-Json

$vmDvdDrivesObject = @(Get-VMDvdDrive -VMName '{{.VMName}}' -ControllerLocation {{.ControllerLocation}} -ControllerNumber {{.ControllerNumber}} )

if (!$vmDvdDrivesObject){
	throw "VM dvd drive does not exist - {{.ControllerLocation}} {{.ControllerNumber}}"
}

$SetVmDvdDriveArgs = @{}
$SetVmDvdDriveArgs.VMName=$vmDvdDrivesObject.VMName
$SetVmDvdDriveArgs.ControllerLocation=$vmDvdDrivesObject.ControllerLocation
$SetVmDvdDriveArgs.ControllerNumber=$vmDvdDrivesObject.ControllerNumber
$SetVmDvdDriveArgs.ToControllerLocation=$vmDvdDrive.ControllerLocation
$SetVmDvdDriveArgs.ToControllerNumber=$vmDvdDrive.ControllerNumber
$SetVmDvdDriveArgs.ResourcePoolName=$vmDvdDrive.ResourcePoolName
$SetVmDvdDriveArgs.Path=$vmDvdDrive.Path
$SetVmDvdDriveArgs.AllowUnverifiedPaths=$true

Set-VMDvdDrive @SetVmDvdDriveArgs

`))

func (c *HypervClient) UpdateVMDvdDrive(
	vmName string,
	controllerNumber int,
	controllerLocation int,
	toControllerNumber int,
	toControllerLocation int,
	path string,
	resourcePoolName string,
) (err error) {

	vmDvdDriveJson, err := json.Marshal(vmDvdDrive{
		VMName:             vmName,
		ControllerNumber:   toControllerNumber,
		ControllerLocation: toControllerLocation,
		Path:               path,
		ResourcePoolName:   resourcePoolName,
	})

	err = c.runFireAndForgetScript(updateVMDvdDriveTemplate, updateVMDvdDriveArgs{
		VMName:             vmName,
		ControllerNumber:   controllerNumber,
		ControllerLocation: controllerLocation,
		VmDvdDriveJson:     string(vmDvdDriveJson),
	})

	return err
}

type deleteVMDvdDriveArgs struct {
	VMName             string
	ControllerNumber   int
	ControllerLocation int
}

var deleteVMDvdDriveTemplate = template.Must(template.New("DeleteVMDvdDrive").Parse(`
$ErrorActionPreference = 'Stop'

@(Get-VMDvdDrive -VMName '{{.VMName}}' -ControllerNumber {{.ControllerNumber}} -ControllerLocation {{.ControllerLocation}}) | Remove-VMDvdDrive -Force
`))

func (c *HypervClient) DeleteVMDvdDrive(vmName string, controllerNumber int, controllerLocation int) (err error) {
	err = c.runFireAndForgetScript(deleteVMDvdDriveTemplate, deleteVMDvdDriveArgs{
		VMName:             vmName,
		ControllerNumber:   controllerNumber,
		ControllerLocation: controllerLocation,
	})

	return err
}

func (c *HypervClient) CreateOrUpdateVMDvdDrives(vmName string, dvdDrives []vmDvdDrive) (err error) {
	currentDvdDrives, err := c.GetVMDvdDrives(vmName)
	if err != nil {
		return err
	}

	currentDvdDrivesLength := len(currentDvdDrives)
	desiredDvdDrivesLength := len(dvdDrives)

	for i := currentDvdDrivesLength - 1; i > desiredDvdDrivesLength-1; i-- {
		currentDvdDrive := currentDvdDrives[i]
		err = c.DeleteVMDvdDrive(vmName, currentDvdDrive.ControllerNumber, currentDvdDrive.ControllerLocation)
		if err != nil {
			return err
		}
	}

	if currentDvdDrivesLength > desiredDvdDrivesLength {
		currentDvdDrivesLength = desiredDvdDrivesLength
	}

	for i := 0; i <= currentDvdDrivesLength-1; i++ {
		currentDvdDrive := currentDvdDrives[i]
		dvdDrive := dvdDrives[i]

		err = c.UpdateVMDvdDrive(
			vmName,
			currentDvdDrive.ControllerNumber,
			currentDvdDrive.ControllerLocation,
			dvdDrive.ControllerNumber,
			dvdDrive.ControllerLocation,
			dvdDrive.Path,
			dvdDrive.ResourcePoolName,
		)
		if err != nil {
			return err
		}
	}

	for i := currentDvdDrivesLength - 1 + 1; i <= desiredDvdDrivesLength-1; i++ {
		dvdDrive := dvdDrives[i]
		err = c.CreateVMDvdDrive(
			vmName,
			dvdDrive.ControllerNumber,
			dvdDrive.ControllerLocation,
			dvdDrive.Path,
			dvdDrive.ResourcePoolName,
		)

		if err != nil {
			return err
		}
	}

	return nil
}
