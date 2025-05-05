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

	clientLayer "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/metadata"
)

var ipBlockQuotaIpBlockVisibilityValues = []string{
	model.IpBlockQuota_IP_BLOCK_VISIBILITY_PRIVATE,
	model.IpBlockQuota_IP_BLOCK_VISIBILITY_EXTERNAL,
}

var ipBlockQuotaIpBlockAddressTypeValues = []string{
	model.IpBlockQuota_IP_BLOCK_ADDRESS_TYPE_IPV4,
	model.IpBlockQuota_IP_BLOCK_ADDRESS_TYPE_IPV6,
}

var ipBlockQuotaPathExample = getMultitenancyPathExample("/infra/limits/[limit]")

var ipBlockQuotaSchema = map[string]*metadata.ExtendedSchema{
	"nsx_id":       metadata.GetExtendedSchema(getNsxIDSchema()),
	"path":         metadata.GetExtendedSchema(getPathSchema()),
	"display_name": metadata.GetExtendedSchema(getDisplayNameSchema()),
	"description":  metadata.GetExtendedSchema(getDescriptionSchema()),
	"revision":     metadata.GetExtendedSchema(getRevisionSchema()),
	"tag":          metadata.GetExtendedSchema(getTagsSchema()),
	"context":      metadata.GetExtendedSchema(getContextSchema(false, false, false)),
	"quota": {
		Schema: schema.Schema{
			Type:     schema.TypeList,
			MaxItems: 1,
			Required: true,
			Elem: &metadata.ExtendedResource{
				Schema: map[string]*metadata.ExtendedSchema{
					"single_ip_cidrs": {
						Schema: schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Default:  -1,
						},
						Metadata: metadata.Metadata{
							SchemaType:   "int",
							SdkFieldName: "SingleIpCidrs",
						},
					},
					"other_cidrs": {
						Schema: schema.Schema{
							Type:     schema.TypeList,
							MaxItems: 1,
							Required: true,
							Elem: &metadata.ExtendedResource{
								Schema: map[string]*metadata.ExtendedSchema{
									"mask": {
										Schema: schema.Schema{
											Type:     schema.TypeString,
											Optional: true,
										},
										Metadata: metadata.Metadata{
											SchemaType:   "string",
											SdkFieldName: "Mask",
										},
									},
									"total_count": {
										Schema: schema.Schema{
											Type:     schema.TypeInt,
											Optional: true,
											Default:  -1,
										},
										Metadata: metadata.Metadata{
											SchemaType:   "int",
											SdkFieldName: "TotalCount",
										},
									},
								},
							},
						},
						Metadata: metadata.Metadata{
							SchemaType:   "struct",
							SdkFieldName: "OtherCidrs",
							ReflectType:  reflect.TypeOf(model.OtherCidrsMsg{}),
						},
					},
					"ip_block_paths": {
						Schema: schema.Schema{
							Type: schema.TypeList,
							Elem: &metadata.ExtendedSchema{
								Schema: schema.Schema{
									Type:         schema.TypeString,
									ValidateFunc: validatePolicyPath(),
								},
								Metadata: metadata.Metadata{
									SchemaType: "string",
								},
							},
							Optional: true,
						},
						Metadata: metadata.Metadata{
							SchemaType:   "list",
							SdkFieldName: "IpBlockPaths",
						},
					},
					"ip_block_visibility": {
						Schema: schema.Schema{
							Type:         schema.TypeString,
							ValidateFunc: validation.StringInSlice(ipBlockQuotaIpBlockVisibilityValues, false),
							Required:     true,
						},
						Metadata: metadata.Metadata{
							SchemaType:   "string",
							SdkFieldName: "IpBlockVisibility",
						},
					},
					"ip_block_address_type": {
						Schema: schema.Schema{
							Type:         schema.TypeString,
							ValidateFunc: validation.StringInSlice(ipBlockQuotaIpBlockAddressTypeValues, false),
							Required:     true,
						},
						Metadata: metadata.Metadata{
							SchemaType:   "string",
							SdkFieldName: "IpBlockAddressType",
						},
					},
				},
			},
		},
		Metadata: metadata.Metadata{
			SchemaType:      "struct",
			SdkFieldName:    "Quota",
			ReflectType:     reflect.TypeOf(model.IpBlockQuota{}),
			PolymorphicType: metadata.PolymorphicTypeFlatten,
			TypeIdentifier:  metadata.ResourceTypeTypeIdentifier,
			BindingType:     model.IpBlockQuotaBindingType(),
			ResourceType:    model.Quota_RESOURCE_TYPE_IPBLOCKQUOTA,
		},
	},
}

