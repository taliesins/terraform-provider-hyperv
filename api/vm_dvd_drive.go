package api

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ExpandDvdDrives(d *schema.ResourceData) ([]VmDvdDrive, error) {
	expandedDvdDrives := make([]VmDvdDrive, 0)

	if v, ok := d.GetOk("dvd_drives"); ok {
		dvdDrives := v.([]interface{})
		for _, dvdDrive := range dvdDrives {
			dvdDrive, ok := dvdDrive.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] dvd_drives should be a Hash - was '%+v'", dvdDrive)
			}

			expandedDvdDrive := VmDvdDrive{
				ControllerNumber:   dvdDrive["controller_number"].(int),
				ControllerLocation: dvdDrive["controller_location"].(int),
				Path:               dvdDrive["path"].(string),
				ResourcePoolName:   dvdDrive["resource_pool_name"].(string),
			}

			expandedDvdDrives = append(expandedDvdDrives, expandedDvdDrive)
		}
	}

	return expandedDvdDrives, nil
}

func FlattenDvdDrives(dvdDrives *[]VmDvdDrive) []interface{} {
	if dvdDrives == nil || len(*dvdDrives) < 1 {
		return nil
	}

	flattenedDvdDrives := make([]interface{}, 0)

	for _, dvdDrive := range *dvdDrives {
		flattenedDvdDrive := make(map[string]interface{})
		flattenedDvdDrive["controller_number"] = dvdDrive.ControllerNumber
		flattenedDvdDrive["controller_location"] = dvdDrive.ControllerLocation
		flattenedDvdDrive["path"] = dvdDrive.Path
		flattenedDvdDrive["resource_pool_name"] = dvdDrive.ResourcePoolName
		flattenedDvdDrives = append(flattenedDvdDrives, flattenedDvdDrive)
	}

	return flattenedDvdDrives
}

type VmDvdDrive struct {
	VmName             string
	ControllerNumber   int
	ControllerLocation int
	Path               string
	// AllowUnverifiedPaths bool no way of checking if its turned on so always turn on
	ResourcePoolName string
}

type HypervVmDvdDriveClient interface {
	CreateVmDvdDrive(
		ctx context.Context,
		vmName string,
		controllerNumber int,
		controllerLocation int,
		path string,
		resourcePoolName string,
	) (err error)
	GetVmDvdDrives(ctx context.Context, vmName string) (result []VmDvdDrive, err error)
	UpdateVmDvdDrive(
		ctx context.Context,
		vmName string,
		controllerNumber int,
		controllerLocation int,
		toControllerNumber int,
		toControllerLocation int,
		path string,
		resourcePoolName string,
	) (err error)
	DeleteVmDvdDrive(ctx context.Context, vmName string, controllerNumber int, controllerLocation int) (err error)
	CreateOrUpdateVmDvdDrives(ctx context.Context, vmName string, dvdDrives []VmDvdDrive) (err error)
}
