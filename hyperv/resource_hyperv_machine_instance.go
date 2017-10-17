package hyperv

import (
	"fmt"
	"log"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/taliesins/terraform-provider-hyperv/api"
)

func resourceHyperVMachineInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceHyperVMachineInstanceCreate,
		Read:   resourceHyperVMachineInstanceRead,
		Update: resourceHyperVMachineInstanceUpdate,
		Delete: resourceHyperVMachineInstanceDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"generation": {
				Type:     schema.TypeInt,
				Optional: true,
				Default: 1,
				ForceNew: true,
			},

			"allow_unverified_paths": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"automatic_critical_error_action": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			"automatic_critical_error_action_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			"automatic_start_action": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			"automatic_start_delay": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			"automatic_stop_action": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			"checkpoint_type": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			"dynamic_memory": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"guest_controlled_cache_types": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"high_memory_mapped_io_space": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			"lock_on_disconnect": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			"low_memory_mapped_io_space": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			"memory_maximum_bytes": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			"memory_minimum_bytes": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			"memory_startup_bytes": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			"notes": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},

			"processor_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},

			"smart_paging_file_path": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},

			"snapshot_file_location": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},

			"static_memory": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceHyperVMachineInstanceCreate(d *schema.ResourceData, meta interface{}) (err error) {
	log.Printf("[INFO][hyperv] creating hyperv machine: %#v", d)
	c := meta.(*api.HypervClient)

	name := ""

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		return fmt.Errorf("name argument is required")
	}

	generation := (d.Get("generation")).(int)
	allowUnverifiedPaths := (d.Get("allow_unverified_paths")).(bool)
	automaticCriticalErrorAction := api.CriticalErrorAction((d.Get("automatic_critical_error_action")).(int))
	automaticCriticalErrorActionTimeout := int32((d.Get("automatic_critical_error_action_timeout")).(int))
	automaticStartAction := api.StartAction((d.Get("automatic_start_action")).(int))
	automaticStartDelay := int32((d.Get("automatic_start_delay")).(int))
	automaticStopAction := api.StopAction((d.Get("automatic_stop_action")).(int))
	checkpointType := api.CheckpointType((d.Get("checkpoint_type")).(int))
	dynamicMemory := (d.Get("dynamic_memory")).(bool)
	guestControlledCacheTypes := (d.Get("guest_controlled_cache_types")).(bool)
	highMemoryMappedIoSpace := int64((d.Get("high_memory_mapped_io_space")).(int))
	lockOnDisconnect := api.OnOffState((d.Get("lock_on_disconnect")).(int))
	lowMemoryMappedIoSpace := int32((d.Get("low_memory_mapped_io_space")).(int))
	memoryMaximumBytes := int64((d.Get("memory_maximum_bytes")).(int))
	memoryMinimumBytes := int64((d.Get("memory_minimum_bytes")).(int))
	memoryStartupBytes := int64((d.Get("memory_startup_bytes")).(int))
	notes := (d.Get("notes")).(string)
	processorCount := int64((d.Get("processor_count")).(int))
	smartPagingFilePath := (d.Get("smart_paging_file_path")).(string)
	snapshotFileLocation := (d.Get("snapshot_file_location")).(string)
	staticMemory := (d.Get("static_memory")).(bool)

	err = c.CreateVM(name, generation, allowUnverifiedPaths, automaticCriticalErrorAction, automaticCriticalErrorActionTimeout, automaticStartAction, automaticStartDelay, automaticStopAction, checkpointType, dynamicMemory, guestControlledCacheTypes, highMemoryMappedIoSpace, lockOnDisconnect, lowMemoryMappedIoSpace, memoryMaximumBytes, memoryMinimumBytes, memoryStartupBytes, notes, processorCount, smartPagingFilePath, snapshotFileLocation, staticMemory)

	if err != nil {
		return err
	}

	d.SetId(name)
	log.Printf("[INFO][hyperv] created hyperv machine: %#v", d)

	return  nil
}

