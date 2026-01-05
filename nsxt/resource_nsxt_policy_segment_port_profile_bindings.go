package nsxt

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func resourceNsxtPolicySegmentPortProfileBindings() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicySegmentPortProfileBindingsCreate,
		Read:   resourceNsxtPolicySegmentPortProfileBindingsRead,
		Update: resourceNsxtPolicySegmentPortProfileBindingsUpdate,
		Delete: resourceNsxtPolicySegmentPortProfileBindingsDelete,
		Importer: &schema.ResourceImporter{
			State: resourceNsxtPolicySegmentPortProfileBindingsImporter,
		},
		Schema: map[string]*schema.Schema{
			"context": getContextSchema(false, false, false),
			"segment_port_path": {
				Type:        schema.TypeString,
				Description: "Policy path of the segment port",
				Required:    true,
				ForceNew:    true,
			},
			"discovery_profile": {
				Type:        schema.TypeList,
				Description: "IP and MAC discovery profiles for this segment port",
				Elem:        getPolicySegmentDiscoveryProfilesSchema(),
				Optional:    true,
				MaxItems:    1,
			},
			"qos_profile": {
				Type:        schema.TypeList,
				Description: "QoS profiles for this segment port",
				Elem:        getPolicySegmentQosProfilesSchema(),
				Optional:    true,
				MaxItems:    1,
			},
			"security_profile": {
				Type:        schema.TypeList,
				Description: "Security profiles for this segment port",
				Elem:        getPolicySegmentSecurityProfilesSchema(),
				Optional:    true,
				MaxItems:    1,
			},
		},
	}
}

func resourceNsxtPolicySegmentPortProfileBindingsCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	context := getSessionContext(d, m)
	segmentPortPath := d.Get("segment_port_path").(string)
	segmentPortID := getPolicyIDFromPath(segmentPortPath)
	segmentPath, err := getPolicySegmentPathFromPortPath(segmentPortPath)
	if err != nil {
		return fmt.Errorf("Error parsing Segment Port Path: %v", err)
	}

	segmentPort, err := getSegmentPort(segmentPath, segmentPortID, context, connector)
	if err != nil {
		return fmt.Errorf("Error getting Segment Port %s: %v", segmentPortID, err)
	}

	obj, err := policySegmentPortProfileBindingsResourceToInfraStruct(segmentPort, d, false)
	if err != nil {
		return err
	}

	err = policyInfraPatch(context, obj, connector, false)
	if err != nil {
		return handleCreateError("SegmentPortProfileBindings", segmentPortID, err)
	}

	d.SetId(segmentPortID)

	return resourceNsxtPolicySegmentPortProfileBindingsRead(d, m)
}

func resourceNsxtPolicySegmentPortProfileBindingsRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	segmentPortPath := d.Get("segment_port_path").(string)
	segmentPortID := getPolicyIDFromPath(segmentPortPath)
	segmentPath, err := getPolicySegmentPathFromPortPath(segmentPortPath)
	if err != nil {
		return fmt.Errorf("Error parsing Segment Port Path: %v", err)
	}

	_, err = getSegmentPort(segmentPath, segmentPortID, getSessionContext(d, m), connector)
	if err != nil {
		if isNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error getting Segment Port: %v", err)
	}

	err = nsxtPolicySegmentPortProfilesRead(d, m)
	if err != nil {
		return err
	}

	return nil
}

func resourceNsxtPolicySegmentPortProfileBindingsUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	context := getSessionContext(d, m)

	segmentPortPath := d.Get("segment_port_path").(string)
	segmentPortID := getPolicyIDFromPath(segmentPortPath)
	segmentPath, err := getPolicySegmentPathFromPortPath(segmentPortPath)
	if err != nil {
		return fmt.Errorf("Error parsing Segment Port Path: %v", err)
	}
	segmentPort, err := getSegmentPort(segmentPath, segmentPortID, context, connector)
	if err != nil {
		return fmt.Errorf("Error getting Segment Port: %v", err)
	}

	obj, err := policySegmentPortProfileBindingsResourceToInfraStruct(segmentPort, d, false)
	if err != nil {
		return err
	}

	err = policyInfraPatch(context, obj, connector, false)
	if err != nil {
		return handleUpdateError("SegmentPortProfileBindings", segmentPortID, err)
	}

	return resourceNsxtPolicySegmentPortProfileBindingsRead(d, m)
}

