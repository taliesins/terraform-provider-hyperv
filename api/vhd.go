package api

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
)

type VhdType int

const (
	VhdType_Unknown      VhdType = 0
	VhdType_Fixed        VhdType = 2
	VhdType_Dynamic      VhdType = 3
	VhdType_Differencing VhdType = 4
)

var VhdType_name = map[VhdType]string{
	VhdType_Unknown:      "Unknown",
	VhdType_Fixed:        "Fixed",
	VhdType_Dynamic:      "Dynamic",
	VhdType_Differencing: "Differencing",
}

var VhdType_value = map[string]VhdType{
	"unknown":      VhdType_Unknown,
	"fixed":        VhdType_Fixed,
	"dynamic":      VhdType_Dynamic,
	"differencing": VhdType_Differencing,
}

func (x VhdType) String() string {
	return VhdType_name[x]
}

func ToVhdType(x string) VhdType {
	if integerValue, err := strconv.Atoi(x); err == nil {
		return VhdType(integerValue)
	}

	return VhdType_value[strings.ToLower(x)]
}

func (d *VhdType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(d.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (d *VhdType) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		var i int
		err2 := json.Unmarshal(b, &i)
		if err2 == nil {
			*d = VhdType(i)
			return nil
		}

		return err
	}
	*d = ToVhdType(s)
	return nil
}

type VhdFormat int

const (
	VhdFormat_Unknown VhdFormat = 0
	VhdFormat_VHD     VhdFormat = 2 //extension ".Vhd"
	VhdFormat_VHDX    VhdFormat = 3 //extension ".vhdx"
	VhdFormat_VHDSet  VhdFormat = 4 //extension ".vhds"
)

var VhdFormat_name = map[VhdFormat]string{
	VhdFormat_Unknown: "Unknown",
	VhdFormat_VHD:     "VHD",
	VhdFormat_VHDX:    "VHDX",
	VhdFormat_VHDSet:  "VHDSet",
}

var VhdFormat_value = map[string]VhdFormat{
	"unknown": VhdFormat_Unknown,
	"Vhd":     VhdFormat_VHD,
	"vhdx":    VhdFormat_VHDX,
	"vhdset":  VhdFormat_VHDSet,
}

func (x VhdFormat) String() string {
	return VhdFormat_name[x]
}

func ToVhdFormat(x string) VhdFormat {
	if integerValue, err := strconv.Atoi(x); err == nil {
		return VhdFormat(integerValue)
	}

	return VhdFormat_value[strings.ToLower(x)]
}

func (d *VhdFormat) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(d.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (d *VhdFormat) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		var i int
		err2 := json.Unmarshal(b, &i)
		if err2 == nil {
			*d = VhdFormat(i)
			return nil
		}

		return err
	}
	*d = ToVhdFormat(s)
	return nil
}

type Vhd struct {
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
	VhdType                 VhdType
	VhdFormat               VhdFormat
}

type HypervVhdClient interface {
	CreateOrUpdateVhd(path string, source string, sourceVm string, sourceDisk int, vhdType VhdType, parentPath string, size uint64, blockSize uint32, logicalSectorSize uint32, physicalSectorSize uint32) (err error)
	ResizeVhd(path string, size uint64) (err error)
	GetVhd(path string) (result Vhd, err error)
	DeleteVhd(path string) (err error)
}
