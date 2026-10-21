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

func resourceNsxtVpcLBUdpMonitorProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtVpcLBUdpMonitorProfileCreate,
		Read:   resourceNsxtVpcLBUdpMonitorProfileRead,
		Update: resourceNsxtVpcLBUdpMonitorProfileUpdate,
		Delete: resourceNsxtVpcLBUdpMonitorProfileDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtVersionCheckImporter("9.0.0", "VPC LB UDP Monitor Profile", getVpcPathResourceImporter(vpcLbMonitorProfilePathExample)),
		},

		Schema: map[string]*schema.Schema{
			"nsx_id":       getNsxIDSchema(),
			"path":         getPathSchema(),
			"display_name": getDisplayNameSchema(),
			"description":  getDescriptionSchema(),
			"revision":     getRevisionSchema(),
			"tag":          getTagsSchema(),
			"context":      getContextSchemaExtended(true, false, true, true),
			"receive": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The expected data string to be received from the response, can be anywhere in the response",
			},
			"send": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The data to be sent to the monitored server.",
			},
			"fall_count":   getLbMonitorFallCountSchema(),
			"interval":     getLbMonitorIntervalSchema(),
			"monitor_port": getPolicyLbMonitorPortSchema(),
			"rise_count":   getLbMonitorRiseCountSchema(),
			"timeout":      getLbMonitorTimeoutSchema(),
		},
	}
}

func resourceNsxtVpcLBUdpMonitorProfilePatch(d *schema.ResourceData, m interface{}, id string) error {
	connector := getPolicyConnector(m)
	converter := bindings.NewTypeConverter()

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsWithProviderManagedDefaults(d, m)
	receive := d.Get("receive").(string)
	send := d.Get("send").(string)
	fallCount := int64(d.Get("fall_count").(int))
	interval := int64(d.Get("interval").(int))
	monitorPort := int64(d.Get("monitor_port").(int))
	riseCount := int64(d.Get("rise_count").(int))
	timeout := int64(d.Get("timeout").(int))
	resourceType := model.LBMonitorProfile_RESOURCE_TYPE_LBUDPMONITORPROFILE

	obj := model.LBUdpMonitorProfile{
		DisplayName:  &displayName,
		Description:  &description,
		Tags:         tags,
		Receive:      &receive,
		Send:         &send,
		FallCount:    &fallCount,
		Interval:     &interval,
		MonitorPort:  &monitorPort,
		RiseCount:    &riseCount,
		Timeout:      &timeout,
		ResourceType: resourceType,
	}

	log.Printf("[INFO] Patching VPC LBUdpMonitorProfile with ID %s", id)

	dataValue, errs := converter.ConvertToVapi(obj, model.LBUdpMonitorProfileBindingType())
	if errs != nil {
		return fmt.Errorf("Error converting VPC LBUdpMonitorProfile %s", errs[0])
	}

	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbMonitorProfilesClient(sessionContext, connector)
	return c.Patch(parents[0], parents[1], parents[2], id, dataValue.(*data.StructValue))
}

func resourceNsxtVpcLBUdpMonitorProfileCreate(d *schema.ResourceData, m interface{}) error {
	if !util.NsxVersionHigherOrEqual("9.0.0") {
		return fmt.Errorf("VPC LB UDP Monitor Profile resource requires NSX version 9.0.0 or higher")
	}

	// Initialize resource Id and verify this ID is not yet used
	id, err := getOrGenerateID2(d, m, resourceNsxtVpcLBMonitorProfileExists)
	if err != nil {
		return err
	}

	err = resourceNsxtVpcLBUdpMonitorProfilePatch(d, m, id)
	if err != nil {
		return handleCreateError("VPC LBUdpMonitorProfile", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtVpcLBUdpMonitorProfileRead(d, m)
}

func resourceNsxtVpcLBUdpMonitorProfileRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VPC LBUdpMonitorProfile ID")
	}

	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbMonitorProfilesClient(sessionContext, connector)
	obj, err := c.Get(parents[0], parents[1], parents[2], id)
	if err != nil {
		return handleReadError(d, "VPC LBUdpMonitorProfile", id, err)
	}

	converter := bindings.NewTypeConverter()
	baseObj, errs := converter.ConvertToGolang(obj, model.LBUdpMonitorProfileBindingType())
	if len(errs) > 0 {
		return fmt.Errorf("Error converting VPC LBUdpMonitorProfile %s", errs[0])
	}
	lbUDPMonitor := baseObj.(model.LBUdpMonitorProfile)

	d.Set("display_name", lbUDPMonitor.DisplayName)
	d.Set("description", lbUDPMonitor.Description)
	setPolicyTagsInSchema(d, lbUDPMonitor.Tags)
	d.Set("nsx_id", id)
	d.Set("path", lbUDPMonitor.Path)
	d.Set("revision", lbUDPMonitor.Revision)

	d.Set("receive", lbUDPMonitor.Receive)
	d.Set("send", lbUDPMonitor.Send)
	d.Set("fall_count", lbUDPMonitor.FallCount)
	d.Set("interval", lbUDPMonitor.Interval)
	d.Set("monitor_port", lbUDPMonitor.MonitorPort)
	d.Set("rise_count", lbUDPMonitor.RiseCount)
	d.Set("timeout", lbUDPMonitor.Timeout)

	return nil
}

func resourceNsxtVpcLBUdpMonitorProfileUpdate(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VPC LBUdpMonitorProfile ID")
	}

	err := resourceNsxtVpcLBUdpMonitorProfilePatch(d, m, id)
	if err != nil {
		return handleUpdateError("VPC LBUdpMonitorProfile", id, err)
	}

	return resourceNsxtVpcLBUdpMonitorProfileRead(d, m)
}

func resourceNsxtVpcLBUdpMonitorProfileDelete(d *schema.ResourceData, m interface{}) error {
	return resourceNsxtVpcLBMonitorProfileDelete(d, m)
}
