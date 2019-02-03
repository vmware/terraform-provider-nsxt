/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	api "github.com/vmware/go-vmware-nsxt"
	"os"
	"strings"
	"testing"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"nsxt": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	var requiredVariables = []string{"NSXT_USERNAME", "NSXT_PASSWORD", "NSXT_MANAGER_HOST", "NSXT_ALLOW_UNVERIFIED_SSL"}
	for _, element := range requiredVariables {
		if v := os.Getenv(element); v == "" {
			str := fmt.Sprintf("%s must be set for acceptance tests", element)
			t.Fatal(str)
		}
	}
}

func testAccGetClient() (*api.APIClient, error) {
	client, ok := testAccProvider.Meta().(*api.APIClient)
	if ok {
		return client, nil
	}

	// Try to create a temporary client using the tests configuration
	// This is necessary since the test PreCheck is called before the client is initialized.
	insecure := false
	if v := strings.ToLower(os.Getenv("NSXT_ALLOW_UNVERIFIED_SSL")); v != "false" && v != "0" {
		insecure = true
	}

	cfg := api.Configuration{
		BasePath:   "/api/v1",
		Host:       os.Getenv("NSXT_MANAGER_HOST"),
		Scheme:     "https",
		UserAgent:  "terraform-provider-nsxt/1.0",
		UserName:   os.Getenv("NSXT_USERNAME"),
		Password:   os.Getenv("NSXT_PASSWORD"),
		RemoteAuth: false,
		Insecure:   insecure,
	}

	return api.NewAPIClient(&cfg)
}

func testAccNSXVersion(t *testing.T, requiredVersion string) {
	client, err := testAccGetClient()
	if err != nil {
		t.Skipf("Skipping non-NSX provider. No NSX client")
		return
	}

	nsxVersion := getNSXVersion(client)
	if nsxVersion < requiredVersion {
		t.Skipf("This test can only run in NSX %s or above (Current version %s)", requiredVersion, nsxVersion)
	}
}
