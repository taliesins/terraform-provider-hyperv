package api

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
	"text/template"
)

type VMVhdType int

const (
	VMVhdType_Unknown      VMVhdType = 0
	VMVhdType_Fixed        VMVhdType = 2
	VMVhdType_Dynamic      VMVhdType = 3
	VMVhdType_Differencing VMVhdType = 4
)

var VMVhdType_name = map[VMVhdType]string{
	VMVhdType_Unknown:      "Unknown",
	VMVhdType_Fixed:        "Fixed",
	VMVhdType_Dynamic:      "Dynamic",
	VMVhdType_Differencing: "Differencing",
}

var VMVhdType_value = map[string]VMVhdType{
	"unknown":      VMVhdType_Unknown,
	"fixed":        VMVhdType_Fixed,
	"dynamic":      VMVhdType_Dynamic,
	"differencing": VMVhdType_Differencing,
}

func (x VMVhdType) String() string {
	return VMVhdType_name[x]
}

func ToVMVhdType(x string) VMVhdType {
	if integerValue, err := strconv.Atoi(x); err == nil {
		return VMVhdType(integerValue)
	}

	return VMVhdType_value[strings.ToLower(x)]
}

func (d *VMVhdType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(d.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (d *VMVhdType) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		var i int
		err2 := json.Unmarshal(b, &i)
		if err2 == nil {
			*d = VMVhdType(i)
			return nil
		}

		return err
	}
	*d = ToVMVhdType(s)
	return nil
}

type VMVhdFormat int

const (
	VMVhdFormat_Unknown VMVhdFormat = 0
	VMVhdFormat_VHD     VMVhdFormat = 2 //extension ".vhd"
	VMVhdFormat_VHDX    VMVhdFormat = 3 //extension ".vhdx"
	VMVhdFormat_VHDSet  VMVhdFormat = 4 //extension ".vhds"
)

var VMVhdFormat_name = map[VMVhdFormat]string{
	VMVhdFormat_Unknown: "Unknown",
	VMVhdFormat_VHD:     "VHD",
	VMVhdFormat_VHDX:    "VHDX",
	VMVhdFormat_VHDSet:  "VHDSet",
}

var VMVhdFormat_value = map[string]VMVhdFormat{
	"unknown": VMVhdFormat_Unknown,
	"vhd":     VMVhdFormat_VHD,
	"vhdx":    VMVhdFormat_VHDX,
	"vhdset":  VMVhdFormat_VHDSet,
}

func (x VMVhdFormat) String() string {
	return VMVhdFormat_name[x]
}

func ToVMVhdFormat(x string) VMVhdFormat {
	if integerValue, err := strconv.Atoi(x); err == nil {
		return VMVhdFormat(integerValue)
	}

	return VMVhdFormat_value[strings.ToLower(x)]
}

func (d *VMVhdFormat) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(d.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (d *VMVhdFormat) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		var i int
		err2 := json.Unmarshal(b, &i)
		if err2 == nil {
			*d = VMVhdFormat(i)
			return nil
		}

		return err
	}
	*d = ToVMVhdFormat(s)
	return nil
}

type vmVhd struct {
	Path                    string
	BlockSize               uint32
	LogicalSectorSize       uint32
	PhysicalSectorSize      uint32
	ParentPath              string
	FileSize                uint64
	Size                    uint64
	MinimumSize             uint64
	Attached                bool
	DiskNumber              int
	Number                  int
	FragmentationPercentage int
	Alignment               int
	DiskIdentifier          string
	VhdType                 VMVhdType
	VhdFormat               VMVhdFormat
}

type createOrUpdateVMVhdArgs struct {
	SourcePath string
	SourceUrl  string
	SourceDisk int
	VhdJson   string
}

