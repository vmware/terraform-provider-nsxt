// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/proxy"
)

func resourceNsxtProxyConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtProxyConfigCreate,
		Read:   resourceNsxtProxyConfigRead,
		Update: resourceNsxtProxyConfigUpdate,
		Delete: resourceNsxtProxyConfigDelete,
		Importer: &schema.ResourceImporter{
			State: resourceNsxtProxyConfigImporter,
		},
		Description: "Global internet proxy configuration. This is a singleton resource - only one instance should exist per NSX-T environment.",

		Schema: map[string]*schema.Schema{
			"nsx_id": getNsxIDSchema(),
			"path":   getPathSchema(),
			"display_name": {
				Type:        schema.TypeString,
				Description: "Display name (read-only, managed by NSX)",
				Computed:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description (not supported by NSX for proxy config)",
				Optional:    true,
				Computed:    true,
			},
			"revision": getRevisionSchema(),
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Enable proxy configuration",
				Optional:    true,
				Default:     false,
			},
			"scheme": {
				Type:         schema.TypeString,
				Description:  "Proxy scheme (HTTP or HTTPS)",
				Optional:     true,
				Default:      "HTTP",
				ValidateFunc: validation.StringInSlice([]string{"HTTP", "HTTPS"}, false),
			},
			"host": {
				Type:        schema.TypeString,
				Description: "Proxy server host (IP address or FQDN)",
				Optional:    true,
			},
			"port": {
				Type:         schema.TypeInt,
				Description:  "Proxy server port",
				Optional:     true,
				Default:      3128,
				ValidateFunc: validation.IntBetween(1, 65535),
			},
			"username": {
				Type:        schema.TypeString,
				Description: "Username for proxy authentication",
				Optional:    true,
			},
			"password": {
				Type:        schema.TypeString,
				Description: "Password for proxy authentication",
				Optional:    true,
				Sensitive:   true,
			},
			"certificate_id": {
				Type:        schema.TypeString,
				Description: "Certificate ID for HTTPS proxy (from trust-management API). Required when scheme is HTTPS.",
				Optional:    true,
			},
			"test_connection_url": {
				Type:        schema.TypeString,
				Description: "URL to test proxy connectivity (e.g., https://www.vmware.com)",
				Optional:    true,
				Default:     "https://www.vmware.com",
			},
		},
	}
}

func resourceNsxtProxyConfigCreate(d *schema.ResourceData, m interface{}) error {
	// Proxy Config is a singleton resource - use fixed ID
	id := "TelemetryConfigIdentifier"

	// Check if user provided a custom ID, but warn that it will be ignored
	if userID := d.Get("nsx_id").(string); userID != "" && userID != id {
		log.Printf("[WARN] Proxy Config is a singleton resource. Custom nsx_id '%s' will be ignored, using '%s'", userID, id)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	log.Printf("[INFO] Creating Proxy Config with ID %s", id)

	// Perform the update operation to configure the proxy
	return resourceNsxtProxyConfigUpdate(d, m)
}

func resourceNsxtProxyConfigRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Proxy Config ID")
	}

	log.Printf("[INFO] Reading Proxy Config with ID %s", id)

	// Create proxy config client
	client := proxy.NewConfigClient(connector)

	// Get proxy configuration
	proxyConfig, err := client.Get()
	if err != nil {
		return handleReadError(d, "ProxyConfig", id, err)
	}

	// Set attributes from API response
	d.Set("nsx_id", id)
	d.Set("display_name", proxyConfig.DisplayName)
	d.Set("description", proxyConfig.Description)
	d.Set("revision", proxyConfig.Revision)
	d.Set("enabled", proxyConfig.Enabled)
	d.Set("scheme", proxyConfig.Scheme)
	d.Set("host", proxyConfig.Host)
	d.Set("port", proxyConfig.Port)
	d.Set("username", proxyConfig.Username)
	d.Set("certificate_id", proxyConfig.CertificateId)
	d.Set("test_connection_url", proxyConfig.TestConnectionUrl)

	// Set path
	path := "/api/v1/proxy/config"
	d.Set("path", path)

	// Note: Password is write-only and not returned by the API
	// Note: Tags are not supported by this NSX Manager API endpoint

	return nil
}

func resourceNsxtProxyConfigUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Proxy Config ID")
	}

	// Create proxy config client
	client := proxy.NewConfigClient(connector)

	// First, read the current configuration to get the revision
	currentConfig, err := client.Get()
	if err != nil {
		return fmt.Errorf("Error reading current Proxy Config: %v", err)
	}

	// Build the proxy configuration
	enabled := d.Get("enabled").(bool)
	scheme := d.Get("scheme").(string)
	host := d.Get("host").(string)
	port := int64(d.Get("port").(int))
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	certificateID := d.Get("certificate_id").(string)
	testConnectionURL := d.Get("test_connection_url").(string)

	log.Printf("[INFO] Updating Proxy Config with ID %s", id)
	log.Printf("[DEBUG] Proxy Config - Enabled: %v, Scheme: %s, Host: %s, Port: %d, Username: %s",
		enabled, scheme, host, port, username)

	// Validation: if proxy is enabled, host is required
	if enabled && host == "" {
		return fmt.Errorf("host is required when enabled is true")
	}

	// NSX requires a valid host even when disabled, provide a placeholder
	if !enabled && host == "" {
		host = "proxy.example.com"
	}

	// Build the Proxy object
	// Note: NSX doesn't allow custom display_name for proxy config singleton
	// We preserve the existing display_name from NSX
	proxyConfig := model.Proxy{
		Enabled:           &enabled,
		Scheme:            &scheme,
		Host:              &host,
		Port:              &port,
		TestConnectionUrl: &testConnectionURL,
		Revision:          currentConfig.Revision,
		DisplayName:       currentConfig.DisplayName, // Preserve NSX display name
		Id:                currentConfig.Id,          // Preserve NSX ID
	}

	// Set optional fields
	if username != "" {
		proxyConfig.Username = &username
	}
	if password != "" {
		proxyConfig.Password = &password
	}
	if certificateID != "" {
		proxyConfig.CertificateId = &certificateID
	}

	// Set description if provided
	description := d.Get("description").(string)
	if description != "" {
		proxyConfig.Description = &description
	}

	// Note: Tags are not supported by this NSX Manager API endpoint

	// Update proxy configuration
	_, err = client.Update(proxyConfig)
	if err != nil {
		return handleUpdateError("ProxyConfig", id, err)
	}

	return resourceNsxtProxyConfigRead(d, m)
}

func resourceNsxtProxyConfigDelete(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Proxy Config ID")
	}

	log.Printf("[INFO] Deleting Proxy Config with ID %s", id)
	log.Printf("[WARN] Proxy Config is a singleton resource and cannot be truly deleted from NSX-T")
	log.Printf("[INFO] The proxy configuration will be disabled instead of deleted")

	// Proxy Config cannot be truly deleted from NSX-T, only disabled
	// Create proxy config client
	client := proxy.NewConfigClient(connector)

	// Read current config to get revision
	currentConfig, err := client.Get()
	if err != nil {
		log.Printf("[WARN] Error reading Proxy Config during delete: %v", err)
		// Continue with deletion from state even if read fails
		d.SetId("")
		return nil
	}

	// Disable the proxy
	enabled := false
	scheme := "HTTP"
	host := ""
	port := int64(3128)
	testURL := "https://www.vmware.com"

	proxyConfig := model.Proxy{
		Enabled:           &enabled,
		Scheme:            &scheme,
		Host:              &host,
		Port:              &port,
		TestConnectionUrl: &testURL,
		Revision:          currentConfig.Revision,
	}

	// Note: NSX requires host and port even when disabled, so we provide defaults
	_, err = client.Update(proxyConfig)
	if err != nil {
		log.Printf("[WARN] Error disabling Proxy Config: %v", err)
		// Continue with state cleanup even if disable fails
	}

	// Clear the Terraform state
	d.SetId("")

	return nil
}

func resourceNsxtProxyConfigImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	// Proxy config is a singleton with a fixed ID
	// Accept any import ID and set it to the correct singleton ID
	id := "TelemetryConfigIdentifier"
	d.SetId(id)

	// Read the current configuration
	err := resourceNsxtProxyConfigRead(d, m)
	if err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
