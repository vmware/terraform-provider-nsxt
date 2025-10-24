package nsxt

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"

	segment "github.com/vmware/terraform-provider-nsxt/api/infra/segments"
	t1_segment "github.com/vmware/terraform-provider-nsxt/api/infra/tier_1s/segments"
	gm_port_profiles "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/global_infra/segments/ports"
	gm_t1_port_profiles "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/global_infra/tier_1s/segments/ports"
	gm_model "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	port_profiles "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/segments/ports"
	t1_port_profiles "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/tier_1s/segments/ports"

	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func policySegmentPortResourceToInfraStruct(context utl.SessionContext, id string, d *schema.ResourceData, isDestroy bool) (model.Infra, error) {
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getPolicyTagsFromSchema(d)
	revision := int64(d.Get("revision").(int))
	segmentPath := d.Get("segment_path").(string)

	// SegmentPort
	resourceType := "SegmentPort"
	obj := model.SegmentPort{
		Id:           &id,
		DisplayName:  &displayName,
		Tags:         tags,
		Revision:     &revision,
		ResourceType: &resourceType,
	}
	if description != "" {
		obj.Description = &description
	}

	err := nsxtPolicySegmentPortProfilesSetInStruct(d, &obj)
	if err != nil {
		return model.Infra{}, err
	}

	childSegmentPort := model.ChildSegmentPort{
		SegmentPort:     &obj,
		ResourceType:    "ChildSegmentPort",
		MarkedForDelete: &isDestroy,
	}

	// Segment
	child, err := vAPIConversion(childSegmentPort, model.ChildSegmentPortBindingType())
	if err != nil {
		return model.Infra{}, fmt.Errorf("Error handling the SegmentPort hierarchial API construction : %v", err)
	}
	segmentChildren := []*data.StructValue{child}
	segmentId := getSegmentIdFromSegPath(segmentPath)
	segmentTargetType := "Segment"
	childSegment := model.ChildResourceReference{
		Id:           &segmentId,
		ResourceType: "ChildResourceReference",
		TargetType:   &segmentTargetType,
		Children:     segmentChildren,
	}

	// Tier1
	child, err = vAPIConversion(childSegment, model.ChildResourceReferenceBindingType())
	if err != nil {
		return model.Infra{}, fmt.Errorf("Error handling the Tier1 gw hierarchial API construction : %v", err)
	}
	if isT1Segment(segmentPath) {
		t1Children := []*data.StructValue{child}
		t1GwId := getT1IdFromSegPath(segmentPath)
		tier1TargetType := "Tier1"
		childTier1Gw := model.ChildResourceReference{
			Id:           &t1GwId,
			ResourceType: "ChildResourceReference",
			TargetType:   &tier1TargetType,
			Children:     t1Children,
		}

		child, err = vAPIConversion(childTier1Gw, model.ChildResourceReferenceBindingType())
		if err != nil {
			return model.Infra{}, fmt.Errorf("Error handling the Infra hierarchial API construction : %v", err)
		}
	}
	// Infra
	infraType := "Infra"
	infraStruct := model.Infra{
		Children:     []*data.StructValue{child},
		ResourceType: &infraType,
	}

	return infraStruct, nil
}

func vAPIConversion(golangValue interface{}, bindingType bindings.BindingType) (*data.StructValue, error) {
	converter := bindings.NewTypeConverter()
	dataValue, errors := converter.ConvertToVapi(golangValue, bindingType)
	if errors != nil {
		return &data.StructValue{}, fmt.Errorf("Error converting Segment Child: %v", errors[0])
	}
	return dataValue.(*data.StructValue), nil
}

func nsxtPolicySegmentPortProfilesSetInStruct(d *schema.ResourceData, obj *model.SegmentPort) error {
	children := []*data.StructValue{}

	child, err := nsxtPolicyPortDiscoveryProfileSetInStruct(d)
	if err != nil {
		return err
	}

	if child != nil {
		children = append(children, child)
	}

	child, err = nsxtPolicyPortQosProfileSetInStruct(d)
	if err != nil {
		return err
	}

	if child != nil {
		children = append(children, child)
	}

	child, err = nsxtPolicyPortSecurityProfileSetInStruct(d)
	if err != nil {
		return err
	}

	if child != nil {
		children = append(children, child)
	}

	obj.Children = children
	return nil

}

