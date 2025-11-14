// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func resourceNsxtPolicyParentGatewayPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyParentGatewayPolicyCreate,
		Read:   resourceNsxtPolicyParentGatewayPolicyRead,
		Update: resourceNsxtPolicyParentGatewayPolicyUpdate,
		Delete: resourceNsxtPolicyGatewayPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtDomainResourceImporter,
		},

		Schema: getPolicyGatewayPolicySchema(false, false),
	}
}

func resourceNsxtPolicyParentGatewayPolicyCreate(d *schema.ResourceData, m interface{}) error {
	return resourceNsxtPolicyGatewayPolicyGeneralCreate(d, m, false)
}

func resourceNsxtPolicyParentGatewayPolicyRead(d *schema.ResourceData, m interface{}) error {
	return resourceNsxtPolicyGatewayPolicyGeneralRead(d, m, false)
}

func resourceNsxtPolicyParentGatewayPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceNsxtPolicyGatewayPolicyGeneralUpdate(d, m, false)
}
