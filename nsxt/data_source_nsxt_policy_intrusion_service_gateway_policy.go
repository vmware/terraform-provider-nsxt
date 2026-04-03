// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func dataSourceNsxtPolicyIntrusionServiceGatewayPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyIntrusionServiceGatewayPolicyRead,
		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"domain":       getDataSourceDomainNameSchema(),
			"category": {
				Type:         schema.TypeString,
				Description:  "Category",
				ValidateFunc: validation.StringInSlice(gatewayPolicyCategoryValues, false),
				Optional:     true,
				Computed:     true,
			},
		},
	}
}

func listIntrusionServiceGatewayPolicies(d *schema.ResourceData, domain string, m interface{}, connector client.Connector) ([]model.IdsGatewayPolicy, error) {
	client := cliIntrusionServiceGatewayPoliciesClient(getSessionContext(d, m), connector)
	if client == nil {
		return nil, policyResourceNotSupportedError()
	}

	var results []model.IdsGatewayPolicy
	var cursor *string
	total := 0

	for {
		includeMarkForDeleteObjectsParam := false
		policies, err := client.List(domain, cursor, &includeMarkForDeleteObjectsParam, nil, nil, nil, nil, nil)
		if err != nil {
			return results, err
		}
		results = append(results, policies.Results...)
		if total == 0 && policies.ResultCount != nil {
			total = int(*policies.ResultCount)
		}

		cursor = policies.Cursor
		if len(results) >= total {
			return results, nil
		}
	}
}

func dataSourceNsxtPolicyIntrusionServiceGatewayPolicyRead(d *schema.ResourceData, m interface{}) error {
	domain := d.Get("domain").(string)
	objID := d.Get("id").(string)
	objName := d.Get("display_name").(string)
	category := d.Get("category").(string)

	if objID != "" {
		connector := getPolicyConnector(m)
		client := cliIntrusionServiceGatewayPoliciesClient(getSessionContext(d, m), connector)
		if client == nil {
			return policyResourceNotSupportedError()
		}
		obj, err := client.Get(domain, objID)
		if err != nil {
			if isNotFoundError(err) {
				return fmt.Errorf("Intrusion Service Gateway Policy with ID %s was not found", objID)
			}
			return fmt.Errorf("Error while reading Intrusion Service Gateway Policy %s: %v", objID, err)
		}
		d.SetId(*obj.Id)
		d.Set("display_name", obj.DisplayName)
		d.Set("description", obj.Description)
		d.Set("path", obj.Path)
		d.Set("category", obj.Category)
		return nil
	}

	if objName == "" && category == "" {
		return fmt.Errorf("Intrusion Service Gateway Policy id, display name or category must be specified")
	}

	connector := getPolicyConnector(m)
	policies, err := listIntrusionServiceGatewayPolicies(d, domain, m, connector)
	if err != nil {
		return fmt.Errorf("Error while reading Intrusion Service Gateway Policies: %v", err)
	}

	var perfectMatch []model.IdsGatewayPolicy
	var prefixMatch []model.IdsGatewayPolicy
	for _, policy := range policies {
		if category != "" && policy.Category != nil && category != *policy.Category {
			continue
		}
		if objName != "" && policy.DisplayName != nil {
			if strings.HasPrefix(*policy.DisplayName, objName) {
				prefixMatch = append(prefixMatch, policy)
			}
			if *policy.DisplayName == objName {
				perfectMatch = append(perfectMatch, policy)
			}
		} else {
			prefixMatch = append(prefixMatch, policy)
		}
	}

	var obj model.IdsGatewayPolicy
	if len(perfectMatch) > 0 {
		if len(perfectMatch) > 1 {
			return fmt.Errorf("Found multiple Intrusion Service Gateway Policies with name '%s'", objName)
		}
		obj = perfectMatch[0]
	} else if len(prefixMatch) > 0 {
		if len(prefixMatch) > 1 {
			return fmt.Errorf("Found multiple Intrusion Service Gateway Policies with name starting with '%s' and category '%s'", objName, category)
		}
		obj = prefixMatch[0]
	} else {
		return fmt.Errorf("Intrusion Service Gateway Policy with name '%s' and category '%s' was not found", objName, category)
	}

	d.SetId(*obj.Id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("path", obj.Path)
	d.Set("category", obj.Category)
	return nil
}
