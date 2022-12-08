//go:build integration
// +build integration

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestHyperVDataSourceNetworkSwitch(t *testing.T) {
	// Skip if -short flag exist
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	name := fmt.Sprintf("wan_%d", randInt())

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testHyperDataSourceVNetworkSwitchConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.hyperv_network_switch.this", "name", name),
				),
			},
		},
	})
}

func testHyperDataSourceVNetworkSwitchConfig(name string) string {
	return fmt.Sprintf(`
resource "hyperv_network_switch" "this" {
	name = "%s"
}

data "hyperv_network_switch" "this" {
	name = hyperv_network_switch.this.name
}
	`, escapeForHcl(name))
}
