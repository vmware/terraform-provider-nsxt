/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"bytes"
	"fmt"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/vmware/go-vmware-nsxt"
	api "github.com/vmware/go-vmware-nsxt"
	"github.com/vmware/go-vmware-nsxt/common"
	"github.com/vmware/go-vmware-nsxt/manager"
	"log"
	"net/http"
)

var adminStateValues = []string{"UP", "DOWN"}
var nsxVersion = ""

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

func stringList2Interface(list []string) []interface{} {
	vs := make([]interface{}, 0, len(list))
	for _, v := range list {
		vs = append(vs, v)
	}
	return vs
}

func interface2Int32List(configured []interface{}) []int32 {
	vs := make([]int32, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(int)
		if ok {
			vs = append(vs, int32(val))
		}
	}
	return vs
}

func int32List2Interface(list []int32) []interface{} {
	vs := make([]interface{}, 0, len(list))
	for _, v := range list {
		vs = append(vs, int(v))
	}
	return vs
}

func getStringListFromSchemaSet(d *schema.ResourceData, schemaAttrName string) []string {
	return interface2StringList(d.Get(schemaAttrName).(*schema.Set).List())
}

func getStringListFromSchemaList(d *schema.ResourceData, schemaAttrName string) []string {
	return interface2StringList(d.Get(schemaAttrName).([]interface{}))
}

func intList2int64List(configured []interface{}) []int64 {
	vs := make([]int64, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(int)
		if ok {
			vs = append(vs, int64(val))
		}
	}
	return vs
}

func getRevisionSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeInt,
		Description: "The _revision property describes the current revision of the resource. To prevent clients from overwriting each other's changes, PUT operations must include the current _revision of the resource, which clients should obtain by issuing a GET operation. If the _revision provided in a PUT request is missing or stale, the operation will be rejected",
		Computed:    true,
	}
}

// utilities to define & handle tags
func getTagsSchemaInternal(required bool, forceNew bool) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "Set of opaque identifiers meaningful to the user",
		Optional:    !required,
		Required:    required,
		ForceNew:    forceNew,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"scope": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: forceNew,
				},
				"tag": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: forceNew,
				},
			},
		},
	}
}

// utilities to define & handle tags
func getTagsSchema() *schema.Schema {
	return getTagsSchemaInternal(false, false)
}

func getRequiredTagsSchema() *schema.Schema {
	return getTagsSchemaInternal(true, false)
}

func getTagsSchemaForceNew() *schema.Schema {
	return getTagsSchemaInternal(false, true)
}

func getCustomizedTagsFromSchema(d *schema.ResourceData, schemaName string) []common.Tag {
	tags := d.Get(schemaName).(*schema.Set).List()
	tagList := make([]common.Tag, 0)
	for _, tag := range tags {
		data := tag.(map[string]interface{})
		elem := common.Tag{
			Scope: data["scope"].(string),
			Tag:   data["tag"].(string)}

		tagList = append(tagList, elem)
	}
	return tagList
}

func setCustomizedTagsInSchema(d *schema.ResourceData, tags []common.Tag, schemaName string) error {
	var tagList []map[string]string
	for _, tag := range tags {
		elem := make(map[string]string)
		elem["scope"] = tag.Scope
		elem["tag"] = tag.Tag
		tagList = append(tagList, elem)
	}
	err := d.Set(schemaName, tagList)
	return err
}
func getTagsFromSchema(d *schema.ResourceData) []common.Tag {
	return getCustomizedTagsFromSchema(d, "tag")
}

func setTagsInSchema(d *schema.ResourceData, tags []common.Tag) error {
	return setCustomizedTagsInSchema(d, tags, "tag")
}

// utilities to define & handle switching profiles
func getSwitchingProfileIdsSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "List of IDs of switching profiles (of various types) to be associated with this object. Default switching profiles will be used if not specified",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"key": {
					Type:        schema.TypeString,
					Description: "The resource type of this profile",
					Required:    true,
				},
				"value": {
					Type:        schema.TypeString,
					Description: "The ID of this profile",
					Required:    true,
				},
			},
		},
	}
}

func getSwitchingProfileIdsFromSchema(d *schema.ResourceData) []manager.SwitchingProfileTypeIdEntry {
	profiles := d.Get("switching_profile_id").(*schema.Set).List()
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

func setSwitchingProfileIdsInSchema(d *schema.ResourceData, nsxClient *nsxt.APIClient, profiles []manager.SwitchingProfileTypeIdEntry) error {
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
	err := d.Set("switching_profile_id", profileList)
	return err
}

// utilities to define & handle address bindings
func getAddressBindingsSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "Address bindings for the Logical switch",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"ip_address": {
					Type:         schema.TypeString,
					Description:  "A single IP address or a subnet cidr",
					Optional:     true,
					ValidateFunc: validateSingleIP(),
				},
				"mac_address": {
					Type:        schema.TypeString,
					Description: "A single MAC address",
					Optional:    true,
				},
				"vlan": {
					Type:        schema.TypeInt,
					Description: "A single vlan tag value",
					Optional:    true,
				},
			},
		},
	}
}

