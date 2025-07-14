// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/metadata"
	tgwclientLayer "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/orgs/projects/transit_gateways/ipsec_vpn_services"
)

var tgwIPSecVpnSessionPathExample = getMultitenancyPathExample("/orgs/[org]/projects/[project]/transit-gateways/[transit-gateway]/ipsec-vpn-sessions/[ipsec-vpn-session]")

var transitGatewayIPSecVpnSessionSchema = map[string]*metadata.ExtendedSchema{
	"nsx_id":                     IPSecVpnSessionSchema["nsx_id"],
	"path":                       IPSecVpnSessionSchema["path"],
	"display_name":               IPSecVpnSessionSchema["display_name"],
	"description":                IPSecVpnSessionSchema["description"],
	"revision":                   IPSecVpnSessionSchema["revision"],
	"tag":                        IPSecVpnSessionSchema["tag"],
	"tunnel_profile_path":        IPSecVpnSessionSchema["tunnel_profile_path"],
	"local_endpoint_path":        IPSecVpnSessionSchema["local_endpoint_path"],
	"ike_profile_path":           IPSecVpnSessionSchema["ike_profile_path"],
	"dpd_profile_path":           IPSecVpnSessionSchema["dpd_profile_path"],
	"vpn_type":                   IPSecVpnSessionSchema["vpn_type"],
	"compliance_suite":           IPSecVpnSessionSchema["compliance_suite"],
	"connection_initiation_mode": IPSecVpnSessionSchema["connection_initiation_mode"],
	"authentication_mode":        IPSecVpnSessionSchema["authentication_mode"],
	"enabled":                    IPSecVpnSessionSchema["enabled"],
	"psk":                        IPSecVpnSessionSchema["psk"],
	"peer_id":                    IPSecVpnSessionSchema["peer_id"],
	"peer_address":               IPSecVpnSessionSchema["peer_address"],
	"ip_addresses":               IPSecVpnSessionSchema["ip_addresses"],
	"rule":                       IPSecVpnSessionSchema["rule"],
	"prefix_length":              IPSecVpnSessionSchema["prefix_length"],
	"direction":                  IPSecVpnSessionSchema["direction"],
	"max_segment_size":           IPSecVpnSessionSchema["max_segment_size"],
	"parent_path":                metadata.GetExtendedSchema(getPolicyPathSchema(true, true, "Policy path of the parent")),
}

func resourceNsxtPolicyTransitGatewayIPSecVpnSession() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyTGWIPSecVpnSessionCreate,
		Read:   resourceNsxtPolicyTGWIPSecVpnSessionRead,
		Update: resourceNsxtPolicyTGWIPSecVpnSessionUpdate,
		Delete: resourceNsxtPolicyTGWIPSecVpnSessionDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtParentPathResourceImporter,
		},
		Schema: metadata.GetSchemaFromExtendedSchema(transitGatewayIPSecVpnSessionSchema),
	}
}

func resourceNsxtPolicyTGWIPSecVpnSessionExists(sessionContext utl.SessionContext, parentPath string, id string, connector client.Connector) (bool, error) {
	var err error
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 4, tgwIPSecVpnSessionPathExample)
	if pathErr != nil {
		return false, pathErr
	}
	client := tgwclientLayer.NewSessionsClient(connector)
	_, err = client.Get(parents[0], parents[1], parents[2], parents[3], id)
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("error retrieving TransitGatewayIPSecVpnSession resource", err)
}

func resourceNsxtPolicyTGWIPSecVpnSessionCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id, err := getOrGenerateIDWithParent(d, m, resourceNsxtPolicyTGWIPSecVpnSessionExists)
	if err != nil {
		return err
	}

	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 4, tgwIPSecVpnSessionPathExample)
	if pathErr != nil {
		return pathErr
	}

	obj, err := getIPSecVPNSessionFromSchema(d)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating TransitGatewayIPSecVpnSession with ID %s", id)

	client := tgwclientLayer.NewSessionsClient(connector)
	err = client.Patch(parents[0], parents[1], parents[2], parents[3], id, obj)
	if err != nil {
		return handleCreateError("TransitGatewayIPSecVpnSession", id, err)
	}
	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyTGWIPSecVpnSessionRead(d, m)
}

func resourceNsxtPolicyTGWIPSecVpnSessionRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining TransitGatewayIPSecVpnSession ID")
	}

	client := tgwclientLayer.NewSessionsClient(connector)
	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 4, tgwIPSecVpnSessionPathExample)
	if pathErr != nil {
		return pathErr
	}

	obj, err := client.Get(parents[0], parents[1], parents[2], parents[3], id)
	if err != nil {
		return handleReadError(d, "TransitGatewayIPSecVpnSession", id, err)
	}

	return setIPSecVpnSessionResourceData(d, obj)
}

func resourceNsxtPolicyTGWIPSecVpnSessionUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining TransitGatewayIPSecVpnSession ID")
	}

	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 4, tgwIPSecVpnSessionPathExample)
	if pathErr != nil {
		return pathErr
	}

	obj, err := getIPSecVPNSessionFromSchema(d)
	if err != nil {
		return err
	}

	client := tgwclientLayer.NewSessionsClient(connector)
	err = client.Patch(parents[0], parents[1], parents[2], parents[3], id, obj)
	if err != nil {
		return handleUpdateError("pooja client update call error TransitGatewayIPSecVpnSession", id, err)
	}

	d.Set("nsx_id", id)
	return resourceNsxtPolicyTGWIPSecVpnSessionRead(d, m)
}

func resourceNsxtPolicyTGWIPSecVpnSessionDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining TransitGatewayIPSecVpnSession ID")
	}

	connector := getPolicyConnector(m)
	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 4, tgwIPSecVpnSessionPathExample)
	if pathErr != nil {
		return pathErr
	}

	client := tgwclientLayer.NewSessionsClient(connector)
	err := client.Delete(parents[0], parents[1], parents[2], parents[3], id)

	if err != nil {
		return handleDeleteError("TransitGatewayIPSecVpnSession", id, err)
	}

	return nil
}
