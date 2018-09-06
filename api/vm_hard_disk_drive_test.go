package api

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestSerializeVmHardDiskDrive(t *testing.T) {
	vmHardDiskDriveJson, err := json.Marshal(vmHardDiskDrive{
		Path:  `C:\data\VirtualMachines\web_server\Virtual Hard Disks\MobyLinuxVM.vhdx`,
		OverrideCacheAttributes:  0,
		ControllerLocation:  0,
		ControllerNumber:  0,
		DiskNumber:  4294967295,
		QosPolicyId:  "00000000-0000-0000-0000-000000000000",
		MinimumIops:  0,
		SupportPersistentReservations:  false,
		ControllerType:  0,
		ResourcePoolName:  "Primordial",
		MaximumIops:  0,
	})

	if err != nil {
		t.Errorf("Unable to deserialize vmHardDiskDrive: %s", err.Error())
	}

	vmHardDiskDriveJsonString := string(vmHardDiskDriveJson)

	if vmHardDiskDriveJsonString == "" {
		t.Errorf("Unable to deserialize vmHardDiskDrive: %s", err.Error())
	}

	if !strings.Contains(vmHardDiskDriveJsonString,`C:\\data\\VirtualMachines\\web_server\\Virtual Hard Disks\\MobyLinuxVM.vhdx`) {
		t.Errorf("Path does not match")
	}
}

func TestDeserializeVmHardDiskDrive(t *testing.T){
	var vmHardDiskDriveJson = `
{
    "Path":  "C:\\data\\VirtualMachines\\web_server\\Virtual Hard Disks\\MobyLinuxVM.vhdx",
    "OverrideCacheAttributes":  0,
    "ControllerLocation":  0,
    "ControllerNumber":  0,
    "DiskNumber":  4294967295,
    "QosPolicyId":  "00000000-0000-0000-0000-000000000000",
    "MinimumIops":  0,
    "SupportPersistentReservations":  false,
    "ControllerType":  0,
    "ResourcePoolName":  "Primordial",
    "MaximumIops":  0
}
`

	var vmHardDiskDrive vmHardDiskDrive
	err := json.Unmarshal([]byte(vmHardDiskDriveJson), &vmHardDiskDrive)
	if err != nil {
		t.Errorf("Unable to deserialize vmHardDiskDrive: %s", err.Error())
	}

	if vmHardDiskDrive.Path != `C:\data\VirtualMachines\web_server\Virtual Hard Disks\MobyLinuxVM.vhdx` {
		t.Errorf("Path does not match")
	}
}

func TestDeserializeVmHardDiskDrives(t *testing.T){
	var vmHardDiskDrivesJson = `
[
{
    "Path":  "C:\\data\\VirtualMachines\\web_server\\Virtual Hard Disks\\MobyLinuxVM.vhdx",
    "OverrideCacheAttributes":  0,
    "ControllerLocation":  0,
    "ControllerNumber":  0,
    "DiskNumber":  4294967295,
    "QosPolicyId":  "00000000-0000-0000-0000-000000000000",
    "MinimumIops":  0,
    "SupportPersistentReservations":  false,
    "ControllerType":  0,
    "ResourcePoolName":  "Primordial",
    "MaximumIops":  0
},
{
    "Path":  "C:\\data\\VirtualMachines\\web_server\\Virtual Hard Disks\\MobyLinuxVM2.vhdx",
    "OverrideCacheAttributes":  0,
    "ControllerLocation":  0,
    "ControllerNumber":  0,
    "DiskNumber":  4294967295,
    "QosPolicyId":  "00000000-0000-0000-0000-000000000000",
    "MinimumIops":  0,
    "SupportPersistentReservations":  false,
    "ControllerType":  0,
    "ResourcePoolName":  "Primordial",
    "MaximumIops":  0
}
]
`

	var vmHardDiskDrives = make([]vmHardDiskDrive, 0)

	err := json.Unmarshal([]byte(vmHardDiskDrivesJson), &vmHardDiskDrives)
	if err != nil {
		t.Errorf("Unable to deserialize vmHardDiskDrives: %s", err.Error())
	}

	if len(vmHardDiskDrives) != 2 {
		t.Errorf("Array does not have 2 elements")
	}

	if vmHardDiskDrives[0].Path != `C:\data\VirtualMachines\web_server\Virtual Hard Disks\MobyLinuxVM.vhdx` {
		t.Errorf("Path does not match")
	}
}
