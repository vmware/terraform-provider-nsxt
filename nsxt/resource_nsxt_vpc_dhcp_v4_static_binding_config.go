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
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	clientLayer "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/orgs/projects/vpcs/subnets"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/metadata"
)

var dhcpV4StaticBindingConfigResourceTypeValues = []string{
	model.DhcpV4StaticBindingConfig_RESOURCE_TYPE_DHCP_V4_STATIC_BINDING_CONFIG,
	model.DhcpV4StaticBindingConfig_RESOURCE_TYPE_DHCP_V6_STATIC_BINDING_CONFIG,
}

var dhcpV4StaticBindingConfigSchema = map[string]*metadata.ExtendedSchema{
	"nsx_id":       metadata.GetExtendedSchema(getNsxIDSchema()),
	"path":         metadata.GetExtendedSchema(getPathSchema()),
	"display_name": metadata.GetExtendedSchema(getDisplayNameSchema()),
	"description":  metadata.GetExtendedSchema(getDescriptionSchema()),
	"revision":     metadata.GetExtendedSchema(getRevisionSchema()),
	"tag":          metadata.GetExtendedSchema(getTagsSchema()),
	"parent_path":  metadata.GetExtendedSchema(getPolicyPathSchema(true, true, "Policy path of the parent")),
	"resource_type": {
		Schema: schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validation.StringInSlice(dhcpV4StaticBindingConfigResourceTypeValues, false),
			Required:     true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "string",
			SdkFieldName: "ResourceType",
		},
	},
	"gateway_address": {
		Schema: schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validation.IsIPv4Address,
			Optional:     true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "string",
			SdkFieldName: "GatewayAddress",
		},
	},
	"host_name": {
		Schema: schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "string",
			SdkFieldName: "HostName",
		},
	},
	"mac_address": {
		Schema: schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validation.IsMACAddress,
			Optional:     true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "string",
			SdkFieldName: "MacAddress",
		},
	},
	"lease_time": {
		Schema: schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			Default:  86400,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "int",
			SdkFieldName: "LeaseTime",
		},
	},
	"ip_address": {
		Schema: schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validation.IsIPv4Address,
			Optional:     true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "string",
			SdkFieldName: "IpAddress",
		},
	},
	"options": {
		Schema: schema.Schema{
			Type:     schema.TypeList,
			MaxItems: 1,
			Elem: &metadata.ExtendedResource{
				Schema: map[string]*metadata.ExtendedSchema{
					"option121": {
						Schema: schema.Schema{
							Type:     schema.TypeList,
							MaxItems: 1,
							Elem: &metadata.ExtendedResource{
								Schema: map[string]*metadata.ExtendedSchema{
									"static_route": {
										Schema: schema.Schema{
											Type: schema.TypeList,
											Elem: &metadata.ExtendedResource{
												Schema: map[string]*metadata.ExtendedSchema{
													"next_hop": {
														Schema: schema.Schema{
															Type:         schema.TypeString,
															ValidateFunc: validateSingleIP(),
															Optional:     true,
														},
														Metadata: metadata.Metadata{
															SchemaType:   "string",
															SdkFieldName: "NextHop",
														},
													},
													"network": {
														Schema: schema.Schema{
															Type:         schema.TypeString,
															ValidateFunc: validateCidrOrIPOrRange(),
															Optional:     true,
														},
														Metadata: metadata.Metadata{
															SchemaType:   "string",
															SdkFieldName: "Network",
														},
													},
												},
											},
											Optional: true,
										},
										Metadata: metadata.Metadata{
											SchemaType:   "list",
											SdkFieldName: "StaticRoutes",
										},
									},
								},
							},
							Optional: true,
						},
						Metadata: metadata.Metadata{
							SchemaType:   "struct",
							SdkFieldName: "Option121",
							ReflectType:  reflect.TypeOf(model.DhcpOption121{}),
						},
					},
					"other": {
						Schema: schema.Schema{
							Type: schema.TypeList,
							Elem: &metadata.ExtendedResource{
								Schema: map[string]*metadata.ExtendedSchema{
									"code": {
										Schema: schema.Schema{
											Type:     schema.TypeInt,
											Optional: true,
										},
										Metadata: metadata.Metadata{
											SchemaType:   "int",
											SdkFieldName: "Code",
										},
									},
									"values": {
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
											SdkFieldName: "Values",
										},
									},
								},
							},
							Optional: true,
						},
						Metadata: metadata.Metadata{
							SchemaType:   "list",
							SdkFieldName: "Others",
						},
					},
				},
			},
			Optional: true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "struct",
			SdkFieldName: "Options",
			ReflectType:  reflect.TypeOf(model.DhcpV4Options{}),
		},
	},
}

