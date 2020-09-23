package powershell

import (
	"bytes"
	"testing"
)

func TestEscapeQuotesOfCommandLineTemplate(t *testing.T) {
	command := `& { if (Test-Path variable:global:ProgressPreference){$ProgressPreference='SilentlyContinue'};;&"C:/Windows/Temp/Test.ps1";exit $LastExitCode }`

	var executePowershellFromCommandLineTemplateRendered bytes.Buffer
	err := executePowershellFromCommandLineTemplate.Execute(&executePowershellFromCommandLineTemplateRendered, executePowershellFromCommandLineTemplateOptions{
		Powershell: command,
	})

	if err != nil {
		t.Errorf("Unable to render command line template: %s", err.Error())
	}

	commandLine := executePowershellFromCommandLineTemplateRendered.String()

	if commandLine != `powershell -NoProfile -ExecutionPolicy Bypass "& { if (Test-Path variable:global:ProgressPreference){$ProgressPreference='SilentlyContinue'};;&\"C:/Windows/Temp/Test.ps1\";exit $LastExitCode }"` {
		t.Errorf("Command line template output not as expected: %s", err.Error())
	}
}
