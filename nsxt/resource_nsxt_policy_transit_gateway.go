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
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	"github.com/vmware/terraform-provider-nsxt/api/orgs/projects"
	transitgateways "github.com/vmware/terraform-provider-nsxt/api/orgs/projects/transit_gateways"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/metadata"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var cliTransitGatewaysClient = projects.NewTransitGatewaysClient
var cliTransitGatewayCentralizedConfigsClient = transitgateways.NewCentralizedConfigsClient
var cliTransitGatewayRoutingConfigsClient = transitgateways.NewRoutingConfigClient
var cliTransitGatewayBgpClient = transitgateways.NewBgpClient

const centralizedConfigID = "default"
const routingConfigID = "default"

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
					"failover_mode": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(policyFailOverModeValues, false),
					},
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
	"advanced_config": {
		Schema: schema.Schema{
			Type:        schema.TypeList,
			Description: "Advanced configuration",
			Optional:    true,
			Computed:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"forwarding_up_timer": {
						Type:         schema.TypeInt,
						Optional:     true,
						Computed:     true,
						ValidateFunc: validation.IntBetween(0, 100),
					},
				},
			},
		},
		Metadata: metadata.Metadata{Skip: true},
	},
	"redistribution_config": {
		Schema: schema.Schema{
			Type:        schema.TypeList,
			Description: "Route redistribution configuration",
			Optional:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"rule": {
						Type:        schema.TypeList,
						Description: "List of route redistribution rules (1–5)",
						Required:    true,
						MinItems:    1,
						MaxItems:    5,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"types": {
									Type:        schema.TypeList,
									Description: "Route types to redistribute",
									Required:    true,
									MinItems:    1,
									Elem: &schema.Schema{
										Type: schema.TypeString,
										ValidateFunc: validation.StringInSlice([]string{
											model.TransitGatewayRedistributionRule_ROUTE_REDISTRIBUTION_TYPES_PUBLIC,
											model.TransitGatewayRedistributionRule_ROUTE_REDISTRIBUTION_TYPES_TGW_PRIVATE,
											model.TransitGatewayRedistributionRule_ROUTE_REDISTRIBUTION_TYPES_TGW_STATIC_ROUTE,
											model.TransitGatewayRedistributionRule_ROUTE_REDISTRIBUTION_TYPES_CONNECTED_SUBNET,
										}, false),
									},
								},
								"route_map_path": {
									Type:         schema.TypeString,
									Description:  "Policy path of a route map to apply for filtering routes",
									Optional:     true,
									ValidateFunc: validatePolicyPath(),
								},
							},
						},
					},
				},
			},
		},
		Metadata: metadata.Metadata{Skip: true},
	},
	"bgp_config": {
		Schema: schema.Schema{
			Type:        schema.TypeList,
			Description: "BGP routing configuration",
			Optional:    true,
			Computed:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: getPolicyBGPConfigSchema(),
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
	client := cliTransitGatewaysClient(sessionContext, connector)
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

func getCentralizedConfigFromSchema(d *schema.ResourceData, _ interface{}) (*model.CentralizedConfig, error) {
	v := d.Get("centralized_config")
	if v == nil {
		return nil, nil
	}
	l, ok := v.([]interface{})
	if !ok || len(l) == 0 {
		return nil, nil
	}
	m, ok := l[0].(map[string]interface{})
	if !ok {
		return nil, nil
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
	if fm, ok := m["failover_mode"].(string); ok && fm != "" {
		if util.NsxVersionLower("9.2.0") {
			return nil, fmt.Errorf("centralized_config.failover_mode is only supported with NSX 9.2.0 and above")
		}
		cfg.FailoverMode = &fm
	}
	return cfg, nil
}

func setCentralizedConfigInSchema(cfg *model.CentralizedConfig, _ interface{}) []interface{} {
	if cfg == nil {
		return nil
	}
	m := make(map[string]interface{})
	if cfg.HaMode != nil {
		m["ha_mode"] = *cfg.HaMode
	}
	m["edge_cluster_paths"] = cfg.EdgeClusterPaths
	if util.NsxVersionHigherOrEqual("9.2.0") && cfg.FailoverMode != nil {
		m["failover_mode"] = *cfg.FailoverMode
	}

	return []interface{}{m}
}

// buildTransitGatewayRoutingConfigChildren builds H-API child structs for TransitGatewayRoutingConfig.
func buildTransitGatewayRoutingConfigChildren(config *model.TransitGatewayRoutingConfig, markDelete bool) ([]*data.StructValue, error) {
	converter := bindings.NewTypeConverter()
	id := routingConfigID
	child := model.ChildTransitGatewayRoutingConfig{
		Id:           &id,
		ResourceType: model.ChildTransitGatewayRoutingConfig__TYPE_IDENTIFIER,
	}
	if markDelete {
		boolTrue := true
		child.MarkedForDelete = &boolTrue
		child.TransitGatewayRoutingConfig = &model.TransitGatewayRoutingConfig{Id: &id}
	} else if config != nil {
		child.TransitGatewayRoutingConfig = config
		if child.TransitGatewayRoutingConfig.Id == nil {
			child.TransitGatewayRoutingConfig.Id = &id
		}
	} else {
		return nil, nil
	}
	dataValue, errs := converter.ConvertToVapi(child, model.ChildTransitGatewayRoutingConfigBindingType())
	if len(errs) > 0 {
		return nil, errs[0]
	}
	return []*data.StructValue{dataValue.(*data.StructValue)}, nil
}

func getRedistributionConfigFromSchema(d *schema.ResourceData) *model.TransitGatewayRedistributionConfig {
	list := d.Get("redistribution_config").([]interface{})
	if len(list) == 0 {
		return nil
	}
	raw, ok := list[0].(map[string]interface{})
	if !ok || raw == nil {
		return nil
	}
	rulesRaw, ok := raw["rule"].([]interface{})
	if !ok || len(rulesRaw) == 0 {
		return nil
	}
	var rules []model.TransitGatewayRedistributionRule
	for _, r := range rulesRaw {
		ruleMap, ok := r.(map[string]interface{})
		if !ok {
			continue
		}
		rule := model.TransitGatewayRedistributionRule{
			RouteRedistributionTypes: interfaceListToStringList(ruleMap["types"].([]interface{})),
		}
		if routeMapPath, ok := ruleMap["route_map_path"].(string); ok && routeMapPath != "" {
			rule.RouteMapPath = &routeMapPath
		}
		rules = append(rules, rule)
	}
	return &model.TransitGatewayRedistributionConfig{Rules: rules}
}

func setRedistributionConfigInSchema(d *schema.ResourceData, config *model.TransitGatewayRedistributionConfig) error {
	if config == nil || len(config.Rules) == 0 {
		return d.Set("redistribution_config", nil)
	}
	var rules []interface{}
	for _, r := range config.Rules {
		ruleMap := map[string]interface{}{
			"types":          r.RouteRedistributionTypes,
			"route_map_path": "",
		}
		if r.RouteMapPath != nil {
			ruleMap["route_map_path"] = *r.RouteMapPath
		}
		rules = append(rules, ruleMap)
	}
	return d.Set("redistribution_config", []interface{}{
		map[string]interface{}{"rule": rules},
	})
}

func getTransitGatewayRoutingConfigFromSchema(d *schema.ResourceData) (*model.TransitGatewayRoutingConfig, error) {
	var forwardingUpTimer *int64
	list := d.Get("advanced_config").([]interface{})
	if len(list) > 0 {
		raw, ok := list[0].(map[string]interface{})
		if ok && raw != nil {
			if v, ok := raw["forwarding_up_timer"]; ok && v != nil {
				timer := int64(v.(int))
				forwardingUpTimer = &timer
			}
		}
	}
	redistributionConfig := getRedistributionConfigFromSchema(d)
	if forwardingUpTimer == nil && redistributionConfig == nil {
		return nil, nil
	}
	id := routingConfigID
	return &model.TransitGatewayRoutingConfig{
		Id:                   &id,
		ResourceType:         strPtr("TransitGatewayRoutingConfig"),
		ForwardingUpTimer:    forwardingUpTimer,
		RedistributionConfig: redistributionConfig,
	}, nil
}

// transitGatewayRoutingConfigForUpdate merges forwarding_up_timer from schema into the existing routing config from the API (preserves redistribution and revision).
func transitGatewayRoutingConfigForUpdate(d *schema.ResourceData, sessionContext utl.SessionContext, connector client.Connector, parents []string, tgwID string) (*model.TransitGatewayRoutingConfig, error) {
	rcClient := cliTransitGatewayRoutingConfigsClient(sessionContext, connector)
	if rcClient == nil {
		return nil, nil
	}
	fromSchema, err := getTransitGatewayRoutingConfigFromSchema(d)
	if err != nil {
		return nil, err
	}
	if fromSchema == nil {
		return nil, nil
	}
	existing, err := rcClient.Get(parents[0], parents[1], tgwID)
	if err != nil {
		if isNotFoundError(err) {
			return fromSchema, nil
		}
		return nil, err
	}
	if fromSchema.ForwardingUpTimer != nil {
		existing.ForwardingUpTimer = fromSchema.ForwardingUpTimer
	}
	existing.RedistributionConfig = fromSchema.RedistributionConfig
	return &existing, nil
}

func setTransitGatewayAdvancedConfigInSchema(d *schema.ResourceData, rc *model.TransitGatewayRoutingConfig) error {
	if rc == nil || rc.ForwardingUpTimer == nil {
		return d.Set("advanced_config", nil)
	}
	configMap := map[string]interface{}{
		"forwarding_up_timer": int(*rc.ForwardingUpTimer),
	}
	return d.Set("advanced_config", []interface{}{configMap})
}

// buildTGWBgpConfigChildren wraps a BgpRoutingConfig (or delete marker) as a ChildBgpRoutingConfig
// for inclusion in TransitGateway.Children during H-API calls.
// The non-delete path delegates to initPolicyTier0ChildBgpConfig (shared with Tier-0 gateway).
func buildTGWBgpConfigChildren(config *model.BgpRoutingConfig, markDelete bool) ([]*data.StructValue, error) {
	if markDelete {
		converter := bindings.NewTypeConverter()
		id := "bgp"
		boolTrue := true
		child := model.ChildBgpRoutingConfig{
			Id:               &id,
			ResourceType:     "ChildBgpRoutingConfig",
			MarkedForDelete:  &boolTrue,
			BgpRoutingConfig: &model.BgpRoutingConfig{Id: &id},
		}
		dataValue, errs := converter.ConvertToVapi(child, model.ChildBgpRoutingConfigBindingType())
		if len(errs) > 0 {
			return nil, errs[0]
		}
		return []*data.StructValue{dataValue.(*data.StructValue)}, nil
	}
	if config == nil {
		return nil, nil
	}
	dataValue, err := initPolicyTier0ChildBgpConfig(config)
	if err != nil {
		return nil, err
	}
	return []*data.StructValue{dataValue}, nil
}

func getTGWBgpConfigFromSchema(d *schema.ResourceData) *model.BgpRoutingConfig {
	bgpConfig := d.Get("bgp_config").([]interface{})
	if len(bgpConfig) == 0 || bgpConfig[0] == nil {
		return nil
	}
	cfgMap := bgpConfig[0].(map[string]interface{})
	cfg := resourceNsxtPolicyTier0GatewayBGPConfigSchemaToStruct(cfgMap, false, d.Id())
	return &cfg
}

func setTGWBgpConfigInSchema(d *schema.ResourceData, bgpConfig *model.BgpRoutingConfig) error {
	if bgpConfig == nil {
		return d.Set("bgp_config", make([]map[string]interface{}, 0))
	}
	cfgMap := initPolicyTier0BGPConfigMap(bgpConfig)
	return d.Set("bgp_config", []interface{}{cfgMap})
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

	// Attach centralized_config and routing_config as H-API children in same transaction
	var childStructs []*data.StructValue
	cc, err := getCentralizedConfigFromSchema(d, m)
	if err != nil {
		return err
	}
	if cc != nil {
		ch, err := buildCentralizedConfigChildren(cc, false)
		if err != nil {
			return handleCreateError("TransitGateway CentralizedConfig", id, err)
		}
		childStructs = append(childStructs, ch...)
	}
	rc, err := getTransitGatewayRoutingConfigFromSchema(d)
	if err != nil {
		return handleCreateError("TransitGateway RoutingConfig", id, err)
	}
	if rc != nil {
		ch, err := buildTransitGatewayRoutingConfigChildren(rc, false)
		if err != nil {
			return handleCreateError("TransitGateway RoutingConfig", id, err)
		}
		childStructs = append(childStructs, ch...)
	}
	if bgpCfg := getTGWBgpConfigFromSchema(d); bgpCfg != nil {
		ch, err := buildTGWBgpConfigChildren(bgpCfg, false)
		if err != nil {
			return handleCreateError("TransitGateway BgpConfig", id, err)
		}
		childStructs = append(childStructs, ch...)
	}
	if len(childStructs) > 0 {
		obj.Children = childStructs
	}

	sessionContext := getSessionContext(d, m)
	orgRootClient := cliOrgRootClient(sessionContext, connector)
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
	client := cliTransitGatewaysClient(sessionContext, connector)
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
	ccClient := cliTransitGatewayCentralizedConfigsClient(sessionContext, connector)
	if ccClient != nil {
		cc, err := ccClient.Get(parents[0], parents[1], id, centralizedConfigID)
		if err == nil {
			d.Set("centralized_config", setCentralizedConfigInSchema(&cc, m))
		} else if !isNotFoundError(err) {
			return handleReadError(d, "TransitGateway CentralizedConfig", id, err)
		} else {
			d.Set("centralized_config", nil)
		}
	}

	rcClient := cliTransitGatewayRoutingConfigsClient(sessionContext, connector)
	if rcClient != nil {
		rc, err := rcClient.Get(parents[0], parents[1], id)
		if err == nil {
			if err := setTransitGatewayAdvancedConfigInSchema(d, &rc); err != nil {
				return handleReadError(d, "TransitGateway RoutingConfig", id, err)
			}
			if err := setRedistributionConfigInSchema(d, rc.RedistributionConfig); err != nil {
				return handleReadError(d, "TransitGateway RoutingConfig", id, err)
			}
		} else if !isNotFoundError(err) {
			return handleReadError(d, "TransitGateway RoutingConfig", id, err)
		} else {
			if err := d.Set("advanced_config", nil); err != nil {
				return handleReadError(d, "TransitGateway RoutingConfig", id, err)
			}
			if err := d.Set("redistribution_config", nil); err != nil {
				return handleReadError(d, "TransitGateway RoutingConfig", id, err)
			}
		}
	}

	bgpClient := cliTransitGatewayBgpClient(sessionContext, connector)
	if bgpClient != nil {
		bgpConfig, err := bgpClient.Get(parents[0], parents[1], id)
		if err == nil {
			if err := setTGWBgpConfigInSchema(d, &bgpConfig); err != nil {
				return handleReadError(d, "TransitGateway BgpConfig", id, err)
			}
		} else if !isNotFoundError(err) {
			return handleReadError(d, "TransitGateway BgpConfig", id, err)
		} else {
			if err := setTGWBgpConfigInSchema(d, nil); err != nil {
				return handleReadError(d, "TransitGateway BgpConfig", id, err)
			}
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

	// When centralized_config or advanced_config changed, attach H-API children for same-transaction update
	sessionContext := getSessionContext(d, m)
	var childStructs []*data.StructValue

	if d.HasChange("centralized_config") {
		newCC, ccErr := getCentralizedConfigFromSchema(d, m)
		if ccErr != nil {
			return handleUpdateError("TransitGateway CentralizedConfig", id, ccErr)
		}
		if newCC != nil {
			ccClient := cliTransitGatewayCentralizedConfigsClient(sessionContext, connector)
			if ccClient != nil {
				if existing, getErr := ccClient.Get(parents[0], parents[1], id, centralizedConfigID); getErr == nil {
					newCC.Revision = existing.Revision
				}
			}
			ch, err := buildCentralizedConfigChildren(newCC, false)
			if err != nil {
				return handleUpdateError("TransitGateway CentralizedConfig", id, err)
			}
			childStructs = append(childStructs, ch...)
		} else {
			ch, err := buildCentralizedConfigChildren(nil, true)
			if err != nil {
				return handleUpdateError("TransitGateway CentralizedConfig", id, err)
			}
			childStructs = append(childStructs, ch...)
		}
	}

	if d.HasChange("advanced_config") || d.HasChange("redistribution_config") {
		rc, rcErr := transitGatewayRoutingConfigForUpdate(d, sessionContext, connector, parents, id)
		if rcErr != nil {
			return handleUpdateError("TransitGateway RoutingConfig", id, rcErr)
		}
		if rc != nil {
			ch, err := buildTransitGatewayRoutingConfigChildren(rc, false)
			if err != nil {
				return handleUpdateError("TransitGateway RoutingConfig", id, err)
			}
			childStructs = append(childStructs, ch...)
		}
	}

	if d.HasChange("bgp_config") {
		newBgp := getTGWBgpConfigFromSchema(d)
		if newBgp != nil {
			bgpAPIClient := cliTransitGatewayBgpClient(sessionContext, connector)
			if bgpAPIClient != nil {
				if existing, getErr := bgpAPIClient.Get(parents[0], parents[1], id); getErr == nil {
					newBgp.Revision = existing.Revision
				}
			}
			ch, err := buildTGWBgpConfigChildren(newBgp, false)
			if err != nil {
				return handleUpdateError("TransitGateway BgpConfig", id, err)
			}
			childStructs = append(childStructs, ch...)
		} else {
			ch, err := buildTGWBgpConfigChildren(nil, true)
			if err != nil {
				return handleUpdateError("TransitGateway BgpConfig", id, err)
			}
			childStructs = append(childStructs, ch...)
		}
	}

	if len(childStructs) > 0 {
		obj.Children = childStructs
	}

	orgRootClient := cliOrgRootClient(sessionContext, connector)
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

	orgRootClient := cliOrgRootClient(sessionContext, connector)
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
