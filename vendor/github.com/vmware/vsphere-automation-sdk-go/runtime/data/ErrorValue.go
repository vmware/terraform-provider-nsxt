/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package data

import "errors"

/**
 * <code>DataValue</code> implementation for vAPI errors.
 */
type ErrorValue struct {
	name   string
	fields map[string]DataValue
}

func NewErrorValue(name string, fields map[string]DataValue) *ErrorValue {
	if fields == nil {
		return &ErrorValue{name: name, fields: make(map[string]DataValue)}
	}
	return &ErrorValue{name: name, fields: fields}
}

func (errorValue *ErrorValue) Type() DataType {
	return ERROR
}

func (errorValue *ErrorValue) Name() string {
	return errorValue.name
}

func (errorValue *ErrorValue) SetField(field string, value DataValue) {
	errorValue.fields[field] = value
}

func (errorValue *ErrorValue) Field(field string) (DataValue, error) {
	if dataValue, ok := errorValue.fields[field]; ok {
		return dataValue, nil
	}
	return nil, errors.New("vapi.data.error.getfield.unknown" + field)
}
func (errorValue *ErrorValue) Fields() map[string]DataValue {
	return errorValue.fields
}

func (errorValue *ErrorValue) HasField(field string) bool {
	if _, ok := errorValue.fields[field]; ok {
		return true
	}
	return false
}
