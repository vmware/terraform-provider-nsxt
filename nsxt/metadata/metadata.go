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
				shallowCopy.Elem = &elem.Schema
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
func StructToSchema(elem reflect.Value, d *schema.ResourceData, metadata map[string]*ExtendedSchema, parent string, parentMap map[string]interface{}) {
	for key, item := range metadata {
		if item.Metadata.Skip {
			continue
		}

		log.Printf("[INFO] inspecting key %s", key)
		if len(parent) > 0 {
			log.Printf("[INFO] parent %s key %s", parent, key)
		}
		if elem.FieldByName(item.Metadata.SdkFieldName).IsNil() {
			log.Printf("[INFO] skip key %s with nil value", key)
			continue
		}
		if item.Metadata.SchemaType == "struct" {
			nestedObj := elem.FieldByName(item.Metadata.SdkFieldName)
			nestedSchema := make(map[string]interface{})
			childElem := item.Schema.Elem.(*ExtendedResource)
			StructToSchema(nestedObj.Elem(), d, childElem.Schema, key, nestedSchema)
			log.Printf("[INFO] assigning struct %+v to %s", nestedObj, key)
			var nestedSlice []map[string]interface{}
			nestedSlice = append(nestedSlice, nestedSchema)
			if len(parent) > 0 {
				parentMap[key] = nestedSlice
			} else {
				d.Set(key, nestedSlice)
			}
		} else if item.Metadata.SchemaType == "list" || item.Metadata.SchemaType == "set" {
			if _, ok := item.Schema.Elem.(*ExtendedSchema); ok {
				// List of string, bool, int
				nestedSlice := elem.FieldByName(item.Metadata.SdkFieldName)
				log.Printf("[INFO] assigning slice %v to %s", nestedSlice.Interface(), key)
				if len(parent) > 0 {
					parentMap[key] = nestedSlice.Interface()
				} else {
					d.Set(key, nestedSlice.Interface())
				}
			} else if childElem, ok := item.Schema.Elem.(*ExtendedResource); ok {
				// List of struct
				sliceElem := elem.FieldByName(item.Metadata.SdkFieldName)
				var nestedSlice []map[string]interface{}
				for i := 0; i < sliceElem.Len(); i++ {
					nestedSchema := make(map[string]interface{})
					StructToSchema(sliceElem.Index(i), d, childElem.Schema, key, nestedSchema)
					nestedSlice = append(nestedSlice, nestedSchema)
					log.Printf("[INFO] appending slice item %+v to %s", nestedSchema, key)
				}
				if len(parent) > 0 {
					parentMap[key] = nestedSlice
				} else {
					d.Set(key, nestedSlice)
				}
			}
		} else {
			if len(parent) > 0 {
				log.Printf("[INFO] assigning nested value %+v to %s", elem.FieldByName(item.Metadata.SdkFieldName).Interface(), key)
				parentMap[key] = elem.FieldByName(item.Metadata.SdkFieldName).Interface()
			} else {
				log.Printf("[INFO] assigning value %+v to %s", elem.FieldByName(item.Metadata.SdkFieldName).Interface(), key)
				d.Set(key, elem.FieldByName(item.Metadata.SdkFieldName).Interface())
			}
		}
	}
}

// SchemaToStruct converts terraform schema to NSX model struct
// currently supports nested subtype and trivial types
func SchemaToStruct(elem reflect.Value, d *schema.ResourceData, metadata map[string]*ExtendedSchema, parent string, parentMap map[string]interface{}) {
	for key, item := range metadata {
		if item.Metadata.ReadOnly || item.Metadata.Skip {
			continue
		}
		if item.Metadata.IntroducedInVersion != "" && util.NsxVersionLower(item.Metadata.IntroducedInVersion) {
			continue
		}

		log.Printf("[INFO] inspecting key %s with type %s", key, item.Metadata.SchemaType)
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
			var itemList []interface{}
			if len(parent) > 0 {
				itemList = parentMap[key].([]interface{})
			} else {
				itemList = d.Get(key).([]interface{})
			}
			if len(itemList) == 0 {
				continue
			}
			nestedSchema := itemList[0].(map[string]interface{})

			childElem := item.Schema.Elem.(*ExtendedResource)
			SchemaToStruct(nestedObj.Elem(), d, childElem.Schema, key, nestedSchema)
			log.Printf("[INFO] assigning struct %v to %s", nestedObj, key)
			elem.FieldByName(item.Metadata.SdkFieldName).Set(nestedObj)
		}
		if item.Metadata.SchemaType == "list" || item.Metadata.SchemaType == "set" {
			var itemList []interface{}
			if len(parent) > 0 {
				itemList = parentMap[key].([]interface{})
			} else {
				itemList = d.Get(key).([]interface{})
			}
			if len(itemList) == 0 {
				continue
			}

			// List of string, bool, int
			if childElem, ok := item.Schema.Elem.(*ExtendedSchema); ok {
				sliceElem := elem.FieldByName(item.Metadata.SdkFieldName)
				switch childElem.Metadata.SchemaType {
				case "string":
					sliceElem.Set(
						reflect.MakeSlice(reflect.TypeOf([]string{}), len(itemList), len(itemList)))
				case "bool":
					sliceElem.Set(
						reflect.MakeSlice(reflect.TypeOf([]bool{}), len(itemList), len(itemList)))
				case "int":
					sliceElem.Set(
						reflect.MakeSlice(reflect.TypeOf([]int64{}), len(itemList), len(itemList)))
				}

				for i, v := range itemList {
					sliceElem.Index(i).Set(reflect.ValueOf(v))
					log.Printf("[INFO] appending %v to %s", v, key)
				}
			}

			// List of struct
			if childElem, ok := item.Schema.Elem.(*ExtendedResource); ok {
				sliceElem := elem.FieldByName(item.Metadata.SdkFieldName)
				sliceElem.Set(
					reflect.MakeSlice(reflect.SliceOf(item.Metadata.ReflectType), len(itemList), len(itemList)))
				for i, childItem := range itemList {
					nestedObj := reflect.New(item.Metadata.ReflectType)
					nestedSchema := childItem.(map[string]interface{})
					SchemaToStruct(nestedObj.Elem(), d, childElem.Schema, key, nestedSchema)
					sliceElem.Index(i).Set(nestedObj.Elem())
					log.Printf("[INFO] appending %+v to %s", nestedObj.Elem(), key)
				}
			}
		}
	}
}