func getAddressBindingsFromSchema(d *schema.ResourceData) []manager.PacketAddressClassifier {
	bindings := d.Get("address_binding").(*schema.Set).List()
	var bindingList []manager.PacketAddressClassifier
	for _, binding := range bindings {
		data := binding.(map[string]interface{})
		elem := manager.PacketAddressClassifier{
			IpAddress:  data["ip_address"].(string),
			MacAddress: data["mac_address"].(string),
			Vlan:       int64(data["vlan"].(int)),
		}

		bindingList = append(bindingList, elem)
	}
	return bindingList
}

func setAddressBindingsInSchema(d *schema.ResourceData, bindings []manager.PacketAddressClassifier) error {
	var bindingList []map[string]interface{}
	for _, binding := range bindings {
		elem := make(map[string]interface{})
		elem["ip_address"] = binding.IpAddress
		elem["mac_address"] = binding.MacAddress
		elem["vlan"] = binding.Vlan
		bindingList = append(bindingList, elem)
	}
	err := d.Set("address_binding", bindingList)
	return err
}

func getResourceReferencesSchema(required bool, computed bool, validTargetTypes []string, description string) *schema.Schema {
	return getResourceReferencesSchemaByType(required, computed, validTargetTypes, true, description, 0)
}

func getSingleResourceReferencesSchema(required bool, computed bool, validTargetTypes []string, description string) *schema.Schema {
	return getResourceReferencesSchemaByType(required, computed, validTargetTypes, true, description, 1)
}

func getResourceReferencesSetSchema(required bool, computed bool, validTargetTypes []string, description string) *schema.Schema {
	return getResourceReferencesSchemaByType(required, computed, validTargetTypes, false, description, 0)
}

func getResourceReferencesSchemaByType(required bool, computed bool, validTargetTypes []string, isList bool, description string, maxItems int) *schema.Schema {
	schType := schema.TypeSet
	if isList {
		schType = schema.TypeList
	}

	return &schema.Schema{
		Type:        schType,
		Required:    required,
		Optional:    !required,
		Computed:    computed,
		MaxItems:    maxItems,
		Description: description,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"is_valid": {
					Type:        schema.TypeBool,
					Description: "A boolean flag which will be set to false if the referenced NSX resource has been deleted",
					Computed:    true,
				},
				"target_display_name": {
					Type:        schema.TypeString,
					Description: "Display name of the NSX resource",
					Computed:    true,
				},
				"target_id": {
					Type:        schema.TypeString,
					Description: "Identifier of the NSX resource",
					Optional:    true,
				},
				"target_type": {
					Type:         schema.TypeString,
					Description:  "Type of the NSX resource",
					Optional:     true,
					ValidateFunc: validation.StringInSlice(validTargetTypes, false),
				},
			},
		},
	}
}

func getSingleResourceReference(references []interface{}) *common.ResourceReference {
	for _, reference := range references {
		data := reference.(map[string]interface{})
		// only 1 ref is allowed so return the first 1
		elem := common.ResourceReference{
			IsValid:           data["is_valid"].(bool),
			TargetDisplayName: data["target_display_name"].(string),
			TargetId:          data["target_id"].(string),
			TargetType:        data["target_type"].(string),
		}
		return &elem
	}
	return nil
}

func getResourceReferences(references []interface{}) []common.ResourceReference {
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

func getResourceReferencesFromSchemaSet(d *schema.ResourceData, schemaAttrName string) []common.ResourceReference {
	references := d.Get(schemaAttrName).(*schema.Set).List()
	return getResourceReferences(references)
}

func returnResourceReferences(references []common.ResourceReference) []map[string]interface{} {
	var referenceList []map[string]interface{}
	for _, reference := range references {
		elem := make(map[string]interface{})
		elem["is_valid"] = reference.IsValid
		elem["target_display_name"] = reference.TargetDisplayName
		elem["target_id"] = reference.TargetId
		elem["target_type"] = reference.TargetType
		referenceList = append(referenceList, elem)
	}
	return referenceList
}

func resourceReferenceHash(v interface{}) int {
	var buf bytes.Buffer

	if v != nil {
		m := v.(map[string]interface{})
		buf.WriteString(fmt.Sprintf("%s-%s", m["target_type"], m["target_id"]))
	}
	return hashcode.String(buf.String())
}

func returnResourceReferencesSet(references []common.ResourceReference) *schema.Set {
	var referenceList []interface{}
	for _, reference := range references {
		elem := make(map[string]interface{})
		elem["is_valid"] = reference.IsValid
		elem["target_display_name"] = reference.TargetDisplayName
		elem["target_id"] = reference.TargetId
		elem["target_type"] = reference.TargetType
		referenceList = append(referenceList, elem)
	}

	s := schema.NewSet(resourceReferenceHash, referenceList)
	return s
}

func setResourceReferencesInSchema(d *schema.ResourceData, references []common.ResourceReference, schemaAttrName string) error {
	referenceList := returnResourceReferences(references)
	err := d.Set(schemaAttrName, referenceList)
	return err
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
		elem := manager.ServiceBinding{ServiceId: &ref}
		bindingList = append(bindingList, elem)
	}
	return bindingList
}

