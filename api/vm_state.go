package api

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
	"text/template"
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

type vmState struct {
	State VmState
}

type getVMStateArgs struct {
	VMName string
}

var getVMStateTemplate = template.Must(template.New("GetVMState").Parse(`
$ErrorActionPreference = 'Stop'
$vmName = '{{.VMName}}'

$vmStateObject = Get-VM | ?{$_.Name -eq $vmName } | %{ @{
	State=$_.State;
}}

if ($vmStateObject) {
	$vmState = ConvertTo-Json -InputObject $vmStateObject
	$vmState
} else {
	"{}"
}
`))

func (c *HypervClient) GetVMState(vmName string) (result vmState, err error) {
	err = c.runScriptWithResult(getVMStateTemplate, getVMStateArgs{
		VMName: vmName,
	}, &result)

	return result, err
}

type updateVMStateArgs struct {
	VMName        string
	Timeout       uint32
	RetryInterval uint32
	VmStateJson   string
}

var updateVMStateTemplate = template.Must(template.New("UpdateVMState").Parse(`
$ErrorActionPreference = 'Stop'

function Test-VmStateRequiresManualIntervention($state){
    $states = @([Microsoft.HyperV.PowerShell.VMState]::Other, 
        [Microsoft.HyperV.PowerShell.VMState]::RunningCritical
        [Microsoft.HyperV.PowerShell.VMState]::OffCritical, 
        [Microsoft.HyperV.PowerShell.VMState]::StoppingCritical,
        [Microsoft.HyperV.PowerShell.VMState]::SavedCritical,
        [Microsoft.HyperV.PowerShell.VMState]::PausedCritical,
        [Microsoft.HyperV.PowerShell.VMState]::StartingCritical,
        [Microsoft.HyperV.PowerShell.VMState]::ResetCritical,
        [Microsoft.HyperV.PowerShell.VMState]::SavingCritical,
        [Microsoft.HyperV.PowerShell.VMState]::PausingCritical,
        [Microsoft.HyperV.PowerShell.VMState]::ResumingCritical,
        [Microsoft.HyperV.PowerShell.VMState]::FastSavedCritical,
        [Microsoft.HyperV.PowerShell.VMState]::FastSavingCritical
        )
	   
    return $states -contains $state 
}

function Test-IsInNotInFinalTransitionState($State){
    $states = @([Microsoft.HyperV.PowerShell.VMState]::Other,
		[Microsoft.HyperV.PowerShell.VMState]::Stopping,
		[Microsoft.HyperV.PowerShell.VMState]::Saved,
		[Microsoft.HyperV.PowerShell.VMState]::Starting,
		[Microsoft.HyperV.PowerShell.VMState]::Reset,
		[Microsoft.HyperV.PowerShell.VMState]::Saving,
		[Microsoft.HyperV.PowerShell.VMState]::Pausing,
		[Microsoft.HyperV.PowerShell.VMState]::Resuming,
		[Microsoft.HyperV.PowerShell.VMState]::FastSaved,
		[Microsoft.HyperV.PowerShell.VMState]::FastSaving,
		[Microsoft.HyperV.PowerShell.VMState]::ForceShutdown,
		[Microsoft.HyperV.PowerShell.VMState]::ForceReboot,
        [Microsoft.HyperV.PowerShell.VMState]::StoppingCritical,
        [Microsoft.HyperV.PowerShell.VMState]::SavedCritical,
        [Microsoft.HyperV.PowerShell.VMState]::StartingCritical,
        [Microsoft.HyperV.PowerShell.VMState]::ResetCritical,
        [Microsoft.HyperV.PowerShell.VMState]::SavingCritical,
        [Microsoft.HyperV.PowerShell.VMState]::PausingCritical,
        [Microsoft.HyperV.PowerShell.VMState]::ResumingCritical,
        [Microsoft.HyperV.PowerShell.VMState]::FastSavedCritical,
        [Microsoft.HyperV.PowerShell.VMState]::FastSavingCritical
        )
	   
    return $states -contains $State 
}

function Wait-IsInInFinalTransitionState($Name, $Timeout, $RetryInterval){
	$timer = [Diagnostics.Stopwatch]::StartNew()
	while (($timer.Elapsed.TotalSeconds -lt $Timeout) -and (Test-IsInNotInFinalTransitionState (Get-VM -name $Name).state)) { 
		Start-Sleep -Seconds $RetryInterval
	}
	$timer.Stop()

	if ($timer.Elapsed.TotalSeconds -gt $Timeout) {
		throw 'Timeout while waiting for vm $($Name) to reach final transition state'
	} 
}

Get-Vm | Out-Null
$vm = '{{.VmStateJson}}' | ConvertFrom-Json
$vmName = '{{.VMName}}'
$state = [Microsoft.HyperV.PowerShell.VMState]$vm.State
$vmObject = Get-VM | ?{$_.Name -eq $vmName}
$timeout = {{.Timeout}}
$retryInterval = {{.RetryInterval}}

if (!$vmObject){
	throw "VM does not exist - $($vmName)"
}

if (Test-VmStateRequiresManualIntervention -State $vmObject.State) {
	throw "VM $($vmName) requires manual intervention as it is in state $($vmObject.State)"
}

Wait-IsInInFinalTransitionState -Name $vmName -Timeout $timeout -RetryInterval $retryInterval

if ($vmObject -eq $state) {
} elseif ($state -eq [Microsoft.HyperV.PowerShell.VMState]::Running) {
	if ($vmObject.State -eq [Microsoft.HyperV.PowerShell.VMState]::Off) { 
		Start-VM -Name $vmName
	} elseif ($vmObject.State -eq [Microsoft.HyperV.PowerShell.VMState]::Off) { 
		Resume-VM -Name $vmName
	} else {
		throw "Unable to change VM $($vmName) state $($vmObject.State) to Running state"
	}
} elseif ($state -eq [Microsoft.HyperV.PowerShell.VMState]::Off) { 
	if ($vmObject.State -eq [Microsoft.HyperV.PowerShell.VMState]::Running -or $vmObject.State -eq [Microsoft.HyperV.PowerShell.VMState]::Paused) { 
		Stop-VM -Name $vmName -force
	} else {
		throw "Unable to change VM $($vmName) state $($vmObject.State) to Off state"
	}
} elseif ($state -eq [Microsoft.HyperV.PowerShell.VMState]::Paused) {
	if ($vmObject.State -eq [Microsoft.HyperV.PowerShell.VMState]::Running) { 
		Suspend-VM -Name $vmName
	} else {
		throw "Unable to change VM $($vmName) state $($vmObject.State) to Paused state"
	}	
}

`))

func (c *HypervClient) UpdateVMState(
	vmName string,
	timeout uint32,
	retryInterval uint32,
	state VmState,
) (err error) {

	vmStateJson, err := json.Marshal(vmState{
		State: state,
	})

	err = c.runFireAndForgetScript(updateVMStateTemplate, updateVMStateArgs{
		VMName:        vmName,
		Timeout:       timeout,
		RetryInterval: retryInterval,
		VmStateJson:   string(vmStateJson),
	})

	return err
}

