// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

var vpcGroupPathExample = "/orgs/[org]/projects/[project]/vpcs/[vpc]/groups/[group]"

func resourceNsxtVPCGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtVPCGroupCreate,
		Read:   resourceNsxtVPCGroupRead,
		Update: resourceNsxtVPCGroupUpdate,
		Delete: resourceNsxtVPCGroupDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtVersionCheckImporter("9.0.0", "VPC Group", getVpcPathResourceImporter(vpcGroupPathExample)),
		},

		Schema: getPolicyGroupSchema(false),
	}
}

func resourceNsxtVPCGroupCreate(d *schema.ResourceData, m interface{}) error {
	if !util.NsxVersionHigherOrEqual("9.0.0") {
		return fmt.Errorf("VPC Group resource requires NSX version 9.0.0 or higher")
	}
	if isConfigScopedCacheMode() {
		_ = d.Set("tag", initPolicyTagsSet(getPolicyTagsWithProviderManagedDefaults(d, m)))
	}
	if err := resourceNsxtPolicyGroupGeneralCreate(d, m, false); err != nil {
		return err
	}
	InvalidateCacheForResourceType("group")
	return nil
}

func resourceNsxtVPCGroupRead(d *schema.ResourceData, m interface{}) error {
	if !isCacheEnabledForRead(d) {
		return resourceNsxtPolicyGroupGeneralRead(d, m, false)
	}

	connector := getPolicyConnector(m)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Group ID")
	}

	var obj *model.Group
	var err error
	obj, _, _, err = CacheAwareResourceRead[model.Group](
		d,
		m,
		connector,
		id,
		"group",
		model.GroupBindingType(),
		func() (*model.Group, error) {
			client := cliGroupsClient(getSessionContext(d, m), connector)
			if client == nil {
				return nil, policyResourceNotSupportedError()
			}
			readObj, readErr := client.Get("", id)
			if readErr != nil {
				return nil, readErr
			}
			return &readObj, nil
		},
		func(patchObj *model.Group) error {
			client := cliGroupsClient(getSessionContext(d, m), connector)
			if client == nil {
				return policyResourceNotSupportedError()
			}
			return client.Patch("", id, *patchObj)
		},
	)

	if err != nil {
		return handleReadError(d, "Group", id, err)
	}

	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("path", obj.Path)
	d.Set("revision", obj.Revision)
	groupType := ""
	if len(obj.GroupType) > 0 && util.NsxVersionHigherOrEqual("3.2.0") {
		groupType = obj.GroupType[0]
		d.Set("group_type", groupType)
	}
	criteria, conditions, err := fromGroupExpressionData(obj.Expression)
	if err != nil {
		return err
	}
	d.Set("criteria", criteria)
	d.Set("conjunction", conditions)
	identityGroups, err := getIdentityGroupsData(obj.ExtendedExpression)
	if err != nil {
		return err
	}
	var extendedCriteria []map[string]interface{}
	if len(identityGroups) > 0 {
		identityGroupsMap := make(map[string]interface{})
		identityGroupsMap["identity_group"] = identityGroups
		extendedCriteria = append(extendedCriteria, identityGroupsMap)
	}
	d.Set("extended_criteria", extendedCriteria)

	return nil
}

func resourceNsxtVPCGroupUpdate(d *schema.ResourceData, m interface{}) error {
	if isConfigScopedCacheMode() {
		_ = d.Set("tag", initPolicyTagsSet(getPolicyTagsWithProviderManagedDefaults(d, m)))
	}
	if err := resourceNsxtPolicyGroupGeneralUpdate(d, m, false); err != nil {
		return err
	}
	InvalidateCacheForResourceType("group")
	return nil
}

func resourceNsxtVPCGroupDelete(d *schema.ResourceData, m interface{}) error {
	return resourceNsxtPolicyGroupGeneralDelete(d, m, false)
}