func setServiceBindingsInSchema(d *schema.ResourceData, serviceBindings []manager.ServiceBinding, schemaAttrName string) error {
	var referenceList []map[string]interface{}
	for _, binding := range serviceBindings {
		elem := make(map[string]interface{})
		elem["is_valid"] = binding.ServiceId.IsValid
		elem["target_display_name"] = binding.ServiceId.TargetDisplayName
		elem["target_id"] = binding.ServiceId.TargetId
		elem["target_type"] = binding.ServiceId.TargetType
		referenceList = append(referenceList, elem)
	}
	err := d.Set(schemaAttrName, referenceList)
	return err
}

func getAdminStateSchema() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "Represents Desired state of the object",
		Default:      "UP",
		ValidateFunc: validation.StringInSlice(adminStateValues, false),
	}
}

func getIDSetSchema(description string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: description,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
	}
}

func getIPRangesSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "List of IP Ranges",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"start": {
					Type:         schema.TypeString,
					ValidateFunc: validateSingleIP(),
					Required:     true,
				},
				"end": {
					Type:         schema.TypeString,
					ValidateFunc: validateSingleIP(),
					Required:     true,
				},
			},
		},
	}
}

func getIPRangesFromSchema(d *schema.ResourceData) []manager.IpPoolRange {
	ranges := d.Get("ip_range").([]interface{})
	var rangeList []manager.IpPoolRange
	for _, r := range ranges {
		data := r.(map[string]interface{})
		elem := manager.IpPoolRange{
			Start: data["start"].(string),
			End:   data["end"].(string)}

		rangeList = append(rangeList, elem)
	}
	return rangeList
}

func setIPRangesInSchema(d *schema.ResourceData, ranges []manager.IpPoolRange) error {
	var rangeList []map[string]string
	for _, r := range ranges {
		elem := make(map[string]string)
		elem["start"] = r.Start
		elem["end"] = r.End
		rangeList = append(rangeList, elem)
	}
	return d.Set("ip_range", rangeList)
}

func makeResourceReference(resourceType string, resourceID string) *common.ResourceReference {
	return &common.ResourceReference{
		TargetType: resourceType,
		TargetId:   resourceID,
	}
}

func getNSXVersion(nsxClient *api.APIClient) (string, error) {
	nodeProperties, resp, err := nsxClient.NsxComponentAdministrationApi.ReadNodeProperties(nsxClient.Context)

	if err != nil || resp.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("Failed to retrieve NSX version (%s). Please check connectivity and authentication settings of the provider", err)

	}
	log.Printf("[DEBUG] NSX version is %s", nodeProperties.NodeVersion)
	return nodeProperties.NodeVersion, nil
}

func initNSXVersion(nsxClient *api.APIClient) error {
	var err error
	nsxVersion, err = getNSXVersion(nsxClient)
	return err
}

func nsxVersionLower(ver string) bool {

	requestedVersion, err1 := version.NewVersion(ver)
	currentVersion, err2 := version.NewVersion(nsxVersion)
	if err1 != nil || err2 != nil {
		log.Printf("[ERROR] Failed perform version check for version %s", ver)
		return true
	}
	return currentVersion.LessThan(requestedVersion)
}

func nsxVersionHigherOrEqual(ver string) bool {

	requestedVersion, err1 := version.NewVersion(ver)
	currentVersion, err2 := version.NewVersion(nsxVersion)
	if err1 != nil || err2 != nil {
		log.Printf("[ERROR] Failed perform version check for version %s", ver)
		return false
	}
	return currentVersion.Compare(requestedVersion) >= 0
}

func resourceNotSupportedError() error {
	return fmt.Errorf("This resource is not supported with given provider settings")
}

func dataSourceNotSupportedError() error {
	return fmt.Errorf("This data source is not supported with given provider settings")
}
