package hyperv

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/taliesins/terraform-provider-hyperv/api"
)

func dataSourceHyperVMachineInstance() *schema.Resource {
	return &schema.Resource{
		Read:   resourceHyperVMachineInstanceRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"generation": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      2,
				ValidateFunc: IntInSlice([]int{1, 2}),
				ForceNew:     true,
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

			"integration_services": {
				Type:             schema.TypeMap,
				Optional:         true,
				DefaultFunc:      api.DefaultVmIntegrationServices,
				DiffSuppressFunc: api.DiffSuppressVmIntegrationServices,
				Elem:             schema.TypeBool,
			},

			"state": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      api.VmState_name[api.VmState_Running],
				ValidateFunc: stringKeyInMap(api.VmState_SettableValue, true),
			},

			"wait_for_state_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  120,
			},

			"wait_for_state_poll_period": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  2,
			},

			"wait_for_ips_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  300,
			},

			"wait_for_ips_poll_period": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5,
			},

			"vm_firmware": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				DefaultFunc: api.DefaultVmFirmwares,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_secure_boot": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      api.OnOffState_name[api.OnOffState_On],
							ValidateFunc: stringKeyInMap(api.OnOffState_value, true),
						},

						"secure_boot_template": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "MicrosoftWindows",
						},

						"preferred_network_boot_protocol": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      api.IPProtocolPreference_name[api.IPProtocolPreference_IPv4],
							ValidateFunc: stringKeyInMap(api.IPProtocolPreference_value, true),
						},

						"console_mode": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      api.ConsoleModeType_name[api.ConsoleModeType_Default],
							ValidateFunc: stringKeyInMap(api.ConsoleModeType_value, true),
						},

						"pause_after_boot_failure": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      api.OnOffState_name[api.OnOffState_Off],
							ValidateFunc: stringKeyInMap(api.OnOffState_value, true),
						},
					},
				},
			},

			"vm_processor": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				DefaultFunc: api.DefaultVmProcessors,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"compatibility_for_migration_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},

						"compatibility_for_older_operating_systems_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},

						"hw_thread_count_per_core": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},

						"maximum": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      100,
							ValidateFunc: ValueOrIntBetween(100, 0, 100),
						},

						"reserve": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      0,
							ValidateFunc: ValueOrIntBetween(0, 0, 100),
						},

						"relative_weight": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      100,
							ValidateFunc: ValueOrIntBetween(100, 0, 10000),
						},

						"maximum_count_per_numa_node": {
							Type:             schema.TypeInt,
							Optional:         true,
							Default:          0, //Dynamic value
							DiffSuppressFunc: api.DiffSuppressVmProcessorMaximumCountPerNumaNode,
						},

						"maximum_count_per_numa_socket": {
							Type:             schema.TypeInt,
							Optional:         true,
							Default:          0, //Dynamic value
							DiffSuppressFunc: api.DiffSuppressVmProcessorMaximumCountPerNumaSocket,
						},

						"enable_host_resource_protection": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},

						"expose_virtualization_extensions": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},

			"network_adaptors": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"switch_name": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
							ForceNew: false,
						},
						"management_os": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"is_legacy": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
							ForceNew: true,
						},
						"dynamic_mac_address": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"static_mac_address": {
							Type:             schema.TypeString,
							Optional:         true,
							Default:          "",
							DiffSuppressFunc: api.DiffSuppressVmStaticMacAddress,
						},
						"mac_address_spoofing": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      api.OnOffState_name[api.OnOffState_Off],
							ValidateFunc: stringKeyInMap(api.OnOffState_value, true),
						},
						"dhcp_guard": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      api.OnOffState_name[api.OnOffState_Off],
							ValidateFunc: stringKeyInMap(api.OnOffState_value, true),
						},
						"router_guard": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      api.OnOffState_name[api.OnOffState_Off],
							ValidateFunc: stringKeyInMap(api.OnOffState_value, true),
						},
						"port_mirroring": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      api.PortMirroring_name[api.PortMirroring_None],
							ValidateFunc: stringKeyInMap(api.PortMirroring_value, true),
						},
						"ieee_priority_tag": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      api.OnOffState_name[api.OnOffState_Off],
							ValidateFunc: stringKeyInMap(api.OnOffState_value, true),
						},
						"vmq_weight": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      100,
							ValidateFunc: validation.IntBetween(0, 100),
						},
						"iov_queue_pairs_requested": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      1,
							ValidateFunc: validation.IntBetween(1, 4294967295),
						},
						"iov_interrupt_moderation": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      api.IovInterruptModerationValue_name[api.IovInterruptModerationValue_Off],
							ValidateFunc: stringKeyInMap(api.IovInterruptModerationValue_value, true),
						},
						"iov_weight": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      100,
							ValidateFunc: validation.IntBetween(0, 100),
						},
						"ipsec_offload_maximum_security_association": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  512,
						},
						"maximum_bandwidth": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
						"minimum_bandwidth_absolute": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
						"minimum_bandwidth_weight": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      0,
							ValidateFunc: validation.IntBetween(0, 100),
						},
						"mandatory_feature_id": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
						"resource_pool_name": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						"test_replica_pool_name": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						"test_replica_switch_name": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						"virtual_subnet_id": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      0,
							ValidateFunc: ValueOrIntBetween(0, 4096, 16777215),
						},
						"allow_teaming": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      api.OnOffState_name[api.OnOffState_On],
							ValidateFunc: stringKeyInMap(api.OnOffState_value, true),
						},
						"not_monitored_in_cluster": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"storm_limit": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
						"dynamic_ip_address_limit": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
						"device_naming": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      api.OnOffState_name[api.OnOffState_Off],
							ValidateFunc: stringKeyInMap(api.OnOffState_value, true),
						},
						"fix_speed_10g": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      api.OnOffState_name[api.OnOffState_Off],
							ValidateFunc: stringKeyInMap(api.OnOffState_value, true),
						},
						"packet_direct_num_procs": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
						"packet_direct_moderation_count": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
						"packet_direct_moderation_interval": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
						"vrss_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"vmmq_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"vmmq_queue_pairs": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  16,
						},
						"vlan_access": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"vlan_id": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
						"wait_for_ips": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"ip_addresses": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The current list of IP addresses on this virtual machine.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},

			"dvd_drives": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"controller_number": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"controller_location": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"path": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						"resource_pool_name": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
					},
				},
			},

			"hard_disk_drives": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"controller_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      api.ControllerType_name[api.ControllerType_Scsi],
							ValidateFunc: stringKeyInMap(api.ControllerType_value, true),
						},
						"controller_number": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"controller_location": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"path": {
							Type:             schema.TypeString,
							Optional:         true,
							Default:          "",
							DiffSuppressFunc: api.DiffSuppressVmHardDiskPath,
						},
						"disk_number": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  MaxUint32,
						},
						"resource_pool_name": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "Primordial",
						},
						"support_persistent_reservations": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"maximum_iops": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
						"minimum_iops": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
						"qos_policy_id": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "00000000-0000-0000-0000-000000000000",
						},
						"override_cache_attributes": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      api.CacheAttributes_name[api.CacheAttributes_Default],
							ValidateFunc: stringKeyInMap(api.CacheAttributes_value, true),
						},
					},
				},
			},
		},
	}
}
