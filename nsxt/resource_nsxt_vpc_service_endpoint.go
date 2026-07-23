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

var cliVpcServiceEndpointsClient = vpcs.NewVpcServiceEndpointsClient

var vpcServiceEndpointIpTypeValues = []string{
	model.VpcServiceEndpoint_SERVICE_ENDPOINT_IP_TYPE_WORKLOAD,
}

var vpcServiceEndpointSchema = map[string]*metadata.ExtendedSchema{
	"nsx_id":       metadata.GetExtendedSchema(getNsxIDSchema()),
	"path":         metadata.GetExtendedSchema(getPathSchema()),
	"display_name": metadata.GetExtendedSchema(getDisplayNameSchema()),
	"description":  metadata.GetExtendedSchema(getDescriptionSchema()),
	"revision":     metadata.GetExtendedSchema(getRevisionSchema()),
	"tag":          metadata.GetExtendedSchema(getTagsSchema()),
	"context":      metadata.GetExtendedSchema(getContextSchemaExtended(true, false, true, true)),
	"service_endpoint_ip": {
		Schema: schema.Schema{
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateSingleIP(),
			Description:  "IP address of the VM providing the service. Must be a valid IPv4 address.",
		},
		Metadata: metadata.Metadata{
			SchemaType:   "string",
			SdkFieldName: "ServiceEndpointIp",
		},
	},
	"service_endpoint_ip_type": {
		Schema: schema.Schema{
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      model.VpcServiceEndpoint_SERVICE_ENDPOINT_IP_TYPE_WORKLOAD,
			ValidateFunc: validation.StringInSlice(vpcServiceEndpointIpTypeValues, false),
			Description:  "Indicates what type of resource the IP is assigned to.",
		},
		Metadata: metadata.Metadata{
			SchemaType:   "string",
			SdkFieldName: "ServiceEndpointIpType",
		},
	},
}

var vpcServiceEndpointPathExample = "/orgs/[org]/projects/[project]/vpcs/[vpc]/vpc-service-endpoints/[endpoint]"

func resourceNsxtVpcServiceEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtVpcServiceEndpointCreate,
		Read:   resourceNsxtVpcServiceEndpointRead,
		Update: resourceNsxtVpcServiceEndpointUpdate,
		Delete: resourceNsxtVpcServiceEndpointDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtVersionCheckImporter("9.2.0", "VPC Service Endpoint", getVpcPathResourceImporter(vpcServiceEndpointPathExample)),
		},
		Schema: metadata.GetSchemaFromExtendedSchema(vpcServiceEndpointSchema),
	}
}

func resourceNsxtVpcServiceEndpointExists(sessionContext utl.SessionContext, id string, connector client.Connector) (bool, error) {
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcServiceEndpointsClient(sessionContext, connector)
	_, err := c.Get(parents[0], parents[1], parents[2], id)
	if err == nil {
		return true, nil
	}
	if isNotFoundError(err) {
		return false, nil
	}
	return false, logAPIError("Error retrieving resource", err)
}

func resourceNsxtVpcServiceEndpointCreate(d *schema.ResourceData, m interface{}) error {
	if !util.NsxVersionHigherOrEqual("9.2.0") {
		return fmt.Errorf("VPC Service Endpoint resource requires NSX version 9.2.0 or higher")
	}
	connector := getPolicyConnector(m)

	id, err := getOrGenerateID2(d, m, resourceNsxtVpcServiceEndpointExists)
	if err != nil {
		return err
	}

	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)

	obj := model.VpcServiceEndpoint{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, vpcServiceEndpointSchema, "", nil); err != nil {
		return err
	}

	log.Printf("[INFO] Creating VpcServiceEndpoint with ID %s", id)

	c := cliVpcServiceEndpointsClient(sessionContext, connector)
	err = c.Patch(parents[0], parents[1], parents[2], id, obj)
	if err != nil {
		return handleCreateError("VpcServiceEndpoint", id, err)
	}
	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtVpcServiceEndpointRead(d, m)
}

func resourceNsxtVpcServiceEndpointRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VpcServiceEndpoint ID")
	}

	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcServiceEndpointsClient(sessionContext, connector)
	obj, err := c.Get(parents[0], parents[1], parents[2], id)
	if err != nil {
		return handleReadError(d, "VpcServiceEndpoint", id, err)
	}

	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("revision", obj.Revision)
	d.Set("path", obj.Path)

	elem := reflect.ValueOf(&obj).Elem()
	return metadata.StructToSchema(elem, d, vpcServiceEndpointSchema, "", nil)
}

func resourceNsxtVpcServiceEndpointUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VpcServiceEndpoint ID")
	}

	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)
	revision := int64(d.Get("revision").(int))

	obj := model.VpcServiceEndpoint{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
		Revision:    &revision,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, vpcServiceEndpointSchema, "", nil); err != nil {
		return err
	}

	c := cliVpcServiceEndpointsClient(sessionContext, connector)
	_, err := c.Update(parents[0], parents[1], parents[2], id, obj)
	if err != nil {
		d.Partial(true)
		return handleUpdateError("VpcServiceEndpoint", id, err)
	}

	return resourceNsxtVpcServiceEndpointRead(d, m)
}

func resourceNsxtVpcServiceEndpointDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VpcServiceEndpoint ID")
	}

	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)

	c := cliVpcServiceEndpointsClient(sessionContext, connector)
	err := c.Delete(parents[0], parents[1], parents[2], id)
	if err != nil {
		return handleDeleteError("VpcServiceEndpoint", id, err)
	}

	return nil
}
