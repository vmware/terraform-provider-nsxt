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
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var vpcLbAppProfilePathExample = "/orgs/[org]/projects/[project]/vpcs/[vpc]/vpc-lb-app-profiles/[profile]"

func resourceNsxtVpcLBHttpApplicationProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtVpcLBHttpApplicationProfileCreate,
		Read:   resourceNsxtVpcLBHttpApplicationProfileRead,
		Update: resourceNsxtVpcLBHttpApplicationProfileUpdate,
		Delete: resourceNsxtVpcLBHttpApplicationProfileDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtVersionCheckImporter("9.0.0", "VPC LB HTTP Application Profile", getVpcPathResourceImporter(vpcLbAppProfilePathExample)),
		},

		Schema: map[string]*schema.Schema{
			"nsx_id":       getNsxIDSchema(),
			"path":         getPathSchema(),
			"display_name": getDisplayNameSchema(),
			"description":  getDescriptionSchema(),
			"revision":     getRevisionSchema(),
			"tag":          getTagsSchema(),
			"context":      getContextSchemaExtended(true, false, true, true),
			"http_redirect_to": {
				Type:        schema.TypeString,
				Description: "A URL that incoming requests for that virtual server can be temporarily redirected to, If a website is temporarily down or has moved",
				Optional:    true,
			},
			"http_redirect_to_https": {
				Type:        schema.TypeBool,
				Description: "A boolean flag which reflects whether the client will automatically be redirected to use SSL",
				Optional:    true,
				Default:     false,
			},
			"idle_timeout": {
				Type:         schema.TypeInt,
				Description:  "Timeout in seconds to specify how long an HTTP application can remain idle",
				Optional:     true,
				Default:      15,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"request_body_size": {
				Type:        schema.TypeInt,
				Description: "Maximum request body size in bytes (Unlimited if not specified)",
				Optional:    true,
			},
			"request_header_size": {
				Type:         schema.TypeInt,
				Description:  "Maximum request header size in bytes. Requests with larger header size will be processed as best effort whereas a request with header below this specified size is guaranteed to be processed",
				Optional:     true,
				Default:      1024,
				ValidateFunc: validation.IntBetween(1, 65536),
			},
			"response_buffering": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "A boolean flag indicating whether the response received by LB from the backend will be saved into the buffers",
				Default:     false,
			},
			"response_header_size": {
				Type:         schema.TypeInt,
				Description:  "Maximum request header size in bytes. Requests with larger header size will be processed as best effort whereas a request with header below this specified size is guaranteed to be processed",
				Optional:     true,
				Default:      4096,
				ValidateFunc: validation.IntBetween(1, 65536),
			},
			"response_timeout": {
				Type:         schema.TypeInt,
				Description:  "Number of seconds waiting for the server response before the connection is closed",
				Optional:     true,
				Default:      60,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"server_keep_alive": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "A boolean flag indicating whether the backend connection will be kept alive for client connection. If server_keep_alive is true, it means the backend connection will keep alive for the client connection. Every client connection is tied 1:1 with the corresponding server-side connection. If server_keep_alive is false, it means the backend connection won't keep alive for the client connection.",
				Default:     false,
			},
			"x_forwarded_for": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(lBHttpProfileXForwardedForValues, false),
				Optional:     true,
				Description:  "When X-Forwarded-For is configured, X-Forwarded-Proto and X-Forwarded-Port information is added automatically into request header",
			},
		},
	}
}

func resourceNsxtVpcLBHttpApplicationProfilePatch(d *schema.ResourceData, m interface{}, id string) error {
	connector := getPolicyConnector(m)
	converter := bindings.NewTypeConverter()

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsWithProviderManagedDefaults(d, m)
	httpRedirectTo := d.Get("http_redirect_to").(string)
	httpRedirectToHTTPS := d.Get("http_redirect_to_https").(bool)
	idleTimeout := int64(d.Get("idle_timeout").(int))
	requestBodySize := int64(d.Get("request_body_size").(int))
	requestHeaderSize := int64(d.Get("request_header_size").(int))
	responseBuffering := d.Get("response_buffering").(bool)
	responseHeaderSize := int64(d.Get("response_header_size").(int))
	responseTimeout := int64(d.Get("response_timeout").(int))
	serverKeepAlive := d.Get("server_keep_alive").(bool)
	xForwardedFor := d.Get("x_forwarded_for").(string)
	resourceType := model.LBAppProfile_RESOURCE_TYPE_LBHTTPPROFILE
	obj := model.LBHttpProfile{
		DisplayName:         &displayName,
		Description:         &description,
		Tags:                tags,
		HttpRedirectToHttps: &httpRedirectToHTTPS,
		IdleTimeout:         &idleTimeout,
		RequestHeaderSize:   &requestHeaderSize,
		ResponseBuffering:   &responseBuffering,
		ResponseHeaderSize:  &responseHeaderSize,
		ResponseTimeout:     &responseTimeout,
		ServerKeepAlive:     &serverKeepAlive,
		ResourceType:        resourceType,
	}
	if requestBodySize > 0 {
		obj.RequestBodySize = &requestBodySize
	}
	if len(xForwardedFor) > 0 {
		obj.XForwardedFor = &xForwardedFor
	}
	if len(httpRedirectTo) > 0 {
		obj.HttpRedirectTo = &httpRedirectTo
	}

	log.Printf("[INFO] Patching VPC LBHttpProfile with ID %s", id)
	dataValue, errs := converter.ConvertToVapi(obj, model.LBHttpProfileBindingType())
	if errs != nil {
		return fmt.Errorf("Error converting VPC LBHttpProfile %s", errs[0])
	}

	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbAppProfilesClient(sessionContext, connector)
	return c.Patch(parents[0], parents[1], parents[2], id, dataValue.(*data.StructValue))
}

