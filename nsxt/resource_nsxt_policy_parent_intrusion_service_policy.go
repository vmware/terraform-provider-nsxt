// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func resourceNsxtPolicyParentIntrusionServicePolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyParentIntrusionServicePolicyCreate,
		Read:   resourceNsxtPolicyParentIntrusionServicePolicyRead,
		Update: resourceNsxtPolicyParentIntrusionServicePolicyUpdate,
		Delete: resourceNsxtPolicyIntrusionServicePolicyDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtDomainResourceImporter,
		},
		Schema: getPolicySecurityPolicySchema(true, true, false, false),
	}
}

func resourceNsxtPolicyParentIntrusionServicePolicyCreate(d *schema.ResourceData, m interface{}) error {
	return resourceNsxtPolicyIntrusionServicePolicyGeneralCreate(d, m, false)
}

func resourceNsxtPolicyParentIntrusionServicePolicyRead(d *schema.ResourceData, m interface{}) error {
	return resourceNsxtPolicyIntrusionServicePolicyGeneralRead(d, m, false)
}

func resourceNsxtPolicyParentIntrusionServicePolicyUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceNsxtPolicyIntrusionServicePolicyGeneralUpdate(d, m, false)
}
