/* Copyright © 2020 VMware, Inc. All Rights Reserved.
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

	"github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

var macDiscoveryProfileMacLimitPolicyValues = []string{
	model.MacDiscoveryProfile_MAC_LIMIT_POLICY_ALLOW,
	model.MacDiscoveryProfile_MAC_LIMIT_POLICY_DROP,
}

type testdata struct {
	createValue interface{}
	updateValue interface{}
}

type metadata struct {
	// we need a separate schema type, in addition to terraform SDK type,
	// in order to distinguish between single subclause and a list of entries
	schemaType   string
	readOnly     bool
	sdkFieldName string
	// This attribute is parent path for the object
	isParentPath        bool
	introducedInVersion string
	// skip handling of this attribute - it will be done manually
	skip        bool
	reflectType reflect.Type
	testData    testdata
}

type extendedSchema struct {
	s schema.Schema
	m metadata
}

type extendedResource struct {
	Schema map[string]*extendedSchema
}

// a helper to convert terraform sdk schema to extended schema
func getExtendedSchema(sch *schema.Schema) *extendedSchema {
	shallowCopy := *sch
	return &extendedSchema{
		s: shallowCopy,
		m: metadata{
			skip: true,
		},
	}
}

// get terraform sdk schema from extended schema definition
func getSchemaFromExtendedSchema(ext map[string]*extendedSchema) map[string]*schema.Schema {
	result := make(map[string]*schema.Schema)

	for key, value := range ext {
		log.Printf("[INFO] inspecting schema key %s, value %v", key, value)
		shallowCopy := value.s
		if (value.s.Type == schema.TypeList) || (value.s.Type == schema.TypeSet) {
			elem, ok := shallowCopy.Elem.(*extendedSchema)
			if ok {
				shallowCopy.Elem = elem.s
			} else {
				elem, ok := shallowCopy.Elem.(*extendedResource)
				if ok {
					shallowCopy.Elem = &schema.Resource{
						Schema: getSchemaFromExtendedSchema(elem.Schema),
					}
				}
			}
		}
		// TODO: deepcopy needed?
		result[key] = &shallowCopy
	}

	return result
}

var macDiscoveryProfileSchema = map[string]*extendedSchema{
	"nsx_id":       getExtendedSchema(getNsxIDSchema()),
	"path":         getExtendedSchema(getPathSchema()),
	"display_name": getExtendedSchema(getDisplayNameSchema()),
	"description":  getExtendedSchema(getDescriptionSchema()),
	"revision":     getExtendedSchema(getRevisionSchema()),
	"tag":          getExtendedSchema(getTagsSchema()),
	"context":      getExtendedSchema(getContextSchema(false, false)),
	"mac_change_enabled": {
		s: schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		m: metadata{
			schemaType:   "bool",
			sdkFieldName: "MacChangeEnabled",
			testData: testdata{
				createValue: "true",
				updateValue: "false",
			},
		},
	},
	"mac_learning_enabled": {
		s: schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		m: metadata{
			schemaType:   "bool",
			sdkFieldName: "MacLearningEnabled",
			testData: testdata{
				createValue: "true",
				updateValue: "false",
			},
		},
	},
	"mac_limit": {
		s: schema.Schema{
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(0, 4096),
			Default:      4096,
		},
		m: metadata{
			schemaType:   "int",
			sdkFieldName: "MacLimit",
			testData: testdata{
				createValue: "20",
				updateValue: "50",
			},
		},
	},
	"mac_limit_policy": {
		s: schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validation.StringInSlice(macDiscoveryProfileMacLimitPolicyValues, false),
			Optional:     true,
			Default:      "ALLOW",
		},
		m: metadata{
			schemaType:   "string",
			sdkFieldName: "MacLimitPolicy",
			testData: testdata{
				createValue: "ALLOW",
				updateValue: "DROP",
			},
		},
	},
	"remote_overlay_mac_limit": {
		s: schema.Schema{
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(2048, 8192),
			Default:      2048,
		},
		m: metadata{
			schemaType:   "int",
			sdkFieldName: "RemoteOverlayMacLimit",
			testData: testdata{
				createValue: "2048",
				updateValue: "4096",
			},
		},
	},
	"unknown_unicast_flooding_enabled": {
		s: schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		m: metadata{
			schemaType:          "bool",
			sdkFieldName:        "UnknownUnicastFloodingEnabled",
			introducedInVersion: "4.0.0",
			testData: testdata{
				createValue: "true",
				updateValue: "false",
			},
		},
	},
}

func resourceNsxtPolicyMacDiscoveryProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyMacDiscoveryProfileCreate,
		Read:   resourceNsxtPolicyMacDiscoveryProfileRead,
		Update: resourceNsxtPolicyMacDiscoveryProfileUpdate,
		Delete: resourceNsxtPolicyMacDiscoveryProfileDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtPolicyPathResourceImporter,
		},
		Schema: getSchemaFromExtendedSchema(macDiscoveryProfileSchema),
	}
}

func resourceNsxtPolicyMacDiscoveryProfileExists(sessionContext utl.SessionContext, id string, connector client.Connector) (bool, error) {
	var err error
	client := infra.NewMacDiscoveryProfilesClient(sessionContext, connector)
	_, err = client.Get(id)
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving resource", err)
}

// convert NSX model struct to terraform schema
// currently supports nested subtype and trivial types
// TODO - support a list of structs
func structToSchema(elem reflect.Value, d *schema.ResourceData, metadata map[string]*extendedSchema, parent string, parentMap map[string]interface{}) {
	for key, item := range metadata {
		if item.m.skip {
			continue
		}

		log.Printf("[INFO] inspecting key %s", key)
		if len(parent) > 0 {
			log.Printf("[INFO] parent %s key %s", parent, key)
		}
		if item.m.schemaType == "struct" {
			nestedObj := elem.FieldByName(item.m.sdkFieldName)
			nestedSchema := make(map[string]interface{})
			childElem := item.s.Elem.(*extendedResource)
			structToSchema(nestedObj.Elem(), d, childElem.Schema, key, nestedSchema)
			log.Printf("[INFO] assigning struct %v to %s", nestedObj, key)
			// TODO - get the schema from nested type if parent in present
			var nestedSlice []map[string]interface{}
			nestedSlice = append(nestedSlice, nestedSchema)
			if len(parent) > 0 {
				parentMap[key] = nestedSlice
			} else {
				d.Set(key, nestedSlice)
			}
		} else {
			if len(parent) > 0 {
				log.Printf("[INFO] assigning nested value %v to %s", elem.FieldByName(item.m.sdkFieldName).Interface(), key)
				parentMap[key] = elem.FieldByName(item.m.sdkFieldName).Interface()
			} else {
				log.Printf("[INFO] assigning value %v to %s", elem.FieldByName(item.m.sdkFieldName).Interface(), key)
				d.Set(key, elem.FieldByName(item.m.sdkFieldName).Interface())
			}
		}
	}
}

// convert terraform schema to NSX model struct
// currently supports nested subtype and trivial types
// TODO - support a list of structs
func schemaToStruct(elem reflect.Value, d *schema.ResourceData, metadata map[string]*extendedSchema, parent string, parentMap map[string]interface{}) {
	for key, item := range metadata {
		if item.m.readOnly || item.m.skip {
			continue
		}
		if item.m.introducedInVersion != "" && nsxVersionLower(item.m.introducedInVersion) {
			continue
		}

		log.Printf("[INFO] inspecting key %s", key)
		if len(parent) > 0 {
			log.Printf("[INFO] parent %s key %s", parent, key)
		}
		if item.m.schemaType == "string" {
			var value string
			if len(parent) > 0 {
				value = parentMap[key].(string)
			} else {
				value = d.Get(key).(string)
			}
			log.Printf("[INFO] assigning string %v to %s", value, key)
			elem.FieldByName(item.m.sdkFieldName).Set(reflect.ValueOf(&value))
		}
		if item.m.schemaType == "bool" {
			var value bool
			if len(parent) > 0 {
				value = parentMap[key].(bool)
			} else {
				value = d.Get(key).(bool)
			}
			log.Printf("[INFO] assigning bool %v to %s", value, key)
			elem.FieldByName(item.m.sdkFieldName).Set(reflect.ValueOf(&value))
		}
		if item.m.schemaType == "int" {
			var value int64
			if len(parent) > 0 {
				value = int64(parentMap[key].(int))
			} else {
				value = int64(d.Get(key).(int))
			}
			log.Printf("[INFO] assigning int %v to %s", value, key)
			elem.FieldByName(item.m.sdkFieldName).Set(reflect.ValueOf(&value))
		}
		if item.m.schemaType == "struct" {
			nestedObj := reflect.New(item.m.reflectType)
			/*
				// Helper for list of structs
				slice := reflect.MakeSlice(reflect.TypeOf(nestedObj), 1, 1)
			*/
			nestedSchemaList := d.Get(key).([]interface{})
			if len(nestedSchemaList) == 0 {
				continue
			}
			nestedSchema := nestedSchemaList[0].(map[string]interface{})

			childElem := item.s.Elem.(*extendedResource)
			schemaToStruct(nestedObj.Elem(), d, childElem.Schema, key, nestedSchema)
			log.Printf("[INFO] assigning struct %v to %s", nestedObj, key)
			elem.FieldByName(item.m.sdkFieldName).Set(nestedObj)
			// TODO - get the schema from nested type if parent in present
		}
	}
}

