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
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	"github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/metadata"
)

var cliDistributedVxlanConnectionsClient = infra.NewDistributedVxlanConnectionsClient

var distributedVxlanConnectionSchema = map[string]*metadata.ExtendedSchema{
	"nsx_id":       metadata.GetExtendedSchema(getNsxIDSchema()),
	"path":         metadata.GetExtendedSchema(getPathSchema()),
	"display_name": metadata.GetExtendedSchema(getDisplayNameSchema()),
	"description":  metadata.GetExtendedSchema(getDescriptionSchema()),
	"revision":     metadata.GetExtendedSchema(getRevisionSchema()),
	"tag":          metadata.GetExtendedSchema(getTagsSchema()),
	"connectivity_type": {
		Schema: schema.Schema{
			Type:         schema.TypeString,
			Optional:     true,
			Default:      model.DistributedVxlanConnection_CONNECTIVITY_TYPE_EVPN,
			Description:  "Connectivity type for the distributed VXLAN connection. Currently only L3_EVPN is supported.",
			ValidateFunc: validation.StringInSlice([]string{model.DistributedVxlanConnection_CONNECTIVITY_TYPE_EVPN}, false),
		},
		Metadata: metadata.Metadata{
			SchemaType:   "string",
			SdkFieldName: "ConnectivityType",
		},
	},
	"l3_vni": {
		Schema: schema.Schema{
			Type:        schema.TypeInt,
			Required:    true,
			Description: "L3 VNI for overlay traffic.",
		},
		Metadata: metadata.Metadata{
			SchemaType:   "int",
			SdkFieldName: "L3Vni",
		},
	},
	"route_controller_path": {
		Schema: schema.Schema{
			Type:         schema.TypeString,
			Required:     true,
			Description:  "Policy path of the route controller.",
			ValidateFunc: validatePolicyPath(),
		},
		Metadata: metadata.Metadata{
			SchemaType:   "string",
			SdkFieldName: "RouteControllerPath",
		},
	},
	"route_distinguisher": {
		Schema: schema.Schema{
			Type:         schema.TypeString,
			Required:     true,
			Description:  "Route distinguisher for the VXLAN connection.",
			ValidateFunc: validateIPorASNPair,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "string",
			SdkFieldName: "RouteDistinguisher",
		},
	},
}

const distributedVxlanConnectionPathExample = "/infra/distributed-vxlan-connections/[connection]"

func getDistributedVxlanConnectionRouteTargetSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Required:    true,
		MinItems:    1,
		MaxItems:    1,
		Description: "Route targets for the distributed VXLAN connection.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"address_family": {
					Type:         schema.TypeString,
					Optional:     true,
					Default:      model.VrfRouteTargets_ADDRESS_FAMILY_EVPN,
					Description:  "Address family for route targets. Currently only L2VPN_EVPN is supported.",
					ValidateFunc: validation.StringInSlice([]string{model.VrfRouteTargets_ADDRESS_FAMILY_EVPN}, false),
				},
				"import_targets": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "List of import route targets.",
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: validateASNPair,
					},
				},
				"export_targets": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "List of export route targets.",
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: validateASNPair,
					},
				},
			},
		},
	}
}

func resourceNsxtPolicyDistributedVxlanConnection() *schema.Resource {
	s := metadata.GetSchemaFromExtendedSchema(distributedVxlanConnectionSchema)
	s["route_target"] = getDistributedVxlanConnectionRouteTargetSchema()
	return &schema.Resource{
		Create: resourceNsxtPolicyDistributedVxlanConnectionCreate,
		Read:   resourceNsxtPolicyDistributedVxlanConnectionRead,
		Update: resourceNsxtPolicyDistributedVxlanConnectionUpdate,
		Delete: resourceNsxtPolicyDistributedVxlanConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: getPolicyPathResourceImporter(distributedVxlanConnectionPathExample),
		},
		Schema: s,
	}
}

func resourceNsxtPolicyDistributedVxlanConnectionExists(id string, connector client.Connector, isGlobalManager bool) (bool, error) {
	var err error

	sessionContext := utl.SessionContext{ClientType: utl.Local}
	client := cliDistributedVxlanConnectionsClient(sessionContext, connector)
	if client == nil {
		return false, fmt.Errorf("unsupported client type")
	}
	_, err = client.Get(id)
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving resource", err)
}

