package nsxt

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/metadata"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	tgwclientLayer "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/orgs/projects/transit_gateways"
)

var twIpsecVpnParentPathExample = "/orgs/[org]/projects/[project]/transit-gateways/[transit-gateway]"

var transitGatewayIpsecVpnServiceSchema = getIPSecVpnServiceCommonSchema(false, true)

// add it in the provider.go
func resourceNsxtPolicyTransitGatewayIpsecVpnService() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyTGWIPSecVpnServicesCreate,
		Read:   resourceNsxtPolicyTGWIPSecVpnServicesRead,
		Update: resourceNsxtPolicyTGWIPSecVpnServicesUpdate,
		Delete: resourceNsxtPolicyTGWIPSecVpnServicesDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtParentPathResourceImporter,
		},
		Schema: metadata.GetSchemaFromExtendedSchema(transitGatewayIpsecVpnServiceSchema),
	}
}

func resourceNsxtPolicyTGWIPSecVpnServicesCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	id, err := getOrGenerateIDWithParent(d, m, resourceNsxtPolicyTGWIPSecVpnServicesExists)
	if err != nil {
		return err
	}
	ipSecVpnService := getIpsecVpnServiceObject(id, d)

	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, twIpsecVpnParentPathExample)
	if pathErr != nil {
		return pathErr
	}

	client := tgwclientLayer.NewIpsecVpnServicesClient(connector)
	err = client.Patch(parents[0], parents[1], parents[2], id, ipSecVpnService)
	if err != nil {
		return handleCreateError("TransitGatewayIPSecVpnService", id, err)
	}
	d.SetId(id)
	d.Set("nsx_id", id)
	return resourceNsxtPolicyTGWIPSecVpnServicesRead(d, m)
}

func resourceNsxtPolicyTGWIPSecVpnServicesRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining TransitGatewayIPSecVpnService ID")
	}

	client := tgwclientLayer.NewIpsecVpnServicesClient(connector)
	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, twIpsecVpnParentPathExample)
	if pathErr != nil {
		return pathErr
	}
	obj, err := client.Get(parents[0], parents[1], parents[2], id)
	if err != nil {
		return handleReadError(d, "TransitGatewayIPSecVpnService", id, err)
	}
	setIpsecVPNServices(id, d, obj)

	return nil
}

func resourceNsxtPolicyTGWIPSecVpnServicesUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining TransitGatewayIPSecVpnService ID")
	}

	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, twIpsecVpnParentPathExample)
	if pathErr != nil {
		return pathErr
	}
	ipSecVpnService := getIpsecVpnServiceObject(id, d)
	revision := int64(d.Get("revision").(int))
	ipSecVpnService.Revision = &revision

	client := tgwclientLayer.NewIpsecVpnServicesClient(connector)
	_, err := client.Update(parents[0], parents[1], parents[2], id, ipSecVpnService)
	if err != nil {
		return handleUpdateError("TGWIPSecVpnService", id, err)
	}

	return resourceNsxtPolicyTGWIPSecVpnServicesRead(d, m)
}

func resourceNsxtPolicyTGWIPSecVpnServicesDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining TransitGatewayIPSecVpnService ID")
	}

	connector := getPolicyConnector(m)
	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, twIpsecVpnParentPathExample)
	if pathErr != nil {
		return pathErr
	}

	client := tgwclientLayer.NewIpsecVpnServicesClient(connector)
	err := client.Delete(parents[0], parents[1], parents[2], id)

	if err != nil {
		return handleDeleteError("TransitGatewayIPSecVpnService", id, err)
	}

	return nil
}

func resourceNsxtPolicyTGWIPSecVpnServicesExists(sessionContext utl.SessionContext, parentPath string, id string, connector client.Connector) (bool, error) {
	var err error
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, twIpsecVpnParentPathExample)
	if pathErr != nil {
		return false, pathErr
	}
	client := tgwclientLayer.NewIpsecVpnServicesClient(connector)
	_, err = client.Get(parents[0], parents[1], parents[2], id)
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}
	return false, logAPIError("error retrieving TransitGatewayIPSecVpnService resource", err)
}
