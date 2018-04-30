/* Copyright © 2018 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vmware/go-vmware-nsxt/loadbalancer"
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

func isDataRequired(protocol string) bool {
	if protocol == "udp" {
		return true
	}

	return false
}

func getSendDescription(protocol string) string {
	if protocol == "tcp" {
		return "If both send and receive are not specified, then just a TCP connection is established (3-way handshake) to validate server is healthy, no data is sent."
	}

	return "The data to be sent to the monitored server."
}

// The only differences between tcp and udp monitors are required vs. optional data fields,
// and their descriptions
func getLbL4MonitorSchema(protocol string) map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
		"tag":          getTagsSchema(),
		"fall_count":   getLbMonitorFallCountSchema(),
		"interval":     getLbMonitorIntervalSchema(),
		"monitor_port": getLbMonitorPortSchema(),
		"rise_count":   getLbMonitorRiseCountSchema(),
		"timeout":      getLbMonitorTimeoutSchema(),
		"receive": &schema.Schema{
			Type:        schema.TypeString,
			Description: "Expected data, if specified, can be anywhere in the response and it has to be a string, regular expressions are not supported",
			Optional:    !isDataRequired(protocol),
			Required:    isDataRequired(protocol),
		},
		"send": &schema.Schema{
			Type:        schema.TypeString,
			Description: getSendDescription(protocol),
			Optional:    !isDataRequired(protocol),
			Required:    isDataRequired(protocol),
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
				"name": &schema.Schema{
					Type:        schema.TypeString,
					Description: "Header name",
					Required:    true,
				},
				"value": &schema.Schema{
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
