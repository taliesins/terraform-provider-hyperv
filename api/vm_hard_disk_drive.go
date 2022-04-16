package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"path/filepath"
	"strconv"
	"strings"
)

type ControllerType int

const (
	ControllerType_Ide  ControllerType = 0
	ControllerType_Scsi ControllerType = 1
)

var ControllerType_name = map[ControllerType]string{
	ControllerType_Ide:  "Ide",
	ControllerType_Scsi: "Scsi",
}

var ControllerType_value = map[string]ControllerType{
	"ide":  ControllerType_Ide,
	"scsi": ControllerType_Scsi,
}

func (x ControllerType) String() string {
	return ControllerType_name[x]
}

func ToControllerType(x string) ControllerType {
	if integerValue, err := strconv.Atoi(x); err == nil {
		return ControllerType(integerValue)
	}
	return ControllerType_value[strings.ToLower(x)]
}

func (d *ControllerType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(d.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (d *ControllerType) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		var i int
		err2 := json.Unmarshal(b, &i)
		if err2 == nil {
			*d = ControllerType(i)
			return nil
		}

		return err
	}
	*d = ToControllerType(s)
	return nil
}

type CacheAttributes int

const (
	CacheAttributes_Default                 CacheAttributes = 0
	CacheAttributes_WriteCacheEnabled       CacheAttributes = 1
	CacheAttributes_WriteCacheAndFUAEnabled CacheAttributes = 2
	CacheAttributes_WriteCacheDisabled      CacheAttributes = 3
)

var CacheAttributes_name = map[CacheAttributes]string{
	CacheAttributes_Default:                 "Default",
	CacheAttributes_WriteCacheEnabled:       "WriteCacheEnabled",
	CacheAttributes_WriteCacheAndFUAEnabled: "WriteCacheAndFUAEnabled",
	CacheAttributes_WriteCacheDisabled:      "WriteCacheDisabled",
}

var CacheAttributes_value = map[string]CacheAttributes{
	"default":                 CacheAttributes_Default,
	"writecacheenabled":       CacheAttributes_WriteCacheEnabled,
	"writecacheandfuaenabled": CacheAttributes_WriteCacheAndFUAEnabled,
	"writecachedisabled":      CacheAttributes_WriteCacheDisabled,
}

func (x CacheAttributes) String() string {
	return CacheAttributes_name[x]
}

func ToCacheAttributes(x string) CacheAttributes {
	if integerValue, err := strconv.Atoi(x); err == nil {
		return CacheAttributes(integerValue)
	}
	return CacheAttributes_value[strings.ToLower(x)]
}

func (d *CacheAttributes) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(d.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (d *CacheAttributes) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		var i int
		err2 := json.Unmarshal(b, &i)
		if err2 == nil {
			*d = CacheAttributes(i)
			return nil
		}

		return err
	}
	*d = ToCacheAttributes(s)
	return nil
}

func DiffSuppressVmHardDiskPath(key, old, new string, d *schema.ResourceData) bool {
	log.Printf("[DEBUG] '[%s]' Comparing old value '[%v]' with new value '[%v]' ", key, old, new)
	if new == "" {
		//We have not explicitly set a value, so allow any value as we are not tracking it
		return true
	}

	if new == old {
		return true
	}

	//Ignore snapshots otherwise it will change from "c:\\vhdx\\web_server_g2_B63C9D15-F9A3-4F63-A896-FFD80BC7754C.avhdx" -> "c:\\vhdx\\web_server_g2.vhdx"
	oldExtension := strings.ToLower(filepath.Ext(old))
	newExtension := strings.ToLower(filepath.Ext(new))
	if oldExtension == ".avhdx" && newExtension == ".vhdx" {
		newName := new[0 : len(new)-len(newExtension)]
		return strings.HasPrefix(old, newName+"_")
	}

	return false
}

