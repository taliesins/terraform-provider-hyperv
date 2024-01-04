package api

import (
	"bytes"
	"context"
	"encoding/json"
	"strconv"
	"strings"
)

type IsoMediaType int

const (
	IsoMediaType_UNKNOWN             IsoMediaType = 0
	IsoMediaType_CDROM               IsoMediaType = 0x1
	IsoMediaType_CDR                 IsoMediaType = 0x2
	IsoMediaType_CDRW                IsoMediaType = 0x3
	IsoMediaType_DVDROM              IsoMediaType = 0x4
	IsoMediaType_DVDRAM              IsoMediaType = 0x5
	IsoMediaType_DVDPLUSR            IsoMediaType = 0x6
	IsoMediaType_DVDPLUSRW           IsoMediaType = 0x7
	IsoMediaType_DVDPLUSR_DUALLAYER  IsoMediaType = 0x8
	IsoMediaType_DVDDASHR            IsoMediaType = 0x9
	IsoMediaType_DVDDASHRW           IsoMediaType = 0xa
	IsoMediaType_DVDDASHR_DUALLAYER  IsoMediaType = 0xb
	IsoMediaType_DISK                IsoMediaType = 0xc
	IsoMediaType_DVDPLUSRW_DUALLAYER IsoMediaType = 0xd
	IsoMediaType_HDDVDROM            IsoMediaType = 0xe
	IsoMediaType_HDDVDR              IsoMediaType = 0xf
	IsoMediaType_HDDVDRAM            IsoMediaType = 0x10
	IsoMediaType_BDROM               IsoMediaType = 0x11
	IsoMediaType_BDR                 IsoMediaType = 0x12
	IsoMediaType_BDRE                IsoMediaType = 0x13
	//IsoMediaType_MAX                 IsoMediaType = 0x13
)

var IsoMediaType_name = map[IsoMediaType]string{
	IsoMediaType_UNKNOWN:             "unknown",
	IsoMediaType_CDROM:               "cdrom",
	IsoMediaType_CDR:                 "cdr",
	IsoMediaType_CDRW:                "cdrw",
	IsoMediaType_DVDROM:              "dvdrom",
	IsoMediaType_DVDRAM:              "dvdram",
	IsoMediaType_DVDPLUSR:            "dvdplusr",
	IsoMediaType_DVDPLUSRW:           "dvdplusrw",
	IsoMediaType_DVDPLUSR_DUALLAYER:  "dvdplusr_duallayer",
	IsoMediaType_DVDDASHR:            "dvddashr",
	IsoMediaType_DVDDASHRW:           "dvddashrw",
	IsoMediaType_DVDDASHR_DUALLAYER:  "dvddashr_duallayer",
	IsoMediaType_DISK:                "disk",
	IsoMediaType_DVDPLUSRW_DUALLAYER: "dvdplusrw_duallayer",
	IsoMediaType_HDDVDROM:            "hddvdrom",
	IsoMediaType_HDDVDR:              "hddvdr",
	IsoMediaType_HDDVDRAM:            "hddvdram",
	IsoMediaType_BDROM:               "bdrom",
	IsoMediaType_BDR:                 "bdr",
	IsoMediaType_BDRE:                "bdre",
	//IsoMediaType_MAX: "max",
}

var IsoMediaType_value = map[string]IsoMediaType{
	"unknown":             IsoMediaType_UNKNOWN,
	"cdrom":               IsoMediaType_CDROM,
	"cdr":                 IsoMediaType_CDR,
	"cdrw":                IsoMediaType_CDRW,
	"dvdrom":              IsoMediaType_DVDROM,
	"dvdram":              IsoMediaType_DVDRAM,
	"dvdplusr":            IsoMediaType_DVDPLUSR,
	"dvdplusrw":           IsoMediaType_DVDPLUSRW,
	"dvdplusr_duallayer":  IsoMediaType_DVDPLUSR_DUALLAYER,
	"dvddashr":            IsoMediaType_DVDDASHR,
	"dvddashrw":           IsoMediaType_DVDDASHRW,
	"dvddashr_duallayer":  IsoMediaType_DVDDASHR_DUALLAYER,
	"disk":                IsoMediaType_DISK,
	"dvdplusrw_duallayer": IsoMediaType_DVDPLUSRW_DUALLAYER,
	"hddvdrom":            IsoMediaType_HDDVDROM,
	"hddvdr":              IsoMediaType_HDDVDR,
	"hddvdram":            IsoMediaType_HDDVDRAM,
	"bdrom":               IsoMediaType_BDROM,
	"bdr":                 IsoMediaType_BDR,
	"bdre":                IsoMediaType_BDRE,
	//"max": IsoMediaType_MAX,
}

