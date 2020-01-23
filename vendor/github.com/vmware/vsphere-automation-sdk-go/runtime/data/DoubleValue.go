/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package data

import "fmt"

type DoubleValue struct {
	value float64
}

func NewDoubleValue(value float64) *DoubleValue {
	return &DoubleValue{value: value}
}
func (doubleValue *DoubleValue) Type() DataType {
	return DOUBLE
}

func (doubleValue *DoubleValue) Value() float64 {
	return doubleValue.value
}

func (doubleValue *DoubleValue) String() string {
	return fmt.Sprintf("%g", doubleValue.value)
}
