/* Copyright © 2020 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"net/http"
	"testing"
)

const (
	testAccLicenseName = "nsxt_license.test"
)

func TestAccResourceNsxtLicense(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccOnlyLocalManager(t); testAccPreCheck(t); testAccEnvDefined(t, "NSXT_LICENSE_KEY") },
		Providers:    testAccProviders,
		CheckDestroy: testAccNSXLicenseCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNSXLicenseCreateTemplate(),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXLicenseExists(),
					resource.TestCheckResourceAttrSet(testAccLicenseName, "license_key"),
				),
			},
		},
	})
}

func testAccNSXLicenseExists() resource.TestCheckFunc {
	return func(state *terraform.State) error {
		exists, err := checkNSXLicenseExists(state)
		if err != nil {
			return err
		}
		if !*exists {
			return fmt.Errorf("License %s does not exist", testAccLicenseName)
		}
		return nil
	}
}

func checkNSXLicenseExists(state *terraform.State) (*bool, error) {
	exists := false
	licenseKey, err := getLicenseKeyByResourceName(state)
	if err != nil {
		return &exists, err
	}

	nsxClient := testAccProvider.Meta().(nsxtClients).NsxtClient
	if nsxClient == nil {
		return nil, resourceNotSupportedError()
	}
	license, responseCode, err := nsxClient.LicensingApi.GetLicenseByKey(nsxClient.Context, licenseKey)
	if responseCode.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error while checking if license exists. license key: %s, Error: %v", licenseKey, err)
	}
	if license.LicenseKey == licenseKey {
		exists = true
	}

	return &exists, nil
}

func getLicenseKeyByResourceName(state *terraform.State) (string, error) {
	rsLicense, ok := state.RootModule().Resources[testAccLicenseName]
	if !ok {
		return "", fmt.Errorf("license resource %s not found in resources", testAccLicenseName)
	}

	licenseKey := rsLicense.Primary.ID
	if licenseKey == "" {
		return "", fmt.Errorf("license resource ID not set in resources ")
	}
	return licenseKey, nil
}

func testAccNSXLicenseCreateTemplate() string {
	return fmt.Sprintf(`
resource "nsxt_license" "test" {
  license_key = "%s"
}`, getTestLicenseKey())
}

func testAccNSXLicenseCheckDestroy(state *terraform.State) error {
	nsxClient := testAccProvider.Meta().(nsxtClients).NsxtClient
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_license" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		license, responseCode, err := nsxClient.LicensingApi.GetLicenseByKey(nsxClient.Context, resourceID)
		if err != nil {
			if responseCode.StatusCode != http.StatusOK {
				return nil
			}
			return fmt.Errorf("Error while retrieving License Key %s. Error: %v", resourceID, err)
		}

		if getTestLicenseKey() == license.LicenseKey {
			return fmt.Errorf("License key %s still exists", getTestLicenseKey())
		}
	}
	return nil
}
