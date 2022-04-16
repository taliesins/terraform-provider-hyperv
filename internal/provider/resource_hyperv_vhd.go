package provider

import (
	"context"
	"log"
	"path"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/taliesins/terraform-provider-hyperv/api"
)

var defaultVVhdTimeoutDuration = time.Minute * 30

func resourceHyperVVhd() *schema.Resource {
	return &schema.Resource{
		Description: "This Hyper-V resource allows you to manage VHDs.",
		Timeouts: &schema.ResourceTimeout{
			Default: &defaultVVhdTimeoutDuration,
		},
		CreateContext: resourceHyperVVhdCreate,
		ReadContext:   resourceHyperVVhdRead,
		UpdateContext: resourceHyperVVhdUpdate,
		DeleteContext: resourceHyperVVhdDelete,
		Schema: map[string]*schema.Schema{
			"path": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					extension := path.Ext(newValue)
					computedPath := strings.TrimSuffix(newValue, extension)

					//Ignore differencing
					if strings.HasPrefix(oldValue, computedPath) && strings.HasSuffix(oldValue, extension) {
						return true
					}

					if oldValue == newValue {
						return true
					}

					return false
				},
				Description: "Path to the new virtual hard disk file(s) that is being created or being copied to. If a filename or relative path is specified, the new virtual hard disk path is calculated relative to the current working directory. Depending on the source selected, the path will be used to determine where to copy source vhd/vhdx/vhds file to.",
			},
			"source": {
				Type:     schema.TypeString,
				Optional: true,
				ConflictsWith: []string{
					"source_vm",
					"parent_path",
					"source_disk",
				},
				Description: "This field is mutually exclusive with the fields `source_vm`, `parent_path`, `source_disk`. This value can be a url or a path (including wildcards). Box, Zip and 7z files will automatically be expanded. The destination folder will be the directory portion of the path. If expanded files have a folder called `Virtual Machines`, then the `Virtual Machines` folder will be used instead of the entire archive contents. ",
			},
			"source_vm": {
				Type:     schema.TypeString,
				Optional: true,
				ConflictsWith: []string{
					"source",
					"parent_path",
					"source_disk",
				},
				Description: "This field is mutually exclusive with the fields `source`, `parent_path`, `source_disk`. This value is the name of the vm to copy the vhds from.",
			},
			"source_disk": {
				Type:     schema.TypeInt,
				Optional: true,
				ConflictsWith: []string{
					"source",
					"source_vm",
					"parent_path",
				},
				Description: "This field is mutually exclusive with the fields `source`, `source_vm`, `parent_path`. Specifies the physical disk to be used as the source for the virtual hard disk to be created.",
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
				Description: "This field is mutually exclusive with the fields `source`, `source_vm`, `parent_path`. Valid values to use are `Unknown`, `Fixed`, `Dynamic`, `Differencing`.",
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
				Description: "This field is mutually exclusive with the fields `source`, `source_vm`, `source_disk`, `size`. Specifies the path to the parent of the differencing disk to be created (this parameter may be specified only for the creation of a differencing disk).",
			},
			"size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
				ConflictsWith: []string{
					"parent_path",
				},
				Description: "This field is mutually exclusive with the field `parent_path`. The maximum size, in bytes, of the virtual hard disk to be created.",
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
				Description: "This field is mutually exclusive with the fields `source`, `source_vm`, `parent_path`. Specifies the block size, in bytes, of the virtual hard disk to be created.",
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
				Description:      "This field is mutually exclusive with the fields `source`, `source_vm`, `parent_path`. Specifies the logical sector size, in bytes, of the virtual hard disk to be created. Valid values to use are `0`, `512`, `4096`.",
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
				Description: "This field is mutually exclusive with the fields	`source`, `source_vm`, `parent_path`. Specifies the physical sector size, in bytes. Valid values to use are `0`, `512`, `4096`.",
			},
			"exists": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Does virtual machine exist.",
			},
		},

		CustomizeDiff: customizeDiffForVhd,
	}
}

