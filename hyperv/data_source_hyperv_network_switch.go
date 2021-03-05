package hyperv

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/taliesins/terraform-provider-hyperv/api"
)

func dataSourceHyperVNetworkSwitch() *schema.Resource {
	return &schema.Resource{
		Read:   resourceHyperVNetworkSwitchRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"notes": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},

			"allow_management_os": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true, //This is tied to the switch type used. internal=true;private=false;external=true or false
			},

			"enable_embedded_teaming": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},

			"enable_iov": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},

			"enable_packet_direct": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},

			"minimum_bandwidth_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      api.VMSwitchBandwidthMode_name[api.VMSwitchBandwidthMode_None],
				ValidateFunc: stringKeyInMap(api.VMSwitchBandwidthMode_value, true),
				ForceNew:     true,
			},

			"switch_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      api.VMSwitchType_name[api.VMSwitchType_Internal],
				ValidateFunc: stringKeyInMap(api.VMSwitchType_value, true),
			},

			"net_adapter_names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},

			"default_flow_minimum_bandwidth_absolute": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			"default_flow_minimum_bandwidth_weight": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntBetween(0, 100),
			},

			"default_queue_vmmq_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"default_queue_vmmq_queue_pairs": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  16,
			},

			"default_queue_vrss_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}
