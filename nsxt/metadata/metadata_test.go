package metadata

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	StringField      *string
	BoolField        *bool
	IntField         *int64
	StringFieldNil   *string
	BoolFieldNil     *bool
	IntFieldNil      *int64
	StructField      *testNestedStruct
	StringListField  []string
	StructList       []testNestedStruct
	DeepNestedStruct *testDeepNestedStruct
}

type testNestedStruct struct {
	StringField *string
	BoolField   *bool
	IntField    *int64
}

type testDeepNestedStruct struct {
	StringField *string
	BoolList    []bool
	StructList  []testNestedStruct
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
	"string_list": {
		Type: schema.TypeList,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	"struct_list": {
		Type: schema.TypeList,
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
	"deep_nested_struct": {
		Type:     schema.TypeList,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"string_field": {
					Type: schema.TypeString,
				},
				"bool_list": {
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeBool,
					},
				},
				"struct_list": {
					Type: schema.TypeList,
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
			},
		},
	},
}

func basicStringSchema(sdkName string, optional bool) *ExtendedSchema {
	return &ExtendedSchema{
		Schema: schema.Schema{
			Type:     schema.TypeString,
			Optional: optional,
		},
		Metadata: Metadata{
			SchemaType:   "string",
			SdkFieldName: sdkName,
		},
	}
}

func basicBoolSchema(sdkName string, optional bool) *ExtendedSchema {
	return &ExtendedSchema{
		Schema: schema.Schema{
			Type:     schema.TypeBool,
			Optional: optional,
		},
		Metadata: Metadata{
			SchemaType:   "bool",
			SdkFieldName: sdkName,
		},
	}
}

func basicIntSchema(sdkName string, optional bool) *ExtendedSchema {
	return &ExtendedSchema{
		Schema: schema.Schema{
			Type:     schema.TypeInt,
			Optional: optional,
		},
		Metadata: Metadata{
			SchemaType:   "int",
			SdkFieldName: sdkName,
		},
	}
}

func basicStructSchema() schema.Schema {
	return schema.Schema{
		Type:     schema.TypeList,
		MaxItems: 1,
		Elem: &ExtendedResource{
			Schema: map[string]*ExtendedSchema{
				"string_field": basicStringSchema("StringField", false),
				"bool_field":   basicBoolSchema("BoolField", false),
				"int_field":    basicIntSchema("IntField", false),
			},
		},
	}
}

func mixedStructSchema() schema.Schema {
	return schema.Schema{
		Type:     schema.TypeList,
		MaxItems: 1,
		Elem: &ExtendedResource{
			Schema: map[string]*ExtendedSchema{
				"string_field": basicStringSchema("StringField", false),
				"bool_list": {
					Schema: schema.Schema{
						Type: schema.TypeList,
						Elem: basicBoolSchema("BoolList", false),
					},
					Metadata: Metadata{
						SchemaType:   "list",
						SdkFieldName: "BoolList",
					},
				},
				"struct_list": {
					Schema: basicStructSchema(),
					Metadata: Metadata{
						SchemaType:   "list",
						SdkFieldName: "StructList",
						ReflectType:  reflect.TypeOf(testNestedStruct{}),
					},
				},
			},
		},
	}
}

var testExtendedSchema = map[string]*ExtendedSchema{
	"string_field":     basicStringSchema("StringField", false),
	"bool_field":       basicBoolSchema("BoolField", false),
	"int_field":        basicIntSchema("IntField", false),
	"string_field_nil": basicStringSchema("StringFieldNil", true),
	"bool_field_nil":   basicBoolSchema("BoolFieldNil", true),
	"int_field_nil":    basicIntSchema("IntFieldNil", true),
	"struct_field": {
		Schema: basicStructSchema(),
		Metadata: Metadata{
			SchemaType:   "struct",
			SdkFieldName: "StructField",
			ReflectType:  reflect.TypeOf(testNestedStruct{}),
		},
	},
	"string_list": {
		Schema: schema.Schema{
			Type: schema.TypeList,
			Elem: basicStringSchema("StringListField", false),
		},
		Metadata: Metadata{
			SchemaType:   "list",
			SdkFieldName: "StringListField",
		},
	},
	"struct_list": {
		Schema: basicStructSchema(),
		Metadata: Metadata{
			SchemaType:   "list",
			SdkFieldName: "StructList",
			ReflectType:  reflect.TypeOf(testNestedStruct{}),
		},
	},
	"deep_nested_struct": {
		Schema: mixedStructSchema(),
		Metadata: Metadata{
			SchemaType:   "struct",
			SdkFieldName: "DeepNestedStruct",
			ReflectType:  reflect.TypeOf(testDeepNestedStruct{}),
		},
	},
}

