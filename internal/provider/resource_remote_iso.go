package provider

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/taliesins/terraform-provider-hyperv/api"
)

const (
	ReadRemoteIsoTimeout   = 1 * time.Minute
	CreateRemoteIsoTimeout = 5 * time.Minute
	UpdateRemoteIsoTimeout = 5 * time.Minute
	DeleteRemoteIsoTimeout = 1 * time.Minute
)

func resourceRemoteIso() *schema.Resource {
	return &schema.Resource{
		Description: "This resource allows you to manage ISOs.",
		Timeouts: &schema.ResourceTimeout{
			Read:   schema.DefaultTimeout(ReadRemoteIsoTimeout),
			Create: schema.DefaultTimeout(CreateRemoteIsoTimeout),
			Update: schema.DefaultTimeout(UpdateRemoteIsoTimeout),
			Delete: schema.DefaultTimeout(DeleteRemoteIsoTimeout),
		},
		CreateContext: resourceRemoteIsoCreate,
		ReadContext:   resourceRemoteIsoRead,
		UpdateContext: resourceRemoteIsoUpdate,
		DeleteContext: resourceRemoteIsoDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"sourceDirectoryPath": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Path to directory for files to be copied into iso.",
			},
			"sourceBootFilePath": {
				Type:        schema.TypeString,
				Required:    false,
				Description: "Path to boot file to be copied into iso.",
			},
			"destinationIsoPath": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Remote path for iso.",
			},
			"excludeList": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "The regex paths to exclude when including files for iso.",
			},
			"media": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Media type for iso.",
			},
			"title": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Title for iso.",
			},

			"path": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					extension := path.Ext(newValue)
					computedPath := strings.TrimSuffix(newValue, extension)

					// Ignore differencing
					if strings.HasPrefix(strings.ToLower(oldValue), strings.ToLower(computedPath)) && strings.HasSuffix(strings.ToLower(oldValue), strings.ToLower(extension)) {
						return true
					}

					if strings.EqualFold(oldValue, newValue) {
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

			"exists": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Does virtual disk exist.",
			},
		},

		CustomizeDiff: customizeDiffForRemoteIso,
	}
}

func customizeDiffForRemoteIso(ctx context.Context, diff *schema.ResourceDiff, i interface{}) error {
	path := diff.Get("path").(string)

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			// file does not exist
			diff.SetNew("exists", false)
			return nil
		} else {
			// other error
			return err
		}
	}

	return nil
}

func resourceRemoteIsoCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO][remote-iso][create] creating remote iso: %#v", d)
	c := meta.(api.Client)

	path := ""

	if v, ok := d.GetOk("path"); ok {
		path = v.(string)
	} else {
		return diag.Errorf("[ERROR][remote-iso][create] path argument is required")
	}

	if d.IsNewResource() {
		existing, err := c.RemoteIsoExists(ctx, path)
		if err != nil {
			return diag.FromErr(fmt.Errorf("checking for existing %s: %+v", path, err))
		}

		if existing.Exists {
			return diag.FromErr(fmt.Errorf("A resource with the ID %q already exists - to be managed via Terraform this resource needs to be imported into the State. Please see the resource documentation for %q for more information.\n terraform import %s.<resource name> %s", path, "hyperv_vhd", "hyperv_vhd", path))
		}
	}

	source := (d.Get("source")).(string)
	sourceVm := (d.Get("source_vm")).(string)
	sourceDisk := (d.Get("source_disk")).(int)
	vhdType := api.ToRemoteIsoType((d.Get("vhd_type")).(string))
	parentPath := (d.Get("parent_path")).(string)
	size := uint64((d.Get("size")).(int))
	blockSize := uint32((d.Get("block_size")).(int))
	logicalSectorSize := uint32((d.Get("logical_sector_size")).(int))
	physicalSectorSize := uint32((d.Get("physical_sector_size")).(int))

	err := c.CreateOrUpdateRemoteIso(ctx, path, source, sourceVm, sourceDisk, vhdType, parentPath, size, blockSize, logicalSectorSize, physicalSectorSize)

	if err != nil {
		return diag.FromErr(err)
	}

	if size > 0 && parentPath == "" {
		// Update vhd size
		err = c.ResizeRemoteIso(ctx, path, size)

		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(path)
	log.Printf("[INFO][remote-iso][create] created remote iso: %#v", d)

	return resourceRemoteIsoRead(ctx, d, meta)
}

func resourceRemoteIsoRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO][remote-iso][read] reading remote iso: %#v", d)
	c := meta.(api.Client)

	path := d.Id()

	vhd, err := c.GetRemoteIso(ctx, path)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO][remote-iso][read] retrieved vhd: %+v", vhd)

	if err := d.Set("path", vhd.Path); err != nil {
		return diag.FromErr(err)
	}

	if vhd.Path == "" {
		log.Printf("[INFO][remote-iso][read] unable to retrieved vhd: %+v", path)
		if err := d.Set("exists", false); err != nil {
			return diag.FromErr(err)
		}
	} else {
		log.Printf("[INFO][remote-iso][read] retrieved vhd: %+v", path)
		if err := d.Set("exists", true); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := d.Set("vhd_type", api.RemoteIsoType_name[vhd.RemoteIsoType]); err != nil {
		return diag.FromErr(err)
	}

	if vhd.RemoteIsoType == api.RemoteIsoType_Differencing {
		if err := d.Set("parent_path", vhd.ParentPath); err != nil {
			return diag.FromErr(err)
		}
	} else {
		if err := d.Set("size", vhd.Size); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("block_size", vhd.BlockSize); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("logical_sector_size", vhd.LogicalSectorSize); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("physical_sector_size", vhd.PhysicalSectorSize); err != nil {
			return diag.FromErr(err)
		}
	}

	log.Printf("[INFO][remote-iso][read] read remote iso: %#v", d)

	return nil
}

func resourceRemoteIsoUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO][remote-iso][update] updating remote iso: %#v", d)
	c := meta.(api.Client)

	path := d.Id()

	source := (d.Get("source")).(string)
	sourceVm := (d.Get("source_vm")).(string)
	sourceDisk := (d.Get("source_disk")).(int)
	vhdType := api.ToRemoteIsoType((d.Get("vhd_type")).(string))
	parentPath := (d.Get("parent_path")).(string)
	size := uint64((d.Get("size")).(int))
	blockSize := uint32((d.Get("block_size")).(int))
	logicalSectorSize := uint32((d.Get("logical_sector_size")).(int))
	physicalSectorSize := uint32((d.Get("physical_sector_size")).(int))

	exists := (d.Get("exists")).(bool)

	if !exists || d.HasChange("path") || d.HasChange("source") || d.HasChange("source_vm") || d.HasChange("source_disk") || d.HasChange("parent_path") {
		// delete it as its changed
		err := c.CreateOrUpdateRemoteIso(ctx, path, source, sourceVm, sourceDisk, vhdType, parentPath, size, blockSize, logicalSectorSize, physicalSectorSize)

		if err != nil {
			return diag.FromErr(err)
		}
	}

	if size > 0 && parentPath == "" {
		if !exists || d.HasChange("size") {
			// Update vhd size
			err := c.ResizeRemoteIso(ctx, path, size)

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	log.Printf("[INFO][remote-iso][update] updated remote iso: %#v", d)

	return resourceRemoteIsoRead(ctx, d, meta)
}

func resourceRemoteIsoDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO][remote-iso][delete] deleting remote iso: %#v", d)

	c := meta.(api.Client)

	path := d.Id()

	err := c.DeleteIsoImage(ctx, path)

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO][remote-iso][delete] deleted remote iso: %#v", d)
	return nil
}
