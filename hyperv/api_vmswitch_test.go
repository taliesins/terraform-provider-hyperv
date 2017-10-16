package hyperv

import (
	"testing"
	"encoding/json"
)

func TestSerializeVmSwitch(t *testing.T) {
	vmSwitchJson, err := json.Marshal(vmSwitch{
		Name:"test",
		Notes:"test notes",
		AllowManagementOS:true,
		EmbeddedTeamingEnabled:true,
		IovEnabled:true,
		PacketDirectEnabled:false,
		BandwidthReservationMode:VMSwitchBandwidthMode_Weight,
		SwitchType:VMSwitchType_Internal,
		NetAdapterNames:[]string{"wan", "lan"},
		DefaultQueueVrssEnabled:true,
		DefaultQueueVmmqQueuePairs:0,
	})

	if err != nil {
		t.Errorf("Unable to deserialize vm switch: %s", err.Error())
	}

	vmSwitchJsonString := string(vmSwitchJson)

	if vmSwitchJsonString == "" {
		t.Errorf("Unable to deserialize vm switch: %s", err.Error())
	}
}

func TestDeserializeVmSwitch(t *testing.T){
	var vmSwitchJson = `
{
    "BandwidthReservationMode":  2,
    "NetAdapterInterfaceDescriptions":  [
                                            "Dell Wireless 1830 802.11ac"
                                        ],
    "Notes":  "test notes",
    "AllowManagementOS":  true,
    "Name":  "test",
    "SwitchType":  2,
    "IovEnabled":  false,
    "EmbeddedTeamingEnabled":  false,
    "PacketDirectEnabled":  false
}
`

	var vmSwitch vmSwitch
	err := json.Unmarshal([]byte(vmSwitchJson), &vmSwitch)

	if err != nil {
		t.Errorf("Unable to deserialize vm switch: %s", err.Error())
	}
}