func customizeDiffForVhd(ctx context.Context, diff *schema.ResourceDiff, i interface{}) error {
	path := diff.Get("path").(string)

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			// file does not exist
			diff.SetNewComputed("exists")
			return nil
		} else {
			// other error
			return err
		}
	}

	return nil
}

func resourceHyperVVhdCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	log.Printf("[INFO][hyperv][create] creating hyperv vhd: %#v", d)
	c := meta.(api.Client)

	path := ""

	if v, ok := d.GetOk("path"); ok {
		path = v.(string)
	} else {
		return diag.Errorf("[ERROR][hyperv][create] path argument is required")
	}

	source := (d.Get("source")).(string)
	sourceVm := (d.Get("source_vm")).(string)
	sourceDisk := (d.Get("source_disk")).(int)
	vhdType := api.ToVhdType((d.Get("vhd_type")).(string))
	parentPath := (d.Get("parent_path")).(string)
	size := uint64((d.Get("size")).(int))
	blockSize := uint32((d.Get("block_size")).(int))
	logicalSectorSize := uint32((d.Get("logical_sector_size")).(int))
	physicalSectorSize := uint32((d.Get("physical_sector_size")).(int))

	err := c.CreateOrUpdateVhd(path, source, sourceVm, sourceDisk, vhdType, parentPath, size, blockSize, logicalSectorSize, physicalSectorSize)

	if err != nil {
		return diag.FromErr(err)
	}

	if size > 0 && parentPath == "" {
		//Update vhd size
		err = c.ResizeVhd(path, size)

		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(path)

	log.Printf("[INFO][hyperv][create] created hyperv vhd: %#v", d)

	return resourceHyperVVhdRead(ctx, d, meta)
}

func resourceHyperVVhdRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO][hyperv][read] reading hyperv vhd: %#v", d)
	c := meta.(api.Client)

	path := ""

	if v, ok := d.GetOk("path"); ok {
		path = v.(string)
	} else {
		return diag.Errorf("[ERROR][hyperv][read] path argument is required")
	}

	vhd, err := c.GetVhd(path)
	if err != nil {
		return diag.FromErr(err)
	}

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

func resourceHyperVVhdUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO][hyperv][update] updating hyperv vhd: %#v", d)
	c := meta.(api.Client)

	path := ""

	if v, ok := d.GetOk("path"); ok {
		path = v.(string)
	} else {
		return diag.Errorf("[ERROR][hyperv][update] path argument is required")
	}

	source := (d.Get("source")).(string)
	sourceVm := (d.Get("source_vm")).(string)
	sourceDisk := (d.Get("source_disk")).(int)
	vhdType := api.ToVhdType((d.Get("vhd_type")).(string))
	parentPath := (d.Get("parent_path")).(string)
	size := uint64((d.Get("size")).(int))
	blockSize := uint32((d.Get("block_size")).(int))
	logicalSectorSize := uint32((d.Get("logical_sector_size")).(int))
	physicalSectorSize := uint32((d.Get("physical_sector_size")).(int))

	exists := (d.Get("exists")).(bool)

	if !exists || d.HasChange("path") || d.HasChange("source") || d.HasChange("source_vm") || d.HasChange("source_disk") || d.HasChange("parent_path") {
		//delete it as its changed
		err := c.CreateOrUpdateVhd(path, source, sourceVm, sourceDisk, vhdType, parentPath, size, blockSize, logicalSectorSize, physicalSectorSize)

		if err != nil {
			return diag.FromErr(err)
		}
	}

	if size > 0 && parentPath == "" {
		if !exists || d.HasChange("size") {
			//Update vhd size
			err := c.ResizeVhd(path, size)

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	log.Printf("[INFO][hyperv][update] updated hyperv vhd: %#v", d)

	return resourceHyperVVhdRead(ctx, d, meta)
}

func resourceHyperVVhdDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO][hyperv][delete] deleting hyperv vhd: %#v", d)

	c := meta.(api.Client)

	path := ""

	if v, ok := d.GetOk("path"); ok {
		path = v.(string)
	} else {
		return diag.Errorf("[ERROR][hyperv][delete] path argument is required")
	}

	err := c.DeleteVhd(path)

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO][hyperv][delete] deleted hyperv vhd: %#v", d)
	return nil
}
