/* Copyright © 2023 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/trust_management"
)

var accTestPrincipleIdentityCreateAttributes = map[string]string{
	"name":      getAccTestResourceName(),
	"node_id":   "node-2",
	"role_path": "/orgs/default",
	"role":      "org_admin",
}

func TestAccResourceNsxtPrincipleIdentity_basic(t *testing.T) {
	testResourceName := "nsxt_principle_identity.test"
	certPem, _, err := testAccGenerateTLSKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPrincipleIdentityCheckDestroy(state, accTestPrincipleIdentityCreateAttributes["name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPrincipleIdentityCreate(certPem),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPrincipleIdentityExists(accTestPrincipleIdentityCreateAttributes["name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "name", accTestPrincipleIdentityCreateAttributes["name"]),
					resource.TestCheckResourceAttr(testResourceName, "node_id", accTestPrincipleIdentityCreateAttributes["node_id"]),
					resource.TestCheckResourceAttr(testResourceName, "roles_for_path.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "roles_for_path.0.path", accTestPrincipleIdentityCreateAttributes["role_path"]),
					resource.TestCheckResourceAttr(testResourceName, "roles_for_path.0.role.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "roles_for_path.0.role.0.role", accTestPrincipleIdentityCreateAttributes["role"]),

					resource.TestCheckResourceAttrSet(testResourceName, "certificate_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "certificate_pem"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPrincipleIdentity_import_basic(t *testing.T) {
	testResourceName := "nsxt_principle_identity.test"
	certPem, _, err := testAccGenerateTLSKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPrincipleIdentityCheckDestroy(state, accTestPrincipleIdentityCreateAttributes["name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPrincipleIdentityCreate(certPem),
			},
			{
				ResourceName:            testResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"certificate_pem"},
			},
		},
	})
}

func testAccNsxtPrincipleIdentityExists(name string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("PrincipleIdentity resource %s not found in resources", resourceName)
		}

		piID := rs.Primary.Attributes["id"]
		if piID == "" {
			return fmt.Errorf("PrincipleIdentity resource ID not set in resources")
		}
		tmClient := trust_management.NewPrincipalIdentitiesClient(connector)
		_, err := tmClient.Get(piID)
		if err != nil {
			if isNotFoundError(err) {
				return fmt.Errorf("PrincipleIdentity %s does not exist", name)
			}
		}

		return err
	}
}

func testAccNsxtPrincipleIdentityCheckDestroy(state *terraform.State, name string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_principle_identity" {
			continue
		}

		piID := rs.Primary.Attributes["id"]
		if piID == "" {
			return fmt.Errorf("PrincipleIdentity resource ID not set in resources")
		}
		tmClient := trust_management.NewPrincipalIdentitiesClient(connector)
		_, err := tmClient.Get(piID)
		if err != nil {
			if isNotFoundError(err) {
				return nil
			}
			return err
		}
		return fmt.Errorf("PrincipleIdentity %s still exists", name)
	}
	return nil
}

func testAccNsxtPrincipleIdentityCreate(certPem string) string {
	attrMap := accTestPrincipleIdentityCreateAttributes
	return fmt.Sprintf(`
resource "nsxt_principle_identity" "test" {
    certificate_pem = <<-EOT
%s
    EOT
    name            = "%s"
    node_id         = "%s"

    roles_for_path {
        path = "%s"

        role {
            role = "%s"
        }
    }
}`, certPem, attrMap["name"], attrMap["node_id"], attrMap["role_path"], attrMap["role"])
}