func resourceNsxtPolicyMacDiscoveryProfileCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	// Initialize resource Id and verify this ID is not yet used
	id, err := getOrGenerateID2(d, m, resourceNsxtPolicyMacDiscoveryProfileExists)
	if err != nil {
		return err
	}

	// TODO - consider including standard object attributes in the schema
	tags := getPolicyTagsFromSchema(d)
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)

	obj := model.MacDiscoveryProfile{
		Tags:        tags,
		DisplayName: &displayName,
		Description: &description,
	}

	elem := reflect.ValueOf(&obj).Elem()
	schemaToStruct(elem, d, macDiscoveryProfileSchema, "", nil)

	// Create the resource using PATCH
	log.Printf("[INFO] Creating MacDiscoveryProfile with ID %s", id)
	boolFalse := false
	client := infra.NewMacDiscoveryProfilesClient(getSessionContext(d, m), connector)
	err = client.Patch(id, obj, &boolFalse)
	if err != nil {
		return handleCreateError("MacDiscoveryProfile", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyMacDiscoveryProfileRead(d, m)
}

func resourceNsxtPolicyMacDiscoveryProfileRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining MacDiscoveryProfile ID")
	}

	client := infra.NewMacDiscoveryProfilesClient(getSessionContext(d, m), connector)
	obj, err := client.Get(id)
	if err != nil {
		return handleReadError(d, "MacDiscoveryProfile", id, err)
	}

	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("revision", obj.Revision)
	d.Set("path", obj.Path)

	elem := reflect.ValueOf(&obj).Elem()
	structToSchema(elem, d, macDiscoveryProfileSchema, "", nil)

	return nil
}

