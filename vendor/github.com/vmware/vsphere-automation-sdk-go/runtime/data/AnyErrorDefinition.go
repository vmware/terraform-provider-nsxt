/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package data

type AnyErrorDefinition struct {
}

func NewAnyErrorDefinition() AnyErrorDefinition {
	return AnyErrorDefinition{}
}

func (a AnyErrorDefinition) ValidInstanceOf(value DataValue) bool {
	var errors = a.Validate(value)
	return len(errors) == 0
}

func (a AnyErrorDefinition) CompleteValue(value DataValue) {
	// AnyError has no fields, so it is a NO-OP
}

func (a AnyErrorDefinition) String() string {
	return a.Type().String()
}

func (a AnyErrorDefinition) Type() DataType {
	return ANY_ERROR
}

func (a AnyErrorDefinition) Validate(value DataValue) []error {
	if value != nil && value.Type() == ERROR {
		return nil
	}
	return a.Type().Validate(value)
}
