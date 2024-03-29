package provider

import (
	"context"
	"fmt"
	log "log"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/taliesins/terraform-provider-hyperv/api"
)

const (
	ReadIsoImageTimeout   = 1 * time.Minute
	CreateIsoImageTimeout = 5 * time.Minute
	UpdateIsoImageTimeout = 5 * time.Minute
	DeleteIsoImageTimeout = 1 * time.Minute
)

func resourceHyperVIsoImage() *schema.Resource {
	return &schema.Resource{
		Description: "This resource allows you to manage ISOs.",
		Timeouts: &schema.ResourceTimeout{
			Read:   schema.DefaultTimeout(ReadIsoImageTimeout),
			Create: schema.DefaultTimeout(CreateIsoImageTimeout),
			Update: schema.DefaultTimeout(UpdateIsoImageTimeout),
			Delete: schema.DefaultTimeout(DeleteIsoImageTimeout),
		},
		CreateContext: resourceHyperVIsoImageCreate,
		ReadContext:   resourceHyperVIsoImageRead,
		UpdateContext: resourceHyperVIsoImageUpdate,
		DeleteContext: resourceHyperVIsoImageDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"source_iso_file_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Local iso file path.",
				ConflictsWith: []string{
					"source_zip_file_path",
					"source_zip_file_path_hash",
					"source_boot_file_path",
					"source_boot_file_path_hash",
				},
			},
			"source_iso_file_path_hash": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Hash of local iso file.",
				ConflictsWith: []string{
					"source_zip_file_path",
					"source_zip_file_path_hash",
					"source_boot_file_path",
					"source_boot_file_path_hash",
				},
			},
			"source_zip_file_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Local zip file path.",
				ConflictsWith: []string{
					"source_iso_file_path",
					"source_iso_file_path_hash",
				},
			},
			"source_zip_file_path_hash": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Hash of local zip file.",
				ConflictsWith: []string{
					"source_iso_file_path",
					"source_iso_file_path_hash",
				},
			},
			"source_boot_file_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Local boot file path.",
				ConflictsWith: []string{
					"source_iso_file_path",
					"source_iso_file_path_hash",
				},
			},
			"source_boot_file_path_hash": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Hash of local boot file.",
				ConflictsWith: []string{
					"source_iso_file_path",
					"source_iso_file_path_hash",
				},
			},
			"destination_iso_file_path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Remote iso file path.",
				ForceNew:    true,
			},
			"destination_zip_file_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Remote zip file path. This defaults to `$env:temp\\{filename(source_zip_file_path)}`",
			},
			"destination_boot_file_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Remote boot file path. This defaults to `$env:temp\\{filename(source_boot_file_path)}`",
			},
			"iso_media_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.IsoMediaType_name[api.IsoMediaType_DVDPLUSRW_DUALLAYER],
				ValidateDiagFunc: StringKeyInMap(api.IsoMediaType_value, false),
				Description:      "Media type for iso. Valid values to use are `unknown`, `cdrom`, `cdr`, `cdrw`, `dvdrom`, `dvdram`, `dvdplusr`, `dvdplusrw`, `dvdplusr_duallayer`, `dvddashr`, `dvddashrw`, `dvddashr_duallayer`, `disk`, `dvdplusrw_duallayer`, `hddvdrom`, `hddvdr`, `hddvdram`, `bdrom`, `bdr`, `bdre`.",
			},
			"iso_file_system_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.IsoFileSystemType_name[api.IsoFileSystemType_Unknown],
				ValidateDiagFunc: StringKeyInMap(api.IsoFileSystemType_value, false),
				Description:      "File system type for iso. Valid values to use are `none`, `iso9660`, `joliet`, `iso9660|joliet`, `udf`, `joliet|udf`, `iso9660|joliet|udf`, `unknown`.",
			},
			"volume_name": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "UNTITLED",
				ValidateDiagFunc: AllowedIsoVolumeName(),
				Description:      "Volume name for iso. Must be 15 characters or less. Characters must be `A` through `Z`, `0` through `9` or `_` (underscore)",
			},
			"resolve_destination_iso_file_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The remote iso file path that was used.",
			},
			"resolve_destination_zip_file_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The remote zip file path that was used.",
			},
			"resolve_destination_boot_file_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The remote boot file path that was used.",
			},
		},
	}
}

