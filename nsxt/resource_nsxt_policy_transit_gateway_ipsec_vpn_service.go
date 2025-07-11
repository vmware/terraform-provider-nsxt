package nsxt

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/metadata"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	clientLayer "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/orgs/projects/transit_gateways"
)

var twIpsecVpnParentPathExample = "/orgs/[org]/projects/[project]/transit-gateways/[transit-gateway]"

var transitGatewayIpsecVpnServiceSchema = map[string]*metadata.ExtendedSchema{
	// Copy all fields from the base IPSecVpnServiceSchema
	"nsx_id":        IPSecVpnServiceSchema["nsx_id"],
	"path":          IPSecVpnServiceSchema["path"],
	"display_name":  IPSecVpnServiceSchema["display_name"],
	"description":   IPSecVpnServiceSchema["description"],
	"revision":      IPSecVpnServiceSchema["revision"],
	"tag":           IPSecVpnServiceSchema["tag"],
	"enabled":       IPSecVpnServiceSchema["enabled"],
	"ha_sync":       IPSecVpnServiceSchema["ha_sync"],
	"ike_log_level": IPSecVpnServiceSchema["ike_log_level"],
	"bypass_rule":   IPSecVpnServiceSchema["bypass_rule"],
	"parent_path":   metadata.GetExtendedSchema(getPolicyPathSchema(true, true, "Policy path of the parent")),
}

// add it in the provider.go
func resourceNsxtPolicyTransitGatewayIpsecVpn() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyTransitGatewayIPSecVpnServicesCreate,
		Read:   resourceNsxtPolicyTransitGatewayIPSecVpnServicesRead,
		Update: resourceNsxtPolicyTransitGatewayIPSecVpnServicesUpdate,
		Delete: resourceNsxtPolicyTransitGatewayIPSecVpnServicesDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtParentPathResourceImporter,
		},
		Schema: metadata.GetSchemaFromExtendedSchema(transitGatewayIpsecVpnServiceSchema),
	}
}

func resourceNsxtPolicyTransitGatewayIPSecVpnServicesCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	id, err := getOrGenerateIDWithParent(d, m, resourceNsxtPolicyTransitGatewayIPSecVpnServicesExists)
	if err != nil {
		return err
	}
	ipSecVpnService := getIpsecVpnServiceObject(id, d)

	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, twIpsecVpnParentPathExample)
	if pathErr != nil {
		return pathErr
	}

	client := clientLayer.NewIpsecVpnServicesClient(connector)
	_, err = client.Update(parents[0], parents[1], parents[2], id, ipSecVpnService)
	if err != nil {
		return handleUpdateError("Transit Gateway IPSecVpnService", id, err)
	}
	d.SetId(id)
	d.Set("nsx_id", id)
	return resourceNsxtPolicyTransitGatewayIPSecVpnServicesRead(d, m)
}

func resourceNsxtPolicyTransitGatewayIPSecVpnServicesRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining Transit Gateway IPSecVpnService ID")
	}

	client := clientLayer.NewIpsecVpnServicesClient(connector)
	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, twIpsecVpnParentPathExample)
	if pathErr != nil {
		return pathErr
	}
	obj, err := client.Get(parents[0], parents[1], parents[2], id)
	if err != nil {
		return handleReadError(d, "Transit Gateway IPSecVpnService", id, err)
	}
	setIpsecVPNServices(id, d, obj)

	return nil
}

func resourceNsxtPolicyTransitGatewayIPSecVpnServicesUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining Transit Gateway IPSecVpnService ID")
	}

	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, twIpsecVpnParentPathExample)
	if pathErr != nil {
		return pathErr
	}
	ipSecVpnService := getIpsecVpnServiceObject(id, d)
	revision := int64(d.Get("revision").(int))
	ipSecVpnService.Revision = &revision

	client := clientLayer.NewIpsecVpnServicesClient(connector)
	_, err := client.Update(parents[0], parents[1], parents[2], id, ipSecVpnService)
	if err != nil {
		return handleUpdateError("Transit Gateway IPSecVpnService", id, err)
	}

	return resourceNsxtPolicyTransitGatewayIPSecVpnServicesRead(d, m)
}

func resourceNsxtPolicyTransitGatewayIPSecVpnServicesDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining Transit Gateway IPSecVpnService ID")
	}

	connector := getPolicyConnector(m)
	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, twIpsecVpnParentPathExample)
	if pathErr != nil {
		return pathErr
	}

	client := clientLayer.NewIpsecVpnServicesClient(connector)
	err := client.Delete(parents[0], parents[1], parents[2], id)

	if err != nil {
		return handleDeleteError("Transit Gateway IPSecVpnService", id, err)
	}

	return nil
}

func resourceNsxtPolicyTransitGatewayIPSecVpnServicesExists(sessionContext utl.SessionContext, parentPath string, id string, connector client.Connector) (bool, error) {
	var err error
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, twIpsecVpnParentPathExample)
	if pathErr != nil {
		return false, pathErr
	}
	client := clientLayer.NewIpsecVpnServicesClient(connector)
	_, err = client.Get(parents[0], parents[1], parents[2], id)
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}
	return false, logAPIError("error retrieving Transit Gateway IPSecVpnService resource", err)
}
