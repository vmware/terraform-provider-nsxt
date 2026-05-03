// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

func dataSourceNsxtPolicyIntrusionServicePolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyIntrusionServicePolicyRead,
		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"context":      getContextSchema(false, false, false),
			"domain":       getDataSourceDomainNameSchema(),
			"category": {
				Type:         schema.TypeString,
				Description:  "Category",
				ValidateFunc: validation.StringInSlice(idpsDfwPolicyCategoryValues, false),
				Optional:     true,
				Computed:     true,
			},
			"stateful": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "When it is stateful, the state of the network connects are tracked and a stateful packet inspection is performed. Note: Intrusion Service Policies are always stateful.",
			},
		},
	}
}

func listIntrusionServicePolicies(context utl.SessionContext, domain string, m interface{}) ([]model.IdsSecurityPolicy, error) {
	connector := getPolicyConnector(m)
	client := cliIntrusionServicePoliciesClient(context, connector)
	if client == nil {
		return nil, policyResourceNotSupportedError()
	}

	var results []model.IdsSecurityPolicy
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

func dataSourceNsxtPolicyIntrusionServicePolicyRead(d *schema.ResourceData, m interface{}) error {
	domain := d.Get("domain").(string)
	objID := d.Get("id").(string)
	objName := d.Get("display_name").(string)
	category := d.Get("category").(string)
	context := getSessionContext(d, m)

	if objID != "" {
		connector := getPolicyConnector(m)
		client := cliIntrusionServicePoliciesClient(context, connector)
		if client == nil {
			return policyResourceNotSupportedError()
		}
		obj, err := client.Get(domain, objID)
		if err != nil {
			if isNotFoundError(err) {
				return fmt.Errorf("Intrusion Service Policy with ID %s was not found", objID)
			}
			return fmt.Errorf("Error while reading Intrusion Service Policy %s: %v", objID, err)
		}
		d.SetId(*obj.Id)
		d.Set("display_name", obj.DisplayName)
		d.Set("description", obj.Description)
		d.Set("path", obj.Path)
		d.Set("domain", getDomainFromResourcePath(*obj.Path))
		d.Set("category", obj.Category)
		d.Set("stateful", obj.Stateful != nil && *obj.Stateful)
		return nil
	}

	if objName == "" && category == "" {
		return fmt.Errorf("Intrusion Service Policy id, display name or category must be specified")
	}

	policies, err := listIntrusionServicePolicies(context, domain, m)
	if err != nil {
		return fmt.Errorf("Error while reading Intrusion Service Policies: %v", err)
	}

	var perfectMatch []model.IdsSecurityPolicy
	var prefixMatch []model.IdsSecurityPolicy
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

	var obj model.IdsSecurityPolicy
	if len(perfectMatch) > 0 {
		if len(perfectMatch) > 1 {
			return fmt.Errorf("Found multiple Intrusion Service Policies with name '%s'", objName)
		}
		obj = perfectMatch[0]
	} else if len(prefixMatch) > 0 {
		if len(prefixMatch) > 1 {
			return fmt.Errorf("Found multiple Intrusion Service Policies with name starting with '%s' and category '%s'", objName, category)
		}
		obj = prefixMatch[0]
	} else {
		return fmt.Errorf("Intrusion Service Policy with name '%s' and category '%s' was not found", objName, category)
	}

	d.SetId(*obj.Id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("path", obj.Path)
	d.Set("domain", getDomainFromResourcePath(*obj.Path))
	d.Set("category", obj.Category)
	d.Set("stateful", obj.Stateful != nil && *obj.Stateful)
	return nil
}