func resourceNsxtVpcLBHttpApplicationProfileCreate(d *schema.ResourceData, m interface{}) error {
	if !util.NsxVersionHigherOrEqual("9.0.0") {
		return fmt.Errorf("VPC LB HTTP Application Profile resource requires NSX version 9.0.0 or higher")
	}

	// Initialize resource Id and verify this ID is not yet used
	id, err := getOrGenerateID2(d, m, resourceNsxtVpcLBAppProfileExists)
	if err != nil {
		return err
	}

	err = resourceNsxtVpcLBHttpApplicationProfilePatch(d, m, id)
	if err != nil {
		return handleCreateError("VPC LBHttpProfile", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtVpcLBHttpApplicationProfileRead(d, m)
}

func resourceNsxtVpcLBHttpApplicationProfileRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	converter := bindings.NewTypeConverter()

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VPC LBHttpProfile ID")
	}

	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbAppProfilesClient(sessionContext, connector)
	obj, err := c.Get(parents[0], parents[1], parents[2], id)
	if err != nil {
		return handleReadError(d, "VPC LBHttpProfile", id, err)
	}

	baseObj, errs := converter.ConvertToGolang(obj, model.LBHttpProfileBindingType())
	if len(errs) > 0 {
		return fmt.Errorf("VPC LBAppProfile with id %s is not of type LBHttpProfile %s", id, errs[0])
	}
	lbHTTPProfile := baseObj.(model.LBHttpProfile)

	d.Set("display_name", lbHTTPProfile.DisplayName)
	d.Set("description", lbHTTPProfile.Description)
	setPolicyTagsInSchema(d, lbHTTPProfile.Tags)
	d.Set("nsx_id", id)
	d.Set("path", lbHTTPProfile.Path)
	d.Set("revision", lbHTTPProfile.Revision)

	d.Set("http_redirect_to", lbHTTPProfile.HttpRedirectTo)
	d.Set("http_redirect_to_https", lbHTTPProfile.HttpRedirectToHttps)
	d.Set("idle_timeout", lbHTTPProfile.IdleTimeout)
	d.Set("request_body_size", lbHTTPProfile.RequestBodySize)
	d.Set("request_header_size", lbHTTPProfile.RequestHeaderSize)
	d.Set("response_buffering", lbHTTPProfile.ResponseBuffering)
	d.Set("response_header_size", lbHTTPProfile.ResponseHeaderSize)
	d.Set("response_timeout", lbHTTPProfile.ResponseTimeout)
	d.Set("server_keep_alive", lbHTTPProfile.ServerKeepAlive)
	d.Set("x_forwarded_for", lbHTTPProfile.XForwardedFor)

	return nil
}

func resourceNsxtVpcLBHttpApplicationProfileUpdate(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VPC LBHttpProfile ID")
	}

	err := resourceNsxtVpcLBHttpApplicationProfilePatch(d, m, id)
	if err != nil {
		return handleUpdateError("VPC LBHttpProfile", id, err)
	}

	return resourceNsxtVpcLBHttpApplicationProfileRead(d, m)
}

func resourceNsxtVpcLBHttpApplicationProfileDelete(d *schema.ResourceData, m interface{}) error {
	return resourceNsxtVpcLBAppProfileDelete(d, m)
}

// resourceNsxtVpcLBAppProfileExists is shared by all VPC LB application profile resources
func resourceNsxtVpcLBAppProfileExists(sessionContext utl.SessionContext, id string, connector client.Connector) (bool, error) {
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbAppProfilesClient(sessionContext, connector)
	_, err := c.Get(parents[0], parents[1], parents[2], id)
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving resource VPC LBAppProfile", err)
}

// resourceNsxtVpcLBAppProfileDelete is shared by all VPC LB application profile resources
func resourceNsxtVpcLBAppProfileDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VPC LBAppProfile ID")
	}

	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbAppProfilesClient(sessionContext, connector)

	force := true
	err := c.Delete(parents[0], parents[1], parents[2], id, &force)
	if err != nil {
		return handleDeleteError("VPC LBAppProfile", id, err)
	}
	return nil
}
