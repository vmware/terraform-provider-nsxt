/* Copyright © 2018 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	api "github.com/vmware/go-vmware-nsxt"
	"github.com/vmware/go-vmware-nsxt/loadbalancer"
	"log"
	"net/http"
)

var xForwardedValues = []string{"INSERT", "REPLACE"}

func resourceNsxtLbHttpApplicationProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtLbHttpApplicationProfileCreate,
		Read:   resourceNsxtLbHttpApplicationProfileRead,
		Update: resourceNsxtLbHttpApplicationProfileUpdate,
		Delete: resourceNsxtLbHttpApplicationProfileDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"revision": getRevisionSchema(),
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Description of this resource",
				Optional:    true,
			},
			"display_name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The display name of this resource. Defaults to ID if not set",
				Optional:    true,
				Computed:    true,
			},
			"tag": getTagsSchema(),
			"http_redirect_to": &schema.Schema{
				Type:        schema.TypeString,
				Description: "A URL that incoming requests for that virtual server can be temporarily redirected to, If a website is temporarily down or has moved",
				Optional:    true,
			},
			"http_redirect_to_https": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "A boolean flag which reflects whether the client will automatically be redirected to use SSL",
				Optional:    true,
				Default:     false,
			},
			"idle_timeout": &schema.Schema{
				Type:         schema.TypeInt,
				Description:  "Timeout in seconds to specify how long an HTTP application can remain idle",
				Optional:     true,
				Default:      15,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"ntlm": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "A boolean flag which reflects whether NTLM challenge/response methodology will be used over HTTP",
				Optional:    true,
				Default:     false,
			},
			"request_body_size": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "Maximum request body size in bytes (Unlimited if not specified)",
				Optional:    true,
			},
			"request_header_size": &schema.Schema{
				Type:         schema.TypeInt,
				Description:  "Maximum request header size in bytes. Requests with larger header size will be processed as best effort whereas a request with header below this specified size is guaranteed to be processed",
				Optional:     true,
				Default:      1024,
				ValidateFunc: validation.IntBetween(1, 65536),
			},
			"response_timeout": &schema.Schema{
				Type:         schema.TypeInt,
				Description:  "Number of seconds waiting for the server response before the connection is closed",
				Optional:     true,
				Default:      60,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"x_forwarded_for": &schema.Schema{
				Type:         schema.TypeString,
				Description:  "When this value is set, the x_forwarded_for header in the incoming request will be inserted or replaced",
				Optional:     true,
				ValidateFunc: validation.StringInSlice(xForwardedValues, false),
			},
		},
	}
}

func resourceNsxtLbHttpApplicationProfileCreate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	httpRedirectTo := d.Get("http_redirect_to").(string)
	httpRedirectToHttps := d.Get("http_redirect_to_https").(bool)
	idleTimeout := int64(d.Get("idle_timeout").(int))
	ntlm := d.Get("ntlm").(bool)
	requestBodySize := int64(d.Get("request_body_size").(int))
	requestHeaderSize := int64(d.Get("request_header_size").(int))
	responseTimeout := int64(d.Get("response_timeout").(int))
	xForwardedFor := d.Get("x_forwarded_for").(string)
	lbHttpApplicationProfile := loadbalancer.LbHttpProfile{
		Description:         description,
		DisplayName:         displayName,
		Tags:                tags,
		HttpRedirectTo:      httpRedirectTo,
		HttpRedirectToHttps: httpRedirectToHttps,
		IdleTimeout:         idleTimeout,
		Ntlm:                ntlm,
		RequestBodySize:     requestBodySize,
		RequestHeaderSize:   requestHeaderSize,
		ResponseTimeout:     responseTimeout,
		XForwardedFor:       xForwardedFor,
	}

	lbHttpApplicationProfile, resp, err := nsxClient.ServicesApi.CreateLoadBalancerHttpProfile(nsxClient.Context, lbHttpApplicationProfile)

	if err != nil {
		return fmt.Errorf("Error during LbHttpApplicationProfile create: %v", err)
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected status returned during LbHttpApplicationProfile create: %v", resp.StatusCode)
	}
	d.SetId(lbHttpApplicationProfile.Id)

	return resourceNsxtLbHttpApplicationProfileRead(d, m)
}

func resourceNsxtLbHttpApplicationProfileRead(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	lbHttpApplicationProfile, resp, err := nsxClient.ServicesApi.ReadLoadBalancerHttpProfile(nsxClient.Context, id)
	if err != nil {
		return fmt.Errorf("Error during LbHttpApplicationProfile read: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[DEBUG] LbHttpApplicationProfile %s not found", id)
		d.SetId("")
		return nil
	}
	d.Set("revision", lbHttpApplicationProfile.Revision)
	d.Set("description", lbHttpApplicationProfile.Description)
	d.Set("display_name", lbHttpApplicationProfile.DisplayName)
	setTagsInSchema(d, lbHttpApplicationProfile.Tags)
	d.Set("http_redirect_to", lbHttpApplicationProfile.HttpRedirectTo)
	d.Set("http_redirect_to_https", lbHttpApplicationProfile.HttpRedirectToHttps)
	d.Set("idle_timeout", lbHttpApplicationProfile.IdleTimeout)
	d.Set("ntlm", lbHttpApplicationProfile.Ntlm)
	d.Set("request_body_size", lbHttpApplicationProfile.RequestBodySize)
	d.Set("request_header_size", lbHttpApplicationProfile.RequestHeaderSize)
	d.Set("response_timeout", lbHttpApplicationProfile.ResponseTimeout)
	d.Set("x_forwarded_for", lbHttpApplicationProfile.XForwardedFor)

	return nil
}

func resourceNsxtLbHttpApplicationProfileUpdate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	revision := int32(d.Get("revision").(int))
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	httpRedirectTo := d.Get("http_redirect_to").(string)
	httpRedirectToHttps := d.Get("http_redirect_to_https").(bool)
	idleTimeout := int64(d.Get("idle_timeout").(int))
	ntlm := d.Get("ntlm").(bool)
	requestBodySize := int64(d.Get("request_body_size").(int))
	requestHeaderSize := int64(d.Get("request_header_size").(int))
	responseTimeout := int64(d.Get("response_timeout").(int))
	xForwardedFor := d.Get("x_forwarded_for").(string)
	lbHttpApplicationProfile := loadbalancer.LbHttpProfile{
		Revision:            revision,
		Description:         description,
		DisplayName:         displayName,
		Tags:                tags,
		HttpRedirectTo:      httpRedirectTo,
		HttpRedirectToHttps: httpRedirectToHttps,
		IdleTimeout:         idleTimeout,
		Ntlm:                ntlm,
		RequestBodySize:     requestBodySize,
		RequestHeaderSize:   requestHeaderSize,
		ResponseTimeout:     responseTimeout,
		XForwardedFor:       xForwardedFor,
	}

	lbHttpApplicationProfile, resp, err := nsxClient.ServicesApi.UpdateLoadBalancerHttpProfile(nsxClient.Context, id, lbHttpApplicationProfile)

	if err != nil || resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Error during LbHttpApplicationProfile update: %v", err)
	}

	return resourceNsxtLbHttpApplicationProfileRead(d, m)
}

func resourceNsxtLbHttpApplicationProfileDelete(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	resp, err := nsxClient.ServicesApi.DeleteLoadBalancerApplicationProfile(nsxClient.Context, id)
	if err != nil {
		return fmt.Errorf("Error during LbHttpApplicationProfile delete: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[DEBUG] LbHttpApplicationProfile %s not found", id)
		d.SetId("")
	}
	return nil
}
