/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package data

/*
 *Type definition for vAPI secret built-in type
 */
type SecretDefinition struct{}

func (s SecretDefinition) Type() DataType {
	return SECRET
}

func NewSecretDefinition() SecretDefinition {
	var instance = SecretDefinition{}
	return instance
}

func (s SecretDefinition) ValidInstanceOf(value DataValue) bool {
	var errors = s.Validate(value)
	if len(errors) != 0 {
		return false
	}
	return true
}

func (s SecretDefinition) Validate(value DataValue) []error {
	return s.Type().Validate(value)
}

func (s SecretDefinition) CompleteValue(value DataValue) {
}

func (s SecretDefinition) String() string {
	return SECRET.String()
}

func (s SecretDefinition) NewValue(value string) *SecretValue {
	return NewSecretValue(value)
}
