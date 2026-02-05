// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/metadata"
)

var lbCookiePersistenceProfileCookieModeValues = []string{
	model.LBCookiePersistenceProfile_COOKIE_MODE_INSERT,
	model.LBCookiePersistenceProfile_COOKIE_MODE_PREFIX,
	model.LBCookiePersistenceProfile_COOKIE_MODE_REWRITE,
}

var lbPersistenceProfilePathExample = getMultitenancyPathExample("/infra/lb-persistence-profiles/[profile]")

var lbCookiePersistenceProfileSchema = map[string]*metadata.ExtendedSchema{
	"nsx_id":       metadata.GetExtendedSchema(getNsxIDSchema()),
	"path":         metadata.GetExtendedSchema(getPathSchema()),
	"display_name": metadata.GetExtendedSchema(getDisplayNameSchema()),
	"description":  metadata.GetExtendedSchema(getDescriptionSchema()),
	"revision":     metadata.GetExtendedSchema(getRevisionSchema()),
	"tag":          metadata.GetExtendedSchema(getTagsSchema()),
	"persistence_shared": {
		Schema: schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "bool",
			SdkFieldName: "PersistenceShared",
		},
	},
	"cookie_name": {
		Schema: schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "NSXLB",
		},
		Metadata: metadata.Metadata{
			SchemaType:   "string",
			SdkFieldName: "CookieName",
		},
	},
	"cookie_mode": {
		Schema: schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validation.StringInSlice(lbCookiePersistenceProfileCookieModeValues, false),
			Optional:     true,
			Default:      model.LBCookiePersistenceProfile_COOKIE_MODE_INSERT,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "string",
			SdkFieldName: "CookieMode",
		},
	},
	"cookie_fallback": {
		Schema: schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "bool",
			SdkFieldName: "CookieFallback",
		},
	},
	"cookie_garble": {
		Schema: schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "bool",
			SdkFieldName: "CookieGarble",
		},
	},
	"cookie_domain": {
		Schema: schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "string",
			SdkFieldName: "CookieDomain",
			OmitIfEmpty:  true,
		},
	},
	"cookie_path": {
		Schema: schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "string",
			SdkFieldName: "CookiePath",
		},
	},
	"session_cookie_time": {
		Schema: schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			ConflictsWith: []string{"persistence_cookie_time"},
			Elem: &metadata.ExtendedResource{
				Schema: map[string]*metadata.ExtendedSchema{
					"max_idle": {
						Schema: schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
						Metadata: metadata.Metadata{
							SchemaType:   "int",
							SdkFieldName: "CookieMaxIdle",
						},
					},
					"max_life": {
						Schema: schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
						Metadata: metadata.Metadata{
							SchemaType:   "int",
							SdkFieldName: "CookieMaxLife",
						},
					},
				},
			},
			Optional: true,
		},
		Metadata: metadata.Metadata{
			SchemaType:      "struct",
			SdkFieldName:    "CookieTime",
			ReflectType:     reflect.TypeOf(model.LBSessionCookieTime{}),
			PolymorphicType: metadata.PolymorphicTypeFlatten,
			TypeIdentifier:  metadata.StandardTypeIdentifier,
			BindingType:     model.LBSessionCookieTimeBindingType(),
			ResourceType:    model.LBCookieTime_TYPE_LBSESSIONCOOKIETIME,
		},
	},
	"persistence_cookie_time": {
		Schema: schema.Schema{
			Type:     schema.TypeList,
			MaxItems: 1,
			Elem: &metadata.ExtendedResource{
				Schema: map[string]*metadata.ExtendedSchema{
					"max_idle": {
						Schema: schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
						Metadata: metadata.Metadata{
							SchemaType:   "int",
							SdkFieldName: "CookieMaxIdle",
						},
					},
				},
			},
			Optional: true,
		},
		Metadata: metadata.Metadata{
			SchemaType:      "struct",
			SdkFieldName:    "CookieTime",
			ReflectType:     reflect.TypeOf(model.LBSessionCookieTime{}),
			PolymorphicType: metadata.PolymorphicTypeFlatten,
			TypeIdentifier:  metadata.StandardTypeIdentifier,
			BindingType:     model.LBPersistenceCookieTimeBindingType(),
			ResourceType:    model.LBCookieTime_TYPE_LBPERSISTENCECOOKIETIME,
		},
	},
	"cookie_httponly": {
		Schema: schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "bool",
			SdkFieldName: "CookieHttponly",
		},
	},
	"cookie_secure": {
		Schema: schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "bool",
			SdkFieldName: "CookieSecure",
		},
	},
}

