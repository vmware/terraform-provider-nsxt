/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package data

/**
 * <code>DataValue</code> implementation for vAPI errors.
 */
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
