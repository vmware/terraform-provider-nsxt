package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccResourceNsxtPolicyGatewayPolicyRule_basic(t *testing.T) {
	testAccResourceNsxtPolicyGatewayPolicyRuleBasic(t, false, func() {
		testAccPreCheck(t)
	})
}

func testAccResourceNsxtPolicyGatewayPolicyRuleBasic(t *testing.T, withContext bool, preCheck func()) {
	ruleName := getAccTestResourceName()
	policyName := getAccTestResourceName()
	direction := "IN"
	proto := "IPV4"
	tag := "abc"
	action := "REJECT"
	description := "Terraform provisioned gateway Policy Rule"
	policyResourceName := "nsxt_policy_gateway_policy"
	testPolicyResourceName := fmt.Sprintf("%s.test", policyResourceName)
	testruleResourceName := "nsxt_policy_gateway_policy_rule.rule1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  preCheck,
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyGatewayPolicyRuleCheckDestroy(state, ruleName)
		},
		Steps: []resource.TestStep{
			{
				Config:             testAccNsxtPolicyGatewayRule(policyResourceName, policyName, direction, proto, tag, withContext, ruleName, action, description),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyGatewayPolicyExists(testPolicyResourceName, defaultDomain),
					testAccNsxtPolicyGatewayPolicyRuleExists(testruleResourceName),
					resource.TestCheckResourceAttr(testruleResourceName, "display_name", ruleName),
					resource.TestCheckResourceAttr(testruleResourceName, "description", description),
					resource.TestCheckResourceAttr(testruleResourceName, "action", action),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyGatewayPolicyRule_importBasic(t *testing.T) {

	name := getAccTestResourceName()
	locked := "true"
	seqNum := "1"
	tcpStrict := "true"

	testruleResourceName := "nsxt_policy_gateway_policy_rule.rule1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyGatewayPolicyRuleCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyParentGatewayPolicyTemplate(false, name, locked, seqNum, tcpStrict),
			},
			{
				ResourceName:      testruleResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccResourceNsxtPolicyImportIDRetriever(testruleResourceName),
			},
		},
	})
}

func testAccNsxtPolicyGatewayPolicyRuleCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_gateway_policy_rule" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		exists, err := resourceNsxtPolicyGatewayPolicyRuleExists(testAccGetSessionContext(), resourceID, rs.Primary.Attributes["policy_path"], connector)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("Policy GatewayPolicy %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyGatewayPolicyRuleExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy GatewayPolicy resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy GatewayPolicy resource ID not set in resources")
		}

		exists, err := resourceNsxtPolicyGatewayPolicyRuleExists(testAccGetSessionContext(), resourceID, rs.Primary.Attributes["policy_path"], connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Error while retrieving policy GatewayPolicy ID %s", resourceID)
		}
		return nil
	}
}

func testAccNsxtPolicyGatewayRule(policyResourceName, policyName, direction, proto, tag string, withContext bool, ruleName, action, description string) string {
	return testAccNsxtPolicyGatewayPolicyWithRule(policyResourceName, policyName, direction, proto, tag, withContext) + fmt.Sprintf(`
	
resource "nsxt_policy_gateway_policy_rule" "rule1" {
	display_name       = "%s"
	description        = "%s"
	policy_path        = nsxt_policy_gateway_policy.test.path
	sequence_number    = 1
	action             = "%s"
	logged             = true
	scope              = [nsxt_policy_tier1_gateway.gwt1test.path]
}
	
	`, ruleName, description, action)
}
