package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"text/template"

	pool "github.com/jolestar/go-commons-pool/v2"
	"github.com/masterzen/winrm"
	"github.com/taliesins/terraform-provider-hyperv/powershell"
)

type HypervClient struct {
	WinRmClientPool  *pool.ObjectPool
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

	command := scriptRendered.String()

	ctx := context.Background()
	winrmClient, err := c.WinRmClientPool.BorrowObject(ctx)

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Running fire and forget script:\n%s\n", command)

	_, _, _, err = powershell.RunPowershell(winrmClient.(*winrm.Client), c.ElevatedUser, c.ElevatedPassword, c.Vars, command)

	err2 := c.WinRmClientPool.ReturnObject(ctx, winrmClient)

	if err != nil {
		return err
	}

	if err2 != nil {
		return err2
	}

	return nil
}

func (c *HypervClient) runScriptWithResult(script *template.Template, args interface{}, result interface{}) (err error) {
	var scriptRendered bytes.Buffer
	err = script.Execute(&scriptRendered, args)

	if err != nil {
		return err
	}

	command := scriptRendered.String()

	ctx := context.Background()
	winrmClient, err := c.WinRmClientPool.BorrowObject(ctx)

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Running script with result:\n%s\n", command)

	_, stdout, _, err := powershell.RunPowershell(winrmClient.(*winrm.Client), c.ElevatedUser, c.ElevatedPassword, c.Vars, command)

	err2 := c.WinRmClientPool.ReturnObject(ctx, winrmClient)

	if err != nil {
		return err
	}

	if err2 != nil {
		return err2
	}

	stdout = strings.TrimSpace(stdout)

	err = json.Unmarshal([]byte(stdout), &result)
	if err != nil {
		return fmt.Errorf("%s\n%s", err, stdout)
	}

	return nil
}
