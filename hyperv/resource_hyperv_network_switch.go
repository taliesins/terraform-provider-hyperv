package hyperv

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/taliesins/terraform-provider-hyperv/api"
)

func resourceHyperVNetworkSwitch() *schema.Resource {
	return &schema.Resource{
		Create: resourceHyperVNetworkSwitchCreate,
		Read:   resourceHyperVNetworkSwitchRead,
		Update: resourceHyperVNetworkSwitchUpdate,
		Delete: resourceHyperVNetworkSwitchDelete,

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

func resourceHyperVNetworkSwitchCreate(d *schema.ResourceData, meta interface{}) (err error) {

	log.Printf("[INFO][hyperv][create] creating hyperv switch: %#v", d)
	c := meta.(*api.HypervClient)

	switchName := ""

	if v, ok := d.GetOk("name"); ok {
		switchName = v.(string)
	} else {
		return fmt.Errorf("[ERROR][hyperv][create] name argument is required")
	}

	notes := (d.Get("notes")).(string)
	allowManagementOS := (d.Get("allow_management_os")).(bool)
	embeddedTeamingEnabled := (d.Get("enable_embedded_teaming")).(bool)
	iovEnabled := (d.Get("enable_iov")).(bool)
	packetDirectEnabled := (d.Get("enable_packet_direct")).(bool)
	bandwidthReservationMode := api.ToVMSwitchBandwidthMode((d.Get("minimum_bandwidth_mode")).(string))
	switchType := api.ToVMSwitchType((d.Get("switch_type")).(string))
	netAdapterNames := []string{}
	if raw, ok := d.GetOk("net_adapter_names"); ok {
		for _, v := range raw.([]interface{}) {
			netAdapterNames = append(netAdapterNames, v.(string))
		}
	}
	defaultFlowMinimumBandwidthAbsolute := int64((d.Get("default_flow_minimum_bandwidth_absolute")).(int))
	defaultFlowMinimumBandwidthWeight := int64((d.Get("default_flow_minimum_bandwidth_weight")).(int))
	defaultQueueVmmqEnabled := (d.Get("default_queue_vmmq_enabled")).(bool)
	defaultQueueVmmqQueuePairs := int32((d.Get("default_queue_vmmq_queue_pairs")).(int))
	defaultQueueVrssEnabled := (d.Get("default_queue_vrss_enabled")).(bool)

	if switchType == api.VMSwitchType_Private {
		if allowManagementOS {
			return fmt.Errorf("[ERROR][hyperv][create] Unable to set AllowManagementOS to true if switch type is private")
		}

		if len(netAdapterNames) > 0 {
			return fmt.Errorf("[ERROR][hyperv][create] Unable to set NetAdapterNames when switch type is private")
		}
	} else if switchType == api.VMSwitchType_Internal {
		if !allowManagementOS {
			return fmt.Errorf("[ERROR][hyperv][create] Unable to set AllowManagementOS to false if switch type is internal")
		}

		if len(netAdapterNames) > 0 {
			return fmt.Errorf("[ERROR][hyperv][create] Unable to set NetAdapterNames when switch type is internal")
		}
	} else if switchType == api.VMSwitchType_External {
		if len(netAdapterNames) < 1 {
			return fmt.Errorf("[ERROR][hyperv][create] Must specify NetAdapterNames if switch type is external")
		}
	}

	if bandwidthReservationMode == api.VMSwitchBandwidthMode_Absolute {
		if defaultFlowMinimumBandwidthWeight != 0 {
			return fmt.Errorf("[ERROR][hyperv][create] Unable to set DefaultFlowMinimumBandwidthWeight if bandwidth reservation mode is absolute")
		}
		if defaultFlowMinimumBandwidthAbsolute < 0 {
			return fmt.Errorf("[ERROR][hyperv][create] Bandwidth absolute must be 0 or greater")
		}
	} else if bandwidthReservationMode == api.VMSwitchBandwidthMode_Weight || (bandwidthReservationMode == api.VMSwitchBandwidthMode_Default && (!iovEnabled)) {
		if defaultFlowMinimumBandwidthAbsolute != 0 {
			return fmt.Errorf("[ERROR][hyperv][create] Unable to set DefaultFlowMinimumBandwidthAbsolute if bandwidth reservation mode is weight")
		}
		if defaultFlowMinimumBandwidthWeight < 1 || defaultFlowMinimumBandwidthWeight > 100 {
			return fmt.Errorf("[ERROR][hyperv][create] Bandwidth weight must be between 1 and 100")
		}
	} else {
		if defaultFlowMinimumBandwidthWeight != 0 {
			return fmt.Errorf("[ERROR][hyperv][create] Unable to set DefaultFlowMinimumBandwidthWeight if bandwidth reservation mode is none")
		}
		if defaultFlowMinimumBandwidthAbsolute != 0 {
			return fmt.Errorf("[ERROR][hyperv][create] Unable to set DefaultFlowMinimumBandwidthAbsolute if bandwidth reservation mode is none")
		}
	}

	if defaultQueueVmmqQueuePairs < 1 {
		return fmt.Errorf("[ERROR][hyperv][create] defaultQueueVmmqQueuePairs must be greater then 0")
	}

	err = c.CreateVMSwitch(switchName, notes, allowManagementOS, embeddedTeamingEnabled, iovEnabled, packetDirectEnabled, bandwidthReservationMode, switchType, netAdapterNames, defaultFlowMinimumBandwidthAbsolute, defaultFlowMinimumBandwidthWeight, defaultQueueVmmqEnabled, defaultQueueVmmqQueuePairs, defaultQueueVrssEnabled)

	if err != nil {
		return err
	}

	d.SetId(switchName)
	log.Printf("[INFO][hyperv][create] created hyperv switch: %#v", d)

	return resourceHyperVNetworkSwitchRead(d, meta)
}

func resourceHyperVNetworkSwitchRead(d *schema.ResourceData, meta interface{}) (err error) {
	log.Printf("[INFO][hyperv][read] reading hyperv switch: %#v", d)
	c := meta.(*api.HypervClient)

	switchName := ""

	if v, ok := d.GetOk("name"); ok {
		switchName = v.(string)
	} else {
		return fmt.Errorf("[ERROR][hyperv][read] name argument is required")
	}

	s, err := c.GetVMSwitch(switchName)

	if err != nil {
		return err
	}

	log.Printf("[INFO][hyperv][read] retrieved network switch: %+v", s)

	if s.Name != switchName {
		d.SetId("")
		log.Printf("[INFO][hyperv][read] unable to read hyperv switch as it does not exist: %#v", switchName)
		return nil
	}

	d.Set("notes", s.Notes)
	d.Set("allow_management_os", s.AllowManagementOS)
	d.Set("enable_embedded_teaming", s.EmbeddedTeamingEnabled)
	d.Set("enable_iov", s.IovEnabled)
	d.Set("enable_packet_direct", s.PacketDirectEnabled)
	d.Set("minimum_bandwidth_mode", s.BandwidthReservationMode.String())
	d.Set("switch_type", s.SwitchType.String())
	d.Set("net_adapter_names", s.NetAdapterNames)
	d.Set("default_flow_minimum_bandwidth_absolute", s.DefaultFlowMinimumBandwidthAbsolute)
	d.Set("default_flow_minimum_bandwidth_weight", s.DefaultFlowMinimumBandwidthWeight)
	d.Set("default_queue_vmmq_enabled", s.DefaultQueueVmmqEnabled)
	d.Set("default_queue_vmmq_queue_pairs", s.DefaultQueueVmmqQueuePairs)
	d.Set("default_queue_vrss_enabled", s.DefaultQueueVrssEnabled)

	if s.SwitchType == api.VMSwitchType_Private {
		if s.AllowManagementOS {
			return fmt.Errorf("[ERROR][hyperv][read] Unable to set AllowManagementOS to true if switch type is private")
		}

		if len(s.NetAdapterNames) > 0 {
			return fmt.Errorf("[ERROR][hyperv][read] Unable to set NetAdapterNames when switch type is private")
		}
	} else if s.SwitchType == api.VMSwitchType_Internal {
		if !s.AllowManagementOS {
			return fmt.Errorf("[ERROR][hyperv][read] Unable to set AllowManagementOS to false if switch type is internal")
		}

		if len(s.NetAdapterNames) > 0 {
			return fmt.Errorf("[ERROR][hyperv][read] Unable to set NetAdapterNames when switch type is internal")
		}
	} else if s.SwitchType == api.VMSwitchType_External {
		if len(s.NetAdapterNames) < 1 {
			return fmt.Errorf("[ERROR][hyperv][read] Must specify NetAdapterNames if switch type is external")
		}
	}

	if s.BandwidthReservationMode == api.VMSwitchBandwidthMode_Absolute {
		if s.DefaultFlowMinimumBandwidthWeight != 0 {
			return fmt.Errorf("[ERROR][hyperv][read] DefaultFlowMinimumBandwidthWeight should be 0 if bandwidth reservation mode is absolute")
		}
		if s.DefaultFlowMinimumBandwidthAbsolute < 0 {
			return fmt.Errorf("[ERROR][hyperv][read] Bandwidth absolute must be 0 or greater")
		}
	} else if s.BandwidthReservationMode == api.VMSwitchBandwidthMode_Weight || (s.BandwidthReservationMode == api.VMSwitchBandwidthMode_Default && (!s.IovEnabled)) {
		if s.DefaultFlowMinimumBandwidthAbsolute != 0 {
			return fmt.Errorf("[ERROR][hyperv][read] DefaultFlowMinimumBandwidthAbsolute should be 0 if bandwidth reservation mode is weight")
		}
		if s.DefaultFlowMinimumBandwidthWeight < 1 || s.DefaultFlowMinimumBandwidthWeight > 100 {
			return fmt.Errorf("[ERROR][hyperv][read] Bandwidth weight must be between 1 and 100")
		}
	} else {
		if s.DefaultFlowMinimumBandwidthWeight != 0 {
			return fmt.Errorf("[ERROR][hyperv][read] DefaultFlowMinimumBandwidthWeight should be 0 if bandwidth reservation mode is none")
		}
		if s.DefaultFlowMinimumBandwidthAbsolute != 0 {
			return fmt.Errorf("[ERROR][hyperv][read] DefaultFlowMinimumBandwidthAbsolute should be 0 if bandwidth reservation mode is none")
		}
	}

	if s.DefaultQueueVmmqQueuePairs < 1 {
		return fmt.Errorf("[ERROR][hyperv][read] defaultQueueVmmqQueuePairs must be greater then 0")
	}

	if err != nil {
		return err
	}

	log.Printf("[INFO][hyperv][read] read hyperv switch: %#v", d)

	return nil
}

func resourceHyperVNetworkSwitchUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	log.Printf("[INFO][hyperv][update] updating hyperv switch: %#v", d)
	c := meta.(*api.HypervClient)

	switchName := ""

	if v, ok := d.GetOk("name"); ok {
		switchName = v.(string)
	} else {
		return fmt.Errorf("[ERROR][hyperv][update] name argument is required")
	}

	notes := (d.Get("notes")).(string)
	allowManagementOS := (d.Get("allow_management_os")).(bool)
	//embeddedTeamingEnabled := (d.Get("enable_embedded_teaming")).(bool)
	iovEnabled := (d.Get("enable_iov")).(bool)
	//packetDirectEnabled := (d.Get("enable_packet_direct")).(bool)
	bandwidthReservationMode := api.ToVMSwitchBandwidthMode((d.Get("minimum_bandwidth_mode")).(string))
	switchType := api.ToVMSwitchType((d.Get("switch_type")).(string))
	netAdapterNames := []string{}
	if raw, ok := d.GetOk("net_adapter_names"); ok {
		for _, v := range raw.([]interface{}) {
			netAdapterNames = append(netAdapterNames, v.(string))
		}
	}
	defaultFlowMinimumBandwidthAbsolute := int64((d.Get("default_flow_minimum_bandwidth_absolute")).(int))
	defaultFlowMinimumBandwidthWeight := int64((d.Get("default_flow_minimum_bandwidth_weight")).(int))
	defaultQueueVmmqEnabled := (d.Get("default_queue_vmmq_enabled")).(bool)
	defaultQueueVmmqQueuePairs := int32((d.Get("default_queue_vmmq_queue_pairs")).(int))
	defaultQueueVrssEnabled := (d.Get("default_queue_vrss_enabled")).(bool)

	if switchType == api.VMSwitchType_Private {
		if allowManagementOS {
			return fmt.Errorf("[ERROR][hyperv][update] Unable to set AllowManagementOS to true if switch type is private")
		}

		if len(netAdapterNames) > 0 {
			return fmt.Errorf("[ERROR][hyperv][update] Unable to set NetAdapterNames when switch type is private")
		}
	} else if switchType == api.VMSwitchType_Internal {
		if !allowManagementOS {
			return fmt.Errorf("[ERROR][hyperv][update] Unable to set AllowManagementOS to false if switch type is internal")
		}

		if len(netAdapterNames) > 0 {
			return fmt.Errorf("[ERROR][hyperv][update] Unable to set NetAdapterNames when switch type is internal")
		}
	} else if switchType == api.VMSwitchType_External {
		if len(netAdapterNames) < 1 {
			return fmt.Errorf("[ERROR][hyperv][update] Must specify NetAdapterNames if switch type is external")
		}
	}

	if bandwidthReservationMode == api.VMSwitchBandwidthMode_Absolute {
		if defaultFlowMinimumBandwidthWeight != 0 {
			return fmt.Errorf("[ERROR][hyperv][update] Unable to set DefaultFlowMinimumBandwidthWeight if bandwidth reservation mode is absolute")
		}
		if defaultFlowMinimumBandwidthAbsolute < 0 {
			return fmt.Errorf("[ERROR][hyperv][update] Bandwidth absolute must be 0 or greater")
		}
	} else if bandwidthReservationMode == api.VMSwitchBandwidthMode_Weight || (bandwidthReservationMode == api.VMSwitchBandwidthMode_Default && (!iovEnabled)) {
		if defaultFlowMinimumBandwidthAbsolute != 0 {
			return fmt.Errorf("[ERROR][hyperv][update] Unable to set DefaultFlowMinimumBandwidthAbsolute if bandwidth reservation mode is weight")
		}
		if defaultFlowMinimumBandwidthWeight < 1 || defaultFlowMinimumBandwidthWeight > 100 {
			return fmt.Errorf("[ERROR][hyperv][update] Bandwidth weight must be between 1 and 100")
		}
	} else {
		if defaultFlowMinimumBandwidthWeight != 0 {
			return fmt.Errorf("[ERROR][hyperv][update] Unable to set DefaultFlowMinimumBandwidthWeight if bandwidth reservation mode is none")
		}
		if defaultFlowMinimumBandwidthAbsolute != 0 {
			return fmt.Errorf("[ERROR][hyperv][update] Unable to set DefaultFlowMinimumBandwidthAbsolute if bandwidth reservation mode is none")
		}
	}

	if defaultQueueVmmqQueuePairs < 1 {
		return fmt.Errorf("[ERROR][hyperv][update] defaultQueueVmmqQueuePairs must be greater then 0")
	}

	err = c.UpdateVMSwitch(switchName, notes, allowManagementOS, switchType, netAdapterNames, defaultFlowMinimumBandwidthAbsolute, defaultFlowMinimumBandwidthWeight, defaultQueueVmmqEnabled, defaultQueueVmmqQueuePairs, defaultQueueVrssEnabled)

	if err != nil {
		return err
	}

	log.Printf("[INFO][hyperv][update] updated hyperv switch: %#v", d)

	return resourceHyperVNetworkSwitchRead(d, meta)
}

func resourceHyperVNetworkSwitchDelete(d *schema.ResourceData, meta interface{}) (err error) {
	log.Printf("[INFO][hyperv][delete] deleting hyperv switch: %#v", d)

	c := meta.(*api.HypervClient)

	switchName := ""

	if v, ok := d.GetOk("name"); ok {
		switchName = v.(string)
	} else {
		return fmt.Errorf("[ERROR][hyperv][delete] name argument is required")
	}

	err = c.DeleteVMSwitch(switchName)

	if err != nil {
		return err
	}

	log.Printf("[INFO][hyperv][delete] deleted hyperv switch: %#v", d)
	return nil
}
