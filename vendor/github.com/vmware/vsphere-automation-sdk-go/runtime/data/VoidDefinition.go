/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package data

type VoidDefinition struct{}

func NewVoidDefinition() VoidDefinition {
	return VoidDefinition{}
}

func (voidDefinition VoidDefinition) Type() DataType {
	return VOID
}

func (voidDefinition VoidDefinition) ValidInstanceOf(value DataValue) bool {
	var errors = voidDefinition.Validate(value)
	if len(errors) != 0 {
		return false
	}
	return true
}

func (voidDefinition VoidDefinition) Validate(value DataValue) []error {
	return voidDefinition.Type().Validate(value)
}

func (voidDefinition VoidDefinition) CompleteValue(value DataValue) {

}

func (voidDefinition VoidDefinition) String() string {
	return VOID.String()
}

func (voidDefinition VoidDefinition) NewValue() *VoidValue {
	return NewVoidValue()
}
