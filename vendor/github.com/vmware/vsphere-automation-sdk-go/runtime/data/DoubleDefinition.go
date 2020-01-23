/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package data

/**
 * Type definition for vAPI double built-in type
 */
type DoubleDefinition struct{}

func NewDoubleDefinition() DoubleDefinition {
	return DoubleDefinition{}
}
func (d DoubleDefinition) Type() DataType {
	return DOUBLE
}

func (d DoubleDefinition) ValidInstanceOf(value DataValue) bool {
	var errors = d.Validate(value)
	return len(errors) == 0
}
func (d DoubleDefinition) Validate(value DataValue) []error {
	return d.Type().Validate(value)
}
func (d DoubleDefinition) CompleteValue(value DataValue) {

}
func (d DoubleDefinition) String() string {
	return DOUBLE.String()
}
