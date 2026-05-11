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
)

func dataSourceNsxtPolicyParentIntrusionServiceGatewayPolicy() *schema.Resource {
	// NOTE: This parent data source intentionally excludes the "rule" field to provide lightweight policy metadata retrieval without embedded rules
	// Use nsxt_policy_intrusion_service_gateway_policy for complete policy with rules
	return &schema.Resource{
		Read: dataSourceNsxtPolicyParentIntrusionServiceGatewayPolicyRead,
		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"domain":       getDataSourceDomainNameSchema(),
			"category": {
				Type:         schema.TypeString,
				Description:  "Category",
				ValidateFunc: validation.StringInSlice(idpsGatewayPolicyCategoryValues, false),
				Optional:     true,
				Computed:     true,
			},
			"stateful": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "When it is stateful, the state of the network connects are tracked and a stateful packet inspection is performed. Note: Intrusion Service Gateway Policies are always stateful.",
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

func dataSourceNsxtPolicyParentIntrusionServiceGatewayPolicyRead(d *schema.ResourceData, m interface{}) error {
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
				return fmt.Errorf("Parent Intrusion Service Gateway Policy with ID %s was not found", objID)
			}
			return fmt.Errorf("Error while reading Parent Intrusion Service Gateway Policy %s: %v", objID, err)
		}
		d.SetId(*obj.Id)
		setParentIntrusionServiceGatewayPolicyDataSourceSchema(d, obj)
		return nil
	}

	if objName == "" && category == "" {
		return fmt.Errorf("Parent Intrusion Service Gateway Policy id, display name or category must be specified")
	}

	connector := getPolicyConnector(m)
	policies, err := listIntrusionServiceGatewayPolicies(d, domain, m, connector)
	if err != nil {
		return fmt.Errorf("Error while reading Parent Intrusion Service Gateway Policies: %v", err)
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
			return fmt.Errorf("Found multiple Parent Intrusion Service Gateway Policies with name '%s'", objName)
		}
		obj = perfectMatch[0]
	} else if len(prefixMatch) > 0 {
		if len(prefixMatch) > 1 {
			return fmt.Errorf("Found multiple Parent Intrusion Service Gateway Policies with name starting with '%s' and category '%s'", objName, category)
		}
		obj = prefixMatch[0]
	} else {
		return fmt.Errorf("Parent Intrusion Service Gateway Policy with name '%s' and category '%s' was not found", objName, category)
	}

	d.SetId(*obj.Id)
	setParentIntrusionServiceGatewayPolicyDataSourceSchema(d, obj)
	return nil
}

func setParentIntrusionServiceGatewayPolicyDataSourceSchema(d *schema.ResourceData, obj model.IdsGatewayPolicy) {
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
