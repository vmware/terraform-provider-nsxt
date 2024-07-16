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
	clientLayer "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/orgs/projects/vpcs"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/metadata"
)

var vpcSubnetAccessModeValues = []string{
	model.VpcSubnet_ACCESS_MODE_PRIVATE,
	model.VpcSubnet_ACCESS_MODE_PUBLIC,
	model.VpcSubnet_ACCESS_MODE_ISOLATED,
}

var vpcSubnetSchema = map[string]*metadata.ExtendedSchema{
	"nsx_id":       metadata.GetExtendedSchema(getNsxIDSchema()),
	"path":         metadata.GetExtendedSchema(getPathSchema()),
	"display_name": metadata.GetExtendedSchema(getDisplayNameSchema()),
	"description":  metadata.GetExtendedSchema(getDescriptionSchema()),
	"revision":     metadata.GetExtendedSchema(getRevisionSchema()),
	"tag":          metadata.GetExtendedSchema(getTagsSchema()),
	"context":      metadata.GetExtendedSchema(getContextSchema(true, false, true)),
	"advanced_config": {
		Schema: schema.Schema{
			Type:     schema.TypeList,
			MaxItems: 1,
			Elem: &metadata.ExtendedResource{
				Schema: map[string]*metadata.ExtendedSchema{
					"static_ip_allocation": {
						Schema: schema.Schema{
							Type:     schema.TypeList,
							MaxItems: 1,
							Elem: &metadata.ExtendedResource{
								Schema: map[string]*metadata.ExtendedSchema{
									"enabled": {
										Schema: schema.Schema{
											Type:     schema.TypeBool,
											Optional: true,
										},
										Metadata: metadata.Metadata{
											SchemaType:   "bool",
											SdkFieldName: "Enabled",
										},
									},
								},
							},
							Optional: true,
						},
						Metadata: metadata.Metadata{
							SchemaType:   "struct",
							SdkFieldName: "StaticIpAllocation",
							ReflectType:  reflect.TypeOf(model.StaticIpAllocation{}),
						},
					},
				},
			},
			Optional: true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "struct",
			SdkFieldName: "AdvancedConfig",
			ReflectType:  reflect.TypeOf(model.SubnetAdvancedConfig{}),
		},
	},
	"ipv4_subnet_size": {
		Schema: schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "int",
			SdkFieldName: "Ipv4SubnetSize",
			OmitIfEmpty:  true,
		},
	},
	"ip_addresses": {
		Schema: schema.Schema{
			Type: schema.TypeList,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional: true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "array",
			SdkFieldName: "IpAddresses",
		},
	},
	"access_mode": {
		Schema: schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validation.StringInSlice(vpcSubnetAccessModeValues, false),
			Optional:     true,
			Default:      model.VpcSubnet_ACCESS_MODE_PRIVATE,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "string",
			SdkFieldName: "AccessMode",
		},
	},
	"dhcp_config": {
		Schema: schema.Schema{
			Type:     schema.TypeList,
			MaxItems: 1,
			Computed: true,
			Elem: &metadata.ExtendedResource{
				Schema: map[string]*metadata.ExtendedSchema{
					"dhcp_relay_config_path": {
						Schema: schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						Metadata: metadata.Metadata{
							SchemaType:   "string",
							SdkFieldName: "DhcpRelayConfigPath",
						},
					},
					"dns_client_config": {
						Schema: schema.Schema{
							Type:     schema.TypeList,
							MaxItems: 1,
							Elem: &metadata.ExtendedResource{
								Schema: map[string]*metadata.ExtendedSchema{
									"dns_server_ips": {
										Schema: schema.Schema{
											Type: schema.TypeList,
											Elem: &schema.Schema{
												Type: schema.TypeString,
											},
											Optional: true,
										},
										Metadata: metadata.Metadata{
											SchemaType:   "array",
											SdkFieldName: "DnsServerIps",
										},
									},
								},
							},
							Optional: true,
						},
						Metadata: metadata.Metadata{
							SchemaType:   "struct",
							SdkFieldName: "DnsClientConfig",
							ReflectType:  reflect.TypeOf(model.DnsClientConfig{}),
						},
					},
					"enable_dhcp": {
						Schema: schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						Metadata: metadata.Metadata{
							SchemaType:   "bool",
							SdkFieldName: "EnableDhcp",
						},
					},
					"static_pool_config": {
						Schema: schema.Schema{
							Type:     schema.TypeList,
							MaxItems: 1,
							Elem: &metadata.ExtendedResource{
								Schema: map[string]*metadata.ExtendedSchema{
									"ipv4_pool_size": {
										Schema: schema.Schema{
											Type:     schema.TypeInt,
											Optional: true,
										},
										Metadata: metadata.Metadata{
											SchemaType:   "int",
											SdkFieldName: "Ipv4PoolSize",
										},
									},
								},
							},
							Optional: true,
						},
						Metadata: metadata.Metadata{
							SchemaType:   "struct",
							SdkFieldName: "StaticPoolConfig",
							ReflectType:  reflect.TypeOf(model.StaticPoolConfig{}),
						},
					},
				},
			},
			Optional: true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "struct",
			SdkFieldName: "DhcpConfig",
			ReflectType:  reflect.TypeOf(model.VpcSubnetDhcpConfig{}),
		},
	},
}

