package metadata

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	StringField    *string
	BoolField      *bool
	IntField       *int64
	StringFieldNil *string
	BoolFieldNil   *bool
	IntFieldNil    *int64
	StructField    *testNestedStruct
}

type testNestedStruct struct {
	StringField *string
	BoolField   *bool
	IntField    *int64
}

var testSchema = map[string]*schema.Schema{
	"string_field": {
		Type: schema.TypeString,
	},
	"bool_field": {
		Type: schema.TypeBool,
	},
	"int_field": {
		Type: schema.TypeInt,
	},
	"string_field_nil": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"bool_field_nil": {
		Type:     schema.TypeBool,
		Optional: true,
	},
	"int_field_nil": {
		Type:     schema.TypeInt,
		Optional: true,
	},
	"struct_field": {
		Type:     schema.TypeList,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"string_field": {
					Type: schema.TypeString,
				},
				"bool_field": {
					Type: schema.TypeBool,
				},
				"int_field": {
					Type: schema.TypeInt,
				},
			},
		},
	},
}

func basicStringSchema(sdkName string) *ExtendedSchema {
	return &ExtendedSchema{
		Schema: schema.Schema{
			Type: schema.TypeString,
		},
		Metadata: Metadata{
			SchemaType:   "string",
			SdkFieldName: sdkName,
		},
	}
}

func basicBoolSchema(sdkName string) *ExtendedSchema {
	return &ExtendedSchema{
		Schema: schema.Schema{
			Type: schema.TypeBool,
		},
		Metadata: Metadata{
			SchemaType:   "bool",
			SdkFieldName: sdkName,
		},
	}
}

func basicIntSchema(sdkName string) *ExtendedSchema {
	return &ExtendedSchema{
		Schema: schema.Schema{
			Type: schema.TypeInt,
		},
		Metadata: Metadata{
			SchemaType:   "int",
			SdkFieldName: sdkName,
		},
	}
}

var testExtendedSchema = map[string]*ExtendedSchema{
	"string_field":     basicStringSchema("StringField"),
	"bool_field":       basicBoolSchema("BoolField"),
	"int_field":        basicIntSchema("IntField"),
	"string_field_nil": basicStringSchema("StringFieldNil"),
	"bool_field_nil":   basicBoolSchema("BoolFieldNil"),
	"int_field_nil":    basicIntSchema("IntFieldNil"),
	"struct_field": {
		Schema: schema.Schema{
			Type:     schema.TypeList,
			MaxItems: 1,
			Elem: &ExtendedResource{
				Schema: map[string]*ExtendedSchema{
					"string_field": basicStringSchema("StringField"),
					"bool_field":   basicBoolSchema("BoolField"),
					"int_field":    basicIntSchema("IntField"),
				},
			},
		},
		Metadata: Metadata{
			SchemaType:   "struct",
			SdkFieldName: "StructField",
			ReflectType:  reflect.TypeOf(testNestedStruct{}),
		},
	},
}

func TestStructToSchema(t *testing.T) {
	testStr := "test_string"
	testBool := true
	testInt := int64(123)
	obj := testStruct{
		StringField: &testStr,
		BoolField:   &testBool,
		IntField:    &testInt,
		StructField: &testNestedStruct{
			StringField: &testStr,
			BoolField:   &testBool,
			IntField:    &testInt,
		},
	}
	d := schema.TestResourceDataRaw(
		t, testSchema, map[string]interface{}{})

	elem := reflect.ValueOf(&obj).Elem()
	StructToSchema(elem, d, testExtendedSchema, "", nil)

	// Base types
	assert.Equal(t, "test_string", d.Get("string_field").(string))
	assert.Equal(t, true, d.Get("bool_field").(bool))
	assert.Equal(t, 123, d.Get("int_field").(int))

	// Zero values
	_, ok := d.GetOk("string_field_nil")
	assert.False(t, ok)
	_, ok = d.GetOk("bool_field_nil")
	assert.False(t, ok)
	_, ok = d.GetOk("int_field_nil")
	assert.False(t, ok)

	// Nested struct
	nestedObj := d.Get("struct_field").([]interface{})[0].(map[string]interface{})
	assert.Equal(t, "test_string", nestedObj["string_field"].(string))
	assert.Equal(t, true, nestedObj["bool_field"].(bool))
	assert.Equal(t, 123, nestedObj["int_field"].(int))
}

func TestSchemaToStruct(t *testing.T) {
	d := schema.TestResourceDataRaw(
		t, testSchema, map[string]interface{}{
			"string_field": "test_string",
			"bool_field":   true,
			"int_field":    100,
			"struct_field": []interface{}{
				map[string]interface{}{
					"string_field": "nested_string",
					"bool_field":   true,
					"int_field":    1,
				},
			},
		})

	obj := testStruct{}
	elem := reflect.ValueOf(&obj).Elem()
	SchemaToStruct(elem, d, testExtendedSchema, "", nil)

	// Base types
	assert.Equal(t, "test_string", *obj.StringField)
	assert.Equal(t, true, *obj.BoolField)
	assert.Equal(t, int64(100), *obj.IntField)

	// Zero values
	assert.Nil(t, obj.StringFieldNil)
	assert.Nil(t, obj.BoolFieldNil)
	assert.Nil(t, obj.IntFieldNil)

	// Nested struct
	assert.Equal(t, "nested_string", *obj.StructField.StringField)
	assert.Equal(t, true, *obj.StructField.BoolField)
	assert.Equal(t, int64(1), *obj.StructField.IntField)
}
