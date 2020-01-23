/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package data

/**
 *
 * An OptionalDefinition instance defines an optional type with
 * a specific element type.
 *
 */
type OptionalDefinition struct {
	//cannot be nil
	optionalElementType DataDefinition
}

func NewOptionalDefinition(optionalElementType DataDefinition) OptionalDefinition {
	return OptionalDefinition{optionalElementType: optionalElementType}
}

func (optionalDefinition OptionalDefinition) Type() DataType {
	return OPTIONAL
}

func (optionalDefinition OptionalDefinition) ElementType() DataDefinition {
	return optionalDefinition.optionalElementType
}

func (optionalDefinition OptionalDefinition) Validate(value DataValue) []error {
	var result = optionalDefinition.Type().Validate(value)
	if result != nil {
		// Error is returned for nil value. But optional value can be nil in permissive mode?
		return result
	}
	var optionalValue = value.(*OptionalValue)

	if optionalValue.IsSet() {
		var subErrors = optionalDefinition.ElementType().Validate(optionalValue.Value())
		if subErrors != nil {
			result = append(result, subErrors...)
		}
	} else {
		// nulls are valid for optional so this is a
		// no-op for a good reason
	}
	return result
}

func (optionalDefinition OptionalDefinition) CompleteValue(value DataValue) {
	if value != nil && value.Type() == OPTIONAL {
		var optValue = value.(*OptionalValue)
		optionalDefinition.ElementType().CompleteValue(optValue.Value())
	}
}

func (optionalDefinition OptionalDefinition) String() string {
	//TODO:sreeshas
	// implement this
	return ""
}

func (optionalDefinition OptionalDefinition) ValidInstanceOf(value DataValue) bool {
	var errors = optionalDefinition.Validate(value)
	if len(errors) != 0 {
		return false
	}
	return true
}

func (optionalDefinition OptionalDefinition) NewValue(value DataValue) *OptionalValue {
	return NewOptionalValue(value)
}