func resourceNsxtPolicyMacDiscoveryProfileUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining MacDiscoveryProfile ID")
	}

	// Read the rest of the configured parameters
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getPolicyTagsFromSchema(d)

	revision := int64(d.Get("revision").(int))

	obj := model.MacDiscoveryProfile{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
		Revision:    &revision,
	}

	elem := reflect.ValueOf(&obj).Elem()
	schemaToStruct(elem, d, macDiscoveryProfileSchema, "", nil)

	// Update the resource using PATCH
	boolFalse := false
	client := infra.NewMacDiscoveryProfilesClient(getSessionContext(d, m), connector)
	_, err := client.Update(id, obj, &boolFalse)
	if err != nil {
		return handleUpdateError("MacDiscoveryProfile", id, err)
	}

	return resourceNsxtPolicyMacDiscoveryProfileRead(d, m)
}

func resourceNsxtPolicyMacDiscoveryProfileDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining MacDiscoveryProfile ID")
	}

	connector := getPolicyConnector(m)
	boolFalse := false
	client := infra.NewMacDiscoveryProfilesClient(getSessionContext(d, m), connector)
	err := client.Delete(id, &boolFalse)

	if err != nil {
		return handleDeleteError("MacDiscoveryProfile", id, err)
	}

	return nil
}
