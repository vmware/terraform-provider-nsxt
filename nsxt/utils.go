/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vmware/go-vmware-nsxt"
	"github.com/vmware/go-vmware-nsxt/common"
	"github.com/vmware/go-vmware-nsxt/manager"
)

func interface2StringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, val)
		}
	}
	return vs
}

// TODO: ???
func stringList2Interface(list []string) []interface{} {
	vs := make([]interface{}, 0, len(list))
	for _, v := range list {
		vs = append(vs, v)
	}
	return vs
}

func getStringListFromSchemaSet(d *schema.ResourceData, schemaAttrName string) []string {
	return interface2StringList(d.Get(schemaAttrName).(*schema.Set).List())
}

func getRevisionSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeInt,
		Computed: true,
	}
}

func getSystemOwnedSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Indicates system owned resource",
		Computed:    true,
	}
}

// utilities to define & handle tags
func getTagsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"scope": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
				"tag": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	}
}

func getTagsFromSchema(d *schema.ResourceData) []common.Tag {
	tags := d.Get("tags").(*schema.Set).List()
	var tagList []common.Tag
	for _, tag := range tags {
		data := tag.(map[string]interface{})
		elem := common.Tag{
			Scope: data["scope"].(string),
			Tag:   data["tag"].(string)}

		tagList = append(tagList, elem)
	}
	return tagList
}

func setTagsInSchema(d *schema.ResourceData, tags []common.Tag) {
	var tagList []map[string]string
	for _, tag := range tags {
		elem := make(map[string]string)
		elem["scope"] = tag.Scope
		elem["tag"] = tag.Tag
		tagList = append(tagList, elem)
	}
	d.Set("tags", tagList)
}

// utilities to define & handle switching profiles
func getSwitchingProfileIdsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"key": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
				"value": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	}
}

func getSwitchingProfileIdsFromSchema(d *schema.ResourceData) []manager.SwitchingProfileTypeIdEntry {
	profiles := d.Get("switching_profile_ids").(*schema.Set).List()
	var profileList []manager.SwitchingProfileTypeIdEntry
	for _, profile := range profiles {
		data := profile.(map[string]interface{})
		elem := manager.SwitchingProfileTypeIdEntry{
			Key:   data["key"].(string),
			Value: data["value"].(string)}

		profileList = append(profileList, elem)
	}
	return profileList
}

func setSwitchingProfileIdsInSchema(d *schema.ResourceData, nsxClient *nsxt.APIClient, profiles []manager.SwitchingProfileTypeIdEntry) {
	var profileList []map[string]string
	for _, profile := range profiles {
		// ignore system owned profiles
		obj, _, _ := nsxClient.LogicalSwitchingApi.GetSwitchingProfile(nsxClient.Context, profile.Value)
		if obj.SystemOwned {
			continue
		}

		elem := make(map[string]string)
		elem["key"] = profile.Key
		elem["value"] = profile.Value
		profileList = append(profileList, elem)
	}
	d.Set("switching_profile_ids", profileList)
}

// utilities to define & handle address bindings
func getAddressBindingsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"ip_address": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
				},
				"mac_address": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
				},
				"vlan": &schema.Schema{
					Type:     schema.TypeInt,
					Optional: true,
				},
			},
		},
	}
}

func getAddressBindingsFromSchema(d *schema.ResourceData) []manager.PacketAddressClassifier {
	bindings := d.Get("address_bindings").(*schema.Set).List()
	var bindingList []manager.PacketAddressClassifier
	for _, binding := range bindings {
		data := binding.(map[string]interface{})
		elem := manager.PacketAddressClassifier{
			IpAddress:  data["ip_address"].(string),
			MacAddress: data["mac_address"].(string),
			Vlan:       data["vlan"].(int64),
		}

		bindingList = append(bindingList, elem)
	}
	return bindingList
}

func setAddressBindingsInSchema(d *schema.ResourceData, bindings []manager.PacketAddressClassifier) {
	var bindingList []map[string]interface{}
	for _, binding := range bindings {
		elem := make(map[string]interface{})
		elem["ip_address"] = binding.IpAddress
		elem["mac_address"] = binding.MacAddress
		elem["vlan"] = binding.Vlan
		bindingList = append(bindingList, elem)
	}
	d.Set("address_bindings", bindingList)
}

