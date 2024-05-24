package api

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ExpandGpuAdapters(d *schema.ResourceData) ([]VmGpuAdapter, error) {
	expandedGpuAdapters := make([]VmGpuAdapter, 0)

	if v, ok := d.GetOk("gpu_adapters"); ok {
		gpuAdapters := v.([]interface{})
		for _, gpuAdapter := range gpuAdapters {
			gpuAdapter, ok := gpuAdapter.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("[ERROR][hyperv] gpu_adapters should be a Hash - was '%+v'", gpuAdapter)
			}

			log.Printf("[DEBUG] gpuAdapter =  [%+v]", gpuAdapter)

			min_partition_encode, _ := strconv.ParseInt(gpuAdapter["min_partition_encode"].(string), 10, 64)
			max_partition_encode, _ := strconv.ParseInt(gpuAdapter["max_partition_encode"].(string), 10, 64)
			optimal_partition_encode, _ := strconv.ParseInt(gpuAdapter["optimal_partition_encode"].(string), 10, 64)

			expandedGpuAdapter := VmGpuAdapter{
				InstancePath:            gpuAdapter["device_path_name"].(string),
				MinPartitionVram:        int32(gpuAdapter["min_partition_vram"].(int)),
				MaxPartitionVram:        int32(gpuAdapter["max_partition_vram"].(int)),
				OptimalPartitionVram:    int32(gpuAdapter["optimal_partition_vram"].(int)),
				MinPartitionEncode:      int64(min_partition_encode),
				MaxPartitionEncode:      int64(max_partition_encode),
				OptimalPartitionEncode:  int64(optimal_partition_encode),
				MinPartitionDecode:      int32(gpuAdapter["min_partition_decode"].(int)),
				MaxPartitionDecode:      int32(gpuAdapter["max_partition_decode"].(int)),
				OptimalPartitionDecode:  int32(gpuAdapter["optimal_partition_decode"].(int)),
				MinPartitionCompute:     int32(gpuAdapter["min_partition_compute"].(int)),
				MaxPartitionCompute:     int32(gpuAdapter["max_partition_compute"].(int)),
				OptimalPartitionCompute: int32(gpuAdapter["optimal_partition_compute"].(int)),
			}

			expandedGpuAdapters = append(expandedGpuAdapters, expandedGpuAdapter)
		}
	}

	return expandedGpuAdapters, nil
}

func FlattenGpuAdapters(gpuAdapters *[]VmGpuAdapter) []interface{} {
	if gpuAdapters == nil || len(*gpuAdapters) < 1 {
		return nil
	}

	flattenedGpuAdapters := make([]interface{}, 0)

	for _, gpuAdapter := range *gpuAdapters {
		flattenedGpuAdapter := make(map[string]interface{})

		min_partition_encode := strconv.Itoa(int(gpuAdapter.MinPartitionEncode))
		max_partition_encode := strconv.Itoa(int(gpuAdapter.MaxPartitionEncode))
		optimal_partition_encode := strconv.Itoa(int(gpuAdapter.OptimalPartitionEncode))

		flattenedGpuAdapter["device_path_name"] = gpuAdapter.InstancePath
		flattenedGpuAdapter["min_partition_vram"] = gpuAdapter.MinPartitionVram
		flattenedGpuAdapter["max_partition_vram"] = gpuAdapter.MaxPartitionVram
		flattenedGpuAdapter["optimal_partition_vram"] = gpuAdapter.OptimalPartitionVram
		flattenedGpuAdapter["min_partition_encode"] = min_partition_encode
		flattenedGpuAdapter["max_partition_encode"] = max_partition_encode
		flattenedGpuAdapter["optimal_partition_encode"] = optimal_partition_encode
		flattenedGpuAdapter["min_partition_decode"] = gpuAdapter.MinPartitionDecode
		flattenedGpuAdapter["max_partition_decode"] = gpuAdapter.MaxPartitionDecode
		flattenedGpuAdapter["optimal_partition_decode"] = gpuAdapter.OptimalPartitionDecode
		flattenedGpuAdapter["min_partition_compute"] = gpuAdapter.MinPartitionCompute
		flattenedGpuAdapter["max_partition_compute"] = gpuAdapter.MaxPartitionCompute
		flattenedGpuAdapter["optimal_partition_compute"] = gpuAdapter.OptimalPartitionCompute

		flattenedGpuAdapters = append(flattenedGpuAdapters, flattenedGpuAdapter)
	}

	return flattenedGpuAdapters
}

type VmGpuAdapter struct {
	VmName                  string
	InstancePath            string
	MinPartitionVram        int32
	MaxPartitionVram        int32
	OptimalPartitionVram    int32
	MinPartitionEncode      int64
	MaxPartitionEncode      int64
	OptimalPartitionEncode  int64
	MinPartitionDecode      int32
	MaxPartitionDecode      int32
	OptimalPartitionDecode  int32
	MinPartitionCompute     int32
	MaxPartitionCompute     int32
	OptimalPartitionCompute int32
}

type HypervGpuAdapterClient interface {
	CreateVmGpuAdapter(
		ctx context.Context,
		vmName string,
		devicePathName string,
		minPartitionVram int32,
		maxPartitionVram int32,
		optimalPartitionVram int32,
		minPartitionEncode int64,
		maxPartitionEncode int64,
		optimalPartitionEncode int64,
		minPartitionDecode int32,
		maxPartitionDecode int32,
		optimalPartitionDecode int32,
		minPartitionCompute int32,
		maxPartitionCompute int32,
		optimalPartitionCompute int32,
	) (err error)
	GetVmGpuAdapters(ctx context.Context, vmName string) (result []VmGpuAdapter, err error)
	UpdateVmGpuAdapter(
		ctx context.Context,
		vmName string,
		instancePath string,
		minPartitionVram int32,
		maxPartitionVram int32,
		optimalPartitionVram int32,
		minPartitionEncode int64,
		maxPartitionEncode int64,
		optimalPartitionEncode int64,
		minPartitionDecode int32,
		maxPartitionDecode int32,
		optimalPartitionDecode int32,
		minPartitionCompute int32,
		maxPartitionCompute int32,
		optimalPartitionCompute int32,
	) (err error)
	DeleteVmGpuAdapter(ctx context.Context, vmName string, instancePath string) (err error)
	CreateOrUpdateVmGpuAdapters(ctx context.Context, vmName string, gpuAdapters []VmGpuAdapter) (err error)
}