func nsxtPolicyPortDiscoveryProfileSetInStruct(d *schema.ResourceData) (*data.StructValue, error) {
	segmentProfileMapID := "default"

	ipDiscoveryProfilePath := ""
	macDiscoveryProfilePath := ""
	revision := int64(0)
	shouldDelete := false
	oldProfiles, newProfiles := d.GetChange("discovery_profile")
	newProfilesList := newProfiles.([]interface{})

	if len(newProfilesList) > 0 && newProfilesList[0] != nil {
		profileMap := newProfilesList[0].(map[string]interface{})

		ipDiscoveryProfilePath = profileMap["ip_discovery_profile_path"].(string)
		macDiscoveryProfilePath = profileMap["mac_discovery_profile_path"].(string)
		if len(profileMap["binding_map_path"].(string)) > 0 {
			segmentProfileMapID = getPolicyIDFromPath(profileMap["binding_map_path"].(string))
		}

		revision = int64(profileMap["revision"].(int))
	} else {
		if len(oldProfiles.([]interface{})) == 0 {
			return nil, nil
		}
		segmentProfileMapID, revision = getOldProfileDataForRemoval(oldProfiles)
		shouldDelete = true
	}

	resourceType := "PortDiscoveryProfileBindingMap"
	discoveryMap := model.PortDiscoveryProfileBindingMap{
		ResourceType: &resourceType,
		Id:           &segmentProfileMapID,
	}

	if len(oldProfiles.([]interface{})) > 0 {
		// This is an update
		discoveryMap.Revision = &revision
	}

	if len(ipDiscoveryProfilePath) > 0 {
		discoveryMap.IpDiscoveryProfilePath = &ipDiscoveryProfilePath
	}

	if len(macDiscoveryProfilePath) > 0 {
		discoveryMap.MacDiscoveryProfilePath = &macDiscoveryProfilePath
	}

	childConfig := model.ChildPortDiscoveryProfileBindingMap{
		ResourceType:                   "ChildPortDiscoveryProfileBindingMap",
		PortDiscoveryProfileBindingMap: &discoveryMap,
		Id:                             &segmentProfileMapID,
		MarkedForDelete:                &shouldDelete,
	}

	converter := bindings.NewTypeConverter()
	dataValue, errors := converter.ConvertToVapi(childConfig, model.ChildPortDiscoveryProfileBindingMapBindingType())
	if errors != nil {
		return nil, fmt.Errorf("Error converting child segment discovery map: %v", errors[0])
	}

	return dataValue.(*data.StructValue), nil
}

func nsxtPolicyPortQosProfileSetInStruct(d *schema.ResourceData) (*data.StructValue, error) {
	segmentProfileMapID := "default"

	qosProfilePath := ""
	revision := int64(0)
	oldProfiles, newProfiles := d.GetChange("qos_profile")
	newProfilesList := newProfiles.([]interface{})
	shouldDelete := false
	if len(newProfilesList) > 0 && newProfilesList[0] != nil {
		profileMap := newProfilesList[0].(map[string]interface{})

		qosProfilePath = profileMap["qos_profile_path"].(string)
		if len(profileMap["binding_map_path"].(string)) > 0 {
			segmentProfileMapID = getPolicyIDFromPath(profileMap["binding_map_path"].(string))
		}

		revision = int64(profileMap["revision"].(int))
	} else {
		if len(oldProfiles.([]interface{})) == 0 {
			return nil, nil
		}
		// Profile should be deleted
		segmentProfileMapID, revision = getOldProfileDataForRemoval(oldProfiles)
		shouldDelete = true
	}

	resourceType := "PortQoSProfileBindingMap"
	qosMap := model.PortQosProfileBindingMap{
		ResourceType: &resourceType,
		Id:           &segmentProfileMapID,
	}

	if len(oldProfiles.([]interface{})) > 0 {
		// This is an update
		qosMap.Revision = &revision
	}

	if len(qosProfilePath) > 0 {
		qosMap.QosProfilePath = &qosProfilePath
	}

	childConfig := model.ChildPortQosProfileBindingMap{
		ResourceType:             "ChildPortQosProfileBindingMap",
		PortQosProfileBindingMap: &qosMap,
		Id:                       &segmentProfileMapID,
		MarkedForDelete:          &shouldDelete,
	}

	converter := bindings.NewTypeConverter()
	dataValue, errors := converter.ConvertToVapi(childConfig, model.ChildPortQosProfileBindingMapBindingType())
	if errors != nil {
		return nil, fmt.Errorf("Error converting child segment QoS map: %v", errors[0])
	}

	return dataValue.(*data.StructValue), nil
}

