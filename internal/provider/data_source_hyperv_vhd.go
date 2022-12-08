package provider

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/taliesins/terraform-provider-hyperv/api"
)

func dataSourceHyperVVhd() *schema.Resource {
	return &schema.Resource{
		Description: "Get information about existing vhd/vhdx/vhds.",
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(ReadVhdTimeout),
		},
		ReadContext: datasourceHyperVVhdRead,
		Schema: map[string]*schema.Schema{
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Path to the existing virtual hard disk file(s) that is being created or being copied to. If a filename or relative path is specified, the virtual hard disk path is calculated relative to the current working directory. Depending on the source selected, the path will be used to determine where to copy source vhd/vhdx/vhds file to.",
			},
			"source": {
				Type:     schema.TypeString,
				Optional: true,
				ConflictsWith: []string{
					"source_vm",
					"parent_path",
					"source_disk",
				},
				Description: "This value can be a url or a path (including wildcards). Box, Zip and 7z files will automatically be expanded. The destination folder will be the directory portion of the path. If expanded files have a folder called `Virtual Machines`, then the `Virtual Machines` folder will be used instead of the entire archive contents. ",
			},
			"source_vm": {
				Type:     schema.TypeString,
				Optional: true,
				ConflictsWith: []string{
					"source",
					"parent_path",
					"source_disk",
				},
				Description: "This value is the name of the vm to copy the vhds from.",
			},
			"source_disk": {
				Type:     schema.TypeInt,
				Optional: true,
				ConflictsWith: []string{
					"source",
					"source_vm",
					"parent_path",
				},
				Description: "Specifies the physical disk to be used as the source for the virtual hard disk to be created.",
			},
			"vhd_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.VhdType_name[api.VhdType_Dynamic],
				ValidateDiagFunc: stringKeyInMap(api.VhdType_value, true),
				ConflictsWith: []string{
					"source",
					"source_vm",
				},
				Description: "VHD type to use. Valid values to use are `Unknown`, `Fixed`, `Dynamic`, `Differencing`.",
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
				Description: "Specifies the path to the parent of the differencing disk to be created (this parameter may be specified only for the creation of a differencing disk)",
			},
			"size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
				ConflictsWith: []string{
					"parent_path",
				},
				Description: "The maximum size, in bytes, of the virtual hard disk to be created.",
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
				Description: "Specifies the block size, in bytes, of the virtual hard disk to be created.",
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
				ValidateDiagFunc: IntInSlice([]int{0, 512, 4096}),
				Description:      "Specifies the logical sector size, in bytes, of the virtual hard disk to be created. Valid values to use are `0`, `512`, `4096`.",
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
				ValidateDiagFunc: IntInSlice([]int{0, 512, 4096}),
				Description:      "Specifies the physical sector size, in bytes. Valid values to use are `0`, `512`, `4096`.",
			},
			"exists": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "",
			},
		},
	}
}

func datasourceHyperVVhdRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO][hyperv][read] reading hyperv vhd: %#v", d)
	c := meta.(api.Client)

	path := ""

	if v, ok := d.GetOk("path"); ok {
		path = v.(string)
	} else {
		return diag.Errorf("[ERROR][hyperv][read] path argument is required")
	}

	vhd, err := c.GetVhd(ctx, path)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO][hyperv][read] retrieved vhd: %+v", vhd)

	if err := d.Set("path", vhd.Path); err != nil {
		return diag.FromErr(err)
	}

	if vhd.Path != "" {
		log.Printf("[INFO][hyperv][read] unable to retrieved vhd: %+v", path)
		if err := d.Set("exists", false); err != nil {
			return diag.FromErr(err)
		}
	} else {
		log.Printf("[INFO][hyperv][read] retrieved vhd: %+v", path)
		if err := d.Set("exists", true); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(path)

	log.Printf("[INFO][hyperv][read] read hyperv vhd: %#v", d)

	return nil
}
