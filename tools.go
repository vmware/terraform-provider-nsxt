package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vmware/go-vmware-nsxt/common"
	"github.com/vmware/go-vmware-nsxt/manager"
)

func Interface2StringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, val)
		}
	}
	return vs
}

func StringList2Interface(list []string) []interface{} {
	vs := make([]interface{}, 0, len(list))
	for _, v := range list {
		vs = append(vs, v)
	}
	return vs
}

func GetRevisionSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeInt,
		Computed: true,
	}
}

// utilities to define & handle tags
func GetTagsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"Scope": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
				"Tag": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	}
}

func GetTagsFromSchema(d *schema.ResourceData) []common.Tag {
	tags := d.Get("Tags").([]interface{})
	var tagList []common.Tag
	for _, tag := range tags {
		data := tag.(map[string]interface{})
		elem := common.Tag{
			Scope: data["Scope"].(string),
			Tag:   data["Tag"].(string)}

		tagList = append(tagList, elem)
	}
	return tagList
}

func SetTagsInSchema(d *schema.ResourceData, tags []common.Tag) {
	var tagList []map[string]string
	for _, tag := range tags {
		elem := make(map[string]string)
		elem["Scope"] = tag.Scope
		elem["Tag"] = tag.Tag
		tagList = append(tagList, elem)
	}
	d.Set("Tags", tagList)
}

// utilities to define & handle switching profiles
func GetSwitchingProfileIdsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"Key": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
				"Value": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	}
}

func GetSwitchingProfileIdsFromSchema(d *schema.ResourceData) []manager.SwitchingProfileTypeIdEntry {
	profiles := d.Get("SwitchingProfileIds").([]interface{})
	var profileList []manager.SwitchingProfileTypeIdEntry
	for _, profile := range profiles {
		data := profile.(map[string]interface{})
		elem := manager.SwitchingProfileTypeIdEntry{
			Key:   data["Key"].(string),
			Value: data["Value"].(string)}

		profileList = append(profileList, elem)
	}
	return profileList
}

func SetSwitchingProfileIdsInSchema(d *schema.ResourceData, profiles []manager.SwitchingProfileTypeIdEntry) {
	var profileList []map[string]string
	for _, profile := range profiles {
		elem := make(map[string]string)
		elem["Key"] = profile.Key
		elem["Value"] = profile.Value
		profileList = append(profileList, elem)
	}
	d.Set("SwitchingProfileIds", profileList)
}

// utilities to define & handle address bindings
func GetAddressBindingsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"IpAddress": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
				},
				"MacAddress": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
				},
				"Vlan": &schema.Schema{
					Type:     schema.TypeInt,
					Optional: true,
				},
			},
		},
	}
}

func GetAddressBindingsFromSchema(d *schema.ResourceData) []manager.PacketAddressClassifier {
	bindings := d.Get("AddressBindings").([]interface{})
	var bindingList []manager.PacketAddressClassifier
	for _, binding := range bindings {
		data := binding.(map[string]interface{})
		elem := manager.PacketAddressClassifier{
			IpAddress:  data["IpAddress"].(string),
			MacAddress: data["MacAddress"].(string),
			Vlan:       data["Vlan"].(int64),
		}

		bindingList = append(bindingList, elem)
	}
	return bindingList
}

func SetAddressBindingsInSchema(d *schema.ResourceData, bindings []manager.PacketAddressClassifier) {
	var bindingList []map[string]interface{}
	for _, binding := range bindings {
		elem := make(map[string]interface{})
		elem["IpAddress"] = binding.IpAddress
		elem["MacAddress"] = binding.MacAddress
		elem["Vlan"] = binding.Vlan
		bindingList = append(bindingList, elem)
	}
	d.Set("AddressBindings", bindingList)
}
