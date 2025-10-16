// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"
	"slices"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	clientLayer "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var l7AccessProfileDefaultActionValues = []string{
	model.L7AccessProfile_DEFAULT_ACTION_ALLOW,
	model.L7AccessProfile_DEFAULT_ACTION_REJECT,
	model.L7AccessProfile_DEFAULT_ACTION_REJECT_WITH_RESPONSE,
}

var l7AccessProfileActionValues = []string{
	model.L7AccessEntry_ACTION_ALLOW,
	model.L7AccessEntry_ACTION_REJECT,
	model.L7AccessEntry_ACTION_REJECT_WITH_RESPONSE,
}

var l7AccessProfileAttributeSourceValues = []string{
	model.L7AccessAttributes_ATTRIBUTE_SOURCE_SYSTEM,
	model.L7AccessAttributes_ATTRIBUTE_SOURCE_CUSTOM,
}

var l7AccessProfileSubAttributeKeyValues = []string{
	model.PolicySubAttributes_KEY_TLS_CIPHER_SUITE,
	model.PolicySubAttributes_KEY_TLS_VERSION,
	model.PolicySubAttributes_KEY_CIFS_SMB_VERSION,
}

var l7AccessProfileKeyAttributeValues = []string{
	model.L7AccessAttributes_KEY_APP_ID,
	model.L7AccessAttributes_KEY_DOMAIN_NAME,
	model.L7AccessAttributes_KEY_URL_CATEGORY,
	model.L7AccessAttributes_KEY_URL_REPUTATION,
	model.L7AccessAttributes_KEY_CUSTOM_URL,
}

var l2AccessProfilePathExample = getMultitenancyPathExample("/infra/l7-access-profiles/[profile]")

func resourceNsxtPolicyL7AccessProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyL7AccessProfileCreate,
		Read:   resourceNsxtPolicyL7AccessProfileRead,
		Update: resourceNsxtPolicyL7AccessProfileUpdate,
		Delete: resourceNsxtPolicyL7AccessProfileDelete,
		Importer: &schema.ResourceImporter{
			State: getPolicyPathOrIDResourceImporter(l2AccessProfilePathExample),
		},
		Schema: map[string]*schema.Schema{
			"nsx_id":       getNsxIDSchema(),
			"path":         getPathSchema(),
			"display_name": getDisplayNameSchema(),
			"description":  getDescriptionSchema(),
			"revision":     getRevisionSchema(),
			"tag":          getTagsSchema(),
			"context":      getContextSchema(false, false, false),
			"default_action": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(l7AccessProfileDefaultActionValues, false),
			},
			"default_action_logged": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"l7_access_entry": {
				Type:        schema.TypeList,
				Description: "Array of Policy L7 Access Profile entries",
				Optional:    true,
				MaxItems:    1000,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nsx_id":       getFlexNsxIDSchema(false),
						"display_name": getOptionalDisplayNameSchema(true),
						"description":  getDescriptionSchema(),
						"path":         getPathSchema(),
						"revision":     getRevisionSchema(),
						"action": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(l7AccessProfileActionValues, false),
						},
						"attribute": {
							Type:        schema.TypeList,
							Description: "Array of Policy L7 Access Profile attributes",
							Optional:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"attribute_source": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice(l7AccessProfileAttributeSourceValues, false),
										Description:  "Source of attribute value i.e whether system defined or custom value",
									},
									"custom_url_partial_match": {
										Type:        schema.TypeBool,
										Optional:    true,
										Computed:    true,
										Description: "true value would be treated as a partial match for custom url",
									},
									"description": getDescriptionSchema(),
									"is_alg_type": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Is the value ALG type",
									},
									"key": {
										Type:         schema.TypeString,
										Required:     true,
										Description:  "Key for attribute",
										ValidateFunc: validation.StringInSlice(l7AccessProfileKeyAttributeValues, false),
									},
									"metadata": getKeyValuePairListSchema(),
									"sub_attribute": {
										Type:        schema.TypeList,
										Description: "Reference to sub attributes for the attribute",
										Optional:    true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:         schema.TypeString,
													Required:     true,
													Description:  "Key for sub attribute",
													ValidateFunc: validation.StringInSlice(l7AccessProfileSubAttributeKeyValues, false),
												},
												"values": {
													Type:        schema.TypeList,
													Description: "Value for sub attribute key",
													MinItems:    1,
													Required:    true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
									"values": {
										Type:        schema.TypeList,
										Description: "Value for attribute key",
										MinItems:    1,
										Required:    true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"disabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Flag to deactivate the entry",
						},
						"logged": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Enable logging flag",
						},
						"sequence_number": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Policy L7 Access Entry Order",
						},
					},
				},
			},
		},
	}
}

