package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccResourceNsxtPolicyGatewayPolicyRule_basic(t *testing.T) {
	testAccResourceNsxtPolicyGatewayPolicyBasic(t, false, func() {
		testAccPreCheck(t)
	})
}

func testAccResourceNsxtPolicyGatewayPolicyRuleBasic(t *testing.T, withContext bool, preCheck func()) {
	policyName := getAccTestResourceName()
	ruleName := getAccTestResourceName()
	updatedName := getAccTestResourceName()
	resourceName := "nsxt_policy_gateway_policy"
	testResourceName := fmt.Sprintf("%s.test", resourceName)
	direction1 := "IN"
	proto1 := "IPV4"
	tag1 := "abc"
	action := "REJECT"
	description := "Terraform provisioned gateway Policy Rule"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  preCheck,
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyGatewayPolicyCheckDestroy(state, updatedName, defaultDomain)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyGatewayRule(policyName, updatedName, direction1, proto1, tag1, withContext, ruleName, action, description),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyGatewayPolicyExists(testResourceName, defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "display_name", ruleName),
					resource.TestCheckResourceAttr(testResourceName, "description", description),
					resource.TestCheckResourceAttr(testResourceName, "action", action),
				),
			},
		},
	})
}

func testAccNsxtPolicyGatewayRule(policyName, updatedName, direction1, proto1, tag1 string, withContext bool, ruleName, action, description string) string {

	return testAccNsxtPolicyGatewayPolicyWithRule(policyName, updatedName, direction1, proto1, tag1, withContext) + fmt.Sprintf(`

resource "nsxt_policy_gateway_policy_rule" "rule1" {
  display_name       = "%s"
  description        = "%s"
  policy_path        = data.nsxt_policy_gateway_policy.test.path
  sequence_number    = 1
  action             = "%s"
  logged             = true
  scope              = [data.nsxt_policy_tier1_gateway.gwt1test.path]
}
`, ruleName, description, action)
}
