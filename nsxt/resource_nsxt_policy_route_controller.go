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

	"github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/metadata"
)

var cliRouteControllersClient = infra.NewRouteControllersClient
var cliRouteControllerBgpClient = infra.NewRouteControllerBgpClient

var routeControllerBgpGracefulRestartModes = []string{
	model.BgpGracefulRestartConfig_MODE_DISABLE,
	model.BgpGracefulRestartConfig_MODE_GR_AND_HELPER,
	model.BgpGracefulRestartConfig_MODE_HELPER_ONLY,
}

var routeControllerSchema = map[string]*metadata.ExtendedSchema{
	"nsx_id":       metadata.GetExtendedSchema(getNsxIDSchema()),
	"path":         metadata.GetExtendedSchema(getPathSchema()),
	"display_name": metadata.GetExtendedSchema(getDisplayNameSchema()),
	"description":  metadata.GetExtendedSchema(getDescriptionSchema()),
	"revision":     metadata.GetExtendedSchema(getRevisionSchema()),
	"tag":          metadata.GetExtendedSchema(getTagsSchema()),
	"ha_mode": {
		Schema: schema.Schema{
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			Description:  "High-availability mode for route controller. Currently only ACTIVE_STANDBY is supported.",
			ValidateFunc: validation.StringInSlice([]string{model.RouteController_HA_MODE_STANDBY}, false),
		},
		Metadata: metadata.Metadata{
			SchemaType:   "string",
			SdkFieldName: "HaMode",
			OmitIfEmpty:  true,
		},
	},
	"virtual_network_appliance_cluster_path": {
		Schema: schema.Schema{
			Type:         schema.TypeString,
			Optional:     true,
			Description:  "Policy path for the virtual network appliance cluster.",
			ValidateFunc: validatePolicyPath(),
		},
		Metadata: metadata.Metadata{
			SchemaType:   "string",
			SdkFieldName: "VirtualNetworkApplianceClusterPath",
			OmitIfEmpty:  true,
		},
	},
}

const routeControllerPathExample = "/infra/route-controllers/[controller]"

func getRouteControllerBgpConfigSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "BGP routing configuration for the route controller",
		Optional:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"path":     getPathSchema(),
				"revision": getRevisionSchema(),
				"ecmp": {
					Type:        schema.TypeBool,
					Description: "Flag to enable ECMP",
					Optional:    true,
					Default:     true,
				},
				"local_as_num": {
					Type:         schema.TypeString,
					Description:  "BGP AS number in ASPLAIN/ASDOT format",
					Required:     true,
					ValidateFunc: validateASPlainOrDot,
				},
				"multipath_relax": {
					Type:        schema.TypeBool,
					Description: "Flag to enable BGP multipath relax option",
					Optional:    true,
					Computed:    true,
				},
				"graceful_restart_mode": {
					Type:         schema.TypeString,
					Description:  "BGP graceful restart configuration mode",
					ValidateFunc: validation.StringInSlice(routeControllerBgpGracefulRestartModes, false),
					Optional:     true,
					Default:      model.BgpGracefulRestartConfig_MODE_HELPER_ONLY,
				},
				"graceful_restart_timer": {
					Type:         schema.TypeInt,
					Description:  "BGP graceful restart timer in seconds",
					Optional:     true,
					Default:      policyBGPGracefulRestartTimerDefault,
					ValidateFunc: validation.IntBetween(1, 3600),
				},
				"peer_route_convergence_timer": {
					Type:         schema.TypeInt,
					Description:  "Extra time in seconds the router waits before sending UP notification after the peer session is established",
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntAtLeast(0),
				},
			},
		},
	}
}