func resourceNsxtDhcpV4StaticBindingConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtDhcpV4StaticBindingConfigCreate,
		Read:   resourceNsxtDhcpV4StaticBindingConfigRead,
		Update: resourceNsxtDhcpV4StaticBindingConfigUpdate,
		Delete: resourceNsxtDhcpV4StaticBindingConfigDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtParentPathResourceImporter,
		},
		Schema: metadata.GetSchemaFromExtendedSchema(dhcpV4StaticBindingConfigSchema),
	}
}

func resourceNsxtDhcpV4StaticBindingConfigExists(sessionContext utl.SessionContext, parentPath string, id string, connector client.Connector) (bool, error) {
	var err error
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 4)
	if pathErr != nil {
		return false, pathErr
	}
	client := clientLayer.NewDhcpV4StaticBindingConfigsClient(connector)
	_, err = client.Get(parents[0], parents[1], parents[2], parents[3], id)
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving resource", err)
}

func resourceNsxtDhcpV4StaticBindingConfigCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id, err := getOrGenerateIDWithParent(d, m, resourceNsxtDhcpV4StaticBindingConfigExists)
	if err != nil {
		return err
	}

	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 4)
	if pathErr != nil {
		return pathErr
	}
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)

	obj := model.DhcpV4StaticBindingConfig{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, dhcpV4StaticBindingConfigSchema, "", nil); err != nil {
		return err
	}

	log.Printf("[INFO] Creating DhcpV4StaticBindingConfig with ID %s", id)

	client := clientLayer.NewDhcpV4StaticBindingConfigsClient(connector)
	err = client.Patch(parents[0], parents[1], parents[2], parents[3], id, obj)
	if err != nil {
		return handleCreateError("DhcpV4StaticBindingConfig", id, err)
	}
	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtDhcpV4StaticBindingConfigRead(d, m)
}

func resourceNsxtDhcpV4StaticBindingConfigRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining DhcpV4StaticBindingConfig ID")
	}

	client := clientLayer.NewDhcpV4StaticBindingConfigsClient(connector)
	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 4)
	if pathErr != nil {
		return pathErr
	}
	obj, err := client.Get(parents[0], parents[1], parents[2], parents[3], id)
	if err != nil {
		return handleReadError(d, "DhcpV4StaticBindingConfig", id, err)
	}

	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("revision", obj.Revision)
	d.Set("path", obj.Path)

	elem := reflect.ValueOf(&obj).Elem()
	return metadata.StructToSchema(elem, d, dhcpV4StaticBindingConfigSchema, "", nil)
}

func resourceNsxtDhcpV4StaticBindingConfigUpdate(d *schema.ResourceData, m interface{}) error {

	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining DhcpV4StaticBindingConfig ID")
	}

	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 4)
	if pathErr != nil {
		return pathErr
	}
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getPolicyTagsFromSchema(d)

	revision := int64(d.Get("revision").(int))

	obj := model.DhcpV4StaticBindingConfig{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
		Revision:    &revision,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, dhcpV4StaticBindingConfigSchema, "", nil); err != nil {
		return err
	}
	client := clientLayer.NewDhcpV4StaticBindingConfigsClient(connector)
	_, err := client.Update(parents[0], parents[1], parents[2], parents[3], id, obj)
	if err != nil {
		return handleUpdateError("DhcpV4StaticBindingConfig", id, err)
	}

	return resourceNsxtDhcpV4StaticBindingConfigRead(d, m)
}

func resourceNsxtDhcpV4StaticBindingConfigDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining DhcpV4StaticBindingConfig ID")
	}

	connector := getPolicyConnector(m)
	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 4)
	if pathErr != nil {
		return pathErr
	}

	client := clientLayer.NewDhcpV4StaticBindingConfigsClient(connector)
	err := client.Delete(parents[0], parents[1], parents[2], parents[3], id)

	if err != nil {
		return handleDeleteError("DhcpV4StaticBindingConfig", id, err)
	}

	return nil
}
