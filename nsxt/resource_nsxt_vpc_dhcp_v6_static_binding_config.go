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
	"github.com/vmware/terraform-provider-nsxt/api/orgs/projects/vpcs/subnets"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/metadata"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var cliVpcSubnetDhcpV6StaticBindingConfigsClient = subnets.NewDhcpStaticBindingConfigsClient

var dhcpV6StaticBindingConfigSchema = map[string]*metadata.ExtendedSchema{
	"nsx_id":       metadata.GetExtendedSchema(getNsxIDSchema()),
	"path":         metadata.GetExtendedSchema(getPathSchema()),
	"display_name": metadata.GetExtendedSchema(getDisplayNameSchema()),
	"description":  metadata.GetExtendedSchema(getDescriptionSchema()),
	"revision":     metadata.GetExtendedSchema(getRevisionSchema()),
	"tag":          metadata.GetExtendedSchema(getTagsSchema()),
	"parent_path":  metadata.GetExtendedSchema(getPolicyPathSchema(true, true, "Policy path of the parent")),
	"dns_nameservers": {
		Schema: schema.Schema{
			Type:        schema.TypeList,
			Description: "DNS nameservers for the DHCPv6 client.",
			Optional:    true,
			Elem: &metadata.ExtendedSchema{
				Schema: schema.Schema{
					Type: schema.TypeString,
				},
				Metadata: metadata.Metadata{
					SchemaType: "string",
				},
			},
		},
		Metadata: metadata.Metadata{
			SchemaType:   "list",
			SdkFieldName: "DnsNameservers",
		},
	},
	"domain_names": {
		Schema: schema.Schema{
			Type:        schema.TypeList,
			Description: "Domain names assigned to the client host.",
			Optional:    true,
			Elem: &metadata.ExtendedSchema{
				Schema: schema.Schema{
					Type: schema.TypeString,
				},
				Metadata: metadata.Metadata{
					SchemaType: "string",
				},
			},
		},
		Metadata: metadata.Metadata{
			SchemaType:   "list",
			SdkFieldName: "DomainNames",
		},
	},
	"ip_addresses": {
		Schema: schema.Schema{
			Type:        schema.TypeList,
			Description: "IPv6 addresses assigned to the client.",
			Optional:    true,
			Elem: &metadata.ExtendedSchema{
				Schema: schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsIPv6Address,
				},
				Metadata: metadata.Metadata{
					SchemaType: "string",
				},
			},
		},
		Metadata: metadata.Metadata{
			SchemaType:   "list",
			SdkFieldName: "IpAddresses",
		},
	},
	"lease_time": {
		Schema: *getDhcpLeaseTimeSchema(),
		Metadata: metadata.Metadata{
			SchemaType:   "int",
			SdkFieldName: "LeaseTime",
		},
	},
	"preferred_time": {
		Schema: *getDhcpPreferredTimeSchema(),
		Metadata: metadata.Metadata{
			SchemaType:   "int",
			SdkFieldName: "PreferredTime",
			OmitIfEmpty:  true,
		},
	},
	"mac_address": {
		Schema: schema.Schema{
			Type:         schema.TypeString,
			Description:  "MAC address of the client host.",
			Optional:     true,
			ValidateFunc: validation.IsMACAddress,
		},
		Metadata: metadata.Metadata{
			SchemaType:   "string",
			SdkFieldName: "MacAddress",
		},
	},
	"ntp_servers": {
		Schema: schema.Schema{
			Type:        schema.TypeList,
			Description: "NTP servers as FQDN or IPv6 address.",
			Optional:    true,
			Elem: &metadata.ExtendedSchema{
				Schema: schema.Schema{
					Type: schema.TypeString,
				},
				Metadata: metadata.Metadata{
					SchemaType: "string",
				},
			},
		},
		Metadata: metadata.Metadata{
			SchemaType:   "list",
			SdkFieldName: "NtpServers",
		},
	},
	"sntp_servers": {
		Schema: schema.Schema{
			Type:        schema.TypeList,
			Description: "SNTP server IPv6 addresses.",
			Optional:    true,
			Elem: &metadata.ExtendedSchema{
				Schema: schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsIPv6Address,
				},
				Metadata: metadata.Metadata{
					SchemaType: "string",
				},
			},
		},
		Metadata: metadata.Metadata{
			SchemaType:   "list",
			SdkFieldName: "SntpServers",
		},
	},
}

func resourceNsxtVpcSubnetDhcpV6StaticBindingConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtVpcSubnetDhcpV6StaticBindingConfigCreate,
		Read:   resourceNsxtVpcSubnetDhcpV6StaticBindingConfigRead,
		Update: resourceNsxtVpcSubnetDhcpV6StaticBindingConfigUpdate,
		Delete: resourceNsxtVpcSubnetDhcpV6StaticBindingConfigDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtVersionCheckImporter("9.2.0", "VPC DHCP v6 Static Binding", nsxtParentPathResourceImporter),
		},
		Schema: metadata.GetSchemaFromExtendedSchema(dhcpV6StaticBindingConfigSchema),
	}
}

func resourceNsxtVpcSubnetDhcpV6StaticBindingConfigExists(sessionContext utl.SessionContext, parentPath string, id string, connector client.Connector) (bool, error) {
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 4, vpcSubnetPathExample)
	if pathErr != nil {
		return false, pathErr
	}
	client := cliVpcSubnetDhcpV6StaticBindingConfigsClient(sessionContext, connector)
	_, err := client.Get(parents[0], parents[1], parents[2], parents[3], id)
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving resource", err)
}

