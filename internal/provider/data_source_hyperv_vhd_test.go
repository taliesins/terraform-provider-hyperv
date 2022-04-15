package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestHyperVDataSourceVhd(t *testing.T) {
	path := fmt.Sprintf("C:\\data\\VirtualMachines\\web_server\\Virtual Hard Disks\\web_server_%d.vhdx", randInt())

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testHyperDataSourceVVhdConfig(path),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.hyperv_vhd.web_server", "path", path),
				),
			},
		},
	})
}

func testHyperDataSourceVVhdConfig(path string) string {
	return fmt.Sprintf(`
data "hyperv_vhd" "web_server" {
	path = "%s"
}
	`, escapeForHcl(path))
}
