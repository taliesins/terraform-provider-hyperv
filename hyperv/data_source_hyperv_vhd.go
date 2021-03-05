package hyperv

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/taliesins/terraform-provider-hyperv/api"
)

func dataSourceHyperVVhd() *schema.Resource {
	return &schema.Resource{
		Read: resourceHyperVVhdRead,

		Schema: map[string]*schema.Schema{
			"path": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source": {
				Type:     schema.TypeString,
				Optional: true,
				ConflictsWith: []string{
					"source_vm",
					"parent_path",
					"source_disk",
				},
			},
			"source_vm": {
				Type:     schema.TypeString,
				Optional: true,
				ConflictsWith: []string{
					"source",
					"parent_path",
					"source_disk",
				},
			},
			"source_disk": {
				Type:     schema.TypeInt,
				Optional: true,
				ConflictsWith: []string{
					"source",
					"source_vm",
					"parent_path",
				},
			},
			"vhd_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      api.VhdType_name[api.VhdType_Dynamic],
				ValidateFunc: stringKeyInMap(api.VhdType_value, true),
				ConflictsWith: []string{
					"source",
					"source_vm",
				},
			},
			"parent_path": {
				Type:     schema.TypeString,
				Optional: true,
				ConflictsWith: []string{
					"source",
					"source_vm",
					"source_disk",
					"size",
				},
			},
			"size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
				ConflictsWith: []string{
					"parent_path",
				},
			},
			"block_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
				ConflictsWith: []string{
					"source",
					"source_vm",
					"parent_path",
				},
			},
			"logical_sector_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
				ConflictsWith: []string{
					"source",
					"source_vm",
					"parent_path",
				},
				ValidateFunc: IntInSlice([]int{0, 512, 4096}),
			},
			"physical_sector_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
				ConflictsWith: []string{
					"source",
					"source_vm",
					"parent_path",
				},
				ValidateFunc: IntInSlice([]int{0, 512, 4096}),
			},
			"exists": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}
