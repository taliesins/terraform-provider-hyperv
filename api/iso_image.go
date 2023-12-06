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
	IsoMediaType_Unknown             IsoMediaType = 0
	IsoMediaType_CDROM               IsoMediaType = 1
	IsoMediaType_CDR                 IsoMediaType = 2
	IsoMediaType_CDRW                IsoMediaType = 3
	IsoMediaType_DVDROM              IsoMediaType = 4
	IsoMediaType_DVDRAM              IsoMediaType = 5
	IsoMediaType_DVDPLUSR            IsoMediaType = 6
	IsoMediaType_DVDPLUSRW           IsoMediaType = 7
	IsoMediaType_DVDPLUSR_DUALLAYER  IsoMediaType = 8
	IsoMediaType_DVDDASHR            IsoMediaType = 9
	IsoMediaType_DVDDASHRW           IsoMediaType = 10
	IsoMediaType_DVDDASHR_DUALLAYER  IsoMediaType = 11
	IsoMediaType_DISK                IsoMediaType = 12
	IsoMediaType_DVDPLUSRW_DUALLAYER IsoMediaType = 13
	IsoMediaType_HDDVDROM            IsoMediaType = 14
	IsoMediaType_HDDVDR              IsoMediaType = 15
	IsoMediaType_HDDVDRAM            IsoMediaType = 16
	IsoMediaType_BDROM               IsoMediaType = 17
	IsoMediaType_BDR                 IsoMediaType = 18
	IsoMediaType_BDRE                IsoMediaType = 19
)

var IsoMediaType_name = map[IsoMediaType]string{
	IsoMediaType_Unknown:             "UNKNOWN",
	IsoMediaType_CDROM:               "CDROM",
	IsoMediaType_CDR:                 "CDR",
	IsoMediaType_CDRW:                "CDRW",
	IsoMediaType_DVDROM:              "DVDROM",
	IsoMediaType_DVDRAM:              "DVDRAM",
	IsoMediaType_DVDPLUSR:            "DVDPLUSR",
	IsoMediaType_DVDPLUSRW:           "DVDPLUSRW",
	IsoMediaType_DVDPLUSR_DUALLAYER:  "DVDPLUSR_DUALLAYER",
	IsoMediaType_DVDDASHR:            "DVDDASHR",
	IsoMediaType_DVDDASHRW:           "DVDDASHRW",
	IsoMediaType_DVDDASHR_DUALLAYER:  "DVDDASHR_DUALLAYER",
	IsoMediaType_DISK:                "DISK",
	IsoMediaType_DVDPLUSRW_DUALLAYER: "DVDPLUSRW_DUALLAYER",
	IsoMediaType_HDDVDROM:            "HDDVDROM",
	IsoMediaType_HDDVDR:              "HDDVDR",
	IsoMediaType_HDDVDRAM:            "HDDVDRAM",
	IsoMediaType_BDROM:               "BDROM",
	IsoMediaType_BDR:                 "BDR",
	IsoMediaType_BDRE:                "BDRE",
}

var IsoMediaType_value = map[string]IsoMediaType{
	"UNKNOWN":             IsoMediaType_Unknown,
	"CDROM":               IsoMediaType_CDROM,
	"CDR":                 IsoMediaType_CDR,
	"CDRW":                IsoMediaType_CDRW,
	"DVDROM":              IsoMediaType_DVDROM,
	"DVDRAM":              IsoMediaType_DVDRAM,
	"DVDPLUSR":            IsoMediaType_DVDPLUSR,
	"DVDPLUSRW":           IsoMediaType_DVDPLUSRW,
	"DVDPLUSR_DUALLAYER":  IsoMediaType_DVDPLUSR_DUALLAYER,
	"DVDDASHR":            IsoMediaType_DVDDASHR,
	"DVDDASHRW":           IsoMediaType_DVDDASHRW,
	"DVDDASHR_DUALLAYER":  IsoMediaType_DVDDASHR_DUALLAYER,
	"DISK":                IsoMediaType_DISK,
	"DVDPLUSRW_DUALLAYER": IsoMediaType_DVDPLUSRW_DUALLAYER,
	"HDDVDROM":            IsoMediaType_HDDVDROM,
	"HDDVDR":              IsoMediaType_HDDVDR,
	"HDDVDRAM":            IsoMediaType_HDDVDRAM,
	"BDROM":               IsoMediaType_BDROM,
	"BDR":                 IsoMediaType_BDR,
	"BDRE":                IsoMediaType_BDRE,
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

type RemoteIsoClient interface {
	//IsoImageExists(ctx context.Context, path string) (result VhdExists, err error)
	CreateOrUpdateIsoImage(ctx context.Context, sourceDirectoryPath string, sourceBootFilePath string, destinationIsoPath string, excludeList []string, media IsoMediaType, title string) (err error)
	//ResizeIsoImage(ctx context.Context, path string, size uint64) (err error)
	//GetIsoImage(ctx context.Context, path string) (result Vhd, err error)
	DeleteIsoImage(ctx context.Context, path string) (err error)
}
