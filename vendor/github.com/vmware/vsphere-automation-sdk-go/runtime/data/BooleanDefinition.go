/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package data

/**
 *	 Type definition for vAPI boolean built-in type
 */
type BooleanDefinition struct{}

func NewBooleanDefinition() BooleanDefinition {
	return BooleanDefinition{}
}

func (b BooleanDefinition) Type() DataType {
	return BOOLEAN
}

func (b BooleanDefinition) ValidInstanceOf(value DataValue) bool {
	var errors = b.Validate(value)
	return len(errors) == 0
}

func (b BooleanDefinition) Validate(value DataValue) []error {
	return b.Type().Validate(value)
}
func (b BooleanDefinition) CompleteValue(value DataValue) {

}

func (b BooleanDefinition) String() string {
	return BOOLEAN.String()
}

func (b BooleanDefinition) Equals(other DataDefinition) bool {
	return other.Type() == BOOLEAN
}

func (b BooleanDefinition) NewValue(value bool) *BooleanValue {
	return NewBooleanValue(value)
}