func resourceNsxtVpcSubnet() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtVpcSubnetCreate,
		Read:   resourceNsxtVpcSubnetRead,
		Update: resourceNsxtVpcSubnetUpdate,
		Delete: resourceNsxtVpcSubnetDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtVPCPathResourceImporter,
		},
		Schema: metadata.GetSchemaFromExtendedSchema(vpcSubnetSchema),
	}
}

func resourceNsxtVpcSubnetExists(sessionContext utl.SessionContext, id string, connector client.Connector) (bool, error) {
	var err error
	parents := getVpcParentsFromContext(sessionContext)
	client := clientLayer.NewSubnetsClient(connector)
	_, err = client.Get(parents[0], parents[1], parents[2], id)
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving resource", err)
}

func resourceNsxtVpcSubnetCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id, err := getOrGenerateID2(d, m, resourceNsxtVpcSubnetExists)
	if err != nil {
		return err
	}

	parents := getVpcParentsFromContext(getSessionContext(d, m))
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)

	obj := model.VpcSubnet{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, vpcSubnetSchema, "", nil); err != nil {
		return err
	}

	log.Printf("[INFO] Creating VpcSubnet with ID %s", id)

	client := clientLayer.NewSubnetsClient(connector)
	err = client.Patch(parents[0], parents[1], parents[2], id, obj)
	if err != nil {
		return handleCreateError("VpcSubnet", id, err)
	}
	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtVpcSubnetRead(d, m)
}

func resourceNsxtVpcSubnetRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VpcSubnet ID")
	}

	client := clientLayer.NewSubnetsClient(connector)
	parents := getVpcParentsFromContext(getSessionContext(d, m))
	obj, err := client.Get(parents[0], parents[1], parents[2], id)
	if err != nil {
		return handleReadError(d, "VpcSubnet", id, err)
	}

	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("revision", obj.Revision)
	d.Set("path", obj.Path)

	elem := reflect.ValueOf(&obj).Elem()
	return metadata.StructToSchema(elem, d, vpcSubnetSchema, "", nil)
}

func resourceNsxtVpcSubnetUpdate(d *schema.ResourceData, m interface{}) error {

	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VpcSubnet ID")
	}

	parents := getVpcParentsFromContext(getSessionContext(d, m))
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getPolicyTagsFromSchema(d)

	revision := int64(d.Get("revision").(int))

	obj := model.VpcSubnet{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
		Revision:    &revision,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, vpcSubnetSchema, "", nil); err != nil {
		return err
	}
	client := clientLayer.NewSubnetsClient(connector)
	_, err := client.Update(parents[0], parents[1], parents[2], id, obj)
	if err != nil {
		return handleUpdateError("VpcSubnet", id, err)
	}

	return resourceNsxtVpcSubnetRead(d, m)
}

func resourceNsxtVpcSubnetDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining VpcSubnet ID")
	}

	connector := getPolicyConnector(m)
	parents := getVpcParentsFromContext(getSessionContext(d, m))

	client := clientLayer.NewSubnetsClient(connector)
	err := client.Delete(parents[0], parents[1], parents[2], id)

	if err != nil {
		return handleDeleteError("VpcSubnet", id, err)
	}

	return nil
}
