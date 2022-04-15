package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestHyperVDataSourceNetworkSwitch(t *testing.T) {
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
data "hyperv_network_switch" "this" {
  name = "%s"
}
	`, escapeForHcl(name))
}
