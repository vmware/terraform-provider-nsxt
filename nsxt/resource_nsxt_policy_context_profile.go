/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	gm_infra "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/global_infra"
	gm_cont_prof "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/global_infra/context_profiles"
	gm_model "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	cont_prof "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/context_profiles"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

var attributeKeyMap = map[string]string{
	"app_id":       model.PolicyAttributes_KEY_APP_ID,
	"domain_name":  model.PolicyAttributes_KEY_DOMAIN_NAME,
	"url_category": model.PolicyAttributes_KEY_URL_CATEGORY,
}

var attributeReverseKeyMap = map[string]string{
	model.PolicyAttributes_KEY_APP_ID:       "app_id",
	model.PolicyAttributes_KEY_DOMAIN_NAME:  "domain_name",
	model.PolicyAttributes_KEY_URL_CATEGORY: "url_category",
}

var subAttributeKeyMap = map[string]string{
	"tls_cipher_suite": model.PolicySubAttributes_KEY_TLS_CIPHER_SUITE,
	"tls_version":      model.PolicySubAttributes_KEY_TLS_VERSION,
	"cifs_smb_version": model.PolicySubAttributes_KEY_CIFS_SMB_VERSION,
}

var subAttributeReverseKeyMap = map[string]string{
	model.PolicySubAttributes_KEY_TLS_CIPHER_SUITE: "tls_cipher_suite",
	model.PolicySubAttributes_KEY_TLS_VERSION:      "tls_version",
	model.PolicySubAttributes_KEY_CIFS_SMB_VERSION: "cifs_smb_version",
}

func resourceNsxtPolicyContextProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyContextProfileCreate,
		Read:   resourceNsxtPolicyContextProfileRead,
		Update: resourceNsxtPolicyContextProfileUpdate,
		Delete: resourceNsxtPolicyContextProfileDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"nsx_id":       getNsxIDSchema(),
			"path":         getPathSchema(),
			"display_name": getDisplayNameSchema(),
			"description":  getDescriptionSchema(),
			"revision":     getRevisionSchema(),
			"tag":          getTagsSchema(),
			"app_id":       getContextProfilePolicyAppIDAttributesSchema(),
			"domain_name":  getContextProfilePolicyOtherAttributesSchema(),
			"url_category": getContextProfilePolicyOtherAttributesSchema(),
		},
	}
}

func getContextProfilePolicyAppIDAttributesSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"description": getDescriptionSchema(),
				"is_alg_type": {
					Type:        schema.TypeBool,
					Description: "Whether the app_id value is ALG type or not",
					Computed:    true,
				},
				"value": {
					Type:        schema.TypeSet,
					Description: "Values for attribute key",
					Required:    true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"sub_attribute": getPolicyAttributeSubAttributeSchema(),
			},
		},
	}
}

func getContextProfilePolicyOtherAttributesSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"description": getDescriptionSchema(),
				"value": {
					Type:        schema.TypeSet,
					Description: "Values for attribute key",
					Required:    true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}
}

func getPolicyAttributeSubAttributeSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"tls_cipher_suite": getPolicyAttributeSubAttributeValueSchema("tls_cipher_suite"),
				"tls_version":      getPolicyAttributeSubAttributeValueSchema("tls_version"),
				"cifs_smb_version": getPolicyAttributeSubAttributeValueSchema("cifs_smb_version"),
			},
		},
	}
}
func getPolicyAttributeSubAttributeValueSchema(subAttributeKey string) *schema.Schema {
	description := fmt.Sprintf("Values for sub attribute key %s", subAttributeKey)
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: description,
		Optional:    true,
		MinItems:    1,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
}

func resourceNsxtPolicyContextProfileExists(id string, connector *client.RestConnector, isGlobalManager bool) (bool, error) {
	var err error
	if isGlobalManager {
		client := gm_infra.NewDefaultContextProfilesClient(connector)
		_, err = client.Get(id)
	} else {
		client := infra.NewDefaultContextProfilesClient(connector)
		_, err = client.Get(id)
	}

	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving ContextProfile", err)
}

func resourceNsxtPolicyContextProfileCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	// Initialize resource Id and verify this ID is not yet used
	id, err := getOrGenerateID(d, m, resourceNsxtPolicyContextProfileExists)
	if err != nil {
		return err
	}

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	attributesStructList := make([]model.PolicyAttributes, 0)
	for key := range attributeKeyMap {
		attributes := d.Get(key).(*schema.Set).List()
		if len(attributes) > 0 {
			err = checkAttributesValid(attributes, m, key)
			if err != nil {
				return err
			}
			attributeStructList, err := constructAttributesModelList(attributes, key)
			if err != nil {
				return err
			}
			attributesStructList = append(attributesStructList, attributeStructList...)
		}
	}
	if len(attributesStructList) == 0 {
		return fmt.Errorf("At least one attribute should be set")
	}

	tags := getPolicyTagsFromSchema(d)

	obj := model.PolicyContextProfile{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
		Attributes:  attributesStructList,
	}

	// Create the resource using PATCH
	log.Printf("[INFO] Creating ContextProfile with ID %s", id)
	if isPolicyGlobalManager(m) {
		client := gm_infra.NewDefaultContextProfilesClient(connector)
		gmObj, err1 := convertModelBindingType(obj, model.PolicyContextProfileBindingType(), gm_model.PolicyContextProfileBindingType())
		if err1 != nil {
			return err1
		}
		err = client.Patch(id, gmObj.(gm_model.PolicyContextProfile))
	} else {
		client := infra.NewDefaultContextProfilesClient(connector)
		err = client.Patch(id, obj)
	}
	if err != nil {
		return handleCreateError("ContextProfile", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyContextProfileRead(d, m)
}

func resourceNsxtPolicyContextProfileRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining ContextProfile ID")
	}

	var obj model.PolicyContextProfile
	if isPolicyGlobalManager(m) {
		client := gm_infra.NewDefaultContextProfilesClient(connector)
		gmObj, err := client.Get(id)
		if err != nil {
			return handleReadError(d, "ContextProfile", id, err)
		}
		rawObj, err := convertModelBindingType(gmObj, gm_model.PolicyContextProfileBindingType(), model.PolicyContextProfileBindingType())
		if err != nil {
			return err
		}
		obj = rawObj.(model.PolicyContextProfile)
	} else {
		var err error
		client := infra.NewDefaultContextProfilesClient(connector)
		obj, err = client.Get(id)
		if err != nil {
			return handleReadError(d, "ContextProfile", id, err)
		}
	}

	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("path", obj.Path)
	d.Set("revision", obj.Revision)
	fillAttributesInSchema(d, obj.Attributes)

	return nil
}

func resourceNsxtPolicyContextProfileUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining ContextProfile ID")
	}

	// Read the rest of the configured parameters
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	attributesStructList := make([]model.PolicyAttributes, 0)
	var err error
	for key := range attributeKeyMap {
		attributes := d.Get(key).(*schema.Set).List()
		err := checkAttributesValid(attributes, m, key)
		if err != nil {
			return err
		}
		attributeStructList, err := constructAttributesModelList(attributes, key)
		if err != nil {
			return err
		}
		attributesStructList = append(attributesStructList, attributeStructList...)
	}
	tags := getPolicyTagsFromSchema(d)

	obj := model.PolicyContextProfile{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
		Attributes:  attributesStructList,
	}

	// Update the resource using PATCH
	if isPolicyGlobalManager(m) {
		rawObj, err1 := convertModelBindingType(obj, model.PolicyContextProfileBindingType(), gm_model.PolicyContextProfileBindingType())
		if err1 != nil {
			return err1
		}
		gmObj := rawObj.(gm_model.PolicyContextProfile)
		client := gm_infra.NewDefaultContextProfilesClient(connector)
		err = client.Patch(id, gmObj)
	} else {
		client := infra.NewDefaultContextProfilesClient(connector)
		err = client.Patch(id, obj)
	}

	if err != nil {
		return handleUpdateError("ContextProfile", id, err)
	}

	return resourceNsxtPolicyContextProfileRead(d, m)
}

func resourceNsxtPolicyContextProfileDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining ContextProfile ID")
	}

	connector := getPolicyConnector(m)
	var err error
	force := true
	if isPolicyGlobalManager(m) {
		client := gm_infra.NewDefaultContextProfilesClient(connector)
		err = client.Delete(id, &force)
	} else {
		client := infra.NewDefaultContextProfilesClient(connector)
		err = client.Delete(id, &force)
	}
	if err != nil {
		return handleDeleteError("ContextProfile", id, err)
	}

	return nil
}

func checkAttributesValid(attributes []interface{}, m interface{}, key string) error {

	if key == "app_id" {
		err := validateSubAttributes(attributes)
		if err != nil {
			return err
		}
	}
	var attrClient interface{}
	connector := getPolicyConnector(m)
	isPolicyGlobalManager := isPolicyGlobalManager(m)
	if isPolicyGlobalManager {
		attrClient = gm_cont_prof.NewDefaultAttributesClient(connector)
	} else {
		attrClient = cont_prof.NewDefaultAttributesClient(connector)
	}
	attributeValues, err := listAttributesWithKey(attributeKeyMap[key], attrClient, isPolicyGlobalManager)
	if err != nil {
		return err
	}
	for _, attribute := range attributes {
		attributeMap := attribute.(map[string]interface{})
		values := interface2StringList(attributeMap["value"].(*schema.Set).List())
		if !containsElements(values, attributeValues) {
			err := fmt.Errorf("Attribute values %s are not valid for attribute type %s", values, key)
			return err
		}
	}
	return nil
}

func validateSubAttributes(attributes []interface{}) error {
	// Validates that sub_attribute keys only present in an attribute with one value
	for _, attribute := range attributes {
		attributeMap := attribute.(map[string]interface{})
		values := attributeMap["value"].(*schema.Set).List()
		if len(attributeMap["sub_attribute"].(*schema.Set).List()) > 0 && len(values) > 1 {
			err := fmt.Errorf("Multiple values found for attribute. Sub-attributes are only applicable to an attribute with a single value")
			return err
		}
	}
	return nil
}

func listAttributesWithKey(attributeKey string, attributeClient interface{}, isPolicyGlobalManager bool) ([]string, error) {
	// returns a list of attribute values
	policyAttributes := make([]string, 0)
	policyContextProfileListResult, err := listContextProfileWithKey(&attributeKey, attributeClient, isPolicyGlobalManager)
	if err != nil {
		return policyAttributes, err
	}
	for _, policyContextProfile := range policyContextProfileListResult.Results {
		for _, attribute := range policyContextProfile.Attributes {
			policyAttributes = append(policyAttributes, attribute.Value...)
		}
	}
	return policyAttributes, nil
}

func listContextProfileWithKey(attributeKey *string, attributeClient interface{}, isPolicyGlobalManager bool) (model.PolicyContextProfileListResult, error) {
	var policyContextProfileListResult model.PolicyContextProfileListResult
	includeMarkForDeleteObjectsParam := false
	if isPolicyGlobalManager {
		client := attributeClient.(*gm_cont_prof.DefaultAttributesClient)
		gmPolicyContextProfileListResult, err := client.List(attributeKey, nil, &includeMarkForDeleteObjectsParam, nil, nil, nil, nil)
		if err != nil {
			return policyContextProfileListResult, err
		}
		lmPolicyContextProfileListResult, err := convertModelBindingType(gmPolicyContextProfileListResult, gm_model.PolicyContextProfileListResultBindingType(), model.PolicyContextProfileListResultBindingType())
		if err != nil {
			return policyContextProfileListResult, err
		}
		policyContextProfileListResult = lmPolicyContextProfileListResult.(model.PolicyContextProfileListResult)
		return policyContextProfileListResult, err
	}
	var err error
	client := attributeClient.(*cont_prof.DefaultAttributesClient)
	policyContextProfileListResult, err = client.List(attributeKey, nil, &includeMarkForDeleteObjectsParam, nil, nil, nil, nil)
	return policyContextProfileListResult, err
}

