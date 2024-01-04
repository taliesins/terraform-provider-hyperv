package api

import (
	"encoding/json"
	"testing"
)

func TestSerializeIsoImage(t *testing.T) {
	isoImageJson, err := json.Marshal(IsoImage{
		SourceIsoFilePath:              "",
		SourceIsoFilePathHash:          "",
		SourceBootFilePath:             "boot.img",
		SourceBootFilePathHash:         "654321",
		SourceZipFilePath:              "bootstrap.zip",
		SourceZipFilePathHash:          "123456",
		DestinationIsoFilePath:         "bootstrap.iso",
		DestinationZipFilePath:         "",
		DestinationBootFilePath:        "",
		FileSystem:                     IsoFileSystemType_ISO9660_or_Joliet,
		Media:                          IsoMediaType_DVDPLUSRW_DUALLAYER,
		VolumeName:                     "test",
		ResolveDestinationIsoFilePath:  "$env:temp\\bootstrap.iso",
		ResolveDestinationZipFilePath:  "$env:temp\\bootstrap.zip",
		ResolveDestinationBootFilePath: "$env:temp\\boot.img",
	})

	if err != nil {
		t.Errorf("Unable to deserialize iso image: %s", err.Error())
	}

	isoImageJsonString := string(isoImageJson)

	if isoImageJsonString == "" {
		t.Errorf("Unable to deserialize iso image: %s", err.Error())
	}
}

func TestDeserializeIsoImage(t *testing.T) {
	var isoImageJson = `
{
	"SourceIsoFilePath":"",
	"SourceIsoFilePathHash":"",
	"SourceZipFilePath":"bootstrap.zip",
	"SourceZipFilePathHash":"123456",
	"SourceBootFilePath":"boot.img",
	"SourceBootFilePathHash":"654321",
	"DestinationIsoFilePath":"bootstrap.iso",
	"DestinationZipFilePath":"",
	"DestinationBootFilePath":"",
	"Media":13,
	"FileSystem":3,
	"VolumeName":"test",
	"ResolveDestinationIsoFilePath":"$env:temp\\bootstrap.iso",
	"ResolveDestinationZipFilePath":"$env:temp\\bootstrap.zip",
	"ResolveDestinationBootFilePath":"$env:temp\\boot.img"
}
`

	var isoImage IsoImage
	err := json.Unmarshal([]byte(isoImageJson), &isoImage)

	if err != nil {
		t.Errorf("Unable to deserialize iso image: %s", err.Error())
	}
}
