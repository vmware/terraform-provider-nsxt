/* Copyright © 2018 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	api "github.com/vmware/go-vmware-nsxt"
	"github.com/vmware/go-vmware-nsxt/loadbalancer"
	"log"
	"net/http"
)

var cookieModeTypes = []string{"INSERT", "PREFIX", "REWRITE"}

func resourceNsxtLbCookiePersistenceProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtLbCookiePersistenceProfileCreate,
		Read:   resourceNsxtLbCookiePersistenceProfileRead,
		Update: resourceNsxtLbCookiePersistenceProfileUpdate,
		Delete: resourceNsxtLbCookiePersistenceProfileDelete,
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
			"tag": getTagsSchema(),
			"cookie_mode": &schema.Schema{
				Type:         schema.TypeString,
				Description:  "The cookie persistence mode",
				Optional:     true,
				Default:      "INSERT",
				ValidateFunc: validation.StringInSlice(cookieModeTypes, false),
			},
			"cookie_name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The name of the cookie",
				Required:    true,
			},
			"persistence_shared": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "????If persistence shared flag is not set in the cookie persistence profile bound to a virtual server, it defaults to cookie persistence that is private to each virtual server and is qualified by the pool. This is accomplished by load balancer inserting a cookie with name in the format &lt;name&gt;.&lt;virtual_server_id&gt;.&lt;pool_id&gt;. If persistence shared flag is set in the cookie persistence profile, in cookie insert mode, cookie persistence could be shared across multiple virtual servers that are bound to the same pools. The cookie name would be changed to &lt;name&gt;.&lt;profile-id&gt;.&lt;pool-id&gt;. If persistence shared flag is not set in the sourceIp persistence profile bound to a virtual server, each virtual server that the profile is bound to maintains its own private persistence table. If persistence shared flag is set in the sourceIp persistence profile, all virtual servers the profile is bound to share the same persistence table",
				Optional:    true,
				Default:     false,
			},
			"cookie_fallback": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "A boolean flag which reflects whether once the server points by this cookie is down, a new server is selected, or the requests will be rejected",
				Optional:    true,
				Default:     true,
			},
			"cookie_garble": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "A boolean flag which reflects whether the cookie value (server IP and port) would be encrypted or in plain text",
				Optional:    true,
				Default:     true,
			},
			"cookie_domain": &schema.Schema{
				Type:        schema.TypeString,
				Description: "HTTP cookie domain (for INSERT mode only)",
				Optional:    true,
				// TODO(asarfaty) add validation on cookie domain values
			},
			"cookie_path": &schema.Schema{
				Type:        schema.TypeString,
				Description: "HTTP cookie path (for INSERT mode only)",
				Optional:    true,
			},
			// TODO(asarfaty) add cookie_time too, after SDK support
		},
	}
}

func resourceNsxtLbCookiePersistenceProfileCreate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	persistenceShared := d.Get("persistence_shared").(bool)
	cookieDomain := d.Get("cookie_domain").(string)
	cookieFallback := d.Get("cookie_fallback").(bool)
	cookieGarble := d.Get("cookie_garble").(bool)
	cookieMode := d.Get("cookie_mode").(string)
	cookieName := d.Get("cookie_name").(string)
	cookiePath := d.Get("cookie_path").(string)
	//cookieTime := d.Get("cookie_time").(*LbCookieTime)
	lbCookiePersistenceProfile := loadbalancer.LbCookiePersistenceProfile{
		Description:       description,
		DisplayName:       displayName,
		Tags:              tags,
		PersistenceShared: persistenceShared,
		CookieDomain:      cookieDomain,
		CookieFallback:    cookieFallback,
		CookieGarble:      cookieGarble,
		CookieMode:        cookieMode,
		CookieName:        cookieName,
		CookiePath:        cookiePath,
		//CookieTime: cookieTime,
	}

	lbCookiePersistenceProfile, resp, err := nsxClient.ServicesApi.CreateLoadBalancerCookiePersistenceProfile(nsxClient.Context, lbCookiePersistenceProfile)

	if err != nil {
		return fmt.Errorf("Error during LbCookiePersistenceProfile create: %v", err)
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected status returned during LbCookiePersistenceProfile create: %v", resp.StatusCode)
	}
	d.SetId(lbCookiePersistenceProfile.Id)

	return resourceNsxtLbCookiePersistenceProfileRead(d, m)
}

func resourceNsxtLbCookiePersistenceProfileRead(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	lbCookiePersistenceProfile, resp, err := nsxClient.ServicesApi.ReadLoadBalancerCookiePersistenceProfile(nsxClient.Context, id)
	if err != nil {
		return fmt.Errorf("Error during LbCookiePersistenceProfile read: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[DEBUG] LbCookiePersistenceProfile %s not found", id)
		d.SetId("")
		return nil
	}
	d.Set("revision", lbCookiePersistenceProfile.Revision)
	d.Set("description", lbCookiePersistenceProfile.Description)
	d.Set("display_name", lbCookiePersistenceProfile.DisplayName)
	setTagsInSchema(d, lbCookiePersistenceProfile.Tags)
	d.Set("persistence_shared", lbCookiePersistenceProfile.PersistenceShared)
	d.Set("cookie_domain", lbCookiePersistenceProfile.CookieDomain)
	d.Set("cookie_fallback", lbCookiePersistenceProfile.CookieFallback)
	d.Set("cookie_garble", lbCookiePersistenceProfile.CookieGarble)
	d.Set("cookie_mode", lbCookiePersistenceProfile.CookieMode)
	d.Set("cookie_name", lbCookiePersistenceProfile.CookieName)
	d.Set("cookie_path", lbCookiePersistenceProfile.CookiePath)
	//d.Set("cookie_time", lbCookiePersistenceProfile.CookieTime)

	return nil
}

func resourceNsxtLbCookiePersistenceProfileUpdate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	revision := int32(d.Get("revision").(int))
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	persistenceShared := d.Get("persistence_shared").(bool)
	cookieDomain := d.Get("cookie_domain").(string)
	cookieFallback := d.Get("cookie_fallback").(bool)
	cookieGarble := d.Get("cookie_garble").(bool)
	cookieMode := d.Get("cookie_mode").(string)
	cookieName := d.Get("cookie_name").(string)
	cookiePath := d.Get("cookie_path").(string)
	//cookieTime := d.Get("cookie_time").(*LbCookieTime)
	lbCookiePersistenceProfile := loadbalancer.LbCookiePersistenceProfile{
		Revision:          revision,
		Description:       description,
		DisplayName:       displayName,
		Tags:              tags,
		PersistenceShared: persistenceShared,
		CookieDomain:      cookieDomain,
		CookieFallback:    cookieFallback,
		CookieGarble:      cookieGarble,
		CookieMode:        cookieMode,
		CookieName:        cookieName,
		CookiePath:        cookiePath,
		//CookieTime: cookieTime,
	}

	lbCookiePersistenceProfile, resp, err := nsxClient.ServicesApi.UpdateLoadBalancerCookiePersistenceProfile(nsxClient.Context, id, lbCookiePersistenceProfile)

	if err != nil || resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Error during LbCookiePersistenceProfile update: %v", err)
	}

	return resourceNsxtLbCookiePersistenceProfileRead(d, m)
}

func resourceNsxtLbCookiePersistenceProfileDelete(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	resp, err := nsxClient.ServicesApi.DeleteLoadBalancerPersistenceProfile(nsxClient.Context, id)
	if err != nil {
		return fmt.Errorf("Error during LbCookiePersistenceProfile delete: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[DEBUG] LbCookiePersistenceProfile %s not found", id)
		d.SetId("")
	}
	return nil
}
