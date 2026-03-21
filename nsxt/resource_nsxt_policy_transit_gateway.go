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
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxt "github.com/vmware/vsphere-automation-sdk-go/services/nsxt"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/orgs/projects"
	transitgateways "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/orgs/projects/transit_gateways"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/metadata"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var cliTransitGatewaysClient = projects.NewTransitGatewaysClient
var cliTransitGatewayCentralizedConfigsClient = transitgateways.NewCentralizedConfigsClient

const centralizedConfigID = "default"

var transitGatewaySchema = map[string]*metadata.ExtendedSchema{
	"nsx_id":       metadata.GetExtendedSchema(getNsxIDSchema()),
	"path":         metadata.GetExtendedSchema(getPathSchema()),
	"display_name": metadata.GetExtendedSchema(getDisplayNameSchema()),
	"description":  metadata.GetExtendedSchema(getDescriptionSchema()),
	"revision":     metadata.GetExtendedSchema(getRevisionSchema()),
	"tag":          metadata.GetExtendedSchema(getTagsSchema()),
	"context":      metadata.GetExtendedSchema(getContextSchemaExtended(true, false, false, true)),
	"transit_subnets": {
		Schema: schema.Schema{
			Type: schema.TypeList,
			Elem: &metadata.ExtendedSchema{
				Schema: schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateIPCidr(),
				},
				Metadata: metadata.Metadata{
					SchemaType: "string",
				},
			},
			Optional: true,
			Computed: true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "list",
			SdkFieldName: "TransitSubnets",
		},
	},
	"centralized_config": {
		Schema: schema.Schema{
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"ha_mode": {
						Type:         schema.TypeString,
						ValidateFunc: validation.StringInSlice([]string{model.CentralizedConfig_HA_MODE_ACTIVE, model.CentralizedConfig_HA_MODE_STANDBY}, false),
						Optional:     true,
						Default:      model.CentralizedConfig_HA_MODE_ACTIVE,
					},
					"edge_cluster_paths": {
						Type:     schema.TypeList,
						Optional: true,
						Elem: &schema.Schema{
							Type:         schema.TypeString,
							ValidateFunc: validatePolicyPath(),
						},
					},
				},
			},
		},
		Metadata: metadata.Metadata{Skip: true},
	},
	"span": {
		Schema: schema.Schema{
			Type:     schema.TypeList,
			MaxItems: 1,
			Computed: true,
			Elem: &metadata.ExtendedResource{
				Schema: map[string]*metadata.ExtendedSchema{
					"cluster_based_span": {
						Schema: schema.Schema{
							Type:     schema.TypeList,
							MaxItems: 1,
							Computed: true,
							ExactlyOneOf: []string{
								"span.0.cluster_based_span",
								"span.0.zone_based_span",
							},
							Optional: true,
							Elem: &metadata.ExtendedResource{
								Schema: map[string]*metadata.ExtendedSchema{
									"span_path": {
										Schema: schema.Schema{
											Type:         schema.TypeString,
											Optional:     true,
											Computed:     true,
											ValidateFunc: validatePolicyPath(),
										},
										Metadata: metadata.Metadata{
											Skip: true,
										},
									},
								},
							},
						},
						Metadata: metadata.Metadata{
							Skip: true,
						},
					},
					"zone_based_span": {
						Schema: schema.Schema{
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &metadata.ExtendedResource{
								Schema: map[string]*metadata.ExtendedSchema{
									"zone_external_ids": {
										Schema: schema.Schema{
											Type:     schema.TypeList,
											Required: true,
											MinItems: 0,
											MaxItems: 10,
											Elem: &metadata.ExtendedSchema{
												Schema: schema.Schema{
													Type: schema.TypeString,
												},
												Metadata: metadata.Metadata{
													Skip: true,
												},
											},
										},
										Metadata: metadata.Metadata{
											Skip: true,
										},
									},
								},
							},
						},
						Metadata: metadata.Metadata{
							Skip: true,
						},
					},
				},
			},
			Optional: true,
		},
		Metadata: metadata.Metadata{
			Skip: true,
		},
	},
}

func resourceNsxtPolicyTransitGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyTransitGatewayCreate,
		Read:   resourceNsxtPolicyTransitGatewayRead,
		Update: resourceNsxtPolicyTransitGatewayUpdate,
		Delete: resourceNsxtPolicyTransitGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: getPolicyPathResourceImporter(transitGatewayPathExample),
		},
		Schema: metadata.GetSchemaFromExtendedSchema(transitGatewaySchema),
	}
}

func resourceNsxtPolicyTransitGatewayExists(sessionContext utl.SessionContext, id string, connector client.Connector) (bool, error) {
	var err error
	parents := getVpcParentsFromContext(sessionContext)
	client := cliTransitGatewaysClient(connector)
	_, err = client.Get(parents[0], parents[1], id)
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving resource", err)
}

func getSpanFromSchema(iSpan interface{}) (*data.StructValue, error) {
	if len(iSpan.([]interface{})) == 0 {
		return nil, nil
	}
	converter := bindings.NewTypeConverter()

	// We're limiting to one span of any kind in the schema
	span := iSpan.([]interface{})[0].(map[string]interface{})
	if len(span["cluster_based_span"].([]interface{})) > 0 {
		cbs := span["cluster_based_span"].([]interface{})[0].(map[string]interface{})
		spanPath := cbs["span_path"].(string)
		clusterBasedSpan := model.ClusterBasedSpan{
			SpanPath: &spanPath,
			Type_:    model.BaseSpan_TYPE_CLUSTERBASEDSPAN,
		}
		dataValue, errs := converter.ConvertToVapi(clusterBasedSpan, model.ClusterBasedSpanBindingType())
		if errs != nil {
			return nil, errs[0]
		}
		return dataValue.(*data.StructValue), nil
	}
	if len(span["zone_based_span"].([]interface{})) > 0 {
		zbs := span["zone_based_span"].([]interface{})[0].(map[string]interface{})
		zoneExternalIds := interfaceListToStringList(zbs["zone_external_ids"].([]interface{}))
		zoneBasedSpan := model.ZoneBasedSpan{
			ZoneExternalIds: zoneExternalIds,
			Type_:           model.BaseSpan_TYPE_ZONEBASEDSPAN,
		}
		dataValue, errs := converter.ConvertToVapi(zoneBasedSpan, model.ZoneBasedSpanBindingType())
		if errs != nil {
			return nil, errs[0]
		}
		return dataValue.(*data.StructValue), nil
	}

	return nil, nil
}

func getCentralizedConfigFromSchema(v interface{}) *model.CentralizedConfig {
	if v == nil {
		return nil
	}
	l, ok := v.([]interface{})
	if !ok || len(l) == 0 {
		return nil
	}
	m, ok := l[0].(map[string]interface{})
	if !ok {
		return nil
	}
	cfg := &model.CentralizedConfig{
		ResourceType: strPtr("CentralizedConfig"),
	}
	if ha, ok := m["ha_mode"].(string); ok && ha != "" {
		cfg.HaMode = &ha
	}
	if paths, ok := m["edge_cluster_paths"].([]interface{}); ok && len(paths) > 0 {
		var s []string
		for _, p := range paths {
			if str, ok := p.(string); ok {
				s = append(s, str)
			}
		}
		cfg.EdgeClusterPaths = s
	}
	return cfg
}

func setCentralizedConfigInSchema(cfg *model.CentralizedConfig) []interface{} {
	if cfg == nil {
		return nil
	}
	m := make(map[string]interface{})
	if cfg.HaMode != nil {
		m["ha_mode"] = *cfg.HaMode
	}
	m["edge_cluster_paths"] = cfg.EdgeClusterPaths

	return []interface{}{m}
}

