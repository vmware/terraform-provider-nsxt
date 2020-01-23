/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package data

/**
 * DataDefinition for dynamic structures.
 */
type DynamicStructDefinition struct {
	validDataTypes []DataType
}

func NewDynamicStructDefinition() DynamicStructDefinition {
	var validDataTypes = []DataType{STRUCTURE}
	return DynamicStructDefinition{validDataTypes: validDataTypes}
}

func (dynamicStructure DynamicStructDefinition) Type() DataType {
	return DYNAMIC_STRUCTURE
}

func (dynamicStructure DynamicStructDefinition) ValidInstanceOf(value DataValue) bool {
	var errors = dynamicStructure.Validate(value)
	return len(errors) == 0
}
func (d DynamicStructDefinition) Validate(value DataValue) []error {
	if value != nil && value.Type() == STRUCTURE {
		return nil
	}
	return d.Type().Validate(value)
}
func (dynamicStructure DynamicStructDefinition) CompleteValue(value DataValue) {
	//do nothing
}
func (dynamicStructure DynamicStructDefinition) String() string {
	return dynamicStructure.Type().String()
}