func getResourceReferencesSchema(required bool, computed bool) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: required,
		Optional: !required,
		Computed: computed,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"is_valid": &schema.Schema{
					Type:     schema.TypeBool,
					Computed: true,
				},
				"target_display_name": &schema.Schema{
					Type:     schema.TypeString,
					Computed: true,
				},
				"target_id": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
				},
				"target_type": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}
}

func getResourceReferencesFromSchema(d *schema.ResourceData, schemaAttrName string) []common.ResourceReference {
	references := d.Get(schemaAttrName).([]interface{})
	var referenceList []common.ResourceReference
	for _, reference := range references {
		data := reference.(map[string]interface{})
		elem := common.ResourceReference{
			IsValid:           data["is_valid"].(bool),
			TargetDisplayName: data["target_display_name"].(string),
			TargetId:          data["target_id"].(string),
			TargetType:        data["target_type"].(string),
		}

		referenceList = append(referenceList, elem)
	}
	return referenceList
}

func setResourceReferencesInSchema(d *schema.ResourceData, references []common.ResourceReference, schemaAttrName string) {
	var referenceList []map[string]interface{}
	for _, reference := range references {
		elem := make(map[string]interface{})
		elem["is_valid"] = reference.IsValid
		elem["target_display_name"] = reference.TargetDisplayName
		elem["target_id"] = reference.TargetId
		elem["target_type"] = reference.TargetType
		referenceList = append(referenceList, elem)
	}
	d.Set(schemaAttrName, referenceList)
}

func getServiceBindingsFromSchema(d *schema.ResourceData, schemaAttrName string) []manager.ServiceBinding {
	references := d.Get(schemaAttrName).([]interface{})
	var bindingList []manager.ServiceBinding
	for _, reference := range references {
		data := reference.(map[string]interface{})
		ref := common.ResourceReference{
			IsValid:           data["is_valid"].(bool),
			TargetDisplayName: data["target_display_name"].(string),
			TargetId:          data["target_id"].(string),
			TargetType:        data["target_type"].(string),
		}
		elem := manager.ServiceBinding{&ref}
		bindingList = append(bindingList, elem)
	}
	return bindingList
}

func setServiceBindingsInSchema(d *schema.ResourceData, serviceBindings []manager.ServiceBinding, schemaAttrName string) {
	var referenceList []map[string]interface{}
	for _, binding := range serviceBindings {
		elem := make(map[string]interface{})
		elem["is_valid"] = binding.ServiceId.IsValid
		elem["target_display_name"] = binding.ServiceId.TargetDisplayName
		elem["target_id"] = binding.ServiceId.TargetId
		elem["target_type"] = binding.ServiceId.TargetType
		referenceList = append(referenceList, elem)
	}
	d.Set(schemaAttrName, referenceList)
}

func getIpSubnetsSchema(required bool, computed bool) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: !required,
		Required: required,
		Computed: computed,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"ip_addresses": &schema.Schema{
					Type:     schema.TypeList,
					Optional: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"prefix_length": &schema.Schema{
					Type:     schema.TypeInt,
					Optional: true,
				},
			},
		},
	}
}

func getIpSubnetsFromSchema(d *schema.ResourceData) []manager.IpSubnet {
	subnets := d.Get("subnets").([]interface{})
	var subnetList []manager.IpSubnet
	for _, subnet := range subnets {
		data := subnet.(map[string]interface{})
		elem := manager.IpSubnet{
			IpAddresses:  interface2StringList(data["ip_addresses"].([]interface{})),
			PrefixLength: int64(data["prefix_length"].(int)),
		}

		subnetList = append(subnetList, elem)
	}
	return subnetList
}

func setIpSubnetsInSchema(d *schema.ResourceData, subnets []manager.IpSubnet) {
	var subnetList []map[string]interface{}
	for _, subnet := range subnets {
		elem := make(map[string]interface{})
		elem["ip_addresses"] = stringList2Interface(subnet.IpAddresses)
		elem["prefix_length"] = subnet.PrefixLength
		subnetList = append(subnetList, elem)
	}
	d.Set("subnets", subnetList)
}

func makeResourceReference(resourceType string, resourceId string) *common.ResourceReference {
	return &common.ResourceReference{
		TargetType: resourceType,
		TargetId:   resourceId,
	}
}
