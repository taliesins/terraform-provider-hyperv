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
				Default:  1,
				ForceNew: true,
			},

			"automatic_critical_error_action": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      api.CriticalErrorAction_name[api.CriticalErrorAction_Pause],
				ValidateFunc: stringKeyInMap(api.CriticalErrorAction_value, true),
			},

			"automatic_critical_error_action_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
			},

			"automatic_start_action": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      api.StartAction_name[api.StartAction_StartIfRunning],
				ValidateFunc: stringKeyInMap(api.StartAction_value, true),
			},

			"automatic_start_delay": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			"automatic_stop_action": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      api.StopAction_name[api.StopAction_Save],
				ValidateFunc: stringKeyInMap(api.StopAction_value, true),
			},

			"checkpoint_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      api.CheckpointType_name[api.CheckpointType_Production],
				ValidateFunc: stringKeyInMap(api.CheckpointType_value, true),
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
				Default:  536870912,
			},

			"lock_on_disconnect": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      api.OnOffState_name[api.OnOffState_Off],
				ValidateFunc: stringKeyInMap(api.OnOffState_value, true),
			},

			"low_memory_mapped_io_space": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  134217728,
			},

			"memory_maximum_bytes": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1099511627776,
			},

			"memory_minimum_bytes": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  536870912,
			},

			"memory_startup_bytes": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  536870912,
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
				Default:  `C:\ProgramData\Microsoft\Windows\Hyper-V`,
			},

			"snapshot_file_location": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  `C:\ProgramData\Microsoft\Windows\Hyper-V`,
			},

			"static_memory": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceHyperVMachineInstanceCreate(d *schema.ResourceData, meta interface{}) (err error) {
	log.Printf("[INFO][hyperv][create] creating hyperv machine: %#v", d)
	c := meta.(*api.HypervClient)

	name := ""

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		return fmt.Errorf("[ERROR][hyperv][create] name argument is required")
	}

	generation := (d.Get("generation")).(int)
	automaticCriticalErrorAction := api.ToCriticalErrorAction((d.Get("automatic_critical_error_action")).(string))
	automaticCriticalErrorActionTimeout := int32((d.Get("automatic_critical_error_action_timeout")).(int))
	automaticStartAction := api.ToStartAction((d.Get("automatic_start_action")).(string))
	automaticStartDelay := int32((d.Get("automatic_start_delay")).(int))
	automaticStopAction := api.ToStopAction((d.Get("automatic_stop_action")).(string))
	checkpointType := api.ToCheckpointType((d.Get("checkpoint_type")).(string))
	dynamicMemory := (d.Get("dynamic_memory")).(bool)
	guestControlledCacheTypes := (d.Get("guest_controlled_cache_types")).(bool)
	highMemoryMappedIoSpace := int64((d.Get("high_memory_mapped_io_space")).(int))
	lockOnDisconnect := api.ToOnOffState((d.Get("lock_on_disconnect")).(string))
	lowMemoryMappedIoSpace := int32((d.Get("low_memory_mapped_io_space")).(int))
	memoryMaximumBytes := int64((d.Get("memory_maximum_bytes")).(int))
	memoryMinimumBytes := int64((d.Get("memory_minimum_bytes")).(int))
	memoryStartupBytes := int64((d.Get("memory_startup_bytes")).(int))
	notes := (d.Get("notes")).(string)
	processorCount := int64((d.Get("processor_count")).(int))
	smartPagingFilePath := (d.Get("smart_paging_file_path")).(string)
	snapshotFileLocation := (d.Get("snapshot_file_location")).(string)
	staticMemory := (d.Get("static_memory")).(bool)

	if dynamicMemory && staticMemory {
		return fmt.Errorf("[ERROR][hyperv][create] Dynamic and static can't be both selected at the same time")
	}

	if !dynamicMemory && !staticMemory {
		return fmt.Errorf("[ERROR][hyperv][create] Either dynamic or static must be selected")
	}

	err = c.CreateVM(name, generation, automaticCriticalErrorAction, automaticCriticalErrorActionTimeout, automaticStartAction, automaticStartDelay, automaticStopAction, checkpointType, dynamicMemory, guestControlledCacheTypes, highMemoryMappedIoSpace, lockOnDisconnect, lowMemoryMappedIoSpace, memoryMaximumBytes, memoryMinimumBytes, memoryStartupBytes, notes, processorCount, smartPagingFilePath, snapshotFileLocation, staticMemory)

	if err != nil {
		return err
	}

	d.SetId(name)
	log.Printf("[INFO][hyperv][create] created hyperv machine: %#v", d)

	return nil
}

