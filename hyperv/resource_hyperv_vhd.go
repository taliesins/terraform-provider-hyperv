package hyperv

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/taliesins/terraform-provider-hyperv/api"
	"os"
)

func resourceHyperVVhd() *schema.Resource {
	return &schema.Resource{
		Create: resourceHyperVVhdCreate,
		Read:   resourceHyperVVhdRead,
		Update: resourceHyperVVhdUpdate,
		Delete: resourceHyperVVhdDelete,

		Schema: map[string]*schema.Schema{
			"path": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source_path": {
				Type:     schema.TypeString,
				Optional: true,
				ConflictsWith: []string{
					"source_url",
					"parent_path",
					"source_disk",
				},
			},
			"source_url": {
				Type:     schema.TypeString,
				Optional: true,
				ConflictsWith: []string{
					"source_path",
					"parent_path",
					"source_disk",
				},
			},
			"source_disk": {
				Type:     schema.TypeInt,
				Optional: true,
				ConflictsWith: []string{
					"source_path",
					"source_url",
					"parent_path",
				},
			},
			"vhd_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      api.VMVhdType_name[api.VMVhdType_Dynamic],
				ValidateFunc: stringKeyInMap(api.VMVhdType_value, true),
				ConflictsWith: []string{
					"source_path",
					"source_url",
					"parent_path",
				},
			},
			"parent_path": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{
					"source_path",
					"source_url",
					"source_disk",
					"size",
				},
			},
			"size": {
				Type:          schema.TypeInt,
				Optional:      true,
				Default:       0,
				ConflictsWith: []string{
					"parent_path",
				},
			},
			"block_size": {
				Type:          schema.TypeInt,
				Optional:      true,
				Default:       0,
				ConflictsWith: []string{
					"source_path",
					"source_url",
					"parent_path",
				},
			},
			"logical_sector_size": {
				Type:          schema.TypeInt,
				Optional:      true,
				Default:       0,
				ConflictsWith: []string{
					"source_path",
					"source_url",
					"parent_path",
				},
				ValidateFunc: IntInSlice([]int{512, 4096}),
			},
			"physical_sector_size": {
				Type:          schema.TypeInt,
				Optional:      true,
				Default:       0,
				ConflictsWith: []string{
					"source_path",
					"source_url",
					"parent_path",
				},
				ValidateFunc: IntInSlice([]int{512, 4096}),
			},
			"exists": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},

		CustomizeDiff: customizeDiffForVhd,
	}
}

func customizeDiffForVhd(diff *schema.ResourceDiff, i interface{}) error {
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

func resourceHyperVVhdCreate(d *schema.ResourceData, meta interface{}) (err error) {

	log.Printf("[INFO][hyperv][create] creating hyperv vhdx: %#v", d)
	c := meta.(*api.HypervClient)

	path := ""

	if v, ok := d.GetOk("path"); ok {
		path = v.(string)
	} else {
		return fmt.Errorf("[ERROR][hyperv][create] path argument is required")
	}

	sourcePath := (d.Get("source_path")).(string)
	sourceUrl := (d.Get("source_url")).(string)
	sourceDisk := (d.Get("source_disk")).(int)
	vhdType := api.ToVMVhdType((d.Get("vhd_type")).(string))
	parentPath := (d.Get("parent_path")).(string)
	size := (d.Get("size")).(uint64)
	blockSize := (d.Get("block_size")).(uint32)
	logicalSectorSize := (d.Get("logical_sector_size")).(uint32)
	physicalSectorSize := (d.Get("physical_sector_size")).(uint32)

	err = c.CreateOrUpdateVMVhd(path, sourcePath, sourceUrl, sourceDisk, vhdType, parentPath, size, blockSize, logicalSectorSize, physicalSectorSize)

	if err != nil {
		return err
	}

	if size > 0 && parentPath == "" {
		//Update vhdx size
		err = c.ResizeVMVhd(path, size)

		if err != nil {
			return err
		}
	}

	d.SetId(path)

	log.Printf("[INFO][hyperv][create] created hyperv vhdx: %#v", d)

	return resourceHyperVVhdRead(d, meta)
}

func resourceHyperVVhdRead(d *schema.ResourceData, meta interface{}) (err error) {
	log.Printf("[INFO][hyperv][read] reading hyperv vhdx: %#v", d)
	c := meta.(*api.HypervClient)

	path := ""

	if v, ok := d.GetOk("path"); ok {
		path = v.(string)
	} else {
		return fmt.Errorf("[ERROR][hyperv][read] path argument is required")
	}

	vmVhdx, err := c.GetVMVhd(path)
	if err != nil {
		return err
	}

	if vmVhdx.Path != "" {
		log.Printf("[INFO][hyperv][read] unable to retrieved vhdx: %+v", path)
		d.Set("exists", false)
	} else {
		log.Printf("[INFO][hyperv][read] retrieved vhdx: %+v", path)
		d.Set("exists", true)
	}

	log.Printf("[INFO][hyperv][read] read hyperv vhdx: %#v", d)

	return nil
}

func resourceHyperVVhdUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	log.Printf("[INFO][hyperv][update] updating hyperv vhdx: %#v", d)
	c := meta.(*api.HypervClient)

	path := ""

	if v, ok := d.GetOk("path"); ok {
		path = v.(string)
	} else {
		return fmt.Errorf("[ERROR][hyperv][update] path argument is required")
	}

	sourcePath := (d.Get("source_path")).(string)
	sourceUrl := (d.Get("source_url")).(string)
	sourceDisk := (d.Get("source_disk")).(int)
	vhdType := api.ToVMVhdType((d.Get("vhd_type")).(string))
	parentPath := (d.Get("parent_path")).(string)
	size := (d.Get("size")).(uint64)
	blockSize := (d.Get("block_size")).(uint32)
	logicalSectorSize := (d.Get("logical_sector_size")).(uint32)
	physicalSectorSize := (d.Get("physical_sector_size")).(uint32)

	exists := (d.Get("exists")).(bool)

	if !exists || d.HasChange("path") || d.HasChange("source_path") || d.HasChange("source_url") || d.HasChange("source_disk") || d.HasChange("parent_path")  {
		//delete it as its changed
		err = c.CreateOrUpdateVMVhd(path, sourcePath, sourceUrl, sourceDisk, vhdType, parentPath, size, blockSize, logicalSectorSize, physicalSectorSize)

		if err != nil {
			return err
		}
	}

	if size > 0 && parentPath == "" {
		if !exists || d.HasChange("size") {
			//Update vhdx size
			err = c.ResizeVMVhd(path, size)

			if err != nil {
				return err
			}
		}
	}

	log.Printf("[INFO][hyperv][update] updated hyperv vhdx: %#v", d)

	return resourceHyperVVhdRead(d, meta)
}

func resourceHyperVVhdDelete(d *schema.ResourceData, meta interface{}) (err error) {
	log.Printf("[INFO][hyperv][delete] deleting hyperv vhdx: %#v", d)

	c := meta.(*api.HypervClient)

	path := ""

	if v, ok := d.GetOk("name"); ok {
		path = v.(string)
	} else {
		return fmt.Errorf("[ERROR][hyperv][delete] path argument is required")
	}

	err = c.DeleteVMVhd(path)

	if err != nil {
		return err
	}

	log.Printf("[INFO][hyperv][delete] deleted hyperv vhdx: %#v", d)
	return nil
}