func resourceNsxtVpcSubnetDhcpV6StaticBindingConfigCreate(d *schema.ResourceData, m interface{}) error {
	if !util.NsxVersionHigherOrEqual("9.2.0") {
		return fmt.Errorf("VPC Subnet DHCP v6 Static Binding Config resource requires NSX version 9.2.0 or higher")
	}
	connector := getPolicyConnector(m)

	id, err := getOrGenerateIDWithParent(d, m, resourceNsxtVpcSubnetDhcpV6StaticBindingConfigExists)
	if err != nil {
		return err
	}

	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 4, vpcSubnetPathExample)
	if pathErr != nil {
		return pathErr
	}
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)

	obj := model.DhcpV6StaticBindingConfig{
		DisplayName:  &displayName,
		Description:  &description,
		Tags:         tags,
		ResourceType: model.DhcpStaticBindingConfig_RESOURCE_TYPE_DHCPV6STATICBINDINGCONFIG,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, dhcpV6StaticBindingConfigSchema, "", nil); err != nil {
		return err
	}

	log.Printf("[INFO] Creating DhcpV6StaticBindingConfig with ID %s", id)

	converter := bindings.NewTypeConverter()
	convObj, convErrs := converter.ConvertToVapi(obj, model.DhcpV6StaticBindingConfigBindingType())
	if convErrs != nil {
		return convErrs[0]
	}

	sessionContext := getSessionContext(d, m)
	client := cliVpcSubnetDhcpV6StaticBindingConfigsClient(sessionContext, connector)
	err = client.Patch(parents[0], parents[1], parents[2], parents[3], id, convObj.(*data.StructValue))
	if err != nil {
		return handleCreateError("DhcpV6StaticBindingConfig", id, err)
	}
	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtVpcSubnetDhcpV6StaticBindingConfigRead(d, m)
}

func resourceNsxtVpcSubnetDhcpV6StaticBindingConfigRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining DhcpV6StaticBindingConfig ID")
	}

	parentPath := d.Get("parent_path").(string)
	sessionContext := getParentContext(d, m, parentPath)
	client := cliVpcSubnetDhcpV6StaticBindingConfigsClient(sessionContext, connector)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 4, vpcSubnetPathExample)
	if pathErr != nil {
		return pathErr
	}
	dhcpObj, err := client.Get(parents[0], parents[1], parents[2], parents[3], id)
	if err != nil {
		return handleReadError(d, "DhcpV6StaticBindingConfig", id, err)
	}

	converter := bindings.NewTypeConverter()
	convObj, errs := converter.ConvertToGolang(dhcpObj, model.DhcpV6StaticBindingConfigBindingType())
	if errs != nil {
		return errs[0]
	}
	obj := convObj.(model.DhcpV6StaticBindingConfig)
	if obj.ResourceType != model.DhcpV6StaticBindingConfig__TYPE_IDENTIFIER {
		return handleReadError(d, "DhcpV6 Static Binding Config", id, fmt.Errorf("Unexpected ResourceType %s", obj.ResourceType))
	}

	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("revision", obj.Revision)
	d.Set("path", obj.Path)

	elem := reflect.ValueOf(&obj).Elem()
	return metadata.StructToSchema(elem, d, dhcpV6StaticBindingConfigSchema, "", nil)
}

func resourceNsxtVpcSubnetDhcpV6StaticBindingConfigUpdate(d *schema.ResourceData, m interface{}) error {

	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining DhcpV6StaticBindingConfig ID")
	}

	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 4, vpcSubnetPathExample)
	if pathErr != nil {
		return pathErr
	}
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getPolicyTagsFromSchema(d)

	revision := int64(d.Get("revision").(int))

	obj := model.DhcpV6StaticBindingConfig{
		DisplayName:  &displayName,
		Description:  &description,
		Tags:         tags,
		Revision:     &revision,
		ResourceType: model.DhcpStaticBindingConfig_RESOURCE_TYPE_DHCPV6STATICBINDINGCONFIG,
	}

	elem := reflect.ValueOf(&obj).Elem()
	if err := metadata.SchemaToStruct(elem, d, dhcpV6StaticBindingConfigSchema, "", nil); err != nil {
		return err
	}

	converter := bindings.NewTypeConverter()
	convObj, convErrs := converter.ConvertToVapi(obj, model.DhcpV6StaticBindingConfigBindingType())
	if convErrs != nil {
		return convErrs[0]
	}

	sessionContext := getParentContext(d, m, parentPath)
	client := cliVpcSubnetDhcpV6StaticBindingConfigsClient(sessionContext, connector)
	_, err := client.Update(parents[0], parents[1], parents[2], parents[3], id, convObj.(*data.StructValue))
	if err != nil {
		d.Partial(true)
		return handleUpdateError("DhcpV6StaticBindingConfig", id, err)
	}

	return resourceNsxtVpcSubnetDhcpV6StaticBindingConfigRead(d, m)
}

func resourceNsxtVpcSubnetDhcpV6StaticBindingConfigDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining DhcpV6StaticBindingConfig ID")
	}

	connector := getPolicyConnector(m)
	parentPath := d.Get("parent_path").(string)
	parents, pathErr := parseStandardPolicyPathVerifySize(parentPath, 4, vpcSubnetPathExample)
	if pathErr != nil {
		return pathErr
	}

	sessionContext := getParentContext(d, m, parentPath)
	client := cliVpcSubnetDhcpV6StaticBindingConfigsClient(sessionContext, connector)
	err := client.Delete(parents[0], parents[1], parents[2], parents[3], id)

	if err != nil {
		return handleDeleteError("DhcpV6StaticBindingConfig", id, err)
	}

	return nil
}
