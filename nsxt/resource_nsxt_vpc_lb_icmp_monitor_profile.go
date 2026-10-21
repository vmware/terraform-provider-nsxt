// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func resourceNsxtVpcLBIcmpMonitorProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtVpcLBIcmpMonitorProfileCreate,
		Read:   resourceNsxtVpcLBIcmpMonitorProfileRead,
		Update: resourceNsxtVpcLBIcmpMonitorProfileUpdate,
		Delete: resourceNsxtVpcLBIcmpMonitorProfileDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtVersionCheckImporter("9.0.0", "VPC LB ICMP Monitor Profile", getVpcPathResourceImporter(vpcLbMonitorProfilePathExample)),
		},

		Schema: map[string]*schema.Schema{
			"nsx_id":       getNsxIDSchema(),
			"path":         getPathSchema(),
			"display_name": getDisplayNameSchema(),
			"description":  getDescriptionSchema(),
			"revision":     getRevisionSchema(),
			"tag":          getTagsSchema(),
			"context":      getContextSchemaExtended(true, false, true, true),
			"data_length": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The data size (in byte) of the ICMP healthcheck packet",
			},
			"fall_count": getLbMonitorFallCountSchema(),
			"interval":   getLbMonitorIntervalSchema(),
			"rise_count": getLbMonitorRiseCountSchema(),
			"timeout":    getLbMonitorTimeoutSchema(),
		},
	}
}

func resourceNsxtVpcLBIcmpMonitorProfilePatch(d *schema.ResourceData, m interface{}, id string) error {
	connector := getPolicyConnector(m)
	converter := bindings.NewTypeConverter()

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsWithProviderManagedDefaults(d, m)
	dataLength := int64(d.Get("data_length").(int))
	fallCount := int64(d.Get("fall_count").(int))
	interval := int64(d.Get("interval").(int))
	riseCount := int64(d.Get("rise_count").(int))
	timeout := int64(d.Get("timeout").(int))
	resourceType := model.LBMonitorProfile_RESOURCE_TYPE_LBICMPMONITORPROFILE

	obj := model.LBIcmpMonitorProfile{
		DisplayName:  &displayName,
		Description:  &description,
		Tags:         tags,
		DataLength:   &dataLength,
		FallCount:    &fallCount,
		Interval:     &interval,
		RiseCount:    &riseCount,
		Timeout:      &timeout,
		ResourceType: resourceType,
	}

	log.Printf("[INFO] Patching VPC LBIcmpMonitorProfile with ID %s", id)
	dataValue, errs := converter.ConvertToVapi(obj, model.LBIcmpMonitorProfileBindingType())
	if errs != nil {
		return fmt.Errorf("VPC LBMonitorProfile %s is not of type LBIcmpMonitorProfile %s", id, errs[0])
	}
	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbMonitorProfilesClient(sessionContext, connector)
	return c.Patch(parents[0], parents[1], parents[2], id, dataValue.(*data.StructValue))
}

func resourceNsxtVpcLBIcmpMonitorProfileCreate(d *schema.ResourceData, m interface{}) error {
	if !util.NsxVersionHigherOrEqual("9.0.0") {
		return fmt.Errorf("VPC LB ICMP Monitor Profile resource requires NSX version 9.0.0 or higher")
	}

	// Initialize resource Id and verify this ID is not yet used
	id, err := getOrGenerateID2(d, m, resourceNsxtVpcLBMonitorProfileExists)
	if err != nil {
		return err
	}

	err = resourceNsxtVpcLBIcmpMonitorProfilePatch(d, m, id)
	if err != nil {
		return handleCreateError("VPC LBIcmpMonitorProfile", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtVpcLBIcmpMonitorProfileRead(d, m)
}

func resourceNsxtVpcLBIcmpMonitorProfileRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	converter := bindings.NewTypeConverter()

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VPC LBIcmpMonitorProfile ID")
	}

	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbMonitorProfilesClient(sessionContext, connector)
	obj, err := c.Get(parents[0], parents[1], parents[2], id)
	if err != nil {
		return handleReadError(d, "VPC LBHttpsMonitorProfile", id, err)
	}

	baseObj, errs := converter.ConvertToGolang(obj, model.LBIcmpMonitorProfileBindingType())
	if len(errs) > 0 {
		return fmt.Errorf("Error converting VPC LBIcmpMonitorProfile %s", errs[0])
	}
	lbICMPMonitor := baseObj.(model.LBIcmpMonitorProfile)

	d.Set("display_name", lbICMPMonitor.DisplayName)
	d.Set("description", lbICMPMonitor.Description)
	setPolicyTagsInSchema(d, lbICMPMonitor.Tags)
	d.Set("nsx_id", id)
	d.Set("path", lbICMPMonitor.Path)
	d.Set("revision", lbICMPMonitor.Revision)

	d.Set("data_length", lbICMPMonitor.DataLength)
	d.Set("fall_count", lbICMPMonitor.FallCount)
	d.Set("interval", lbICMPMonitor.Interval)
	d.Set("rise_count", lbICMPMonitor.RiseCount)
	d.Set("timeout", lbICMPMonitor.Timeout)

	return nil
}

func resourceNsxtVpcLBIcmpMonitorProfileUpdate(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VPC LBIcmpMonitorProfile ID")
	}

	err := resourceNsxtVpcLBIcmpMonitorProfilePatch(d, m, id)
	if err != nil {
		return handleUpdateError("VPC LBIcmpMonitorProfile", id, err)
	}

	return resourceNsxtVpcLBIcmpMonitorProfileRead(d, m)
}

func resourceNsxtVpcLBIcmpMonitorProfileDelete(d *schema.ResourceData, m interface{}) error {
	return resourceNsxtVpcLBMonitorProfileDelete(d, m)
}
