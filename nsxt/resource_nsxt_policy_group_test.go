package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/domains"
	"testing"
)

func TestAccResourceNsxtPolicyGroup_basicImport(t *testing.T) {
	name := fmt.Sprintf("test-nsx-policy-group-ipaddrs")
	testResourceName := "nsxt_policy_group.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyGroupCheckDestroy(state, name, defaultDomain)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyGroupIPAddressImportTemplate(name),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceNsxtPolicyGroup_singleIPAddressCriteria(t *testing.T) {
	name := fmt.Sprintf("test-nsx-policy-group-ipaddrs")
	testResourceName := "nsxt_policy_group.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyGroupCheckDestroy(state, name, defaultDomain)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyGroupIPAddressCreateTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyGroupExists(testResourceName, defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "domain", defaultDomain),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckNoResourceAttr(testResourceName, "conjunction"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "criteria.#", "1"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyGroup_multipleIPAddressCriteria(t *testing.T) {
	name := fmt.Sprintf("test-nsx-policy-group-ipaddrs")
	updatedName := fmt.Sprintf("%s-update", name)
	testResourceName := "nsxt_policy_group.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyGroupCheckDestroy(state, name, defaultDomain)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyGroupIPAddressMultipleCreateTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyGroupExists(testResourceName, defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "domain", defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "conjunction.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "criteria.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "2"),
				),
			},
			{
				Config: testAccNsxtPolicyGroupIPAddressMultipleUpdateTemplate(updatedName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyGroupExists(testResourceName, defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatedName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttr(testResourceName, "domain", defaultDomain),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "conjunction.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "criteria.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyGroup_pathCriteria(t *testing.T) {
	name := fmt.Sprintf("test-nsx-policy-group-paths")
	testResourceName := "nsxt_policy_group.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyGroupCheckDestroy(state, name, defaultDomain)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyGroupPathsCreateTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyGroupExists(testResourceName, defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "domain", defaultDomain),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckNoResourceAttr(testResourceName, "conjunction"),
					resource.TestCheckResourceAttr(testResourceName, "criteria.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "criteria.0.path_expression.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "criteria.0.path_expression.0.member_paths.#", "2"),
				),
			},
			{
				Config: testAccNsxtPolicyGroupPathsUpdateTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyGroupExists(testResourceName, defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "domain", defaultDomain),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckNoResourceAttr(testResourceName, "conjunction"),
					resource.TestCheckResourceAttr(testResourceName, "criteria.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "criteria.0.path_expression.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "criteria.0.path_expression.0.member_paths.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyGroupPathsPrerequisites(),
			},
		},
	})
}

func TestAccResourceNsxtPolicyGroup_nestedCriteria(t *testing.T) {
	name := fmt.Sprintf("test-nsx-policy-group-nested")
	updatedName := fmt.Sprintf("%s-update", name)
	testResourceName := "nsxt_policy_group.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyGroupCheckDestroy(state, name, defaultDomain)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyGroupNestedConditionCreateTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyGroupExists(testResourceName, defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckNoResourceAttr(testResourceName, "conjunction"),
					resource.TestCheckResourceAttr(testResourceName, "domain", defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "criteria.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "criteria.0.condition.#", "2"),
				),
			},
			{
				Config: testAccNsxtPolicyGroupNestedConditionUpdateTemplate(updatedName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyGroupExists(testResourceName, defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatedName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttr(testResourceName, "domain", defaultDomain),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckNoResourceAttr(testResourceName, "conjunction"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "criteria.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "criteria.0.condition.#", "2"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyGroup_multipleCriteria(t *testing.T) {
	name := fmt.Sprintf("test-nsx-policy-group-multiple")
	updatedName := fmt.Sprintf("%s-update", name)
	testResourceName := "nsxt_policy_group.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyGroupCheckDestroy(state, name, defaultDomain)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyGroupMultipleConditionCreateTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyGroupExists(testResourceName, defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "domain", defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "conjunction.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "criteria.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "criteria.0.condition.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "criteria.1.condition.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyGroupMultipleConditionUpdateTemplate(updatedName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyGroupExists(testResourceName, defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatedName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "domain", defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "conjunction.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "criteria.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "criteria.0.condition.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "criteria.1.condition.#", "1"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyGroup_multipleNestedCriteria(t *testing.T) {
	name := fmt.Sprintf("test-nsx-policy-group-multiple-nested")
	testResourceName := "nsxt_policy_group.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyGroupCheckDestroy(state, name, defaultDomain)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyGroupMultipleNestedCreateTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyGroupExists(testResourceName, defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "domain", defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "conjunction.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "criteria.#", "3"),
					resource.TestCheckResourceAttr(testResourceName, "criteria.0.condition.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "criteria.1.condition.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "criteria.2.ipaddress_expression.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyGroupMultipleNestedUpdateTemplateDeleteIPAddrs(name),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyGroupExists(testResourceName, defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "domain", defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "conjunction.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "criteria.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "criteria.0.condition.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "criteria.1.condition.#", "2"),
				),
			},
			{
				Config: testAccNsxtPolicyGroupMultipleNestedUpdateTemplateDeleteNestedConds(name),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyGroupExists(testResourceName, defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "domain", defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "conjunction.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "criteria.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "criteria.0.condition.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "criteria.1.condition.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyGroupMultipleNestedUpdateTemplateDeleteLastCriteria(name),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyGroupExists(testResourceName, defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "domain", defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "conjunction.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "criteria.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "criteria.0.condition.#", "1"),
				),
			},
		},
	})
}

