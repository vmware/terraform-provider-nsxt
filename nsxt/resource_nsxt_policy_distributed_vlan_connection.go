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
	clientLayer "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	"github.com/vmware/terraform-provider-nsxt/nsxt/metadata"
)

var distributedVlanConnectionSchema = map[string]*metadata.ExtendedSchema{
	"nsx_id":       metadata.GetExtendedSchema(getNsxIDSchema()),
	"path":         metadata.GetExtendedSchema(getPathSchema()),
	"display_name": metadata.GetExtendedSchema(getDisplayNameSchema()),
	"description":  metadata.GetExtendedSchema(getDescriptionSchema()),
	"revision":     metadata.GetExtendedSchema(getRevisionSchema()),
	"tag":          metadata.GetExtendedSchema(getTagsSchema()),
	"vlan_id": {
		Schema: schema.Schema{
			Type:     schema.TypeInt,
			Required: true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "int",
			SdkFieldName: "VlanId",
		},
	},
	"subnet_exclusive_config": metadata.GetExtendedSchema(subnetExclusiveConfigSchema()),
	"gateway_addresses": {
		Schema: schema.Schema{
			Type: schema.TypeList,
			Elem: &metadata.ExtendedSchema{
				Schema: schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateIPCidr(),
				},
				Metadata: metadata.Metadata{
					SchemaType: "string",
				},
			},
			Required: true,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "list",
			SdkFieldName: "GatewayAddresses",
		},
	},
}

func subnetExclusiveConfigSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "Subnet exclusive config",
		Optional:    true,
		ForceNew:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"ip_block_path": {
					Type:        schema.TypeString,
					Description: "Policy path of external IP block. This IP block must be marked as reserved for VLAN extension.",
					Optional:    true,
					ForceNew:    true,
				},
				"vlan_extension": vlanExtensionSchema(),
			},
		},
	}

}

func vlanExtensionSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "Specifies whether VLAN extension and VPC gateway connectivity are enabled for the VPC subnet.",
		Optional:    true,
		ForceNew:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"vpc_gateway_connection_enable": {
					Type:        schema.TypeBool,
					Description: "This configuration controls whether the VLAN extension subnet connects to the VPC gateway.",
					Optional:    true,
					ForceNew:    true,
				},
			},
		},
	}

}

const distributedVlanConnectionPathExample = "/infra/distributed-vlan-connections/[connection]"

func resourceNsxtPolicyDistributedVlanConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyDistributedVlanConnectionCreate,
		Read:   resourceNsxtPolicyDistributedVlanConnectionRead,
		Update: resourceNsxtPolicyDistributedVlanConnectionUpdate,
		Delete: resourceNsxtPolicyDistributedVlanConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: getPolicyPathResourceImporter(distributedVlanConnectionPathExample),
		},
		Schema: metadata.GetSchemaFromExtendedSchema(distributedVlanConnectionSchema),
	}
}

func resourceNsxtPolicyDistributedVlanConnectionExists(id string, connector client.Connector, isGlobalManager bool) (bool, error) {
	var err error

	client := clientLayer.NewDistributedVlanConnectionsClient(connector)
	_, err = client.Get(id)
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving resource", err)
}

func resourceNsxtPolicyDistributedVlanConnectionCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id, err := getOrGenerateID(d, m, resourceNsxtPolicyDistributedVlanConnectionExists)
	if err != nil {
		return err
	}

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)

	obj := model.DistributedVlanConnection{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
	}

	subnetExclusiveConfig := model.SubnetExclusiveConfig{}
	subnetExclusiveConfigArgUnparsed := d.Get("subnet_exclusive_config").([]interface{})
	if len(subnetExclusiveConfigArgUnparsed) > 0 || subnetExclusiveConfigArgUnparsed[0] != nil {
		tempSubnetExclusiveConfig, mapExists := subnetExclusiveConfigArgUnparsed[0].(map[string]interface{})
		if ipBlockPath, ok := tempSubnetExclusiveConfig["ip_block_path"]; mapExists && ok {
			ipBlockPathStr := ipBlockPath.(string)
			subnetExclusiveConfig.IpBlockPath = &ipBlockPathStr
		}
		if vlanExtensionSet, ok := tempSubnetExclusiveConfig["vlan_extension"]; ok && len(vlanExtensionSet.([]interface{})) > 0 {
			vlanExtension := vlanExtensionSet.([]interface{})[0]
			if vpcGatewayConnectionEnableMap, ok := vlanExtension.(map[string]interface{}); ok {
				vpcGatewayConnectionEnable, BoolVarExists := vpcGatewayConnectionEnableMap["vpc_gateway_connection_enable"]
				if vpcGatewayConnectionEnablePtr, ok := vpcGatewayConnectionEnable.(bool); BoolVarExists && ok {
					subnetExclusiveConfig.VlanExtension = &model.VlanExtension{
						VpcGatewayConnectionEnable: &vpcGatewayConnectionEnablePtr,
					}
				}
			}
		}
		obj.SubnetExclusiveConfig = &subnetExclusiveConfig
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, distributedVlanConnectionSchema, "", nil); err != nil {
		return err
	}

	log.Printf("[INFO] Creating DistributedVlanConnection with ID %s", id)

	client := clientLayer.NewDistributedVlanConnectionsClient(connector)
	err = client.Patch(id, obj)
	if err != nil {
		return handleCreateError("DistributedVlanConnection", id, err)
	}
	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyDistributedVlanConnectionRead(d, m)
}

