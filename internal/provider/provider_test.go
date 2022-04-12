package provider

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/masterzen/winrm"
	"github.com/taliesins/terraform-provider-hyperv/powershell"
)

func TestGetRemotePath(t *testing.T) {
	config := &Config{
		User:          "Administrator",
		Password:      "",
		Host:          "localhost",
		Port:          5986,
		HTTPS:         true,
		Insecure:      true,
		NTLM:          true,
		TLSServerName: "",
		CACert:        nil,
		Key:           nil,
		Cert:          nil,
		ScriptPath:    "C:/Temp/terraform_%RAND%.cmd",
		Timeout:       "30s",
	}

	client, err := config.Client()
	if err != nil {
		t.Errorf("Unable to get client for HyperV: %s", err.Error())
	}

	ctx := context.Background()
	winrmClient, err := client.WinRmClientPool.BorrowObject(ctx)

	if err != nil {
		t.Errorf("Unable to borrow winrm client for HyperV: %s", err.Error())
	}

	tempFile := fmt.Sprintf("terraform-%s", powershell.TimeOrderedUUID())
	tempPath := fmt.Sprintf(`%s\%s`, `$env:TEMP`, tempFile)
	log.Printf("Resolving remote temp path of [%s]", tempPath)
	tempPath, err = powershell.ResolvePath(winrmClient.(*winrm.Client), tempPath)
	err2 := client.WinRmClientPool.ReturnObject(ctx, winrmClient)

	if err != nil {
		t.Errorf("Unable to resolve remote path: %s", err.Error())
	}

	if err2 != nil {
		t.Errorf("Unable to release winrm client: %s", err2.Error())
	}
	log.Printf("Remote temp path resolved to [%s]", tempPath)
}
