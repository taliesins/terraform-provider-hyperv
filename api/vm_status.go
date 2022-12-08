package api

import (
	"bytes"
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type VmState int

const (
	VmState_Other              VmState = 1
	VmState_Running            VmState = 2
	VmState_Off                VmState = 3
	VmState_Stopping           VmState = 4
	VmState_Saved              VmState = 6
	VmState_Paused             VmState = 9
	VmState_Starting           VmState = 10
	VmState_Reset              VmState = 11
	VmState_Saving             VmState = 32773
	VmState_Pausing            VmState = 32776
	VmState_Resuming           VmState = 32777
	VmState_FastSaved          VmState = 32779
	VmState_FastSaving         VmState = 32780
	VmState_ForceShutdown      VmState = 32781
	VmState_ForceReboot        VmState = 32782
	VmState_RunningCritical    VmState = 32783
	VmState_OffCritical        VmState = 32784
	VmState_StoppingCritical   VmState = 32785
	VmState_SavedCritical      VmState = 32786
	VmState_PausedCritical     VmState = 32787
	VmState_StartingCritical   VmState = 32788
	VmState_ResetCritical      VmState = 32789
	VmState_SavingCritical     VmState = 32790
	VmState_PausingCritical    VmState = 32791
	VmState_ResumingCritical   VmState = 32792
	VmState_FastSavedCritical  VmState = 32793
	VmState_FastSavingCritical VmState = 32794
)

var VmState_name = map[VmState]string{
	VmState_Other:              "Other",
	VmState_Running:            "Running",
	VmState_Off:                "Off",
	VmState_Stopping:           "Stopping",
	VmState_Saved:              "Saved",
	VmState_Paused:             "Paused",
	VmState_Starting:           "Starting",
	VmState_Reset:              "Reset",
	VmState_Saving:             "Saving",
	VmState_Pausing:            "Pausing",
	VmState_Resuming:           "Resuming",
	VmState_FastSaved:          "FastSaved",
	VmState_FastSaving:         "FastSaving",
	VmState_ForceShutdown:      "ForceShutdown",
	VmState_ForceReboot:        "ForceReboot",
	VmState_RunningCritical:    "RunningCritical",
	VmState_OffCritical:        "OffCritical",
	VmState_StoppingCritical:   "StoppingCritical",
	VmState_SavedCritical:      "SavedCritical",
	VmState_PausedCritical:     "PausedCritical",
	VmState_StartingCritical:   "StartingCritical",
	VmState_ResetCritical:      "ResetCritical",
	VmState_SavingCritical:     "SavingCritical",
	VmState_PausingCritical:    "PausingCritical",
	VmState_ResumingCritical:   "ResumingCritical",
	VmState_FastSavedCritical:  "FastSavedCritical",
	VmState_FastSavingCritical: "FastSavingCritical",
}

var VmState_SettableValue = map[string]VmState{
	"running": VmState_Running,
	"off":     VmState_Off,
}

var VmState_value = map[string]VmState{
	"other":              VmState_Other,
	"running":            VmState_Running,
	"off":                VmState_Off,
	"stopping":           VmState_Stopping,
	"saved":              VmState_Saved,
	"paused":             VmState_Paused,
	"starting":           VmState_Starting,
	"reset":              VmState_Reset,
	"saving":             VmState_Saving,
	"pausing":            VmState_Pausing,
	"resuming":           VmState_Resuming,
	"fastsaved":          VmState_FastSaved,
	"fastsaving":         VmState_FastSaving,
	"forceshutdown":      VmState_ForceShutdown,
	"forcereboot":        VmState_ForceReboot,
	"runningcritical":    VmState_RunningCritical,
	"offcritical":        VmState_OffCritical,
	"stoppingcritical":   VmState_StoppingCritical,
	"savedcritical":      VmState_SavedCritical,
	"pausedcritical":     VmState_PausedCritical,
	"startingcritical":   VmState_StartingCritical,
	"resetcritical":      VmState_ResetCritical,
	"savingcritical":     VmState_SavingCritical,
	"pausingcritical":    VmState_PausingCritical,
	"resumingcritical":   VmState_ResumingCritical,
	"fastsavedcritical":  VmState_FastSavedCritical,
	"fastsavingcritical": VmState_FastSavingCritical,
}

func (x VmState) String() string {
	return VmState_name[x]
}

func ToVmState(x string) VmState {
	if integerValue, err := strconv.Atoi(x); err == nil {
		return VmState(integerValue)
	}
	return VmState_value[strings.ToLower(x)]
}

func (d *VmState) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(d.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (d *VmState) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		var i int
		err2 := json.Unmarshal(b, &i)
		if err2 == nil {
			*d = VmState(i)
			return nil
		}

		return err
	}
	*d = ToVmState(s)
	return nil
}

type VmStatus struct {
	State VmState
}

func ExpandVmStateWaitForState(d *schema.ResourceData) (uint32, uint32, error) {
	waitForIpsTimeout := uint32((d.Get("wait_for_state_timeout")).(int))
	waitForIpsPollPeriod := uint32((d.Get("wait_for_state_poll_period")).(int))

	return waitForIpsTimeout, waitForIpsPollPeriod, nil
}

type HypervVmStatusClient interface {
	GetVmStatus(ctx context.Context, vmName string) (result VmStatus, err error)
	UpdateVmStatus(
		ctx context.Context,
		vmName string,
		timeout uint32,
		pollPeriod uint32,
		state VmState,
	) (err error)
}
