// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNsxtPolicyCertificate_basic(t *testing.T) {
	name := getTestCertificateName(false)
	testResourceName := "data.nsxt_policy_certificate.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccEnvDefined(t, "NSXT_TEST_CERTIFICATE_NAME")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyCertificateReadTemplate(false, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
		},
	})
}

func testAccNsxtPolicyCertificateReadTemplate(withContext bool, name string) string {
	context := ""
	if withContext {
		context = testAccNsxtPolicyMultitenancyContext()
	}
	return fmt.Sprintf(`
data "nsxt_policy_certificate" "test" {
%s
  display_name = "%s"
}`, context, name)
}

func TestAccDataSourceContextBasedNsxtPolicyCertificate_basic(t *testing.T) {
	name := getTestCertificateName(false)
	testResourceName := "data.nsxt_policy_certificate.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccEnvDefined(t, "NSXT_TEST_LOCAL_DEV")
			testAccEnvDefined(t, "NSXT_TEST_CERTIFICATE_NAME")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyCertificateReadTemplate(true, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
		},
	})
}