func nsxtPolicyPortSecurityProfileSetInStruct(d *schema.ResourceData) (*data.StructValue, error) {
	segmentProfileMapID := "default"

	spoofguardProfilePath := ""
	securityProfilePath := ""
	revision := int64(0)
	oldProfiles, newProfiles := d.GetChange("security_profile")
	newProfilesList := newProfiles.([]interface{})
	shouldDelete := false
	if len(newProfilesList) > 0 && newProfilesList[0] != nil {
		profileMap := newProfilesList[0].(map[string]interface{})

		spoofguardProfilePath = profileMap["spoofguard_profile_path"].(string)
		securityProfilePath = profileMap["security_profile_path"].(string)
		if len(profileMap["binding_map_path"].(string)) > 0 {
			segmentProfileMapID = getPolicyIDFromPath(profileMap["binding_map_path"].(string))
		}

		revision = int64(profileMap["revision"].(int))
	} else {
		if len(oldProfiles.([]interface{})) == 0 {
			return nil, nil
		}
		// Profile should be deleted
		segmentProfileMapID, revision = getOldProfileDataForRemoval(oldProfiles)
		shouldDelete = true
	}

	resourceType := "PortSecurityProfileBindingMap"
	securityMap := model.PortSecurityProfileBindingMap{
		ResourceType: &resourceType,
		Id:           &segmentProfileMapID,
	}

	if len(oldProfiles.([]interface{})) > 0 {
		// This is an update
		securityMap.Revision = &revision
	}

	if len(spoofguardProfilePath) > 0 {
		securityMap.SpoofguardProfilePath = &spoofguardProfilePath
	}

	if len(securityProfilePath) > 0 {
		securityMap.SegmentSecurityProfilePath = &securityProfilePath
	}

	childConfig := model.ChildPortSecurityProfileBindingMap{
		ResourceType:                  "ChildPortSecurityProfileBindingMap",
		PortSecurityProfileBindingMap: &securityMap,
		Id:                            &segmentProfileMapID,
		MarkedForDelete:               &shouldDelete,
	}

	converter := bindings.NewTypeConverter()
	dataValue, errors := converter.ConvertToVapi(childConfig, model.ChildPortSecurityProfileBindingMapBindingType())
	if errors != nil {
		return nil, fmt.Errorf("Error converting child segment security map: %v", errors[0])
	}

	return dataValue.(*data.StructValue), nil
}

func resourceNsxtPolicySegmentPortExists(d *schema.ResourceData, context utl.SessionContext, connector client.Connector) func(context utl.SessionContext, id string, connector client.Connector) (bool, error) {
	segmentPath := d.Get("segment_path").(string)
	return func(context utl.SessionContext, id string, connector client.Connector) (bool, error) {
		_, err := getSegmentPort(segmentPath, id, context, connector)
		if err == nil {
			return true, nil
		}
		if isNotFoundError(err) {
			return false, nil
		}
		return false, logAPIError("Error retrieving Segment Port", err)
	}
}

func getSegmentPort(segmentPath, segmentPortId string, context utl.SessionContext, connector client.Connector) (model.SegmentPort, error) {
	var segPort model.SegmentPort
	var err error
	segmentId := getSegmentIdFromSegPath(segmentPath)
	t1Id := getT1IdFromSegPath(segmentPath)
	if isT1Segment(segmentPath) {
		if t1Id == "" {
			return model.SegmentPort{}, fmt.Errorf("Error getting the tier1 gateway ID : %v", err)
		}
		portsT1Client := t1_segment.NewPortsClient(context, connector)
		segPort, err = portsT1Client.Get(t1Id, segmentId, segmentPortId)
	} else {
		portsClient := segment.NewPortsClient(context, connector)
		segPort, err = portsClient.Get(segmentId, segmentPortId)
	}
	return segPort, err
}

func isT1Segment(segmentPath string) bool {
	pathSplit := strings.Split(segmentPath, "/")
	if len(pathSplit) >= 3 && pathSplit[len(pathSplit)-4] == "tier-1s" {
		return true
	}
	return false
}