func (x IsoMediaType) String() string {
	return IsoMediaType_name[x]
}

func ToIsoMediaType(x string) IsoMediaType {
	if integerValue, err := strconv.Atoi(x); err == nil {
		return IsoMediaType(integerValue)
	}

	return IsoMediaType_value[strings.ToLower(x)]
}

func (d *IsoMediaType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(d.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (d *IsoMediaType) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		var i int
		err2 := json.Unmarshal(b, &i)
		if err2 == nil {
			*d = IsoMediaType(i)
			return nil
		}

		return err
	}
	*d = ToIsoMediaType(s)
	return nil
}

type IsoFileSystemType int

const (
	IsoFileSystemType_None              IsoFileSystemType = 0
	IsoFileSystemType_ISO9660           IsoFileSystemType = 0x1
	IsoFileSystemType_Joliet            IsoFileSystemType = 0x2
	IsoFileSystemType_ISO9660_or_Joliet IsoFileSystemType = 0x3
	IsoFileSystemType_UDF               IsoFileSystemType = 4
	IsoFileSystemType_Joliet_or_UDF     IsoFileSystemType = 0x6
	IsoFileSystemType_ALL               IsoFileSystemType = 0x7
	IsoFileSystemType_Unknown           IsoFileSystemType = 0x40000000
)

var IsoFileSystemType_name = map[IsoFileSystemType]string{
	IsoFileSystemType_None:              "none",
	IsoFileSystemType_ISO9660:           "iso9660",
	IsoFileSystemType_Joliet:            "joliet",
	IsoFileSystemType_ISO9660_or_Joliet: "iso9660|joliet",
	IsoFileSystemType_UDF:               "udf",
	IsoFileSystemType_Joliet_or_UDF:     "joliet|udf",
	IsoFileSystemType_ALL:               "iso9660|joliet|udf",
	IsoFileSystemType_Unknown:           "unknown",
}

var IsoFileSystemType_value = map[string]IsoFileSystemType{
	"none":               IsoFileSystemType_None,
	"iso9660":            IsoFileSystemType_ISO9660,
	"joliet":             IsoFileSystemType_Joliet,
	"iso9660|joliet":     IsoFileSystemType_ISO9660_or_Joliet,
	"udf":                IsoFileSystemType_UDF,
	"joliet|udf":         IsoFileSystemType_Joliet_or_UDF,
	"iso9660|joliet|udf": IsoFileSystemType_ALL,
	"unknown":            IsoFileSystemType_Unknown,
}

func (x IsoFileSystemType) String() string {
	return IsoFileSystemType_name[x]
}

func ToIsoFileSystemType(x string) IsoFileSystemType {
	if integerValue, err := strconv.Atoi(x); err == nil {
		return IsoFileSystemType(integerValue)
	}

	return IsoFileSystemType_value[strings.ToLower(x)]
}

func (d *IsoFileSystemType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(d.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (d *IsoFileSystemType) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		var i int
		err2 := json.Unmarshal(b, &i)
		if err2 == nil {
			*d = IsoFileSystemType(i)
			return nil
		}

		return err
	}
	*d = ToIsoFileSystemType(s)
	return nil
}

type IsoImage struct {
	SourceIsoFilePath              string
	SourceIsoFilePathHash          string
	SourceZipFilePath              string
	SourceZipFilePathHash          string
	SourceBootFilePath             string
	SourceBootFilePathHash         string
	DestinationIsoFilePath         string
	DestinationZipFilePath         string
	DestinationBootFilePath        string
	Media                          IsoMediaType
	FileSystem                     IsoFileSystemType
	VolumeName                     string
	ResolveDestinationIsoFilePath  string
	ResolveDestinationZipFilePath  string
	ResolveDestinationBootFilePath string
}

type HypervIsoImageClient interface {
	RemoteFileExists(ctx context.Context, path string) (exists bool, err error)
	RemoteFileDelete(ctx context.Context, path string) (err error)
	RemoteFileUpload(ctx context.Context, filePath string, remoteFilePath string) (err error)

	CreateOrUpdateIsoImage(ctx context.Context, sourceIsoFilePath string, sourceIsoFilePathHash string, sourceZipFilePath string, sourceZipFilePathHash string, sourceBootFilePath string, sourceBootFilePathHash string, destinationIsoFilePath string, destinationZipFilePath string, destinationBootFilePath string, media IsoMediaType, fileSystem IsoFileSystemType, volumeName string, resolveDestinationIsoFilePath string, resolveDestinationZipFilePath string, resolveDestinationBootFilePath string) (err error)
	GetIsoImage(ctx context.Context, resolveDestinationIsoFilePath string) (result IsoImage, err error)
}
