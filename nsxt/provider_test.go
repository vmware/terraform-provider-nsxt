/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"crypto/tls"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	api "github.com/vmware/go-vmware-nsxt"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/security"
	"net/http"
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
	if os.Getenv("NSXT_MANAGER_HOST") == "" {
		return nil, fmt.Errorf("NSXT_MANAGER_HOST is not set in environment")
	}

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
	if nsxVersion == "" {
		client, err := testAccGetClient()
		if err != nil {
			t.Skipf("Skipping non-NSX provider. No NSX client")
			return
		}

		initNSXVersion(client)
	}

	if nsxVersionLower(requiredVersion) {
		t.Skipf("This test can only run in NSX %s or above (Current version %s)", requiredVersion, nsxVersion)
	}
}

func testAccNSXVersionLessThan(t *testing.T, requiredVersion string) {
	if nsxVersion == "" {
		client, err := testAccGetClient()
		if err != nil {
			t.Skipf("Skipping non-NSX provider. No NSX client")
			return
		}

		initNSXVersion(client)
	}

	if nsxVersionHigherOrEqual(requiredVersion) {
		t.Skipf("This test can only run in NSX below %s (Current version %s)", requiredVersion, nsxVersion)
	}
}

func testAccGetPolicyConnector() (*client.RestConnector, error) {
	if os.Getenv("NSXT_MANAGER_HOST") == "" {
		return nil, fmt.Errorf("NSXT_MANAGER_HOST is not set in environment")
	}

	// Try to create a temporary client using the tests configuration
	// This is necessary since the test PreCheck is called before the client is initialized.
	insecure := false
	if v := strings.ToLower(os.Getenv("NSXT_ALLOW_UNVERIFIED_SSL")); v != "false" && v != "0" {
		insecure = true
	}

	hostIP := os.Getenv("NSXT_MANAGER_HOST")
	host := fmt.Sprintf("https://%s", hostIP)
	username := os.Getenv("NSXT_USERNAME")
	password := os.Getenv("NSXT_PASSWORD")

	//TODO: add error handling
	securityCtx := core.NewSecurityContextImpl()
	securityCtx.SetProperty(security.AUTHENTICATION_SCHEME_ID, security.USER_PASSWORD_SCHEME_ID)
	securityCtx.SetProperty(security.USER_KEY, username)
	securityCtx.SetProperty(security.PASSWORD_KEY, password)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
		Proxy:           http.ProxyFromEnvironment,
	}
	httpClient := http.Client{Transport: tr}
	connector := client.NewRestConnector(host, httpClient)
	connector.SetSecurityContext(securityCtx)

	return connector, nil
}
