// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	"github.com/vmware/terraform-provider-nsxt/api/orgs/projects/vpcs"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var cliVpcLbsClient = vpcs.NewVpcLbsClient

var vpcLbServicePathExample = "/orgs/[org]/projects/[project]/vpcs/[vpc]/vpc-lbs/[service]"

func resourceNsxtVpcLBService() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtVpcLBServiceCreate,
		Read:   resourceNsxtVpcLBServiceRead,
		Update: resourceNsxtVpcLBServiceUpdate,
		Delete: resourceNsxtVpcLBServiceDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtVersionCheckImporter("9.0.0", "VPC LB Service", getVpcPathResourceImporter(vpcLbServicePathExample)),
		},

		Schema: map[string]*schema.Schema{
			"nsx_id":       getNsxIDSchema(),
			"path":         getPathSchema(),
			"display_name": getDisplayNameSchema(),
			"description":  getDescriptionSchema(),
			"revision":     getRevisionSchema(),
			"tag":          getTagsSchema(),
			"context":      getContextSchemaExtended(true, false, true, true),
			"connectivity_path": {
				Type:        schema.TypeString,
				Description: "Policy path of the VPC which the LB service is connected to",
				Computed:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Flag to enable the Service",
				Optional:    true,
				Default:     true,
			},
			"error_log_level": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(lBServiceErrorLogLevelValues, false),
				Description:  "Log level for Load Balancer Service messages",
				Optional:     true,
				Default:      model.LBService_ERROR_LOG_LEVEL_INFO,
			},
			"size": {
				Type:         schema.TypeString,
				Description:  "Load Balancer Service size",
				ValidateFunc: validation.StringInSlice(lBServiceSizeValues, false),
				Optional:     true,
				Default:      model.LBService_SIZE_SMALL,
			},
		},
	}
}

func resourceNsxtVpcLBServiceExists(sessionContext utl.SessionContext, id string, connector client.Connector) (bool, error) {
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbsClient(sessionContext, connector)
	_, err := c.Get(parents[0], parents[1], parents[2], id)
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving resource", err)
}

func resourceNsxtVpcLBServiceCreate(d *schema.ResourceData, m interface{}) error {
	if !util.NsxVersionHigherOrEqual("9.0.0") {
		return fmt.Errorf("VPC LB Service resource requires NSX version 9.0.0 or higher")
	}
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbsClient(sessionContext, connector)

	// Initialize resource Id and verify this ID is not yet used
	id, err := getOrGenerateID2(d, m, resourceNsxtVpcLBServiceExists)
	if err != nil {
		return err
	}

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsWithProviderManagedDefaults(d, m)
	enabled := d.Get("enabled").(bool)
	errorLogLevel := d.Get("error_log_level").(string)
	size := d.Get("size").(string)

	obj := model.LBService{
		DisplayName:   &displayName,
		Description:   &description,
		Tags:          tags,
		Enabled:       &enabled,
		ErrorLogLevel: &errorLogLevel,
		Size:          &size,
	}

	// Create the resource using PATCH
	log.Printf("[INFO] Creating VPC LBService with ID %s", id)
	err = c.Patch(parents[0], parents[1], parents[2], id, obj, nil)
	if err != nil {
		return handleCreateError("VPC LBService", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)
	return resourceNsxtVpcLBServiceRead(d, m)
}

func resourceNsxtVpcLBServiceRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbsClient(sessionContext, connector)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VPC LBService ID")
	}

	obj, err := c.Get(parents[0], parents[1], parents[2], id)
	if err != nil {
		return handleReadError(d, "VPC LBService", id, err)
	}

	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("path", obj.Path)
	d.Set("revision", obj.Revision)

	d.Set("connectivity_path", obj.ConnectivityPath)
	d.Set("enabled", obj.Enabled)
	d.Set("error_log_level", obj.ErrorLogLevel)
	d.Set("size", obj.Size)

	return nil
}

func resourceNsxtVpcLBServiceUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbsClient(sessionContext, connector)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VPC LBService ID")
	}

	// Read the rest of the configured parameters
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getPolicyTagsWithProviderManagedDefaults(d, m)
	enabled := d.Get("enabled").(bool)
	errorLogLevel := d.Get("error_log_level").(string)
	size := d.Get("size").(string)

	obj := model.LBService{
		DisplayName:   &displayName,
		Description:   &description,
		Tags:          tags,
		Enabled:       &enabled,
		ErrorLogLevel: &errorLogLevel,
		Size:          &size,
	}

	// Update the resource using PATCH
	err := c.Patch(parents[0], parents[1], parents[2], id, obj, nil)
	if err != nil {
		return handleUpdateError("VPC LBService", id, err)
	}
	return resourceNsxtVpcLBServiceRead(d, m)
}

func resourceNsxtVpcLBServiceDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VPC LBService ID")
	}

	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbsClient(sessionContext, connector)

	force := true
	err := c.Delete(parents[0], parents[1], parents[2], id, &force)
	if err != nil {
		return handleDeleteError("VPC LBService", id, err)
	}

	return nil
}
