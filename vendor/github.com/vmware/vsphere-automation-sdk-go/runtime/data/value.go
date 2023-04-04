/* Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/log"
	"strconv"
)

type DataValue interface {
	Type() DataType
}

type BlobValue struct {
	value []byte
}

func NewBlobValue(value []byte) *BlobValue {
	return &BlobValue{value: value}
}

func (b *BlobValue) Type() DataType {
	return BLOB
}

func (b *BlobValue) Value() []byte {
	return b.value
}

type BooleanValue struct {
	value bool
}

func NewBooleanValue(value bool) *BooleanValue {
	return &BooleanValue{value: value}
}

func (b *BooleanValue) Type() DataType {
	return BOOLEAN
}

func (b *BooleanValue) Value() bool {
	return b.value
}

type DoubleValue struct {
	value float64
}

func NewDoubleValue(value float64) *DoubleValue {
	return &DoubleValue{value: value}
}
func (doubleValue *DoubleValue) Type() DataType {
	return DOUBLE
}

func (doubleValue *DoubleValue) Value() float64 {
	return doubleValue.value
}

func (doubleValue *DoubleValue) String() string {
	return fmt.Sprintf("%g", doubleValue.value)
}

type ErrorValue struct {
	StructValue
}

func NewErrorValue(name string, fields map[string]DataValue) *ErrorValue {
	if fields == nil {
		return &ErrorValue{StructValue{name: name, fields: make(map[string]DataValue)}}
	}
	return &ErrorValue{StructValue{name: name, fields: fields}}
}

func (errorValue *ErrorValue) Type() DataType {
	return ERROR
}

type IntegerValue struct {
	//Should it map to int, int8, int64?
	//this is int64 because reflection picks biggest size available.
	value int64
}

func NewIntegerValue(value int64) *IntegerValue {
	return &IntegerValue{value: value}
}

func (integerValue *IntegerValue) Type() DataType {
	return INTEGER
}

func (integerValue *IntegerValue) Value() int64 {
	return integerValue.value
}

func (integerValue *IntegerValue) String() string {
	return strconv.FormatInt(integerValue.value, 10)
}

func (integerValue *IntegerValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(integerValue.value)
}

type ListValue struct {
	list  []DataValue
	isMap bool
}

func NewListValue() *ListValue {
	var list = make([]DataValue, 0)
	return &ListValue{list: list}
}

func (listValue *ListValue) Type() DataType {
	return LIST
}

func (listValue *ListValue) Add(value DataValue) {
	listValue.list = append(listValue.list, value)
}

func (listValue *ListValue) Get(index int) DataValue {
	return listValue.list[index]
}

func (listValue *ListValue) IsEmpty() bool {
	return len(listValue.list) == 0
}

func (listValue *ListValue) List() []DataValue {
	return listValue.list
}

// MarkAsMap is used with the rest protocol to differentiate
// maps converted to listValues and normal lists
func (listValue *ListValue) MarkAsMap() {
	listValue.isMap = true
}

// IsMap is used with the REST protocol to determine whether the
// list value is a presentation of a map or of a list.
func (listValue *ListValue) IsMap() bool {
	return listValue.isMap
}

type OptionalValue struct {
	value DataValue
}

func NewOptionalValue(value DataValue) *OptionalValue {
	return &OptionalValue{value: value}
}

func NewOptionalValueString(value string) *OptionalValue {
	var stringValue = NewStringValue(value)
	return NewOptionalValue(stringValue)
}

func (optionalValue *OptionalValue) Type() DataType {
	return OPTIONAL
}

func (optionalValue *OptionalValue) IsSet() bool {
	return optionalValue.value != nil
}

func (optionalValue *OptionalValue) Value() DataValue {
	return optionalValue.value
}

func (optionalValue *OptionalValue) String() (string, error) {
	var value = optionalValue.Value()
	if optionalValue.IsSet() {
		if value.Type() == STRING {
			var stringValue = value.(*StringValue)
			return stringValue.Value(), nil
		} else {
			return "", errors.New("vapi.data.optional.getvalue.mismatch")
		}
	} else {
		return "", errors.New("vapi.data.optional.getvalue.unset")
	}

}

func (optionalValue *OptionalValue) Struct() (*StructValue, error) {
	var value = optionalValue.Value()
	if optionalValue.IsSet() {
		if value.Type() == STRUCTURE {
			var structValue = value.(*StructValue)
			return structValue, nil
		} else {
			return nil, errors.New("vapi.data.optional.getvalue.mismatch")
		}
	} else {
		return nil, errors.New("vapi.data.optional.getvalue.unset")
	}
}

/**
 * Value of type secret, which is intended to represent sensitive
 * information, like passwords.
 *
 *
 * In addition the actual content will not be returned by the {@link #String()}
 * as a precaution for avoiding accidental displaying or logging it.
 */
type SecretValue struct {
	value string
}

func NewSecretValue(value string) *SecretValue {
	return &SecretValue{value: value}
}

func (s *SecretValue) Type() DataType {
	return SECRET
}

func (s *SecretValue) Value() string {
	return s.value
}

type StringValue struct {
	value string
}

func NewStringValue(value string) *StringValue {
	var stringValue = &StringValue{value: value}
	return stringValue
}

func (stringValue *StringValue) Type() DataType {
	return STRING
}

func (stringValue *StringValue) Value() string {
	return stringValue.value
}

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
	for key := range structValue.fields {
		result = append(result, key)
	}
	return result
}

func (structValue *StructValue) List(fieldName string) (*ListValue, error) {
	var value, err = structValue.Field(fieldName)
	if err != nil {
		log.Errorf("Error getting list for %s", fieldName)
		return nil, err
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
	var value, err = structValue.Field(fieldName)
	if err != nil {
		log.Errorf("Error getting struct for %s", fieldName)
		return nil, err
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
	var value, err = structValue.Field(fieldName)
	if err != nil {
		return result, err
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
	var value, err = structValue.Field(fieldName)
	if err != nil {
		return nil, err
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

type VoidValue struct {
}

func NewVoidValue() *VoidValue {
	return &VoidValue{}
}

func (voidValue *VoidValue) Value() DataValue {
	return nil
}

func (voidValue *VoidValue) Type() DataType {
	return VOID
}
