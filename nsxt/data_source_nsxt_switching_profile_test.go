/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/vmware/go-vmware-nsxt/manager"
	"net/http"
	"testing"
)

func TestAccDataSourceNsxtSwitchingProfile_basic(t *testing.T) {
	profileName := "terraform_test_profile"
	profileType := "QosSwitchingProfile"
	testResourceName := "data.nsxt_switching_profile.test"
	var s *terraform.State

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if err := testAccDataSourceNsxtSwitchingProfileCreate(profileName, profileType); err != nil {
						panic(err)
					}
				},
				Config: testAccNSXSwitchingProfileReadTemplate(profileName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", profileName),
					resource.TestCheckResourceAttr(testResourceName, "description", profileName),
					resource.TestCheckResourceAttr(testResourceName, "resource_type", profileType),
					copyStatePtr(&s),
				),
			},
			{
				PreConfig: func() {
					if err := testAccDataSourceNsxtSwitchingProfileDelete(s, testResourceName); err != nil {
						panic(err)
					}
				},
				Config: testAccNSXSwitchingNoProfileTemplate(),
			},
		},
	})
}

func testAccDataSourceNsxtSwitchingProfileCreate(profileName string, profileType string) error {
	nsxClient := testAccGetClient()
	profile := manager.BaseSwitchingProfile{
		DisplayName:  profileName,
		ResourceType: profileType,
		Description:  profileName,
	}
	profile, responseCode, err := nsxClient.LogicalSwitchingApi.CreateSwitchingProfile(nsxClient.Context, profile)
	if err != nil {
		return fmt.Errorf("Error during SwitchingProfile creation: %v", err)
	}

	if responseCode.StatusCode != http.StatusCreated {
		return fmt.Errorf("Unexpected status returned during SwitchingProfile creation: %v", responseCode.StatusCode)
	}
	return nil
}

func testAccDataSourceNsxtSwitchingProfileDelete(state *terraform.State, resourceName string) error {
	rs, ok := state.RootModule().Resources[resourceName]
	if !ok {
		return fmt.Errorf("NSX SwitchingProfile data source %s not found", resourceName)
	}

	resourceID := rs.Primary.ID
	if resourceID == "" {
		return fmt.Errorf("NSX SwitchingProfile data source ID not set")
	}
	nsxClient := testAccGetClient()
	localVarOptionals := make(map[string]interface{})
	responseCode, err := nsxClient.LogicalSwitchingApi.DeleteSwitchingProfile(nsxClient.Context, resourceID, localVarOptionals)
	if err != nil {
		return fmt.Errorf("Error during SwitchingProfile deletion: %v", err)
	}

	if responseCode.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected status returned during SwitchingProfile deletion: %v", responseCode.StatusCode)
	}
	return nil
}

func testAccNSXSwitchingProfileReadTemplate(profileName string) string {
	return fmt.Sprintf(`
data "nsxt_switching_profile" "test" {
     display_name = "%s"
}`, profileName)
}

func testAccNSXSwitchingNoProfileTemplate() string {
	return fmt.Sprintf(` `)
}