func testAccNsxtPolicyGroupExists(resourceName string, domainName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
		nsxClient := domains.NewDefaultGroupsClient(connector)

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy Group resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy Group resource ID not set in resources")
		}

		_, err := nsxClient.Get(domainName, resourceID)
		if err != nil {
			return fmt.Errorf("Error while retrieving policy Group ID %s. Error: %v", resourceID, err)
		}

		return nil
	}
}

func testAccNsxtPolicyGroupCheckDestroy(state *terraform.State, displayName string, domainName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	nsxClient := domains.NewDefaultGroupsClient(connector)
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_group" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		_, err := nsxClient.Get(domainName, resourceID)
		if err == nil {
			return fmt.Errorf("Policy Group %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyGroupIPAddressImportTemplate(name string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_group" "test" {
  display_name = "%s"
  description  = "Acceptance Test"
}
`, name)
}

func testAccNsxtPolicyGroupIPAddressCreateTemplate(name string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_group" "test" {
  display_name = "%s"
  description  = "Acceptance Test"

  criteria {
    ipaddress_expression {
	  ip_addresses = ["111.1.1.1", "222.2.2.2"]
	}
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }

  tag {
    scope = "scope2"
    tag   = "tag2"
  }
}
`, name)
}

func testAccNsxtPolicyGroupIPAddressMultipleCreateTemplate(name string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_group" "test" {
  display_name = "%s"
  description  = "Acceptance Test"

  criteria {
    ipaddress_expression {
	  ip_addresses = ["111.1.1.1", "222.2.2.2"]
	}
  }

  conjunction {
	operator = "OR"
  }

  criteria {
    ipaddress_expression {
	  ip_addresses = ["111.1.1.3", "222.2.2.3"]
	}
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }

  tag {
    scope = "scope2"
    tag   = "tag2"
  }
}
`, name)
}

func testAccNsxtPolicyGroupPathsPrerequisites() string {
	return fmt.Sprintf(`
data "nsxt_policy_transport_zone" "test"{
  display_name = "%s"
}

resource "nsxt_policy_segment" "test-1" {
  display_name        = "group-test-1"
  transport_zone_path = data.nsxt_policy_transport_zone.test.path
}

resource "nsxt_policy_segment" "test-2" {
  display_name        = "group-test-1"
  transport_zone_path = data.nsxt_policy_transport_zone.test.path
}`, getOverlayTransportZoneName())

}

func testAccNsxtPolicyGroupPathsCreateTemplate(name string) string {
	return testAccNsxtPolicyGroupPathsPrerequisites() + fmt.Sprintf(`

resource "nsxt_policy_group" "test" {
  display_name = "%s"
  description  = "Acceptance Test"

  criteria {
    path_expression {
      member_paths = [nsxt_policy_segment.test-1.path, nsxt_policy_segment.test-2.path]
    }
  }
}
`, name)
}

func testAccNsxtPolicyGroupPathsUpdateTemplate(name string) string {
	return testAccNsxtPolicyGroupPathsPrerequisites() + fmt.Sprintf(`

resource "nsxt_policy_group" "test" {
  display_name = "%s"
  description  = "Acceptance Test"

  criteria {
    path_expression {
      member_paths = [nsxt_policy_segment.test-1.path]
    }
  }
}
`, name)
}

func testAccNsxtPolicyGroupIPAddressMultipleUpdateTemplate(name string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_group" "test" {
  display_name = "%s"
  description  = "Acceptance Test"

  criteria {
    ipaddress_expression {
	  ip_addresses = ["111.1.1.1", "222.2.2.2"]
	}
  }

  conjunction {
	operator = "OR"
  }

  criteria {
    ipaddress_expression {
	  ip_addresses = ["111.1.1.4", "222.2.2.4"]
	}
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}
`, name)
}