func resourceNsxtPolicyRouteController() *schema.Resource {
	rcSchema := metadata.GetSchemaFromExtendedSchema(routeControllerSchema)
	rcSchema["bgp_config"] = getRouteControllerBgpConfigSchema()
	return &schema.Resource{
		Create: resourceNsxtPolicyRouteControllerCreate,
		Read:   resourceNsxtPolicyRouteControllerRead,
		Update: resourceNsxtPolicyRouteControllerUpdate,
		Delete: resourceNsxtPolicyRouteControllerDelete,
		Importer: &schema.ResourceImporter{
			State: getPolicyPathResourceImporter(routeControllerPathExample),
		},
		Schema: rcSchema,
	}
}

func resourceNsxtPolicyRouteControllerExists(id string, connector client.Connector, isGlobalManager bool) (bool, error) {
	sessionContext := utl.SessionContext{ClientType: utl.Local}
	c := cliRouteControllersClient(sessionContext, connector)
	if c == nil {
		return false, fmt.Errorf("unsupported client type")
	}
	_, err := c.Get(id)
	if err == nil {
		return true, nil
	}
	if isNotFoundError(err) {
		return false, nil
	}
	return false, logAPIError("Error retrieving resource", err)
}

// initRouteControllerBgpChild wraps a RouteControllerBgpRoutingConfig as a
// ChildRouteControllerBgpRoutingConfig for inclusion in RouteController.Children
// during H-API calls. When markDelete is true a MarkedForDelete tombstone is returned.
func initRouteControllerBgpChild(bgpObj *model.RouteControllerBgpRoutingConfig, markDelete bool) (*data.StructValue, error) {
	converter := bindings.NewTypeConverter()
	bgpID := "bgp"
	child := model.ChildRouteControllerBgpRoutingConfig{
		ResourceType: "ChildRouteControllerBgpRoutingConfig",
	}
	if markDelete {
		boolTrue := true
		child.MarkedForDelete = &boolTrue
		child.RouteControllerBgpRoutingConfig = &model.RouteControllerBgpRoutingConfig{Id: &bgpID}
	} else {
		bgpType := "RouteControllerBgpRoutingConfig"
		bgpObj.Id = &bgpID
		bgpObj.ResourceType = &bgpType
		child.RouteControllerBgpRoutingConfig = bgpObj
	}
	dataValue, errs := converter.ConvertToVapi(child, model.ChildRouteControllerBgpRoutingConfigBindingType())
	if len(errs) > 0 {
		return nil, errs[0]
	}
	return dataValue.(*data.StructValue), nil
}

func routeControllerBgpConfigFromSchema(d *schema.ResourceData) *model.RouteControllerBgpRoutingConfig {
	bgpConfigList := d.Get("bgp_config").([]interface{})
	if len(bgpConfigList) == 0 || bgpConfigList[0] == nil {
		return nil
	}
	cfgMap := bgpConfigList[0].(map[string]interface{})

	ecmp := cfgMap["ecmp"].(bool)
	multipathRelax := cfgMap["multipath_relax"].(bool)
	restartMode := cfgMap["graceful_restart_mode"].(string)
	restartTimer := int64(cfgMap["graceful_restart_timer"].(int))

	restartTimerStruct := model.BgpGracefulRestartTimer{RestartTimer: &restartTimer}
	restartConfigStruct := model.BgpGracefulRestartConfig{
		Mode:  &restartMode,
		Timer: &restartTimerStruct,
	}

	bgpObj := model.RouteControllerBgpRoutingConfig{
		Ecmp:                  &ecmp,
		MultipathRelax:        &multipathRelax,
		GracefulRestartConfig: &restartConfigStruct,
	}

	if localAsNum, ok := cfgMap["local_as_num"].(string); ok && localAsNum != "" {
		bgpObj.LocalAsNum = &localAsNum
	}
	if timer, ok := cfgMap["peer_route_convergence_timer"].(int); ok && timer > 0 {
		t := int64(timer)
		bgpObj.PeerRouteConvergenceTimer = &t
	}

	return &bgpObj
}

