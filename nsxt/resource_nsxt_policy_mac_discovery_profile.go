/* Copyright © 2020 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"log"
	"reflect"
	"strings"

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
	schemaType          string
	readOnly            bool
	sdkFieldName        string
	introducedInVersion string
	skip                bool
	reflectType         reflect.Type
	testData            testdata
}

var macDiscoveryProfileMetadata = map[string]*metadata{
	"mac_change_enabled": {
		schemaType:   "bool",
		sdkFieldName: "MacChangeEnabled",
		testData: testdata{
			createValue: "true",
			updateValue: "false",
		},
	},
	"mac_learning_enabled": {
		schemaType:   "bool",
		sdkFieldName: "MacLearningEnabled",
		testData: testdata{
			createValue: "true",
			updateValue: "false",
		},
	},
	"mac_limit": {
		schemaType:   "int",
		sdkFieldName: "MacLimit",
		testData: testdata{
			createValue: "20",
			updateValue: "50",
		},
	},
	"mac_limit_policy": {
		schemaType:   "string",
		sdkFieldName: "MacLimitPolicy",
		testData: testdata{
			createValue: "ALLOW",
			updateValue: "DROP",
		},
	},
	"remote_overlay_mac_limit": {
		schemaType:   "int",
		sdkFieldName: "RemoteOverlayMacLimit",
		testData: testdata{
			createValue: "2048",
			updateValue: "4096",
		},
	},
	"unknown_unicast_flooding_enabled": {
		schemaType:          "bool",
		sdkFieldName:        "UnknownUnicastFloodingEnabled",
		introducedInVersion: "4.0.0",
		testData: testdata{
			createValue: "true",
			updateValue: "false",
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

		Schema: map[string]*schema.Schema{
			"nsx_id":       getNsxIDSchema(),
			"path":         getPathSchema(),
			"display_name": getDisplayNameSchema(),
			"description":  getDescriptionSchema(),
			"revision":     getRevisionSchema(),
			"tag":          getTagsSchema(),
			"context":      getContextSchema(false, false),
			"mac_change_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"mac_learning_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"mac_limit": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 4096),
				Default:      4096,
			},
			"mac_limit_policy": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(macDiscoveryProfileMacLimitPolicyValues, false),
				Optional:     true,
				Default:      "ALLOW",
			},
			"remote_overlay_mac_limit": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(2048, 8192),
				Default:      2048,
			},
			"unknown_unicast_flooding_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
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

func structToSchema(elem reflect.Value, d *schema.ResourceData, metadata map[string]*metadata, parent string, parentMap map[string]interface{}) {
	for key, item := range metadata {
		if len(parent) > 0 && !strings.HasPrefix(key, parent+".") {
			continue
		}
		if item.skip {
			continue
		}

		log.Printf("[INFO] inspecting key %s", key)
		if len(parent) > 0 {
			key = key[len(parent)+1:]
			log.Printf("[INFO] parent %s key %s", parent, key)
		} else if strings.Contains(key, ".") {
			continue
		}
		if item.schemaType == "struct" {
			nestedObj := elem.FieldByName(item.sdkFieldName)
			nestedSchema := make(map[string]interface{})
			structToSchema(nestedObj.Elem(), d, metadata, key, nestedSchema)
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
				log.Printf("[INFO] assigning nested value %v to %s", elem.FieldByName(item.sdkFieldName).Interface(), key)
				parentMap[key] = elem.FieldByName(item.sdkFieldName).Interface()
			} else {
				log.Printf("[INFO] assigning value %v to %s", elem.FieldByName(item.sdkFieldName).Interface(), key)
				d.Set(key, elem.FieldByName(item.sdkFieldName).Interface())
			}
		}
	}
}

func schemaToStruct(elem reflect.Value, d *schema.ResourceData, metadata map[string]*metadata, parent string, parentMap map[string]interface{}) {
	for key, item := range metadata {
		if len(parent) > 0 && !strings.HasPrefix(key, parent+".") {
			continue
		}
		if item.readOnly || item.skip {
			continue
		}
		if item.introducedInVersion != "" && nsxVersionLower(item.introducedInVersion) {
			continue
		}

		log.Printf("[INFO] inspecting key %s", key)
		if len(parent) > 0 {
			key = key[len(parent)+1:]
			log.Printf("[INFO] parent %s key %s", parent, key)
		} else if strings.Contains(key, ".") {
			continue
		}
		if item.schemaType == "string" {
			var value string
			if len(parent) > 0 {
				value = parentMap[key].(string)
			} else {
				value = d.Get(key).(string)
			}
			log.Printf("[INFO] assigning string %v to %s", value, key)
			elem.FieldByName(item.sdkFieldName).Set(reflect.ValueOf(&value))
		}
		if item.schemaType == "bool" {
			var value bool
			if len(parent) > 0 {
				value = parentMap[key].(bool)
			} else {
				value = d.Get(key).(bool)
			}
			log.Printf("[INFO] assigning bool %v to %s", value, key)
			elem.FieldByName(item.sdkFieldName).Set(reflect.ValueOf(&value))
		}
		if item.schemaType == "int" {
			var value int64
			if len(parent) > 0 {
				value = int64(parentMap[key].(int))
			} else {
				value = int64(d.Get(key).(int))
			}
			log.Printf("[INFO] assigning int %v to %s", value, key)
			elem.FieldByName(item.sdkFieldName).Set(reflect.ValueOf(&value))
		}
		if item.schemaType == "struct" {
			/*
				TODO - get reflect Type from sdk
				structType := bindingType.(bindings.StructType)
				fieldType := structType.Field(key)
			*/
			nestedObj := reflect.New(item.reflectType)
			/*
				slice := reflect.MakeSlice(reflect.TypeOf(nestedObj), 1, 1)
			*/
			nestedSchemaList := d.Get(key).([]interface{})
			if len(nestedSchemaList) == 0 {
				continue
			}
			nestedSchema := nestedSchemaList[0].(map[string]interface{})

			schemaToStruct(nestedObj.Elem(), d, metadata, key, nestedSchema)
			//elem.FieldByName(item.sdkFieldName).Set(reflect.ValueOf(nestedObj))
			log.Printf("[INFO] assigning struct %v to %s", nestedObj, key)
			elem.FieldByName(item.sdkFieldName).Set(nestedObj)
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

	tags := getPolicyTagsFromSchema(d)
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	/*
		macChangeEnabled := d.Get("mac_change_enabled").(bool)
		macLearningEnabled := d.Get("mac_learning_enabled").(bool)
		macLimit := int64(d.Get("mac_limit").(int))
		macLimitPolicy := d.Get("mac_limit_policy").(string)
		remoteOverlayMacLimit := int64(d.Get("remote_overlay_mac_limit").(int))
		unknownUnicastFloodingEnabled := d.Get("unknown_unicast_flooding_enabled").(bool)
	*/

	obj := model.MacDiscoveryProfile{
		Tags:        tags,
		DisplayName: &displayName,
		Description: &description,
		/*
			MacChangeEnabled:              &macChangeEnabled,
			MacLearningEnabled:            &macLearningEnabled,
			MacLimit:                      &macLimit,
			MacLimitPolicy:                &macLimitPolicy,
			RemoteOverlayMacLimit:         &remoteOverlayMacLimit,
			UnknownUnicastFloodingEnabled: &unknownUnicastFloodingEnabled,
		*/
	}

	elem := reflect.ValueOf(&obj).Elem()
	schemaToStruct(elem, d, macDiscoveryProfileMetadata, "", nil)

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
	/*
		for key, item := range macDiscoveryProfileMetadata {
			d.Set(key, elem.FieldByName(item.sdkFieldName).Interface())
		}
	*/
	structToSchema(elem, d, macDiscoveryProfileMetadata, "", nil)
	/*
		d.Set("display_name", obj.DisplayName)
		d.Set("description", obj.Description)

		d.Set("mac_change_enabled", obj.MacChangeEnabled)
		d.Set("mac_learning_enabled", obj.MacLearningEnabled)
		d.Set("mac_limit", obj.MacLimit)
		d.Set("mac_limit_policy", obj.MacLimitPolicy)
		d.Set("remote_overlay_mac_limit", obj.RemoteOverlayMacLimit)
		d.Set("unknown_unicast_flooding_enabled", obj.UnknownUnicastFloodingEnabled)
	*/

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

	/*
		macChangeEnabled := d.Get("mac_change_enabled").(bool)
		macLearningEnabled := d.Get("mac_learning_enabled").(bool)
		macLimit := int64(d.Get("mac_limit").(int))
		macLimitPolicy := d.Get("mac_limit_policy").(string)
		remoteOverlayMacLimit := int64(d.Get("remote_overlay_mac_limit").(int))
		unknownUnicastFloodingEnabled := d.Get("unknown_unicast_flooding_enabled").(bool)
	*/
	revision := int64(d.Get("revision").(int))

	obj := model.MacDiscoveryProfile{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
		Revision:    &revision,
	}

	elem := reflect.ValueOf(&obj).Elem()
	schemaToStruct(elem, d, macDiscoveryProfileMetadata, "", nil)

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
