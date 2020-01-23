/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package data

/**
 * Type definition for vAPI integer built-in type
 */
type IntegerDefinition struct{}

func NewIntegerDefinition() IntegerDefinition {
	return IntegerDefinition{}
}

func (integerDefinition IntegerDefinition) Type() DataType {
	return INTEGER
}

func (integerDefinition IntegerDefinition) ValidInstanceOf(value DataValue) bool {
	var errors = integerDefinition.Validate(value)
	if len(errors) != 0 {
		return false
	}
	return true
}
func (i IntegerDefinition) Validate(value DataValue) []error {
	return i.Type().Validate(value)
}

func (integerDefinition IntegerDefinition) CompleteValue(value DataValue) {

}

func (integerDefinition IntegerDefinition) String() string {
	return INTEGER.String()
}
