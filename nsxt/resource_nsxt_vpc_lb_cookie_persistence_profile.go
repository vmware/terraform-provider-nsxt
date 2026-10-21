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
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/metadata"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var vpcLbPersistenceProfilePathExample = "/orgs/[org]/projects/[project]/vpcs/[vpc]/vpc-lb-persistence-profiles/[profile]"

// vpcLbSchemaWithContext returns a copy of the given extended schema with the VPC context added.
func vpcLbSchemaWithContext(base map[string]*metadata.ExtendedSchema) map[string]*metadata.ExtendedSchema {
	s := make(map[string]*metadata.ExtendedSchema, len(base)+1)
	for k, v := range base {
		s[k] = v
	}
	s["context"] = metadata.GetExtendedSchema(getContextSchemaExtended(true, false, true, true))
	return s
}

func resourceNsxtVpcLBPersistenceProfileExists(sessionContext utl.SessionContext, id string, connector client.Connector) (bool, error) {
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbPersistenceProfilesClient(sessionContext, connector)
	_, err := c.Get(parents[0], parents[1], parents[2], id)
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving resource", err)
}

// resourceNsxtVpcLBPersistenceProfileDelete is shared by all VPC LB persistence profile resources
func resourceNsxtVpcLBPersistenceProfileDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VPC LBPersistenceProfile ID")
	}

	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbPersistenceProfilesClient(sessionContext, connector)

	force := true
	err := c.Delete(parents[0], parents[1], parents[2], id, &force)
	if err != nil {
		return handleDeleteError("VPC LBPersistenceProfile", id, err)
	}

	return nil
}

var vpcLbCookiePersistenceProfileSchema = vpcLbSchemaWithContext(lbCookiePersistenceProfileSchema)

func resourceNsxtVpcLBCookiePersistenceProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtVpcLBCookiePersistenceProfileCreate,
		Read:   resourceNsxtVpcLBCookiePersistenceProfileRead,
		Update: resourceNsxtVpcLBCookiePersistenceProfileUpdate,
		Delete: resourceNsxtVpcLBPersistenceProfileDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtVersionCheckImporter("9.0.0", "VPC LB Cookie Persistence Profile", getVpcPathResourceImporter(vpcLbPersistenceProfilePathExample)),
		},
		Schema: metadata.GetSchemaFromExtendedSchema(vpcLbCookiePersistenceProfileSchema),
	}
}

func resourceNsxtVpcLBCookiePersistenceProfileCreate(d *schema.ResourceData, m interface{}) error {
	if !util.NsxVersionHigherOrEqual("9.0.0") {
		return fmt.Errorf("VPC LB Cookie Persistence Profile resource requires NSX version 9.0.0 or higher")
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

	obj := model.LBCookiePersistenceProfile{
		DisplayName:  &displayName,
		Description:  &description,
		Tags:         tags,
		ResourceType: model.LBPersistenceProfile_RESOURCE_TYPE_LBCOOKIEPERSISTENCEPROFILE,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, vpcLbCookiePersistenceProfileSchema, "", nil); err != nil {
		return err
	}

	dataValue, errs := converter.ConvertToVapi(obj, model.LBCookiePersistenceProfileBindingType())
	if errs != nil {
		return fmt.Errorf("Profile %s is not of type VPC LBCookiePersistenceProfile %s", id, errs[0])
	}

	log.Printf("[INFO] Creating VPC LBCookiePersistenceProfile with ID %s", id)
	err = c.Patch(parents[0], parents[1], parents[2], id, dataValue.(*data.StructValue))
	if err != nil {
		return handleCreateError("VPC LBCookiePersistenceProfile", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)
	return resourceNsxtVpcLBCookiePersistenceProfileRead(d, m)
}

func resourceNsxtVpcLBCookiePersistenceProfileRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbPersistenceProfilesClient(sessionContext, connector)
	converter := bindings.NewTypeConverter()

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VPC LBCookiePersistenceProfile ID")
	}

	genObj, err := c.Get(parents[0], parents[1], parents[2], id)
	if err != nil {
		return handleReadError(d, "VPC LBCookiePersistenceProfile", id, err)
	}

	baseObj, errs := converter.ConvertToGolang(genObj, model.LBCookiePersistenceProfileBindingType())
	if len(errs) > 0 {
		return fmt.Errorf("Error converting VPC LBCookiePersistenceProfile %s", errs[0])
	}
	obj := baseObj.(model.LBCookiePersistenceProfile)

	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("revision", obj.Revision)
	d.Set("path", obj.Path)

	elem := reflect.ValueOf(&obj).Elem()
	return metadata.StructToSchema(elem, d, vpcLbCookiePersistenceProfileSchema, "", nil)
}

func resourceNsxtVpcLBCookiePersistenceProfileUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbPersistenceProfilesClient(sessionContext, connector)
	converter := bindings.NewTypeConverter()

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VPC LBCookiePersistenceProfile ID")
	}

	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getPolicyTagsWithProviderManagedDefaults(d, m)
	revision := int64(d.Get("revision").(int))

	obj := model.LBCookiePersistenceProfile{
		DisplayName:  &displayName,
		Description:  &description,
		Tags:         tags,
		Revision:     &revision,
		ResourceType: model.LBPersistenceProfile_RESOURCE_TYPE_LBCOOKIEPERSISTENCEPROFILE,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, vpcLbCookiePersistenceProfileSchema, "", nil); err != nil {
		return err
	}

	dataValue, errs := converter.ConvertToVapi(obj, model.LBCookiePersistenceProfileBindingType())
	if errs != nil {
		return fmt.Errorf("Profile %s is not of type VPC LBCookiePersistenceProfile %s", id, errs[0])
	}

	_, err := c.Update(parents[0], parents[1], parents[2], id, dataValue.(*data.StructValue))
	if err != nil {
		return handleUpdateError("VPC LBCookiePersistenceProfile", id, err)
	}

	return resourceNsxtVpcLBCookiePersistenceProfileRead(d, m)
}
