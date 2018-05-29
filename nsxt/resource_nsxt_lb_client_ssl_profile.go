/* Copyright © 2018 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	api "github.com/vmware/go-vmware-nsxt"
	"github.com/vmware/go-vmware-nsxt/loadbalancer"
	"log"
	"net/http"
)

func resourceNsxtLbClientSslProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtLbClientSslProfileCreate,
		Read:   resourceNsxtLbClientSslProfileRead,
		Update: resourceNsxtLbClientSslProfileUpdate,
		Delete: resourceNsxtLbClientSslProfileDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"revision": getRevisionSchema(),
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Description of this resource",
				Optional:    true,
			},
			"display_name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The display name of this resource. Defaults to ID if not set",
				Optional:    true,
				Computed:    true,
			},
			"tag":       getTagsSchema(),
			"ciphers":   getSSLCiphersSchema(),
			"is_secure": getIsSecureSchema(),
			"prefer_server_ciphers": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "During SSL handshake as part of the SSL client Hello client sends an ordered list of ciphers that it can support (or prefers) and typically server selects the first one from the top of that list it can also support. For Perfect Forward Secrecy(PFS), server could override the client's preference",
				Optional:    true,
				Default:     false,
			},
			"protocols": getSSLProtocolsSchema(),
			"session_cache_enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "SSL session caching allows SSL client and server to reuse previously negotiated security parameters avoiding the expensive public key operation during handshake",
				Optional:    true,
				Default:     true,
			},
			"session_cache_timeout": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "Session cache timeout specifies how long the SSL session parameters are held on to and can be reused",
				Optional:    true,
				Default:     300,
			},
		},
	}
}

func resourceNsxtLbClientSslProfileCreate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	ciphers := getStringListFromSchemaSet(d, "ciphers")
	isSecure := d.Get("is_secure").(bool)
	preferServerCiphers := d.Get("prefer_server_ciphers").(bool)
	protocols := getStringListFromSchemaSet(d, "protocols")
	sessionCacheEnabled := d.Get("session_cache_enabled").(bool)
	sessionCacheTimeout := int64(d.Get("session_cache_timeout").(int))
	lbClientSslProfile := loadbalancer.LbClientSslProfile{
		Description:         description,
		DisplayName:         displayName,
		Tags:                tags,
		Ciphers:             ciphers,
		IsSecure:            isSecure,
		PreferServerCiphers: preferServerCiphers,
		Protocols:           protocols,
		SessionCacheEnabled: sessionCacheEnabled,
		SessionCacheTimeout: sessionCacheTimeout,
	}

	lbClientSslProfile, resp, err := nsxClient.ServicesApi.CreateLoadBalancerClientSslProfile(nsxClient.Context, lbClientSslProfile)

	if err != nil {
		return fmt.Errorf("Error during LbClientSslProfile create: %v", err)
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected status returned during LbClientSslProfile create: %v", resp.StatusCode)
	}
	d.SetId(lbClientSslProfile.Id)

	return resourceNsxtLbClientSslProfileRead(d, m)
}

func resourceNsxtLbClientSslProfileRead(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	lbClientSslProfile, resp, err := nsxClient.ServicesApi.ReadLoadBalancerClientSslProfile(nsxClient.Context, id)
	if err != nil {
		return fmt.Errorf("Error during LbClientSslProfile read: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[DEBUG] LbClientSslProfile %s not found", id)
		d.SetId("")
		return nil
	}
	d.Set("revision", lbClientSslProfile.Revision)
	d.Set("description", lbClientSslProfile.Description)
	d.Set("display_name", lbClientSslProfile.DisplayName)
	setTagsInSchema(d, lbClientSslProfile.Tags)
	d.Set("ciphers", lbClientSslProfile.Ciphers)
	d.Set("is_secure", lbClientSslProfile.IsSecure)
	d.Set("prefer_server_ciphers", lbClientSslProfile.PreferServerCiphers)
	d.Set("protocols", lbClientSslProfile.Protocols)
	d.Set("session_cache_enabled", lbClientSslProfile.SessionCacheEnabled)
	d.Set("session_cache_timeout", lbClientSslProfile.SessionCacheTimeout)

	return nil
}

func resourceNsxtLbClientSslProfileUpdate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	revision := int32(d.Get("revision").(int))
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	ciphers := getStringListFromSchemaSet(d, "ciphers")
	isSecure := d.Get("is_secure").(bool)
	preferServerCiphers := d.Get("prefer_server_ciphers").(bool)
	protocols := getStringListFromSchemaSet(d, "protocols")
	sessionCacheEnabled := d.Get("session_cache_enabled").(bool)
	sessionCacheTimeout := int64(d.Get("session_cache_timeout").(int))
	lbClientSslProfile := loadbalancer.LbClientSslProfile{
		Revision:            revision,
		Description:         description,
		DisplayName:         displayName,
		Tags:                tags,
		Ciphers:             ciphers,
		IsSecure:            isSecure,
		PreferServerCiphers: preferServerCiphers,
		Protocols:           protocols,
		SessionCacheEnabled: sessionCacheEnabled,
		SessionCacheTimeout: sessionCacheTimeout,
	}

	lbClientSslProfile, resp, err := nsxClient.ServicesApi.UpdateLoadBalancerClientSslProfile(nsxClient.Context, id, lbClientSslProfile)

	if err != nil || resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Error during LbClientSslProfile update: %v", err)
	}

	return resourceNsxtLbClientSslProfileRead(d, m)
}

func resourceNsxtLbClientSslProfileDelete(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	resp, err := nsxClient.ServicesApi.DeleteLoadBalancerClientSslProfile(nsxClient.Context, id)
	if err != nil {
		return fmt.Errorf("Error during LbClientSslProfile delete: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[DEBUG] LbClientSslProfile %s not found", id)
		d.SetId("")
	}
	return nil
}