func TestGetSchemaFromExtendedSchema(t *testing.T) {
	obs := GetSchemaFromExtendedSchema(testExtendedSchema)
	assertSchemaEqual(t, testSchema, obs)
}

func assertSchemaEqual(t *testing.T, expected, actual map[string]*schema.Schema) {
	for k, v := range actual {
		assert.Contains(t, expected, k, "unexpected schema key %s", k)
		expectedSchema := expected[k]
		assert.Equal(t, expectedSchema.Type, v.Type,
			"expected type %s, actual schema %s", expectedSchema.Type, v.Type)
		if v.Type == schema.TypeList || v.Type == schema.TypeSet {
			if _, ok := v.Elem.(*schema.Resource); ok {
				assertSchemaEqual(t,
					expectedSchema.Elem.(*schema.Resource).Schema,
					v.Elem.(*schema.Resource).Schema)
			} else if _, ok := v.Elem.(*schema.Schema); ok {
				assert.Equal(t, expectedSchema, v)
			} else {
				assert.Fail(t, fmt.Sprintf("unexpected schema element type %s", reflect.TypeOf(v.Elem)))
			}
		} else {
			assert.Equal(t, expectedSchema, v, "schema for key %s not equal", k)
		}
	}
}

func TestStructToSchema(t *testing.T) {
	testStr := "test_string"
	nestStr1, nestStr2 := "nest_str1", "nest_str2"
	testBool := true
	nestBool1, nestBool2 := true, false
	testInt := int64(123)
	nestInt1, nestInt2 := int64(123), int64(456)
	obj := testStruct{
		StringField: &testStr,
		BoolField:   &testBool,
		IntField:    &testInt,
		StructField: &testNestedStruct{
			StringField: &testStr,
			BoolField:   &testBool,
			IntField:    &testInt,
		},
		StringListField: []string{"string_1", "string_2"},
		StructList: []testNestedStruct{
			{
				StringField: &nestStr1,
				BoolField:   &nestBool1,
				IntField:    &nestInt1,
			},
			{
				StringField: &nestStr2,
				BoolField:   &nestBool2,
				IntField:    &nestInt2,
			},
		},
		DeepNestedStruct: &testDeepNestedStruct{
			StringField: &nestStr2,
			BoolList:    []bool{false, true},
			StructList: []testNestedStruct{
				{
					StringField: &nestStr1,
					BoolField:   &nestBool2,
					IntField:    &nestInt2,
				},
			},
		},
	}
	d := schema.TestResourceDataRaw(
		t, testSchema, map[string]interface{}{})

	elem := reflect.ValueOf(&obj).Elem()
	StructToSchema(elem, d, testExtendedSchema, "", nil)

	t.Run("Base types", func(t *testing.T) {
		assert.Equal(t, "test_string", d.Get("string_field").(string))
		assert.Equal(t, true, d.Get("bool_field").(bool))
		assert.Equal(t, 123, d.Get("int_field").(int))
	})

	t.Run("Zero values", func(t *testing.T) {
		_, ok := d.GetOk("string_field_nil")
		assert.False(t, ok)
		_, ok = d.GetOk("bool_field_nil")
		assert.False(t, ok)
		_, ok = d.GetOk("int_field_nil")
		assert.False(t, ok)
	})

	t.Run("Nested struct", func(t *testing.T) {
		nestedObj := d.Get("struct_field").([]interface{})[0].(map[string]interface{})
		assert.Equal(t, "test_string", nestedObj["string_field"].(string))
		assert.Equal(t, true, nestedObj["bool_field"].(bool))
		assert.Equal(t, 123, nestedObj["int_field"].(int))
	})

	t.Run("Base list", func(t *testing.T) {
		stringSlice := d.Get("string_list").([]interface{})
		assert.Equal(t, []interface{}{"string_1", "string_2"}, stringSlice)
	})

	t.Run("Struct list", func(t *testing.T) {
		structSlice := d.Get("struct_list").([]interface{})
		assert.Equal(t, 2, len(structSlice))
		assert.Equal(t, map[string]interface{}{
			"string_field": nestStr1,
			"bool_field":   nestBool1,
			"int_field":    123,
		}, structSlice[0].(map[string]interface{}))
		assert.Equal(t, map[string]interface{}{
			"string_field": nestStr2,
			"bool_field":   nestBool2,
			"int_field":    456,
		}, structSlice[1].(map[string]interface{}))
	})

	t.Run("Deep nested struct", func(t *testing.T) {
		deepNestedObj := d.Get("deep_nested_struct").([]interface{})[0].(map[string]interface{})
		assert.Equal(t, nestStr2, deepNestedObj["string_field"].(string))
		assert.Equal(t, []interface{}{false, true}, deepNestedObj["bool_list"].([]interface{}))
		nestedStructList := deepNestedObj["struct_list"].([]interface{})
		assert.Equal(t, 1, len(nestedStructList))
		assert.Equal(t, map[string]interface{}{
			"string_field": nestStr1,
			"bool_field":   nestBool2,
			"int_field":    456,
		}, nestedStructList[0].(map[string]interface{}))
	})
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
			"string_list": []interface{}{"string_1", "string_2"},
			"struct_list": []interface{}{
				map[string]interface{}{
					"string_field": "nested_string_1",
					"bool_field":   true,
					"int_field":    1,
				},
				map[string]interface{}{
					"string_field": "nested_string_2",
					"bool_field":   true,
					"int_field":    2,
				},
			},
			"deep_nested_struct": []interface{}{
				map[string]interface{}{
					"string_field": "nested_string_1",
					"bool_list":    []interface{}{true, false},
					"struct_list": []interface{}{
						map[string]interface{}{
							"string_field": "nested_string_2",
							"bool_field":   true,
							"int_field":    1,
						},
					},
				},
			},
		})

	obj := testStruct{}
	elem := reflect.ValueOf(&obj).Elem()
	SchemaToStruct(elem, d, testExtendedSchema, "", nil)

	nestStr1, nestStr2 := "nested_string_1", "nested_string_2"
	trueVal := true
	intVal1, intVal2 := int64(1), int64(2)

	t.Run("Base types", func(t *testing.T) {
		assert.Equal(t, "test_string", *obj.StringField)
		assert.Equal(t, true, *obj.BoolField)
		assert.Equal(t, int64(100), *obj.IntField)
	})

	t.Run("Zero values", func(t *testing.T) {
		assert.Equal(t, "", *obj.StringFieldNil)
		assert.Equal(t, false, *obj.BoolFieldNil)
		assert.Equal(t, int64(0), *obj.IntFieldNil)
	})

	t.Run("Nested struct", func(t *testing.T) {
		assert.Equal(t, "nested_string", *obj.StructField.StringField)
		assert.Equal(t, true, *obj.StructField.BoolField)
		assert.Equal(t, int64(1), *obj.StructField.IntField)
	})

	t.Run("Base list", func(t *testing.T) {
		assert.Equal(t, []string{"string_1", "string_2"}, obj.StringListField)
	})

	t.Run("Struct list", func(t *testing.T) {
		assert.Equal(t, 2, len(obj.StructList))
		assert.EqualValues(t, testNestedStruct{
			StringField: &nestStr1,
			BoolField:   &trueVal,
			IntField:    &intVal1,
		}, obj.StructList[0])
		assert.EqualValues(t, testNestedStruct{
			StringField: &nestStr2,
			BoolField:   &trueVal,
			IntField:    &intVal2,
		}, obj.StructList[1])
	})

	t.Run("Deep nested struct", func(t *testing.T) {
		assert.EqualValues(t, &testDeepNestedStruct{
			StringField: &nestStr1,
			BoolList:    []bool{true, false},
			StructList: []testNestedStruct{{
				StringField: &nestStr2,
				BoolField:   &trueVal,
				IntField:    &intVal1,
			}},
		}, obj.DeepNestedStruct)
	})
}
