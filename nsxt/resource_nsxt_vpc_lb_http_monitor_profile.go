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
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var vpcLbMonitorProfilePathExample = "/orgs/[org]/projects/[project]/vpcs/[vpc]/vpc-lb-monitor-profiles/[profile]"

func resourceNsxtVpcLBHttpMonitorProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtVpcLBHttpMonitorProfileCreate,
		Read:   resourceNsxtVpcLBHttpMonitorProfileRead,
		Update: resourceNsxtVpcLBHttpMonitorProfileUpdate,
		Delete: resourceNsxtVpcLBHttpMonitorProfileDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtVersionCheckImporter("9.0.0", "VPC LB HTTP Monitor Profile", getVpcPathResourceImporter(vpcLbMonitorProfilePathExample)),
		},

		Schema: map[string]*schema.Schema{
			"nsx_id":                getNsxIDSchema(),
			"path":                  getPathSchema(),
			"display_name":          getDisplayNameSchema(),
			"description":           getDescriptionSchema(),
			"revision":              getRevisionSchema(),
			"tag":                   getTagsSchema(),
			"context":               getContextSchemaExtended(true, false, true, true),
			"request_body":          getLbMonitorRequestBodySchema(),
			"request_header":        getLbHTTPHeaderSchema("Array of HTTP request headers"),
			"request_method":        getLbMonitorRequestMethodSchema(),
			"request_url":           getLbMonitorRequestURLSchema(),
			"request_version":       getLbMonitorRequestVersionSchema(),
			"response_body":         getLbMonitorResponseBodySchema(),
			"response_status_codes": getLbMonitorResponseStatusCodesSchema(),
			"fall_count":            getLbMonitorFallCountSchema(),
			"interval":              getLbMonitorIntervalSchema(),
			"rise_count":            getLbMonitorRiseCountSchema(),
			"timeout":               getLbMonitorTimeoutSchema(),
			"monitor_port":          getPolicyLbMonitorPortSchema(),
		},
	}
}

func resourceNsxtVpcLBHttpMonitorProfilePatch(d *schema.ResourceData, m interface{}, id string) error {
	connector := getPolicyConnector(m)
	converter := bindings.NewTypeConverter()

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsWithProviderManagedDefaults(d, m)
	requestBody := d.Get("request_body").(string)
	requestHeaders := getPolicyLbHTTPHeaderFromSchema(d, "request_header")
	requestMethod := d.Get("request_method").(string)
	requestURL := d.Get("request_url").(string)
	requestVersion := d.Get("request_version").(string)
	responseBody := d.Get("response_body").(string)
	responseStatusCodes := interface2Int64List(d.Get("response_status_codes").([]interface{}))
	fallCount := int64(d.Get("fall_count").(int))
	interval := int64(d.Get("interval").(int))
	riseCount := int64(d.Get("rise_count").(int))
	timeout := int64(d.Get("timeout").(int))
	monitorPort := int64(d.Get("monitor_port").(int))

	resourceType := model.LBMonitorProfile_RESOURCE_TYPE_LBHTTPMONITORPROFILE

	obj := model.LBHttpMonitorProfile{
		DisplayName:         &displayName,
		Description:         &description,
		Tags:                tags,
		RequestBody:         &requestBody,
		RequestHeaders:      requestHeaders,
		RequestMethod:       &requestMethod,
		RequestUrl:          &requestURL,
		RequestVersion:      &requestVersion,
		ResponseBody:        &responseBody,
		ResponseStatusCodes: responseStatusCodes,
		FallCount:           &fallCount,
		Interval:            &interval,
		MonitorPort:         &monitorPort,
		RiseCount:           &riseCount,
		Timeout:             &timeout,
		ResourceType:        resourceType,
	}

	log.Printf("[INFO] Patching VPC LBHttpMonitorProfile with ID %s", id)
	dataValue, errs := converter.ConvertToVapi(obj, model.LBHttpMonitorProfileBindingType())
	if errs != nil {
		return fmt.Errorf("VPC LBMonitorProfile %s is not of type LBHttpMonitorProfile %s", id, errs[0])
	}

	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbMonitorProfilesClient(sessionContext, connector)
	return c.Patch(parents[0], parents[1], parents[2], id, dataValue.(*data.StructValue))
}

