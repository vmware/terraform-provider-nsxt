/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package data

import (
	"errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/log"
)

type StructValue struct {
	name   string
	fields map[string]DataValue
}

func NewStructValue(name string, fields map[string]DataValue) *StructValue {
	if fields == nil {
		return &StructValue{name: name, fields: make(map[string]DataValue)}
	}
	return &StructValue{name: name, fields: fields}
}

func (structValue *StructValue) Type() DataType {
	return STRUCTURE
}

func (structValue *StructValue) Name() string {
	return structValue.name
}

func (structValue *StructValue) Fields() map[string]DataValue {
	return structValue.fields
}

func (structValue *StructValue) Field(field string) (DataValue, error) {
	if dataValue, ok := structValue.fields[field]; ok {
		return dataValue, nil
	}
	return nil, errors.New("vapi.data.structure.getfield.unknown " + field)
}

func (structValue *StructValue) FieldNames() []string {
	var result = make([]string, 0)
	for key, _ := range structValue.fields {
		result = append(result, key)
	}
	return result
}

func (structValue *StructValue) List(fieldName string) (*ListValue, error) {
	var value, error = structValue.Field(fieldName)
	if error != nil {
		log.Errorf("Error getting list for %s", fieldName)
		return nil, error
	}
	if value.Type() == LIST {
		var sValue, ok = value.(*ListValue)
		if ok {
			return sValue, nil
		}
	}
	return nil, errors.New("vapi.data.structure.getfield.mismatch " + fieldName)
}

func (structValue *StructValue) Struct(fieldName string) (*StructValue, error) {
	var value, error = structValue.Field(fieldName)
	if error != nil {
		log.Errorf("Error getting struct for %s", fieldName)
		return nil, error
	}
	if value.Type() == STRUCTURE {
		var sValue, ok = value.(*StructValue)
		if ok {
			return sValue, nil
		}
	}
	return nil, errors.New("vapi.data.structure.getfield.mismatch " + fieldName)
}

func (structValue *StructValue) String(fieldName string) (string, error) {

	var result string
	var value, error = structValue.Field(fieldName)
	if error != nil {
		return result, error
	}
	if (value).Type() == STRING {
		var stringValue, ok = (value).(*StringValue)
		if ok {
			return stringValue.Value(), nil
		}
	}
	return result, errors.New("vapi.data.structure.getfield.mismatch. Expected string but got " + (value).Type().String())

}

func (structValue *StructValue) Error(fieldName string) (*ErrorValue, error) {
	var value, error = structValue.Field(fieldName)
	if error != nil {
		return nil, error
	}
	if (value).Type() == ERROR {
		var errorValue, ok = (value).(*ErrorValue)
		if ok {
			return errorValue, nil
		}
	}
	return nil, errors.New("vapi.data.structure.getfield.mismatch. Expected Error but got " + (value).Type().String())

}

func (structValue *StructValue) Optional(field string) (*OptionalValue, error) {
	var value, ok = structValue.fields[field]
	if ok {
		if value.Type() == OPTIONAL {
			return value.(*OptionalValue), nil
		}
	}
	return nil, errors.New("vapi.data.structure.getfield.mismatch " + field + "Optional" + value.Type().String())
}

func (structValue *StructValue) SetStringField(field string, value string) {
	var stringValue = NewStringValue(value)
	structValue.fields[field] = stringValue
}

func (structValue *StructValue) SetField(field string, value DataValue) {
	structValue.fields[field] = value
}

func (structValue *StructValue) HasField(field string) bool {
	var _, ok = structValue.fields[field]
	if ok {
		return true
	}
	return false
}
