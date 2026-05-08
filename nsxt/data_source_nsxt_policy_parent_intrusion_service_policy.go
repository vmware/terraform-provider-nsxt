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

	"github.com/vmware/terraform-provider-nsxt/api/infra/domains"
)

func dataSourceNsxtPolicyParentIntrusionServicePolicy() *schema.Resource {
	// NOTE: This parent data source intentionally excludes the "rule" field to provide lightweight policy metadata retrieval without embedded rules.
	// Use nsxt_policy_intrusion_service_policy for complete policy with rules
	return &schema.Resource{
		Read: dataSourceNsxtPolicyParentIntrusionServicePolicyRead,
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
			"locked": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether a security policy should be locked. If the security policy is locked by a user, then no other user would be able to modify this security policy.",
			},
			"sequence_number": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "This field is used to resolve conflicts between multiple policies that have rules that match the same packet.",
			},
			"comments": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Comments for security policy lock/unlock.",
			},
			"tag":      getTagsSchema(),
			"revision": getRevisionSchema(),
		},
	}
}

func dataSourceNsxtPolicyParentIntrusionServicePolicyRead(d *schema.ResourceData, m interface{}) error {
	domain := d.Get("domain").(string)
	objID := d.Get("id").(string)
	objName := d.Get("display_name").(string)
	category := d.Get("category").(string)
	context := getSessionContext(d, m)

	if objID != "" {
		connector := getPolicyConnector(m)
		client := domains.NewIntrusionServicePoliciesClient(context, connector)
		if client == nil {
			return policyResourceNotSupportedError()
		}
		obj, err := client.Get(domain, objID)
		if err != nil {
			if isNotFoundError(err) {
				return fmt.Errorf("Parent Intrusion Service Policy with ID %s was not found", objID)
			}
			return fmt.Errorf("Error while reading Parent Intrusion Service Policy %s: %v", objID, err)
		}
		d.SetId(*obj.Id)
		setParentIntrusionServicePolicyDataSourceSchema(d, obj)
		return nil
	}

	if objName == "" && category == "" {
		return fmt.Errorf("Parent Intrusion Service Policy id, display name or category must be specified")
	}

	policies, err := listIntrusionServicePolicies(context, domain, m)
	if err != nil {
		return fmt.Errorf("Error while reading Parent Intrusion Service Policies: %v", err)
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
			return fmt.Errorf("Found multiple Parent Intrusion Service Policies with name '%s'", objName)
		}
		obj = perfectMatch[0]
	} else if len(prefixMatch) > 0 {
		if len(prefixMatch) > 1 {
			return fmt.Errorf("Found multiple Parent Intrusion Service Policies with name starting with '%s' and category '%s'", objName, category)
		}
		obj = prefixMatch[0]
	} else {
		return fmt.Errorf("Parent Intrusion Service Policy with name '%s' and category '%s' was not found", objName, category)
	}

	d.SetId(*obj.Id)
	setParentIntrusionServicePolicyDataSourceSchema(d, obj)
	return nil
}

func setParentIntrusionServicePolicyDataSourceSchema(d *schema.ResourceData, obj model.IdsSecurityPolicy) {
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("path", obj.Path)
	d.Set("domain", getDomainFromResourcePath(*obj.Path))
	d.Set("category", obj.Category)
	d.Set("stateful", obj.Stateful != nil && *obj.Stateful)
	d.Set("locked", obj.Locked != nil && *obj.Locked)
	d.Set("comments", obj.Comments)
	if obj.SequenceNumber != nil {
		d.Set("sequence_number", int(*obj.SequenceNumber))
	}
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("revision", obj.Revision)
}
