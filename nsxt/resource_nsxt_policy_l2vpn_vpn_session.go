/* Copyright © 2020 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	//     "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	// gm_infra "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/global_infra"
	// gm_model "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	// "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	l2vpn_services "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/tier_0s/locale_services/l2vpn_services"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func resourceNsxtPolicyL2VPNSession() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyL2VPNSessionCreate,
		Read:   resourceNsxtPolicyL2VPNSessionRead,
		Update: resourceNsxtPolicyL2VPNSessionUpdate,
		Delete: resourceNsxtPolicyL2VPNSessionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"nsx_id":       getNsxIDSchema(),
			"path":         getPathSchema(),
			"display_name": getDisplayNameSchema(),
			"description":  getDescriptionSchema(),
			"revision":     getRevisionSchema(),
			"tag":          getTagsSchema(),
			"tier0_id": {
				Type:        schema.TypeString,
				Description: "Policy path referencing Local endpoint.",
				Optional:    true,
				Default:     "vmc",
			},
			"locale_service": {
				Type:        schema.TypeString,
				Description: "Local_service",
				Optional:    true,
				Default:     "default",
			},
			"service_id": {
				Type:        schema.TypeString,
				Description: "Policy path referencing Local endpoint.",
				Optional:    true,
				Default:     "default",
			},
			"transport_tunnels": {
				Type:        schema.TypeList,
				Description: "List of transport tunnels for redundancy",
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceNsxtPolicyL2VPNSessionCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	Tier0ID := d.Get("tier0_id").(string)
	LocaleService := d.Get("locale_service").(string)
	ServiceID := d.Get("service_id").(string)
	TransportTunnel := getStringListFromSchemaList(d, "transport_tunnels")

	// Initialize resource Id and verify this ID is not yet used
	id := "l2vpn_id"
	var err error
	if err != nil {
		return err
	}

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)

	obj := model.L2VPNSession{
		DisplayName:      &displayName,
		Description:      &description,
		Tags:             tags,
		TransportTunnels: TransportTunnel,
	}

	// Create the resource using PATCH
	log.Printf("[INFO] Creating L2VPNSession with ID %s", id)

	client := l2vpn_services.NewDefaultSessionsClient(connector)
	err = client.Patch(Tier0ID, LocaleService, ServiceID, id, obj)

	if err != nil {
		return handleCreateError("L2VPNSession", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyL2VPNSessionRead(d, m)
}

func resourceNsxtPolicyL2VPNSessionRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	Tier0ID := d.Get("tier0_id").(string)
	LocaleService := d.Get("locale_service").(string)
	ServiceID := d.Get("service_id").(string)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining L2VPNSession ID")
	}

	var obj model.L2VPNSession

	client := l2vpn_services.NewDefaultSessionsClient(connector)
	var err error
	obj, err = client.Get(Tier0ID, LocaleService, ServiceID, id)
	if err != nil {
		return handleReadError(d, "L2VPNSession", id, err)
	}

	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("path", obj.Path)
	d.Set("revision", obj.Revision)

	//    <!SET_OBJ_ATTRS_IN_SCHEMA!>

	return nil
}

func resourceNsxtPolicyL2VPNSessionUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	Tier0ID := d.Get("tier0_id").(string)
	LocaleService := d.Get("locale_service").(string)
	ServiceID := d.Get("service_id").(string)
	TransportTunnel := getStringListFromSchemaList(d, "transport_tunnels")

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining L2VPNSession ID")
	}

	// Read the rest of the configured parameters
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getPolicyTagsFromSchema(d)

	obj := model.L2VPNSession{
		DisplayName:      &displayName,
		Description:      &description,
		Tags:             tags,
		TransportTunnels: TransportTunnel,
	}

	// Update the resource using PATCH
	var err error

	client := l2vpn_services.NewDefaultSessionsClient(connector)
	err = client.Patch(Tier0ID, LocaleService, ServiceID, id, obj)

	if err != nil {
		return handleUpdateError("L2VPNSession", id, err)
	}

	return resourceNsxtPolicyL2VPNSessionRead(d, m)

}

func resourceNsxtPolicyL2VPNSessionDelete(d *schema.ResourceData, m interface{}) error {

	Tier0ID := d.Get("tier0_id").(string)
	LocaleService := d.Get("locale_service").(string)
	ServiceID := d.Get("service_id").(string)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining L2VPNSession ID")
	}

	connector := getPolicyConnector(m)
	var err error

	client := l2vpn_services.NewDefaultSessionsClient(connector)
	err = client.Delete(Tier0ID, LocaleService, ServiceID, id)

	if err != nil {
		return handleDeleteError("L2VPNSession", id, err)
	}

	return nil
}