// buildCentralizedConfigChildren builds H-API child structs for TransitGateway.Children.
// When config is present, returns one ChildCentralizedConfig. When markDelete is true (user removed block), returns child with MarkedForDelete.
func buildCentralizedConfigChildren(config *model.CentralizedConfig, markDelete bool) ([]*data.StructValue, error) {
	converter := bindings.NewTypeConverter()
	id := centralizedConfigID
	child := model.ChildCentralizedConfig{
		Id:           &id,
		ResourceType: "ChildCentralizedConfig",
	}
	if markDelete {
		boolTrue := true
		child.MarkedForDelete = &boolTrue
		child.CentralizedConfig = &model.CentralizedConfig{Id: &id}
	} else if config != nil {
		child.CentralizedConfig = config
		if child.CentralizedConfig.Id == nil {
			child.CentralizedConfig.Id = &id
		}
	} else {
		return nil, nil
	}
	dataValue, errs := converter.ConvertToVapi(child, model.ChildCentralizedConfigBindingType())
	if len(errs) > 0 {
		return nil, errs[0]
	}
	return []*data.StructValue{dataValue.(*data.StructValue)}, nil
}

// createChildOrgWithTransitGateway builds the H-API tree Org -> Project -> TransitGateway for OrgRoot.Patch.
// When markForDelete is true, tgw is ignored and a minimal TransitGateway with MarkedForDelete is used for delete.
func createChildOrgWithTransitGateway(orgID, projectID, tgwID string, tgw *model.TransitGateway, markForDelete bool) (*data.StructValue, error) {
	converter := bindings.NewTypeConverter()
	var childTGW model.ChildTransitGateway
	if markForDelete {
		boolTrue := true
		tgwMinimal := model.TransitGateway{Id: &tgwID, ResourceType: strPtr("TransitGateway")}
		childTGW = model.ChildTransitGateway{
			Id:              &tgwID,
			ResourceType:    model.ChildTransitGateway__TYPE_IDENTIFIER,
			MarkedForDelete: &boolTrue,
			TransitGateway:  &tgwMinimal,
		}
	} else {
		tgw.Id = &tgwID
		childTGW = model.ChildTransitGateway{
			Id:             &tgwID,
			ResourceType:   model.ChildTransitGateway__TYPE_IDENTIFIER,
			TransitGateway: tgw,
		}
	}
	dataValue, errs := converter.ConvertToVapi(childTGW, model.ChildTransitGatewayBindingType())
	if len(errs) > 0 {
		return nil, errs[0]
	}
	targetProject := "Project"
	childProject := model.ChildResourceReference{
		Id:           &projectID,
		ResourceType: "ChildResourceReference",
		TargetType:   &targetProject,
		Children:     []*data.StructValue{dataValue.(*data.StructValue)},
	}
	dataValue, errs = converter.ConvertToVapi(childProject, model.ChildResourceReferenceBindingType())
	if len(errs) > 0 {
		return nil, errs[0]
	}
	targetOrg := "Org"
	childOrg := model.ChildResourceReference{
		Id:           &orgID,
		ResourceType: "ChildResourceReference",
		TargetType:   &targetOrg,
		Children:     []*data.StructValue{dataValue.(*data.StructValue)},
	}
	dataValue, errs = converter.ConvertToVapi(childOrg, model.ChildResourceReferenceBindingType())
	if len(errs) > 0 {
		return nil, errs[0]
	}
	return dataValue.(*data.StructValue), nil
}

func resourceNsxtPolicyTransitGatewayCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id, err := getOrGenerateID2(d, m, resourceNsxtPolicyTransitGatewayExists)
	if err != nil {
		return err
	}

	parents := getVpcParentsFromContext(getSessionContext(d, m))
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags, tagErr := getValidatedTagsFromSchema(d)
	if tagErr != nil {
		return tagErr
	}

	obj := model.TransitGateway{
		DisplayName:  &displayName,
		Description:  &description,
		ResourceType: strPtr("TransitGateway"),
		Tags:         tags,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, transitGatewaySchema, "", nil); err != nil {
		return err
	}

	if util.NsxVersionHigherOrEqual("9.1.0") {
		obj.Span, err = getSpanFromSchema(d.Get("span"))
		if err != nil {
			return handleCreateError("TransitGateway", id, err)
		}

	}

	// Attach centralized_config as H-API child in same transaction
	if cc := getCentralizedConfigFromSchema(d.Get("centralized_config")); cc != nil {
		children, err := buildCentralizedConfigChildren(cc, false)
		if err != nil {
			return handleCreateError("TransitGateway CentralizedConfig", id, err)
		}
		if children != nil {
			obj.Children = children
		}
	}

	orgRootClient := nsxt.NewOrgRootClient(connector)
	if orgRootClient == nil {
		return policyResourceNotSupportedError()
	}
	childOrg, err := createChildOrgWithTransitGateway(parents[0], parents[1], id, &obj, false)
	if err != nil {
		return handleCreateError("TransitGateway", id, err)
	}
	orgRoot := model.OrgRoot{
		ResourceType: strPtr("OrgRoot"),
		Children:     []*data.StructValue{childOrg},
	}
	log.Printf("[INFO] Creating TransitGateway with ID %s via OrgRoot H-API", id)
	if err := orgRootClient.Patch(orgRoot, nil); err != nil {
		return handleCreateError("TransitGateway", id, err)
	}
	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyTransitGatewayRead(d, m)
}

func resourceNsxtPolicyTransitGatewayRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining TransitGateway ID")
	}

	sessionContext := getSessionContext(d, m)
	client := cliTransitGatewaysClient(connector)
	parents := getVpcParentsFromContext(sessionContext)
	obj, err := client.Get(parents[0], parents[1], id)
	if err != nil {
		return handleReadError(d, "TransitGateway", id, err)
	}

	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("revision", obj.Revision)
	d.Set("path", obj.Path)

	if util.NsxVersionHigherOrEqual("9.1.0") && obj.Span != nil {
		span, err := setSpanFromSchema(obj.Span)

		if err != nil {
			return handleReadError(d, "TransitGateway", id, err)
		}
		d.Set("span", span)
	}

	// Read centralized config via H-API
	ccClient := cliTransitGatewayCentralizedConfigsClient(connector)
	if ccClient != nil {
		cc, err := ccClient.Get(parents[0], parents[1], id, centralizedConfigID)
		if err == nil {
			d.Set("centralized_config", setCentralizedConfigInSchema(&cc))
		} else if !isNotFoundError(err) {
			return handleReadError(d, "TransitGateway CentralizedConfig", id, err)
		} else {
			d.Set("centralized_config", nil)
		}
	}

	elem := reflect.ValueOf(&obj).Elem()
	return metadata.StructToSchema(elem, d, transitGatewaySchema, "", nil)
}

func setSpanFromSchema(span *data.StructValue) (interface{}, error) {
	converter := bindings.NewTypeConverter()
	base, errs := converter.ConvertToGolang(span, model.BaseSpanBindingType())

	schemaSpan := make(map[string]interface{})
	if errs != nil {
		return nil, errs[0]
	}
	switch base.(model.BaseSpan).Type_ {
	case model.BaseSpan_TYPE_CLUSTERBASEDSPAN:
		cbs, errs := converter.ConvertToGolang(span, model.ClusterBasedSpanBindingType())
		if errs != nil {
			return nil, errs[0]
		}
		clusterBasedSpan := cbs.(model.ClusterBasedSpan)

		elem := make(map[string]interface{})
		elem["span_path"] = clusterBasedSpan.SpanPath
		schemaSpan["cluster_based_span"] = []interface{}{elem}
		return []interface{}{schemaSpan}, nil

	case model.BaseSpan_TYPE_ZONEBASEDSPAN:
		zbs, errs := converter.ConvertToGolang(span, model.ZoneBasedSpanBindingType())
		if errs != nil {
			return nil, errs[0]
		}
		zoneBasedSpan := zbs.(model.ZoneBasedSpan)

		elem := make(map[string]interface{})
		elem["zone_external_ids"] = zoneBasedSpan.ZoneExternalIds
		schemaSpan["zone_based_span"] = []interface{}{elem}
		return []interface{}{schemaSpan}, nil
	}
	return nil, nil
}

