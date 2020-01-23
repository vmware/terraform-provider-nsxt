/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package data


type StringDefinition struct{}

func NewStringDefinition() StringDefinition {
	var instance = StringDefinition{}
	return instance
}

func (s StringDefinition) Type() DataType {
	return STRING
}

func (s StringDefinition) ValidInstanceOf(value DataValue) bool {
	var errors = s.Validate(value)
	if len(errors) != 0 {
		return false
	}
	return true
}

func (s StringDefinition) Validate(value DataValue) []error {
	return s.Type().Validate(value)
}

func (s StringDefinition) CompleteValue(value DataValue) {

}

func (s StringDefinition) String() string {
	return STRING.String()
}

func (s StringDefinition) NewValue(value string) *StringValue {
	return NewStringValue(value)
}