func resourceNsxtPolicyL7AccessProfileExists(sessionContext utl.SessionContext, id string, connector client.Connector) (bool, error) {
	var err error
	client := clientLayer.NewL7AccessProfilesClient(sessionContext, connector)
	_, err = client.Get(id)
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving resource", err)
}

func populateL7AccessProfileStruct(id string, d *schema.ResourceData, objChildren []*data.StructValue) model.L7AccessProfile {
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)

	defaultAction := d.Get("default_action").(string)
	defaultActionLogged := d.Get("default_action_logged").(bool)

	obj := model.L7AccessProfile{
		Id:                  &id,
		DisplayName:         &displayName,
		Description:         &description,
		Tags:                tags,
		DefaultAction:       &defaultAction,
		DefaultActionLogged: &defaultActionLogged,
		ResourceType:        strPtr("L7AccessProfile"),
		Children:            objChildren,
	}

	return obj
}

func l7AccessProfileToHAPIFormat(obj model.L7AccessProfile) (model.Infra, error) {
	child7AccessProfile := model.ChildL7AccessProfile{
		L7AccessProfile: &obj,
		ResourceType:    "ChildL7AccessProfile",
	}

	converter := bindings.NewTypeConverter()
	dataValue, errors := converter.ConvertToVapi(child7AccessProfile, model.ChildL7AccessProfileBindingType())
	if len(errors) > 0 {
		return model.Infra{}, errors[0]
	}

	infraChildren := []*data.StructValue{dataValue.(*data.StructValue)}
	infraStruct := model.Infra{
		Children:     infraChildren,
		ResourceType: strPtr("Infra"),
	}

	return infraStruct, nil
}

