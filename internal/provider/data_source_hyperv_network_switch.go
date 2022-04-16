package provider

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/taliesins/terraform-provider-hyperv/api"
)

func dataSourceHyperVNetworkSwitch() *schema.Resource {
	return &schema.Resource{
		Description: "Get information about existing virtual network switches.",
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(ReadNetworkSwitchTimeout),
		},
		ReadContext: datasourceHyperVNetworkSwitchRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the name of the switch.",
			},

			"notes": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Specifies a note to be associated with the switch.",
			},

			"allow_management_os": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true, //This is tied to the switch type used. internal=true;private=false;external=true or false
				Description: "Specifies if the HyperV host machine will have access to network switch when created. It provides this access via a virtual adaptor, so you will need to either configure static ips on the virtual adaptor or configure a dhcp on a machine connected to the network switch. This is tied to the switch type used: `internal=true`;`private=false`;`external=true or false`.",
			},

			"enable_embedded_teaming": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Specifies if the HyperV host machine will enable teaming for network switch when created. It allows NIC teaming so that you could support scenarios such as redundant links. ",
			},

			"enable_iov": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Specifies if the HyperV host machine will enable IO virtualization for network switch when created. If your hardware supports it, it enables the virtual machine to talk directly to the NIC.",
			},

			"enable_packet_direct": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Specifies if the HyperV host machine will enable packet direct path for network switch when created. Increases packet throughoutput and reduces the network latency between vms on the switch.",
			},

			"minimum_bandwidth_mode": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.VMSwitchBandwidthMode_name[api.VMSwitchBandwidthMode_None],
				ValidateDiagFunc: stringKeyInMap(api.VMSwitchBandwidthMode_value, true),
				ForceNew:         true,
				Description:      "Valid values to use are `Absolute`, `Default`, `None`, `Weight`. Specifies how minimum bandwidth is to be configured on the virtual switch. If `Absolute` is specified, minimum bandwidth is bits per second. If `Weight` is specified, minimum bandwidth is a value ranging from `1` to `100`. If `None` is specified, minimum bandwidth is disabled on the switch – that is, users cannot configure it on any network adapter connected to the switch. If `Default` is specified, the system will set the mode to Weight, if the switch is not IOV-enabled, or `None` if the switch is IOV-enabled.",
			},

			"switch_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.VMSwitchType_name[api.VMSwitchType_Internal],
				ValidateDiagFunc: stringKeyInMap(api.VMSwitchType_value, true),
				Description:      "Valid values to use are `Internal`, `Private` and `External`. Specifies the type of the switch to be created. ",
			},

			"net_adapter_names": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: " Specifies the name of the network adapter to be bound to the switch. ",
			},

			"default_flow_minimum_bandwidth_absolute": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Specifies the minimum bandwidth, in bits per second, that is allocated to a special category called `default flow`. Any traffic sent by a virtual network adapter that is connected to this virtual switch and does not have minimum bandwidth allocated is filtered into this category. Specify a value for this parameter only if the minimum bandwidth mode on this virtual switch is absolute. By default, the virtual switch allocates 10% of the total bandwidth, which depends on the physical network adapter it binds to, to this category. For example, if a virtual switch binds to a 1 GbE network adapter, this special category can use at least 100 Mbps. If the value is not a multiple of 8, the value is rounded down to the nearest number that is a multiple of 8. For example, a value input as 1234567 is converted to 1234560.",
			},

			"default_flow_minimum_bandwidth_weight": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          0,
				ValidateDiagFunc: IntBetween(0, 100),
				Description:      "Should be a value of `0` or between `1` to `100`. Specifies the minimum bandwidth, in relative weight, that is allocated to a special category called `default flow`. Any traffic sent by a virtual network adapter that is connected to this virtual switch and doesn’t have minimum bandwidth allocated is filtered into this category. Specify a value for this parameter only if the minimum bandwidth mode on this virtual switch is weight. By default, this special category has a weight of 1.",
			},

			"default_queue_vmmq_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Should Virtual Machine Multi-Queue be enabled. With set to true multiple queues are allocated to a single VM with each queue affinitized to a core in the VM.",
			},

			"default_queue_vmmq_queue_pairs": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     16,
				Description: "The number of Virtual Machine Multi-Queues to create for this VM.",
			},

			"default_queue_vrss_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Should Virtual Receive Side Scaling be enabled. This configuration allows the load from a virtual network adapter to be distributed across multiple virtual processors in a virtual machine (VM), allowing the VM to process more network traffic more rapidly than it can with a single logical processor.",
			},
		},
	}
}

func datasourceHyperVNetworkSwitchRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO][hyperv][read] reading hyperv switch: %#v", d)
	c := meta.(api.Client)

	var switchName string

	if v, ok := d.GetOk("name"); ok {
		switchName = v.(string)
	} else {
		return diag.Errorf("[ERROR][hyperv][read] name argument is required")
	}

	s, err := c.GetVMSwitch(switchName)

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO][hyperv][read] retrieved network switch: %+v", s)

	if s.Name != switchName {
		log.Printf("[INFO][hyperv][read] unable to read hyperv switch as it does not exist: %#v", switchName)
		return nil
	}

	if s.SwitchType == api.VMSwitchType_Private {
		if s.AllowManagementOS {
			return diag.Errorf("[ERROR][hyperv][read] Unable to set AllowManagementOS to true if switch type is private")
		}

		if len(s.NetAdapterNames) > 0 {
			return diag.Errorf("[ERROR][hyperv][read] Unable to set NetAdapterNames when switch type is private")
		}
	} else if s.SwitchType == api.VMSwitchType_Internal {
		if !s.AllowManagementOS {
			return diag.Errorf("[ERROR][hyperv][read] Unable to set AllowManagementOS to false if switch type is internal")
		}

		if len(s.NetAdapterNames) > 0 {
			return diag.Errorf("[ERROR][hyperv][read] Unable to set NetAdapterNames when switch type is internal")
		}
	} else if s.SwitchType == api.VMSwitchType_External {
		if len(s.NetAdapterNames) < 1 {
			return diag.Errorf("[ERROR][hyperv][read] Must specify NetAdapterNames if switch type is external")
		}
	}

	if s.BandwidthReservationMode == api.VMSwitchBandwidthMode_Absolute {
		if s.DefaultFlowMinimumBandwidthWeight != 0 {
			return diag.Errorf("[ERROR][hyperv][read] DefaultFlowMinimumBandwidthWeight should be 0 if bandwidth reservation mode is absolute")
		}
		if s.DefaultFlowMinimumBandwidthAbsolute < 0 {
			return diag.Errorf("[ERROR][hyperv][read] Bandwidth absolute must be 0 or greater")
		}
	} else if s.BandwidthReservationMode == api.VMSwitchBandwidthMode_Weight || (s.BandwidthReservationMode == api.VMSwitchBandwidthMode_Default && (!s.IovEnabled)) {
		if s.DefaultFlowMinimumBandwidthAbsolute != 0 {
			return diag.Errorf("[ERROR][hyperv][read] DefaultFlowMinimumBandwidthAbsolute should be 0 if bandwidth reservation mode is weight")
		}
		if s.DefaultFlowMinimumBandwidthWeight < 1 || s.DefaultFlowMinimumBandwidthWeight > 100 {
			return diag.Errorf("[ERROR][hyperv][read] Bandwidth weight must be between 1 and 100")
		}
	} else {
		if s.DefaultFlowMinimumBandwidthWeight != 0 {
			return diag.Errorf("[ERROR][hyperv][read] DefaultFlowMinimumBandwidthWeight should be 0 if bandwidth reservation mode is none")
		}
		if s.DefaultFlowMinimumBandwidthAbsolute != 0 {
			return diag.Errorf("[ERROR][hyperv][read] DefaultFlowMinimumBandwidthAbsolute should be 0 if bandwidth reservation mode is none")
		}
	}

	if s.DefaultQueueVmmqQueuePairs < 1 {
		return diag.Errorf("[ERROR][hyperv][read] defaultQueueVmmqQueuePairs must be greater then 0")
	}

	if err := d.Set("name", s.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("notes", s.Notes); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("allow_management_os", s.AllowManagementOS); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("enable_embedded_teaming", s.EmbeddedTeamingEnabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("enable_iov", s.IovEnabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("enable_packet_direct", s.PacketDirectEnabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("minimum_bandwidth_mode", s.BandwidthReservationMode.String()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("switch_type", s.SwitchType.String()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("net_adapter_names", s.NetAdapterNames); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("default_flow_minimum_bandwidth_absolute", s.DefaultFlowMinimumBandwidthAbsolute); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("default_flow_minimum_bandwidth_weight", s.DefaultFlowMinimumBandwidthWeight); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("default_queue_vmmq_enabled", s.DefaultQueueVmmqEnabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("default_queue_vmmq_queue_pairs", s.DefaultQueueVmmqQueuePairs); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("default_queue_vrss_enabled", s.DefaultQueueVrssEnabled); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(switchName)

	log.Printf("[INFO][hyperv][read] read hyperv switch: %#v", d.State())

	return nil
}
