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

func resourceNsxtVpcLBHttpsMonitorProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtVpcLBHttpsMonitorProfileCreate,
		Read:   resourceNsxtVpcLBHttpsMonitorProfileRead,
		Update: resourceNsxtVpcLBHttpsMonitorProfileUpdate,
		Delete: resourceNsxtVpcLBHttpsMonitorProfileDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtVersionCheckImporter("9.0.0", "VPC LB HTTPS Monitor Profile", getVpcPathResourceImporter(vpcLbMonitorProfilePathExample)),
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
			"server_ssl":            getLbServerSslSchema(),
			"fall_count":            getLbMonitorFallCountSchema(),
			"interval":              getLbMonitorIntervalSchema(),
			"monitor_port":          getPolicyLbMonitorPortSchema(),
			"rise_count":            getLbMonitorRiseCountSchema(),
			"timeout":               getLbMonitorTimeoutSchema(),
		},
	}
}

func resourceNsxtVpcLBHttpsMonitorProfilePatch(d *schema.ResourceData, m interface{}, id string) error {
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
	serverSsl := getLbServerSslFromSchema(d)
	fallCount := int64(d.Get("fall_count").(int))
	interval := int64(d.Get("interval").(int))
	monitorPort := int64(d.Get("monitor_port").(int))
	riseCount := int64(d.Get("rise_count").(int))
	timeout := int64(d.Get("timeout").(int))

	resourceType := model.LBMonitorProfile_RESOURCE_TYPE_LBHTTPSMONITORPROFILE

	obj := model.LBHttpsMonitorProfile{
		DisplayName:             &displayName,
		Description:             &description,
		Tags:                    tags,
		RequestBody:             &requestBody,
		RequestHeaders:          requestHeaders,
		RequestMethod:           &requestMethod,
		RequestUrl:              &requestURL,
		RequestVersion:          &requestVersion,
		ResponseBody:            &responseBody,
		ResponseStatusCodes:     responseStatusCodes,
		ServerSslProfileBinding: serverSsl,
		FallCount:               &fallCount,
		Interval:                &interval,
		MonitorPort:             &monitorPort,
		RiseCount:               &riseCount,
		Timeout:                 &timeout,
		ResourceType:            resourceType,
	}

	log.Printf("[INFO] Patching VPC LBHttpsMonitorProfile with ID %s", id)
	dataValue, errs := converter.ConvertToVapi(obj, model.LBHttpsMonitorProfileBindingType())
	if errs != nil {
		return fmt.Errorf("VPC LBMonitorProfile %s is not of type LBHttpsMonitorProfile %s", id, errs[0])
	}

	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbMonitorProfilesClient(sessionContext, connector)
	return c.Patch(parents[0], parents[1], parents[2], id, dataValue.(*data.StructValue))
}

func resourceNsxtVpcLBHttpsMonitorProfileCreate(d *schema.ResourceData, m interface{}) error {
	if !util.NsxVersionHigherOrEqual("9.0.0") {
		return fmt.Errorf("VPC LB HTTPS Monitor Profile resource requires NSX version 9.0.0 or higher")
	}

	// Initialize resource Id and verify this ID is not yet used
	id, err := getOrGenerateID2(d, m, resourceNsxtVpcLBMonitorProfileExists)
	if err != nil {
		return err
	}

	err = resourceNsxtVpcLBHttpsMonitorProfilePatch(d, m, id)
	if err != nil {
		return handleCreateError("VPC LBHttpsMonitorProfile", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtVpcLBHttpsMonitorProfileRead(d, m)
}

func resourceNsxtVpcLBHttpsMonitorProfileRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	converter := bindings.NewTypeConverter()

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VPC LBHttpsMonitorProfile ID")
	}

	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbMonitorProfilesClient(sessionContext, connector)
	obj, err := c.Get(parents[0], parents[1], parents[2], id)
	if err != nil {
		return handleReadError(d, "VPC LBHttpsMonitorProfile", id, err)
	}

	baseObj, errs := converter.ConvertToGolang(obj, model.LBHttpsMonitorProfileBindingType())
	if len(errs) > 0 {
		return fmt.Errorf("Error converting VPC LBHttpsMonitorProfile %s", errs[0])
	}
	lbHTTPSMonitor := baseObj.(model.LBHttpsMonitorProfile)

	d.Set("revision", lbHTTPSMonitor.Revision)
	d.Set("description", lbHTTPSMonitor.Description)
	d.Set("display_name", lbHTTPSMonitor.DisplayName)
	setPolicyTagsInSchema(d, lbHTTPSMonitor.Tags)
	d.Set("fall_count", lbHTTPSMonitor.FallCount)
	d.Set("interval", lbHTTPSMonitor.Interval)
	d.Set("monitor_port", lbHTTPSMonitor.MonitorPort)
	d.Set("rise_count", lbHTTPSMonitor.RiseCount)
	d.Set("timeout", lbHTTPSMonitor.Timeout)
	d.Set("request_body", lbHTTPSMonitor.RequestBody)
	d.Set("path", lbHTTPSMonitor.Path)
	d.Set("nsx_id", id)
	setPolicyLbHTTPHeaderInSchema(d, "request_header", lbHTTPSMonitor.RequestHeaders)
	d.Set("request_method", lbHTTPSMonitor.RequestMethod)
	d.Set("request_url", lbHTTPSMonitor.RequestUrl)
	d.Set("request_version", lbHTTPSMonitor.RequestVersion)
	d.Set("response_body", lbHTTPSMonitor.ResponseBody)
	d.Set("response_status_codes", int64List2Interface(lbHTTPSMonitor.ResponseStatusCodes))
	if lbHTTPSMonitor.ServerSslProfileBinding != nil {
		setLbServerSslInSchema(d, *lbHTTPSMonitor.ServerSslProfileBinding)
	} else if err := d.Set("server_ssl", nil); err != nil {
		return err
	}
	return nil
}

func resourceNsxtVpcLBHttpsMonitorProfileUpdate(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VPC LBHttpsMonitorProfile ID")
	}

	err := resourceNsxtVpcLBHttpsMonitorProfilePatch(d, m, id)
	if err != nil {
		return handleUpdateError("VPC LBHttpsMonitorProfile", id, err)
	}

	return resourceNsxtVpcLBHttpsMonitorProfileRead(d, m)
}

func resourceNsxtVpcLBHttpsMonitorProfileDelete(d *schema.ResourceData, m interface{}) error {
	return resourceNsxtVpcLBMonitorProfileDelete(d, m)
}
