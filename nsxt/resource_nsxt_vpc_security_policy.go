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

// VPC Security Policy importer with version check
func nsxtVpcSecurityPolicyImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	// Check NSX version compatibility for import
	if !util.NsxVersionHigherOrEqual("9.0.0") {
		return []*schema.ResourceData{d}, fmt.Errorf("VPC Security Policy import requires NSX version 9.0.0 or higher")
	}
	
	// Use the existing VPC path importer logic
	importer := getVpcPathResourceImporter(vpcSecurityPolicyPathExample)
	return importer(d, m)
}

func resourceNsxtVPCSecurityPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtVPCSecurityPolicyCreate,
		Read:   resourceNsxtVPCSecurityPolicyRead,
		Update: resourceNsxtVPCSecurityPolicyUpdate,
		Delete: resourceNsxtVPCSecurityPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtVpcSecurityPolicyImporter,
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

	client := domains.NewSecurityPoliciesClient(getSessionContext(d, m), connector)
	err := client.Delete("", id)

	if err != nil {
		return handleDeleteError("Security Policy", id, err)
	}

	return nil
}