func getSegmentIdFromSegPath(segPortPath string) string {
	pathSplit := strings.Split(segPortPath, "/")
	return pathSplit[len(pathSplit)-1]
}

func getT1IdFromSegPath(segPortPath string) string {
	pathSplit := strings.Split(segPortPath, "/")
	if len(pathSplit) >= 3 && pathSplit[len(pathSplit)-4] == "tier-1s" {
		return pathSplit[len(pathSplit)-3]
	}
	return ""
}

type segmentConfig interface {
	nsxtPolicySegmentPortDiscoveryProfileRead(d *schema.ResourceData, m interface{}) error
	nsxtPolicySegmentPortQosProfileRead(d *schema.ResourceData, m interface{}) error
	nsxtPolicyPortSegmentSecurityProfileRead(d *schema.ResourceData, m interface{}) error
}

type segmentPort struct {
	segmentId string
	portId    string
}

type tier1SegmentPort struct {
	tier1GatewayId string
	ids            *segmentPort
}

func nsxtPolicySegmentPortProfilesRead(d *schema.ResourceData, m interface{}) error {
	var config segmentConfig
	segmentPath := d.Get("segment_path").(string)
	s := segmentPort{
		segmentId: getSegmentIdFromSegPath(segmentPath),
		portId:    d.Id(),
	}

	config = segmentConfig(s)
	if isT1Segment(segmentPath) {
		t := tier1SegmentPort{
			tier1GatewayId: getT1IdFromSegPath(segmentPath),
			ids:            &s,
		}
		config = segmentConfig(t)
	}
	err := config.nsxtPolicySegmentPortDiscoveryProfileRead(d, m)
	if err != nil {
		return err
	}

	err = config.nsxtPolicySegmentPortQosProfileRead(d, m)
	if err != nil {
		return err
	}

	err = config.nsxtPolicyPortSegmentSecurityProfileRead(d, m)
	if err != nil {
		return err
	}

	return nil
}

func (c segmentPort) nsxtPolicySegmentPortDiscoveryProfileRead(d *schema.ResourceData, m interface{}) error {
	errorMessage := "Failed to read Discovery Profile Map for segment port %s: %s"
	connector := getPolicyConnector(m)

	var results model.PortDiscoveryProfileBindingMapListResult
	if isPolicyGlobalManager(m) {
		client := gm_port_profiles.NewPortDiscoveryProfileBindingMapsClient(connector)
		gmResults, err := client.List(c.segmentId, c.portId, nil, nil, nil, nil, nil, nil)
		if err != nil {
			return fmt.Errorf(errorMessage, c.portId, err)
		}
		lmResults, err := convertModelBindingType(gmResults, gm_model.PortDiscoveryProfileBindingMapListResultBindingType(), model.PortDiscoveryProfileBindingMapListResultBindingType())
		if err != nil {
			return err
		}
		results = lmResults.(model.PortDiscoveryProfileBindingMapListResult)
	} else {
		client := port_profiles.NewPortDiscoveryProfileBindingMapsClient(connector)
		if client == nil {
			return policyResourceNotSupportedError()
		}
		var err error
		results, err = client.List(c.segmentId, c.portId, nil, nil, nil, nil, nil, nil)
		if err != nil {
			return fmt.Errorf(errorMessage, c.portId, err)
		}
	}

	config := make(map[string]interface{})
	var configList []map[string]interface{}

	for _, obj := range results.Results {
		config["ip_discovery_profile_path"] = obj.IpDiscoveryProfilePath
		config["mac_discovery_profile_path"] = obj.MacDiscoveryProfilePath
		config["binding_map_path"] = obj.Path
		config["revision"] = obj.Revision
		configList = append(configList, config)
		d.Set("discovery_profile", configList)
		return nil
	}

	return nil
}

