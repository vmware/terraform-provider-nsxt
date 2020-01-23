/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccDataSourceNsxtPolicyCertificate_basic(t *testing.T) {
	name := getTestCertificateName(false)
	testResourceName := "data.nsxt_policy_certificate.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccEnvDefined(t, "NSXT_TEST_CERTIFICATE_NAME") },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyCertificateReadTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
			{
				Config: testAccNsxtPolicyEmptyTemplate(),
			},
		},
	})
}

func testAccNsxtPolicyCertificateReadTemplate(name string) string {
	return fmt.Sprintf(`
data "nsxt_policy_certificate" "test" {
  display_name = "%s"
}`, name)
}
