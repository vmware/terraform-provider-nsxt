/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/vmware/go-vmware-nsxt/manager"
	"net/http"
	"testing"
)

func TestAccDataSourceNsxtFirewallSection_basic(t *testing.T) {
	name := "terraform_test_firewall_section"
	testResourceName := "data.nsxt_firewall_section.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccDataSourceNsxtFirewallSectionDeleteByName(name)
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if err := testAccDataSourceNsxtFirewallSectionCreate(name); err != nil {
						panic(err)
					}
				},
				Config: testAccNSXFirewallSectionReadTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", name),
				),
			},
			{
				Config: testAccNSXNoFirewallSectionTemplate(),
			},
		},
	})
}

func testAccDataSourceNsxtFirewallSectionCreate(name string) error {
	nsxClient, err := testAccGetClient()
	if err != nil {
		return fmt.Errorf("Error during test client initialization: %v", err)
	}

	section := manager.FirewallSection{
		Description: name,
		DisplayName: name,
		SectionType: "LAYER3",
	}

	_, responseCode, err := nsxClient.ServicesApi.AddSection(nsxClient.Context, section, nil)
	if err != nil {
		return fmt.Errorf("Error during firewall section creation: %v", err)
	}

	if responseCode.StatusCode != http.StatusCreated {
		return fmt.Errorf("Unexpected status returned during firewall section creation: %v", responseCode.StatusCode)
	}
	return nil
}

func testAccDataSourceNsxtFirewallSectionDeleteByName(name string) error {
	nsxClient, err := testAccGetClient()
	if err != nil {
		return fmt.Errorf("Error during test client initialization: %v", err)
	}

	// Find the object by name
	objList, _, err := nsxClient.ServicesApi.ListSections(nsxClient.Context, nil)
	if err != nil {
		return fmt.Errorf("Error while reading Firewall sections: %v", err)
	}
	// go over the list to find the correct one
	for _, objInList := range objList.Results {
		if objInList.DisplayName == name {
			responseCode, err := nsxClient.ServicesApi.DeleteSection(nsxClient.Context, objInList.Id, nil)
			if err != nil {
				return fmt.Errorf("Error during firewall section deletion: %v", err)
			}

			if responseCode.StatusCode != http.StatusOK {
				return fmt.Errorf("Unexpected status returned during firewall section deletion: %v", responseCode.StatusCode)
			}
			return nil
		}
	}
	return fmt.Errorf("Error while deleting Firewall section '%s': firewall section not found", name)
}

func testAccNSXFirewallSectionReadTemplate(name string) string {
	return fmt.Sprintf(`
data "nsxt_firewall_section" "test" {
  display_name = "%s"
}`, name)
}

func testAccNSXNoFirewallSectionTemplate() string {
	return fmt.Sprintf(` `)
}
