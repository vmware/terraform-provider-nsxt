/* Copyright © 2022 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

var lbHTTPAppProfileTypeValues = []string{"HTTP"}
var lbHTTPAppProfileTypeMap = map[string]string{
	model.LBAppProfile_RESOURCE_TYPE_LBHTTPPROFILE: "HTTP",
}

var lbHTTPAppProfileXFFValues = []string{"INSERT", "REPLACE"}

func dataSourceNsxtPolicyLBHTTPAppProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyLBHTTPAppProfileRead,

		Schema: map[string]*schema.Schema{
			"id": getDataSourceIDSchema(),
			"type": {
				Type:         schema.TypeString,
				Description:  "Application Profile Type",
				Optional:     true,
				Default:      "HTTP",
				ValidateFunc: validation.StringInSlice(lbHTTPAppProfileTypeValues, false),
			},
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"http_direct_to": {
				Type:         schema.TypeString,
				Description:  "Target URL when virtual server is down",
				Optional:     true,
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},
			"http_redirect_to_https": {
				Type:        schema.TypeBool,
				Description: "Force Redirect via HTTPS",
				Optional:    true,
			},
			"idle_timeout": {
				Type:        schema.TypeInt,
				Description: "Idle timeout",
				Optional:    true,
			},
			"ntlm": {
				Type:        schema.TypeBool,
				Description: "NTLM Authentication",
				Optional:    true,
			},
			"request_body_size": {
				Type:        schema.TypeInt,
				Description: "Request body size",
				Optional:    true,
			},
			"request_header_size": {
				Type:        schema.TypeInt,
				Description: "Request header size",
				Optional:    true,
			},
			"response_buffering": {
				Type:        schema.TypeBool,
				Description: "Response buffering",
				Optional:    true,
			},
			"response_header_size": {
				Type:        schema.TypeInt,
				Description: "Request header size",
				Optional:    true,
			},
			"response_timeout": {
				Type:        schema.TypeInt,
				Description: "Request timeout",
				Optional:    true,
			},
			"server_keep_alive": {
				Type:        schema.TypeBool,
				Description: "Server keep alive",
				Optional:    true,
			},
			"x_forwarded_for": {
				Type:         schema.TypeString,
				Description:  "Insert or replace x_forwarded_for",
				Optional:     true,
				ValidateFunc: validation.StringInSlice(lbHTTPAppProfileXFFValues, false),
			},
		},
	}
}

func policyLbHTTPAppProfileConvert(obj *data.StructValue, requestedType string) (*model.LBHttpProfile, error) {
	converter := bindings.NewTypeConverter()
	converter.SetMode(bindings.REST)

	data, errs := converter.ConvertToGolang(obj, model.LBHttpProfileBindingType())
	if errs != nil {
		return nil, errs[0]
	}

	profile := data.(model.LBHttpProfile)
	profileType, ok := lbHTTPAppProfileTypeMap[profile.ResourceType]
	if !ok {
		return nil, fmt.Errorf("Unknown LB Application Profile type %s", profile.ResourceType)
	}
	if (requestedType != "HTTP") && (requestedType != profileType) {
		return nil, nil
	}
	return &profile, nil
}

func dataSourceNsxtPolicyLBHTTPAppProfileRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := infra.NewLbAppProfilesClient(connector)

	objID := d.Get("id").(string)
	objTypeValue, typeSet := d.GetOk("type")
	objType := objTypeValue.(string)
	objName := d.Get("display_name").(string)
	var result *model.LBHttpProfile
	if objID != "" {
		// Get by id
		objGet, err := client.Get(objID)

		if err != nil {
			return handleDataSourceReadError(d, "LBHttpProfile", objID, err)
		}
		result, err = policyLbHTTPAppProfileConvert(objGet, objType)
		if err != nil {
			return fmt.Errorf("Error while converting LBHttpProfile %s: %v", objID, err)
		}
		if result == nil {
			return fmt.Errorf("LBHttpProfile with ID '%s' and type %s was not found", objID, objType)
		}
	} else if objName == "" && !typeSet {
		return fmt.Errorf("Error obtaining LBHttpProfile ID or name or type during read")
	} else {
		// Get by full name/prefix
		includeMarkForDeleteObjectsParam := false
		objList, err := client.List(nil, &includeMarkForDeleteObjectsParam, nil, nil, nil, nil)
		if err != nil {
			return handleListError("LBHttpProfile", err)
		}
		// go over the list to find the correct one (prefer a perfect match. If not - prefix match)
		var perfectMatch []model.LBHttpProfile
		var prefixMatch []model.LBHttpProfile
		for _, objInList := range objList.Results {
			resourceType, err := objInList.String("resource_type")
			if err != nil {
				return fmt.Errorf("Couldn't read resource_type %s: %v", objID, err)
			}
			if resourceType == "LBHttpProfile" {
				obj, err := policyLbHTTPAppProfileConvert(objInList, objType)
				if err != nil {
					return fmt.Errorf("Error while converting LBHttpProfile %s: %v", objID, err)
				}
				if obj == nil {
					continue
				}
				if objName != "" && obj.DisplayName != nil && strings.HasPrefix(*obj.DisplayName, objName) {
					prefixMatch = append(prefixMatch, *obj)
				}
				if obj.DisplayName != nil && *obj.DisplayName == objName {
					perfectMatch = append(perfectMatch, *obj)
				}
				if objName == "" && typeSet {
					// match only by type
					perfectMatch = append(perfectMatch, *obj)
				}
			}
			if len(perfectMatch) > 0 {
				if len(perfectMatch) > 1 {
					return fmt.Errorf("Found multiple LBHttpProfiles with name '%s'", objName)
				}
				result = &perfectMatch[0]
			} else if len(prefixMatch) > 0 {
				if len(prefixMatch) > 1 {
					return fmt.Errorf("Found multiple LBHttpProfiles with name starting with '%s'", objName)
				}
				result = &prefixMatch[0]
			} else {
				return fmt.Errorf("LBHttpProfile with name '%s' and type %s was not found", objName, objType)
			}
		}
	}

	d.SetId(*result.Id)
	d.Set("display_name", result.DisplayName)
	d.Set("type", lbHTTPAppProfileTypeMap[result.ResourceType])
	d.Set("description", result.Description)
	d.Set("path", result.Path)
	d.Set("http_direct_to", result.HttpRedirectTo)
	d.Set("http_redirect_to_https", result.HttpRedirectToHttps)
	d.Set("idle_timeout", result.IdleTimeout)
	d.Set("ntlm", result.Ntlm)
	d.Set("request_body_size", result.RequestBodySize)
	d.Set("request_header_size", result.RequestHeaderSize)
	d.Set("response_buffering", result.ResponseBuffering)
	d.Set("response_header_size", result.ResponseHeaderSize)
	d.Set("response_timeout", result.ResponseTimeout)
	d.Set("server_keep_alive", result.ServerKeepAlive)
	d.Set("x_forwarded_for", result.XForwardedFor)
	return nil
}