func resourceNsxtPolicyIpBlockQuota() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyIpBlockQuotaCreate,
		Read:   resourceNsxtPolicyIpBlockQuotaRead,
		Update: resourceNsxtPolicyIpBlockQuotaUpdate,
		Delete: resourceNsxtPolicyIpBlockQuotaDelete,
		Importer: &schema.ResourceImporter{
			State: getPolicyPathOrIDResourceImporter(ipBlockQuotaPathExample),
		},
		Schema: metadata.GetSchemaFromExtendedSchema(ipBlockQuotaSchema),
	}
}

func resourceNsxtPolicyIpBlockQuotaExists(sessionContext utl.SessionContext, id string, connector client.Connector) (bool, error) {
	var err error

	client := clientLayer.NewLimitsClient(sessionContext, connector)
	_, err = client.Get(id)
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving resource", err)
}

func resourceNsxtPolicyIpBlockQuotaCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id, err := getOrGenerateID2(d, m, resourceNsxtPolicyIpBlockQuotaExists)
	if err != nil {
		return err
	}

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)

	obj := model.Limit{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, ipBlockQuotaSchema, "", nil); err != nil {
		return err
	}

	log.Printf("[INFO] Creating IpBlockQuota with ID %s", id)

	client := clientLayer.NewLimitsClient(getSessionContext(d, m), connector)
	if client == nil {
		return fmt.Errorf("error creating IpBlockQuota with ID %s - operation is not supported with this backend", id)
	}
	err = client.Patch(id, obj)
	if err != nil {
		return handleCreateError("IpBlockQuota", id, err)
	}
	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyIpBlockQuotaRead(d, m)
}

func resourceNsxtPolicyIpBlockQuotaRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining IpBlockQuota ID")
	}

	client := clientLayer.NewLimitsClient(getSessionContext(d, m), connector)

	obj, err := client.Get(id)
	if err != nil {
		return handleReadError(d, "IpBlockQuota", id, err)
	}

	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("revision", obj.Revision)
	d.Set("path", obj.Path)

	elem := reflect.ValueOf(&obj).Elem()
	return metadata.StructToSchema(elem, d, ipBlockQuotaSchema, "", nil)
}

func resourceNsxtPolicyIpBlockQuotaUpdate(d *schema.ResourceData, m interface{}) error {

	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining IpBlockQuota ID")
	}

	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getPolicyTagsFromSchema(d)

	revision := int64(d.Get("revision").(int))

	obj := model.Limit{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
		Revision:    &revision,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, ipBlockQuotaSchema, "", nil); err != nil {
		return err
	}

	client := clientLayer.NewLimitsClient(getSessionContext(d, m), connector)
	err := client.Patch(id, obj)
	if err != nil {
		return handleUpdateError("IpBlockQuota", id, err)
	}

	return resourceNsxtPolicyIpBlockQuotaRead(d, m)
}

func resourceNsxtPolicyIpBlockQuotaDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining IpBlockQuota ID")
	}

	connector := getPolicyConnector(m)

	client := clientLayer.NewLimitsClient(getSessionContext(d, m), connector)
	err := client.Delete(id)

	if err != nil {
		return handleDeleteError("IpBlockQuota", id, err)
	}

	return nil
}
