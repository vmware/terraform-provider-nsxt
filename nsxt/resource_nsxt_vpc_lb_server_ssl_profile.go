// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/metadata"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var vpcLbServerSslProfilePathExample = "/orgs/[org]/projects/[project]/vpcs/[vpc]/vpc-lb-server-ssl-profiles/[profile]"

// Same fields as the infra LB server SSL profile, with the VPC context
var vpcLbServerSslProfileSchema = vpcLbSchemaWithContext(lbServerSslProfileSchema)

func resourceNsxtVpcLBServerSslProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtVpcLBServerSslProfileCreate,
		Read:   resourceNsxtVpcLBServerSslProfileRead,
		Update: resourceNsxtVpcLBServerSslProfileUpdate,
		Delete: resourceNsxtVpcLBServerSslProfileDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtVersionCheckImporter("9.0.0", "VPC LB Server SSL Profile", getVpcPathResourceImporter(vpcLbServerSslProfilePathExample)),
		},
		Schema: metadata.GetSchemaFromExtendedSchema(vpcLbServerSslProfileSchema),
	}
}

func resourceNsxtVpcLBServerSslProfileExists(sessionContext utl.SessionContext, id string, connector client.Connector) (bool, error) {
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbServerSslProfilesClient(sessionContext, connector)
	_, err := c.Get(parents[0], parents[1], parents[2], id)
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving resource", err)
}

func resourceNsxtVpcLBServerSslProfileCreate(d *schema.ResourceData, m interface{}) error {
	if !util.NsxVersionHigherOrEqual("9.0.0") {
		return fmt.Errorf("VPC LB Server SSL Profile resource requires NSX version 9.0.0 or higher")
	}
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbServerSslProfilesClient(sessionContext, connector)

	id, err := getOrGenerateID2(d, m, resourceNsxtVpcLBServerSslProfileExists)
	if err != nil {
		return err
	}

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsWithProviderManagedDefaults(d, m)

	obj := model.LBServerSslProfile{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, vpcLbServerSslProfileSchema, "", nil); err != nil {
		return err
	}

	log.Printf("[INFO] Creating VPC LBServerSslProfile with ID %s", id)
	err = c.Patch(parents[0], parents[1], parents[2], id, obj)
	if err != nil {
		return handleCreateError("VPC LBServerSslProfile", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)
	return resourceNsxtVpcLBServerSslProfileRead(d, m)
}

func resourceNsxtVpcLBServerSslProfileRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbServerSslProfilesClient(sessionContext, connector)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VPC LBServerSslProfile ID")
	}

	obj, err := c.Get(parents[0], parents[1], parents[2], id)
	if err != nil {
		return handleReadError(d, "VPC LBServerSslProfile", id, err)
	}

	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("revision", obj.Revision)
	d.Set("path", obj.Path)

	elem := reflect.ValueOf(&obj).Elem()
	return metadata.StructToSchema(elem, d, vpcLbServerSslProfileSchema, "", nil)
}

func resourceNsxtVpcLBServerSslProfileUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbServerSslProfilesClient(sessionContext, connector)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VPC LBServerSslProfile ID")
	}

	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getPolicyTagsWithProviderManagedDefaults(d, m)
	revision := int64(d.Get("revision").(int))

	obj := model.LBServerSslProfile{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
		Revision:    &revision,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, vpcLbServerSslProfileSchema, "", nil); err != nil {
		return err
	}

	_, err := c.Update(parents[0], parents[1], parents[2], id, obj)
	if err != nil {
		return handleUpdateError("VPC LBServerSslProfile", id, err)
	}

	return resourceNsxtVpcLBServerSslProfileRead(d, m)
}

func resourceNsxtVpcLBServerSslProfileDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VPC LBServerSslProfile ID")
	}

	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbServerSslProfilesClient(sessionContext, connector)

	force := true
	err := c.Delete(parents[0], parents[1], parents[2], id, &force)
	if err != nil {
		return handleDeleteError("VPC LBServerSslProfile", id, err)
	}

	return nil
}
