//go:build integration
// +build integration

package provider

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestHyperVDataSourceVhd(t *testing.T) {
	// Skip if -short flag exist
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	//tempDirectory := os.TempDir() uses short name ;<
	tempDirectory, _ := filepath.Abs(".")
	path, _ := filepath.Abs(filepath.Join(tempDirectory, fmt.Sprintf("testhypervdatasourcevhd_%d.vhdx", randInt())))

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testHyperDataSourceVVhdConfig(path),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.hyperv_vhd.this", "path", path),
				),
			},
		},
	})
}

func testHyperDataSourceVVhdConfig(path string) string {
	return fmt.Sprintf(`
resource "hyperv_vhd" "this" {
	path = "%s"
	size = 4001792
}

data "hyperv_vhd" "this" {
	path = hyperv_vhd.this.path
}
	`, escapeForHcl(path))
}
