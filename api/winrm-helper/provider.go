package winrm_helper

import "text/template"

type Client interface {
	RunFireAndForgetScript(script *template.Template, args interface{}) error
	RunScriptWithResult(script *template.Template, args interface{}, result interface{}) (err error)
}

type Provider struct {
	Client Client
}
