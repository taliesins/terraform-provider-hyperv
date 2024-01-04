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

func (c *ClientConfig) UploadFile(ctx context.Context, filePath string, remoteFilePath string) (string, error) {
	winrmClient, err := c.WinRmClientPool.BorrowObject(ctx)

	if err != nil {
		return "", err
	}

	log.Printf("[DEBUG] upload file %#v", filePath)

	remoteFilePath, err = powershell.UploadFile(winrmClient.(*winrm.Client), filePath, remoteFilePath)

	err2 := c.WinRmClientPool.ReturnObject(ctx, winrmClient)

	if err != nil {
		return "", err
	}

	if err2 != nil {
		return "", err2
	}

	log.Printf("[DEBUG] uploaded file %#v to %#v", filePath, remoteFilePath)

	return remoteFilePath, nil
}

func (c *ClientConfig) UploadDirectory(ctx context.Context, rootPath string, excludeList []string) (remoteRootPath string, remoteAbsoluteFilePaths []string, err error) {
	winrmClient, err := c.WinRmClientPool.BorrowObject(ctx)

	if err != nil {
		return "", []string{}, err
	}

	log.Printf("[DEBUG] upload directory %#v", rootPath)

	remoteRootPath, remoteAbsoluteFilePaths, err = powershell.UploadDirectory(winrmClient.(*winrm.Client), rootPath, excludeList)

	err2 := c.WinRmClientPool.ReturnObject(ctx, winrmClient)

	if err != nil {
		return "", []string{}, err
	}

	if err2 != nil {
		return "", []string{}, err2
	}

	log.Printf("[DEBUG] uploaded directory %#v to %#v. The following files where uploaded %#v", rootPath, remoteRootPath, remoteAbsoluteFilePaths)

	return remoteRootPath, remoteAbsoluteFilePaths, nil
}

func (c *ClientConfig) FileExists(ctx context.Context, remoteFilePath string) (exists bool, err error) {
	winrmClient, err := c.WinRmClientPool.BorrowObject(ctx)

	if err != nil {
		return false, err
	}

	log.Printf("[DEBUG] check file exists %#v", remoteFilePath)

	result, err := powershell.FileExists(winrmClient.(*winrm.Client), remoteFilePath)

	err2 := c.WinRmClientPool.ReturnObject(ctx, winrmClient)

	if err != nil {
		return false, err
	}

	if err2 != nil {
		return false, err2
	}

	if result {
		log.Printf("[DEBUG] file exists %#v", remoteFilePath)
	} else {
		log.Printf("[DEBUG] file does not exists %#v", remoteFilePath)
	}

	return result, nil
}

func (c *ClientConfig) DirectoryExists(ctx context.Context, remoteDirectoryPath string) (exists bool, err error) {
	winrmClient, err := c.WinRmClientPool.BorrowObject(ctx)

	if err != nil {
		return false, err
	}

	log.Printf("[DEBUG] check directory exists %#v", remoteDirectoryPath)

	result, err := powershell.DirectoryExists(winrmClient.(*winrm.Client), remoteDirectoryPath)

	err2 := c.WinRmClientPool.ReturnObject(ctx, winrmClient)

	if err != nil {
		return false, err
	}

	if err2 != nil {
		return false, err2
	}

	if result {
		log.Printf("[DEBUG] directory exists %#v", remoteDirectoryPath)
	} else {
		log.Printf("[DEBUG] directory does not exists %#v", remoteDirectoryPath)
	}

	return result, nil
}

func (c *ClientConfig) DeleteFileOrDirectory(ctx context.Context, remotePath string) (err error) {
	winrmClient, err := c.WinRmClientPool.BorrowObject(ctx)

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] delete file or directory %#v", remotePath)

	err = powershell.DeleteFileOrDirectory(winrmClient.(*winrm.Client), remotePath)

	err2 := c.WinRmClientPool.ReturnObject(ctx, winrmClient)

	if err != nil {
		return err
	}

	if err2 != nil {
		return err2
	}

	log.Printf("[DEBUG] file or directory deleted %#v", remotePath)

	return nil
}