func (c segmentPort) nsxtPolicySegmentPortQosProfileRead(d *schema.ResourceData, m interface{}) error {
	errorMessage := "Failed to read QoS Profile Map for segment port %s: %s"
	connector := getPolicyConnector(m)
	var results model.PortQosProfileBindingMapListResult
	if isPolicyGlobalManager(m) {
		client := gm_port_profiles.NewPortQosProfileBindingMapsClient(connector)
		gmResults, err := client.List(c.segmentId, c.portId, nil, nil, nil, nil, nil)
		if err != nil {
			return fmt.Errorf(errorMessage, c.portId, err)
		}
		lmResults, err := convertModelBindingType(gmResults, gm_model.PortQosProfileBindingMapListResultBindingType(), model.PortQosProfileBindingMapListResultBindingType())
		if err != nil {
			return err
		}
		results = lmResults.(model.PortQosProfileBindingMapListResult)
	} else {
		client := port_profiles.NewPortQosProfileBindingMapsClient(connector)
		if client == nil {
			return policyResourceNotSupportedError()
		}
		var err error
		results, err = client.List(c.segmentId, c.portId, nil, nil, nil, nil, nil)
		if err != nil {
			return fmt.Errorf(errorMessage, c.portId, err)
		}
	}

	config := make(map[string]interface{})
	var configList []map[string]interface{}

	for _, obj := range results.Results {
		if obj.QosProfilePath != nil && (len(*obj.QosProfilePath) > 0) {
			config["qos_profile_path"] = obj.QosProfilePath
			config["binding_map_path"] = obj.Path
			config["revision"] = obj.Revision
			configList = append(configList, config)
			d.Set("qos_profile", configList)
			return nil
		}
	}

	return nil
}

func (c segmentPort) nsxtPolicyPortSegmentSecurityProfileRead(d *schema.ResourceData, m interface{}) error {
	errorMessage := "Failed to read Security Profile Map for segment port %s: %s"
	connector := getPolicyConnector(m)
	var results model.PortSecurityProfileBindingMapListResult
	if isPolicyGlobalManager(m) {
		client := gm_port_profiles.NewPortSecurityProfileBindingMapsClient(connector)
		gmResults, err := client.List(c.segmentId, c.portId, nil, nil, nil, nil, nil)
		if err != nil {
			return fmt.Errorf(errorMessage, c.portId, err)
		}
		lmResults, err := convertModelBindingType(gmResults, gm_model.PortSecurityProfileBindingMapListResultBindingType(), model.PortSecurityProfileBindingMapListResultBindingType())
		if err != nil {
			return err
		}
		results = lmResults.(model.PortSecurityProfileBindingMapListResult)
	} else {
		client := port_profiles.NewPortSecurityProfileBindingMapsClient(connector)
		if client == nil {
			return policyResourceNotSupportedError()
		}
		var err error
		results, err = client.List(c.segmentId, c.portId, nil, nil, nil, nil, nil)
		if err != nil {
			return fmt.Errorf(errorMessage, c.portId, err)
		}
	}

	config := make(map[string]interface{})
	var configList []map[string]interface{}

	for _, obj := range results.Results {
		config["security_profile_path"] = obj.SegmentSecurityProfilePath
		config["spoofguard_profile_path"] = obj.SpoofguardProfilePath
		config["binding_map_path"] = obj.Path
		config["revision"] = obj.Revision
		configList = append(configList, config)
		d.Set("security_profile", configList)
		return nil
	}

	return nil
}

func (c tier1SegmentPort) nsxtPolicySegmentPortDiscoveryProfileRead(d *schema.ResourceData, m interface{}) error {
	errorMessage := "Failed to read Discovery Profile Map for segment port %s: %s"
	connector := getPolicyConnector(m)

	var results model.PortDiscoveryProfileBindingMapListResult
	if isPolicyGlobalManager(m) {
		client := gm_t1_port_profiles.NewPortDiscoveryProfileBindingMapsClient(connector)
		gmResults, err := client.List(c.tier1GatewayId, c.ids.segmentId, c.ids.portId, nil, nil, nil, nil, nil, nil)
		if err != nil {
			return fmt.Errorf(errorMessage, c.ids.portId, err)
		}
		lmResults, err := convertModelBindingType(gmResults, gm_model.PortDiscoveryProfileBindingMapListResultBindingType(), model.PortDiscoveryProfileBindingMapListResultBindingType())
		if err != nil {
			return err
		}
		results = lmResults.(model.PortDiscoveryProfileBindingMapListResult)
	} else {
		client := t1_port_profiles.NewPortDiscoveryProfileBindingMapsClient(connector)
		if client == nil {
			return policyResourceNotSupportedError()
		}
		var err error
		results, err = client.List(c.tier1GatewayId, c.ids.segmentId, c.ids.portId, nil, nil, nil, nil, nil, nil)
		if err != nil {
			return fmt.Errorf(errorMessage, c.ids.portId, err)
		}
	}

	config := make(map[string]interface{})
	var configList []map[string]interface{}

	for _, obj := range results.Results {
		config["ip_discovery_profile_path"] = obj.IpDiscoveryProfilePath
		config["mac_discovery_profile_path"] = obj.MacDiscoveryProfilePath
		config["binding_map_path"] = obj.Path
		config["revision"] = obj.Revision
		configList = append(configList, config)
		d.Set("discovery_profile", configList)
		return nil
	}

	return nil
}

