// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNsxtPolicyIdpsCustomSignature_basic(t *testing.T) {
	testResourceName := "nsxt_policy_idps_custom_signature.test"
	testDataSourceName := "data.nsxt_policy_idps_custom_signature.test"
	sid := testAccIdpsCustomSignatureSIDBase + 3
	signature := fmt.Sprintf(`alert tcp any any -> any any (msg:"TF Acc DataSource Test"; content:"tf_acc_ds_%d"; nocase; metadata:signature_severity Medium; sid:%d; rev:1;)`, sid, sid)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "4.2.0")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNsxtPolicyIdpsCustomSignatureBasic(signature),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(testDataSourceName, "id", testResourceName, "id"),
					resource.TestCheckResourceAttr(testDataSourceName, "signature_version_id", testAccIdpsCustomSignatureVersionID),
					resource.TestCheckResourceAttrSet(testDataSourceName, "id"),
				),
			},
		},
	})
}

func testAccDataSourceNsxtPolicyIdpsCustomSignatureBasic(signature string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_idps_custom_signature" "test" {
  signature_version_id = %q
  signature            = %q
  publish              = false
}

data "nsxt_policy_idps_custom_signature" "test" {
  id = nsxt_policy_idps_custom_signature.test.id
  depends_on = [nsxt_policy_idps_custom_signature.test]
}`, testAccIdpsCustomSignatureVersionID, signature)
}
