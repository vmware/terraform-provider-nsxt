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

func TestAccDataSourceNsxtSwitchingProfile_basic(t *testing.T) {
	profileName := "terraform_test_profile"
	profileType := "QosSwitchingProfile"
	testResourceName := "data.nsxt_switching_profile.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccDataSourceNsxtSwitchingProfileDeleteByName(profileName)
		},
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
				),
			},
			{
				Config: testAccNSXSwitchingNoProfileTemplate(),
			},
		},
	})
}

func testAccDataSourceNsxtSwitchingProfileCreate(profileName string, profileType string) error {
	nsxClient, err := testAccGetClient()
	if err != nil {
		return fmt.Errorf("Error during test client initialization: %v", err)
	}
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

func testAccDataSourceNsxtSwitchingProfileDeleteByName(profileName string) error {
	nsxClient, err := testAccGetClient()
	if err != nil {
		return fmt.Errorf("Error during test client initialization: %v", err)
	}
	// Look up this profile
	localVarOptionals := make(map[string]interface{})
	localVarOptionals["includeSystemOwned"] = true
	objList, _, err := nsxClient.LogicalSwitchingApi.ListSwitchingProfiles(nsxClient.Context, localVarOptionals)
	if err != nil {
		return fmt.Errorf("Error while reading switching profiles: %v", err)
	}
	// go over the list to find the correct one
	for _, objInList := range objList.Results {
		if objInList.DisplayName == profileName {
			localVarOptionals := make(map[string]interface{})
			_, err := nsxClient.LogicalSwitchingApi.DeleteSwitchingProfile(nsxClient.Context, objInList.Id, localVarOptionals)
			if err != nil {
				return fmt.Errorf("Error during SwitchingProfile deletion: %v", err)
			}
			return nil
		}
	}
	return fmt.Errorf("Switching profile '%s' was not found", profileName)
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
