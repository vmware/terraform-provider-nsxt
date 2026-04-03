// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func resourceNsxtPolicyParentIntrusionServiceGatewayPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyParentIntrusionServiceGatewayPolicyCreate,
		Read:   resourceNsxtPolicyParentIntrusionServiceGatewayPolicyRead,
		Update: resourceNsxtPolicyParentIntrusionServiceGatewayPolicyUpdate,
		Delete: resourceNsxtPolicyIntrusionServiceGatewayPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtDomainResourceImporter,
		},
		Schema: getIntrusionServiceGatewayPolicySchema(false),
	}
}

func resourceNsxtPolicyParentIntrusionServiceGatewayPolicyCreate(d *schema.ResourceData, m interface{}) error {
	return resourceNsxtPolicyIntrusionServiceGatewayPolicyGeneralCreate(d, m, false)
}

func resourceNsxtPolicyParentIntrusionServiceGatewayPolicyRead(d *schema.ResourceData, m interface{}) error {
	return resourceNsxtPolicyIntrusionServiceGatewayPolicyGeneralRead(d, m, false)
}

func resourceNsxtPolicyParentIntrusionServiceGatewayPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceNsxtPolicyIntrusionServiceGatewayPolicyGeneralUpdate(d, m, false)
}
