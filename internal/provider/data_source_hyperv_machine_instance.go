package provider

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/taliesins/terraform-provider-hyperv/api"
)

func dataSourceHyperVMachineInstance() *schema.Resource {
	return &schema.Resource{
		Description: "This Hyper-V data source provides information about existing virtual machine instances.",
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(ReadMachineInstanceTimeout),
		},
		ReadContext: datasourceHyperVMachineInstanceRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the virtual machine.",
			},

			"path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					if newValue == "" {
						return true
					}

					// When specifying path on new-vm it will auto append machine name on the end
					name := d.Get("name").(string)
					computedPath := newValue
					if !strings.HasSuffix(computedPath, "\\") {
						computedPath += "\\"
					}
					computedPath += name

					if strings.EqualFold(computedPath, oldValue) {
						return true
					}

					if strings.EqualFold(oldValue, newValue) {
						return true
					}

					return false
				},
				Description: "The path of the virtual machine.",
			},

			"generation": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          2,
				ValidateDiagFunc: IntInSlice([]int{1, 2}),
				ForceNew:         true,
				Description:      "Specifies the generation, as an integer, for the virtual machine. Valid values to use are `1`, `2`.",
			},

			"automatic_critical_error_action": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.CriticalErrorAction_name[api.CriticalErrorAction_Pause],
				ValidateDiagFunc: StringKeyInMap(api.CriticalErrorAction_value, true),
				Description:      "Specifies the action to take when the VM encounters a critical error, and exceeds the timeout duration specified by the AutomaticCriticalErrorActionTimeout cmdlet. Valid values to use are `Pause`, `None`.",
			},

			"automatic_critical_error_action_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     30,
				Description: "Specifies the amount of time, in minutes, to wait in critical pause before powering off the virtual machine.",
			},

			"automatic_start_action": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.StartAction_name[api.StartAction_StartIfRunning],
				ValidateDiagFunc: StringKeyInMap(api.StartAction_value, true),
				Description:      "Specifies the action the virtual machine is to take upon start. Valid values to use are `Nothing`, `StartIfRunning`, `Start`.",
			},

			"automatic_start_delay": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Specifies the number of seconds by which the virtual machine's start should be delayed.",
			},

			"automatic_stop_action": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.StopAction_name[api.StopAction_Save],
				ValidateDiagFunc: StringKeyInMap(api.StopAction_value, true),
				Description:      "Specifies the action the virtual machine is to take when the virtual machine host shuts down. Valid values to use are `TurnOff`, `Save`, `ShutDown`.",
			},

			"checkpoint_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.CheckpointType_name[api.CheckpointType_Production],
				ValidateDiagFunc: StringKeyInMap(api.CheckpointType_value, true),
				Description:      "Allows you to configure the type of checkpoints created by Hyper-V. If `Disabled` is specified, block creation of checkpoints. If `Standard` is specified, create standard checkpoints. If `Production` is specified, create production checkpoints if supported by guest operating system. Otherwise, create standard checkpoints. If `ProductionOnly` is specified, create production checkpoints if supported by guest operating system. Otherwise, the operation fails. Valid values to use are `Disabled`, `Standard`, `Production`, `ProductionOnly`.",
			},

			"dynamic_memory": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Specifies if machine instance will have dynamic memory enabled.",
			},

			"guest_controlled_cache_types": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Specifies if the machine instance will use guest controlled cache types.",
			},

			"high_memory_mapped_io_space": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     536870912,
				Description: "",
			},

			"lock_on_disconnect": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.OnOffState_name[api.OnOffState_Off],
				ValidateDiagFunc: StringKeyInMap(api.OnOffState_value, true),
				Description:      "Specifies whether virtual machine connection in basic mode locks the console after a user disconnects. Valid values to use are `On`, `Off`.",
			},

			"low_memory_mapped_io_space": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     134217728,
				Description: "",
			},

			"memory_maximum_bytes": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1099511627776,
				Description: "Specifies the maximum amount of memory that the virtual machine is to be allocated. (Applies only to virtual machines using dynamic memory.)",
			},

			"memory_minimum_bytes": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     536870912,
				Description: "Specifies the minimum amount of memory that the virtual machine is to be allocated. (Applies only to virtual machines using dynamic memory.)",
			},

			"memory_startup_bytes": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     536870912,
				Description: "Specifies the amount of memory that the virtual machine is to be allocated upon startup. (If the virtual machine does not use dynamic memory, then this is the static amount of memory to be allocated.)",
			},

			"notes": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Specifies a note to be associated with the machine to be created.",
			},

			"processor_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Specifies the number of virtual processors for the virtual machine.",
			},

			"smart_paging_file_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     `C:\ProgramData\Microsoft\Windows\Hyper-V`,
				Description: "Specifies the folder in which the Smart Paging file is to be stored.",
			},

			"snapshot_file_location": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     `C:\ProgramData\Microsoft\Windows\Hyper-V`,
				Description: "Specifies the folder in which the virtual machine is to store its snapshot files.",
			},

			"static_memory": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Specifies if the machine instance will use static memory.",
			},

			"integration_services": {
				Type:             schema.TypeMap,
				Optional:         true,
				DefaultFunc:      api.DefaultVmIntegrationServices,
				DiffSuppressFunc: api.DiffSuppressVmIntegrationServices,
				Elem:             schema.TypeBool,
				Description:      "A map of all the integration services and if the integration service should be enabled/disabled. Integration services that are not specified will not be enforced.",
			},

			"state": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.VmState_name[api.VmState_Running],
				ValidateDiagFunc: StringKeyInMap(api.VmState_SettableValue, true),
				Description:      "Specifies if the machine instance will be running or off. Valid values to use are `Running`, `Off`.",
			},

			"wait_for_state_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     120,
				Description: "The amount of time in seconds to wait before throwing an exception when trying to change for the virtual machine to the desired state.",
			},

			"wait_for_state_poll_period": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     2,
				Description: "The amount of time in seconds to wait between trying to change for the virtual machine to the desired state.",
			},

			"wait_for_ips_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     300,
				Description: "The amount of time in seconds to wait before throwing an exception when trying to get ip addresses for network cards on the virtual machine.",
			},

			"wait_for_ips_poll_period": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     5,
				Description: "The amount of time in seconds to wait between trying to get ip addresses for network cards on the virtual machine.",
			},

			"vm_processor": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				DefaultFunc: api.DefaultVmProcessors,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"compatibility_for_migration_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Specifies whether the virtual processor's features are to be limited for compatibility when migrating the virtual machine to another host.",
						},

						"compatibility_for_older_operating_systems_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Specifies whether the virtual processor's features are to be limited for compatibility with older operating systems.",
						},

						"hw_thread_count_per_core": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     0,
							Description: "Specifies the number of virtual SMT threads exposed to the virtual machine. Setting this value to 0 indicates the virtual machine will inherit the host's number of threads per core. This setting may not exceed the host's number of threads per core. Note: Windows Server 2016 does not support setting HwThreadCountPerCore to 0. For more details, see Configuring VM SMT settings using PowerShell.",
						},

						"maximum": {
							Type:             schema.TypeInt,
							Optional:         true,
							Default:          100,
							ValidateDiagFunc: ValueOrIntBetween(100, 0, 100),
							Description:      "Specifies the maximum percentage of resources available to the virtual machine processor to be configured. Allowed values range from 0 to 100.",
						},

						"reserve": {
							Type:             schema.TypeInt,
							Optional:         true,
							Default:          0,
							ValidateDiagFunc: ValueOrIntBetween(0, 0, 100),
							Description:      "Specifies the percentage of processor resources to be reserved for this virtual machine. Allowed values range from 0 to 100.",
						},

						"relative_weight": {
							Type:             schema.TypeInt,
							Optional:         true,
							Default:          100,
							ValidateDiagFunc: ValueOrIntBetween(100, 0, 10000),
							Description:      "Specifies the priority for allocating the physical computer's processing power to this virtual machine relative to others. Allowed values range from 1 to 10000.",
						},

						"maximum_count_per_numa_node": {
							Type:             schema.TypeInt,
							Optional:         true,
							Default:          0, // Dynamic value
							DiffSuppressFunc: api.DiffSuppressVmProcessorMaximumCountPerNumaNode,
							Description:      "Specifies the maximum number of processors per NUMA node to be configured for the virtual machine.",
						},

						"maximum_count_per_numa_socket": {
							Type:             schema.TypeInt,
							Optional:         true,
							Default:          0, // Dynamic value
							DiffSuppressFunc: api.DiffSuppressVmProcessorMaximumCountPerNumaSocket,
							Description:      "Specifies the maximum number of sockets per NUMA node to be configured for the virtual machine.",
						},

						"enable_host_resource_protection": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Specifies whether to enable host resource protection on the virtual machine. When enabled, the host will enforce limits on some aspects of the virtual machine's activity, preventing excessive consumption of host compute resources. VM activities controlled by this setting include the VMbus pipe messages associated with a subset of the VM's virtual devices, and intercepts generated by the VM. The virtual devices affected include the video, keyboard, mouse, and dynamic memory VDEVs.",
						},

						"expose_virtualization_extensions": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Specifies whether the hypervisor should expose the presence of virtualization extensions to the virtual machine, which enables support for nested virtualization.",
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
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies the name for the virtual network adapter.",
						},
						"switch_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							ForceNew:    false,
							Description: "Specifies the name of the virtual switch to connect to the new network adapter. If the switch name is not unique, then the operation fails.",
						},
						"management_os": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Specifies the virtual network adapter in the management operating system to be configured.",
						},
						"is_legacy": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							ForceNew:    true,
							Description: "Specifies whether the virtual network adapter is the legacy type.",
						},
						"dynamic_mac_address": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Assigns a dynamically generated MAC address to the virtual network adapter.",
						},
						"static_mac_address": {
							Type:             schema.TypeString,
							Optional:         true,
							Default:          "",
							DiffSuppressFunc: api.DiffSuppressVmStaticMacAddress,
							Description:      "Assigns a specific a MAC addresss to the virtual network adapter.",
						},
						"mac_address_spoofing": {
							Type:             schema.TypeString,
							Optional:         true,
							Default:          api.OnOffState_name[api.OnOffState_Off],
							ValidateDiagFunc: StringKeyInMap(api.OnOffState_value, true),
							Description:      "Specifies whether virtual machines may change the source MAC address in outgoing packets to one not assigned to them. On allows the virtual machine to use a different MAC address. Off only allows the virtual machine to use the MAC address assigned to it. Valid values to use are `On`, `Off`.",
						},
						"dhcp_guard": {
							Type:             schema.TypeString,
							Optional:         true,
							Default:          api.OnOffState_name[api.OnOffState_Off],
							ValidateDiagFunc: StringKeyInMap(api.OnOffState_value, true),
							Description:      "Specifies whether to drop DHCP messages from a virtual machine claiming to be a DHCP server. Valid values to use are `On`, `Off`.",
						},
						"router_guard": {
							Type:             schema.TypeString,
							Optional:         true,
							Default:          api.OnOffState_name[api.OnOffState_Off],
							ValidateDiagFunc: StringKeyInMap(api.OnOffState_value, true),
							Description:      "Specifies whether to drop Router Advertisement and Redirection messages from unauthorized virtual machines. If On is specified, such messages are dropped. If Off is specified, such messages are sent. Valid values to use are `On`, `Off`.",
						},
						"port_mirroring": {
							Type:             schema.TypeString,
							Optional:         true,
							Default:          api.PortMirroring_name[api.PortMirroring_None],
							ValidateDiagFunc: StringKeyInMap(api.PortMirroring_value, true),
							Description:      "Specifies the port mirroring mode for the network adapter to be configured. If a virtual network adapter is configured as Source, every packet it sends or receives is copied and forwarded to a virtual network adapter configured to receive the packets. If a virtual network adapter is configured as Destination, it receives copied packets from the source virtual network adapter. The source and destination virtual network adapters must be connected to the same virtual switch. Specify None to disable the feature. Valid values to use are `None`, `Source`, `Destination`.",
						},
						"ieee_priority_tag": {
							Type:             schema.TypeString,
							Optional:         true,
							Default:          api.OnOffState_name[api.OnOffState_Off],
							ValidateDiagFunc: StringKeyInMap(api.OnOffState_value, true),
							Description:      "Specifies whether IEEE 802.1p tagged packets from the virtual machine should be trusted. If it is on, the IEEE 802.1p tagged packets will be let go as is. If it is off, the priority value is reset to 0. Valid values to use are `On`, `Off`.",
						},
						"vmq_weight": {
							Type:             schema.TypeInt,
							Optional:         true,
							Default:          100,
							ValidateDiagFunc: IntBetween(0, 100),
							Description:      "Specifies whether virtual machine queue (VMQ) is to be enabled on the virtual network adapter. The relative weight describes the affinity of the virtual network adapter to use VMQ. Specify 0 to disable VMQ on the virtual network adapter. Valid values to use are between `1` to `100`.",
						},
						"iov_queue_pairs_requested": {
							Type:             schema.TypeInt,
							Optional:         true,
							Default:          1,
							ValidateDiagFunc: IntBetween(1, 4294967295),
							Description:      "Specifies the number of hardware queue pairs to be allocated to an SR-IOV virtual function. If receive-side scaling (RSS) is required, and if the physical network adapter that binds to the virtual switch supports RSS on SR-IOV virtual functions, then more than one queue pair is required. Valid values to use are between `1` to `4294967295`.",
						},
						"iov_interrupt_moderation": {
							Type:             schema.TypeString,
							Optional:         true,
							Default:          api.IovInterruptModerationValue_name[api.IovInterruptModerationValue_Off],
							ValidateDiagFunc: StringKeyInMap(api.IovInterruptModerationValue_value, true),
							Description:      "Specifies the interrupt moderation value for a single-root I/O virtualization (SR-IOV) virtual function assigned to a virtual network adapter. If Default is chosen, the value is determined by the physical network adapter vendor's setting. If Adaptive is chosen, the interrupt moderation rate will be based on the runtime traffic pattern. Valid values to use are `Default`, `Adaptive`, `Off`, `Low `, `Medium`, `High`.",
						},
						"iov_weight": {
							Type:             schema.TypeInt,
							Optional:         true,
							Default:          100,
							ValidateDiagFunc: IntBetween(0, 100),
							Description:      "Specifies whether single-root I/O virtualization (SR-IOV) is to be enabled on this virtual network adapter. The relative weight sets the affinity of the virtual network adapter to the assigned SR-IOV virtual function. Specify 0 to disable SR-IOV on the virtual network adapter. Valid values to use are between `0` to `100`.",
						},
						"ipsec_offload_maximum_security_association": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     512,
							Description: "Specifies the maximum number of security associations that can be offloaded to the physical network adapter that is bound to the virtual switch and that supports IPSec Task Offload. Specify zero to disable the feature.",
						},
						"maximum_bandwidth": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     0,
							Description: "Specifies the maximum bandwidth, in bits per second, for the virtual network adapter. The specified value is rounded to the nearest multiple of eight. Specify zero to disable the feature.",
						},
						"minimum_bandwidth_absolute": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     0,
							Description: "Specifies the minimum bandwidth, in bits per second, for the virtual network adapter. The specified value is rounded to the nearest multiple of eight. A value larger than 100 Mbps is recommended.",
						},
						"minimum_bandwidth_weight": {
							Type:             schema.TypeInt,
							Optional:         true,
							Default:          0,
							ValidateDiagFunc: IntBetween(0, 100),
							Description:      "Specifies the minimum bandwidth, in terms of relative weight, for the virtual network adapter. The weight describes how much bandwidth to provide to the virtual network adapter relative to other virtual network adapters connected to the same virtual switch. Specify 0 to disable the feature. Valid values to use are between `0` to `100`.",
						},
						"mandatory_feature_id": {
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         schema.HashString,
							Description: "Specifies the unique identifiers of the virtual switch extension features that are required for this virtual network adapter to operate.",
						},
						"resource_pool_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "Specifies the name of the resource pool.",
						},
						"test_replica_pool_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "This parameter applies only to virtual machines that are enabled for replication. It specifies the name of the network resource pool that will be used by this virtual network adapter when its virtual machine is created during a test failover.",
						},
						"test_replica_switch_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "This parameter applies only to virtual machines that are enabled for replication. It specifies the name of the virtual switch to which the virtual network adapter should be connected when its virtual machine is created during a test failover.",
						},
						"virtual_subnet_id": {
							Type:             schema.TypeInt,
							Optional:         true,
							Default:          0,
							ValidateDiagFunc: ValueOrIntBetween(0, 4096, 16777215),
							Description:      "Specifies the virtual subnet ID to use with Hyper-V Network Virtualization. Use 0 to clear this parameter. Valid values to use are `0` or between `4096` to `16777215` (2^24 - 1).",
						},
						"allow_teaming": {
							Type:             schema.TypeString,
							Optional:         true,
							Default:          api.OnOffState_name[api.OnOffState_On],
							ValidateDiagFunc: StringKeyInMap(api.OnOffState_value, true),
							Description:      "Specifies whether the virtual network adapter can be teamed with other network adapters connected to the same virtual switch. Valid values to use are `On`, `Off`.",
						},
						"not_monitored_in_cluster": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Indicates whether to not monitor the network adapter if the virtual machine that it belongs to is part of a cluster. By default, network adapters for clustered virtual machines are monitored.",
						},
						"storm_limit": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     0,
							Description: "Specifies the number of broadcast, multicast, and unknown unicast packets per second a virtual machine is allowed to send through the specified virtual network adapter. Broadcast, multicast, and unknown unicast packets beyond the limit during that one second interval are dropped. A value of zero (0) means there is no limit.",
						},
						"dynamic_ip_address_limit": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     0,
							Description: "Specifies the dynamic IP address limit.",
						},
						"device_naming": {
							Type:             schema.TypeString,
							Optional:         true,
							Default:          api.OnOffState_name[api.OnOffState_Off],
							ValidateDiagFunc: StringKeyInMap(api.OnOffState_value, true),
							Description:      "Specifies whether this adapter uses device naming. Valid values to use are `On`, `Off`.",
						},
						"fix_speed_10g": {
							Type:             schema.TypeString,
							Optional:         true,
							Default:          api.OnOffState_name[api.OnOffState_Off],
							ValidateDiagFunc: StringKeyInMap(api.OnOffState_value, true),
							Description:      "Specifies whether the adapter uses fix speed of 10G. Valid values to use are `On`, `Off`.",
						},
						"packet_direct_num_procs": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     0,
							Description: "Specifies the number of processors to use for virtual switch processing inside of the host.",
						},
						"packet_direct_moderation_count": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     0,
							Description: "Specifies the number of packets to wait for before signaling an interrupt.",
						},
						"packet_direct_moderation_interval": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     0,
							Description: "Specifies the amount of time, in milliseconds, to wait before signaling an interrupt after a packet arrives.",
						},
						"vrss_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Should Virtual Receive Side Scaling be enabled. This configuration allows the load from a virtual network adapter to be distributed across multiple virtual processors in a virtual machine (VM), allowing the VM to process more network traffic more rapidly than it can with a single logical processor.",
						},
						"vmmq_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Should Virtual Machine Multi-Queue be enabled. With set to true multiple queues are allocated to a single VM with each queue affinitized to a core in the VM.",
						},
						"vmmq_queue_pairs": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     16,
							Description: "The number of Virtual Machine Multi-Queues to create for this VM.",
						},
						"vlan_access": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "",
						},
						"vlan_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     0,
							Description: "",
						},
						"wait_for_ips": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Wait for the network card to be assigned an ip address. ",
						},
						"ip_addresses": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The current list of IP addresses on this machine. If HyperV integration tools is not running on the virtual machine, or if the VM is powered off, or has not been assigned an ip address, this list will be empty. ",
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
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Specifies the number of the controller to which the DVD drive is to be added.",
						},
						"controller_location": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Specifies the number of the location on the controller at which the DVD drive is to be added. ",
						},
						"path": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "Specifies the full path to the virtual hard disk file or physical hard disk volume for the added DVD drive.",
						},
						"resource_pool_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "Specifies the friendly name of the ISO resource pool to which this DVD drive is to be associated.",
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
							Type:             schema.TypeString,
							Optional:         true,
							Default:          api.ControllerType_name[api.ControllerType_Scsi],
							ValidateDiagFunc: StringKeyInMap(api.ControllerType_value, true),
							Description:      "Specifies the type of the controller to which the hard disk drive is to be added. Valid values to use are `Ide`, `Scsi`.",
						},
						"controller_number": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Specifies the number of the controller to which the hard disk drive is to be added.",
						},
						"controller_location": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Specifies the number of the location on the controller at which the hard disk drive is to be added.",
						},
						"path": {
							Type:             schema.TypeString,
							Optional:         true,
							Default:          "",
							DiffSuppressFunc: api.DiffSuppressVmHardDiskPath,
							Description:      "Specifies the full path of the hard disk drive file to be added.",
						},
						"disk_number": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     MaxUint32,
							Description: "Specifies the disk number of the offline physical hard drive to be connected as a passthrough disk. If value is 4294967295 then disk number is ignored.",
						},
						"resource_pool_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "Primordial",
							Description: "Specifies the friendly name of the resource pool to which this virtual hard disk is to be associated.",
						},
						"support_persistent_reservations": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Indicates that the hard disk supports SCSI persistent reservation semantics. Specify this parameter when the hard disk is a shared disk that is used by multiple virtual machines.",
						},
						"maximum_iops": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     0,
							Description: "Specifies the maximum normalized I/O operations per second (IOPS) for the hard disk. Hyper-V calculates normalized IOPS as the total size of I/O per second divided by 8 KB. If value is 0 then iops is ignored.",
						},
						"minimum_iops": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     0,
							Description: "Specifies the minimum normalized I/O operations per second (IOPS) for the hard disk. Hyper-V calculates normalized IOPS as the total size of I/O per second divided by 8 KB. If maximum iops value is 0 then iops is ignored.",
						},
						"qos_policy_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "00000000-0000-0000-0000-000000000000",
							Description: "Specifies the unique ID for a storage QoS policy that this cmdlet associates with the hard disk drive. If value is 00000000-0000-0000-0000-000000000000 then qos policy id is ignored.",
						},
						"override_cache_attributes": {
							Type:             schema.TypeString,
							Optional:         true,
							Default:          api.CacheAttributes_name[api.CacheAttributes_Default],
							ValidateDiagFunc: StringKeyInMap(api.CacheAttributes_value, true),
							Description:      "With Default it is equivalent of WriteCacheDisabled. With WriteCacheEnabled write I/O is acknowledged as written before it is committed to stable media. If your internal disks, DAS, SAN, or NAS has a battery backup system that can guarantee clean cache flushes on a power outage, write caching is generally safe. Internal batteries that report their status and/or automatically disable caching are best. UPS-backed systems are sometimes OK, but they are not foolproof. With WriteCacheAndFUAEnabled write I/O is committed to stable media BEFORE the I/O is acknowledged as written. With WriteCacheDisabled when I/O is written it is acknowledged as written as there is no cache in between. Valid values to use are `Default`, `WriteCacheEnabled`, `WriteCacheAndFUAEnabled`, `WriteCacheDisabled`.",
						},
					},
				},
			},

			"vm_firmware": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				//DefaultFunc: api.DefaultVmFirmwares,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"boot_order": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"boot_type": {
										Type:             schema.TypeString,
										Required:         true,
										ValidateDiagFunc: StringKeyInMap(api.Gen2BootType_value, true),
										Description:      "The type of boot device. Valid values to use are `NetworkAdapter`, `HardDiskDrive` and `DvdDrive`.",
									},
									"network_adapter_name": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "",
										DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
											bootTypeKey := strings.Replace(k, "network_adapter_name", "boot_type", 1)
											bootType := d.Get(bootTypeKey).(string)

											if bootType == "" {
												return true
											}

											if newValue == "" || oldValue == newValue {
												return true
											}
											return false
										},
										Description: "Specifies the name of ethernet adapter.",
									},
									"switch_name": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "",
										DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
											bootTypeKey := strings.Replace(k, "switch_name", "boot_type", 1)
											bootType := d.Get(bootTypeKey).(string)
											if bootType == "" {
												return true
											}

											if newValue == "" || oldValue == newValue {
												return true
											}
											return false
										},
										Description: "Specifies the name of ethernet adapter switch.",
									},
									"mac_address": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "",
										DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
											bootTypeKey := strings.Replace(k, "mac_address", "boot_type", 1)
											bootType := d.Get(bootTypeKey).(string)
											if bootType == "" {
												return true
											}

											if newValue == "" || strings.EqualFold(oldValue, newValue) {
												return true
											}
											return false
										},
										Description: "Specifies the mac address of ethernet adapter.",
									},
									"path": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "",
										DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
											bootTypeKey := strings.Replace(k, "path", "boot_type", 1)
											bootType := d.Get(bootTypeKey).(string)
											if bootType == "" {
												return true
											}

											if newValue == "" {
												return true
											}

											// When specifying path on new-vm it will auto append machine name on the end
											name := d.Get("name").(string)
											computedPath := newValue
											if !strings.HasSuffix(computedPath, "\\") {
												computedPath += "\\"
											}
											computedPath += name

											if strings.EqualFold(computedPath, oldValue) {
												return true
											}

											if strings.EqualFold(oldValue, newValue) {
												return true
											}

											return false
										},
										Description: "Specifies the file path of hard disk drive or dvd drive.",
									},
									"controller_number": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  -1,
										DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
											bootTypeKey := strings.Replace(k, "controller_number", "boot_type", 1)
											bootType := d.Get(bootTypeKey).(string)
											if bootType == "" {
												return true
											}

											return newValue == "-1"
										},
										Description: "Specifies the number of the controller to which the hard disk drive or dvd drive.",
									},
									"controller_location": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  -1,
										DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
											bootTypeKey := strings.Replace(k, "controller_location", "boot_type", 1)
											bootType := d.Get(bootTypeKey).(string)
											if bootType == "" {
												return true
											}

											return newValue == "-1"
										},
										Description: "Specifies the number of the location on the controller at which the hard disk drive or dvd drive.",
									},
								},
							},
							Description: "The boot order of the devices that the generation 2 virtual machine should try to use for boot up.",
						},

						"enable_secure_boot": {
							Type:             schema.TypeString,
							Optional:         true,
							Default:          api.OnOffState_name[api.OnOffState_On],
							ValidateDiagFunc: StringKeyInMap(api.OnOffState_value, true),
							Description:      "Specifies whether to enable secure boot. Valid values to use are `On`, `Off`.",
						},

						"secure_boot_template": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "MicrosoftWindows",
							Description: "Specifies the name of the secure boot template. If secure boot is enabled, you must have a valid secure boot template for the guest operating system to start. Example values to use are `MicrosoftWindows`,`MicrosoftUEFICertificateAuthority`, `OpenSourceShieldedVM`.",
						},

						"preferred_network_boot_protocol": {
							Type:             schema.TypeString,
							Optional:         true,
							Default:          api.IPProtocolPreference_name[api.IPProtocolPreference_IPv4],
							ValidateDiagFunc: StringKeyInMap(api.IPProtocolPreference_value, true),
							Description:      "Specifies the IP protocol version to use during a network boot. Valid values to use are `IPv4`, `IPv6`.",
						},

						"console_mode": {
							Type:             schema.TypeString,
							Optional:         true,
							Default:          api.ConsoleModeType_name[api.ConsoleModeType_Default],
							ValidateDiagFunc: StringKeyInMap(api.ConsoleModeType_value, true),
							Description:      "Specifies the console mode type for the virtual machine. This parameter allows a virtual machine to run without graphical user interface. Valid values to use are `Default`, `COM1`, `COM2`, `None`.",
						},

						"pause_after_boot_failure": {
							Type:             schema.TypeString,
							Optional:         true,
							Default:          api.OnOffState_name[api.OnOffState_Off],
							ValidateDiagFunc: StringKeyInMap(api.OnOffState_value, true),
							Description:      "Specifies the behavior of the virtual machine after a start failure. For a value of On, if the virtual machine fails to start correctly from a device, the virtual machine is paused. Valid values to use are `On`, `Off`.",
						},
					},
				},
				Description: "",
			},
		},
	}
}

func datasourceHyperVMachineInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO][hyperv][read] reading hyperv machine: %#v", d)
	client := meta.(api.Client)

	var name string
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		return diag.Errorf("[ERROR][hyperv][read] name argument is required")
	}

	vm, err := client.GetVm(ctx, name)
	if err != nil {
		return diag.FromErr(err)
	}

	var vmFirmwares []api.VmFirmware
	if vm.Generation > 1 {
		vmFirmwares, err = client.GetVmFirmwares(ctx, name)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		vmFirmwares = client.GetNoVmFirmwares(ctx)
	}

	vmProcessors, err := client.GetVmProcessors(ctx, name)
	if err != nil {
		return diag.FromErr(err)
	}

	integrationServices, err := client.GetVmIntegrationServices(ctx, name)
	if err != nil {
		return diag.FromErr(err)
	}

	dvdDrives, err := client.GetVmDvdDrives(ctx, name)
	if err != nil {
		return diag.FromErr(err)
	}

	hardDiskDrives, err := client.GetVmHardDiskDrives(ctx, name)
	if err != nil {
		return diag.FromErr(err)
	}

	vmState, err := client.GetVmStatus(ctx, name)
	if err != nil {
		return diag.FromErr(err)
	}

	networkAdaptersWaitForIps, waitForIpsTimeout, waitForIpsPollPeriod, err := api.ExpandVmNetworkAdapterWaitForIps(d)
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.WaitForVmNetworkAdaptersIps(ctx, name, waitForIpsTimeout, waitForIpsPollPeriod, networkAdaptersWaitForIps)
	if err != nil {
		return diag.FromErr(err)
	}

	networkAdapters, err := client.GetVmNetworkAdapters(ctx, name, networkAdaptersWaitForIps)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO][hyperv][read] retrieved vm: %+v", vm)

	if vm.Name != name {
		log.Printf("[INFO][hyperv][read] unable to read hyperv machine as it does not exist: %#v", name)
		return nil
	}

	if vm.DynamicMemory && vm.StaticMemory {
		return diag.Errorf("[ERROR][hyperv][read] Dynamic and static can't be both selected at the same time")
	}

	if !vm.DynamicMemory && !vm.StaticMemory {
		return diag.Errorf("[ERROR][hyperv][read] Either dynamic or static must be selected")
	}

	flattenedVmFirmwares := api.FlattenVmFirmwares(&vmFirmwares)
	if err := d.Set("vm_firmware", flattenedVmFirmwares); err != nil {
		return diag.Errorf("[DEBUG] Error setting vm_firmware error: %v", err)
	}
	if vm.Generation > 1 {
		log.Printf("[INFO][hyperv][read] vmFirmwares: %v", vmFirmwares)
		log.Printf("[INFO][hyperv][read] flattenedVmFirmwares: %v", flattenedVmFirmwares)
	} else {
		log.Printf("[INFO][hyperv][read] skip vmFirmwares as vm generation is %v", vm.Generation)
		log.Printf("[INFO][hyperv][read] skip flattenedVmFirmwares as vm generation is %v", vm.Generation)
	}

	flattenedVmProcessors := api.FlattenVmProcessors(&vmProcessors)
	if err := d.Set("vm_processor", flattenedVmProcessors); err != nil {
		return diag.Errorf("[DEBUG] Error setting vm_processor error: %v", err)
	}
	log.Printf("[INFO][hyperv][read] vmProcessors: %v", vmProcessors)
	log.Printf("[INFO][hyperv][read] flattenedVmProcessors: %v", flattenedVmProcessors)

	flattenedIntegrationServices := api.FlattenIntegrationServices(&integrationServices)
	if err := d.Set("integration_services", flattenedIntegrationServices); err != nil {
		return diag.Errorf("[DEBUG] Error setting integration_services error: %v", err)
	}

	flattenedDvdDrives := api.FlattenDvdDrives(&dvdDrives)
	if err := d.Set("dvd_drives", flattenedDvdDrives); err != nil {
		return diag.Errorf("[DEBUG] Error setting dvd_drives error: %v", err)
	}
	log.Printf("[INFO][hyperv][read] dvdDrives: %v", dvdDrives)
	log.Printf("[INFO][hyperv][read] flattenedDvdDrives: %v", flattenedDvdDrives)

	flattenedHardDiskDrives := api.FlattenHardDiskDrives(&hardDiskDrives)
	if err := d.Set("hard_disk_drives", flattenedHardDiskDrives); err != nil {
		return diag.Errorf("[DEBUG] Error setting hard_disk_drives error: %v", err)
	}
	log.Printf("[INFO][hyperv][read] hardDiskDrives: %v", hardDiskDrives)
	log.Printf("[INFO][hyperv][read] flattenedHardDiskDrives: %v", flattenedHardDiskDrives)

	flattenedNetworkAdapters := api.FlattenNetworkAdapters(&networkAdapters)
	if err := d.Set("network_adaptors", flattenedNetworkAdapters); err != nil {
		return diag.Errorf("[DEBUG] Error setting network_adaptors error: %v", err)
	}
	log.Printf("[INFO][hyperv][read] networkAdapters: %v", networkAdapters)
	log.Printf("[INFO][hyperv][read] flattenedNetworkAdapters: %v", flattenedNetworkAdapters)

	if err := d.Set("name", vm.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("path", vm.Path); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("generation", vm.Generation); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("automatic_critical_error_action", vm.AutomaticCriticalErrorAction.String()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("automatic_critical_error_action_timeout", vm.AutomaticCriticalErrorActionTimeout); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("automatic_start_action", vm.AutomaticStartAction.String()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("automatic_start_delay", vm.AutomaticStartDelay); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("automatic_stop_action", vm.AutomaticStopAction.String()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("checkpoint_type", vm.CheckpointType.String()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("dynamic_memory", vm.DynamicMemory); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("guest_controlled_cache_types", vm.GuestControlledCacheTypes); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("high_memory_mapped_io_space", vm.HighMemoryMappedIoSpace); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("lock_on_disconnect", vm.LockOnDisconnect.String()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("low_memory_mapped_io_space", vm.LowMemoryMappedIoSpace); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("memory_maximum_bytes", vm.MemoryMaximumBytes); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("memory_minimum_bytes", vm.MemoryMinimumBytes); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("memory_startup_bytes", vm.MemoryStartupBytes); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("notes", vm.Notes); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("processor_count", vm.ProcessorCount); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("smart_paging_file_path", vm.SmartPagingFilePath); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("snapshot_file_location", vm.SnapshotFileLocation); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("static_memory", vm.StaticMemory); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("state", vmState.State.String()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(name)

	log.Printf("[INFO][hyperv][read] read hyperv machine: %#v", d)

	return nil
}