func distributedVxlanConnectionRouteTargetsFromSchema(d *schema.ResourceData) []model.VrfRouteTargets {
	routeTargets := d.Get("route_target").([]interface{})
	if len(routeTargets) == 0 {
		return nil
	}
	rt := routeTargets[0].(map[string]interface{})
	addressFamily := rt["address_family"].(string)
	exportTargets := interface2StringList(rt["export_targets"].([]interface{}))
	importTargets := interface2StringList(rt["import_targets"].([]interface{}))
	targets := model.VrfRouteTargets{
		AddressFamily:      &addressFamily,
		ExportRouteTargets: exportTargets,
		ImportRouteTargets: importTargets,
	}
	return []model.VrfRouteTargets{targets}
}

func distributedVxlanConnectionRouteTargetsToSchema(d *schema.ResourceData, routeTargets []model.VrfRouteTargets) {
	if len(routeTargets) == 0 {
		return
	}
	rt := routeTargets[0]
	elem := make(map[string]interface{})
	elem["address_family"] = rt.AddressFamily
	elem["import_targets"] = rt.ImportRouteTargets
	elem["export_targets"] = rt.ExportRouteTargets
	d.Set("route_target", []interface{}{elem})
}

func resourceNsxtPolicyDistributedVxlanConnectionCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id, err := getOrGenerateID(d, m, resourceNsxtPolicyDistributedVxlanConnectionExists)
	if err != nil {
		return err
	}

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)

	obj := model.DistributedVxlanConnection{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, distributedVxlanConnectionSchema, "", nil); err != nil {
		return err
	}

	obj.RouteTargets = distributedVxlanConnectionRouteTargetsFromSchema(d)

	log.Printf("[INFO] Creating DistributedVxlanConnection with ID %s", id)

	sessionContext := getSessionContext(d, m)
	client := cliDistributedVxlanConnectionsClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}
	err = client.Patch(id, obj)
	if err != nil {
		return handleCreateError("DistributedVxlanConnection", id, err)
	}
	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyDistributedVxlanConnectionRead(d, m)
}

func resourceNsxtPolicyDistributedVxlanConnectionRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining DistributedVxlanConnection ID")
	}

	sessionContext := getSessionContext(d, m)
	client := cliDistributedVxlanConnectionsClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}

	obj, err := client.Get(id)
	if err != nil {
		return handleReadError(d, "DistributedVxlanConnection", id, err)
	}

	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("revision", obj.Revision)
	d.Set("path", obj.Path)

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.StructToSchema(elem, d, distributedVxlanConnectionSchema, "", nil); err != nil {
		return err
	}

	distributedVxlanConnectionRouteTargetsToSchema(d, obj.RouteTargets)

	return nil
}

func resourceNsxtPolicyDistributedVxlanConnectionUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining DistributedVxlanConnection ID")
	}

	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getPolicyTagsFromSchema(d)
	revision := int64(d.Get("revision").(int))

	obj := model.DistributedVxlanConnection{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
		Revision:    &revision,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, distributedVxlanConnectionSchema, "", nil); err != nil {
		return err
	}

	obj.RouteTargets = distributedVxlanConnectionRouteTargetsFromSchema(d)

	sessionContext := getSessionContext(d, m)
	client := cliDistributedVxlanConnectionsClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}
	_, err := client.Update(id, obj)
	if err != nil {
		return handleUpdateError("DistributedVxlanConnection", id, err)
	}

	return resourceNsxtPolicyDistributedVxlanConnectionRead(d, m)
}

func resourceNsxtPolicyDistributedVxlanConnectionDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining DistributedVxlanConnection ID")
	}

	connector := getPolicyConnector(m)

	sessionContext := getSessionContext(d, m)
	client := cliDistributedVxlanConnectionsClient(sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}
	err := client.Delete(id)

	if err != nil {
		return handleDeleteError("DistributedVxlanConnection", id, err)
	}

	return nil
}
