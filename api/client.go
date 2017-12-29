package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"text/template"

	"github.com/masterzen/winrm"
	"github.com/taliesins/terraform-provider-hyperv/powershell"
)

type HypervClient struct {
	WinrmClient      *winrm.Client
	ElevatedUser     string
	ElevatedPassword string
	Vars             string
}

func (c *HypervClient) runFireAndForgetScript(script *template.Template, args interface{}) error {
	var scriptRendered bytes.Buffer
	err := script.Execute(&scriptRendered, args)

	if err != nil {
		return err
	}

	command := string(scriptRendered.Bytes())

	log.Printf("\nRunning fire and forget script:\n%s\n", command)

	_, _, _, err = powershell.RunPowershell(c.WinrmClient, c.ElevatedUser, c.ElevatedPassword, c.Vars, command)

	if err != nil {
		return err
	}

	return nil
}

func (c *HypervClient) runScriptWithResult(script *template.Template, args interface{}, result interface{}) (err error) {
	var scriptRendered bytes.Buffer
	err = script.Execute(&scriptRendered, args)

	if err != nil {
		return err
	}

	command := string(scriptRendered.Bytes())

	log.Printf("\nRunning script with result:\n%s\n", command)

	_, stdout, _, err := powershell.RunPowershell(c.WinrmClient, c.ElevatedUser, c.ElevatedPassword, c.Vars, command)

	if err != nil {
		return err
	}

	stdout = strings.TrimSpace(stdout)

	err = json.Unmarshal([]byte(stdout), &result)
	if err != nil {
		return fmt.Errorf("%s\n%s", err, stdout)
	}

	return nil
}