func (c tier1SegmentPort) nsxtPolicySegmentPortQosProfileRead(d *schema.ResourceData, m interface{}) error {
	errorMessage := "Failed to read QoS Profile Map for segment port %s: %s"
	connector := getPolicyConnector(m)
	var results model.PortQosProfileBindingMapListResult
	if isPolicyGlobalManager(m) {
		client := gm_t1_port_profiles.NewPortQosProfileBindingMapsClient(connector)
		gmResults, err := client.List(c.tier1GatewayId, c.ids.segmentId, c.ids.portId, nil, nil, nil, nil, nil)
		if err != nil {
			return fmt.Errorf(errorMessage, c.ids.portId, err)
		}
		lmResults, err := convertModelBindingType(gmResults, gm_model.PortQosProfileBindingMapListResultBindingType(), model.PortQosProfileBindingMapListResultBindingType())
		if err != nil {
			return err
		}
		results = lmResults.(model.PortQosProfileBindingMapListResult)
	} else {
		client := t1_port_profiles.NewPortQosProfileBindingMapsClient(connector)
		if client == nil {
			return policyResourceNotSupportedError()
		}
		var err error
		results, err = client.List(c.tier1GatewayId, c.ids.segmentId, c.ids.portId, nil, nil, nil, nil, nil)
		if err != nil {
			return fmt.Errorf(errorMessage, c.ids.portId, err)
		}
	}

	config := make(map[string]interface{})
	var configList []map[string]interface{}

	for _, obj := range results.Results {
		if obj.QosProfilePath != nil && (len(*obj.QosProfilePath) > 0) {
			config["qos_profile_path"] = obj.QosProfilePath
			config["binding_map_path"] = obj.Path
			config["revision"] = obj.Revision
			configList = append(configList, config)
			d.Set("qos_profile", configList)
			return nil
		}
	}

	return nil
}

func (c tier1SegmentPort) nsxtPolicyPortSegmentSecurityProfileRead(d *schema.ResourceData, m interface{}) error {
	errorMessage := "Failed to read Security Profile Map for segment port %s: %s"
	connector := getPolicyConnector(m)
	var results model.PortSecurityProfileBindingMapListResult
	if isPolicyGlobalManager(m) {
		client := gm_t1_port_profiles.NewPortSecurityProfileBindingMapsClient(connector)
		gmResults, err := client.List(c.tier1GatewayId, c.ids.segmentId, c.ids.portId, nil, nil, nil, nil, nil)
		if err != nil {
			return fmt.Errorf(errorMessage, c.ids.portId, err)
		}
		lmResults, err := convertModelBindingType(gmResults, gm_model.PortSecurityProfileBindingMapListResultBindingType(), model.PortSecurityProfileBindingMapListResultBindingType())
		if err != nil {
			return err
		}
		results = lmResults.(model.PortSecurityProfileBindingMapListResult)
	} else {
		client := t1_port_profiles.NewPortSecurityProfileBindingMapsClient(connector)
		if client == nil {
			return policyResourceNotSupportedError()
		}
		var err error
		results, err = client.List(c.tier1GatewayId, c.ids.segmentId, c.ids.portId, nil, nil, nil, nil, nil)
		if err != nil {
			return fmt.Errorf(errorMessage, c.ids.portId, err)
		}
	}

	config := make(map[string]interface{})
	var configList []map[string]interface{}

	for _, obj := range results.Results {
		config["security_profile_path"] = obj.SegmentSecurityProfilePath
		config["spoofguard_profile_path"] = obj.SpoofguardProfilePath
		config["binding_map_path"] = obj.Path
		config["revision"] = obj.Revision
		configList = append(configList, config)
		d.Set("security_profile", configList)
		return nil
	}

	return nil
}