func resourceHyperVMachineInstanceRead(d *schema.ResourceData, meta interface{}) (err error) {
	log.Printf("[INFO][hyperv][read] reading hyperv machine: %#v", d)
	c := meta.(*api.HypervClient)

	name := d.Id()

	s, err := c.GetVM(name)

	if err != nil {
		return err
	}

	log.Printf("[INFO][hyperv][read] retrieved vm: %+v", s)

	if s.Name != name {
		d.SetId("")
		log.Printf("[INFO][hyperv][read] unable to read hyperv machine as it does not exist: %#v", name)
		return nil
	}

	d.Set("generation", s.Generation)
	d.Set("automatic_critical_error_action", s.AutomaticCriticalErrorAction.String())
	d.Set("automatic_critical_error_action_timeout", s.AutomaticCriticalErrorActionTimeout)
	d.Set("automatic_start_action", s.AutomaticStartAction.String())
	d.Set("automatic_start_delay", s.AutomaticStartDelay)
	d.Set("automatic_stop_action", s.AutomaticStopAction.String())
	d.Set("checkpoint_type", s.CheckpointType.String())
	d.Set("dynamic_memory", s.DynamicMemory)
	d.Set("guest_controlled_cache_types", s.GuestControlledCacheTypes)
	d.Set("high_memory_mapped_io_space", s.HighMemoryMappedIoSpace)
	d.Set("lock_on_disconnect", s.LockOnDisconnect.String())
	d.Set("low_memory_mapped_io_space", s.LowMemoryMappedIoSpace)
	d.Set("memory_maximum_bytes", s.MemoryMaximumBytes)
	d.Set("memory_minimum_bytes", s.MemoryMinimumBytes)
	d.Set("memory_startup_bytes", s.MemoryStartupBytes)
	d.Set("notes", s.Notes)
	d.Set("processor_count", s.ProcessorCount)
	d.Set("smart_paging_file_path", s.SmartPagingFilePath)
	d.Set("snapshot_file_location", s.SnapshotFileLocation)
	d.Set("static_memory", s.StaticMemory)

	if s.DynamicMemory && s.StaticMemory {
		return fmt.Errorf("[ERROR][hyperv][read] Dynamic and static can't be both selected at the same time")
	}

	if !s.DynamicMemory && !s.StaticMemory {
		return fmt.Errorf("[ERROR][hyperv][read] Either dynamic or static must be selected")
	}

	if err != nil {
		return err
	}

	log.Printf("[INFO][hyperv][read] read hyperv machine: %#v", d)

	return nil
}

func resourceHyperVMachineInstanceUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	log.Printf("[INFO][hyperv][update] updating hyperv machine: %#v", d)
	c := meta.(*api.HypervClient)

	name := d.Id()

	//generation := (d.Get("generation")).(int)
	automaticCriticalErrorAction := api.ToCriticalErrorAction((d.Get("automatic_critical_error_action")).(string))
	automaticCriticalErrorActionTimeout := int32((d.Get("automatic_critical_error_action_timeout")).(int))
	automaticStartAction := api.ToStartAction((d.Get("automatic_start_action")).(string))
	automaticStartDelay := int32((d.Get("automatic_start_delay")).(int))
	automaticStopAction := api.ToStopAction((d.Get("automatic_stop_action")).(string))
	checkpointType := api.ToCheckpointType((d.Get("checkpoint_type")).(string))
	dynamicMemory := (d.Get("dynamic_memory")).(bool)
	guestControlledCacheTypes := (d.Get("guest_controlled_cache_types")).(bool)
	highMemoryMappedIoSpace := int64((d.Get("high_memory_mapped_io_space")).(int))
	lockOnDisconnect := api.ToOnOffState((d.Get("lock_on_disconnect")).(string))
	lowMemoryMappedIoSpace := int32((d.Get("low_memory_mapped_io_space")).(int))
	memoryMaximumBytes := int64((d.Get("memory_maximum_bytes")).(int))
	memoryMinimumBytes := int64((d.Get("memory_minimum_bytes")).(int))
	memoryStartupBytes := int64((d.Get("memory_startup_bytes")).(int))
	notes := (d.Get("notes")).(string)
	processorCount := int64((d.Get("processor_count")).(int))
	smartPagingFilePath := (d.Get("smart_paging_file_path")).(string)
	snapshotFileLocation := (d.Get("snapshot_file_location")).(string)
	staticMemory := (d.Get("static_memory")).(bool)

	if dynamicMemory && staticMemory {
		return fmt.Errorf("[ERROR][hyperv][update] Dynamic and static can't be both selected at the same time")
	}

	if !dynamicMemory && !staticMemory {
		return fmt.Errorf("[ERROR][hyperv][update] Either dynamic or static must be selected")
	}

	err = c.UpdateVM(name, automaticCriticalErrorAction, automaticCriticalErrorActionTimeout, automaticStartAction, automaticStartDelay, automaticStopAction, checkpointType, dynamicMemory, guestControlledCacheTypes, highMemoryMappedIoSpace, lockOnDisconnect, lowMemoryMappedIoSpace, memoryMaximumBytes, memoryMinimumBytes, memoryStartupBytes, notes, processorCount, smartPagingFilePath, snapshotFileLocation, staticMemory)

	if err != nil {
		return err
	}

	log.Printf("[INFO][hyperv][update] updated hyperv machine: %#v", d)

	return nil
}

func resourceHyperVMachineInstanceDelete(d *schema.ResourceData, meta interface{}) (err error) {
	log.Printf("[INFO][hyperv][delete] deleting hyperv machine: %#v", d)

	c := meta.(*api.HypervClient)

	name := d.Id()

	err = c.DeleteVM(name)

	if err != nil {
		return err
	}

	log.Printf("[INFO][hyperv][delete] deleted hyperv machine: %#v", d)
	return nil
}
