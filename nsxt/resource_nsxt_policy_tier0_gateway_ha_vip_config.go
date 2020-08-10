/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	//"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	//gm_infra "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/global_infra"
	//gm_tier0s "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/global_infra/tier_0s"
	//gm_locale_services "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/global_infra/tier_0s/locale_services"
	//gm_model "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	//"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	//"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/tier_0s/locale_services"
	//"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"log"
	"strings"
)

func resourceNsxtPolicyTier0HAVipConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyTier0HAVipConfigCreate,
		Read:   resourceNsxtPolicyTier0HAVipConfigRead,
		Update: resourceNsxtPolicyTier0HAVipConfigUpdate,
		Delete: resourceNsxtPolicyTier0HAVipConfigDelete,
		Importer: &schema.ResourceImporter{
			State: resourceNsxtPolicyTier0HAVipConfigImport,
		},

		Schema: map[string]*schema.Schema{
			"config": {
				Type:        schema.TypeList,
				Description: "Tier0 HA VIP Config",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Description: "Flag to enable this HA VIP config",
							Optional:    true,
							Default:     true,
						},
						"external_interface_paths": {
							Type:        schema.TypeList,
							Description: "paths to Tier0 external interfaces which are to be paired to provide redundancy",
							Required:    true,
							MinItems:    2,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: validatePolicyPath(),
							},
						},
						"vip_subnets": {
							Type:        schema.TypeList,
							Description: "IP address subnets which will be used as floating IP addresses",
							Required:    true,
							MinItems:    1,
							MaxItems:    2,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validateIPCidr(),
							},
						},
					},
				},
			},
			"locale_service_id": {
				Type:        schema.TypeString,
				Description: "Id of associated Gateway Locale Service on NSX",
				Computed:    true,
			},
			"tier0_id": {
				Type:        schema.TypeString,
				Description: "Id of associated Tier0 Gateway on NSX",
				Computed:    true,
			},
		},
	}
}

func resourceNsxtPolicyTier0HAVipConfigCreate(d *schema.ResourceData, m interface{}) error {
	//connector := getPolicyConnector(m)

	// Get the tier0 id and locale-service id from the interfaces
	configs := d.Get("config").([]interface{})
	for _, config := range configs {
		config_data := config.(map[string]interface{})
		interfaces := config_data["external_interface_paths"].(*schema.Set).List()
		example_interface := interfaces[0].(string)
		log.Printf("[ERROR] DEBUG ADIT path %s", example_interface)
		result := strings.Split(example_interface, "/")
		log.Printf("[ERROR] DEBUG ADIT result %v", result)
	}

	id := newUUID()
	d.SetId(id)
	//d.Set("tier0_id", tier0ID)
	//d.Set("locale_service_id", localeServiceID)

	return resourceNsxtPolicyTier0HAVipConfigRead(d, m)
}

func resourceNsxtPolicyTier0HAVipConfigRead(d *schema.ResourceData, m interface{}) error {
	//connector := getPolicyConnector(m)

	id := d.Id()
	tier0ID := d.Get("tier0_id").(string)
	localeServiceID := d.Get("locale_service_id").(string)
	if id == "" || tier0ID == "" || localeServiceID == ""{
		return fmt.Errorf("Error obtaining Tier0 id or Locale Service id")
	}

	return nil
}

func resourceNsxtPolicyTier0HAVipConfigUpdate(d *schema.ResourceData, m interface{}) error {
	//connector := getPolicyConnector(m)

	id := d.Id()
	tier0ID := d.Get("tier0_id").(string)
	localeServiceID := d.Get("locale_service_id").(string)
	if id == "" || tier0ID == "" || localeServiceID == "" {
		return fmt.Errorf("Error obtaining Tier0 id or Locale Service id")
	}

	return resourceNsxtPolicyTier0HAVipConfigRead(d, m)
}

func resourceNsxtPolicyTier0HAVipConfigDelete(d *schema.ResourceData, m interface{}) error {
	//connector := getPolicyConnector(m)

	id := d.Id()
	tier0ID := d.Get("tier0_id").(string)
	localeServiceID := d.Get("locale_service_id").(string)
	if id == "" || tier0ID == "" || localeServiceID == ""{
		return fmt.Errorf("Error obtaining Tier0 id or Locale Service id")
	}

	return nil
}

func resourceNsxtPolicyTier0HAVipConfigImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	importID := d.Id()
	s := strings.Split(importID, "/")
	if len(s) != 2 {
		return nil, fmt.Errorf("Please provide <gateway-id>/<locale-service-id> as an input")
	}

	gwID := s[0]
	localeServiceID := s[1]
	// connector := getPolicyConnector(m)
	// var tier0GW model.Tier0
	// if isPolicyGlobalManager(m) {
	// 	client := gm_infra.NewDefaultTier0sClient(connector)
	// 	gmObj, err1 := client.Get(gwID)
	// 	if err1 != nil {
	// 		return nil, err1
	// 	}

	// 	convertedObj, err2 := convertModelBindingType(gmObj, model.Tier0BindingType(), model.Tier0BindingType())
	// 	if err2 != nil {
	// 		return nil, err2
	// 	}

	// 	tier0GW = convertedObj.(model.Tier0)
	// } else {
	// 	client := infra.NewDefaultTier0sClient(connector)
	// 	var err error
	// 	tier0GW, err = client.Get(gwID)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	d.Set("tier0_id", gwID)
	d.Set("locale_service_id", localeServiceID)
	id := newUUID()
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
