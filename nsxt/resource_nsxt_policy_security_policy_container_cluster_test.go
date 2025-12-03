// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceNsxtPolicySecurityPolicyContainerCluster_basic(t *testing.T) {
	testResourceName := "nsxt_policy_security_policy_container_cluster.test"

	name := getAccTestResourceName()
	updateName := getAccTestResourceName()
	createDesc := "Acceptance tests description"
	updateDesc := "Updated acceptance tests description"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "9.0.0")
			testAccEnvDefined(t, "NSXT_TEST_CONTAINER_CLUSTER")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicySecurityPolicyContainerClusterCheckDestroy(state, updateName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicySecurityPolicyContainerClusterTemplate(name, createDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicySecurityPolicyContainerClusterExists(testResourceName, defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", createDesc),
				),
			},
			{
				Config: testAccNsxtPolicySecurityPolicyContainerClusterTemplate(updateName, updateDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicySecurityPolicyContainerClusterExists(testResourceName, defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateName),
					resource.TestCheckResourceAttr(testResourceName, "description", updateDesc),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicySecurityPolicyContainerCluster_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_security_policy_container_cluster.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "9.0.0")
			testAccEnvDefined(t, "NSXT_TEST_CONTAINER_CLUSTER")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicySecurityPolicyContainerClusterCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicySecurityPolicyContainerClusterTemplate(name, ""),
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

func testAccNsxtPolicySecurityPolicyContainerClusterExists(resourceName string, domainName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy SecurityPolicyContainerCluster resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy SecurityPolicyContainerCluster resource ID not set in resources")
		}
		policyPath := rs.Primary.Attributes["policy_path"]
		exists, err := resourceNsxtPolicySecurityPolicyContainerClusterExists(resourceID, connector, policyPath)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Error while retrieving policy SecurityPolicyContainerCluster ID %s", resourceID)
		}
		return nil
	}
}

func testAccNsxtPolicySecurityPolicyContainerClusterCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_security_policy_container_cluster" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		policyPath := rs.Primary.Attributes["policy_path"]
		exists, err := resourceNsxtPolicySecurityPolicyContainerClusterExists(resourceID, connector, policyPath)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("Policy SecurityPolicyContainerCluster %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicySecurityPolicyContainerClusterTemplate(name, description string) string {
	locked := "true"
	seqNum := "1"
	tcpStrict := "true"
	return testAccNSXPolicyContainerClusterReadTemplate(os.Getenv("NSXT_TEST_CONTAINER_CLUSTER")) +
		testAccNsxtPolicyParentSecurityPolicyTemplate(false, name, locked, seqNum, tcpStrict) +
		fmt.Sprintf(`
resource "nsxt_policy_security_policy_container_cluster" "test" {
  display_name           = "%s"
  description            = "%s"
  policy_path            = nsxt_policy_parent_security_policy.test.path
  container_cluster_path = data.nsxt_policy_container_cluster.test.path
}`, name, description)
}
