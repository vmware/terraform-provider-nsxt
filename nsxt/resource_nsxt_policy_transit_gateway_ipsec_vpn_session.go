// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"

	ipsecvpnservices "github.com/vmware/terraform-provider-nsxt/api/orgs/projects/transit_gateways/ipsec_vpn_services"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/metadata"
)

var tgwIPSecVpnSessionPathExample = getMultitenancyPathExample("/orgs/[org]/projects/[project]/transit-gateways/[transit-gateway]/ipsec-vpn-sessions/[ipsec-vpn-session]")

var transitGatewayIPSecVpnSessionSchema = getIPSecVpnSessionCommon(true)

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
	client := ipsecvpnservices.NewSessionsClient(sessionContext, connector)
	if client == nil {
		return false, fmt.Errorf("unsupported client type")
	}
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

	sessionContext := getParentContext(d, m, parentPath)
	client := ipsecvpnservices.NewSessionsClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}
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

	parentPath := d.Get("parent_path").(string)
	sessionContext := getParentContext(d, m, parentPath)
	client := ipsecvpnservices.NewSessionsClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}
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

	sessionContext := getParentContext(d, m, parentPath)
	client := ipsecvpnservices.NewSessionsClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}
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

	sessionContext := getParentContext(d, m, parentPath)
	client := ipsecvpnservices.NewSessionsClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}
	err := client.Delete(parents[0], parents[1], parents[2], parents[3], id)

	if err != nil {
		return handleDeleteError("TransitGatewayIPSecVpnSession", id, err)
	}

	return nil
}
