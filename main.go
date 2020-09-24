package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/taliesins/terraform-provider-hyperv/hyperv"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: hyperv.Provider,
	})
}
