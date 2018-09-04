package api

import (
	"testing"
	"encoding/json"
)

func TestSerializeVmNetworkAdapter(t *testing.T) {
	vmNetworkAdapterJson, err := json.Marshal(vmNetworkAdapter{
		Name:"test",
	})

	if err != nil {
		t.Errorf("Unable to deserialize vmNetworkAdapter: %s", err.Error())
	}

	vmNetworkAdapterJsonString := string(vmNetworkAdapterJson)

	if vmNetworkAdapterJsonString == "" {
		t.Errorf("Unable to deserialize vmNetworkAdapter: %s", err.Error())
	}
}

func TestDeserializeVmNetworkAdapter(t *testing.T){
	var vmNetworkAdapterJson = `
{
    "Name":  "TestMachine"
}
`

	var vmNetworkAdapter vmNetworkAdapter
	err := json.Unmarshal([]byte(vmNetworkAdapterJson), &vmNetworkAdapter)
	if err != nil {
		t.Errorf("Unable to deserialize vmNetworkAdapter: %s", err.Error())
	}
}