func routeControllerToInfraStruct(d *schema.ResourceData, id string, markBgpDelete bool) (model.Infra, error) {
	var infraStruct model.Infra
	converter := bindings.NewTypeConverter()

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)
	rcType := "RouteController"

	rcObj := model.RouteController{
		DisplayName:  &displayName,
		Description:  &description,
		Tags:         tags,
		Id:           &id,
		ResourceType: &rcType,
	}

	// OmitIfEmpty metadata on ha_mode and virtual_network_appliance_cluster_path
	// ensures SchemaToStruct leaves those fields nil when the user omits them.
	elem := reflect.ValueOf(&rcObj).Elem()
	if err := metadata.SchemaToStruct(elem, d, routeControllerSchema, "", nil); err != nil {
		return infraStruct, err
	}

	// Embed BGP config as a child in the H-API payload so the entire operation is atomic.
	var rcChildren []*data.StructValue
	bgpObj := routeControllerBgpConfigFromSchema(d)
	if bgpObj != nil {
		bgpChild, err := initRouteControllerBgpChild(bgpObj, false)
		if err != nil {
			return infraStruct, fmt.Errorf("error building BGP child: %v", err)
		}
		rcChildren = append(rcChildren, bgpChild)
	} else if markBgpDelete {
		bgpChild, err := initRouteControllerBgpChild(nil, true)
		if err != nil {
			return infraStruct, fmt.Errorf("error building BGP delete child: %v", err)
		}
		rcChildren = append(rcChildren, bgpChild)
	}

	if len(rcChildren) > 0 {
		rcObj.Children = rcChildren
	}

	childRC := model.ChildRouteController{
		RouteController: &rcObj,
		ResourceType:    "ChildRouteController",
	}
	rcDataValue, errors := converter.ConvertToVapi(childRC, model.ChildRouteControllerBindingType())
	if errors != nil {
		return infraStruct, fmt.Errorf("error converting ChildRouteController: %v", errors[0])
	}

	infraType := "Infra"
	infraStruct = model.Infra{
		Children:     []*data.StructValue{rcDataValue.(*data.StructValue)},
		ResourceType: &infraType,
	}
	return infraStruct, nil
}

func resourceNsxtPolicyRouteControllerCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	context := getSessionContext(d, m)

	id, err := getOrGenerateID(d, m, resourceNsxtPolicyRouteControllerExists)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating RouteController with ID %s using H-API", id)

	infraObj, err := routeControllerToInfraStruct(d, id, false)
	if err != nil {
		return handleCreateError("RouteController", id, err)
	}

	if err = policyInfraPatch(context, infraObj, connector, false); err != nil {
		return handleCreateError("RouteController", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyRouteControllerRead(d, m)
}

func setRouteControllerBgpConfigInSchema(d *schema.ResourceData, bgpObj *model.RouteControllerBgpRoutingConfig) error {
	cfgMap := make(map[string]interface{})

	if bgpObj.Path != nil {
		cfgMap["path"] = *bgpObj.Path
	}
	if bgpObj.Revision != nil {
		cfgMap["revision"] = int(*bgpObj.Revision)
	}
	if bgpObj.Ecmp != nil {
		cfgMap["ecmp"] = *bgpObj.Ecmp
	}
	if bgpObj.MultipathRelax != nil {
		cfgMap["multipath_relax"] = *bgpObj.MultipathRelax
	}
	if bgpObj.LocalAsNum != nil {
		cfgMap["local_as_num"] = *bgpObj.LocalAsNum
	}
	if bgpObj.PeerRouteConvergenceTimer != nil {
		cfgMap["peer_route_convergence_timer"] = int(*bgpObj.PeerRouteConvergenceTimer)
	}

	cfgMap["graceful_restart_mode"] = model.BgpGracefulRestartConfig_MODE_HELPER_ONLY
	cfgMap["graceful_restart_timer"] = policyBGPGracefulRestartTimerDefault
	if bgpObj.GracefulRestartConfig != nil {
		if bgpObj.GracefulRestartConfig.Mode != nil {
			cfgMap["graceful_restart_mode"] = *bgpObj.GracefulRestartConfig.Mode
		}
		if bgpObj.GracefulRestartConfig.Timer != nil && bgpObj.GracefulRestartConfig.Timer.RestartTimer != nil {
			cfgMap["graceful_restart_timer"] = int(*bgpObj.GracefulRestartConfig.Timer.RestartTimer)
		}
	}

	return d.Set("bgp_config", []map[string]interface{}{cfgMap})
}

func resourceNsxtPolicyRouteControllerRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining RouteController ID")
	}

	sessionContext := getSessionContext(d, m)
	rcClient := cliRouteControllersClient(sessionContext, connector)
	if rcClient == nil {
		return fmt.Errorf("unsupported client type")
	}

	obj, err := rcClient.Get(id)
	if err != nil {
		return handleReadError(d, "RouteController", id, err)
	}

	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("revision", obj.Revision)
	d.Set("path", obj.Path)

	if obj.HaMode != nil {
		d.Set("ha_mode", *obj.HaMode)
	}
	if obj.VirtualNetworkApplianceClusterPath != nil {
		d.Set("virtual_network_appliance_cluster_path", *obj.VirtualNetworkApplianceClusterPath)
	}

	bgpClient := cliRouteControllerBgpClient(sessionContext, connector)
	if bgpClient != nil {
		bgpObj, bgpErr := bgpClient.Get(id)
		if bgpErr == nil {
			if err = setRouteControllerBgpConfigInSchema(d, &bgpObj); err != nil {
				return err
			}
		} else if !isNotFoundError(bgpErr) {
			return bgpErr
		}
	}

	return nil
}

func resourceNsxtPolicyRouteControllerUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	context := getSessionContext(d, m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining RouteController ID")
	}

	log.Printf("[INFO] Updating RouteController with ID %s using H-API", id)

	// When bgp_config is removed, send a MarkedForDelete child to clean it up atomically.
	markBgpDelete := false
	if d.HasChange("bgp_config") {
		bgpList := d.Get("bgp_config").([]interface{})
		if len(bgpList) == 0 || bgpList[0] == nil {
			markBgpDelete = true
		}
	}

	infraObj, err := routeControllerToInfraStruct(d, id, markBgpDelete)
	if err != nil {
		return handleUpdateError("RouteController", id, err)
	}

	if err = policyInfraPatch(context, infraObj, connector, false); err != nil {
		return handleUpdateError("RouteController", id, err)
	}

	return resourceNsxtPolicyRouteControllerRead(d, m)
}

func resourceNsxtPolicyRouteControllerDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining RouteController ID")
	}

	// Use H-API (Infra.Patch with MarkedForDelete) so NSX cascade-deletes any
	// BGP child atomically. A direct DELETE call fails with error 500030 when
	// the RC has a bgp_config child attached.
	converter := bindings.NewTypeConverter()
	boolTrue := true
	rcType := "RouteController"

	rcObj := model.RouteController{
		Id:           &id,
		ResourceType: &rcType,
	}

	childRC := model.ChildRouteController{
		MarkedForDelete: &boolTrue,
		RouteController: &rcObj,
		ResourceType:    "ChildRouteController",
	}

	rcDataValue, errors := converter.ConvertToVapi(childRC, model.ChildRouteControllerBindingType())
	if errors != nil {
		return fmt.Errorf("error converting ChildRouteController: %v", errors[0])
	}

	infraType := "Infra"
	infraObj := model.Infra{
		Children:     []*data.StructValue{rcDataValue.(*data.StructValue)},
		ResourceType: &infraType,
	}

	log.Printf("[INFO] Deleting RouteController with ID %s via H-API", id)
	if err := policyInfraPatch(getSessionContext(d, m), infraObj, getPolicyConnector(m), false); err != nil {
		return handleDeleteError("RouteController", id, err)
	}
	return nil
}
