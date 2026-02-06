// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/proxy"
)

func dataSourceNsxtProxyConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtProxyConfigRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceExtendedDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether proxy configuration is enabled",
				Computed:    true,
			},
			"scheme": {
				Type:        schema.TypeString,
				Description: "Proxy scheme (HTTP or HTTPS)",
				Computed:    true,
			},
			"host": {
				Type:        schema.TypeString,
				Description: "Proxy server host (IP address or FQDN)",
				Computed:    true,
			},
			"port": {
				Type:        schema.TypeInt,
				Description: "Proxy server port",
				Computed:    true,
			},
			"username": {
				Type:        schema.TypeString,
				Description: "Username for proxy authentication",
				Computed:    true,
			},
			"certificate_id": {
				Type:        schema.TypeString,
				Description: "Certificate ID for HTTPS proxy",
				Computed:    true,
			},
			"test_connection_url": {
				Type:        schema.TypeString,
				Description: "URL to test proxy connectivity",
				Computed:    true,
			},
		},
	}
}

func dataSourceNsxtProxyConfigRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	// Proxy Config is a singleton resource with fixed ID
	id := "TelemetryConfigIdentifier"

	log.Printf("[INFO] Reading Proxy Config data source with ID %s", id)

	// Create proxy config client
	client := proxy.NewConfigClient(connector)

	// Get proxy configuration
	proxyConfig, err := client.Get()
	if err != nil {
		return err
	}

	// Set the ID and attributes
	d.SetId(id)
	d.Set("display_name", proxyConfig.DisplayName)
	d.Set("description", proxyConfig.Description)
	d.Set("enabled", proxyConfig.Enabled)
	d.Set("scheme", proxyConfig.Scheme)
	d.Set("host", proxyConfig.Host)
	d.Set("port", proxyConfig.Port)
	d.Set("username", proxyConfig.Username)
	d.Set("certificate_id", proxyConfig.CertificateId)
	d.Set("test_connection_url", proxyConfig.TestConnectionUrl)
	d.Set("path", "/api/v1/proxy/config")

	return nil
}
