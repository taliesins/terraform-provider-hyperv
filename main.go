package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/tidalf/terraform-provider-hyperv/hyperv"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: hyperv.Provider,
	})
}
