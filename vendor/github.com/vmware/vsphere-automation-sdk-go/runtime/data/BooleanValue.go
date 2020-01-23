/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package data

type BooleanValue struct {
	value bool
}

func NewBooleanValue(value bool) *BooleanValue {
	return &BooleanValue{value: value}
}

func (b *BooleanValue) Type() DataType {
	return BOOLEAN
}

func (b *BooleanValue) Value() bool {
	return (b.value)
}
