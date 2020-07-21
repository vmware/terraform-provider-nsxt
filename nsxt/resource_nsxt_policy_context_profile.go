/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	gm_infra "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/global_infra"
	gm_cont_prof "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/global_infra/context_profiles"
	gm_model "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	cont_prof "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/context_profiles"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"log"
)

var attributeDataTypes = []string{
	model.PolicyAttributes_DATATYPE_STRING,
}

var attributeKeys = []string{
	model.PolicyAttributes_KEY_APP_ID,
	model.PolicyAttributes_KEY_DOMAIN_NAME,
	model.PolicyAttributes_KEY_URL_CATEGORY,
}

var subAttributeDataTypes = []string{
	model.PolicySubAttributes_DATATYPE_STRING,
}

var subAttributeKeys = []string{
	model.PolicySubAttributes_KEY_CIFS_SMB_VERSION,
	model.PolicySubAttributes_KEY_TLS_CIPHER_SUITE,
	model.PolicySubAttributes_KEY_TLS_VERSION,
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
			"attribute":    getContextProfilePolicyAttributesSchema(),
		},
	}
}

func resourceNsxtPolicyContextProfileExists(id string, connector *client.RestConnector, isGlobalManager bool) bool {
	var err error
	if isGlobalManager {
		client := gm_infra.NewDefaultContextProfilesClient(connector)
		_, err = client.Get(id)
	} else {
		client := infra.NewDefaultContextProfilesClient(connector)
		_, err = client.Get(id)
	}

	if err == nil {
		return true
	}

	if isNotFoundError(err) {
		return false
	}

	logAPIError("Error retrieving ContextProfile", err)

	return false
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
	attributes := d.Get("attribute").(*schema.Set).List()
	err = checkAttributesValid(attributes, m)
	if err != nil {
		return err
	}
	attributesStructList, err := constructAttributesModelList(attributes)
	if err != nil {
		return err
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
	d.Set("attribute", fillAttributesInSchema(obj.Attributes))
	return nil
}

func resourceNsxtPolicyContextProfileUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining ContextProfile ID")
	}

	// Read the rest of the configured parameters
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getPolicyTagsFromSchema(d)
	attributes := d.Get("attribute").(*schema.Set).List()
	err := checkAttributesValid(attributes, m)
	if err != nil {
		return err
	}
	var attributesStructList []model.PolicyAttributes
	attributesStructList, err = constructAttributesModelList(attributes)
	if err != nil {
		return err
	}

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

func getContextProfilePolicyAttributesSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Required: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"description": {
					Type:        schema.TypeString,
					Description: "Description for attribute value",
					Computed:    true,
				},
				"data_type": {
					Type:         schema.TypeString,
					Description:  "Data type of attribute",
					Required:     true,
					ValidateFunc: validation.StringInSlice(attributeDataTypes, false),
				},
				"is_alg_type": {
					Type:        schema.TypeBool,
					Description: "Whether the APP_ID value is ALG type or not",
					Computed:    true,
				},
				"key": {
					Type:         schema.TypeString,
					Description:  "Key for attribute",
					Required:     true,
					ValidateFunc: validation.StringInSlice(attributeKeys, false),
				},
				"value": {
					Type:        schema.TypeSet,
					Description: "Values for attribute key",
					Required:    true,
					MinItems:    1,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"sub_attribute": getPolicyAttributeSubAttributesSchema(),
			},
		},
	}
}

func checkAttributesValid(attributes []interface{}, m interface{}) error {

	var attrClient interface{}
	connector := getPolicyConnector(m)
	isPolicyGlobalManager := isPolicyGlobalManager(m)
	if isPolicyGlobalManager {
		attrClient = gm_cont_prof.NewDefaultAttributesClient(connector)
	} else {
		attrClient = cont_prof.NewDefaultAttributesClient(connector)
	}
	attributeMap, err := validateAndConstructAttributesMap(attributes)
	if err != nil {
		return err
	}

	for key, values := range attributeMap {
		attributeValues, err := listAttributesWithKey(&key, attrClient, isPolicyGlobalManager)
		if err != nil {
			return err
		}
		if len(attributeValues) == 0 {
			// Theoretically impossible as the attributes are pre-set on NSX
			err := fmt.Errorf("No attribute values are available for attribute type %s", key)
			return err
		}
		if !containsElements(values, attributeValues) {
			err := fmt.Errorf("Attribute values %s are not valid for attribute type %s", values, key)
			return err
		}
	}
	return nil
}