func winPath(path string) string {
	if len(path) == 0 {
		return path
	}

	if strings.Contains(path, " ") {
		path = fmt.Sprintf("'%s'", strings.Trim(path, "'\""))
	}

	return strings.ReplaceAll(path, "/", "\\")
}

func ensureFileStateCreate(ctx context.Context, d *schema.ResourceData, c api.Client, name string) (string, error) {
	sourceFilePathKey := fmt.Sprintf("source_%s_file_path", name)
	destinationFilePathKey := fmt.Sprintf("destination_%s_file_path", name)

	sourceFilePath := (d.Get(sourceFilePathKey)).(string)
	destinationFilePath := (d.Get(destinationFilePathKey)).(string)

	resolveDestinationFilePath := destinationFilePath
	if sourceFilePath != "" && resolveDestinationFilePath == "" {
		resolveDestinationFilePath = winPath(filepath.Join(`$env:TEMP`, filepath.Base(sourceFilePath)))
	}

	if resolveDestinationFilePath != "" {
		log.Printf("[INFO][iso-image][create] check if file exists: %#v", resolveDestinationFilePath)
		resolveDestinationFilePathExists, err := c.RemoteFileExists(ctx, resolveDestinationFilePath)
		if err != nil {
			return "", fmt.Errorf("checking for existing %s: %+v", resolveDestinationFilePath, err)
		}
		if resolveDestinationFilePathExists {
			return "", fmt.Errorf("A resource with the ID %q already exists - to be managed via Terraform this resource needs to be imported into the State. Please see the resource documentation for %q for more information.\n terraform import %s.<resource name> %s", resolveDestinationFilePath, "remote_iso", "remote_iso", resolveDestinationFilePath)
		}

		if sourceFilePath != "" {
			err = c.RemoteFileUpload(ctx, sourceFilePath, resolveDestinationFilePath)
			if err != nil {
				return resolveDestinationFilePath, err
			}
		}
	}

	// if d.IsNewResource() {
	// }

	return resolveDestinationFilePath, nil
}

func resourceHyperVIsoImageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO][iso-image][create] creating remote iso: %#v", d)
	c := meta.(api.Client)

	sourceIsoFilePath := (d.Get("source_iso_file_path")).(string)
	sourceIsoFilePathHash := (d.Get("source_iso_file_path_hash")).(string)
	sourceZipFilePath := (d.Get("source_zip_file_path")).(string)
	sourceZipFilePathHash := (d.Get("source_zip_file_path_hash")).(string)
	sourceBootFilePath := (d.Get("source_boot_file_path")).(string)
	sourceBootFilePathHash := (d.Get("source_boot_file_path_hash")).(string)
	destinationIsoFilePath := (d.Get("destination_iso_file_path")).(string)
	destinationZipFilePath := (d.Get("destination_zip_file_path")).(string)
	destinationBootFilePath := (d.Get("destination_boot_file_path")).(string)
	media := api.ToIsoMediaType((d.Get("iso_media_type")).(string))
	fileSystem := api.ToIsoFileSystemType((d.Get("iso_file_system_type")).(string))
	volumeName := (d.Get("volume_name")).(string)

	if destinationIsoFilePath == "" {
		return diag.Errorf("[ERROR][iso-image][create] path argument is required")
	}

	resolveDestinationIsoFilePath, err := ensureFileStateCreate(ctx, d, c, "iso")
	if err != nil {
		return diag.FromErr(err)
	}

	resolveDestinationZipFilePath, err := ensureFileStateCreate(ctx, d, c, "zip")
	if err != nil {
		return diag.FromErr(err)
	}

	resolveDestinationBootFilePath, err := ensureFileStateCreate(ctx, d, c, "boot")
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.CreateOrUpdateIsoImage(ctx, sourceIsoFilePath, sourceIsoFilePathHash, sourceZipFilePath, sourceZipFilePathHash, sourceBootFilePath, sourceBootFilePathHash, destinationIsoFilePath, destinationZipFilePath, destinationBootFilePath, media, fileSystem, volumeName, resolveDestinationIsoFilePath, resolveDestinationZipFilePath, resolveDestinationBootFilePath)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(destinationIsoFilePath)
	log.Printf("[INFO][iso-image][create] created remote iso: %#v", d)

	return resourceHyperVIsoImageRead(ctx, d, meta)
}

func resourceHyperVIsoImageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO][iso-image][read] reading remote iso: %#v", d)
	c := meta.(api.Client)

	destinationIsoFilePath := d.Id()

	isoImage, err := c.GetIsoImage(ctx, destinationIsoFilePath)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO][iso-image][read] retrieved isoImage: %+v", isoImage)

	if err := d.Set("source_iso_file_path", isoImage.SourceIsoFilePath); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("source_iso_file_path_hash", isoImage.SourceIsoFilePathHash); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("source_zip_file_path", isoImage.SourceZipFilePath); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("source_zip_file_path_hash", isoImage.SourceZipFilePathHash); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("source_boot_file_path", isoImage.SourceBootFilePath); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("source_boot_file_path_hash", isoImage.SourceBootFilePathHash); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("destination_iso_file_path", destinationIsoFilePath); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("destination_zip_file_path", isoImage.DestinationZipFilePath); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("destination_boot_file_path", isoImage.DestinationBootFilePath); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("iso_media_type", api.IsoMediaType_name[isoImage.Media]); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("iso_file_system_type", api.IsoFileSystemType_name[isoImage.FileSystem]); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("volume_name", isoImage.VolumeName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("resolve_destination_iso_file_path", isoImage.ResolveDestinationIsoFilePath); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("resolve_destination_zip_file_path", isoImage.ResolveDestinationZipFilePath); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("resolve_destination_boot_file_path", isoImage.ResolveDestinationBootFilePath); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO][iso-image][read] read remote iso: %#v", d)

	return nil
}