func resourceNsxtVpcLBHttpMonitorProfileCreate(d *schema.ResourceData, m interface{}) error {
	if !util.NsxVersionHigherOrEqual("9.0.0") {
		return fmt.Errorf("VPC LB HTTP Monitor Profile resource requires NSX version 9.0.0 or higher")
	}

	// Initialize resource Id and verify this ID is not yet used
	id, err := getOrGenerateID2(d, m, resourceNsxtVpcLBMonitorProfileExists)
	if err != nil {
		return err
	}

	err = resourceNsxtVpcLBHttpMonitorProfilePatch(d, m, id)
	if err != nil {
		return handleCreateError("VPC LBHttpMonitorProfile", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtVpcLBHttpMonitorProfileRead(d, m)
}

func resourceNsxtVpcLBHttpMonitorProfileRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	converter := bindings.NewTypeConverter()

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VPC LBHttpMonitorProfile ID")
	}

	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbMonitorProfilesClient(sessionContext, connector)
	obj, err := c.Get(parents[0], parents[1], parents[2], id)
	if err != nil {
		return handleReadError(d, "VPC LBHttpMonitorProfile", id, err)
	}

	baseObj, errs := converter.ConvertToGolang(obj, model.LBHttpMonitorProfileBindingType())
	if len(errs) > 0 {
		return fmt.Errorf("Error converting VPC LBHttpMonitorProfile %s", errs[0])
	}

	lbHTTPMonitor := baseObj.(model.LBHttpMonitorProfile)

	d.Set("revision", lbHTTPMonitor.Revision)
	d.Set("description", lbHTTPMonitor.Description)
	d.Set("display_name", lbHTTPMonitor.DisplayName)
	setPolicyTagsInSchema(d, lbHTTPMonitor.Tags)
	d.Set("fall_count", lbHTTPMonitor.FallCount)
	d.Set("interval", lbHTTPMonitor.Interval)
	d.Set("monitor_port", lbHTTPMonitor.MonitorPort)
	d.Set("rise_count", lbHTTPMonitor.RiseCount)
	d.Set("timeout", lbHTTPMonitor.Timeout)
	d.Set("request_body", lbHTTPMonitor.RequestBody)
	d.Set("path", lbHTTPMonitor.Path)
	d.Set("nsx_id", id)
	setPolicyLbHTTPHeaderInSchema(d, "request_header", lbHTTPMonitor.RequestHeaders)
	d.Set("request_method", lbHTTPMonitor.RequestMethod)
	d.Set("request_url", lbHTTPMonitor.RequestUrl)
	d.Set("request_version", lbHTTPMonitor.RequestVersion)
	d.Set("response_body", lbHTTPMonitor.ResponseBody)
	d.Set("response_status_codes", int64List2Interface(lbHTTPMonitor.ResponseStatusCodes))

	return nil
}

func resourceNsxtVpcLBHttpMonitorProfileUpdate(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VPC LBHttpMonitorProfile ID")
	}

	err := resourceNsxtVpcLBHttpMonitorProfilePatch(d, m, id)
	if err != nil {
		return handleUpdateError("VPC LBHttpMonitorProfile", id, err)
	}

	return resourceNsxtVpcLBHttpMonitorProfileRead(d, m)
}

func resourceNsxtVpcLBHttpMonitorProfileDelete(d *schema.ResourceData, m interface{}) error {
	return resourceNsxtVpcLBMonitorProfileDelete(d, m)
}

// resourceNsxtVpcLBMonitorProfileExists is shared by all VPC LB monitor profile resources
func resourceNsxtVpcLBMonitorProfileExists(sessionContext utl.SessionContext, id string, connector client.Connector) (bool, error) {
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbMonitorProfilesClient(sessionContext, connector)
	_, err := c.Get(parents[0], parents[1], parents[2], id)
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving resource VPC LBMonitorProfile", err)
}

// resourceNsxtVpcLBMonitorProfileDelete is shared by all VPC LB monitor profile resources
func resourceNsxtVpcLBMonitorProfileDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VPC LBMonitorProfile ID")
	}

	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbMonitorProfilesClient(sessionContext, connector)

	force := true
	err := c.Delete(parents[0], parents[1], parents[2], id, &force)
	if err != nil {
		return handleDeleteError("VPC LBMonitorProfile", id, err)
	}
	return nil
}
