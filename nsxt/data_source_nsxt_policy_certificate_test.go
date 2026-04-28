// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDataSourceNsxtPolicyCertificate_basic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "data.nsxt_policy_certificate.test"
	var accTestTlsCertID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(*terraform.State) error {
			return testAccTlsCertificateDeleteGlobal(accTestTlsCertID)
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					id, err := testAccTlsCertificateSelfSignedCreateGlobal(name)
					if err != nil {
						t.Fatal(err)
					}
					accTestTlsCertID = id
				},
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
	name := getAccTestResourceName()
	testResourceName := "data.nsxt_policy_certificate.test"
	var accTestTlsCertID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			if testAccMultitenancyProjectID() == "" {
				t.Skipf("This test requires NSXT_PROJECT_ID or NSXT_VPC_PROJECT_ID for project-scoped TLS certificate API")
			}
		},
		Providers: testAccProviders,
		CheckDestroy: func(*terraform.State) error {
			return testAccTlsCertificateDeleteProject(accTestTlsCertID)
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					id, err := testAccTlsCertificateSelfSignedCreateProject(name)
					if err != nil {
						t.Fatal(err)
					}
					accTestTlsCertID = id
				},
				Config: testAccNsxtPolicyCertificateReadTemplate(true, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
		},
	})
}
