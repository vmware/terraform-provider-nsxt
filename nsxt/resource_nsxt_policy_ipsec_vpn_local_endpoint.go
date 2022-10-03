/* Copyright © 2020 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	t0_service "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/tier_0s/locale_services/ipsec_vpn_services"
	t1_service "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/tier_1s/locale_services/ipsec_vpn_services"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func resourceNsxtPolicyIPSecVpnLocalEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyIPSecVpnLocalEndpointCreate,
		Read:   resourceNsxtPolicyIPSecVpnLocalEndpointRead,
		Update: resourceNsxtPolicyIPSecVpnLocalEndpointUpdate,
		Delete: resourceNsxtPolicyIPSecVpnLocalEndpointDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtVPNServiceResourceImporter,
		},

		Schema: map[string]*schema.Schema{
			"nsx_id":       getNsxIDSchema(),
			"path":         getPathSchema(),
			"display_name": getDisplayNameSchema(),
			"description":  getDescriptionSchema(),
			"revision":     getRevisionSchema(),
			"tag":          getTagsSchema(),
			"service_path": getPolicyPathSchema(true, true, "Policy path for IPSec VPN service"),
			"local_address": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsIPv4Address,
			},
			"certificate_path": getPolicyPathSchema(false, false, "Policy path referencing site certificate"),
			"local_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"trust_ca_paths": {
				Type:     schema.TypeSet,
				Elem:     getElemPolicyPathSchema(),
				Optional: true,
			},
			"trust_crl_paths": {
				Type:     schema.TypeSet,
				Elem:     getElemPolicyPathSchema(),
				Optional: true,
			},
		},
	}
}

func resourceNsxtPolicyIPSecVpnLocalEndpointExistsOnService(id string, connector *client.RestConnector, servicePath string) (bool, error) {
	isT0, gwID, localeServiceID, serviceID, err := parseIPSecVPNServicePolicyPath(servicePath)
	if err != nil {
		return false, err
	}
	if isT0 {
		client := t0_service.NewLocalEndpointsClient(connector)
		_, err = client.Get(gwID, localeServiceID, serviceID, id)
	} else {
		client := t1_service.NewLocalEndpointsClient(connector)
		_, err = client.Get(gwID, localeServiceID, serviceID, id)
	}
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving resource", err)
}

func resourceNsxtPolicyIPSecVpnLocalEndpointExists(servicePath string) func(id string, connector *client.RestConnector, isGlobalManager bool) (bool, error) {
	return func(id string, connector *client.RestConnector, isGlobalManager bool) (bool, error) {
		return resourceNsxtPolicyIPSecVpnLocalEndpointExistsOnService(id, connector, servicePath)
	}
}

func ipSecVpnLocalEndpointInitStruct(d *schema.ResourceData) model.IPSecVpnLocalEndpoint {
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)
	certificatePath := d.Get("certificate_path").(string)
	localAddress := d.Get("local_address").(string)
	localID := d.Get("local_id").(string)
	trustCaPaths := getStringListFromSchemaSet(d, "trust_ca_paths")
	trustCrlPaths := getStringListFromSchemaSet(d, "trust_crl_paths")

	obj := model.IPSecVpnLocalEndpoint{
		DisplayName:  &displayName,
		Description:  &description,
		Tags:         tags,
		LocalAddress: &localAddress,
	}

	if len(localID) > 0 {
		obj.LocalId = &localID
	}

	if len(certificatePath) > 0 {
		obj.CertificatePath = &certificatePath
	}

	if len(trustCaPaths) > 0 {
		obj.TrustCaPaths = trustCaPaths
	}

	if len(trustCrlPaths) > 0 {
		obj.TrustCrlPaths = trustCrlPaths
	}

	return obj
}

func resourceNsxtPolicyIPSecVpnLocalEndpointCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	servicePath := d.Get("service_path").(string)
	isT0, gwID, localeServiceID, serviceID, err := parseIPSecVPNServicePolicyPath(servicePath)
	if err != nil {
		return err
	}

	id, err := getOrGenerateID(d, m, resourceNsxtPolicyIPSecVpnLocalEndpointExists(servicePath))
	if err != nil {
		return err
	}

	obj := ipSecVpnLocalEndpointInitStruct(d)

	log.Printf("[INFO] Creating IPSecVpnLocalEndpoint with ID %s", id)
	if isT0 {
		client := t0_service.NewLocalEndpointsClient(connector)
		err = client.Patch(gwID, localeServiceID, serviceID, id, obj)
	} else {
		client := t1_service.NewLocalEndpointsClient(connector)
		err = client.Patch(gwID, localeServiceID, serviceID, id, obj)
	}
	if err != nil {
		return handleCreateError("IPSecVpnLocalEndpoint", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyIPSecVpnLocalEndpointRead(d, m)
}

func resourceNsxtPolicyIPSecVpnLocalEndpointRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining IPSecVpnLocalEndpoint ID")
	}

	servicePath := d.Get("service_path").(string)
	isT0, gwID, localeServiceID, serviceID, err := parseIPSecVPNServicePolicyPath(servicePath)
	if err != nil {
		return err
	}

	var obj model.IPSecVpnLocalEndpoint
	if isT0 {
		client := t0_service.NewLocalEndpointsClient(connector)
		obj, err = client.Get(gwID, localeServiceID, serviceID, id)
	} else {
		client := t1_service.NewLocalEndpointsClient(connector)
		obj, err = client.Get(gwID, localeServiceID, serviceID, id)
	}
	if err != nil {
		return handleReadError(d, "IPSecVpnLocalEndpoint", id, err)
	}

	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("path", obj.Path)
	d.Set("revision", obj.Revision)
	d.Set("certificate_path", obj.CertificatePath)
	d.Set("local_address", obj.LocalAddress)
	d.Set("local_id", obj.LocalId)
	d.Set("trust_ca_paths", obj.TrustCaPaths)
	d.Set("trust_crl_paths", obj.TrustCrlPaths)

	return nil
}

func resourceNsxtPolicyIPSecVpnLocalEndpointUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining IPSecVpnLocalEndpoint ID")
	}

	servicePath := d.Get("service_path").(string)
	isT0, gwID, localeServiceID, serviceID, err := parseIPSecVPNServicePolicyPath(servicePath)
	if err != nil {
		return err
	}

	revision := int64(d.Get("revision").(int))

	obj := ipSecVpnLocalEndpointInitStruct(d)
	obj.Revision = &revision

	log.Printf("[INFO] Updating IPSecVpnLocalEndpoint with ID %s", id)
	if isT0 {
		client := t0_service.NewLocalEndpointsClient(connector)
		err = client.Patch(gwID, localeServiceID, serviceID, id, obj)
	} else {
		client := t1_service.NewLocalEndpointsClient(connector)
		err = client.Patch(gwID, localeServiceID, serviceID, id, obj)
	}
	if err != nil {
		return handleUpdateError("IPSecVpnLocalEndpoint", id, err)
	}

	return resourceNsxtPolicyIPSecVpnLocalEndpointRead(d, m)
}

func resourceNsxtPolicyIPSecVpnLocalEndpointDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining IPSecVpnLocalEndpoint ID")
	}

	servicePath := d.Get("service_path").(string)
	isT0, gwID, localeServiceID, serviceID, err := parseIPSecVPNServicePolicyPath(servicePath)
	if err != nil {
		return err
	}
	connector := getPolicyConnector(m)
	if isT0 {
		client := t0_service.NewLocalEndpointsClient(connector)
		err = client.Delete(gwID, localeServiceID, serviceID, id)
	} else {
		client := t1_service.NewLocalEndpointsClient(connector)
		err = client.Delete(gwID, localeServiceID, serviceID, id)
	}

	if err != nil {
		return handleDeleteError("IPSecVpnLocalEndpoint", id, err)
	}

	return nil
}

func nsxtVPNServiceResourceImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	importID := d.Id()
	err := fmt.Errorf("VPN Local Endpoint Path expected, got %s", importID)
	// path format example: /infra/tier-1s/aaa/locale-services/default/ipsec-vpn-services/bbb/local-endpoints/ccc
	s := strings.Split(importID, "/")
	if len(s) < 10 {
		return []*schema.ResourceData{d}, err
	}

	endpointID := s[9]
	d.SetId(endpointID)
	s = strings.Split(importID, "/local-endpoints/")
	if len(s) != 2 {
		return []*schema.ResourceData{d}, err
	}

	d.Set("service_path", s[0])

	return []*schema.ResourceData{d}, nil
}
