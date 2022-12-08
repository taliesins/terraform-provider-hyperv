package api

import (
	"bytes"
	"context"
	"encoding/json"
	"strconv"
	"strings"
)

type VMSwitchBandwidthMode int

const (
	VMSwitchBandwidthMode_Default  VMSwitchBandwidthMode = 0
	VMSwitchBandwidthMode_Weight   VMSwitchBandwidthMode = 1
	VMSwitchBandwidthMode_Absolute VMSwitchBandwidthMode = 2
	VMSwitchBandwidthMode_None     VMSwitchBandwidthMode = 3
)

var VMSwitchBandwidthMode_name = map[VMSwitchBandwidthMode]string{
	VMSwitchBandwidthMode_Default:  "Default",
	VMSwitchBandwidthMode_Weight:   "Weight",
	VMSwitchBandwidthMode_Absolute: "Absolute",
	VMSwitchBandwidthMode_None:     "None",
}

var VMSwitchBandwidthMode_value = map[string]VMSwitchBandwidthMode{
	"default":  VMSwitchBandwidthMode_Default,
	"weight":   VMSwitchBandwidthMode_Weight,
	"absolute": VMSwitchBandwidthMode_Absolute,
	"none":     VMSwitchBandwidthMode_None,
}

func (x VMSwitchBandwidthMode) String() string {
	return VMSwitchBandwidthMode_name[x]
}

func ToVMSwitchBandwidthMode(x string) VMSwitchBandwidthMode {
	if integerValue, err := strconv.Atoi(x); err == nil {
		return VMSwitchBandwidthMode(integerValue)
	}

	return VMSwitchBandwidthMode_value[strings.ToLower(x)]
}

func (d *VMSwitchBandwidthMode) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(d.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (d *VMSwitchBandwidthMode) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		var i int
		err2 := json.Unmarshal(b, &i)
		if err2 == nil {
			*d = VMSwitchBandwidthMode(i)
			return nil
		}

		return err
	}
	*d = ToVMSwitchBandwidthMode(s)
	return nil
}

type VMSwitchType int

const (
	VMSwitchType_Private  VMSwitchType = 0
	VMSwitchType_Internal VMSwitchType = 1
	VMSwitchType_External VMSwitchType = 2
)

var VMSwitchType_name = map[VMSwitchType]string{
	VMSwitchType_Private:  "Private",
	VMSwitchType_Internal: "Internal",
	VMSwitchType_External: "External",
}

var VMSwitchType_value = map[string]VMSwitchType{
	"private":  VMSwitchType_Private,
	"internal": VMSwitchType_Internal,
	"external": VMSwitchType_External,
}

func (x VMSwitchType) String() string {
	return VMSwitchType_name[x]
}

func ToVMSwitchType(x string) VMSwitchType {
	if integerValue, err := strconv.Atoi(x); err == nil {
		return VMSwitchType(integerValue)
	}

	return VMSwitchType_value[strings.ToLower(x)]
}

func (d *VMSwitchType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(d.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (d *VMSwitchType) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		var i int
		err2 := json.Unmarshal(b, &i)
		if err2 == nil {
			*d = VMSwitchType(i)
			return nil
		}

		return err
	}
	*d = ToVMSwitchType(s)
	return nil
}

type VmSwitchExists struct {
	Exists bool
}

type VmSwitch struct {
	Name                                string
	Notes                               string
	AllowManagementOS                   bool
	EmbeddedTeamingEnabled              bool
	IovEnabled                          bool
	PacketDirectEnabled                 bool
	BandwidthReservationMode            VMSwitchBandwidthMode
	SwitchType                          VMSwitchType
	NetAdapterNames                     []string
	DefaultFlowMinimumBandwidthAbsolute int64
	DefaultFlowMinimumBandwidthWeight   int64
	DefaultQueueVmmqEnabled             bool
	DefaultQueueVmmqQueuePairs          int32
	DefaultQueueVrssEnabled             bool
}

type HypervVmSwitchClient interface {
	VMSwitchExists(ctx context.Context, name string) (result VmSwitchExists, err error)
	CreateVMSwitch(
		ctx context.Context,
		name string,
		notes string,
		allowManagementOS bool,
		embeddedTeamingEnabled bool,
		iovEnabled bool,
		packetDirectEnabled bool,
		bandwidthReservationMode VMSwitchBandwidthMode,
		switchType VMSwitchType,
		netAdapterNames []string,
		defaultFlowMinimumBandwidthAbsolute int64,
		defaultFlowMinimumBandwidthWeight int64,
		defaultQueueVmmqEnabled bool,
		defaultQueueVmmqQueuePairs int32,
		defaultQueueVrssEnabled bool,
	) (err error)
	GetVMSwitch(ctx context.Context, name string) (result VmSwitch, err error)
	UpdateVMSwitch(
		ctx context.Context,
		name string,
		notes string,
		allowManagementOS bool,
		// embeddedTeamingEnabled bool,
		// iovEnabled bool,
		// packetDirectEnabled bool,
		// bandwidthReservationMode VMSwitchBandwidthMode,
		switchType VMSwitchType,
		netAdapterNames []string,
		defaultFlowMinimumBandwidthAbsolute int64,
		defaultFlowMinimumBandwidthWeight int64,
		defaultQueueVmmqEnabled bool,
		defaultQueueVmmqQueuePairs int32,
		defaultQueueVrssEnabled bool,
	) (err error)
	DeleteVMSwitch(ctx context.Context, name string) (err error)
}
