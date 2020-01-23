/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package data

import (
	"errors"
)

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
