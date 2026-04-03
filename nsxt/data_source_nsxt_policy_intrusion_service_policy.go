// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
		return nil
	}

	if objName == "" {
		return fmt.Errorf("Intrusion Service Policy id or display name must be specified")
	}

	policies, err := listIntrusionServicePolicies(context, domain, m)
	if err != nil {
		return fmt.Errorf("Error while reading Intrusion Service Policies: %v", err)
	}

	var perfectMatch []model.IdsSecurityPolicy
	var prefixMatch []model.IdsSecurityPolicy
	for _, policy := range policies {
		if policy.DisplayName != nil {
			if strings.HasPrefix(*policy.DisplayName, objName) {
				prefixMatch = append(prefixMatch, policy)
			}
			if *policy.DisplayName == objName {
				perfectMatch = append(perfectMatch, policy)
			}
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
			return fmt.Errorf("Found multiple Intrusion Service Policies with name starting with '%s'", objName)
		}
		obj = prefixMatch[0]
	} else {
		return fmt.Errorf("Intrusion Service Policy with name '%s' was not found", objName)
	}

	d.SetId(*obj.Id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("path", obj.Path)
	return nil
}
