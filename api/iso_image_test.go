package api

import (
	"encoding/json"
	"testing"
)

func TestSerializeIsoImage(t *testing.T) {
	isoImageJson, err := json.Marshal(IsoImage{
		VolumeName:          "test",
		FileSystem:          IsoFileSystemType_ISO9660_or_Joliet,
		Media:               IsoMediaType_DVDPLUSRW_DUALLAYER,
		DestinationFilePath: "bootstrap.iso",
		SourceBootFilePath:  "image.",
		SourceDirectoryPath: "bootstrap",
		ExcludeList: []string{
			"do_not_include.json",
		},
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
	"SourceDirectoryPath":"bootstrap",
	"SourceBootFilePath":"image.",
	"DestinationFilePath":"bootstrap.iso",
	"ExcludeList":[
		"do_not_include.json"
	],
	"Media":13,
	"FileSystem":3,
	"VolumeName":"test"
}
`

	var isoImage IsoImage
	err := json.Unmarshal([]byte(isoImageJson), &isoImage)

	if err != nil {
		t.Errorf("Unable to deserialize iso image: %s", err.Error())
	}
}
