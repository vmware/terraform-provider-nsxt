/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	api "github.com/vmware/go-vmware-nsxt"
	"log"
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
	var requiredVariables = []string{"NSX_USERNAME", "NSX_PASSWORD", "NSX_MANAGER_HOST", "NSX_ALLOW_UNVERIFIED_SSL"}
	for _, element := range requiredVariables {
		if v := os.Getenv(element); v == "" {
			str = fmt.Sprintf("%s must be set for acceptance tests", element)
			t.Fatal(str)
		}
	}
}

func testAccGetClient() *api.APIClient {
	client, ok := testAccProvider.Meta().(*api.APIClient)
	if ok {
		return client
	}
	// Try to create a temporary client using the tests configuration
	// This is necessary since the test PreCheck is called before the client is initialized.
	insecure := false
	if strings.ToLower(os.Getenv("NSX_ALLOW_UNVERIFIED_SSL")) == "true" {
		insecure = true
	}
	cfg := api.Configuration{
		BasePath:  "/api/v1",
		Host:      os.Getenv("NSX_MANAGER_HOST"),
		Scheme:    "https",
		UserAgent: "terraform-provider-nsxt/1.0",
		UserName:  os.Getenv("NSX_USERNAME"),
		Password:  os.Getenv("NSX_PASSWORD"),
		Insecure:  insecure,
	}

	newClient, err := api.NewAPIClient(&cfg)
	if err != nil {
		log.Printf("Failed to create a test client: %s", err)
		return nil
	}
	return newClient
}

func testAccNSXVersion(t *testing.T, requiredVersion string) {
	client := testAccGetClient()
	if client == nil {
		t.Skipf("Skipping non-NSX provider. No NSX client")
		return
	}

	nsxVersion := getNSXVersion(client)
	if nsxVersion < requiredVersion {
		t.Skipf("This test can only run in NSX %s or above (Current version %s)", requiredVersion, nsxVersion)
	}
}
