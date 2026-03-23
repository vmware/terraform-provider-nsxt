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

func TestAccDataSourceNsxtPolicyIdpsSignatureDiff_basic(t *testing.T) {
	testDataSourceName := "data.nsxt_policy_idps_signature_diff.test"
	// Use a SID in a separate range so path segment (e.g. 1006000000) does not collide with resource tests (1005999991+).
	sid := 6000000
	signature := fmt.Sprintf(`alert tcp any any -> any any (msg:"TF Acc Diff Test"; content:"tf_acc_diff_%d"; nocase; metadata:signature_severity Low; sid:%d; rev:1;)`, sid, sid)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "4.2.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyIdpsCustomSignatureCheckDestroy(state)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNsxtPolicyIdpsSignatureDiffBasic(signature),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testDataSourceName, "signature_version_id", testAccIdpsCustomSignatureVersionID),
					resource.TestCheckResourceAttr(testDataSourceName, "id", testAccIdpsCustomSignatureVersionID),
					resource.TestCheckResourceAttrSet(testDataSourceName, "newly_added_signatures.#"),
					resource.TestCheckResourceAttrSet(testDataSourceName, "deleted_signatures.#"),
					resource.TestCheckResourceAttrSet(testDataSourceName, "existing_signatures.#"),
				),
			},
		},
	})
}

func testAccDataSourceNsxtPolicyIdpsSignatureDiffBasic(signature string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_idps_custom_signature" "test" {
  signature_version_id = %q
  signature            = %q
  publish              = false
}

data "nsxt_policy_idps_signature_diff" "test" {
  signature_version_id = %q
  depends_on = [nsxt_policy_idps_custom_signature.test]
}`, testAccIdpsCustomSignatureVersionID, signature, testAccIdpsCustomSignatureVersionID)
}
