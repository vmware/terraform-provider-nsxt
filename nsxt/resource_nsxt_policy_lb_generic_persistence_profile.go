// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	"github.com/vmware/terraform-provider-nsxt/api/infra"
	"github.com/vmware/terraform-provider-nsxt/nsxt/metadata"
)

var lbGenericPersistenceProfileSchema = map[string]*metadata.ExtendedSchema{
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
	"timeout": {
		Schema: schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			Default:  300,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "int",
			SdkFieldName: "Timeout",
		},
	},
	"ha_persistence_mirroring_enabled": {
		Schema: schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "bool",
			SdkFieldName: "HaPersistenceMirroringEnabled",
		},
	},
}

func resourceNsxtPolicyLBGenericPersistenceProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyLBGenericPersistenceProfileCreate,
		Read:   resourceNsxtPolicyLBGenericPersistenceProfileRead,
		Update: resourceNsxtPolicyLBGenericPersistenceProfileUpdate,
		Delete: resourceNsxtPolicyLBGenericPersistenceProfileDelete,
		Importer: &schema.ResourceImporter{
			State: getPolicyPathOrIDResourceImporter(lbPersistenceProfilePathExample),
		},
		Schema: metadata.GetSchemaFromExtendedSchema(lbGenericPersistenceProfileSchema),
	}
}

func resourceNsxtPolicyLBGenericPersistenceProfileCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	converter := bindings.NewTypeConverter()

	id, err := getOrGenerateID(d, m, resourceNsxtPolicyLBPersistenceProfileExists)
	if err != nil {
		return err
	}

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)

	obj := model.LBGenericPersistenceProfile{
		DisplayName:  &displayName,
		Description:  &description,
		Tags:         tags,
		ResourceType: model.LBPersistenceProfile_RESOURCE_TYPE_LBGENERICPERSISTENCEPROFILE,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, lbGenericPersistenceProfileSchema, "", nil); err != nil {
		return err
	}

	dataValue, errs := converter.ConvertToVapi(obj, model.LBGenericPersistenceProfileBindingType())
	if errs != nil {
		return fmt.Errorf("Profile %s is not of type LBGenericPersistenceProfile %s", id, errs[0])
	}

	log.Printf("[INFO] Creating LBGenericPersistenceProfile with ID %s", id)

	sessionContext := getSessionContext(d, m)
	client := infra.NewLbPersistenceProfilesClient(sessionContext, connector)
	err = client.Patch(id, dataValue.(*data.StructValue))
	if err != nil {
		return handleCreateError("LBGenericPersistenceProfile", id, err)
	}
	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyLBGenericPersistenceProfileRead(d, m)
}

func resourceNsxtPolicyLBGenericPersistenceProfileRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	converter := bindings.NewTypeConverter()

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining LBGenericPersistenceProfile ID")
	}

	sessionContext := getSessionContext(d, m)
	client := infra.NewLbPersistenceProfilesClient(sessionContext, connector)

	genObj, err := client.Get(id)
	if err != nil {
		return handleReadError(d, "LBGenericPersistenceProfile", id, err)
	}

	baseObj, errs := converter.ConvertToGolang(genObj, model.LBGenericPersistenceProfileBindingType())
	if len(errs) > 0 {
		return fmt.Errorf("Error converting LBGenericPersistenceProfile %s", errs[0])
	}
	obj := baseObj.(model.LBGenericPersistenceProfile)

	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("revision", obj.Revision)
	d.Set("path", obj.Path)

	elem := reflect.ValueOf(&obj).Elem()
	return metadata.StructToSchema(elem, d, lbGenericPersistenceProfileSchema, "", nil)
}

func resourceNsxtPolicyLBGenericPersistenceProfileUpdate(d *schema.ResourceData, m interface{}) error {

	connector := getPolicyConnector(m)
	converter := bindings.NewTypeConverter()

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining LBGenericPersistenceProfile ID")
	}

	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getPolicyTagsFromSchema(d)

	revision := int64(d.Get("revision").(int))

	obj := model.LBGenericPersistenceProfile{
		DisplayName:  &displayName,
		Description:  &description,
		Tags:         tags,
		ResourceType: model.LBPersistenceProfile_RESOURCE_TYPE_LBGENERICPERSISTENCEPROFILE,
		Revision:     &revision,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, lbGenericPersistenceProfileSchema, "", nil); err != nil {
		return err
	}
	dataValue, errs := converter.ConvertToVapi(obj, model.LBGenericPersistenceProfileBindingType())
	if errs != nil {
		return fmt.Errorf("Profile %s is not of type LBGenericPersistenceProfile %s", id, errs[0])
	}

	sessionContext := getSessionContext(d, m)
	client := infra.NewLbPersistenceProfilesClient(sessionContext, connector)
	_, err := client.Update(id, dataValue.(*data.StructValue))
	if err != nil {
		return handleUpdateError("LBGenericPersistenceProfile", id, err)
	}

	return resourceNsxtPolicyLBGenericPersistenceProfileRead(d, m)
}

func resourceNsxtPolicyLBGenericPersistenceProfileDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining LBGenericPersistenceProfile ID")
	}

	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	client := infra.NewLbPersistenceProfilesClient(sessionContext, connector)
	err := client.Delete(id, nil)

	if err != nil {
		return handleDeleteError("LBGenericPersistenceProfile", id, err)
	}

	return nil
}
