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

	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/metadata"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
	clientLayer "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/orgs/projects/transit_gateways"
)

var connectivityPolicyConnectivityScopeValues = []string{
	model.ConnectivityPolicy_CONNECTIVITY_SCOPE_ISOLATED,
	model.ConnectivityPolicy_CONNECTIVITY_SCOPE_COMMUNITY,
	model.ConnectivityPolicy_CONNECTIVITY_SCOPE_PROMISCUOUS,
}

var connectivityPolicyPathExample = getMultitenancyPathExample("orgs/[org]/projects/[project]/transit-gateways/[transit-gateway]/connectivity-policies/[policy]")

var connectivityPolicySchema = map[string]*metadata.ExtendedSchema{
	"nsx_id":       metadata.GetExtendedSchema(getNsxIDSchema()),
	"path":         metadata.GetExtendedSchema(getPathSchema()),
	"display_name": metadata.GetExtendedSchema(getDisplayNameSchema()),
	"description":  metadata.GetExtendedSchema(getDescriptionSchema()),
	"revision":     metadata.GetExtendedSchema(getRevisionSchema()),
	"tag":          metadata.GetExtendedSchema(getTagsSchema()),
	"parent_path":  metadata.GetExtendedSchema(getPolicyPathSchema(true, true, "Policy path of the parent")),
	"group_path": {
		Schema: schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validatePolicyPath(),
			Required:     true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "string",
			SdkFieldName: "Group",
		},
	},
	"connectivity_scope": {
		Schema: schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validation.StringInSlice(connectivityPolicyConnectivityScopeValues, false),
			Optional:     true,
			ForceNew:     true,
			Default:      model.ConnectivityPolicy_CONNECTIVITY_SCOPE_COMMUNITY,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "string",
			SdkFieldName: "ConnectivityScope",
		},
	},
	"internal_id": {
		Schema: schema.Schema{
			Type:     schema.TypeInt,
			Computed: true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "int",
			SdkFieldName: "InternalId",
			ReadOnly:     true,
		},
	},
}

func resourceNsxtPolicyConnectivityPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyConnectivityPolicyCreate,
		Read:   resourceNsxtPolicyConnectivityPolicyRead,
		Update: resourceNsxtPolicyConnectivityPolicyUpdate,
		Delete: resourceNsxtPolicyConnectivityPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtParentPathResourceImporter,
		},
		Schema: metadata.GetSchemaFromExtendedSchema(connectivityPolicySchema),
	}
}

func resourceNsxtPolicyConnectivityPolicyExists(sessionContext utl.SessionContext, parentPath string, id string, connector client.Connector) (bool, error) {
	var err error
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, connectivityPolicyPathExample)
	if pathErr != nil {
		return false, pathErr
	}
	client := clientLayer.NewConnectivityPoliciesClient(connector)
	_, err = client.Get(parents[0], parents[1], parents[2], id)
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving resource", err)
}

func resourceNsxtPolicyConnectivityPolicyCreate(d *schema.ResourceData, m interface{}) error {

	if !util.NsxVersionHigherOrEqual("9.1.0") {
		return fmt.Errorf("policy Connectivity Policy Create resource requires NSX version 9.1.0 or higher")
	}
	connector := getPolicyConnector(m)

	id, err := getOrGenerateIDWithParent(d, m, resourceNsxtPolicyConnectivityPolicyExists)
	if err != nil {
		return err
	}

	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, connectivityPolicyPathExample)
	if pathErr != nil {
		return pathErr
	}
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)

	obj := model.ConnectivityPolicy{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, connectivityPolicySchema, "", nil); err != nil {
		return err
	}

	log.Printf("[INFO] Creating ConnectivityPolicy with ID %s", id)

	client := clientLayer.NewConnectivityPoliciesClient(connector)
	err = client.Patch(parents[0], parents[1], parents[2], id, obj)
	if err != nil {
		return handleCreateError("ConnectivityPolicy", id, err)
	}
	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyConnectivityPolicyRead(d, m)
}

func resourceNsxtPolicyConnectivityPolicyRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining ConnectivityPolicy ID")
	}

	client := clientLayer.NewConnectivityPoliciesClient(connector)
	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, connectivityPolicyPathExample)
	if pathErr != nil {
		return pathErr
	}
	obj, err := client.Get(parents[0], parents[1], parents[2], id)
	if err != nil {
		return handleReadError(d, "ConnectivityPolicy", id, err)
	}

	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("revision", obj.Revision)
	d.Set("path", obj.Path)

	elem := reflect.ValueOf(&obj).Elem()
	return metadata.StructToSchema(elem, d, connectivityPolicySchema, "", nil)
}

func resourceNsxtPolicyConnectivityPolicyUpdate(d *schema.ResourceData, m interface{}) error {

	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining ConnectivityPolicy ID")
	}

	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, connectivityPolicyPathExample)
	if pathErr != nil {
		return pathErr
	}
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getPolicyTagsFromSchema(d)

	revision := int64(d.Get("revision").(int))

	obj := model.ConnectivityPolicy{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
		Revision:    &revision,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, connectivityPolicySchema, "", nil); err != nil {
		return err
	}
	client := clientLayer.NewConnectivityPoliciesClient(connector)
	_, err := client.Update(parents[0], parents[1], parents[2], id, obj)
	if err != nil {
		return handleUpdateError("ConnectivityPolicy", id, err)
	}

	return resourceNsxtPolicyConnectivityPolicyRead(d, m)
}

func resourceNsxtPolicyConnectivityPolicyDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining ConnectivityPolicy ID")
	}

	connector := getPolicyConnector(m)
	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 3, connectivityPolicyPathExample)
	if pathErr != nil {
		return pathErr
	}

	client := clientLayer.NewConnectivityPoliciesClient(connector)
	err := client.Delete(parents[0], parents[1], parents[2], id)

	if err != nil {
		return handleDeleteError("ConnectivityPolicy", id, err)
	}

	return nil
}
