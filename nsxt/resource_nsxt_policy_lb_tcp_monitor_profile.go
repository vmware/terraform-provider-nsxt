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
)

func resourceNsxtPolicyLBTcpMonitorProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyLBTcpMonitorProfileCreate,
		Read:   resourceNsxtPolicyLBTcpMonitorProfileRead,
		Update: resourceNsxtPolicyLBTcpMonitorProfileUpdate,
		Delete: resourceNsxtPolicyLBTcpMonitorProfileDelete,
		Importer: &schema.ResourceImporter{
			State: getPolicyPathOrIDResourceImporter(lbMonitorProfilePathExample),
		},

		Schema: map[string]*schema.Schema{
			"nsx_id":       getNsxIDSchema(),
			"path":         getPathSchema(),
			"display_name": getDisplayNameSchema(),
			"description":  getDescriptionSchema(),
			"revision":     getRevisionSchema(),
			"tag":          getTagsSchema(),
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

func resourceNsxtPolicyLBTcpMonitorProfilePatch(d *schema.ResourceData, m interface{}, id string) error {
	connector := getPolicyConnector(m)
	converter := bindings.NewTypeConverter()

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	var tags []model.Tag
	if isConfigScopedCacheMode() {
		tags = getPolicyTagsWithProviderManagedDefaults(d, m)
	} else {
		tags = getPolicyTagsFromSchema(d)
	}
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

	log.Printf("[INFO] Patching LBTcpMonitorProfile with ID %s", id)

	dataValue, errs := converter.ConvertToVapi(obj, model.LBTcpMonitorProfileBindingType())
	if errs != nil {
		return fmt.Errorf("LBMonitorProfile %s is not of type LBTcpMonitorProfile %s", id, errs[0])
	}

	sessionContext := getSessionContext(d, m)
	client := cliLbMonitorProfilesClient(sessionContext, connector)
	return client.Patch(id, dataValue.(*data.StructValue))
}

func resourceNsxtPolicyLBTcpMonitorProfileCreate(d *schema.ResourceData, m interface{}) error {

	// Initialize resource Id and verify this ID is not yet used
	id, err := getOrGenerateID(d, m, resourceNsxtPolicyLBMonitorProfileExistsWrapper)
	if err != nil {
		return err
	}

	err = resourceNsxtPolicyLBTcpMonitorProfilePatch(d, m, id)
	if err != nil {
		return handleCreateError("LBTcpMonitorProfile", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)
	MarkPostWriteAndInvalidateCacheForResourceType("LBTcpMonitorProfile", d)

	return resourceNsxtPolicyLBTcpMonitorProfileRead(d, m)
}

func resourceNsxtPolicyLBTcpMonitorProfileRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining LBTcpMonitorProfile ID")
	}

	sessionContext := getSessionContext(d, m)
	client := cliLbMonitorProfilesClient(sessionContext, connector)
	var lbTCPMonitor *model.LBTcpMonitorProfile
	var err error
	if isCacheEnabledForRead(d) {
		lbTCPMonitor, _, _, err = CacheAwareResourceRead[model.LBTcpMonitorProfile](
			d,
			m,
			connector,
			id,
			"LBTcpMonitorProfile",
			model.LBTcpMonitorProfileBindingType(),
			func() (*model.LBTcpMonitorProfile, error) {
				converter := bindings.NewTypeConverter()
				obj, readErr := client.Get(id)
				if readErr != nil {
					return nil, readErr
				}
				baseObj, errs := converter.ConvertToGolang(obj, model.LBTcpMonitorProfileBindingType())
				if len(errs) > 0 {
					return nil, fmt.Errorf("Error converting LBTcpMonitorProfile %s", errs[0])
				}
				typed := baseObj.(model.LBTcpMonitorProfile)
				return &typed, nil
			},
			func(patchObj *model.LBTcpMonitorProfile) error {
				return resourceNsxtPolicyLBTcpMonitorProfilePatch(d, m, id)
			},
		)
	} else {
		converter := bindings.NewTypeConverter()
		obj, readErr := client.Get(id)
		if readErr != nil {
			err = readErr
		} else {
			baseObj, errs := converter.ConvertToGolang(obj, model.LBTcpMonitorProfileBindingType())
			if len(errs) > 0 {
				err = fmt.Errorf("Error converting LBTcpMonitorProfile %s", errs[0])
			} else {
				typed := baseObj.(model.LBTcpMonitorProfile)
				lbTCPMonitor = &typed
			}
		}
	}
	if err != nil {
		return handleReadError(d, "LBHTcpMonitorProfile", id, err)
	}

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

func resourceNsxtPolicyLBTcpMonitorProfileUpdate(d *schema.ResourceData, m interface{}) error {

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining LBTcpMonitorProfile ID")
	}

	err := resourceNsxtPolicyLBTcpMonitorProfilePatch(d, m, id)
	if err != nil {
		return handleUpdateError("LBTcpMonitorProfile", id, err)
	}
	MarkPostWriteAndInvalidateCacheForResourceType("LBTcpMonitorProfile", d)

	return resourceNsxtPolicyLBTcpMonitorProfileRead(d, m)
}

func resourceNsxtPolicyLBTcpMonitorProfileDelete(d *schema.ResourceData, m interface{}) error {
	return resourceNsxtPolicyLBMonitorProfileDelete(d, m)
}
