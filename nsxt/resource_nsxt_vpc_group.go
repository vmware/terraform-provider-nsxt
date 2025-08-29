// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var vpcGroupPathExample = "/orgs/[org]/projects/[project]/vpcs/[vpc]/groups/[group]"

// VPC Group importer with version check
func nsxtVpcGroupImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	// Check NSX version compatibility for import
	if !util.NsxVersionHigherOrEqual("9.0.0") {
		return []*schema.ResourceData{d}, fmt.Errorf("VPC Group import requires NSX version 9.0.0 or higher")
	}
	
	// Use the existing VPC path importer logic
	importer := getVpcPathResourceImporter(vpcGroupPathExample)
	return importer(d, m)
}

func resourceNsxtVPCGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtVPCGroupCreate,
		Read:   resourceNsxtVPCGroupRead,
		Update: resourceNsxtVPCGroupUpdate,
		Delete: resourceNsxtVPCGroupDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtVpcGroupImporter,
		},

		Schema: getPolicyGroupSchema(false),
	}
}

func resourceNsxtVPCGroupCreate(d *schema.ResourceData, m interface{}) error {
	if !util.NsxVersionHigherOrEqual("9.0.0") {
		return fmt.Errorf("VPC Group resource requires NSX version 9.0.0 or higher")
	}
	return resourceNsxtPolicyGroupGeneralCreate(d, m, false)
}

func resourceNsxtVPCGroupRead(d *schema.ResourceData, m interface{}) error {
	return resourceNsxtPolicyGroupGeneralRead(d, m, false)
}

func resourceNsxtVPCGroupUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceNsxtPolicyGroupGeneralUpdate(d, m, false)
}

func resourceNsxtVPCGroupDelete(d *schema.ResourceData, m interface{}) error {
	return resourceNsxtPolicyGroupGeneralDelete(d, m, false)
}
