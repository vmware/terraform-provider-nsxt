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

	"github.com/vmware/terraform-provider-nsxt/nsxt/metadata"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var vpcLbGenericPersistenceProfileSchema = vpcLbSchemaWithContext(lbGenericPersistenceProfileSchema)

func resourceNsxtVpcLBGenericPersistenceProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtVpcLBGenericPersistenceProfileCreate,
		Read:   resourceNsxtVpcLBGenericPersistenceProfileRead,
		Update: resourceNsxtVpcLBGenericPersistenceProfileUpdate,
		Delete: resourceNsxtVpcLBPersistenceProfileDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtVersionCheckImporter("9.0.0", "VPC LB Generic Persistence Profile", getVpcPathResourceImporter(vpcLbPersistenceProfilePathExample)),
		},
		Schema: metadata.GetSchemaFromExtendedSchema(vpcLbGenericPersistenceProfileSchema),
	}
}

func resourceNsxtVpcLBGenericPersistenceProfileCreate(d *schema.ResourceData, m interface{}) error {
	if !util.NsxVersionHigherOrEqual("9.0.0") {
		return fmt.Errorf("VPC LB Generic Persistence Profile resource requires NSX version 9.0.0 or higher")
	}
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbPersistenceProfilesClient(sessionContext, connector)
	converter := bindings.NewTypeConverter()

	id, err := getOrGenerateID2(d, m, resourceNsxtVpcLBPersistenceProfileExists)
	if err != nil {
		return err
	}

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsWithProviderManagedDefaults(d, m)

	obj := model.LBGenericPersistenceProfile{
		DisplayName:  &displayName,
		Description:  &description,
		Tags:         tags,
		ResourceType: model.LBPersistenceProfile_RESOURCE_TYPE_LBGENERICPERSISTENCEPROFILE,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, vpcLbGenericPersistenceProfileSchema, "", nil); err != nil {
		return err
	}

	dataValue, errs := converter.ConvertToVapi(obj, model.LBGenericPersistenceProfileBindingType())
	if errs != nil {
		return fmt.Errorf("Profile %s is not of type VPC LBGenericPersistenceProfile %s", id, errs[0])
	}

	log.Printf("[INFO] Creating VPC LBGenericPersistenceProfile with ID %s", id)
	err = c.Patch(parents[0], parents[1], parents[2], id, dataValue.(*data.StructValue))
	if err != nil {
		return handleCreateError("VPC LBGenericPersistenceProfile", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)
	return resourceNsxtVpcLBGenericPersistenceProfileRead(d, m)
}

func resourceNsxtVpcLBGenericPersistenceProfileRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbPersistenceProfilesClient(sessionContext, connector)
	converter := bindings.NewTypeConverter()

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VPC LBGenericPersistenceProfile ID")
	}

	genObj, err := c.Get(parents[0], parents[1], parents[2], id)
	if err != nil {
		return handleReadError(d, "VPC LBGenericPersistenceProfile", id, err)
	}

	baseObj, errs := converter.ConvertToGolang(genObj, model.LBGenericPersistenceProfileBindingType())
	if len(errs) > 0 {
		return fmt.Errorf("Error converting VPC LBGenericPersistenceProfile %s", errs[0])
	}
	obj := baseObj.(model.LBGenericPersistenceProfile)

	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("revision", obj.Revision)
	d.Set("path", obj.Path)

	elem := reflect.ValueOf(&obj).Elem()
	return metadata.StructToSchema(elem, d, vpcLbGenericPersistenceProfileSchema, "", nil)
}

func resourceNsxtVpcLBGenericPersistenceProfileUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbPersistenceProfilesClient(sessionContext, connector)
	converter := bindings.NewTypeConverter()

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VPC LBGenericPersistenceProfile ID")
	}

	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getPolicyTagsWithProviderManagedDefaults(d, m)
	revision := int64(d.Get("revision").(int))

	obj := model.LBGenericPersistenceProfile{
		DisplayName:  &displayName,
		Description:  &description,
		Tags:         tags,
		Revision:     &revision,
		ResourceType: model.LBPersistenceProfile_RESOURCE_TYPE_LBGENERICPERSISTENCEPROFILE,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, vpcLbGenericPersistenceProfileSchema, "", nil); err != nil {
		return err
	}

	dataValue, errs := converter.ConvertToVapi(obj, model.LBGenericPersistenceProfileBindingType())
	if errs != nil {
		return fmt.Errorf("Profile %s is not of type VPC LBGenericPersistenceProfile %s", id, errs[0])
	}

	_, err := c.Update(parents[0], parents[1], parents[2], id, dataValue.(*data.StructValue))
	if err != nil {
		return handleUpdateError("VPC LBGenericPersistenceProfile", id, err)
	}

	return resourceNsxtVpcLBGenericPersistenceProfileRead(d, m)
}
