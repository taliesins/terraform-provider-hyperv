package winrm_helper

import (
	"context"
	"text/template"
)

type Client interface {
	RunFireAndForgetScript(ctx context.Context, script *template.Template, args interface{}) error
	RunScriptWithResult(ctx context.Context, script *template.Template, args interface{}, result interface{}) (err error)
	UploadFile(ctx context.Context, filePath string, remoteFilePath string) (resolvedRemoteFilePath string, err error)
	UploadDirectory(ctx context.Context, rootPath string, excludeList []string) (remoteRootPath string, remoteAbsoluteFilePaths []string, err error)
	FileExists(ctx context.Context, remoteFilePath string) (exists bool, err error)
	DirectoryExists(ctx context.Context, remoteDirectoryPath string) (exists bool, err error)
	DeleteFileOrDirectory(ctx context.Context, remotePath string) (err error)
}

type Provider struct {
	Client Client
}
