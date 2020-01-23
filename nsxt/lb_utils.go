/* Copyright © 2018 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/vmware/go-vmware-nsxt/loadbalancer"
	"log"
	"net/http"
)

// Helpers for common LB monitor schema settings
func getLbMonitorFallCountSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeInt,
		Description: "Number of consecutive checks that must fail before marking it down",
		Optional:    true,
		Default:     3,
	}
}

func getLbMonitorIntervalSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeInt,
		Description: "The frequency at which the system issues the monitor check (in seconds)",
		Optional:    true,
		Default:     5,
	}
}

func getLbMonitorPortSchema() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Description:  "If the monitor port is specified, it would override pool member port setting for healthcheck. A port range is not supported",
		Optional:     true,
		ValidateFunc: validateSinglePort(),
	}
}

func getLbMonitorRiseCountSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeInt,
		Description: "Number of consecutive checks that must pass before marking it up",
		Optional:    true,
		Default:     3,
	}
}

func getLbMonitorTimeoutSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeInt,
		Description: "Number of seconds the target has to respond to the monitor request",
		Optional:    true,
		Default:     15,
	}
}

func getLbMonitorRequestBodySchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Description: "String to send as HTTP health check request body. Valid only for certain HTTP methods like POST",
		Optional:    true,
	}
}

func getLbMonitorRequestMethodSchema() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Description:  "Health check method for HTTP monitor type",
		Optional:     true,
		ValidateFunc: validation.StringInSlice([]string{"GET", "HEAD", "OPTIONS", "POST", "PUT"}, false),
		Default:      "GET",
	}
}

func getLbMonitorRequestURLSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Description: "URL used for HTTP monitor",
		Optional:    true,
		Default:     "/",
	}
}

func getLbMonitorRequestVersionSchema() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Description:  "HTTP request version",
		Optional:     true,
		ValidateFunc: validation.StringInSlice([]string{"HTTP_VERSION_1_0", "HTTP_VERSION_1_1"}, false),
		Default:      "HTTP_VERSION_1_1",
	}
}

func getLbMonitorResponseBodySchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Description: "If HTTP specified, healthcheck HTTP response body is matched against the specified string (regular expressions not supported), and succeeds only if there is a match",
		Optional:    true,
	}
}

func getLbMonitorResponseStatusCodesSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "The HTTP response status code should be a valid HTTP status code",
		Elem: &schema.Schema{
			Type:         schema.TypeInt,
			ValidateFunc: validation.IntBetween(100, 505),
		},
		Optional: true,
		Computed: true,
	}
}

func isLbMonitorDataRequired(protocol string) bool {
	if protocol == "udp" {
		return true
	}

	return false
}

func getLbMonitorSendDescription(protocol string) string {
	if protocol == "tcp" {
		return "If both send and receive are not specified, then just a TCP connection is established (3-way handshake) to validate server is healthy, no data is sent."
	}

	return "The data to be sent to the monitored server."
}

// The only differences between tcp and udp monitors are required vs. optional data fields,
// and their descriptions
func getLbL4MonitorSchema(protocol string) map[string]*schema.Schema {
	dataRequired := isLbMonitorDataRequired(protocol)

	return map[string]*schema.Schema{
		"revision": getRevisionSchema(),
		"description": {
			Type:        schema.TypeString,
			Description: "Description of this resource",
			Optional:    true,
		},
		"display_name": {
			Type:        schema.TypeString,
			Description: "The display name of this resource. Defaults to ID if not set",
			Optional:    true,
			Computed:    true,
		},
		"tag":          getTagsSchema(),
		"fall_count":   getLbMonitorFallCountSchema(),
		"interval":     getLbMonitorIntervalSchema(),
		"monitor_port": getLbMonitorPortSchema(),
		"rise_count":   getLbMonitorRiseCountSchema(),
		"timeout":      getLbMonitorTimeoutSchema(),
		"receive": {
			Type:        schema.TypeString,
			Description: "Expected data, if specified, can be anywhere in the response and it has to be a string, regular expressions are not supported",
			Optional:    !dataRequired,
			Required:    dataRequired,
		},
		"send": {
			Type:        schema.TypeString,
			Description: getLbMonitorSendDescription(protocol),
			Optional:    !dataRequired,
			Required:    dataRequired,
		},
	}
}

func getLbHTTPHeaderSchema(description string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: description,
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Description: "Header name",
					Required:    true,
				},
				"value": {
					Type:        schema.TypeString,
					Description: "Header value",
					Required:    true,
				},
			},
		},
	}
}

func getLbHTTPHeaderFromSchema(d *schema.ResourceData, attrName string) []loadbalancer.LbHttpRequestHeader {
	headers := d.Get(attrName).(*schema.Set).List()
	var headerList []loadbalancer.LbHttpRequestHeader
	for _, header := range headers {
		data := header.(map[string]interface{})
		elem := loadbalancer.LbHttpRequestHeader{
			HeaderName:  data["name"].(string),
			HeaderValue: data["value"].(string)}

		headerList = append(headerList, elem)
	}
	return headerList
}

func setLbHTTPHeaderInSchema(d *schema.ResourceData, attrName string, headers []loadbalancer.LbHttpRequestHeader) {
	var headerList []map[string]string
	for _, header := range headers {
		elem := make(map[string]string)
		elem["name"] = header.HeaderName
		elem["value"] = header.HeaderValue
		headerList = append(headerList, elem)
	}
	d.Set(attrName, headerList)
}

func resourceNsxtLbMonitorDelete(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(nsxtClients).NsxtClient
	if nsxClient == nil {
		return resourceNotSupportedError()
	}

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	resp, err := nsxClient.ServicesApi.DeleteLoadBalancerMonitor(nsxClient.Context, id)
	if err != nil {
		return fmt.Errorf("Error during LbMonitor delete: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[DEBUG] LbMonitor %s not found", id)
		d.SetId("")
	}
	return nil
}
