// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/terraform-provider-nsxt/api/orgs/projects/vpcs"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/metadata"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var cliVpcEndpointsClient = vpcs.NewVpcEndpointsClient

var vpcEndpointSchema = map[string]*metadata.ExtendedSchema{
	"nsx_id":       metadata.GetExtendedSchema(getNsxIDSchema()),
	"path":         metadata.GetExtendedSchema(getPathSchema()),
	"display_name": metadata.GetExtendedSchema(getDisplayNameSchema()),
	"description":  metadata.GetExtendedSchema(getDescriptionSchema()),
	"revision":     metadata.GetExtendedSchema(getRevisionSchema()),
	"tag":          metadata.GetExtendedSchema(getTagsSchema()),
	"context":      metadata.GetExtendedSchema(getContextSchemaExtended(true, false, true, true)),
	"vpc_service_endpoint": {
		Schema: schema.Schema{
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validatePolicyPath(),
			Description:  "Policy path to the VPC service endpoint being consumed.",
		},
		Metadata: metadata.Metadata{
			SchemaType:   "string",
			SdkFieldName: "VpcServiceEndpoint",
		},
	},
	"ip_allocation_path": {
		Schema: schema.Schema{
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validatePolicyPath(),
			Description:  "Policy path to the VPC IP address allocation that supplies the client endpoint IP address.",
		},
		Metadata: metadata.Metadata{
			SchemaType:   "string",
			SdkFieldName: "IpAllocationPath",
		},
	},
}

var vpcEndpointPathExample = "/orgs/[org]/projects/[project]/vpcs/[vpc]/vpc-endpoints/[endpoint]"

func resourceNsxtVpcEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtVpcEndpointCreate,
		Read:   resourceNsxtVpcEndpointRead,
		Update: resourceNsxtVpcEndpointUpdate,
		Delete: resourceNsxtVpcEndpointDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtVersionCheckImporter("9.2.0", "VPC Endpoint", getVpcPathResourceImporter(vpcEndpointPathExample)),
		},
		Schema: metadata.GetSchemaFromExtendedSchema(vpcEndpointSchema),
	}
}

func resourceNsxtVpcEndpointExists(sessionContext utl.SessionContext, id string, connector client.Connector) (bool, error) {
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcEndpointsClient(sessionContext, connector)
	_, err := c.Get(parents[0], parents[1], parents[2], id)
	if err == nil {
		return true, nil
	}
	if isNotFoundError(err) {
		return false, nil
	}
	return false, logAPIError("Error retrieving resource", err)
}

func resourceNsxtVpcEndpointCreate(d *schema.ResourceData, m interface{}) error {
	if !util.NsxVersionHigherOrEqual("9.2.0") {
		return fmt.Errorf("VPC Endpoint resource requires NSX version 9.2.0 or higher")
	}
	connector := getPolicyConnector(m)

	id, err := getOrGenerateID2(d, m, resourceNsxtVpcEndpointExists)
	if err != nil {
		return err
	}

	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)

	obj := model.VpcEndpoint{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, vpcEndpointSchema, "", nil); err != nil {
		return err
	}

	log.Printf("[INFO] Creating VpcEndpoint with ID %s", id)

	c := cliVpcEndpointsClient(sessionContext, connector)
	err = c.Patch(parents[0], parents[1], parents[2], id, obj)
	if err != nil {
		return handleCreateError("VpcEndpoint", id, err)
	}
	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtVpcEndpointRead(d, m)
}

func resourceNsxtVpcEndpointRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VpcEndpoint ID")
	}

	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcEndpointsClient(sessionContext, connector)
	obj, err := c.Get(parents[0], parents[1], parents[2], id)
	if err != nil {
		return handleReadError(d, "VpcEndpoint", id, err)
	}

	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("revision", obj.Revision)
	d.Set("path", obj.Path)

	elem := reflect.ValueOf(&obj).Elem()
	return metadata.StructToSchema(elem, d, vpcEndpointSchema, "", nil)
}

func resourceNsxtVpcEndpointUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VpcEndpoint ID")
	}

	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)
	revision := int64(d.Get("revision").(int))

	obj := model.VpcEndpoint{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
		Revision:    &revision,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, vpcEndpointSchema, "", nil); err != nil {
		return err
	}

	c := cliVpcEndpointsClient(sessionContext, connector)
	_, err := c.Update(parents[0], parents[1], parents[2], id, obj)
	if err != nil {
		d.Partial(true)
		return handleUpdateError("VpcEndpoint", id, err)
	}

	return resourceNsxtVpcEndpointRead(d, m)
}

func resourceNsxtVpcEndpointDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VpcEndpoint ID")
	}

	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)

	c := cliVpcEndpointsClient(sessionContext, connector)
	err := c.Delete(parents[0], parents[1], parents[2], id)
	if err != nil {
		return handleDeleteError("VpcEndpoint", id, err)
	}

	return nil
}