func resourceNsxtPolicyLBCookiePersistenceProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyLBCookiePersistenceProfileCreate,
		Read:   resourceNsxtPolicyLBCookiePersistenceProfileRead,
		Update: resourceNsxtPolicyLBCookiePersistenceProfileUpdate,
		Delete: resourceNsxtPolicyLBPersistenceProfileDelete,
		Importer: &schema.ResourceImporter{
			State: getPolicyPathOrIDResourceImporter(lbPersistenceProfilePathExample),
		},
		Schema: metadata.GetSchemaFromExtendedSchema(lbCookiePersistenceProfileSchema),
	}
}

func resourceNsxtPolicyLBPersistenceProfileExists(id string, connector client.Connector, isGlobalManager bool) (bool, error) {
	var err error
	sessionContext := utl.SessionContext{ClientType: utl.Local}
	client := cliLbPersistenceProfilesClient(sessionContext, connector)
	_, err = client.Get(id)
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving resource", err)
}

func resourceNsxtPolicyLBCookiePersistenceProfileCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	converter := bindings.NewTypeConverter()

	id, err := getOrGenerateID(d, m, resourceNsxtPolicyLBPersistenceProfileExists)
	if err != nil {
		return err
	}

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)

	obj := model.LBCookiePersistenceProfile{
		DisplayName:  &displayName,
		Description:  &description,
		Tags:         tags,
		ResourceType: model.LBPersistenceProfile_RESOURCE_TYPE_LBCOOKIEPERSISTENCEPROFILE,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, lbCookiePersistenceProfileSchema, "", nil); err != nil {
		return err
	}

	dataValue, errs := converter.ConvertToVapi(obj, model.LBCookiePersistenceProfileBindingType())
	if errs != nil {
		return fmt.Errorf("Profile %s is not of type LBCookiepersistenceProfile %s", id, errs[0])
	}

	log.Printf("[INFO] Creating LBCookiePersistenceProfile with ID %s", id)

	sessionContext := getSessionContext(d, m)
	client := cliLbPersistenceProfilesClient(sessionContext, connector)
	err = client.Patch(id, dataValue.(*data.StructValue))
	if err != nil {
		return handleCreateError("LBCookiePersistenceProfile", id, err)
	}
	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyLBCookiePersistenceProfileRead(d, m)
}

func resourceNsxtPolicyLBCookiePersistenceProfileRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	converter := bindings.NewTypeConverter()

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining LBCookiePersistenceProfile ID")
	}

	sessionContext := getSessionContext(d, m)
	client := cliLbPersistenceProfilesClient(sessionContext, connector)

	genObj, err := client.Get(id)
	if err != nil {
		return handleReadError(d, "LBCookiePersistenceProfile", id, err)
	}

	baseObj, errs := converter.ConvertToGolang(genObj, model.LBCookiePersistenceProfileBindingType())
	if len(errs) > 0 {
		return fmt.Errorf("Error converting LBCookiePersistenceProfile %s", errs[0])
	}
	obj := baseObj.(model.LBCookiePersistenceProfile)

	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("revision", obj.Revision)
	d.Set("path", obj.Path)

	elem := reflect.ValueOf(&obj).Elem()
	return metadata.StructToSchema(elem, d, lbCookiePersistenceProfileSchema, "", nil)
}

func resourceNsxtPolicyLBCookiePersistenceProfileUpdate(d *schema.ResourceData, m interface{}) error {

	connector := getPolicyConnector(m)
	converter := bindings.NewTypeConverter()

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining LBCookiePersistenceProfile ID")
	}

	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getPolicyTagsFromSchema(d)

	revision := int64(d.Get("revision").(int))

	obj := model.LBCookiePersistenceProfile{
		DisplayName:  &displayName,
		Description:  &description,
		Tags:         tags,
		Revision:     &revision,
		ResourceType: model.LBPersistenceProfile_RESOURCE_TYPE_LBCOOKIEPERSISTENCEPROFILE,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, lbCookiePersistenceProfileSchema, "", nil); err != nil {
		return err
	}

	dataValue, errs := converter.ConvertToVapi(obj, model.LBCookiePersistenceProfileBindingType())
	if errs != nil {
		return fmt.Errorf("Profile %s is not of type LBCookiepersistenceProfile %s", id, errs[0])
	}

	sessionContext := getSessionContext(d, m)
	client := cliLbPersistenceProfilesClient(sessionContext, connector)
	_, err := client.Update(id, dataValue.(*data.StructValue))
	if err != nil {
		return handleUpdateError("LBCookiePersistenceProfile", id, err)
	}

	return resourceNsxtPolicyLBCookiePersistenceProfileRead(d, m)
}

func resourceNsxtPolicyLBPersistenceProfileDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining LBPersistenceProfile ID")
	}

	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	forceParam := (*bool)(nil)
	client := cliLbPersistenceProfilesClient(sessionContext, connector)
	err := client.Delete(id, forceParam)

	if err != nil {
		return handleDeleteError("LBPersistenceProfile", id, err)
	}

	return nil
}
