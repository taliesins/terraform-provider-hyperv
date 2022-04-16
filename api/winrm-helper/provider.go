package winrm_helper

import (
	"context"
	"text/template"
)

type Client interface {
	RunFireAndForgetScript(ctx context.Context, script *template.Template, args interface{}) error
	RunScriptWithResult(ctx context.Context, script *template.Template, args interface{}, result interface{}) (err error)
}

type Provider struct {
	Client Client
}
