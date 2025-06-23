// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"

	"github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

var visibilityTypes = []string{
	model.IpAddressBlock_VISIBILITY_EXTERNAL,
	model.IpAddressBlock_VISIBILITY_PRIVATE,
}

var ipBlockPathExample = getMultitenancyPathExample("/infra/ip-blocks/[ip-block]")

func resourceNsxtPolicyIPBlock() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyIPBlockCreate,
		Read:   resourceNsxtPolicyIPBlockRead,
		Update: resourceNsxtPolicyIPBlockUpdate,
		Delete: resourceNsxtPolicyIPBlockDelete,
		Importer: &schema.ResourceImporter{
			State: getPolicyPathOrIDResourceImporter(ipBlockPathExample),
		},

		Schema: map[string]*schema.Schema{
			"nsx_id":       getNsxIDSchema(),
			"path":         getPathSchema(),
			"display_name": getDisplayNameSchema(),
			"description":  getDescriptionSchema(),
			"revision":     getRevisionSchema(),
			"tag":          getTagsSchema(),
			"context":      getContextSchema(false, false, false),
			"cidr": {
				Type:         schema.TypeString,
				Description:  "Network address and the prefix length which will be associated with a layer-2 broadcast domain",
				Optional:     true,
				ValidateFunc: validateCidr(),
			},
			"visibility": {
				Type:         schema.TypeString,
				Description:  "Visibility of the Ip Block. Cannot be updated once associated with other resources.",
				Optional:     true,
				ValidateFunc: validation.StringInSlice(visibilityTypes, false),
			},
			"cidr_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Array of contiguous IP address spaces represented by network address and prefix length",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateCidr(),
				},
				ConflictsWith: []string{"cidr"},
			},
			"is_subnet_exclusive": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If this property is set to true, then this block is reserved for direct vlan extension use case",
			},
			"range_list":   getAllocationRangeListSchema(false, "Represents list of IP address ranges in the form of start and end IPs"),
			"reserved_ips": getAllocationRangeListSchema(false, "Represents list of reserved IP address in the form of start and end IPs"),
		},
	}
}

func resourceNsxtPolicyIPBlockExists(sessionContext utl.SessionContext, id string, connector client.Connector) (bool, error) {
	client := infra.NewIpBlocksClient(sessionContext, connector)
	if client == nil {
		return false, policyResourceNotSupportedError()
	}

	_, err := client.Get(id, nil)
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving IP Block", err)
}

func resourceNsxtPolicyIPBlockRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := infra.NewIpBlocksClient(getSessionContext(d, m), connector)
	if client == nil {
		return policyResourceNotSupportedError()
	}

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining IP Block ID")
	}

	block, err := client.Get(id, nil)
	if err != nil {
		return handleReadError(d, "IP Block", id, err)
	}

	d.Set("display_name", block.DisplayName)
	d.Set("description", block.Description)
	setPolicyTagsInSchema(d, block.Tags)
	d.Set("nsx_id", block.Id)
	d.Set("path", block.Path)
	d.Set("revision", block.Revision)
	if util.NsxVersionHigherOrEqual("4.2.0") {
		d.Set("visibility", block.Visibility)
	}
	if util.NsxVersionHigherOrEqual("9.1.0") {
		if len(d.Get("cidr_list").([]interface{})) > 0 {
			d.Set("cidr_list", block.CidrList)
		}
		d.Set("is_subnet_exclusive", block.IsSubnetExclusive)
		d.Set("range_list", setAllocationRangeListInSchema(block.RangeList))
		d.Set("reserved_ips", setAllocationRangeListInSchema(block.ReservedIps))
		if block.Cidr != nil {
			d.Set("cidr", block.Cidr)
		}
	} else {
		d.Set("cidr", block.Cidr)
	}

	return nil
}

func resourceNsxtPolicyIPBlockCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := infra.NewIpBlocksClient(getSessionContext(d, m), connector)
	if client == nil {
		return policyResourceNotSupportedError()
	}

	id, err := getOrGenerateID2(d, m, resourceNsxtPolicyIPBlockExists)
	if err != nil {
		return err
	}

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	cidr := d.Get("cidr").(string)
	visibility := d.Get("visibility").(string)
	tags := getPolicyTagsFromSchema(d)
	cidrList := getStringListFromSchemaList(d, "cidr_list")
	isSubnetExclusive := d.Get("is_subnet_exclusive").(bool)
	rangeList := getAllocationRangeListFromSchema(d.Get("range_list").([]interface{}))
	reservedIPs := getAllocationRangeListFromSchema(d.Get("reserved_ips").([]interface{}))

	obj := model.IpAddressBlock{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
	}
	if util.NsxVersionHigherOrEqual("4.2.0") && len(visibility) > 0 {
		obj.Visibility = &visibility
	}
	if util.NsxVersionHigherOrEqual("9.1.0") && len(cidrList) > 0 {
		obj.CidrList = cidrList
		obj.RangeList = rangeList
		obj.ReservedIps = reservedIPs
	} else if cidr != "" {
		obj.Cidr = &cidr
	}
	if util.NsxVersionHigherOrEqual("9.1.0") {
		obj.IsSubnetExclusive = &isSubnetExclusive
	}

	// Create the resource using PATCH
	log.Printf("[INFO] Creating IP Block with ID %s", id)
	err = client.Patch(id, obj)
	if err != nil {
		return handleCreateError("IP Block", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)
	return resourceNsxtPolicyIPBlockRead(d, m)
}

func resourceNsxtPolicyIPBlockUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := infra.NewIpBlocksClient(getSessionContext(d, m), connector)
	if client == nil {
		return policyResourceNotSupportedError()
	}

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining IP Block ID")
	}

	// Read the rest of the configured parameters
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	cidr := d.Get("cidr").(string)
	visibility := d.Get("visibility").(string)
	revision := int64(d.Get("revision").(int))
	tags := getPolicyTagsFromSchema(d)
	cidrList := getStringListFromSchemaList(d, "cidr_list")
	isSubnetExclusive := d.Get("is_subnet_exclusive").(bool)
	rangeList := getAllocationRangeListFromSchema(d.Get("range_list").([]interface{}))
	reservedIPs := getAllocationRangeListFromSchema(d.Get("reserved_ips").([]interface{}))

	obj := model.IpAddressBlock{
		Id:          &id,
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
		Revision:    &revision,
	}
	if util.NsxVersionHigherOrEqual("4.2.0") && len(visibility) > 0 {
		obj.Visibility = &visibility
	}
	if util.NsxVersionHigherOrEqual("9.1.0") && len(cidrList) > 0 {
		obj.CidrList = cidrList
		obj.RangeList = rangeList
		obj.ReservedIps = reservedIPs
	} else if cidr != "" {
		obj.Cidr = &cidr
	}
	if util.NsxVersionHigherOrEqual("9.1.0") {
		obj.IsSubnetExclusive = &isSubnetExclusive
	}

	_, err := client.Update(id, obj)
	if err != nil {
		return handleUpdateError("IP Block", id, err)
	}
	return resourceNsxtPolicyIPBlockRead(d, m)

}

func resourceNsxtPolicyIPBlockDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining IP Block ID")
	}

	connector := getPolicyConnector(m)
	client := infra.NewIpBlocksClient(getSessionContext(d, m), connector)
	if client == nil {
		return policyResourceNotSupportedError()
	}
	err := client.Delete(id)
	if err != nil {
		return handleDeleteError("IP Block", id, err)
	}

	return nil

}