func resourceHyperVMachineInstanceRead(d *schema.ResourceData, meta interface{}) (err error) {
	log.Printf("[INFO][hyperv] reading hyperv machine: %#v", d)
	c := meta.(*api.HypervClient)

	name := d.Id()

	s, err := c.GetVM(name)

	if err != nil {
		return err
	}

	if s.Name != name {
		d.SetId("")
		log.Printf("[INFO][hyperv] unable to read hyperv machine as it does not exist: %#v", name)
		return nil
	}

	d.Set("generation", s.Generation)
	d.Set("allow_unverified_paths", s.AllowUnverifiedPaths)
	d.Set("automatic_critical_error_action", s.AutomaticCriticalErrorAction)
	d.Set("automatic_critical_error_action_timeout", s.AutomaticCriticalErrorActionTimeout)
	d.Set("automatic_start_action", s.AutomaticStartAction)
	d.Set("automatic_start_delay", s.AutomaticStartDelay)
	d.Set("automatic_stop_action", s.AutomaticStopAction)
	d.Set("checkpoint_type", s.CheckpointType)
	d.Set("dynamic_memory", s.DynamicMemory)
	d.Set("guest_controlled_cache_types", s.GuestControlledCacheTypes)
	d.Set("high_memory_mapped_io_space", s.HighMemoryMappedIoSpace)
	d.Set("lock_on_disconnect", s.LockOnDisconnect)
	d.Set("low_memory_mapped_io_space", s.LowMemoryMappedIoSpace)
	d.Set("memory_maximum_bytes", s.MemoryMaximumBytes)
	d.Set("memory_minimum_bytes", s.MemoryMinimumBytes)
	d.Set("memory_startup_bytes", s.MemoryStartupBytes)
	d.Set("notes", s.Notes)
	d.Set("processor_count", s.ProcessorCount)
	d.Set("smart_paging_file_path", s.SmartPagingFilePath)
	d.Set("snapshot_file_location", s.SnapshotFileLocation)
	d.Set("static_memory", s.StaticMemory)

	if err != nil {
		return err
	}

	log.Printf("[INFO][hyperv] read hyperv machine: %#v", d)

	return nil
}

func resourceHyperVMachineInstanceUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	log.Printf("[INFO][hyperv] updating hyperv machine: %#v", d)
	c := meta.(*api.HypervClient)

	name := d.Id()

	//generation := (d.Get("generation")).(int)
	allowUnverifiedPaths := (d.Get("allow_unverified_paths")).(bool)
	automaticCriticalErrorAction := api.CriticalErrorAction((d.Get("automatic_critical_error_action")).(int))
	automaticCriticalErrorActionTimeout := int32((d.Get("automatic_critical_error_action_timeout")).(int))
	automaticStartAction := api.StartAction((d.Get("automatic_start_action")).(int))
	automaticStartDelay := int32((d.Get("automatic_start_delay")).(int))
	automaticStopAction := api.StopAction((d.Get("automatic_stop_action")).(int))
	checkpointType := api.CheckpointType((d.Get("checkpoint_type")).(int))
	dynamicMemory := (d.Get("dynamic_memory")).(bool)
	guestControlledCacheTypes := (d.Get("guest_controlled_cache_types")).(bool)
	highMemoryMappedIoSpace := int64((d.Get("high_memory_mapped_io_space")).(int))
	lockOnDisconnect := api.OnOffState((d.Get("lock_on_disconnect")).(int))
	lowMemoryMappedIoSpace := int32((d.Get("low_memory_mapped_io_space")).(int))
	memoryMaximumBytes := int64((d.Get("memory_maximum_bytes")).(int))
	memoryMinimumBytes := int64((d.Get("memory_minimum_bytes")).(int))
	memoryStartupBytes := int64((d.Get("memory_startup_bytes")).(int))
	notes := (d.Get("notes")).(string)
	processorCount := int64((d.Get("processor_count")).(int))
	smartPagingFilePath := (d.Get("smart_paging_file_path")).(string)
	snapshotFileLocation := (d.Get("snapshot_file_location")).(string)
	staticMemory := (d.Get("static_memory")).(bool)

	err = c.UpdateVM(name, allowUnverifiedPaths, automaticCriticalErrorAction, automaticCriticalErrorActionTimeout, automaticStartAction, automaticStartDelay, automaticStopAction, checkpointType, dynamicMemory, guestControlledCacheTypes, highMemoryMappedIoSpace, lockOnDisconnect, lowMemoryMappedIoSpace, memoryMaximumBytes, memoryMinimumBytes, memoryStartupBytes, notes, processorCount, smartPagingFilePath, snapshotFileLocation, staticMemory)

	if err != nil {
		return err
	}

	log.Printf("[INFO][hyperv] updated hyperv machine: %#v", d)

	return nil
}

func resourceHyperVMachineInstanceDelete(d *schema.ResourceData, meta interface{}) (err error) {
	log.Printf("[INFO][hyperv] deleting hyperv machine: %#v", d)

	c := meta.(*api.HypervClient)

	name := d.Id()

	err = c.DeleteVM(name)

	if err != nil {
		return err
	}

	log.Printf("[INFO][hyperv] deleted hyperv machine: %#v", d)
	return nil
}