var createOrUpdateVMVhdTemplate = template.Must(template.New("CreateOrUpdateVMVhd").Parse(`
$ErrorActionPreference = 'Stop'

Get-VM | Out-Null
$sourcePath='{{.SourcePath}}'
$sourceUrl='{{.SourceUrl}}'
$sourceDisk={{.SourceDisk}}
$vhd = '{{.VhdJson}}' | ConvertFrom-Json
$vhdType = [Microsoft.Vhd.PowerShell.VhdType]$vhd.VhdType

if (!(Test-Path -Path $vhd.Path)) {
	if ($sourcePath) {
		Copy-Item $sourcePath $vhd.Path
	} elseif ($sourceUrl) {
		(New-Object System.Net.WebClient).DownloadFile($sourceUrl, $vhd.Path)
	} else {
		$NewVhdArgs = @{}
		$NewVhdArgs.Path=$vhd.Path

		if ($sourceDisk) {
			$NewVhdArgs.SourceDisk=$sourceDisk
		} elseif ($vhdType -eq [Microsoft.Vhd.PowerShell.VhdType]::Differencing) {
			$NewVhdArgs.Differencing=$true
			$NewVhdArgs.ParentPath=$vhd.ParentPath
		} else {
			if ($vhdType -eq [Microsoft.Vhd.PowerShell.VhdType]::Dynamic) {
				$NewVhdArgs.Dynamic=$true
			} elseif ($vhdType -eq [Microsoft.Vhd.PowerShell.VhdType]::Fixed) {
				$NewVhdArgs.Fixed=$true
			}

			if ($vhd.Size -gt 0) {
				$NewVhdArgs.SizeBytes=$vhd.Size
			} 

			if ($vhd.BlockSize -gt 0) {
				$NewVhdArgs.BlockSizeBytes=$vhd.BlockSize
			} 

			if ($vhd.LogicalSectorSize -gt 0) {
				$NewVhdArgs.LogicalSectorSizeBytes=$vhd.LogicalSectorSize
			} 

			if ($vhd.PhysicalSectorSize -gt 0) {
				$NewVhdArgs.PhysicalSectorSizeBytes=$vhd.PhysicalSectorSize
			} 
		}

		New-VHD @NewVhdArgs
	}
}
`))

func (c *HypervClient) CreateOrUpdateVMVhd(path string, sourcePath string, sourceUrl string, sourceDisk int, vhdType VMVhdType, parentPath string, size uint64, blockSize uint32, logicalSectorSize uint32, physicalSectorSize uint32) (err error) {
	vhdJson, err := json.Marshal(vmVhd{
		Path:               path,
		VhdType:            vhdType,
		ParentPath:         parentPath,
		Size:               size,
		BlockSize:          blockSize,
		LogicalSectorSize:  logicalSectorSize,
		PhysicalSectorSize: physicalSectorSize,
	})

	err = c.runFireAndForgetScript(createOrUpdateVMVhdTemplate, createOrUpdateVMVhdArgs{
		SourcePath: sourcePath,
		SourceUrl:  sourceUrl,
		SourceDisk: sourceDisk,
		VhdJson:   string(vhdJson),
	})

	return err
}

type resizeVMVhdArgs struct {
	Path string
	Size uint64
}

var resizeVMVhdTemplate = template.Must(template.New("ResizeVMVhd").Parse(`
$ErrorActionPreference = 'Stop'
Resize-VHD –Path '{{.Path}}' –SizeBytes {{.Size}}
`))

func (c *HypervClient) ResizeVMVhd(path string, size uint64) (err error) {
	err = c.runFireAndForgetScript(resizeVMVhdTemplate, resizeVMVhdArgs{
		Path: path,
		Size: size,
	})

	return err
}

type getVMVhdArgs struct {
	Path string
}

var getVMVhdTemplate = template.Must(template.New("GetVMVhd").Parse(`
$ErrorActionPreference = 'Stop'
$path='{{.Path}}'

$vmVhdObject =  Get-VHD -path $path | %{ @{
	Path=$_.Path;
	BlockSize=$_.BlockSize;
	LogicalSectorSize=$_.LogicalSectorSize;
	PhysicalSectorSize=$_.PhysicalSectorSize;
	ParentPath=$_.ParentPath;
	FileSize=$_.FileSize;
	Size=$_.Size;
	MinimumSize=$_.MinimumSize;
	Attached=$_.Attached;
	DiskNumber=$_.DiskNumber;
	Number=$_.Number;
	FragmentationPercentage=$_.FragmentationPercentage;
	Alignment=$_.Alignment;
	DiskIdentifier=$_.DiskIdentifier;
	VhdType=$_.VhdType;
	VhdFormat=$_.VhdFormat;
}}

if ($vmVhdObject){
	$vmVhd = ConvertTo-Json -InputObject $vmVhdObject
	$vmVhd
} else {
	"{}"
}
`))

func (c *HypervClient) GetVMVhd(path string) (result vmVhd, err error) {
	err = c.runScriptWithResult(getVMVhdTemplate, getVMVhdArgs{
		Path: path,
	}, &result)

	return result, err
}

type deleteVMVhdArgs struct {
	Path string
}

var deleteVMVhdTemplate = template.Must(template.New("DeleteVMVhd").Parse(`
$ErrorActionPreference = 'Stop'
if (Test-Path -Path '{{.Path}}') {
	Remove-Item -Path '{{.Path}}' -Force
}
`))

func (c *HypervClient) DeleteVMVhd(path string) (err error) {
	err = c.runFireAndForgetScript(deleteVMVhdTemplate, deleteVMVhdArgs{
		Path: path,
	})

	return err
}
