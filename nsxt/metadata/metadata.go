package metadata

import (
	"log"
	"reflect"

	"github.com/vmware/terraform-provider-nsxt/nsxt/util"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Metadata struct {
	// we need a separate schema type, in addition to terraform SDK type,
	// in order to distinguish between single subclause and a list of entries
	SchemaType   string
	ReadOnly     bool
	SdkFieldName string
	// This attribute is parent path for the object
	IsParentPath        bool
	IntroducedInVersion string
	// skip handling of this attribute - it will be done manually
	Skip        bool
	ReflectType reflect.Type
	TestData    Testdata
}

type ExtendedSchema struct {
	Schema   schema.Schema
	Metadata Metadata
}

type ExtendedResource struct {
	Schema map[string]*ExtendedSchema
}

type Testdata struct {
	CreateValue interface{}
	UpdateValue interface{}
}

// GetExtendedSchema is a helper to convert terraform sdk schema to extended schema
func GetExtendedSchema(sch *schema.Schema) *ExtendedSchema {
	shallowCopy := *sch
	return &ExtendedSchema{
		Schema: shallowCopy,
		Metadata: Metadata{
			Skip: true,
		},
	}
}

// GetSchemaFromExtendedSchema gets terraform sdk schema from extended schema definition
func GetSchemaFromExtendedSchema(ext map[string]*ExtendedSchema) map[string]*schema.Schema {
	result := make(map[string]*schema.Schema)

	for key, value := range ext {
		log.Printf("[INFO] inspecting schema key %s, value %v", key, value)
		shallowCopy := value.Schema
		if (value.Schema.Type == schema.TypeList) || (value.Schema.Type == schema.TypeSet) {
			elem, ok := shallowCopy.Elem.(*ExtendedSchema)
			if ok {
				shallowCopy.Elem = elem.Schema
			} else {
				elem, ok := shallowCopy.Elem.(*ExtendedResource)
				if ok {
					shallowCopy.Elem = &schema.Resource{
						Schema: GetSchemaFromExtendedSchema(elem.Schema),
					}
				}
			}
		}
		// TODO: deepcopy needed?
		result[key] = &shallowCopy
	}

	return result
}

// StructToSchema converts NSX model struct to terraform schema
// currently supports nested subtype and trivial types
// TODO - support a list of structs
func StructToSchema(elem reflect.Value, d *schema.ResourceData, metadata map[string]*ExtendedSchema, parent string, parentMap map[string]interface{}) {
	for key, item := range metadata {
		if item.Metadata.Skip {
			continue
		}

		log.Printf("[INFO] inspecting key %s", key)
		if len(parent) > 0 {
			log.Printf("[INFO] parent %s key %s", parent, key)
		}
		if item.Metadata.SchemaType == "struct" {
			nestedObj := elem.FieldByName(item.Metadata.SdkFieldName)
			nestedSchema := make(map[string]interface{})
			childElem := item.Schema.Elem.(*ExtendedResource)
			StructToSchema(nestedObj.Elem(), d, childElem.Schema, key, nestedSchema)
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
				log.Printf("[INFO] assigning nested value %v to %s", elem.FieldByName(item.Metadata.SdkFieldName).Interface(), key)
				parentMap[key] = elem.FieldByName(item.Metadata.SdkFieldName).Interface()
			} else {
				log.Printf("[INFO] assigning value %v to %s", elem.FieldByName(item.Metadata.SdkFieldName).Interface(), key)
				d.Set(key, elem.FieldByName(item.Metadata.SdkFieldName).Interface())
			}
		}
	}
}

// SchemaToStruct converts terraform schema to NSX model struct
// currently supports nested subtype and trivial types
// TODO - support a list of structs
func SchemaToStruct(elem reflect.Value, d *schema.ResourceData, metadata map[string]*ExtendedSchema, parent string, parentMap map[string]interface{}) {
	for key, item := range metadata {
		if item.Metadata.ReadOnly || item.Metadata.Skip {
			continue
		}
		if item.Metadata.IntroducedInVersion != "" && util.NsxVersionLower(item.Metadata.IntroducedInVersion) {
			continue
		}

		log.Printf("[INFO] inspecting key %s", key)
		if len(parent) > 0 {
			log.Printf("[INFO] parent %s key %s", parent, key)
		}
		if item.Metadata.SchemaType == "string" {
			var value string
			if len(parent) > 0 {
				value = parentMap[key].(string)
			} else {
				value = d.Get(key).(string)
			}
			log.Printf("[INFO] assigning string %v to %s", value, key)
			elem.FieldByName(item.Metadata.SdkFieldName).Set(reflect.ValueOf(&value))
		}
		if item.Metadata.SchemaType == "bool" {
			var value bool
			if len(parent) > 0 {
				value = parentMap[key].(bool)
			} else {
				value = d.Get(key).(bool)
			}
			log.Printf("[INFO] assigning bool %v to %s", value, key)
			elem.FieldByName(item.Metadata.SdkFieldName).Set(reflect.ValueOf(&value))
		}
		if item.Metadata.SchemaType == "int" {
			var value int64
			if len(parent) > 0 {
				value = int64(parentMap[key].(int))
			} else {
				value = int64(d.Get(key).(int))
			}
			log.Printf("[INFO] assigning int %v to %s", value, key)
			elem.FieldByName(item.Metadata.SdkFieldName).Set(reflect.ValueOf(&value))
		}
		if item.Metadata.SchemaType == "struct" {
			nestedObj := reflect.New(item.Metadata.ReflectType)
			/*
				// Helper for list of structs
				slice := reflect.MakeSlice(reflect.TypeOf(nestedObj), 1, 1)
			*/
			nestedSchemaList := d.Get(key).([]interface{})
			if len(nestedSchemaList) == 0 {
				continue
			}
			nestedSchema := nestedSchemaList[0].(map[string]interface{})

			childElem := item.Schema.Elem.(*ExtendedResource)
			SchemaToStruct(nestedObj.Elem(), d, childElem.Schema, key, nestedSchema)
			log.Printf("[INFO] assigning struct %v to %s", nestedObj, key)
			elem.FieldByName(item.Metadata.SdkFieldName).Set(nestedObj)
			// TODO - get the schema from nested type if parent in present
		}
	}
}
