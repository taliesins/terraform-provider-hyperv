package hyperv_winrm

import (
	"encoding/json"
	"github.com/taliesins/terraform-provider-hyperv/api"
	"text/template"
)

type getVmStatusArgs struct {
	VmName string
}

var getVmStatusTemplate = template.Must(template.New("GetVmStatus").Parse(`
$ErrorActionPreference = 'Stop'
$vmName = '{{.VmName}}'

$vmStateObject = Get-VM -Name "$($vmName)*" | ?{$_.Name -eq $vmName } | %{ @{
	State=$_.State;
}}

if ($vmStateObject) {
	$vmState = ConvertTo-Json -InputObject $vmStateObject
	$vmState
} else {
	"{}"
}
`))

func (c *ClientConfig) GetVmStatus(vmName string) (result api.VmStatus, err error) {
	err = c.WinRmClient.RunScriptWithResult(getVmStatusTemplate, getVmStatusArgs{
		VmName: vmName,
	}, &result)

	return result, err
}

type updateVmStatusArgs struct {
	VmName       string
	Timeout      uint32
	PollPeriod   uint32
	VmStatusJson string
}

var updateVmStatusTemplate = template.Must(template.New("UpdateVmStatus").Parse(`
$ErrorActionPreference = 'Stop'

function Test-VmStateRequiresManualIntervention($state){
    $states = @([Microsoft.HyperV.PowerShell.VMState]::Other, 
        [Microsoft.HyperV.PowerShell.VMState]::RunningCritical,
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

function Test-IsNotInFinalTransitionState($State){
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

function Wait-IsInFinalTransitionState($Name, $Timeout, $PollPeriod){
	$timer = [Diagnostics.Stopwatch]::StartNew()
	while (($timer.Elapsed.TotalSeconds -lt $Timeout) -and (Test-IsNotInFinalTransitionState (Get-VM -name $Name).state)) { 
		Start-Sleep -Seconds $PollPeriod
	}
	$timer.Stop()

	if ($timer.Elapsed.TotalSeconds -gt $Timeout) {
		throw 'Timeout while waiting for vm $($Name) to reach final transition state'
	} 
}

Import-Module Hyper-V
$vm = '{{.VmStatusJson}}' | ConvertFrom-Json
$vmName = '{{.VmName}}'
$state = [Microsoft.HyperV.PowerShell.VMState]$vm.State
$vmObject = Get-VM -Name "$($vmName)*" | ?{$_.Name -eq $vmName}
$timeout = {{.Timeout}}
$pollPeriod = {{.PollPeriod}}

if (!$vmObject){
	throw "VM does not exist - $($vmName)"
}

if ($vmObject.State -ne $state) {
    if (Test-VmStateRequiresManualIntervention -State $vmObject.State) {
        throw "VM $($vmName) requires manual intervention as it is in state $($vmObject.State)"
    }

    Wait-IsInFinalTransitionState -Name $vmName -Timeout $timeout -PollPeriod $pollPeriod

    $vmObject = Get-VM -Name "$($vmName)*" | ?{$_.Name -eq $vmName}

    if ($vmObject.State -eq $state) {
    } elseif ($state -eq [Microsoft.HyperV.PowerShell.VMState]::Running) {
        if ($vmObject.State -eq [Microsoft.HyperV.PowerShell.VMState]::Off) { 
            Start-VM -Name $vmName
            Start-Sleep -Seconds $pollPeriod
            Wait-IsInFinalTransitionState -Name $vmName -Timeout $timeout -PollPeriod $pollPeriod
        } elseif ($vmObject.State -eq [Microsoft.HyperV.PowerShell.VMState]::Off) { 
            Resume-VM -Name $vmName
            Start-Sleep -Seconds $pollPeriod
            Wait-IsInFinalTransitionState -Name $vmName -Timeout $timeout -PollPeriod $pollPeriod
        } else {
            throw "Unable to change VM $($vmName) state $($vmObject.State) to Running state"
        }
    } elseif ($state -eq [Microsoft.HyperV.PowerShell.VMState]::Off) { 
        if ($vmObject.State -eq [Microsoft.HyperV.PowerShell.VMState]::Running -or $vmObject.State -eq [Microsoft.HyperV.PowerShell.VMState]::Paused) { 
            Stop-VM -Name $vmName -force
            Start-Sleep -Seconds $pollPeriod
            Wait-IsInFinalTransitionState -Name $vmName -Timeout $timeout -PollPeriod $pollPeriod
        } else {
            throw "Unable to change VM $($vmName) state $($vmObject.State) to Off state"
        }
    } elseif ($state -eq [Microsoft.HyperV.PowerShell.VMState]::Paused) {
        if ($vmObject.State -eq [Microsoft.HyperV.PowerShell.VMState]::Running) { 
            Suspend-VM -Name $vmName
            Start-Sleep -Seconds $pollPeriod
            Wait-IsInFinalTransitionState -Name $vmName -Timeout $timeout -PollPeriod $pollPeriod
        } else {
            throw "Unable to change VM $($vmName) state $($vmObject.State) to Paused state"
        }	
    }
}
`))

func (c *ClientConfig) UpdateVmStatus(
	vmName string,
	timeout uint32,
	pollPeriod uint32,
	state api.VmState,
) (err error) {

	vmStatusJson, err := json.Marshal(api.VmStatus{
		State: state,
	})

	if err != nil {
		return err
	}

	err = c.WinRmClient.RunFireAndForgetScript(updateVmStatusTemplate, updateVmStatusArgs{
		VmName:       vmName,
		Timeout:      timeout,
		PollPeriod:   pollPeriod,
		VmStatusJson: string(vmStatusJson),
	})

	return err
}
