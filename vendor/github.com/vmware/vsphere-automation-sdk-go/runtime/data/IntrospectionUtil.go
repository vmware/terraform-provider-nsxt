/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package data

import "github.com/vmware/vsphere-automation-sdk-go/runtime/lib"

var DATA_DEFINITION = "com.vmware.vapi.std.introspection.operation.data_definition"

//this is not datadefinition interface. Its is type DataDefinition in Operation class.
//So this method takes ErrorDefinition as input and outputs DataValue for d
func ConvertOperationDataDefinitionToDataValue(dataDef DataDefinition) DataValue {
	var result = NewStructValue(DATA_DEFINITION, nil)
	var dataType = dataDef.Type().String()
	var strValue = NewStringValue(dataType)
	result.SetField("type", strValue)
	//TODO
	//refactor this
	if dataDef.Type() == STRUCTURE {
		result.SetStringField("name", (dataDef).(StructDefinition).Name())
	} else if dataDef.Type() == STRUCTURE_REF {
		result.SetStringField("name", (dataDef).(*StructRefDefinition).Name())
	} else if dataDef.Type() == ERROR {
		result.SetStringField("name", (dataDef).(ErrorDefinition).Name())
	} else {
		result.SetField("name", NewOptionalValue(nil))
	}

	if dataDef.Type() == OPTIONAL {
		var elementDef = (dataDef).(OptionalDefinition).ElementType()
		var elementValue = ConvertOperationDataDefinitionToDataValue(elementDef)
		result.SetField("element_definition", NewOptionalValue(elementValue))
	} else if dataDef.Type() == LIST {
		var elementDef = (dataDef).(ListDefinition).ElementType()
		var elementValue = ConvertOperationDataDefinitionToDataValue(elementDef)
		result.SetField("element_definition", NewOptionalValue(elementValue))
	} else {
		result.SetField("element_definition", NewOptionalValue(nil))
	}
	if dataDef.Type() == STRUCTURE {
		var fields = NewListValue()
		var structDef = (dataDef).(StructDefinition)
		for _, fieldName := range structDef.FieldNames() {
			var elementDefinition = structDef.Field(fieldName)
			var fieldPair = NewStructValue(lib.MAP_ENTRY, nil)
			var strValue = NewStringValue(fieldName)
			fieldPair.SetField(lib.MAP_KEY_FIELD, strValue)
			fieldPair.SetField(lib.MAP_VALUE_FIELD, ConvertOperationDataDefinitionToDataValue(elementDefinition))
			fields.Add(fieldPair)
		}
		result.SetField(lib.STRUCT_FIELDS, NewOptionalValue(fields))

	} else if dataDef.Type() == ERROR {
		var fields = NewListValue()
		var errorDefinition = (dataDef).(ErrorDefinition)
		for _, fieldName := range errorDefinition.FieldNames() {
			var elementDefinition = errorDefinition.Field(fieldName)
			var fieldPair = NewStructValue(lib.MAP_ENTRY, nil)
			var strValue = NewStringValue(fieldName)
			fieldPair.SetField(lib.MAP_KEY_FIELD, strValue)
			fieldPair.SetField(lib.MAP_VALUE_FIELD, ConvertOperationDataDefinitionToDataValue(elementDefinition))
			fields.Add(fieldPair)
		}
		result.SetField(lib.STRUCT_FIELDS, NewOptionalValue(fields))
	} else {
		result.SetField(lib.STRUCT_FIELDS, NewOptionalValue(nil))
	}
	return result
}