func resourceNsxtPolicyDistributedVlanConnectionRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining DistributedVlanConnection ID")
	}

	client := clientLayer.NewDistributedVlanConnectionsClient(connector)

	obj, err := client.Get(id)
	if err != nil {
		return handleReadError(d, "DistributedVlanConnection", id, err)
	}

	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("revision", obj.Revision)
	d.Set("path", obj.Path)

	elem := reflect.ValueOf(&obj).Elem()
	return metadata.StructToSchema(elem, d, distributedVlanConnectionSchema, "", nil)
}

func resourceNsxtPolicyDistributedVlanConnectionUpdate(d *schema.ResourceData, m interface{}) error {

	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining DistributedVlanConnection ID")
	}

	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getPolicyTagsFromSchema(d)

	revision := int64(d.Get("revision").(int))

	obj := model.DistributedVlanConnection{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
		Revision:    &revision,
	}
	subnetExclusiveConfig := model.SubnetExclusiveConfig{}
	subnetExclusiveConfigArgUnparsed := d.Get("subnet_exclusive_config").([]interface{})
	if len(subnetExclusiveConfigArgUnparsed) > 0 || subnetExclusiveConfigArgUnparsed[0] != nil {
		tempSubnetExclusiveConfig, mapExists := subnetExclusiveConfigArgUnparsed[0].(map[string]interface{})
		if ipBlockPath, ok := tempSubnetExclusiveConfig["ip_block_path"]; mapExists && ok {
			ipBlockPathStr := ipBlockPath.(string)
			subnetExclusiveConfig.IpBlockPath = &ipBlockPathStr
		}
		if vlanExtensionSet, ok := tempSubnetExclusiveConfig["vlan_extension"]; ok && len(vlanExtensionSet.([]interface{})) > 0 {
			vlanExtension := vlanExtensionSet.([]interface{})[0]
			if vpcGatewayConnectionEnableMap, ok := vlanExtension.(map[string]interface{}); ok {
				vpcGatewayConnectionEnable, BoolVarExists := vpcGatewayConnectionEnableMap["vpc_gateway_connection_enable"]
				if vpcGatewayConnectionEnablePtr, ok := vpcGatewayConnectionEnable.(bool); BoolVarExists && ok {
					subnetExclusiveConfig.VlanExtension = &model.VlanExtension{
						VpcGatewayConnectionEnable: &vpcGatewayConnectionEnablePtr,
					}
				}
			}
		}
		obj.SubnetExclusiveConfig = &subnetExclusiveConfig
	}
	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, distributedVlanConnectionSchema, "", nil); err != nil {
		return err
	}
	client := clientLayer.NewDistributedVlanConnectionsClient(connector)
	_, err := client.Update(id, obj)
	if err != nil {
		return handleUpdateError("DistributedVlanConnection", id, err)
	}

	return resourceNsxtPolicyDistributedVlanConnectionRead(d, m)
}

func resourceNsxtPolicyDistributedVlanConnectionDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining DistributedVlanConnection ID")
	}

	connector := getPolicyConnector(m)

	client := clientLayer.NewDistributedVlanConnectionsClient(connector)
	err := client.Delete(id)

	if err != nil {
		return handleDeleteError("DistributedVlanConnection", id, err)
	}

	return nil
}
