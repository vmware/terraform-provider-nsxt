/* Copyright © 2020 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/vmware/go-vmware-nsxt/licensing"
	"log"
	"net/http"
	"regexp"
)

func resourceNsxtLicense() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtLicenseCreate,
		Read:   resourceNsxtLicenseRead,
		Update: resourceNsxtLicenseCreate,
		Delete: resourceNsxtLicenseDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"license_key": {
				Type:        schema.TypeString,
				Description: "license key",
				Required:    true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(
						"^[A-Z0-9]{5}-[A-Z0-9]{5}-[A-Z0-9]{5}-[A-Z0-9]{5}-[A-Z0-9]{5}$"),
					"Must be a valid nsx license key matching: ^[A-Z0-9]{5}-[A-Z0-9]{5}-[A-Z0-9]{5}-[A-Z0-9]{5}-[A-Z0-9]{5}$"),
			},
		},
	}
}

func resourceNsxtLicenseCreate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(nsxtClients).NsxtClient
	if nsxClient == nil {
		return resourceNotSupportedError()
	}
	key := d.Get("license_key").(string)
	license := licensing.License{
		LicenseKey: key,
	}

	license, resp, err := nsxClient.LicensingApi.CreateLicense(nsxClient.Context, license)

	if err != nil {
		return fmt.Errorf("Error during license create: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected status returned during license create: %v", resp.StatusCode)
	}
	d.SetId(license.LicenseKey)

	return resourceNsxtLicenseRead(d, m)
}

func resourceNsxtLicenseRead(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(nsxtClients).NsxtClient
	if nsxClient == nil {
		return resourceNotSupportedError()
	}
	key := d.Get("license_key").(string)

	_, resp, err := nsxClient.LicensingApi.GetLicenseByKey(nsxClient.Context, key)

	if resp != nil && resp.StatusCode == http.StatusBadRequest {
		log.Printf("[DEBUG] license %s not found", key)
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error during license read: %v", err)
	}

	return nil
}

func resourceNsxtLicenseDelete(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(nsxtClients).NsxtClient
	if nsxClient == nil {
		return resourceNotSupportedError()
	}
	key := d.Get("license_key").(string)

	resp, err := nsxClient.LicensingApi.DeleteLicense(nsxClient.Context, key)
	if err != nil {
		return fmt.Errorf("Error during license delete: %v", err)
	}

	if resp.StatusCode == http.StatusBadRequest {
		log.Printf("[DEBUG] license %s not found", key)
		d.SetId("")
	}
	return nil
}