func validateAndConstructAttributesMap(attributes []interface{}) (map[string][]string, error) {
	// Validate that attribute keys are unique
	res := make(map[string][]string)
	for _, attribute := range attributes {
		attributeMap := attribute.(map[string]interface{})
		key := attributeMap["key"].(string)
		values := interface2StringList(attributeMap["value"].(*schema.Set).List())
		subAttributes := attributeMap["sub_attribute"].(*schema.Set).List()
		// There should be only one value if sub attributes are specified
		if len(subAttributes) > 0 && len(values) > 1 {
			err := fmt.Errorf("Multiple values found for attribute key %s. Sub-attribtes are only applicable to an attribute with a single value", key)
			return nil, err
		}
		if res[key] != nil {
			err := fmt.Errorf("Duplicate attribute key found: %s", key)
			return nil, err
		}
		res[key] = values
	}
	return res, nil
}

func listAttributesWithKey(attributeKey *string, attributeClient interface{}, isPolicyGlobalManager bool) ([]string, error) {
	// returns a list of attribute values
	policyAttributes := make([]string, 0)
	policyContextProfileListResult, err := listContextProfileWithKey(attributeKey, attributeClient, isPolicyGlobalManager)
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
			logAPIError("Error listing Policy Attributes", err)
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

func constructAttributesModelList(rawAttributes []interface{}) ([]model.PolicyAttributes, error) {
	res := make([]model.PolicyAttributes, 0, len(rawAttributes))
	for _, rawAttribute := range rawAttributes {
		attributeMap := rawAttribute.(map[string]interface{})
		dataType := attributeMap["data_type"].(string)
		key := attributeMap["key"].(string)
		values := interface2StringList(attributeMap["value"].(*schema.Set).List())
		subAttributes := attributeMap["sub_attribute"].(*schema.Set).List()
		subAttributesList, err := constructSubAttributeModelList(subAttributes)
		if err != nil {
			return nil, err
		}
		attributeStruct := model.PolicyAttributes{
			Datatype:      &dataType,
			Key:           &key,
			Value:         values,
			SubAttributes: subAttributesList,
		}
		res = append(res, attributeStruct)
	}
	return res, nil
}

func constructSubAttributeModelList(rawSubAttributes []interface{}) ([]model.PolicySubAttributes, error) {
	res := make([]model.PolicySubAttributes, 0, len(rawSubAttributes))
	for _, rawSubAttribute := range rawSubAttributes {
		rawSubAttributeMap := rawSubAttribute.(map[string]interface{})
		dataType := rawSubAttributeMap["data_type"].(string)
		key := rawSubAttributeMap["key"].(string)
		values := interface2StringList(rawSubAttributeMap["value"].(*schema.Set).List())
		subAttributeStruct := model.PolicySubAttributes{
			Datatype: &dataType,
			Key:      &key,
			Value:    values,
		}
		res = append(res, subAttributeStruct)
	}
	return res, nil
}

func fillAttributesInSchema(policyAttributes []model.PolicyAttributes) []map[string]interface{} {
	attributes := make([]map[string]interface{}, 0, len(policyAttributes))
	for _, policyAttribute := range policyAttributes {
		elem := make(map[string]interface{})
		elem["description"] = policyAttribute.Description
		elem["data_type"] = policyAttribute.Datatype
		elem["key"] = policyAttribute.Key
		elem["value"] = policyAttribute.Value
		elem["is_alg_type"] = policyAttribute.IsALGType
		elem["sub_attribute"] = fillSubAttributesInSchema(policyAttribute.SubAttributes)
		attributes = append(attributes, elem)
	}
	return attributes
}

func fillSubAttributesInSchema(policySubAttributes []model.PolicySubAttributes) []map[string]interface{} {
	subAttributes := make([]map[string]interface{}, 0, len(policySubAttributes))
	for _, policySubAttribute := range policySubAttributes {
		elem := make(map[string]interface{})
		elem["data_type"] = policySubAttribute.Datatype
		elem["key"] = policySubAttribute.Key
		elem["value"] = policySubAttribute.Value
		subAttributes = append(subAttributes, elem)
	}
	return subAttributes
}

func getPolicyAttributeSubAttributesSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"data_type": {
					Type:         schema.TypeString,
					Description:  "Data type of sub attribute",
					Required:     true,
					ValidateFunc: validation.StringInSlice(subAttributeDataTypes, false),
				},
				"key": {
					Type:         schema.TypeString,
					Description:  "Key for attribute",
					Required:     true,
					ValidateFunc: validation.StringInSlice(subAttributeKeys, false),
				},
				"value": {
					Type:        schema.TypeSet,
					Description: "Values for attribute key",
					Required:    true,
					MinItems:    1,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}
}
