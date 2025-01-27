/* Copyright © 2024 Broadcom, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	clientLayer "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/metadata"
)

var l7AccessProfileDefaultActionValues = []string{
	model.L7AccessProfile_DEFAULT_ACTION_ALLOW,
	model.L7AccessProfile_DEFAULT_ACTION_REJECT,
	model.L7AccessProfile_DEFAULT_ACTION_REJECT_WITH_RESPONSE,
}

var l7AccessProfileActionValues = []string{
	model.L7AccessEntry_ACTION_ALLOW,
	model.L7AccessEntry_ACTION_REJECT,
	model.L7AccessEntry_ACTION_REJECT_WITH_RESPONSE,
}

var l7AccessProfileAttributeSourceValues = []string{
	model.L7AccessAttributes_ATTRIBUTE_SOURCE_SYSTEM,
	model.L7AccessAttributes_ATTRIBUTE_SOURCE_CUSTOM,
}

var l7AccessProfileSubAttributeKeyValues = []string{
	model.PolicySubAttributes_KEY_TLS_CIPHER_SUITE,
	model.PolicySubAttributes_KEY_TLS_VERSION,
	model.PolicySubAttributes_KEY_CIFS_SMB_VERSION,
}

var l7AccessProfileDatatypeValues = []string{
	model.PolicySubAttributes_DATATYPE_STRING,
}

var l7AccessProfileKeyAttributeValues = []string{
	model.L7AccessAttributes_KEY_APP_ID,
	model.L7AccessAttributes_KEY_DOMAIN_NAME,
	model.L7AccessAttributes_KEY_URL_CATEGORY,
	model.L7AccessAttributes_KEY_URL_REPUTATION,
	model.L7AccessAttributes_KEY_CUSTOM_URL,
}

var l7AccessProfileSchema = map[string]*metadata.ExtendedSchema{
	"nsx_id":       metadata.GetExtendedSchema(getNsxIDSchema()),
	"path":         metadata.GetExtendedSchema(getPathSchema()),
	"display_name": metadata.GetExtendedSchema(getDisplayNameSchema()),
	"description":  metadata.GetExtendedSchema(getDescriptionSchema()),
	"revision":     metadata.GetExtendedSchema(getRevisionSchema()),
	"tag":          metadata.GetExtendedSchema(getTagsSchema()),
	"context":      metadata.GetExtendedSchema(getContextSchema(false, false, false)),
	"default_action_logged": {
		Schema: schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "bool",
			SdkFieldName: "DefaultActionLogged",
		},
	},
	"default_action": {
		Schema: schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validation.StringInSlice(l7AccessProfileDefaultActionValues, false),
			Required:     true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "string",
			SdkFieldName: "DefaultAction",
			OmitIfEmpty:  true,
		},
	},
	"entry_count": {
		Schema: schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "int",
			SdkFieldName: "EntryCount",
		},
	},
	"l7_access_entry": {
		Schema: schema.Schema{
			Type: schema.TypeList,
			Elem: &metadata.ExtendedResource{
				Schema: map[string]*metadata.ExtendedSchema{
					"id": {
						Schema: schema.Schema{
							Type:         schema.TypeString,
							Description:  "NSX ID for this resource",
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(1, 1024),
						},
						Metadata: metadata.Metadata{
							SchemaType:   "string",
							SdkFieldName: "Id",
						},
					},
					"display_name": {
						Schema: schema.Schema{
							Type:         schema.TypeString,
							Description:  "Display name for this resource",
							Computed:     true,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(1, 255),
						},
						Metadata: metadata.Metadata{
							SchemaType:   "string",
							SdkFieldName: "DisplayName",
						},
					},
					"action": {
						Schema: schema.Schema{
							Type:         schema.TypeString,
							ValidateFunc: validation.StringInSlice(l7AccessProfileActionValues, false),
							Optional:     true,
						},
						Metadata: metadata.Metadata{
							SchemaType:   "string",
							SdkFieldName: "Action",
						},
					},
					"attribute": {
						Schema: schema.Schema{
							Type: schema.TypeList,
							Elem: &metadata.ExtendedResource{
								Schema: map[string]*metadata.ExtendedSchema{
									"sub_attribute": {
										Schema: schema.Schema{
											Type: schema.TypeList,
											Elem: &metadata.ExtendedResource{
												Schema: map[string]*metadata.ExtendedSchema{
													"datatype": {
														Schema: schema.Schema{
															Type:         schema.TypeString,
															ValidateFunc: validation.StringInSlice(l7AccessProfileDatatypeValues, false),
															Optional:     true,
														},
														Metadata: metadata.Metadata{
															SchemaType:   "string",
															SdkFieldName: "Datatype",
														},
													},
													"value": {
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
														},
														Metadata: metadata.Metadata{
															SchemaType:   "list",
															SdkFieldName: "Value",
														},
													},
													"key": {
														Schema: schema.Schema{
															Type:         schema.TypeString,
															ValidateFunc: validation.StringInSlice(l7AccessProfileSubAttributeKeyValues, false),
															Optional:     true,
														},
														Metadata: metadata.Metadata{
															SchemaType:   "string",
															SdkFieldName: "Key",
														},
													},
												},
											},
											Optional: true,
										},
										Metadata: metadata.Metadata{
											SchemaType:   "list",
											SdkFieldName: "SubAttributes",
											ReflectType:  reflect.TypeOf(model.PolicySubAttributes{}),
										},
									},
									"attribute_source": {
										Schema: schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.StringInSlice(l7AccessProfileAttributeSourceValues, false),
											Optional:     true,
											Default:      model.L7AccessAttributes_ATTRIBUTE_SOURCE_SYSTEM,
										},
										Metadata: metadata.Metadata{
											SchemaType:   "string",
											SdkFieldName: "AttributeSource",
										},
									},
									"custom_url_partial_match": {
										Schema: schema.Schema{
											Type:     schema.TypeBool,
											Optional: true,
										},
										Metadata: metadata.Metadata{
											SchemaType:   "bool",
											SdkFieldName: "CustomUrlPartialMatch",
										},
									},
									"description": {
										Schema: schema.Schema{
											Type:     schema.TypeString,
											Optional: true,
										},
										Metadata: metadata.Metadata{
											SchemaType:   "string",
											SdkFieldName: "Description",
										},
									},
									"key": {
										Schema: schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.StringInSlice(l7AccessProfileKeyAttributeValues, false),
											Optional:     true,
										},
										Metadata: metadata.Metadata{
											SchemaType:   "string",
											SdkFieldName: "Key",
										},
									},
									"datatype": {
										Schema: schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.StringInSlice(l7AccessProfileDatatypeValues, false),
											Optional:     true,
										},
										Metadata: metadata.Metadata{
											SchemaType:   "string",
											SdkFieldName: "Datatype",
										},
									},
									"is_alg_type": {
										Schema: schema.Schema{
											Type:     schema.TypeBool,
											Optional: true,
										},
										Metadata: metadata.Metadata{
											SchemaType:   "bool",
											SdkFieldName: "IsALGType",
										},
									},
									"value": {
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
										},
										Metadata: metadata.Metadata{
											SchemaType:   "list",
											SdkFieldName: "Value",
										},
									},
									"metadata": {
										Schema: schema.Schema{
											Type: schema.TypeList,
											Elem: &metadata.ExtendedResource{
												Schema: map[string]*metadata.ExtendedSchema{
													"value": {
														Schema: schema.Schema{
															Type:     schema.TypeString,
															Optional: true,
														},
														Metadata: metadata.Metadata{
															SchemaType:   "string",
															SdkFieldName: "Value",
														},
													},
													"key": {
														Schema: schema.Schema{
															Type:     schema.TypeString,
															Optional: true,
														},
														Metadata: metadata.Metadata{
															SchemaType:   "string",
															SdkFieldName: "Key",
														},
													},
												},
											},
											Optional: true,
										},
										Metadata: metadata.Metadata{
											SchemaType:   "list",
											SdkFieldName: "Metadata",
											ReflectType:  reflect.TypeOf(model.ContextProfileAttributesMetadata{}),
										},
									},
								},
							},
							Optional: true,
						},
						Metadata: metadata.Metadata{
							SchemaType:   "list",
							SdkFieldName: "Attributes",
							ReflectType:  reflect.TypeOf(model.L7AccessAttributes{}),
						},
					},
					"disabled": {
						Schema: schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},
						Metadata: metadata.Metadata{
							SchemaType:   "bool",
							SdkFieldName: "Disabled",
						},
					},
					"logged": {
						Schema: schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},
						Metadata: metadata.Metadata{
							SchemaType:   "bool",
							SdkFieldName: "Logged",
						},
					},
					"sequence_number": {
						Schema: schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
						Metadata: metadata.Metadata{
							SchemaType:   "int",
							SdkFieldName: "SequenceNumber",
						},
					},
				},
			},
			Optional: true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "list",
			SdkFieldName: "L7AccessEntries",
			ReflectType:  reflect.TypeOf(model.L7AccessEntry{}),
		},
	},
}

func resourceNsxtPolicyL7AccessProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyL7AccessProfileCreate,
		Read:   resourceNsxtPolicyL7AccessProfileRead,
		Update: resourceNsxtPolicyL7AccessProfileUpdate,
		Delete: resourceNsxtPolicyL7AccessProfileDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtParentPathResourceImporter,
		},
		Schema: metadata.GetSchemaFromExtendedSchema(l7AccessProfileSchema),
	}
}

func resourceNsxtPolicyL7AccessProfileExists(sessionContext utl.SessionContext, id string, connector client.Connector) (bool, error) {
	var err error
	client := clientLayer.NewL7AccessProfilesClient(sessionContext, connector)
	_, err = client.Get(id)
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving resource", err)
}

func resourceNsxtPolicyL7AccessProfileCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id, err := getOrGenerateID2(d, m, resourceNsxtPolicyL7AccessProfileExists)
	if err != nil {
		return err
	}

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)

	obj := model.L7AccessProfile{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, l7AccessProfileSchema, "", nil); err != nil {
		return err
	}

	log.Printf("[INFO] Creating L7AccessProfile with ID %s", id)

	client := clientLayer.NewL7AccessProfilesClient(getSessionContext(d, m), connector)
	_, err = client.Patch(id, obj, nil)
	if err != nil {
		return handleCreateError("L7AccessProfile", id, err)
	}
	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyL7AccessProfileRead(d, m)
}

func resourceNsxtPolicyL7AccessProfileRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining L7AccessProfile ID")
	}

	client := clientLayer.NewL7AccessProfilesClient(getSessionContext(d, m), connector)
	obj, err := client.Get(id)
	if err != nil {
		return handleReadError(d, "L7AccessProfile", id, err)
	}

	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("revision", obj.Revision)
	d.Set("path", obj.Path)

	elem := reflect.ValueOf(&obj).Elem()
	return metadata.StructToSchema(elem, d, l7AccessProfileSchema, "", nil)
}

func resourceNsxtPolicyL7AccessProfileUpdate(d *schema.ResourceData, m interface{}) error {

	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining L7AccessProfile ID")
	}

	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getPolicyTagsFromSchema(d)

	revision := int64(d.Get("revision").(int))

	obj := model.L7AccessProfile{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
		Revision:    &revision,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, l7AccessProfileSchema, "", nil); err != nil {
		return err
	}
	client := clientLayer.NewL7AccessProfilesClient(getSessionContext(d, m), connector)
	_, err := client.Patch(id, obj, nil)
	if err != nil {
		return handleUpdateError("L7AccessProfile", id, err)
	}

	return resourceNsxtPolicyL7AccessProfileRead(d, m)
}

func resourceNsxtPolicyL7AccessProfileDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining L7AccessProfile ID")
	}

	connector := getPolicyConnector(m)

	client := clientLayer.NewL7AccessProfilesClient(getSessionContext(d, m), connector)
	err := client.Delete(id, nil)

	if err != nil {
		return handleDeleteError("L7AccessProfile", id, err)
	}

	return nil
}