func ExpandHardDiskDrives(d *schema.ResourceData) ([]VmHardDiskDrive, error) {
	expandedHardDiskDrives := make([]VmHardDiskDrive, 0)

	if v, ok := d.GetOk("hard_disk_drives"); ok {
		hardDiskDrives := v.([]interface{})

		for _, hardDiskDrive := range hardDiskDrives {
			hardDiskDrive, ok := hardDiskDrive.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] hard_disk_drives should be a Hash - was '%+v'", hardDiskDrive)
			}

			expandedHardDiskDrive := VmHardDiskDrive{
				ControllerType:                ToControllerType(hardDiskDrive["controller_type"].(string)),
				ControllerNumber:              int32(hardDiskDrive["controller_number"].(int)),
				ControllerLocation:            int32(hardDiskDrive["controller_location"].(int)),
				Path:                          hardDiskDrive["path"].(string),
				DiskNumber:                    uint32(hardDiskDrive["disk_number"].(int)),
				ResourcePoolName:              hardDiskDrive["resource_pool_name"].(string),
				SupportPersistentReservations: hardDiskDrive["support_persistent_reservations"].(bool),
				MaximumIops:                   uint64(hardDiskDrive["maximum_iops"].(int)),
				MinimumIops:                   uint64(hardDiskDrive["minimum_iops"].(int)),
				QosPolicyId:                   hardDiskDrive["qos_policy_id"].(string),
				OverrideCacheAttributes:       ToCacheAttributes(hardDiskDrive["override_cache_attributes"].(string)),
			}

			expandedHardDiskDrives = append(expandedHardDiskDrives, expandedHardDiskDrive)
		}
	}

	return expandedHardDiskDrives, nil
}

func FlattenHardDiskDrives(hardDiskDrives *[]VmHardDiskDrive) []interface{} {
	if hardDiskDrives == nil || len(*hardDiskDrives) < 1 {
		return nil
	}

	flattenedHardDiskDrives := make([]interface{}, 0)

	for _, hardDiskDrive := range *hardDiskDrives {
		flattenedHardDiskDrive := make(map[string]interface{})
		flattenedHardDiskDrive["controller_type"] = hardDiskDrive.ControllerType.String()
		flattenedHardDiskDrive["controller_number"] = hardDiskDrive.ControllerNumber
		flattenedHardDiskDrive["controller_location"] = hardDiskDrive.ControllerLocation
		flattenedHardDiskDrive["path"] = hardDiskDrive.Path
		flattenedHardDiskDrive["disk_number"] = hardDiskDrive.DiskNumber
		flattenedHardDiskDrive["resource_pool_name"] = hardDiskDrive.ResourcePoolName
		flattenedHardDiskDrive["support_persistent_reservations"] = hardDiskDrive.SupportPersistentReservations
		flattenedHardDiskDrive["maximum_iops"] = hardDiskDrive.MaximumIops
		flattenedHardDiskDrive["minimum_iops"] = hardDiskDrive.MinimumIops
		flattenedHardDiskDrive["qos_policy_id"] = hardDiskDrive.QosPolicyId
		flattenedHardDiskDrive["override_cache_attributes"] = hardDiskDrive.OverrideCacheAttributes.String()
		flattenedHardDiskDrives = append(flattenedHardDiskDrives, flattenedHardDiskDrive)
	}
	
	return flattenedHardDiskDrives
}

type VmHardDiskDrive struct {
	VmName                        string
	ControllerType                ControllerType
	ControllerNumber              int32
	ControllerLocation            int32
	Path                          string
	DiskNumber                    uint32
	ResourcePoolName              string
	SupportPersistentReservations bool
	MaximumIops                   uint64
	MinimumIops                   uint64
	QosPolicyId                   string
	OverrideCacheAttributes       CacheAttributes
	//AllowUnverifiedPaths          bool no way of checking if its turned on so always turn on
}

type HypervVmHardDiskDriveClient interface {
	CreateVmHardDiskDrive(
		vmName string,
		controllerType ControllerType,
		controllerNumber int32,
		controllerLocation int32,
		path string,
		diskNumber uint32,
		resourcePoolName string,
		supportPersistentReservations bool,
		maximumIops uint64,
		minimumIops uint64,
		qosPolicyId string,
		overrideCacheAttributes CacheAttributes,

	) (err error)
	GetVmHardDiskDrives(vmName string) (result []VmHardDiskDrive, err error)
	UpdateVmHardDiskDrive(
		vmName string,
		controllerNumber int32,
		controllerLocation int32,
		controllerType ControllerType,
		toControllerNumber int32,
		toControllerLocation int32,
		path string,
		diskNumber uint32,
		resourcePoolName string,
		supportPersistentReservations bool,
		maximumIops uint64,
		minimumIops uint64,
		qosPolicyId string,
		overrideCacheAttributes CacheAttributes,
	) (err error)
	DeleteVmHardDiskDrive(vmname string, controllerNumber int32, controllerLocation int32) (err error)
	CreateOrUpdateVmHardDiskDrives(vmName string, hardDiskDrives []VmHardDiskDrive) (err error)
}
