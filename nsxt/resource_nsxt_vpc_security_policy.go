// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/vmware/terraform-provider-nsxt/api/infra/domains"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var vpcSecurityPolicyPathExample = "/orgs/[org]/projects/[project]/vpcs/[vpc]/security-policies/[security-policy]"

func resourceNsxtVPCSecurityPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtVPCSecurityPolicyCreate,
		Read:   resourceNsxtVPCSecurityPolicyRead,
		Update: resourceNsxtVPCSecurityPolicyUpdate,
		Delete: resourceNsxtVPCSecurityPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtVersionCheckImporter("9.0.0", "VPC Security Policy", getVpcPathResourceImporter(vpcSecurityPolicyPathExample)),
		},
		Schema: getPolicySecurityPolicySchema(false, true, true, true),
	}
}

func resourceNsxtVPCSecurityPolicyCreate(d *schema.ResourceData, m interface{}) error {
	if !util.NsxVersionHigherOrEqual("9.0.0") {
		return fmt.Errorf("VPC Security Policy resource requires NSX version 9.0.0 or higher")
	}
	return resourceNsxtPolicySecurityPolicyGeneralCreate(d, m, true, true)
}

func resourceNsxtVPCSecurityPolicyRead(d *schema.ResourceData, m interface{}) error {
	return resourceNsxtPolicySecurityPolicyGeneralRead(d, m, true, true)
}

func resourceNsxtVPCSecurityPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceNsxtPolicySecurityPolicyGeneralUpdate(d, m, true, true)
}

func resourceNsxtVPCSecurityPolicyDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VPC Security Policy id")
	}

	connector := getPolicyConnector(m)

	client := domains.NewSecurityPoliciesClient(commonSessionContext, connector)
	err := client.Delete("", id)

	if err != nil {
		return handleDeleteError("Security Policy", id, err)
	}

	return nil
}
