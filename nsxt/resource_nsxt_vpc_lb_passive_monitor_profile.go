// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func resourceNsxtVpcLBPassiveMonitorProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtVpcLBPassiveMonitorProfileCreate,
		Read:   resourceNsxtVpcLBPassiveMonitorProfileRead,
		Update: resourceNsxtVpcLBPassiveMonitorProfileUpdate,
		Delete: resourceNsxtVpcLBPassiveMonitorProfileDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtVersionCheckImporter("9.0.0", "VPC LB Passive Monitor Profile", getVpcPathResourceImporter(vpcLbMonitorProfilePathExample)),
		},

		Schema: map[string]*schema.Schema{
			"nsx_id":       getNsxIDSchema(),
			"path":         getPathSchema(),
			"display_name": getDisplayNameSchema(),
			"description":  getDescriptionSchema(),
			"revision":     getRevisionSchema(),
			"tag":          getTagsSchema(),
			"context":      getContextSchemaExtended(true, false, true, true),
			"max_fails": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "Number of consecutive failures before a member is considered temporarily unavailable",
				Default:      5,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"timeout": getLbMonitorTimeoutSchema(),
		},
	}
}

func resourceNsxtVpcLBPassiveMonitorProfilePatch(d *schema.ResourceData, m interface{}, id string) error {
	connector := getPolicyConnector(m)
	converter := bindings.NewTypeConverter()

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsWithProviderManagedDefaults(d, m)
	maxFails := int64(d.Get("max_fails").(int))
	timeout := int64(d.Get("timeout").(int))
	resourceType := model.LBMonitorProfile_RESOURCE_TYPE_LBPASSIVEMONITORPROFILE

	obj := model.LBPassiveMonitorProfile{
		DisplayName:  &displayName,
		Description:  &description,
		Tags:         tags,
		MaxFails:     &maxFails,
		Timeout:      &timeout,
		ResourceType: resourceType,
	}

	log.Printf("[INFO] Patching VPC LBPassiveMonitorProfile with ID %s", id)

	dataValue, errs := converter.ConvertToVapi(obj, model.LBPassiveMonitorProfileBindingType())
	if errs != nil {
		return fmt.Errorf("VPC LBMonitorProfile %s is not of type LBPassiveMonitorProfile %s", id, errs[0])
	}

	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbMonitorProfilesClient(sessionContext, connector)
	return c.Patch(parents[0], parents[1], parents[2], id, dataValue.(*data.StructValue))
}

func resourceNsxtVpcLBPassiveMonitorProfileCreate(d *schema.ResourceData, m interface{}) error {
	if !util.NsxVersionHigherOrEqual("9.0.0") {
		return fmt.Errorf("VPC LB Passive Monitor Profile resource requires NSX version 9.0.0 or higher")
	}

	// Initialize resource Id and verify this ID is not yet used
	id, err := getOrGenerateID2(d, m, resourceNsxtVpcLBMonitorProfileExists)
	if err != nil {
		return err
	}

	err = resourceNsxtVpcLBPassiveMonitorProfilePatch(d, m, id)
	if err != nil {
		return handleCreateError("VPC LBPassiveMonitorProfile", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtVpcLBPassiveMonitorProfileRead(d, m)
}

func resourceNsxtVpcLBPassiveMonitorProfileRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	converter := bindings.NewTypeConverter()

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VPC LBPassiveMonitorProfile ID")
	}

	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbMonitorProfilesClient(sessionContext, connector)
	obj, err := c.Get(parents[0], parents[1], parents[2], id)
	if err != nil {
		return handleReadError(d, "VPC LBPassiveMonitorProfile", id, err)
	}

	baseObj, errs := converter.ConvertToGolang(obj, model.LBPassiveMonitorProfileBindingType())
	if len(errs) > 0 {
		return fmt.Errorf("Error converting VPC LBPassiveMonitorProfile %s", errs[0])
	}
	lbPassiveMonitor := baseObj.(model.LBPassiveMonitorProfile)

	d.Set("display_name", lbPassiveMonitor.DisplayName)
	d.Set("description", lbPassiveMonitor.Description)
	setPolicyTagsInSchema(d, lbPassiveMonitor.Tags)
	d.Set("nsx_id", id)
	d.Set("path", lbPassiveMonitor.Path)
	d.Set("revision", lbPassiveMonitor.Revision)
	d.Set("max_fails", lbPassiveMonitor.MaxFails)
	d.Set("timeout", lbPassiveMonitor.Timeout)

	return nil
}

func resourceNsxtVpcLBPassiveMonitorProfileUpdate(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VPC LBPassiveMonitorProfile ID")
	}

	err := resourceNsxtVpcLBPassiveMonitorProfilePatch(d, m, id)
	if err != nil {
		return handleUpdateError("VPC LBPassiveMonitorProfile", id, err)
	}

	return resourceNsxtVpcLBPassiveMonitorProfileRead(d, m)
}

func resourceNsxtVpcLBPassiveMonitorProfileDelete(d *schema.ResourceData, m interface{}) error {
	return resourceNsxtVpcLBMonitorProfileDelete(d, m)
}
