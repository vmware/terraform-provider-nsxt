/* Copyright © 2020 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

var IPSecVpnIkeProfileEncryptionAlgorithms = []string{
	model.IPSecVpnIkeProfile_ENCRYPTION_ALGORITHMS_128,
	model.IPSecVpnIkeProfile_ENCRYPTION_ALGORITHMS_256,
	model.IPSecVpnIkeProfile_ENCRYPTION_ALGORITHMS_GCM_128,
	model.IPSecVpnIkeProfile_ENCRYPTION_ALGORITHMS_GCM_192,
	model.IPSecVpnIkeProfile_ENCRYPTION_ALGORITHMS_256,
}

var IPSecVpnIkeProfileDigestAlgorithms = []string{
	model.IPSecVpnIkeProfile_DIGEST_ALGORITHMS_SHA1,
	model.IPSecVpnIkeProfile_DIGEST_ALGORITHMS_SHA2_256,
	model.IPSecVpnIkeProfile_DIGEST_ALGORITHMS_SHA2_384,
	model.IPSecVpnIkeProfile_DIGEST_ALGORITHMS_SHA2_512,
}

var IPSecVpnIkeProfileDhGroups = []string{
	model.IPSecVpnIkeProfile_DH_GROUPS_GROUP14,
	model.IPSecVpnIkeProfile_DH_GROUPS_GROUP2,
	model.IPSecVpnIkeProfile_DH_GROUPS_GROUP15,
	model.IPSecVpnIkeProfile_DH_GROUPS_GROUP16,
	model.IPSecVpnIkeProfile_DH_GROUPS_GROUP19,
	model.IPSecVpnIkeProfile_DH_GROUPS_GROUP20,
	model.IPSecVpnIkeProfile_DH_GROUPS_GROUP21,
}

var IPSecVpnIkeProfile = []string{
	model.IPSecVpnIkeProfile_IKE_VERSION_V1,
	model.IPSecVpnIkeProfile_IKE_VERSION_V2,
	model.IPSecVpnIkeProfile_IKE_VERSION_FLEX,
}

func resourceNsxtPolicyIpsecVpnIkeProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyIpsecVpnIkeProfileCreate,
		Read:   resourceNsxtPolicyIpsecVpnIkeProfileRead,
		Update: resourceNsxtPolicyIpsecVpnIkeProfileUpdate,
		Delete: resourceNsxtPolicyIpsecVpnIkeProfileDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"nsx_id":       getNsxIDSchema(),
			"path":         getPathSchema(),
			"display_name": getDisplayNameSchema(),
			"description":  getDescriptionSchema(),
			"revision":     getRevisionSchema(),
			"tag":          getTagsSchema(),
			"ike_version": {
				Type:         schema.TypeString,
				Description:  "IKE protocol version to be used. IKE-Flex will initiate IKE-V2 and responds to both IKE-V1 and IKE-V2.",
				ValidateFunc: validation.StringInSlice(IPSecVpnIkeProfile, false),
				Optional:     true,
				Default:      "IKE_V2",
			},
			"encryption_algorithms": {
				Type:        schema.TypeSet,
				Description: "Encryption algorithm is used during Internet Key Exchange(IKE) negotiation. Default is AES_128.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(IPSecVpnIkeProfileEncryptionAlgorithms, false),
				},
				Optional: true,
			},
			"digest_algorithms": {
				Type:        schema.TypeSet,
				Description: "Algorithm to be used for message digest during Internet Key Exchange(IKE) negotiation. Default is SHA2_256.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(IPSecVpnIkeProfileDigestAlgorithms, false),
				},
				Optional: true,
			},
			"dh_groups": {
				Type:        schema.TypeSet,
				Description: "Diffie-Hellman group to be used if PFS is enabled. Default is GROUP14.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(IPSecVpnIkeProfileDhGroups, false),
				},
				Optional: true,
			},
		},
	}
}

func resourceNsxtPolicyIpsecVpnIkeProfileExists(id string, connector *client.RestConnector, isGlobalManager bool) (bool, error) {
	var err error

	client := infra.NewDefaultIpsecVpnIkeProfilesClient(connector)
	_, err = client.Get(id)

	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving resource", err)
}

func resourceNsxtPolicyIpsecVpnIkeProfileCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	// Initialize resource Id and verify this ID is not yet used
	id, err := getOrGenerateID(d, m, resourceNsxtPolicyIpsecVpnIkeProfileExists)
	if err != nil {
		return err
	}

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	ikeVersion := d.Get("ike_version").(string)
	DhGroups := getStringListFromSchemaSet(d, "dh_groups")
	DigestAlgorithms := getStringListFromSchemaSet(d, "digest_algorithms")
	EncryptionAlgorithms := getStringListFromSchemaSet(d, "encryption_algorithms")

	obj := model.IPSecVpnIkeProfile{
		DisplayName:          &displayName,
		Description:          &description,
		IkeVersion:           &ikeVersion,
		DhGroups:             DhGroups,
		DigestAlgorithms:     DigestAlgorithms,
		EncryptionAlgorithms: EncryptionAlgorithms,
	}

	// Create the resource using PATCH
	log.Printf("[INFO] Creating IpsecVpnIkeProfile with ID %s", id)

	client := infra.NewDefaultIpsecVpnIkeProfilesClient(connector)
	err = client.Patch(id, obj)

	if err != nil {
		return handleCreateError("IpsecVpnIkeProfile", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyIpsecVpnIkeProfileRead(d, m)
}

func resourceNsxtPolicyIpsecVpnIkeProfileRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining IpsecVpnIkeProfile ID")
	}

	var obj model.IPSecVpnIkeProfile
	client := infra.NewDefaultIpsecVpnIkeProfilesClient(connector)
	var err error
	obj, err = client.Get(id)
	if err != nil {
		return handleReadError(d, "IpsecVpnIkeProfile", id, err)
	}
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("path", obj.Path)
	d.Set("revision", obj.Revision)
	return nil
}

func resourceNsxtPolicyIpsecVpnIkeProfileUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining IpsecVpnIkeProfile ID")
	}
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	ikeVersion := d.Get("ike_version").(string)
	DhGroups := getStringListFromSchemaSet(d, "dh_groups")
	DigestAlgorithms := getStringListFromSchemaSet(d, "digest_algorithms")
	EncryptionAlgorithms := getStringListFromSchemaSet(d, "encryption_algorithms")
	obj := model.IPSecVpnIkeProfile{
		DisplayName:          &displayName,
		Description:          &description,
		IkeVersion:           &ikeVersion,
		DhGroups:             DhGroups,
		DigestAlgorithms:     DigestAlgorithms,
		EncryptionAlgorithms: EncryptionAlgorithms,
	}
	var err error
	client := infra.NewDefaultIpsecVpnIkeProfilesClient(connector)
	err = client.Patch(id, obj)

	if err != nil {
		return handleUpdateError("IpsecVpnIkeProfile", id, err)
	}

	return resourceNsxtPolicyIpsecVpnIkeProfileRead(d, m)

}

func resourceNsxtPolicyIpsecVpnIkeProfileDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining IpsecVpnIkeProfile ID")
	}

	connector := getPolicyConnector(m)
	var err error
	client := infra.NewDefaultIpsecVpnIkeProfilesClient(connector)
	err = client.Delete(id)

	if err != nil {
		return handleDeleteError("IpsecVpnIkeProfile", id, err)
	}

	return nil
}
