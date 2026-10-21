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

func resourceNsxtVpcLBTcpMonitorProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtVpcLBTcpMonitorProfileCreate,
		Read:   resourceNsxtVpcLBTcpMonitorProfileRead,
		Update: resourceNsxtVpcLBTcpMonitorProfileUpdate,
		Delete: resourceNsxtVpcLBTcpMonitorProfileDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtVersionCheckImporter("9.0.0", "VPC LB TCP Monitor Profile", getVpcPathResourceImporter(vpcLbMonitorProfilePathExample)),
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

func resourceNsxtVpcLBTcpMonitorProfilePatch(d *schema.ResourceData, m interface{}, id string) error {
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
	resourceType := model.LBMonitorProfile_RESOURCE_TYPE_LBTCPMONITORPROFILE

	obj := model.LBTcpMonitorProfile{
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

	log.Printf("[INFO] Patching VPC LBTcpMonitorProfile with ID %s", id)

	dataValue, errs := converter.ConvertToVapi(obj, model.LBTcpMonitorProfileBindingType())
	if errs != nil {
		return fmt.Errorf("VPC LBMonitorProfile %s is not of type LBTcpMonitorProfile %s", id, errs[0])
	}

	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbMonitorProfilesClient(sessionContext, connector)
	return c.Patch(parents[0], parents[1], parents[2], id, dataValue.(*data.StructValue))
}

func resourceNsxtVpcLBTcpMonitorProfileCreate(d *schema.ResourceData, m interface{}) error {
	if !util.NsxVersionHigherOrEqual("9.0.0") {
		return fmt.Errorf("VPC LB TCP Monitor Profile resource requires NSX version 9.0.0 or higher")
	}

	// Initialize resource Id and verify this ID is not yet used
	id, err := getOrGenerateID2(d, m, resourceNsxtVpcLBMonitorProfileExists)
	if err != nil {
		return err
	}

	err = resourceNsxtVpcLBTcpMonitorProfilePatch(d, m, id)
	if err != nil {
		return handleCreateError("VPC LBTcpMonitorProfile", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtVpcLBTcpMonitorProfileRead(d, m)
}

func resourceNsxtVpcLBTcpMonitorProfileRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VPC LBTcpMonitorProfile ID")
	}

	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbMonitorProfilesClient(sessionContext, connector)
	obj, err := c.Get(parents[0], parents[1], parents[2], id)
	if err != nil {
		return handleReadError(d, "VPC LBTcpMonitorProfile", id, err)
	}

	converter := bindings.NewTypeConverter()
	baseObj, errs := converter.ConvertToGolang(obj, model.LBTcpMonitorProfileBindingType())
	if len(errs) > 0 {
		return fmt.Errorf("Error converting VPC LBTcpMonitorProfile %s", errs[0])
	}
	lbTCPMonitor := baseObj.(model.LBTcpMonitorProfile)

	d.Set("display_name", lbTCPMonitor.DisplayName)
	d.Set("description", lbTCPMonitor.Description)
	setPolicyTagsInSchema(d, lbTCPMonitor.Tags)
	d.Set("nsx_id", id)
	d.Set("path", lbTCPMonitor.Path)
	d.Set("revision", lbTCPMonitor.Revision)

	d.Set("receive", lbTCPMonitor.Receive)
	d.Set("send", lbTCPMonitor.Send)
	d.Set("fall_count", lbTCPMonitor.FallCount)
	d.Set("interval", lbTCPMonitor.Interval)
	d.Set("monitor_port", lbTCPMonitor.MonitorPort)
	d.Set("rise_count", lbTCPMonitor.RiseCount)
	d.Set("timeout", lbTCPMonitor.Timeout)

	return nil
}

func resourceNsxtVpcLBTcpMonitorProfileUpdate(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VPC LBTcpMonitorProfile ID")
	}

	err := resourceNsxtVpcLBTcpMonitorProfilePatch(d, m, id)
	if err != nil {
		return handleUpdateError("VPC LBTcpMonitorProfile", id, err)
	}

	return resourceNsxtVpcLBTcpMonitorProfileRead(d, m)
}

func resourceNsxtVpcLBTcpMonitorProfileDelete(d *schema.ResourceData, m interface{}) error {
	return resourceNsxtVpcLBMonitorProfileDelete(d, m)
}
