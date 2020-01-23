/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package bindings

import (
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/l10n"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/log"
)

type Validator interface {
	// Validates a struct value
	Validate(structValue *data.StructValue) []error
}

type FieldData struct {
	name       string
	isRequired bool
}

func NewFieldData(fieldName string, isRequired bool) FieldData {
	return FieldData{fieldName, isRequired}
}

type UnionValidator struct {
	//Name of a structure field that represents a union discriminant/tag
	tagName string
	// For each value tag, this map contains the field data for the fields
	// associated with the tag
	// Key is tag value and value is a list of field names and whether or
	// not they are required
	caseMap map[string][]FieldData
	// Name of all fields in the values of caseMap
	allCaseFields []string
}

func NewUnionValidator(tagName string, caseMap map[string][]FieldData) UnionValidator {
	var allFields []string
	for _, fields := range caseMap {
		for _, f := range fields {
			allFields = append(allFields, f.name)
		}
	}
	return UnionValidator{tagName, caseMap, allFields}
}

func (uv UnionValidator) Validate(structVal *data.StructValue) []error {
	if structVal == nil {
		return nil
	}
	var tag string

	// Retrieving value of tag. It is ok for tag field to not be set. So, err is ignored
	if tagField, err := structVal.Field(uv.tagName); err == nil {
		var tagDataVal data.DataValue
		optionalUnset := false
		if tagVal, ok := tagField.(*data.OptionalValue); ok {
			if !tagVal.IsSet() {
				optionalUnset = true
			}
			tagDataVal = tagVal.Value()
		} else {
			tagDataVal = tagField
		}
		if tagVal, ok := tagDataVal.(*data.StringValue); ok {
			tag = tagVal.Value()
		} else if !optionalUnset {
			log.Errorf("Expected StringValue for tag field named %s. But got %s value", uv.tagName, tagDataVal.Type())
			return nil
		}
	}

	// Verify that the fields associated with the tag are present
	allowedFields := map[string]bool{}
	for _, fieldData := range uv.caseMap[tag] {
		allowedFields[fieldData.name] = true
		field, notPresentErr := structVal.Field(fieldData.name)
		requiredNotPresent := false
		if notPresentErr != nil && fieldData.isRequired {
			requiredNotPresent = true
		}

		requiredNotSet := false
		if optionalField, ok := field.(*data.OptionalValue); ok {
			if !optionalField.IsSet() && fieldData.isRequired {
				requiredNotSet = true
			}
		}
		if requiredNotPresent || requiredNotSet {
			var args = map[string]string{
				"structName": structVal.Name(),
				"fieldName":  fieldData.name}
			return []error{l10n.NewRuntimeError("vapi.data.structure.union.missing", args)}
		}
	}

	// Case fields that are not associated with the current tag value
	// should not be present/set
	prohibitedFields := []string{}
	for _, field := range uv.allCaseFields {
		if !allowedFields[field] {
			prohibitedFields = append(prohibitedFields, field)
		}
	}

	for _, fieldName := range prohibitedFields {
		fieldIsSet := false
		if field, _ := structVal.Field(fieldName); field != nil {
			if optionalField, ok := field.(*data.OptionalValue); ok {
				if optionalField.IsSet() {
					fieldIsSet = true
				}
			} else {
				fieldIsSet = true
			}
		}
		if fieldIsSet {
			var args = map[string]string{
				"structName": structVal.Name(),
				"fieldName":  fieldName}
			return []error{l10n.NewRuntimeError("vapi.data.structure.union.extra", args)}
		}
	}
	return nil
}

// HasFieldsOfValidator validator class that validates the data_value has
// required fields of the class specified
type HasFieldsOfValidator struct {
	hasFieldsOfTypes []ReferenceType
	mode ConverterMode
}

func NewHasFieldsOfValidator(hasFieldsOfTypes []ReferenceType, mode ConverterMode) HasFieldsOfValidator {
	return HasFieldsOfValidator{hasFieldsOfTypes: hasFieldsOfTypes, mode: mode}
}

//Validates whether a StructValue satisfies the HasFieldsOf constraint
func (hv HasFieldsOfValidator) Validate(structValue *data.StructValue) []error {
	if structValue == nil {
		return nil
	}

	if hv.hasFieldsOfTypes == nil || len(hv.hasFieldsOfTypes) == 0 {
		// HasFieldsOfValidator is only for dynamic structures that have hasFieldsOfTypes
		return nil
	}

	var errorMsgs []error
	for _, hasTypeRef := range hv.hasFieldsOfTypes {
		hasBindingType := hasTypeRef.Resolve().(StructType)
		converter := NewTypeConverter()
		converter.SetMode(hv.mode)
		_, err := converter.ConvertToGolang(structValue, hasBindingType)
		if err != nil {
			msg := l10n.NewRuntimeError("vapi.data.structure.dynamic.invalid",
				map[string]string{"structName": hasBindingType.Name()})
			err = append(err, msg)
			errorMsgs = append(errorMsgs, err...)
		}
	}
	return errorMsgs
}

// IsOneOfValidator validates a single is-one-of constraint
// on a single field of a structure.
type IsOneOfValidator struct {
	fieldName     string
	allowedValues map[string]bool
}

// Instantiates an IsOneOfValidator with the canonical name of the
// field that should be validated and the allowed values for the field
func NewIsOneOfValidator(fieldName string, allowedValues []string) IsOneOfValidator {
	allowedValuesSet := map[string]bool{}
	for _, val := range allowedValues {
		allowedValuesSet[val] = true
	}
	return IsOneOfValidator{fieldName: fieldName, allowedValues: allowedValuesSet}
}

//Validates whether a StructValue satisfies the HasFieldsOf constraint
func (is IsOneOfValidator) Validate(structValue *data.StructValue) []error {
	// @IsOneOf is only allowed on String fields, and not on Optional<String>
	// ones (this is enforced by the idl toolkit. So it is safe to assume the
	// field must be present
	if structValue == nil || !structValue.HasField(is.fieldName) {
		return nil
	}

	fieldValue, _ := structValue.Field(is.fieldName)
	if valStr, ok := fieldValue.(*data.StringValue); ok {
		if !is.allowedValues[valStr.Value()] {
			msg := l10n.NewRuntimeError("vapi.data.structure.isoneof.value.invalid",
				map[string]string{"value": valStr.Value(), "fieldName": is.fieldName})
			return []error{msg}
		}
	}
	return nil
}
