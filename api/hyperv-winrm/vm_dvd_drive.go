package hyperv_winrm

import (
	"encoding/json"
	"text/template"

	"github.com/taliesins/terraform-provider-hyperv/api"
)

type createVmDvdDriveArgs struct {
	VmDvdDriveJson string
}

var createVmDvdDriveTemplate = template.Must(template.New("CreateVmDvdDrive").Parse(`
$ErrorActionPreference = 'Stop'
Import-Module Hyper-V
$vmDvdDrive = '{{.VmDvdDriveJson}}' | ConvertFrom-Json
if (!$vmDvdDrive.Path){
	$vmDvdDrive.Path = $null
}
$NewVmDvdDriveArgs = @{
	VmName=$vmDvdDrive.VmName
	ControllerNumber=$vmDvdDrive.ControllerNumber
	ControllerLocation=$vmDvdDrive.ControllerLocation
	Path=$vmDvdDrive.Path
	ResourcePoolName=$vmDvdDrive.ResourcePoolName
	AllowUnverifiedPaths=$true
}

Add-VmDvdDrive @NewVmDvdDriveArgs
`))

func (c *ClientConfig) CreateVmDvdDrive(
	vmName string,
	controllerNumber int,
	controllerLocation int,
	path string,
	resourcePoolName string,
) (err error) {

	vmDvdDriveJson, err := json.Marshal(api.VmDvdDrive{
		VmName:             vmName,
		ControllerNumber:   controllerNumber,
		ControllerLocation: controllerLocation,
		Path:               path,
		ResourcePoolName:   resourcePoolName,
	})

	if err != nil {
		return err
	}

	err = c.WinRmClient.RunFireAndForgetScript(createVmDvdDriveTemplate, createVmDvdDriveArgs{
		VmDvdDriveJson: string(vmDvdDriveJson),
	})

	return err
}

type getVmDvdDrivesArgs struct {
	VmName string
}

var getVmDvdDrivesTemplate = template.Must(template.New("GetVmDvdDrives").Parse(`
$ErrorActionPreference = 'Stop'
$vmDvdDrivesObject = @(Get-VMDvdDrive -VmName '{{.VmName}}' | %{ @{
	ControllerNumber=$_.ControllerNumber;
	ControllerLocation=$_.ControllerLocation;
	Path=$_.Path;
	#ControllerType=$_.ControllerType; not able to set it
	#DvdMediaType=$_.DvdMediaType; not able to set it
	ResourcePoolName=$_.PoolName;
}})

if ($vmDvdDrivesObject) {
	$vmDvdDrives = ConvertTo-Json -InputObject $vmDvdDrivesObject
	$vmDvdDrives
} else {
	"[]"
}
`))

func (c *ClientConfig) GetVmDvdDrives(vmName string) (result []api.VmDvdDrive, err error) {
	result = make([]api.VmDvdDrive, 0)

	err = c.WinRmClient.RunScriptWithResult(getVmDvdDrivesTemplate, getVmDvdDrivesArgs{
		VmName: vmName,
	}, &result)

	return result, err
}

type updateVmDvdDriveArgs struct {
	VmName             string
	ControllerNumber   int
	ControllerLocation int
	VmDvdDriveJson     string
}

var updateVmDvdDriveTemplate = template.Must(template.New("UpdateVmDvdDrive").Parse(`
$ErrorActionPreference = 'Stop'
Import-Module Hyper-V
$vmDvdDrive = '{{.VmDvdDriveJson}}' | ConvertFrom-Json

$vmDvdDrivesObject = @(Get-VMDvdDrive -VmName '{{.VmName}}' -ControllerLocation {{.ControllerLocation}} -ControllerNumber {{.ControllerNumber}} )

if (!$vmDvdDrivesObject){
	throw "VM dvd drive does not exist - {{.ControllerLocation}} {{.ControllerNumber}}"
}

$SetVmDvdDriveArgs = @{}
$SetVmDvdDriveArgs.VmName=$vmDvdDrivesObject.VmName
$SetVmDvdDriveArgs.ControllerLocation=$vmDvdDrivesObject.ControllerLocation
$SetVmDvdDriveArgs.ControllerNumber=$vmDvdDrivesObject.ControllerNumber
$SetVmDvdDriveArgs.ToControllerLocation=$vmDvdDrive.ControllerLocation
$SetVmDvdDriveArgs.ToControllerNumber=$vmDvdDrive.ControllerNumber
$SetVmDvdDriveArgs.ResourcePoolName=$vmDvdDrive.ResourcePoolName
$SetVmDvdDriveArgs.Path=$vmDvdDrive.Path
$SetVmDvdDriveArgs.AllowUnverifiedPaths=$true

if (!$SetVmDvdDriveArgs.Path){
	$SetVmDvdDriveArgs.Path = $null
}

Set-VMDvdDrive @SetVmDvdDriveArgs

`))

func (c *ClientConfig) UpdateVmDvdDrive(
	vmName string,
	controllerNumber int,
	controllerLocation int,
	toControllerNumber int,
	toControllerLocation int,
	path string,
	resourcePoolName string,
) (err error) {

	vmDvdDriveJson, err := json.Marshal(api.VmDvdDrive{
		VmName:             vmName,
		ControllerNumber:   toControllerNumber,
		ControllerLocation: toControllerLocation,
		Path:               path,
		ResourcePoolName:   resourcePoolName,
	})

	if err != nil {
		return err
	}

	err = c.WinRmClient.RunFireAndForgetScript(updateVmDvdDriveTemplate, updateVmDvdDriveArgs{
		VmName:             vmName,
		ControllerNumber:   controllerNumber,
		ControllerLocation: controllerLocation,
		VmDvdDriveJson:     string(vmDvdDriveJson),
	})

	return err
}

type deleteVmDvdDriveArgs struct {
	VmName             string
	ControllerNumber   int
	ControllerLocation int
}

var deleteVmDvdDriveTemplate = template.Must(template.New("DeleteVmDvdDrive").Parse(`
$ErrorActionPreference = 'Stop'

@(Get-VMDvdDrive -VmName '{{.VmName}}' -ControllerNumber {{.ControllerNumber}} -ControllerLocation {{.ControllerLocation}}) | Remove-VMDvdDrive -Force
`))

func (c *ClientConfig) DeleteVmDvdDrive(vmName string, controllerNumber int, controllerLocation int) (err error) {
	err = c.WinRmClient.RunFireAndForgetScript(deleteVmDvdDriveTemplate, deleteVmDvdDriveArgs{
		VmName:             vmName,
		ControllerNumber:   controllerNumber,
		ControllerLocation: controllerLocation,
	})

	return err
}

func (c *ClientConfig) CreateOrUpdateVmDvdDrives(vmName string, dvdDrives []api.VmDvdDrive) (err error) {
	currentDvdDrives, err := c.GetVmDvdDrives(vmName)
	if err != nil {
		return err
	}

	currentDvdDrivesLength := len(currentDvdDrives)
	desiredDvdDrivesLength := len(dvdDrives)

	for i := currentDvdDrivesLength - 1; i > desiredDvdDrivesLength-1; i-- {
		currentDvdDrive := currentDvdDrives[i]
		err = c.DeleteVmDvdDrive(vmName, currentDvdDrive.ControllerNumber, currentDvdDrive.ControllerLocation)
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

		err = c.UpdateVmDvdDrive(
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
		err = c.CreateVmDvdDrive(
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
