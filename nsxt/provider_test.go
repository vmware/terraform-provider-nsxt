/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"os"
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
	if v := os.Getenv("NSX_USERNAME"); v == "" {
		t.Fatal("NSX_USERNAME must be set for acceptance tests")
	}

	if v := os.Getenv("NSX_PASSWORD"); v == "" {
		t.Fatal("NSX_PASSWORD must be set for acceptance tests")
	}

	if v := os.Getenv("NSX_MANAGER_HOST"); v == "" {
		t.Fatal("NSX_MANAGER_HOST must be set for acceptance tests")
	}

	if v := os.Getenv("NSX_INSECURE"); v == "" {
		t.Fatal("NSX_INSECURE must be set for acceptance tests")
	}
}