func testAccNsxtPolicyGroupNestedConditionCreateTemplate(name string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_group" "test" {
  display_name = "%s"
  description  = "Acceptance Test"

  criteria {
    condition {
      key         = "Name"
      member_type = "VirtualMachine"
      operator    = "STARTSWITH"
      value       = "publicVM"
    }
    condition {
      key         = "OSName"
      member_type = "VirtualMachine"
      operator    = "CONTAINS"
      value       = "Ubuntu"
    }
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }

  tag {
    scope = "scope2"
    tag   = "tag2"
  }
}
`, name)
}

func testAccNsxtPolicyGroupNestedConditionUpdateTemplate(name string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_group" "test" {
  display_name = "%s"
  description  = "Acceptance Test"

  criteria {
    condition {
      key         = "Name"
      member_type = "VirtualMachine"
      operator    = "STARTSWITH"
      value       = "publicVM"
    }
    condition {
      key         = "Tag"
      member_type = "VirtualMachine"
      operator    = "EQUALS"
      value       = "green"
    }
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}
`, name)
}

func testAccNsxtPolicyGroupMultipleConditionCreateTemplate(name string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_group" "test" {
  display_name = "%s"
  description  = "Acceptance Test"

  criteria {
    condition {
      key         = "Name"
      member_type = "VirtualMachine"
      operator    = "STARTSWITH"
      value       = "publicVM"
    }
  }

  conjunction {
    operator = "AND"
  }

  criteria {
    condition {
      key         = "OSName"
      member_type = "VirtualMachine"
      operator    = "CONTAINS"
      value       = "Ubuntu"
    }
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }

  tag {
    scope = "scope2"
    tag   = "tag2"
  }
}
`, name)
}

func testAccNsxtPolicyGroupMultipleConditionUpdateTemplate(name string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_group" "test" {
  display_name = "%s"
  description  = "Acceptance Test"

  criteria {
    condition {
      key         = "Name"
      member_type = "VirtualMachine"
      operator    = "STARTSWITH"
      value       = "publicVM"
    }
  }

  conjunction {
    operator = "OR"
  }

  criteria {
    condition {
      key         = "Tag"
      member_type = "VirtualMachine"
      operator    = "CONTAINS"
      value       = "public"
    }
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }

}
`, name)
}

func testAccNsxtPolicyGroupMultipleNestedCreateTemplate(name string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_group" "test" {
  display_name = "%s"
  description  = "Acceptance Test"

  criteria {
    condition {
      key         = "Name"
      member_type = "VirtualMachine"
      operator    = "STARTSWITH"
      value       = "publicVM"
    }
    condition {
      key         = "Tag"
      member_type = "VirtualMachine"
      operator    = "EQUALS"
      value       = "green"
    }
  }

  conjunction {
	operator = "OR"
  }

  criteria {
    condition {
      key         = "OSName"
      member_type = "VirtualMachine"
      operator    = "CONTAINS"
      value       = "Ubuntu"
    }
    condition {
      key         = "Tag"
      member_type = "VirtualMachine"
      operator    = "EQUALS"
      value       = "public"
    }
  }

  conjunction {
	operator = "OR"
  }

  criteria {
    ipaddress_expression {
	  ip_addresses = ["111.1.1.4", "222.2.2.4"]
	}
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}
`, name)
}

func testAccNsxtPolicyGroupMultipleNestedUpdateTemplateDeleteIPAddrs(name string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_group" "test" {
  display_name = "%s"
  description  = "Acceptance Test"

  criteria {
    condition {
      key         = "Name"
      member_type = "VirtualMachine"
      operator    = "STARTSWITH"
      value       = "publicVM"
    }
    condition {
      key         = "Tag"
      member_type = "VirtualMachine"
      operator    = "EQUALS"
      value       = "green"
    }
  }

  conjunction {
	operator = "OR"
  }

  criteria {
    condition {
      key         = "OSName"
      member_type = "VirtualMachine"
      operator    = "CONTAINS"
      value       = "Ubuntu"
    }
    condition {
      key         = "Tag"
      member_type = "VirtualMachine"
      operator    = "EQUALS"
      value       = "public"
    }
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}
`, name)
}

func testAccNsxtPolicyGroupMultipleNestedUpdateTemplateDeleteNestedConds(name string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_group" "test" {
  display_name = "%s"
  description  = "Acceptance Test"

  criteria {
    condition {
      key         = "Tag"
      member_type = "VirtualMachine"
      operator    = "EQUALS"
      value       = "green"
    }
  }

  conjunction {
	operator = "OR"
  }

  criteria {
    condition {
      key         = "OSName"
      member_type = "VirtualMachine"
      operator    = "CONTAINS"
      value       = "Ubuntu"
    }
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}
`, name)
}

func testAccNsxtPolicyGroupMultipleNestedUpdateTemplateDeleteLastCriteria(name string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_group" "test" {
  display_name = "%s"
  description  = "Acceptance Test"

  criteria {
    condition {
      key         = "Tag"
      member_type = "VirtualMachine"
      operator    = "EQUALS"
      value       = "green"
    }
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}
`, name)
}