func resourceNsxtPolicySegmentPortProfileBindingsDelete(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	context := getSessionContext(d, m)

	segmentPortPath := d.Get("segment_port_path").(string)
	segmentPortID := getPolicyIDFromPath(segmentPortPath)
	segmentPath, err := getPolicySegmentPathFromPortPath(segmentPortPath)
	if err != nil {
		return fmt.Errorf("Error parsing Segment Port Path: %v", err)
	}
	segmentPort, err := getSegmentPort(segmentPath, segmentPortID, context, connector)
	if err != nil {
		if isNotFoundError(err) {
			return nil
		}
		return fmt.Errorf("Error getting Segment Port: %v", err)
	}

	obj, err := policySegmentPortProfileBindingsResourceToInfraStruct(segmentPort, d, true)
	if err != nil {
		return err
	}

	err = policyInfraPatch(context, obj, connector, false)
	if err != nil {
		return handleDeleteError("SegmentPortProfileBindings", segmentPortID, err)
	}

	return nil
}

func policySegmentPortProfileBindingsResourceToInfraStruct(segmentPort model.SegmentPort, d *schema.ResourceData, isDestroy bool) (model.Infra, error) {
	segmentPortPath := d.Get("segment_port_path").(string)
	markedForDelete := false
	segmentPath, err := getParameterFromPolicyPath("/segments/", "/ports/", segmentPortPath)
	if err != nil {
		return model.Infra{}, err
	}

	if isDestroy {
		// For delete operation, set empty profile bindings
		err = nsxtPolicySegmentPortProfilesSetEmptyInStruct(&segmentPort)
		if err != nil {
			return model.Infra{}, err
		}
	} else {
		err = nsxtPolicySegmentPortProfilesSetInStruct(d, &segmentPort)
		if err != nil {
			return model.Infra{}, err
		}
	}

	childSegmentPort := model.ChildSegmentPort{
		SegmentPort:     &segmentPort,
		ResourceType:    "ChildSegmentPort",
		MarkedForDelete: &markedForDelete,
	}

	// Segment
	child, err := vAPIConversion(childSegmentPort, model.ChildSegmentPortBindingType())
	if err != nil {
		return model.Infra{}, fmt.Errorf("Error handling the SegmentPort hierarchical API construction : %v", err)
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
		return model.Infra{}, fmt.Errorf("Error handling the Tier1 gw hierarchical API construction : %v", err)
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
			return model.Infra{}, fmt.Errorf("Error handling the Infra hierarchical API construction : %v", err)
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

func resourceNsxtPolicySegmentPortProfileBindingsImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	importID := d.Id()
	segmentPortPath := importID

	segmentPortID := getPolicyIDFromPath(segmentPortPath)
	if segmentPortID == "" {
		return nil, fmt.Errorf("invalid segment port path %s: cannot extract port ID", importID)
	}

	segmentPath, err := getPolicySegmentPathFromPortPath(segmentPortPath)
	if err != nil {
		return nil, fmt.Errorf("invalid segment port path %s: %v", importID, err)
	}

	d.Set("segment_port_path", segmentPortPath)

	rd, err := nsxtPolicyPathResourceImporterHelper(d, m)
	if err != nil {
		return nil, err
	}

	d.SetId(segmentPortID)

	connector := getPolicyConnector(m)
	context := getSessionContext(d, m)

	_, err = getSegmentPort(segmentPath, segmentPortID, context, connector)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve segment port %s: %v", segmentPortID, err)
	}

	err = nsxtPolicySegmentPortProfilesRead(d, m)
	if err != nil {
		return nil, fmt.Errorf("failed to read profile bindings for segment port %s: %v", segmentPortID, err)
	}

	return rd, nil
}
