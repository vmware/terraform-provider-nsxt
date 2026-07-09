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
	tgwrouting "github.com/vmware/terraform-provider-nsxt/api/orgs/projects/transit_gateways/routing"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/metadata"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var cliTransitGatewaysClient = projects.NewTransitGatewaysClient
var cliTransitGatewayCentralizedConfigsClient = transitgateways.NewCentralizedConfigsClient
var cliTransitGatewayBgpClient = tgwrouting.NewBgpClient

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
					"failover_mode": {
						Type:         schema.TypeString,
						Optional:     true,
						Computed:     true,
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
											model.TransitGatewayBgpRedistributionRule_ROUTE_REDISTRIBUTION_TYPES_PUBLIC,
											model.TransitGatewayBgpRedistributionRule_ROUTE_REDISTRIBUTION_TYPES_TGW_PRIVATE,
											model.TransitGatewayBgpRedistributionRule_ROUTE_REDISTRIBUTION_TYPES_TGW_STATIC_ROUTE,
											model.TransitGatewayBgpRedistributionRule_ROUTE_REDISTRIBUTION_TYPES_CONNECTED_SUBNET,
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
							Computed: true,
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
	if iSpan == nil {
		return nil, nil
	}
	spanList, ok := iSpan.([]interface{})
	if !ok || len(spanList) == 0 || spanList[0] == nil {
		return nil, nil
	}
	// We're limiting to one span of any kind in the schema
	span, ok := spanList[0].(map[string]interface{})
	if !ok {
		return nil, nil
	}
	converter := bindings.NewTypeConverter()
	// Presence is determined by list length rather than getElemOrEmptyMapFromMap:
	// when zone_external_ids is explicitly set to an empty list, the SDKv2 diff
	// engine emits no attribute entries for the nested block's contents (only the
	// "#" count), so the reconstructed element can be a bare nil even though the
	// block itself is present. Checking len(list) > 0 matches the correctly
	// diffed "#" count and avoids treating a present-but-empty block as absent.
	if cbsList, ok := span["cluster_based_span"].([]interface{}); ok && len(cbsList) > 0 {
		cbs, _ := cbsList[0].(map[string]interface{})
		spanPath, _ := cbs["span_path"].(string)
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
	if zbsList, ok := span["zone_based_span"].([]interface{}); ok && len(zbsList) > 0 {
		var zoneExternalIds []string
		if zbs, ok := zbsList[0].(map[string]interface{}); ok {
			zoneExternalIds = interfaceListToStringList(zbs["zone_external_ids"].([]interface{}))
		}
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

func buildTGWBgpConfigChildren(config *model.TransitGatewayBgpRoutingConfig, markDelete bool) ([]*data.StructValue, error) {
	if markDelete {
		converter := bindings.NewTypeConverter()
		id := "bgp"
		boolTrue := true
		child := model.ChildTransitGatewayBgpRoutingConfig{
			Id:                             &id,
			ResourceType:                   "ChildTransitGatewayBgpRoutingConfig",
			MarkedForDelete:                &boolTrue,
			TransitGatewayBgpRoutingConfig: &model.TransitGatewayBgpRoutingConfig{Id: &id},
		}
		dataValue, errs := converter.ConvertToVapi(child, model.ChildTransitGatewayBgpRoutingConfigBindingType())
		if len(errs) > 0 {
			return nil, errs[0]
		}
		return []*data.StructValue{dataValue.(*data.StructValue)}, nil
	}
	if config == nil {
		return nil, nil
	}
	dataValue, err := initTGWChildBgpConfig(config)
	if err != nil {
		return nil, err
	}
	return []*data.StructValue{dataValue}, nil
}

func initTGWChildBgpConfig(config *model.TransitGatewayBgpRoutingConfig) (*data.StructValue, error) {
	converter := bindings.NewTypeConverter()
	childConfig := model.ChildTransitGatewayBgpRoutingConfig{
		ResourceType:                   "ChildTransitGatewayBgpRoutingConfig",
		TransitGatewayBgpRoutingConfig: config,
	}
	dataValue, errors := converter.ConvertToVapi(childConfig, model.ChildTransitGatewayBgpRoutingConfigBindingType())
	if errors != nil {
		return nil, fmt.Errorf("Error converting child TGW BGP Routing Configuration: %v", errors[0])
	}
	return dataValue.(*data.StructValue), nil
}

func getTGWBgpConfigFromSchema(d *schema.ResourceData) *model.TransitGatewayBgpRoutingConfig {
	bgpConfig := d.Get("bgp_config").([]interface{})
	if len(bgpConfig) == 0 || bgpConfig[0] == nil {
		return nil
	}
	cfgMap := bgpConfig[0].(map[string]interface{})

	revision := int64(cfgMap["revision"].(int))
	ecmp := cfgMap["ecmp"].(bool)
	enabled := cfgMap["enabled"].(bool)
	localAsNum := cfgMap["local_as_num"].(string)
	multipathRelax := cfgMap["multipath_relax"].(bool)
	restartMode := cfgMap["graceful_restart_mode"].(string)
	restartTimer := int64(cfgMap["graceful_restart_timer"].(int))
	staleTimer := int64(cfgMap["graceful_restart_stale_route_timer"].(int))

	var tagStructs []model.Tag
	if cfgMap["tag"] != nil {
		for _, tag := range cfgMap["tag"].(*schema.Set).List() {
			data := tag.(map[string]interface{})
			tagScope := data["scope"].(string)
			tagTag := data["tag"].(string)
			tagStructs = append(tagStructs, model.Tag{Scope: &tagScope, Tag: &tagTag})
		}
	}

	var aggregationStructs []model.RouteAggregationEntry
	for _, agg := range cfgMap["route_aggregation"].([]interface{}) {
		data := agg.(map[string]interface{})
		prefix := data["prefix"].(string)
		summary := data["summary_only"].(bool)
		aggregationStructs = append(aggregationStructs, model.RouteAggregationEntry{
			Prefix:      &prefix,
			SummaryOnly: &summary,
		})
	}

	restartTimerStruct := model.BgpGracefulRestartTimer{
		RestartTimer:    &restartTimer,
		StaleRouteTimer: &staleTimer,
	}
	restartConfigStruct := model.BgpGracefulRestartConfig{
		Mode:  &restartMode,
		Timer: &restartTimerStruct,
	}

	id := "bgp"
	resourceType := "TransitGatewayBgpRoutingConfig"
	cfg := &model.TransitGatewayBgpRoutingConfig{
		Ecmp:                  &ecmp,
		Enabled:               &enabled,
		RouteAggregations:     aggregationStructs,
		ResourceType:          &resourceType,
		Tags:                  tagStructs,
		Id:                    &id,
		Revision:              &revision,
		MultipathRelax:        &multipathRelax,
		GracefulRestartConfig: &restartConfigStruct,
	}
	if len(localAsNum) > 0 {
		cfg.LocalAsNum = &localAsNum
	}
	if interSrIbgp, ok := cfgMap["inter_sr_ibgp"].(bool); ok {
		cfg.InterSrIbgp = &interSrIbgp
	}
	return cfg
}

func setTGWBgpConfigInSchema(d *schema.ResourceData, bgpConfig *model.TransitGatewayBgpRoutingConfig) error {
	if bgpConfig == nil {
		return d.Set("bgp_config", make([]map[string]interface{}, 0))
	}
	cfgMap := make(map[string]interface{})
	cfgMap["revision"] = bgpConfig.Revision
	cfgMap["path"] = bgpConfig.Path
	cfgMap["ecmp"] = bgpConfig.Ecmp
	cfgMap["enabled"] = bgpConfig.Enabled
	cfgMap["inter_sr_ibgp"] = bgpConfig.InterSrIbgp
	cfgMap["local_as_num"] = bgpConfig.LocalAsNum
	cfgMap["multipath_relax"] = bgpConfig.MultipathRelax

	if bgpConfig.GracefulRestartConfig != nil {
		cfgMap["graceful_restart_mode"] = bgpConfig.GracefulRestartConfig.Mode
		if bgpConfig.GracefulRestartConfig.Timer != nil {
			cfgMap["graceful_restart_timer"] = int(*bgpConfig.GracefulRestartConfig.Timer.RestartTimer)
			cfgMap["graceful_restart_stale_route_timer"] = int(*bgpConfig.GracefulRestartConfig.Timer.StaleRouteTimer)
		}
	} else {
		cfgMap["graceful_restart_mode"] = model.BgpGracefulRestartConfig_MODE_HELPER_ONLY
		cfgMap["graceful_restart_timer"] = policyBGPGracefulRestartTimerDefault
		cfgMap["graceful_restart_stale_route_timer"] = policyBGPGracefulRestartStaleRouteTimerDefault
	}

	var tagList []map[string]string
	for _, tag := range bgpConfig.Tags {
		tagList = append(tagList, map[string]string{"scope": *tag.Scope, "tag": *tag.Tag})
	}
	cfgMap["tag"] = tagList

	var aggregationList []map[string]interface{}
	for _, agg := range bgpConfig.RouteAggregations {
		aggregationList = append(aggregationList, map[string]interface{}{
			"prefix":       agg.Prefix,
			"summary_only": *agg.SummaryOnly,
		})
	}
	cfgMap["route_aggregation"] = aggregationList

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
	var tags []model.Tag
	var tagErr error
	if isConfigScopedCacheMode() {
		tags = getPolicyTagsWithProviderManagedDefaults(d, m)
	} else {
		tags, tagErr = getValidatedTagsFromSchema(d)
		if tagErr != nil {
			return tagErr
		}
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

	// Attach centralized_config and bgp_config as H-API children in the same creation transaction.
	// Note: TransitGatewayRoutingConfig cannot be included as an H-API child during TGW creation
	// (API error 500165 "Invalid Dto type hierarchy"). It is patched separately below via direct API.
	// BgpRoutingConfig IS supported as an H-API child and is included here for atomic creation.
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
	MarkPostWriteAndInvalidateCacheForResourceType("TransitGateway", d.Id())

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
	var obj *model.TransitGateway
	var err error
	if isCacheEnabledForRead(d) {
		obj, _, _, err = CacheAwareResourceRead[model.TransitGateway](
			d,
			m,
			connector,
			id,
			"TransitGateway",
			model.TransitGatewayBindingType(),
			func() (*model.TransitGateway, error) {
				readObj, readErr := client.Get(parents[0], parents[1], id)
				if readErr != nil {
					return nil, readErr
				}
				return &readObj, nil
			},
			func(patchObj *model.TransitGateway) error {
				orgRootClient := cliOrgRootClient(sessionContext, connector)
				if orgRootClient == nil {
					return policyResourceNotSupportedError()
				}
				childOrg, childErr := createChildOrgWithTransitGateway(parents[0], parents[1], id, patchObj, false)
				if childErr != nil {
					return childErr
				}
				orgRoot := model.OrgRoot{
					ResourceType: strPtr("OrgRoot"),
					Children:     []*data.StructValue{childOrg},
				}
				return orgRootClient.Patch(orgRoot, nil)
			},
		)
	} else {
		readObj, readErr := client.Get(parents[0], parents[1], id)
		if readErr != nil {
			err = readErr
		} else {
			obj = &readObj
		}
	}
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

	_ = d.Set("advanced_config", nil)
	_ = d.Set("redistribution_config", nil)

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

	elem := reflect.ValueOf(obj).Elem()
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
	var tags []model.Tag
	var tagErr error
	if isConfigScopedCacheMode() {
		tags = getPolicyTagsWithProviderManagedDefaults(d, m)
	} else {
		tags, tagErr = getValidatedTagsFromSchema(d)
		if tagErr != nil {
			return tagErr
		}
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
			return handleUpdateError("TransitGateway", id, err)
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

	if d.HasChange("bgp_config") {
		newBgp := getTGWBgpConfigFromSchema(d)
		bgpAPIClient := cliTransitGatewayBgpClient(sessionContext, connector)
		if bgpAPIClient != nil {
			if newBgp != nil {
				if existing, getErr := bgpAPIClient.Get(parents[0], parents[1], id); getErr == nil {
					newBgp.Revision = existing.Revision
				}
				log.Printf("[INFO] Updating TransitGateway %s BGP config via direct API", id)
				if err := bgpAPIClient.Patch(parents[0], parents[1], id, *newBgp); err != nil {
					return handleUpdateError("TransitGateway BgpConfig", id, err)
				}
			}
		}
	}

	MarkPostWriteAndInvalidateCacheForResourceType("TransitGateway", d.Id())
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