func ensureFileStateUpdate(ctx context.Context, d *schema.ResourceData, c api.Client, name string) (string, bool, error) {
	sourceFilePathKey := fmt.Sprintf("source_%s_file_path", name)
	sourceFilePathHashKey := fmt.Sprintf("source_%s_file_path_hash", name)
	destinationFilePathKey := fmt.Sprintf("destination_%s_file_path", name)

	sourceFilePath := (d.Get(sourceFilePathKey)).(string)
	destinationFilePath := (d.Get(destinationFilePathKey)).(string)

	switch {
	case d.HasChange(sourceFilePathKey):
		o, n := d.GetChange(sourceFilePathKey)
		oldSourceFilePath := o.(string)
		newSourceFilePath := n.(string)

		oldDestinationFilePath := ""
		if d.HasChange(destinationFilePathKey) {
			od, _ := d.GetChange(destinationFilePathKey)
			oldDestinationFilePath = od.(string)
		} else {
			oldDestinationFilePath = destinationFilePath
		}
		resolveOldDestinationFilePath := oldDestinationFilePath
		if resolveOldDestinationFilePath == "" {
			resolveOldDestinationFilePath = winPath(filepath.Join(`$env:TEMP`, filepath.Base(oldSourceFilePath)))
		}

		newDestinationFilePath := ""
		if d.HasChange(destinationFilePathKey) {
			_, nd := d.GetChange(destinationFilePathKey)
			newDestinationFilePath = nd.(string)
		} else {
			newDestinationFilePath = destinationFilePath
		}
		resolveNewDestinationFilePath := newDestinationFilePath
		if resolveNewDestinationFilePath == "" {
			resolveNewDestinationFilePath = winPath(filepath.Join(`$env:TEMP`, filepath.Base(newSourceFilePath)))
		}

		switch {
		case newSourceFilePath == "":
			// must delete the old filename as we have removed the source one
			err := c.RemoteFileDelete(ctx, resolveOldDestinationFilePath)
			return resolveNewDestinationFilePath, true, err
		case oldSourceFilePath == "":
			// must upload file as we have set a new source one
			err := c.RemoteFileUpload(ctx, newSourceFilePath, resolveNewDestinationFilePath)
			return resolveNewDestinationFilePath, true, err
		case d.HasChange(sourceFilePathHashKey):
			// must upload file over existing one as hash has changed
			err := c.RemoteFileUpload(ctx, newSourceFilePath, resolveNewDestinationFilePath)
			return resolveNewDestinationFilePath, true, err
		}
	case sourceFilePath != "" && d.HasChange(destinationFilePathKey):
		o, n := d.GetChange(destinationFilePathKey)
		oldDestinationFilePath := o.(string)
		newDestinationFilePath := n.(string)

		resolveOldDestinationFilePath := oldDestinationFilePath
		if resolveOldDestinationFilePath == "" {
			resolveOldDestinationFilePath = winPath(filepath.Join(`$env:TEMP`, filepath.Base(sourceFilePath)))
		}

		resolveNewDestinationFilePath := newDestinationFilePath
		if resolveNewDestinationFilePath == "" {
			resolveNewDestinationFilePath = winPath(filepath.Join(`$env:TEMP`, filepath.Base(sourceFilePath)))
		}

		if resolveOldDestinationFilePath != resolveNewDestinationFilePath {
			// must delete the old filename as we have renamed it and we are not sure if any other properties have changed
			err := c.RemoteFileDelete(ctx, resolveOldDestinationFilePath)
			if err != nil {
				return resolveNewDestinationFilePath, true, err
			}

			// must upload file as we deleted the old one and we are expecting a file
			err = c.RemoteFileUpload(ctx, sourceFilePath, resolveNewDestinationFilePath)
			return resolveNewDestinationFilePath, true, err
		}
	case sourceFilePath != "" && d.HasChange(sourceFilePathHashKey):
		// must upload file over existing one as hash has changed
		resolveDestinationFilePath := destinationFilePath
		if resolveDestinationFilePath == "" {
			resolveDestinationFilePath = winPath(filepath.Join(`$env:TEMP`, filepath.Base(sourceFilePath)))
		}

		err := c.RemoteFileUpload(ctx, sourceFilePath, resolveDestinationFilePath)
		return resolveDestinationFilePath, true, err
	case sourceFilePath != "":
		resolveDestinationFilePath := destinationFilePath
		if resolveDestinationFilePath == "" {
			resolveDestinationFilePath = winPath(filepath.Join(`$env:TEMP`, filepath.Base(sourceFilePath)))
		}

		log.Printf("[INFO][iso-image][create] check if iso exists: %#v", resolveDestinationFilePath)
		resolveDestinationFilePathExists, err := c.RemoteFileExists(ctx, resolveDestinationFilePath)
		if err != nil {
			return resolveDestinationFilePath, true, fmt.Errorf("checking for existing %s: %+v", resolveDestinationFilePath, err)
		}

		if !resolveDestinationFilePathExists {
			// must upload file over as we are missing expected file
			err = c.RemoteFileUpload(ctx, sourceFilePath, resolveDestinationFilePath)
			return resolveDestinationFilePath, true, err
		}
	}

	resolveDestinationFilePath := destinationFilePath
	if resolveDestinationFilePath == "" {
		resolveDestinationFilePath = winPath(filepath.Join(`$env:TEMP`, filepath.Base(sourceFilePath)))
	}

	return resolveDestinationFilePath, false, nil
}

func resourceHyperVIsoImageUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var err error
	log.Printf("[INFO][iso-image][update] updating remote iso: %#v", d)
	c := meta.(api.Client)

	destinationIsoFilePath := d.Id()

	if destinationIsoFilePath != (d.Get("destination_iso_file_path")).(string) {
		return diag.FromErr(fmt.Errorf("cannot update destination_iso_file_path from %+v to %+v", destinationIsoFilePath, (d.Get("destination_iso_file_path")).(string)))
	}

	sourceIsoFilePath := (d.Get("source_iso_file_path")).(string)
	sourceIsoFilePathHash := (d.Get("source_iso_file_path_hash")).(string)
	sourceZipFilePath := (d.Get("source_zip_file_path")).(string)
	sourceZipFilePathHash := (d.Get("source_zip_file_path_hash")).(string)
	sourceBootFilePath := (d.Get("source_boot_file_path")).(string)
	sourceBootFilePathHash := (d.Get("source_boot_file_path_hash")).(string)
	destinationZipFilePath := (d.Get("destination_zip_file_path")).(string)
	destinationBootFilePath := (d.Get("destination_boot_file_path")).(string)
	media := api.ToIsoMediaType((d.Get("iso_media_type")).(string))
	fileSystem := api.ToIsoFileSystemType((d.Get("iso_file_system_type")).(string))
	volumeName := (d.Get("volume_name")).(string)

	resolveDestinationIsoFilePath, resolveDestinationIsoFilePathChanged, err := ensureFileStateUpdate(ctx, d, c, "iso")
	if err != nil {
		return diag.FromErr(err)
	}

	resolveDestinationZipFilePath, resolveDestinationZipFilePathChanged, err := ensureFileStateUpdate(ctx, d, c, "zip")
	if err != nil {
		return diag.FromErr(err)
	}

	resolveDestinationBootFilePath, resolveDestinationBootFilePathChanged, err := ensureFileStateUpdate(ctx, d, c, "boot")
	if err != nil {
		return diag.FromErr(err)
	}

	recreateDestinationIsoFilePath := false
	if !resolveDestinationIsoFilePathChanged && (resolveDestinationZipFilePathChanged || resolveDestinationBootFilePathChanged) || ((sourceZipFilePath != "" || sourceBootFilePath != "") && (d.HasChange("iso_media_type") || d.HasChange("iso_file_system_type") || d.HasChange("volume_name"))) {
		// must delete the iso file as we need to recreate it as the way it is created has changed
		err = c.RemoteFileDelete(ctx, resolveDestinationIsoFilePath)
		if err != nil {
			return diag.FromErr(err)
		}
		recreateDestinationIsoFilePath = true
	}

	if resolveDestinationIsoFilePathChanged || resolveDestinationZipFilePathChanged || resolveDestinationBootFilePathChanged || recreateDestinationIsoFilePath || d.HasChange("source_iso_file_path") || d.HasChange("source_iso_file_path_hash") || d.HasChange("source_zip_file_path") || d.HasChange("source_zip_file_path_hash") || d.HasChange("source_boot_file_path") || d.HasChange("source_boot_file_path_hash") || d.HasChange("destination_zip_file_path") || d.HasChange("destination_boot_file_path") {
		err = c.CreateOrUpdateIsoImage(ctx, sourceIsoFilePath, sourceIsoFilePathHash, sourceZipFilePath, sourceZipFilePathHash, sourceBootFilePath, sourceBootFilePathHash, destinationIsoFilePath, destinationZipFilePath, destinationBootFilePath, media, fileSystem, volumeName, resolveDestinationIsoFilePath, resolveDestinationZipFilePath, resolveDestinationBootFilePath)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	log.Printf("[INFO][iso-image][update] updated remote iso: %#v", d)

	return resourceHyperVIsoImageRead(ctx, d, meta)
}

func resourceHyperVIsoImageDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(api.Client)

	resolvedDestinationIsoFilePath := (d.Get("resolve_destination_iso_file_path")).(string)
	resolvedDestinationZipFilePath := (d.Get("resolve_destination_zip_file_path")).(string)
	resolvedDestinationBootFilePath := (d.Get("resolve_destination_boot_file_path")).(string)

	if resolvedDestinationIsoFilePath != "" {
		log.Printf("[INFO][iso-image][delete] deleting remote iso file: %#v", d)
		err := c.RemoteFileDelete(ctx, resolvedDestinationIsoFilePath)
		if err != nil {
			return diag.FromErr(err)
		}
		log.Printf("[INFO][iso-image][delete] deleted remote iso file: %#v", d)

		log.Printf("[INFO][iso-image][delete] deleting remote iso metadata file: %#v", d)
		err = c.RemoteFileDelete(ctx, fmt.Sprintf("%s.json", resolvedDestinationIsoFilePath))
		if err != nil {
			return diag.FromErr(err)
		}
		log.Printf("[INFO][iso-image][delete] deleted remote iso metadata file: %#v", d)
	}

	if resolvedDestinationZipFilePath != "" {
		log.Printf("[INFO][iso-image][delete] deleting remote zip file: %#v", d)
		err := c.RemoteFileDelete(ctx, resolvedDestinationZipFilePath)
		if err != nil {
			return diag.FromErr(err)
		}
		log.Printf("[INFO][iso-image][delete] deleted remote zip file: %#v", d)
	}

	if resolvedDestinationBootFilePath != "" {
		log.Printf("[INFO][iso-image][delete] deleting remote boot file: %#v", d)
		err := c.RemoteFileDelete(ctx, resolvedDestinationBootFilePath)
		if err != nil {
			return diag.FromErr(err)
		}
		log.Printf("[INFO][iso-image][delete] deleted remote boot file: %#v", d)
	}

	return nil
}