func constructAttributesModelList(rawAttributes []interface{}, key string) ([]model.PolicyAttributes, error) {
	res := make([]model.PolicyAttributes, 0, len(rawAttributes))
	for _, rawAttribute := range rawAttributes {
		attributeMap := rawAttribute.(map[string]interface{})
		dataType := model.PolicyAttributes_DATATYPE_STRING
		description := attributeMap["description"].(string)
		attrKey := attributeKeyMap[key]
		values := interface2StringList(attributeMap["value"].(*schema.Set).List())
		subAttributesList := make([]model.PolicySubAttributes, 0)
		if key == "app_id" {
			var err error
			subAttributes := attributeMap["sub_attribute"].(*schema.Set).List()
			subAttributesList, err = constructSubAttributeModelList(subAttributes)
			if err != nil {
				return nil, err
			}
		}
		attributeStruct := model.PolicyAttributes{
			Datatype:      &dataType,
			Description:   &description,
			Key:           &attrKey,
			Value:         values,
			SubAttributes: subAttributesList,
		}
		res = append(res, attributeStruct)
	}
	return res, nil
}

func constructSubAttributeModelList(rawSubAttributes []interface{}) ([]model.PolicySubAttributes, error) {
	res := make([]model.PolicySubAttributes, 0)
	dataType := model.PolicySubAttributes_DATATYPE_STRING
	for _, rawSubAttribute := range rawSubAttributes {
		rawSubAttributeMap := rawSubAttribute.(map[string]interface{})
		for key, subAttrKey := range subAttributeKeyMap {
			vals := rawSubAttributeMap[key]
			if vals != nil {
				values := interface2StringList(vals.(*schema.Set).List())
				if len(values) > 0 {
					tmp := subAttrKey
					subAttributeStruct := model.PolicySubAttributes{
						Datatype: &dataType,
						Key:      &tmp,
						Value:    values,
					}
					res = append(res, subAttributeStruct)
				}
			}
		}
	}
	return res, nil
}

func fillAttributesInSchema(d *schema.ResourceData, policyAttributes []model.PolicyAttributes) {
	attributes := make(map[string][]interface{})
	for _, policyAttribute := range policyAttributes {
		elem := make(map[string]interface{})
		key := attributeReverseKeyMap[*policyAttribute.Key]
		elem["description"] = policyAttribute.Description
		elem["value"] = policyAttribute.Value
		if *policyAttribute.Key == model.PolicyAttributes_KEY_APP_ID {
			if len(policyAttribute.SubAttributes) > 0 {
				elem["sub_attribute"] = fillSubAttributesInSchema(policyAttribute.SubAttributes)
			}
			elem["is_alg_type"] = policyAttribute.IsALGType
		}
		attributes[key] = append(attributes[key], elem)
	}
	for key, attributeList := range attributes {
		d.Set(key, attributeList)
	}
}

func fillSubAttributesInSchema(policySubAttributes []model.PolicySubAttributes) []interface{} {
	subAttributes := make(map[string]interface{})
	for _, policySubAttribute := range policySubAttributes {
		key := subAttributeReverseKeyMap[*policySubAttribute.Key]
		subAttributes[key] = policySubAttribute.Value
	}
	res := make([]interface{}, 0, 1)
	res = append(res, subAttributes)
	return res
}