func resourceNsxtPolicyTransitGatewayUpdate(d *schema.ResourceData, m interface{}) error {

	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining TransitGateway ID")
	}

	parents := getVpcParentsFromContext(getSessionContext(d, m))
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags, tagErr := getValidatedTagsFromSchema(d)
	if tagErr != nil {
		return tagErr
	}

	revision := int64(d.Get("revision").(int))

	obj := model.TransitGateway{
		DisplayName:  &displayName,
		Description:  &description,
		ResourceType: strPtr("TransitGateway"),
		Tags:         tags,
		Revision:     &revision,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, transitGatewaySchema, "", nil); err != nil {
		return err
	}

	var err error
	if util.NsxVersionHigherOrEqual("9.1.0") {
		obj.Span, err = getSpanFromSchema(d.Get("span"))
		if err != nil {
			return handleCreateError("TransitGateway", id, err)
		}
	}

	// When centralized_config changed, attach child for same-transaction update
	if d.HasChange("centralized_config") {
		newCC := getCentralizedConfigFromSchema(d.Get("centralized_config"))
		if newCC != nil {
			ccClient := cliTransitGatewayCentralizedConfigsClient(connector)
			if ccClient != nil {
				if existing, getErr := ccClient.Get(parents[0], parents[1], id, centralizedConfigID); getErr == nil {
					newCC.Revision = existing.Revision
				}
			}
			children, err := buildCentralizedConfigChildren(newCC, false)
			if err != nil {
				return handleUpdateError("TransitGateway CentralizedConfig", id, err)
			}
			if children != nil {
				obj.Children = children
			}
		} else {
			children, err := buildCentralizedConfigChildren(nil, true)
			if err != nil {
				return handleUpdateError("TransitGateway CentralizedConfig", id, err)
			}
			if children != nil {
				obj.Children = children
			}
		}
	}

	orgRootClient := nsxt.NewOrgRootClient(connector)
	if orgRootClient == nil {
		return policyResourceNotSupportedError()
	}
	childOrg, err := createChildOrgWithTransitGateway(parents[0], parents[1], id, &obj, false)
	if err != nil {
		return handleUpdateError("TransitGateway", id, err)
	}
	orgRoot := model.OrgRoot{
		ResourceType: strPtr("OrgRoot"),
		Children:     []*data.StructValue{childOrg},
	}
	log.Printf("[INFO] Updating TransitGateway with ID %s via OrgRoot H-API", id)
	if err := orgRootClient.Patch(orgRoot, nil); err != nil {
		return handleUpdateError("TransitGateway", id, err)
	}
	return resourceNsxtPolicyTransitGatewayRead(d, m)
}

func resourceNsxtPolicyTransitGatewayDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining TransitGateway ID")
	}

	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)

	orgRootClient := nsxt.NewOrgRootClient(connector)
	if orgRootClient == nil {
		return policyResourceNotSupportedError()
	}
	childOrg, err := createChildOrgWithTransitGateway(parents[0], parents[1], id, nil, true)
	if err != nil {
		return handleDeleteError("TransitGateway", id, err)
	}
	orgRoot := model.OrgRoot{
		ResourceType: strPtr("OrgRoot"),
		Children:     []*data.StructValue{childOrg},
	}
	log.Printf("[INFO] Deleting TransitGateway with ID %s via OrgRoot H-API", id)
	if err := orgRootClient.Patch(orgRoot, nil); err != nil {
		return handleDeleteError("TransitGateway", id, err)
	}
	return nil
}
