// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	"github.com/vmware/terraform-provider-nsxt/api/orgs/projects/transit_gateways"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/metadata"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var staticRoutesParentPathExample = "/orgs/[org]/projects/[project]/transit-gateways/[transit-gateway]"

var transitGatewayStaticRouteSchema = map[string]*metadata.ExtendedSchema{
	"nsx_id":       metadata.GetExtendedSchema(getNsxIDSchema()),
	"path":         metadata.GetExtendedSchema(getPathSchema()),
	"display_name": metadata.GetExtendedSchema(getDisplayNameSchema()),
	"description":  metadata.GetExtendedSchema(getDescriptionSchema()),
	"revision":     metadata.GetExtendedSchema(getRevisionSchema()),
	"tag":          metadata.GetExtendedSchema(getTagsSchema()),
	"parent_path":  metadata.GetExtendedSchema(getPolicyPathSchema(true, true, "Policy path of the parent")),
	"network": {
		Schema: schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validateCidrOrIPOrRange(),
			Optional:     true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "string",
			SdkFieldName: "Network",
		},
	},
	"next_hop": {
		Schema: schema.Schema{
			Type:     schema.TypeList,
			MaxItems: 1,
			Elem: &metadata.ExtendedResource{
				Schema: map[string]*metadata.ExtendedSchema{
					"admin_distance": {
						Schema: schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1,
						},
						Metadata: metadata.Metadata{
							SchemaType:   "int",
							SdkFieldName: "AdminDistance",
						},
					},
					"scope": {
						Schema: schema.Schema{
							Type: schema.TypeList,
							Elem: &metadata.ExtendedSchema{
								Schema: schema.Schema{
									Type: schema.TypeString,
								},
								Metadata: metadata.Metadata{
									SchemaType: "string",
								},
							},
							Required: true,
						},
						Metadata: metadata.Metadata{
							SchemaType:   "list",
							SdkFieldName: "Scope",
						},
					},
				},
			},
			Optional: true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "list",
			SdkFieldName: "NextHops",
			ReflectType:  reflect.TypeOf(model.RouterNexthop{}),
		},
	},
	"enabled_on_secondary": {
		Schema: schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "bool",
			SdkFieldName: "EnabledOnSecondary",
		},
	},
}

func resourceNsxtPolicyTransitGatewayStaticRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyTransitGatewayStaticRouteCreate,
		Read:   resourceNsxtPolicyTransitGatewayStaticRouteRead,
		Update: resourceNsxtPolicyTransitGatewayStaticRouteUpdate,
		Delete: resourceNsxtPolicyTransitGatewayStaticRouteDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtParentPathResourceImporter,
		},
		Schema: metadata.GetSchemaFromExtendedSchema(transitGatewayStaticRouteSchema),
	}
}

func resourceNsxtPolicyTransitGatewayStaticRouteExists(sessionContext utl.SessionContext, parentPath string, id string, connector client.Connector) (bool, error) {
	var err error
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, staticRoutesParentPathExample)
	if pathErr != nil {
		return false, pathErr
	}
	client := transitgateways.NewStaticRoutesClient(sessionContext, connector)
	if client == nil {
		return false, fmt.Errorf("unsupported client type")
	}
	_, err = client.Get(parents[0], parents[1], parents[2], id)
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving resource", err)
}

func resourceNsxtPolicyTransitGatewayStaticRouteCreate(d *schema.ResourceData, m interface{}) error {
	if !util.NsxVersionHigherOrEqual("9.1.0") {
		return fmt.Errorf("Policy Transit Gateway Static Route resource requires NSX version 9.1.0 or higher")
	}
	connector := getPolicyConnector(m)

	id, err := getOrGenerateIDWithParent(d, m, resourceNsxtPolicyTransitGatewayStaticRouteExists)
	if err != nil {
		return err
	}

	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, staticRoutesParentPathExample)
	if pathErr != nil {
		return pathErr
	}
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)

	obj := model.StaticRoutes{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, transitGatewayStaticRouteSchema, "", nil); err != nil {
		return err
	}

	log.Printf("[INFO] Creating TGW StaticRoute with ID %s", id)

	sessionContext := getSessionContext(d, m)
	client := transitgateways.NewStaticRoutesClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}
	_, err = client.Patch(parents[0], parents[1], parents[2], id, obj)
	if err != nil {
		return handleCreateError("TGW StaticRoute", id, err)
	}
	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyTransitGatewayStaticRouteRead(d, m)
}

func resourceNsxtPolicyTransitGatewayStaticRouteRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining TGW StaticRoute ID")
	}

	sessionContext := getSessionContext(d, m)
	client := transitgateways.NewStaticRoutesClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}
	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, staticRoutesParentPathExample)
	if pathErr != nil {
		return pathErr
	}
	obj, err := client.Get(parents[0], parents[1], parents[2], id)
	if err != nil {
		return handleReadError(d, "TGW StaticRoute", id, err)
	}

	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("revision", obj.Revision)
	d.Set("path", obj.Path)

	elem := reflect.ValueOf(&obj).Elem()
	return metadata.StructToSchema(elem, d, transitGatewayStaticRouteSchema, "", nil)
}

func resourceNsxtPolicyTransitGatewayStaticRouteUpdate(d *schema.ResourceData, m interface{}) error {

	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining TGW StaticRoute ID")
	}

	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, staticRoutesParentPathExample)
	if pathErr != nil {
		return pathErr
	}
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getPolicyTagsFromSchema(d)

	revision := int64(d.Get("revision").(int))

	obj := model.StaticRoutes{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
		Revision:    &revision,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, transitGatewayStaticRouteSchema, "", nil); err != nil {
		return err
	}
	sessionContext := getSessionContext(d, m)
	client := transitgateways.NewStaticRoutesClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}
	_, err := client.Update(parents[0], parents[1], parents[2], id, obj)
	if err != nil {
		return handleUpdateError("TGW StaticRoute", id, err)
	}

	return resourceNsxtPolicyTransitGatewayStaticRouteRead(d, m)
}

func resourceNsxtPolicyTransitGatewayStaticRouteDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining TGW StaticRoute ID")
	}

	connector := getPolicyConnector(m)
	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, staticRoutesParentPathExample)
	if pathErr != nil {
		return pathErr
	}

	sessionContext := getSessionContext(d, m)
	client := transitgateways.NewStaticRoutesClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}
	err := client.Delete(parents[0], parents[1], parents[2], id)

	if err != nil {
		return handleDeleteError("TGW StaticRoute", id, err)
	}

	return nil
}
