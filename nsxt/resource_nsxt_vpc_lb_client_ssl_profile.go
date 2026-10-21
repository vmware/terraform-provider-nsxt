// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var vpcLbClientSslProfilePathExample = "/orgs/[org]/projects/[project]/vpcs/[vpc]/vpc-lb-client-ssl-profiles/[profile]"

func resourceNsxtVpcLBClientSslProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtVpcLBClientSslProfileCreate,
		Read:   resourceNsxtVpcLBClientSslProfileRead,
		Update: resourceNsxtVpcLBClientSslProfileUpdate,
		Delete: resourceNsxtVpcLBClientSslProfileDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtVersionCheckImporter("9.0.0", "VPC LB Client SSL Profile", getVpcPathResourceImporter(vpcLbClientSslProfilePathExample)),
		},

		Schema: map[string]*schema.Schema{
			"nsx_id":       getNsxIDSchema(),
			"path":         getPathSchema(),
			"display_name": getDisplayNameSchema(),
			"description":  getDescriptionSchema(),
			"revision":     getRevisionSchema(),
			"tag":          getTagsSchema(),
			"context":      getContextSchemaExtended(true, false, true, true),
			"cipher_group_label": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(lBClientSslProfileCipherGroupLabelValues, false),
				Optional:     true,
				Default:      model.LBClientSslProfile_CIPHER_GROUP_LABEL_BALANCED,
				Description:  "A label of cipher group which is mostly consumed by GUI. Default value is BALANCED.",
			},
			"ciphers": getSSLCiphersSchema(),
			"is_fips": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "This flag is set to true when all the ciphers and protocols are FIPS compliant. It is set to false when one of the ciphers or protocols are not FIPS compliant.",
			},
			"is_secure": getIsSecureSchema(),
			"prefer_server_ciphers": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "During SSL handshake as part of the SSL client Hello client sends an ordered list of ciphers that it can support (or prefers) and typically server selects the first one from the top of that list it can also support. For Perfect Forward Secrecy(PFS), server could override the client's preference.",
			},
			"protocols": getSSLProtocolsSchema(),
			"session_cache_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "If set to true, SSL session caching allows SSL client and server to reuse previously negotiated security parameters avoiding the expensive public key operation during handshake.",
			},
			"session_cache_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     300,
				Description: "Session cache timeout specifies how long the SSL session parameters are held on to and can be reused. Default value is 300.",
			},
		},
	}
}

func resourceNsxtVpcLBClientSslProfileExists(sessionContext utl.SessionContext, id string, connector client.Connector) (bool, error) {
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbClientSslProfilesClient(sessionContext, connector)
	_, err := c.Get(parents[0], parents[1], parents[2], id)
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving resource", err)
}

func resourceNsxtVpcLBClientSslProfilePatch(d *schema.ResourceData, m interface{}, id string) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbClientSslProfilesClient(sessionContext, connector)

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsWithProviderManagedDefaults(d, m)
	cipherGroupLabel := d.Get("cipher_group_label").(string)
	ciphers := getStringListFromSchemaSet(d, "ciphers")
	preferServerCiphers := d.Get("prefer_server_ciphers").(bool)
	protocols := getStringListFromSchemaSet(d, "protocols")
	sessionCacheEnabled := d.Get("session_cache_enabled").(bool)
	sessionCacheTimeout := int64(d.Get("session_cache_timeout").(int))

	obj := model.LBClientSslProfile{
		DisplayName:         &displayName,
		Description:         &description,
		Tags:                tags,
		CipherGroupLabel:    &cipherGroupLabel,
		Ciphers:             ciphers,
		PreferServerCiphers: &preferServerCiphers,
		Protocols:           protocols,
		SessionCacheEnabled: &sessionCacheEnabled,
		SessionCacheTimeout: &sessionCacheTimeout,
	}

	log.Printf("[INFO] Patching VPC LBClientSslProfile with ID %s", id)
	return c.Patch(parents[0], parents[1], parents[2], id, obj)
}

func resourceNsxtVpcLBClientSslProfileCreate(d *schema.ResourceData, m interface{}) error {
	if !util.NsxVersionHigherOrEqual("9.0.0") {
		return fmt.Errorf("VPC LB Client SSL Profile resource requires NSX version 9.0.0 or higher")
	}

	// Initialize resource Id and verify this ID is not yet used
	id, err := getOrGenerateID2(d, m, resourceNsxtVpcLBClientSslProfileExists)
	if err != nil {
		return err
	}

	err = resourceNsxtVpcLBClientSslProfilePatch(d, m, id)
	if err != nil {
		return handleCreateError("VPC LBClientSslProfile", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)
	return resourceNsxtVpcLBClientSslProfileRead(d, m)
}

func resourceNsxtVpcLBClientSslProfileRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbClientSslProfilesClient(sessionContext, connector)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VPC LBClientSslProfile ID")
	}

	obj, err := c.Get(parents[0], parents[1], parents[2], id)
	if err != nil {
		return handleReadError(d, "VPC LBClientSslProfile", id, err)
	}

	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("path", obj.Path)
	d.Set("revision", obj.Revision)

	d.Set("cipher_group_label", obj.CipherGroupLabel)
	d.Set("ciphers", obj.Ciphers)
	d.Set("is_fips", obj.IsFips)
	d.Set("is_secure", obj.IsSecure)
	d.Set("prefer_server_ciphers", obj.PreferServerCiphers)
	d.Set("protocols", obj.Protocols)
	d.Set("session_cache_enabled", obj.SessionCacheEnabled)
	d.Set("session_cache_timeout", obj.SessionCacheTimeout)

	return nil
}

func resourceNsxtVpcLBClientSslProfileUpdate(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VPC LBClientSslProfile ID")
	}

	err := resourceNsxtVpcLBClientSslProfilePatch(d, m, id)
	if err != nil {
		return handleUpdateError("VPC LBClientSslProfile", id, err)
	}

	return resourceNsxtVpcLBClientSslProfileRead(d, m)
}

func resourceNsxtVpcLBClientSslProfileDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VPC LBClientSslProfile ID")
	}

	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbClientSslProfilesClient(sessionContext, connector)

	force := true
	err := c.Delete(parents[0], parents[1], parents[2], id, &force)
	if err != nil {
		return handleDeleteError("VPC LBClientSslProfile", id, err)
	}

	return nil
}
