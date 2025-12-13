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

var iPSecVpnLocalEndpointPathExample = "/orgs/[org]/projects/[project]/transit-gateways/[transit-gateway]/ipsec-vpn-local-endpoints/[ipsec-vpn-local-endpoint]"

var tGwiPSecVpnLocalEndpointSchema = getIPSecVpnLocalEndpointCommonSchema(true)

func resourceNsxtPolicyTransitGatewayIPSecVpnLocalEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyTGWIPSecVpnLocalEndpointCreate,
		Read:   resourceNsxtPolicyTGWIPSecVpnLocalEndpointRead,
		Update: resourceNsxtPolicyTGWIPSecVpnLocalEndpointUpdate,
		Delete: resourceNsxtPolicyTGWIPSecVpnLocalEndpointDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtParentPathResourceImporter,
		},
		Schema: metadata.GetSchemaFromExtendedSchema(tGwiPSecVpnLocalEndpointSchema),
	}
}

func resourceNsxtPolicyTGWIPSecVpnLocalEndpointExists(sessionContext utl.SessionContext, parentPath string, id string, connector client.Connector) (bool, error) {
	var err error
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 4, iPSecVpnLocalEndpointPathExample)
	if pathErr != nil {
		return false, pathErr
	}
	client := ipsecvpnservices.NewLocalEndpointsClient(sessionContext, connector)
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

	return false, logAPIError("error retrieving TransitGatewayIPSecVpnLocalEndpoint resource", err)
}

func resourceNsxtPolicyTGWIPSecVpnLocalEndpointCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id, err := getOrGenerateIDWithParent(d, m, resourceNsxtPolicyTGWIPSecVpnLocalEndpointExists)
	if err != nil {
		return err
	}

	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 4, tgwIPSecVpnSessionPathExample)
	if pathErr != nil {
		return pathErr
	}

	obj := ipSecVpnLocalEndpointInitStruct(d)

	log.Printf("[INFO] Creating TransitGatewayIPSecVpnLocalEndpoint with ID %s", id)

	sessionContext := getParentContext(d, m, parentPath)
	client := ipsecvpnservices.NewLocalEndpointsClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}

	err = client.Patch(parents[0], parents[1], parents[2], parents[3], id, obj)
	if err != nil {
		return handleCreateError("TransitGatewayIPSecVpnLocalEndpoint", id, err)
	}
	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyTGWIPSecVpnLocalEndpointRead(d, m)
}

func resourceNsxtPolicyTGWIPSecVpnLocalEndpointRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining TransitGatewayIPSecVpnLocalEndpoint ID")
	}

	parentPath := d.Get("parent_path").(string)
	sessionContext := getParentContext(d, m, parentPath)
	client := ipsecvpnservices.NewLocalEndpointsClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 4, tgwIPSecVpnSessionPathExample)
	if pathErr != nil {
		return pathErr
	}
	obj, err := client.Get(parents[0], parents[1], parents[2], parents[3], id)
	if err != nil {
		return handleReadError(d, "TransitGatewayIPSecVpnLocalEndpoint", id, err)
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

func resourceNsxtPolicyTGWIPSecVpnLocalEndpointUpdate(d *schema.ResourceData, m interface{}) error {

	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining TransitGatewayIPSecVpnLocalEndpoint ID")
	}

	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 4, tgwIPSecVpnSessionPathExample)
	if pathErr != nil {
		return pathErr
	}

	obj := ipSecVpnLocalEndpointInitStruct(d)
	revision := int64(d.Get("revision").(int))
	obj.Revision = &revision

	sessionContext := getParentContext(d, m, parentPath)
	client := ipsecvpnservices.NewLocalEndpointsClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}
	_, err := client.Update(parents[0], parents[1], parents[2], parents[3], id, obj)
	if err != nil {
		return handleUpdateError("TransitGatewayIPSecVpnLocalEndpoint", id, err)
	}

	return resourceNsxtPolicyTGWIPSecVpnLocalEndpointRead(d, m)
}

func resourceNsxtPolicyTGWIPSecVpnLocalEndpointDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining TransitGatewayIPSecVpnLocalEndpoint ID")
	}

	connector := getPolicyConnector(m)
	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 4, tgwIPSecVpnSessionPathExample)
	if pathErr != nil {
		return pathErr
	}

	sessionContext := getParentContext(d, m, parentPath)
	client := ipsecvpnservices.NewLocalEndpointsClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}
	err := client.Delete(parents[0], parents[1], parents[2], parents[3], id)

	if err != nil {
		return handleDeleteError("TransitGatewayIPSecVpnLocalEndpoint", id, err)
	}

	return nil
}
