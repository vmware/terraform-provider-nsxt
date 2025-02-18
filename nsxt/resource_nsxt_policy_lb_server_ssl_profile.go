/* Copyright © 2024 Broadcom, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	clientLayer "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	"github.com/vmware/terraform-provider-nsxt/nsxt/metadata"
)

var lbServerSslProfileCipherGroupLabelValues = []string{
	model.LBServerSslProfile_CIPHER_GROUP_LABEL_BALANCED,
	model.LBServerSslProfile_CIPHER_GROUP_LABEL_HIGH_SECURITY,
	model.LBServerSslProfile_CIPHER_GROUP_LABEL_HIGH_COMPATIBILITY,
	model.LBServerSslProfile_CIPHER_GROUP_LABEL_CUSTOM,
}

var lbServerSslProfileProtocolsValues = []string{
	model.LBServerSslProfile_PROTOCOLS_SSL_V2,
	model.LBServerSslProfile_PROTOCOLS_SSL_V3,
	model.LBServerSslProfile_PROTOCOLS_TLS_V1,
	model.LBServerSslProfile_PROTOCOLS_TLS_V1_1,
	model.LBServerSslProfile_PROTOCOLS_TLS_V1_2,
}

var lbServerSslProfilePathExample = getMultitenancyPathExample("/infra/lb-server-ssl-profiles/[profile]")

var lbServerSslProfileSchema = map[string]*metadata.ExtendedSchema{
	"nsx_id":       metadata.GetExtendedSchema(getNsxIDSchema()),
	"path":         metadata.GetExtendedSchema(getPathSchema()),
	"display_name": metadata.GetExtendedSchema(getDisplayNameSchema()),
	"description":  metadata.GetExtendedSchema(getDescriptionSchema()),
	"revision":     metadata.GetExtendedSchema(getRevisionSchema()),
	"tag":          metadata.GetExtendedSchema(getTagsSchema()),
	"context":      metadata.GetExtendedSchema(getContextSchema(false, false, false)),
	"cipher_group_label": {
		Schema: schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validation.StringInSlice(lbServerSslProfileCipherGroupLabelValues, false),
			Required:     true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "string",
			SdkFieldName: "CipherGroupLabel",
		},
	},
	"ciphers": {
		Schema: schema.Schema{
			Type: schema.TypeList,
			Elem: &metadata.ExtendedSchema{
				Schema: schema.Schema{
					Type: schema.TypeString,
				},
				Metadata: metadata.Metadata{
					SchemaType: "string",
				},
			},
			Optional: true,
			Computed: true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "list",
			SdkFieldName: "Ciphers",
		},
	},
	"protocols": {
		Schema: schema.Schema{
			Type: schema.TypeList,
			Elem: &metadata.ExtendedSchema{
				Schema: schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(lbServerSslProfileProtocolsValues, false),
				},
				Metadata: metadata.Metadata{
					SchemaType: "string",
				},
			},
			Optional: true,
			Computed: true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "list",
			SdkFieldName: "Protocols",
		},
	},
	"session_cache_enabled": {
		Schema: schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "bool",
			SdkFieldName: "SessionCacheEnabled",
		},
	},
	"is_secure": {
		Schema: schema.Schema{
			Type:     schema.TypeBool,
			Computed: true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "bool",
			SdkFieldName: "IsSecure",
			ReadOnly:     true,
		},
	},
	"is_fips": {
		Schema: schema.Schema{
			Type:     schema.TypeBool,
			Computed: true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "bool",
			SdkFieldName: "IsFips",
			ReadOnly:     true,
		},
	},
}

func resourceNsxtPolicyLBServerSslProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyLBServerSslProfileCreate,
		Read:   resourceNsxtPolicyLBServerSslProfileRead,
		Update: resourceNsxtPolicyLBServerSslProfileUpdate,
		Delete: resourceNsxtPolicyLBServerSslProfileDelete,
		Importer: &schema.ResourceImporter{
			State: getPolicyPathOrIDResourceImporter(lbServerSslProfilePathExample),
		},
		Schema: metadata.GetSchemaFromExtendedSchema(lbServerSslProfileSchema),
	}
}

func resourceNsxtPolicyLBServerSslProfileExists(id string, connector client.Connector, isGlobalManager bool) (bool, error) {
	var err error

	client := clientLayer.NewLbServerSslProfilesClient(connector)
	_, err = client.Get(id)
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving resource", err)
}

func resourceNsxtPolicyLBServerSslProfileCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id, err := getOrGenerateID(d, m, resourceNsxtPolicyLBServerSslProfileExists)
	if err != nil {
		return err
	}

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)

	obj := model.LBServerSslProfile{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, lbServerSslProfileSchema, "", nil); err != nil {
		return err
	}

	log.Printf("[INFO] Creating LBServerSslProfile with ID %s", id)

	client := clientLayer.NewLbServerSslProfilesClient(connector)
	err = client.Patch(id, obj)
	if err != nil {
		return handleCreateError("LBServerSslProfile", id, err)
	}
	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyLBServerSslProfileRead(d, m)
}

func resourceNsxtPolicyLBServerSslProfileRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining LBServerSslProfile ID")
	}

	client := clientLayer.NewLbServerSslProfilesClient(connector)

	obj, err := client.Get(id)
	if err != nil {
		return handleReadError(d, "LBServerSslProfile", id, err)
	}

	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("revision", obj.Revision)
	d.Set("path", obj.Path)

	elem := reflect.ValueOf(&obj).Elem()
	return metadata.StructToSchema(elem, d, lbServerSslProfileSchema, "", nil)
}

func resourceNsxtPolicyLBServerSslProfileUpdate(d *schema.ResourceData, m interface{}) error {

	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining LBServerSslProfile ID")
	}

	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getPolicyTagsFromSchema(d)

	revision := int64(d.Get("revision").(int))

	obj := model.LBServerSslProfile{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
		Revision:    &revision,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, lbServerSslProfileSchema, "", nil); err != nil {
		return err
	}
	client := clientLayer.NewLbServerSslProfilesClient(connector)
	_, err := client.Update(id, obj)
	if err != nil {
		return handleUpdateError("LBServerSslProfile", id, err)
	}

	return resourceNsxtPolicyLBServerSslProfileRead(d, m)
}

func resourceNsxtPolicyLBServerSslProfileDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining LBServerSslProfile ID")
	}

	connector := getPolicyConnector(m)

	client := clientLayer.NewLbServerSslProfilesClient(connector)
	err := client.Delete(id, nil)

	if err != nil {
		return handleDeleteError("LBServerSslProfile", id, err)
	}

	return nil
}