func childL7AccessEntryFromSchema(id string, obj map[string]interface{}, markedForDeletion bool) (*data.StructValue, error) {
	displayName := obj["display_name"].(string)
	description := obj["description"].(string)
	path := obj["path"].(string)
	revision := int64(obj["revision"].(int))
	action := obj["action"].(string)
	disabled := obj["disabled"].(bool)
	logged := obj["logged"].(bool)
	sequenceNumber := int64(obj["sequence_number"].(int))

	var attributes []model.L7AccessAttributes
	attrs := obj["attribute"].([]interface{})
	for _, a := range attrs {
		attr := a.(map[string]interface{})
		attributeSource := attr["attribute_source"].(string)
		customUrlPartialMatch := attr["custom_url_partial_match"].(bool)
		description = attr["description"].(string)
		isALGType := attr["is_alg_type"].(bool)
		key := attr["key"].(string)

		var metadata []model.ContextProfileAttributesMetadata
		md := attr["metadata"].([]interface{})
		for _, m := range md {
			mdItem := m.(map[string]interface{})
			key := mdItem["key"].(string)
			value := mdItem["value"].(string)
			metadata = append(metadata, model.ContextProfileAttributesMetadata{Key: &key, Value: &value})
		}

		var subAttributes []model.PolicySubAttributes
		subAttrs := attr["sub_attribute"].([]interface{})
		for _, sa := range subAttrs {
			subAttr := sa.(map[string]interface{})
			key := subAttr["key"].(string)
			values := interface2StringList(subAttr["values"].([]interface{}))

			subAttributes = append(subAttributes, model.PolicySubAttributes{
				Key:      &key,
				Value:    values,
				Datatype: strPtr(model.PolicySubAttributes_DATATYPE_STRING),
			})
		}

		values := interface2StringList(attr["values"].([]interface{}))

		attribute := model.L7AccessAttributes{
			AttributeSource: &attributeSource,
			Description:     &description,
			IsALGType:       &isALGType,
			Key:             &key,
			Datatype:        strPtr(model.PolicySubAttributes_DATATYPE_STRING),
			Metadata:        metadata,
			SubAttributes:   subAttributes,
			Value:           values,
		}

		if util.NsxVersionHigherOrEqual("4.0.0") {
			attribute.CustomUrlPartialMatch = &customUrlPartialMatch
		}
		attributes = append(attributes, attribute)
	}

	childL7AccessEntry := model.ChildL7AccessEntry{
		L7AccessEntry: &model.L7AccessEntry{
			Id:             &id,
			DisplayName:    &displayName,
			Description:    &description,
			Path:           &path,
			Revision:       &revision,
			Action:         &action,
			Disabled:       &disabled,
			Logged:         &logged,
			SequenceNumber: &sequenceNumber,
			ResourceType:   strPtr("L7AccessEntry"),
			Attributes:     attributes,
		},
		Id:              &id,
		MarkedForDelete: &markedForDeletion,
		ResourceType:    "ChildL7AccessEntry",
	}

	converter := bindings.NewTypeConverter()
	dataValue, errors := converter.ConvertToVapi(childL7AccessEntry, model.ChildL7AccessEntryBindingType())
	if len(errors) > 0 {
		return nil, errors[0]
	}

	return dataValue.(*data.StructValue), nil
}

func populateL7AccessProfileStructChildren(old, new interface{}) ([]*data.StructValue, error) {
	var retList []*data.StructValue

	newObjects := new.([]interface{})
	var newIDs []string
	for _, o := range newObjects {
		obj := o.(map[string]interface{})
		id := obj["nsx_id"].(string)
		if id == "" {
			id = newUUID()
		}
		newIDs = append(newIDs, id)

		childL7AccessEntry, err := childL7AccessEntryFromSchema(id, obj, false)
		if err != nil {
			return nil, err
		}
		retList = append(retList, childL7AccessEntry)
	}

	oldObjects := old.([]interface{})

	for _, o := range oldObjects {
		obj := o.(map[string]interface{})
		id := obj["nsx_id"].(string)

		if !slices.Contains(newIDs, id) {
			// This entry doesn't exist in the new schema, mark it for deletion
			childL7AccessEntry, err := childL7AccessEntryFromSchema(id, obj, true)
			if err != nil {
				return nil, err
			}
			retList = append(retList, childL7AccessEntry)
		}
	}

	return retList, nil
}

func resourceNsxtPolicyL7AccessProfileCreate(d *schema.ResourceData, m interface{}) error {
	id, err := getOrGenerateID2(d, m, resourceNsxtPolicyL7AccessProfileExists)
	if err != nil {
		return err
	}

	objChildren, err := populateL7AccessProfileStructChildren([]interface{}{}, d.Get("l7_access_entry"))
	if err != nil {
		return handleCreateError("L7AccessProfile", id, err)
	}
	obj := populateL7AccessProfileStruct(id, d, objChildren)
	infraStruct, err := l7AccessProfileToHAPIFormat(obj)
	if err != nil {
		return handleCreateError("L7AccessProfile", id, err)
	}

	log.Printf("[INFO] Creating L7AccessProfile with ID %s", id)

	err = policyInfraPatch(commonSessionContext, infraStruct, getPolicyConnector(m), false)
	if err != nil {
		return handleCreateError("L7AccessProfile", id, err)
	}
	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyL7AccessProfileRead(d, m)
}

