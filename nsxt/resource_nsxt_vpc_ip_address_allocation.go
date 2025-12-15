// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/terraform-provider-nsxt/api/orgs/projects/vpcs"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/metadata"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var cliVpcIpAddressAllocationsClient = vpcs.NewIpAddressAllocationsClient

var vpcIpAddressAllocationIpAddressTypeValues = []string{
	model.VpcIpAddressAllocation_IP_ADDRESS_TYPE_IPV4,
	model.VpcIpAddressAllocation_IP_ADDRESS_TYPE_IPV6,
}

var vpcIpAddressAllocationIpAddressBlockVisibilityValues = []string{
	model.VpcIpAddressAllocation_IP_ADDRESS_BLOCK_VISIBILITY_EXTERNAL,
	model.VpcIpAddressAllocation_IP_ADDRESS_BLOCK_VISIBILITY_PRIVATE,
	model.VpcIpAddressAllocation_IP_ADDRESS_BLOCK_VISIBILITY_PRIVATE_TGW,
}

var vpcIpAddressAllocationSchema = map[string]*metadata.ExtendedSchema{
	"nsx_id":       metadata.GetExtendedSchema(getNsxIDSchema()),
	"path":         metadata.GetExtendedSchema(getPathSchema()),
	"display_name": metadata.GetExtendedSchema(getDisplayNameSchema()),
	"description":  metadata.GetExtendedSchema(getDescriptionSchema()),
	"revision":     metadata.GetExtendedSchema(getRevisionSchema()),
	"tag":          metadata.GetExtendedSchema(getTagsSchema()),
	"context":      metadata.GetExtendedSchema(getContextSchemaExtended(true, false, true, true)),
	"allocation_ips": {
		Schema: schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validateCidrOrIPOrRange(),
			Optional:     true,
			Computed:     true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "string",
			SdkFieldName: "AllocationIps",
			OmitIfEmpty:  true,
		},
	},
	"allocation_size": {
		Schema: schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "int",
			SdkFieldName: "AllocationSize",
			OmitIfEmpty:  true,
		},
	},
	"ip_address_type": {
		Schema: schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validation.StringInSlice(vpcIpAddressAllocationIpAddressTypeValues, false),
			Optional:     true,
			Default:      model.VpcIpAddressAllocation_IP_ADDRESS_TYPE_IPV4,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "string",
			SdkFieldName: "IpAddressType",
		},
	},
	"ip_address_block_visibility": {
		Schema: schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validation.StringInSlice(vpcIpAddressAllocationIpAddressBlockVisibilityValues, false),
			Optional:     true,
			Default:      model.VpcIpAddressAllocation_IP_ADDRESS_BLOCK_VISIBILITY_EXTERNAL,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "string",
			SdkFieldName: "IpAddressBlockVisibility",
		},
	},
	"ip_block": {
		Schema: schema.Schema{
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validatePolicyPath(),
		},
		Metadata: metadata.Metadata{
			SchemaType:   "string",
			SdkFieldName: "IpBlock",
		},
	},
}

var vpcIpAddressAllocationPathExample = "/orgs/[org]/projects/[project]/vpcs/[vpc]/ip-address-allocations/[allocation]"

func resourceNsxtVpcIpAddressAllocation() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtVpcIpAddressAllocationCreate,
		Read:   resourceNsxtVpcIpAddressAllocationRead,
		Update: resourceNsxtVpcIpAddressAllocationUpdate,
		Delete: resourceNsxtVpcIpAddressAllocationDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtVersionCheckImporter("9.0.0", "VPC IP Address Allocation", getVpcPathResourceImporter(vpcIpAddressAllocationPathExample)),
		},
		Schema: metadata.GetSchemaFromExtendedSchema(vpcIpAddressAllocationSchema),
	}
}

func resourceNsxtVpcIpAddressAllocationExists(sessionContext utl.SessionContext, id string, connector client.Connector) (bool, error) {
	var err error
	parents := getVpcParentsFromContext(sessionContext)
	client := cliVpcIpAddressAllocationsClient(sessionContext, connector)
	_, err = client.Get(parents[0], parents[1], parents[2], id)
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving resource", err)
}

func resourceNsxtVpcIpAddressAllocationCreate(d *schema.ResourceData, m interface{}) error {
	if !util.NsxVersionHigherOrEqual("9.0.0") {
		return fmt.Errorf("VPC IP Address Allocation resource requires NSX version 9.0.0 or higher")
	}
	connector := getPolicyConnector(m)

	id, err := getOrGenerateID2(d, m, resourceNsxtVpcIpAddressAllocationExists)
	if err != nil {
		return err
	}

	parents := getVpcParentsFromContext(getSessionContext(d, m))
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)

	obj := model.VpcIpAddressAllocation{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, vpcIpAddressAllocationSchema, "", nil); err != nil {
		return err
	}

	log.Printf("[INFO] Creating VpcIpAddressAllocation with ID %s", id)

	sessionContext := getSessionContext(d, m)
	client := cliVpcIpAddressAllocationsClient(sessionContext, connector)
	err = client.Patch(parents[0], parents[1], parents[2], id, obj)
	if err != nil {
		return handleCreateError("VpcIpAddressAllocation", id, err)
	}
	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtVpcIpAddressAllocationRead(d, m)
}

func resourceNsxtVpcIpAddressAllocationRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VpcIpAddressAllocation ID")
	}

	sessionContext := getSessionContext(d, m)
	client := cliVpcIpAddressAllocationsClient(sessionContext, connector)
	parents := getVpcParentsFromContext(sessionContext)
	obj, err := client.Get(parents[0], parents[1], parents[2], id)
	if err != nil {
		return handleReadError(d, "VpcIpAddressAllocation", id, err)
	}

	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("revision", obj.Revision)
	d.Set("path", obj.Path)

	elem := reflect.ValueOf(&obj).Elem()
	return metadata.StructToSchema(elem, d, vpcIpAddressAllocationSchema, "", nil)
}

func resourceNsxtVpcIpAddressAllocationUpdate(d *schema.ResourceData, m interface{}) error {

	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VpcIpAddressAllocation ID")
	}

	parents := getVpcParentsFromContext(getSessionContext(d, m))
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getPolicyTagsFromSchema(d)

	revision := int64(d.Get("revision").(int))

	obj := model.VpcIpAddressAllocation{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
		Revision:    &revision,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, vpcIpAddressAllocationSchema, "", nil); err != nil {
		return err
	}
	sessionContext := getSessionContext(d, m)
	client := cliVpcIpAddressAllocationsClient(sessionContext, connector)
	_, err := client.Update(parents[0], parents[1], parents[2], id, obj)
	if err != nil {
		// Trigger partial update to avoid terraform updating state based on failed intent
		// TODO - move this into handleUpdateError
		d.Partial(true)
		return handleUpdateError("VpcIpAddressAllocation", id, err)
	}

	return resourceNsxtVpcIpAddressAllocationRead(d, m)
}

func resourceNsxtVpcIpAddressAllocationDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VpcIpAddressAllocation ID")
	}

	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)

	client := cliVpcIpAddressAllocationsClient(sessionContext, connector)
	err := client.Delete(parents[0], parents[1], parents[2], id)

	if err != nil {
		return handleDeleteError("VpcIpAddressAllocation", id, err)
	}

	return nil
}
