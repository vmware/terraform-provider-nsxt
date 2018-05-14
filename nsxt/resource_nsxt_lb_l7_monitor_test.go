/* Copyright © 2018 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/vmware/go-vmware-nsxt"
	"github.com/vmware/go-vmware-nsxt/trust"
	"net/http"
	"testing"
)

var caCertID string
var clientCertID string

func TestAccResourceNsxtLbHTTPMonitor_basic(t *testing.T) {
	testAccResourceNsxtLbL7MonitorBasic(t, "http")
}

func TestAccResourceNsxtLbHTTPSMonitor_basic(t *testing.T) {
	testAccResourceNsxtLbL7MonitorBasic(t, "https")
}

func TestAccResourceNsxtLbHTTPMonitor_importBasic(t *testing.T) {
	testAccResourceNsxtLbL7MonitorImport(t, "http")
}

func TestAccResourceNsxtLbHTTPSMonitor_importBasic(t *testing.T) {
	testAccResourceNsxtLbL7MonitorImport(t, "https")
}

func testAccResourceNsxtLbL7MonitorBasic(t *testing.T, protocol string) {
	name := "test-nsx-monitor"
	testResourceName := fmt.Sprintf("nsxt_lb_%s_monitor.test", protocol)
	requestMethod := "HEAD"
	updatedRequestMethod := "POST"
	body1 := "XXXXXXXXXXXXXXXXXXX"
	body2 := "YYYYYYYYYYYYYYYYYYY"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNSXLbL7MonitorCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNSXLbL7MonitorCreateTemplate(name, protocol, requestMethod, body1, body2),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXLbHTTPSMonitorExists(name, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "test description"),
					resource.TestCheckResourceAttr(testResourceName, "request_header.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "request_url", "/healthcheck"),
					resource.TestCheckResourceAttr(testResourceName, "request_method", requestMethod),
					resource.TestCheckResourceAttr(testResourceName, "request_body", body1),
					resource.TestCheckResourceAttr(testResourceName, "response_body", body2),
					resource.TestCheckResourceAttr(testResourceName, "response_status_codes.#", "3"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNSXLbL7MonitorCreateTemplate(name, protocol, updatedRequestMethod, body2, body1),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXLbHTTPSMonitorExists(name, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "test description"),
					resource.TestCheckResourceAttr(testResourceName, "request_header.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "request_url", "/healthcheck"),
					resource.TestCheckResourceAttr(testResourceName, "request_method", updatedRequestMethod),
					resource.TestCheckResourceAttr(testResourceName, "request_body", body2),
					resource.TestCheckResourceAttr(testResourceName, "response_body", body1),
					resource.TestCheckResourceAttr(testResourceName, "response_status_codes.#", "3"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
		},
	})
}

func TestAccResourceNsxtLbHTTPSMonitor_withAuth(t *testing.T) {
	name := "test-nsx-monitor"
	testResourceName := "nsxt_lb_https_monitor.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			testAccNSXLbHTTPSMonitorDeleteCerts(t)
			return testAccNSXLbL7MonitorCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testAccNSXLbHTTPSMonitorCreateCerts(t)
				},

				Config: testAccNSXLbHTTPSMonitorCreateTemplateWithAuth(name),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXLbHTTPSMonitorExists(name, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "test description"),
					resource.TestCheckResourceAttr(testResourceName, "ciphers.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "protocols.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "certificate_chain_depth", "2"),
					resource.TestCheckResourceAttr(testResourceName, "server_auth", "REQUIRED"),
					resource.TestCheckResourceAttr(testResourceName, "server_auth_ca_ids.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "client_certificate_id"),
				),
			},
			{Config: " "},
		},
	})
}

func testAccResourceNsxtLbL7MonitorImport(t *testing.T, protocol string) {
	name := "test-nsx-monitor"
	testResourceName := fmt.Sprintf("nsxt_lb_%s_monitor.test", protocol)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNSXLbL7MonitorCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNSXLbL7MonitorCreateTemplateTrivial(name, protocol),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNSXLbHTTPSMonitorExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		nsxClient := testAccProvider.Meta().(*nsxt.APIClient)
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("NSX LB monitor resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("NSX LB monitor resource ID not set in resources ")
		}

		monitor, response, err := nsxClient.ServicesApi.ReadLoadBalancerMonitor(nsxClient.Context, resourceID)
		if err != nil {
			return fmt.Errorf("Error while retrieving LB monitor with ID %s. Error: %v", resourceID, err)
		}

		if response.StatusCode != http.StatusOK {
			return fmt.Errorf("Error while checking if LB monitor %s exists. HTTP return code was %d", resourceID, response.StatusCode)
		}

		if displayName == monitor.DisplayName {
			return nil
		}
		return fmt.Errorf("NSX LB monitor %s wasn't found", displayName)
	}
}

func testAccNSXLbL7MonitorCheckDestroy(state *terraform.State, displayName string) error {
	nsxClient := testAccProvider.Meta().(*nsxt.APIClient)
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_lb_https_monitor" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		monitor, response, err := nsxClient.ServicesApi.ReadLoadBalancerMonitor(nsxClient.Context, resourceID)
		if err != nil {
			if response.StatusCode != http.StatusOK {
				return nil
			}
			return fmt.Errorf("Error while retrieving LB https monitor with ID %s. Error: %v", resourceID, err)
		}

		if displayName == monitor.DisplayName {
			return fmt.Errorf("NSX LB https monitor %s still exists", displayName)
		}
	}
	return nil
}

func testAccNSXLbHTTPSMonitorCreateCert(t *testing.T, name string, certPem string, certPK string, certType string) string {
	nsxClient, err := testAccGetClient()
	if err != nil {
		t.Fatal(err)
	}

	object := trust.TrustObjectData{
		DisplayName:  name,
		ResourceType: certType,
		PemEncoded:   certPem,
		PrivateKey:   certPK,
	}

	certList, response, err := nsxClient.NsxComponentAdministrationApi.AddCertificateImport(nsxClient.Context, object)

	if err != nil {
		t.Fatal(fmt.Sprintf("Error while creating %s certificate. Error: %v", certType, err))
	}
	if response.StatusCode != http.StatusCreated {
		t.Fatal(fmt.Errorf("Error while creating %s certificate. HTTP return code %d", certType, response.StatusCode))
	}
	certID := ""
	for _, cert := range certList.Results {
		certID = cert.Id
	}

	return certID
}

func testAccNSXLbHTTPSMonitorDeleteCert(t *testing.T, id string) {
	nsxClient, err := testAccGetClient()
	if err != nil {
		t.Fatal(err)
	}

	response, err := nsxClient.NsxComponentAdministrationApi.DeleteCertificate(nsxClient.Context, id)

	if err != nil {
		t.Fatal(fmt.Sprintf("Error while deleting certificate %s. Error: %v", id, err))
	}
	if response.StatusCode != http.StatusOK {
		t.Fatal(fmt.Errorf("Error while deleting certificate %s. HTTP return code %d", id, response.StatusCode))
	}
}

func testAccNSXLbHTTPSMonitorCreateCerts(t *testing.T) {

	caPem := `-----BEGIN CERTIFICATE-----\nMIICVjCCAb8CAg37MA0GCSqGSIb3DQEBBQUAMIGbMQswCQYDVQQGEwJKUDEOMAwGA1UECBMFVG9reW8xEDAOBgNVBAcTB0NodW8ta3UxETAPBgNVBAoTCEZyYW5rNEREMRgwFgYDVQQLEw9XZWJDZXJ0IFN1cHBvcnQxGDAWBgNVBAMTD0ZyYW5rNEREIFdlYiBDQTEjMCEGCSqGSIb3DQEJARYUc3VwcG9ydEBmcmFuazRkZC5jb20wHhcNMTIwODIyMDUyNzIzWhcNMTcwODIxMDUyNzIzWjBKMQswCQYDVQQGEwJKUDEOMAwGA1UECAwFVG9reW8xETAPBgNVBAoMCEZyYW5rNEREMRgwFgYDVQQDDA93d3cuZXhhbXBsZS5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJAoGBAMYBBrx5PlP0WNI/ZdzD+6Pktmurn+F2kQYbtc7XQh8/LTBvCo+P6iZoLEmUA9e7EXLRxgU1CVqeAi7QcAn9MwBlc8ksFJHB0rtf9pmf8Oza9E0Bynlq/4/Kb1x+d+AyhL7oK9tQwB24uHOueHi1C/iVv8CSWKiYe6hzN1txYe8rAgMBAAEwDQYJKoZIhvcNAQEFBQADgYEAASPdjigJkXCqKWpnZ/Oc75EUcMi6HztaW8abUMlYXPIgkV2F7YanHOB7K4f7OOLjiz8DTPFfjC9UeuErhaA/zzWi8ewMTFZW/WshOrm3fNvcMrMLKtH534JKvcdMg6qIdjTFINIrevnAhf0cwULaebn+lMs8Pdl7y37+sfluVok=\n-----END CERTIFICATE-----\n`

	caCertID = testAccNSXLbHTTPSMonitorCreateCert(t, "ca", caPem, "", "certificate_ca")

	clientCertPem := "-----BEGIN CERTIFICATE-----\n" +
		"MIID6jCCAtKgAwIBAgIJAOtKKdMP6oZcMA0GCSqGSIb3DQEBCwUAMHYxCzAJBgNV\n" +
		"BAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBX\n" +
		"aWRnaXRzIFB0eSBMdGQxETAPBgNVBAMMCHRlc3QuY29tMRwwGgYJKoZIhvcNAQkB\n" +
		"Fg10ZXN0QHRlc3QuY29tMB4XDTE4MDUwNzE3MTkxOVoXDTE5MDUwNzE3MTkxOVow\n" +
		"djELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoMGElu\n" +
		"dGVybmV0IFdpZGdpdHMgUHR5IEx0ZDERMA8GA1UEAwwIdGVzdC5jb20xHDAaBgkq\n" +
		"hkiG9w0BCQEWDXRlc3RAdGVzdC5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAw\n" +
		"ggEKAoIBAQC8GhgSX8nV0QEStgQ208PDXsl62xefAWVYMVjF9+xUWHY6zp/sTTQA\n" +
		"uFGD6DBiyljYYShj85ha7N3m70Xmur5Bbpa7ca010R0oCLPQRcySvsom190OfeUG\n" +
		"LW396gZUjiY9+8OQqI8BrTcNkI1DOpC5fBbbXkvLN4pOIf+/JqZNtr1igQBv2g+E\n" +
		"FCiPxDuWPOagi2xeEQv3SLLBhfF92UAlV6JRziCwaILlC3Zn/6B8E8rpuqVxUShW\n" +
		"wDmLowoT4BZP6/OHRneyEvVi2L/Ucsdk3nwUTe5Q4ojHY21ftW6by6uXYvK4IePC\n" +
		"ZXeTjdNtRX4MaAtJ4iLY7E/9d0BdTNRvAgMBAAGjezB5MAkGA1UdEwQCMAAwLAYJ\n" +
		"YIZIAYb4QgENBB8WHU9wZW5TU0wgR2VuZXJhdGVkIENlcnRpZmljYXRlMB0GA1Ud\n" +
		"DgQWBBTP56rycJh9hpWKOmvTzYg3YEkqiDAfBgNVHSMEGDAWgBTP56rycJh9hpWK\n" +
		"OmvTzYg3YEkqiDANBgkqhkiG9w0BAQsFAAOCAQEAZ90xoDf4gRbh7/PHxbokP69K\n" +
		"uqt7s78JDKrGxCDwiezUhZrdOBwi2r/sg4cWi43uuwCfjGNCd824EQYRaSCjanAn\n" +
		"5OH14KshCOBi66CaWDzJK6v4X/hbrKtUmXvbvUjrQCEHVuLueEQfvJB5/8O0dpEA\n" +
		"xF4MhnSaF2Id5O7tlXnFoXZh0YI9QnTwHLQ4L9+3PS5LqWd0peV1XqgWy0CXwjZF\n" +
		"nEpHq+TGDwLRoAgnoBGrbaFJRmvm+iVU4J76AtV7B3keckVMyMIeBR9CWB7kDm64\n" +
		"86qiRcEGN7V5mMJtxF49l0F01qdOgrictZRf+gMMrtGmX4KkZ6DKrl278HPs7A==\n" +
		"-----END CERTIFICATE-----"

	clientPKPem := "-----BEGIN PRIVATE KEY-----\n" +
		"MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC8GhgSX8nV0QES\n" +
		"tgQ208PDXsl62xefAWVYMVjF9+xUWHY6zp/sTTQAuFGD6DBiyljYYShj85ha7N3m\n" +
		"70Xmur5Bbpa7ca010R0oCLPQRcySvsom190OfeUGLW396gZUjiY9+8OQqI8BrTcN\n" +
		"kI1DOpC5fBbbXkvLN4pOIf+/JqZNtr1igQBv2g+EFCiPxDuWPOagi2xeEQv3SLLB\n" +
		"hfF92UAlV6JRziCwaILlC3Zn/6B8E8rpuqVxUShWwDmLowoT4BZP6/OHRneyEvVi\n" +
		"2L/Ucsdk3nwUTe5Q4ojHY21ftW6by6uXYvK4IePCZXeTjdNtRX4MaAtJ4iLY7E/9\n" +
		"d0BdTNRvAgMBAAECggEAYS8KKNQcr7/gUg6AduNKTXJ3nmX7+kb6WWqFdyL0k09x\n" +
		"JkkDD0+EAqs4UnJwLVpmNMVd3OZKXQ/sRhhxgRgSnDPK5OWCnD+CVODKJl0pqNey\n" +
		"EgeNSqN45IwsO/fhdWZME9Iz5FVyLWeU/gklMwrbIzodhRFfD4uOhXfDbrtFSPix\n" +
		"KGzhDNCn3z/Y9Sml6tUMgCLPisQbxWcQ+bbTHCgXXEdK2yG45j4q3FqBGy8IExou\n" +
		"cxw0aWu8cIqIGlXBQA5SZ78csBisKxEmlHjkV1i7Jn2DIARUluU/nmYKGC4ElWGp\n" +
		"yuURRJhcjO17MkCEr+niTxnTptkHi9e1p4d/k6HCgQKBgQD0Uj/JPr9vIM6Mt3vm\n" +
		"kK0GSTfpMUDDwaAMR+Pip6qocEYD5nZ6rnSuXIqlKyP/YxtQqsDqFdYHbPJet/Lg\n" +
		"3HFKerQ7xHiral95rsP0TWiRtxCHaF1xahgQWlVaYqGDeBcL7EavGJwQy3crFiVo\n" +
		"6wHwbVfNV/hoXHHK7pX4J8M3sQKBgQDFF+P/gM86QrvDbJnH/7PuF6QzcdvvtacH\n" +
		"i+5aTG9YQaUGldQIx2Fjpn48V73/YGK04DTrXqACjE2XxgEu0/B0wAGkOwGZztAr\n" +
		"q54O4AEWgVAyUukarkjdGtgeUAf/UMeXf7D2tGhjyNtORrOyZCdKS3kZ3GfZuEXj\n" +
		"FjhbJbX2HwKBgCMh8Kovq7d/MDRr7hUpmLfer3uI6Zc8sJcTf2GIWrH98xN8gG0D\n" +
		"ySOJiyZVHcgLqFHhO/xtR2mp8PBN408SY/ghzOkLR47erPwCdYsb1n2dpXLTPxyf\n" +
		"9PXlB4EHzdHp4uaEA2YKU+bWWzyG4rpDkPPRxV5x1/ap1HMp+8bDcP8BAoGAJTdQ\n" +
		"pxNUjgTB3bHpC9ndyOyP5eLvC8F6S7OBi215bOngVnD+O7YiTqXGmnBbARjbKppX\n" +
		"g8Y3YqPJlwodeREuC22iIbe+oqNprYVXcCmeKvi6Avai65XTTmTeQEMOb4h6V8IV\n" +
		"0U/ZklYACzTQg7Pjs2Sy9k4nEfZ4w9uTQqrJRDMCgYEAvgB29WCwEJ81pjBIGiYw\n" +
		"FbqLxQtIFiXKvRW+Qr1EKcV89kUYvsH3+TO/hZpf6rBtRf9yLJvw2fR4tVouowbn\n" +
		"/w+mKvKmcm11NOuobla4w9gztptqeUzcW03JmHmcnyHlnnUwpsU/XcgnsToV6kJB\n" +
		"bS6lN3HXkJnFnAzK2BKoZCA=\n" +
		"-----END PRIVATE KEY-----\n"

	clientCertID = testAccNSXLbHTTPSMonitorCreateCert(t, "client", clientCertPem, clientPKPem, "certificate_self_signed")
}

func testAccNSXLbHTTPSMonitorDeleteCerts(t *testing.T) {

	testAccNSXLbHTTPSMonitorDeleteCert(t, caCertID)
	testAccNSXLbHTTPSMonitorDeleteCert(t, clientCertID)
}

func testAccNSXLbL7MonitorCreateTemplate(name string, protocol string, requestMethod string, requestBody string, responseBody string) string {
	return fmt.Sprintf(`
resource "nsxt_lb_%s_monitor" "test" {
  description       = "test description"
  display_name      = "%s"
  request_header {
      name = "header1"
      value = "value1"
  }
  request_header {
      name = "header2"
      value = "value2"
  }
  request_method = "%s"
  request_url = "/healthcheck"
  request_body = "%s"
  response_body = "%s"
  response_status_codes = [200, 304, 404]
  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}
`, protocol, name, requestMethod, requestBody, responseBody)
}

func testAccNSXLbL7MonitorCreateTemplateTrivial(name string, protocol string) string {
	return fmt.Sprintf(`
resource "nsxt_lb_%s_monitor" "test" {
  description = "test description"
}
`, protocol)
}

func testAccNSXLbHTTPSMonitorCreateTemplateWithAuth(name string) string {
	return fmt.Sprintf(`
data "nsxt_certificate" "ca" {
  display_name = "ca"
}

data "nsxt_certificate" "client" {
  display_name = "client"
}

resource "nsxt_lb_https_monitor" "test" {
  description             = "test description"
  display_name            = "%s"
  server_auth             = "REQUIRED"
  server_auth_ca_ids      = ["${data.nsxt_certificate.ca.id}"]
  certificate_chain_depth = 2
  client_certificate_id   = "${data.nsxt_certificate.client.id}"
  protocols               = ["TLS_V1_2"]
  ciphers                 = ["TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256", "TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA384"]
}
`, name)
}
