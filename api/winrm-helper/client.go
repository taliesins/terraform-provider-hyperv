package winrm_helper

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

func New(clientConfig *ClientConfig) (*Provider, error) {
	return &Provider{
		Client: clientConfig,
	}, nil
}

type ClientConfig struct {
	WinRmClientPool  *pool.ObjectPool
	ElevatedUser     string
	ElevatedPassword string
	Vars             string
}

func (c *ClientConfig) RunFireAndForgetScript(ctx context.Context, script *template.Template, args interface{}) error {
	var scriptRendered bytes.Buffer
	err := script.Execute(&scriptRendered, args)

	if err != nil {
		return err
	}

	command := scriptRendered.String()

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

func (c *ClientConfig) RunScriptWithResult(ctx context.Context, script *template.Template, args interface{}, result interface{}) (err error) {
	var scriptRendered bytes.Buffer
	err = script.Execute(&scriptRendered, args)

	if err != nil {
		return err
	}

	command := scriptRendered.String()

	winrmClient, err := c.WinRmClientPool.BorrowObject(ctx)

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Running script with result:\n%s\n", command)

	exitStatus, stdout, stderr, err := powershell.RunPowershell(winrmClient.(*winrm.Client), c.ElevatedUser, c.ElevatedPassword, c.Vars, command)

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
		return fmt.Errorf("exitStatus:%d\nstdOut:%s\nstdErr:%s\nerr:%s\ncommand:%s", exitStatus, stdout, stderr, err, command)
	}

	return nil
}
