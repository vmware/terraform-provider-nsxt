package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccResourceNsxtPolicyParentGatewayPolicy_basic(t *testing.T) {
	testAccResourceNsxtPolicyParentGatewayPolicyBasic(t, false, func() {
		testAccPreCheck(t)
	})
}

func TestAccResourceNsxtPolicyParentGatewayPolicy_multitenancy(t *testing.T) {
	testAccResourceNsxtPolicyParentGatewayPolicyBasic(t, true, func() {
		testAccPreCheck(t)
		testAccOnlyMultitenancy(t)
	})
}

func testAccResourceNsxtPolicyParentGatewayPolicyBasic(t *testing.T, withContext bool, preCheck func()) {
	testResourceName := "nsxt_policy_parent_gateway_policy.test"

	name := getAccTestResourceName()
	updatedName := getAccTestResourceName()
	locked := "true"
	updatedLocked := "false"
	seqNum := "1"
	updatedSeqNum := "2"
	tcpStrict := "true"
	updatedTCPStrict := "false"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  preCheck,
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyParentGatewayPolicyCheckDestroy(state, updatedName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyParentGatewayPolicyTemplate(withContext, name, locked, seqNum, tcpStrict),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyGatewayPolicyExists(testResourceName, defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "locked", locked),
					resource.TestCheckResourceAttr(testResourceName, "sequence_number", seqNum),
					resource.TestCheckResourceAttr(testResourceName, "tcp_strict", tcpStrict),
				),
			},
			{
				Config: testAccNsxtPolicyParentGatewayPolicyTemplate(withContext, updatedName, updatedLocked, updatedSeqNum, updatedTCPStrict),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyGatewayPolicyExists(testResourceName, defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatedName),
					resource.TestCheckResourceAttr(testResourceName, "locked", updatedLocked),
					resource.TestCheckResourceAttr(testResourceName, "sequence_number", updatedSeqNum),
					resource.TestCheckResourceAttr(testResourceName, "tcp_strict", updatedTCPStrict),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyParentGatewayPolicy_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_parent_gateway_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyParentGatewayPolicyCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyParentGatewayPolicyTemplate(false, name, "true", "1", "true"),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccResourceNsxtPolicyImportIDRetriever(testResourceName),
			},
		},
	})
}

func TestAccResourceNsxtPolicyParentGatewayPolicy_importBasic_multitenancy(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_parent_gateway_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyMultitenancy(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyGatewayPolicyCheckDestroy(state, name, defaultDomain)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyParentGatewayPolicyTemplate(true, name, "true", "1", "true"),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccResourceNsxtPolicyImportIDRetriever(testResourceName),
			},
		},
	})
}

func testAccNsxtPolicyParentGatewayPolicyCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_parent_gateway_policy" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		domain := rs.Primary.Attributes["domain"]
		exists, err := resourceNsxtPolicyGatewayPolicyExistsInDomain(testAccGetSessionContext(), resourceID, domain, connector)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("Policy GatewayPolicy %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyT1GwDataSourceTemplate(withContext bool) string {
	context := ""
	if withContext {
		context = testAccNsxtPolicyMultitenancyContext()
	}
	return fmt.Sprintf(`

data "nsxt_policy_edge_cluster" "EC" {
  display_name = "%s"
}

resource "nsxt_policy_tier1_gateway" "t1_gw" {
  %s
  display_name      = "t1-gw-policy-rule-test"
  edge_cluster_path = data.nsxt_policy_edge_cluster.EC.path
}
`, getEdgeClusterName(), context)
}

func testAccNsxtPolicyParentGatewayPolicyTemplate(withContext bool, name, locked, seqNum, tcpStrict string) string {
	context := ""
	if withContext {
		context = testAccNsxtPolicyMultitenancyContext()
	}
	return testAccNsxtPolicyT1GwDataSourceTemplate(withContext) + fmt.Sprintf(`
resource "nsxt_policy_parent_gateway_policy" "test" {
%s
  display_name    = "%s"
  description     = "Acceptance Test"
  category        = "LocalGatewayRules"
  locked          = %s
  sequence_number = %s
  stateful        = true
  tcp_strict      = %s
  
  tag {
    scope = "color"
    tag   = "orange"
  }
}`, context, name, locked, seqNum, tcpStrict) + testAccNsxtPolicyGatewayPolicyRuleTemplate(withContext)
}

func testAccNsxtPolicyGatewayPolicyRuleTemplate(withContext bool) string {
	context := ""
	if withContext {
		context = testAccNsxtPolicyMultitenancyContext()
	}
	return fmt.Sprintf(`

resource "nsxt_policy_gateway_policy_rule" "rule1" {
  %s
  display_name       = "rule2"
  description        = "Terraform provisioned gateway Policy Rule"
  policy_path = nsxt_policy_parent_gateway_policy.test.path
  sequence_number    = 1
  action             = "DROP"
  logged             = true
  scope              = [nsxt_policy_tier1_gateway.t1_gw.path]
}
`, context)
}
