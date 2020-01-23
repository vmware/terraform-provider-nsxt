/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package data

import (
	"github.com/vmware/vsphere-automation-sdk-go/runtime/l10n"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/log"
)

/**
 * Type definition for vAPI error type
 */
type ErrorDefinition struct {
	name   string
	fields map[string]DataDefinition
}

func NewErrorDefinition(name string, fields map[string]DataDefinition) ErrorDefinition {
	if len(name) == 0 {
		log.Error("ErrorDefinition name missing")
	}
	if fields == nil {
		log.Error("Missing fields for ErrorDefinition")
	}
	for key, value := range fields {
		if value == nil {
			log.Error("Missing value for field " + key)
		}
	}
	return ErrorDefinition{name: name, fields: fields}
}

func (errorDefinition ErrorDefinition) Type() DataType {
	return ERROR
}

func (errorDefinition ErrorDefinition) Name() string {
	return errorDefinition.name
}

func (errorDefinition ErrorDefinition) FieldNames() []string {
	var keys = make([]string, len(errorDefinition.fields))
	var i = 0
	for k := range errorDefinition.fields {
		keys[i] = k
		i++
	}
	return keys
}

func (errorDefinition ErrorDefinition) Field(field string) DataDefinition {
	//TODO:
	// zero value of pointer is nil
	// what is the zero value of a structure
	var value, ok = errorDefinition.fields[field]
	if ok {
		return value
	} else {
		return nil
	}

}

func (errorDefinition ErrorDefinition) HasField(field string) bool {
	var _, ok = errorDefinition.fields[field]
	return ok
}

func (errorDefinition ErrorDefinition) Validate(value DataValue) []error {
	var result = errorDefinition.Type().Validate(value)
	if result != nil {
		return result
	}
	var errorValue = value.(*ErrorValue)
	if errorDefinition.Name() != errorValue.Name() {
		var args = map[string]string{
			"actualName":   errorValue.Name(),
			"expectedName": errorDefinition.Name()}
		return []error{l10n.NewRuntimeError("vapi.data.structure.name.mismatch", args)}
	}
	// Make sure no fields are missing and those that are present are valid
	for _, fieldName := range errorDefinition.FieldNames() {
		var fieldDef = errorDefinition.Field(fieldName)
		var fieldVal, err = errorValue.Field(fieldName)
		if err != nil {
			var args = map[string]string{
				"structName": errorDefinition.Name(),
				"fieldName":  fieldName}
			return []error{l10n.NewRuntimeError("vapi.data.structure.field.missing", args)}
		}

		var errors = fieldDef.Validate(fieldVal)
		if errors != nil {
			var args = map[string]string{
				"structName": errorDefinition.Name(),
				"fieldName":  fieldName}
			var msg = l10n.NewRuntimeError("vapi.data.structure.field.invalid", args)
			errors = append(errors, msg)
			return errors
		}
	}
	return nil
}

func (errorDefinition ErrorDefinition) CompleteValue(value DataValue) {
	if value != nil && value.Type() == errorDefinition.Type() {
		var errorValue = value.(*ErrorValue)

		for _, fieldName := range errorDefinition.FieldNames() {
			if !errorValue.HasField(fieldName) {
				var fieldDef = errorDefinition.Field(fieldName)
				if fieldDef.Type() == OPTIONAL {
					errorValue.SetField(fieldName, NewOptionalValue(nil))
				}
			} else {
				var fieldDef = errorDefinition.Field(fieldName)
				//TODO
				// we already check if the field is present in hasfield.
				// so handling this error is unnecessary. we should change the code block and remove call to HasField
				var fieldValue, _ = errorValue.Field(fieldName)
				fieldDef.CompleteValue(fieldValue)
			}
		}
	}

}

func (errorDefinition ErrorDefinition) ValidInstanceOf(value DataValue) bool {
	var errors = errorDefinition.Validate(value)
	if len(errors) != 0 {
		return false
	}
	return true
}

func (errorDefinition ErrorDefinition) String() string {
	return errorDefinition.Type().String()
}

func (errorDefinition ErrorDefinition) NewValue() *ErrorValue {
	return NewErrorValue(errorDefinition.name, nil)
}