func resourceNsxtPolicyL7AccessProfileRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining L7AccessProfile ID")
	}

	client := clientLayer.NewL7AccessProfilesClient(commonSessionContext, connector)
	obj, err := client.Get(id)
	if err != nil {
		return handleReadError(d, "L7AccessProfile", id, err)
	}

	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("revision", obj.Revision)
	d.Set("path", obj.Path)

	d.Set("default_action", obj.DefaultAction)
	d.Set("default_action_logged", obj.DefaultActionLogged)

	var l7entries []interface{}
	for _, l7e := range obj.L7AccessEntries {
		l7entry := make(map[string]interface{})
		l7entry["nsx_id"] = l7e.Id
		l7entry["display_name"] = l7e.DisplayName
		l7entry["description"] = l7e.Description
		l7entry["path"] = l7e.Path
		l7entry["revision"] = l7e.Revision
		l7entry["action"] = l7e.Action
		l7entry["disabled"] = l7e.Disabled
		l7entry["logged"] = l7e.Logged
		l7entry["sequence_number"] = l7e.SequenceNumber

		var attributes []interface{}
		for _, attr := range l7e.Attributes {
			attribute := make(map[string]interface{})
			attribute["attribute_source"] = attr.AttributeSource
			attribute["custom_url_partial_match"] = attr.CustomUrlPartialMatch
			attribute["description"] = attr.Description
			attribute["is_alg_type"] = attr.IsALGType
			attribute["key"] = attr.Key
			attribute["values"] = attr.Value

			var metadata []interface{}
			for _, meta := range attr.Metadata {
				mdItem := make(map[string]interface{})
				mdItem["key"] = meta.Key
				mdItem["value"] = meta.Value
				metadata = append(metadata, mdItem)
			}
			attribute["metadata"] = metadata

			var subAttributes []interface{}
			for _, subAttr := range attr.SubAttributes {
				subAttribute := make(map[string]interface{})
				subAttribute["key"] = subAttr.Key
				subAttribute["values"] = subAttr.Value

				subAttributes = append(subAttributes, subAttribute)
			}
			attribute["sub_attribute"] = subAttributes

			attributes = append(attributes, attribute)
		}
		l7entry["attribute"] = attributes
		l7entries = append(l7entries, l7entry)
	}
	d.Set("l7_access_entry", l7entries)

	return nil
}

func resourceNsxtPolicyL7AccessProfileUpdate(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining L7AccessProfile ID")
	}

	oldEntries, newEntries := d.GetChange("l7_access_entry")
	objChildren, err := populateL7AccessProfileStructChildren(oldEntries, newEntries)
	if err != nil {
		return handleUpdateError("L7AccessProfile", id, err)
	}

	obj := populateL7AccessProfileStruct(id, d, objChildren)
	revision := int64(d.Get("revision").(int))
	obj.Revision = &revision
	infraStruct, err := l7AccessProfileToHAPIFormat(obj)
	if err != nil {
		return handleUpdateError("L7AccessProfile", id, err)
	}

	log.Printf("[INFO] Updating L7AccessProfile with ID %s", id)

	err = policyInfraPatch(commonSessionContext, infraStruct, getPolicyConnector(m), false)
	if err != nil {
		return handleUpdateError("L7AccessProfile", id, err)
	}

	return resourceNsxtPolicyL7AccessProfileRead(d, m)
}

func resourceNsxtPolicyL7AccessProfileDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining L7AccessProfile ID")
	}

	connector := getPolicyConnector(m)

	client := clientLayer.NewL7AccessProfilesClient(commonSessionContext, connector)
	err := client.Delete(id, nil)

	if err != nil {
		return handleDeleteError("L7AccessProfile", id, err)
	}

	return nil
}
