// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNsxtPolicyEdgeTransportNode_basic(t *testing.T) {
	etnName := getEdgeTransportNodeName()
	testResourceName := "data.nsxt_policy_edge_transport_node.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccOnlyLocalManager(t)
			testAccPreCheck(t)
			testAccEnvDefined(t, "NSXT_TEST_EDGE_TRANSPORT_NODE")
			testAccNSXVersion(t, "9.0.0")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNSXEdgeTransportNodeReadTemplate(etnName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", etnName),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttrSet(testResourceName, "unique_id"),
				),
			},
		},
	})
}

func testAccNSXEdgeTransportNodeReadTemplate(name string) string {
	return fmt.Sprintf(`
data "nsxt_policy_edge_transport_node" "test" {
  display_name = "%s"
}`, name)
}
