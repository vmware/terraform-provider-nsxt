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
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	"github.com/vmware/terraform-provider-nsxt/nsxt/metadata"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var networkSpanPathExample = getMultitenancyPathExample("/infra/network-spans/[networkSpan]")

var networkSpanSchema = map[string]*metadata.ExtendedSchema{
	"nsx_id":       metadata.GetExtendedSchema(getNsxIDSchema()),
	"path":         metadata.GetExtendedSchema(getPathSchema()),
	"display_name": metadata.GetExtendedSchema(getDisplayNameSchema()),
	"description":  metadata.GetExtendedSchema(getDescriptionSchema()),
	"revision":     metadata.GetExtendedSchema(getRevisionSchema()),
	"tag":          metadata.GetExtendedSchema(getTagsSchema()),
	"exclusive": {
		Schema: schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "bool",
			SdkFieldName: "Exclusive",
		},
	},
}

func resourceNsxtPolicyNetworkSpan() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyNetworkSpanCreate,
		Read:   resourceNsxtPolicyNetworkSpanRead,
		Update: resourceNsxtPolicyNetworkSpanUpdate,
		Delete: resourceNsxtPolicyNetworkSpanDelete,
		Importer: &schema.ResourceImporter{
			State: getPolicyPathOrIDResourceImporter(networkSpanPathExample),
		},
		Schema: metadata.GetSchemaFromExtendedSchema(networkSpanSchema),
	}
}

func resourceNsxtPolicyNetworkSpanExists(id string, connector client.Connector, isGlobalManager bool) (bool, error) {
	var err error

	client := infra.NewNetworkSpansClient(connector)
	_, err = client.Get(id)
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving resource", err)
}

func resourceNsxtPolicyNetworkSpanCreate(d *schema.ResourceData, m interface{}) error {
	if !util.NsxVersionHigherOrEqual("9.1.0") {
		return fmt.Errorf("Policy NetworkSpan resource requires NSX version 9.1.0 or higher")
	}
	connector := getPolicyConnector(m)

	id, err := getOrGenerateID(d, m, resourceNsxtPolicyNetworkSpanExists)
	if err != nil {
		return err
	}

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)

	obj := model.NetworkSpan{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, networkSpanSchema, "", nil); err != nil {
		return err
	}

	log.Printf("[INFO] Creating NetworkSpan with ID %s", id)

	client := infra.NewNetworkSpansClient(connector)
	err = client.Patch(id, obj)
	if err != nil {
		return handleCreateError("NetworkSpan", id, err)
	}
	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyNetworkSpanRead(d, m)
}

func resourceNsxtPolicyNetworkSpanRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining NetworkSpan ID")
	}

	client := infra.NewNetworkSpansClient(connector)

	obj, err := client.Get(id)
	if err != nil {
		return handleReadError(d, "NetworkSpan", id, err)
	}

	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("revision", obj.Revision)
	d.Set("path", obj.Path)

	elem := reflect.ValueOf(&obj).Elem()
	return metadata.StructToSchema(elem, d, networkSpanSchema, "", nil)
}

func resourceNsxtPolicyNetworkSpanUpdate(d *schema.ResourceData, m interface{}) error {

	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining NetworkSpan ID")
	}

	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getPolicyTagsFromSchema(d)

	revision := int64(d.Get("revision").(int))

	obj := model.NetworkSpan{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
		Revision:    &revision,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, networkSpanSchema, "", nil); err != nil {
		return err
	}
	client := infra.NewNetworkSpansClient(connector)
	_, err := client.Update(id, obj)
	if err != nil {
		return handleUpdateError("NetworkSpan", id, err)
	}

	return resourceNsxtPolicyNetworkSpanRead(d, m)
}

func resourceNsxtPolicyNetworkSpanDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining NetworkSpan ID")
	}

	connector := getPolicyConnector(m)

	client := infra.NewNetworkSpansClient(connector)
	err := client.Delete(id)

	if err != nil {
		return handleDeleteError("NetworkSpan", id, err)
	}

	return nil
}
