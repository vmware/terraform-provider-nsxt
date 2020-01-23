/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package data

/**
 * Value of type secret, which is intended to represent sensitive
 * information, like passwords.
 *
 *
 * In addition the actual content will not be returned by the {@link #String()}
 * as a precaution for avoiding accidental displaying or logging it.
 */
type SecretValue struct {
	value string
}

func NewSecretValue(value string) *SecretValue {
	return &SecretValue{value: value}
}

func (s *SecretValue) Type() DataType {
	return SECRET
}

func (s *SecretValue) Value() string {
	return s.value
}
