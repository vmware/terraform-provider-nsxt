/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package data

import (
	"github.com/vmware/vsphere-automation-sdk-go/runtime/l10n"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/lib"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/log"
)

type StructDefinition struct {
	fields map[string]DataDefinition
	name   string
}

/**
* Constructor.
*
* @param name name for the structure; must not be <code>null</code>
* @param fields field names and definitions for the structure;
*               must not be <code>null</code>
* @throws IllegalArgumentException if some of the arguments or the
*         values in <code>fields</code> is <code>null</code>
 */

func NewStructDefinition(name string, fields map[string]DataDefinition) StructDefinition {
	if len(name) == 0 {
		log.Error("StructDefinition name missing")
	}
	if fields == nil {
		log.Error("Missing fields for StructDefinition")
	}
	for key, value := range fields {
		if value == nil {
			log.Errorf("Missing value for field %s", key)
			return StructDefinition{}
		}
	}
	return StructDefinition{name: name, fields: fields}
}

func (structDefinition StructDefinition) Type() DataType {
	return STRUCTURE
}

func (structDefinition StructDefinition) Name() string {
	return (structDefinition.name)
}

func (structDefinition StructDefinition) FieldNames() []string {
	var keys = make([]string, len(structDefinition.fields))
	var i = 0
	for k := range structDefinition.fields {
		keys[i] = k
		i++
	}
	return keys
}

func (structDefinition StructDefinition) Field(field string) DataDefinition {
	//TODO:
	// zero value of pointer is nil
	// what is the zero value of a structure
	var value, ok = structDefinition.fields[field]
	if ok {
		return value
	} else {
		return nil
	}

}

func (structDefinition StructDefinition) HasField(field string) bool {
	var _, ok = structDefinition.fields[field]
	if ok {
		return true
	}
	return false
}

func (structDefinition StructDefinition) Validate(value DataValue) []error {
	var result = structDefinition.Type().Validate(value)
	if result != nil {
		return result
	}
	var structValue = value.(*StructValue)
	if structDefinition.Name() != structValue.Name() {
		if !(structValue.Name() == lib.MAP_ENTRY || structValue.Name() == lib.OPERATION_INPUT) {
			var args = map[string]string{
				"actualName":   structValue.Name(),
				"expectedName": structDefinition.Name()}
			return []error{l10n.NewRuntimeError("vapi.data.structure.name.mismatch", args)}
		}
	}
	// Make sure no fields are missing and those that are present are valid
	for _, fieldName := range structDefinition.FieldNames() {
		var fieldDef = structDefinition.Field(fieldName)
		var fieldVal, err = structValue.Field(fieldName)
		if err != nil && fieldDef.Type() != OPTIONAL {
			var args = map[string]string{
				"structName": structDefinition.Name(),
				"fieldName":  fieldName}
			return []error{l10n.NewRuntimeError("vapi.data.structure.field.missing", args)}
		}

		var errors = fieldDef.Validate(fieldVal)
		if errors != nil {
			var args = map[string]string{
				"structName": structDefinition.Name(),
				"fieldName":  fieldName}
			var msg = l10n.NewRuntimeError("vapi.data.structure.field.invalid", args)
			errors = append(errors, msg)
			return errors
		}
	}
	return result
}

func (structDefinition StructDefinition) CompleteValue(value DataValue) {
	if value != nil && value.Type() == structDefinition.Type() {
		var structValue = value.(*StructValue)

		for _, fieldName := range structDefinition.FieldNames() {
			if !structValue.HasField(fieldName) {
				var fieldDef = structDefinition.Field(fieldName)
				if fieldDef.Type() == OPTIONAL {
					structValue.SetField(fieldName, NewOptionalValue(nil))
				}
			} else {
				var fieldDef = structDefinition.Field(fieldName)
				var fieldValue, error = structValue.Field(fieldName)
				if error != nil {
					log.Error(error)
				} else {
					fieldDef.CompleteValue(fieldValue)
				}
			}
		}
	}

}

func (structDefinition StructDefinition) ValidInstanceOf(value DataValue) bool {
	var errors = structDefinition.Validate(value)
	if len(errors) != 0 {
		return false
	}
	return true
}

func (structDefinition StructDefinition) String() string {
	return STRUCTURE.String()
}

//TODO: sreeshas
// lot of structdef tests depend on this method. remove it after correcting tests
func (structDefinition StructDefinition) Equals(other DataDefinition) bool {
	if other == nil {
		return false
	}
	if other.Type() != STRUCTURE {
		return false
	}
	var otherStruct = other.(StructDefinition)
	if (structDefinition.name) != (otherStruct.name) {
		return false
	}

	for key, _ := range structDefinition.fields {
		if _, ok := otherStruct.fields[key]; !ok {
			return false
		}
	}
	return true

}

func (structDefinition StructDefinition) NewValue() *StructValue {
	return NewStructValue(structDefinition.name, nil)

}
