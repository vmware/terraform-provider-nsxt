/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package data

import (
	"encoding/json"
	"strconv"
)

type IntegerValue struct {
	//Should it map to int, int8, int64?
	//this is int64 because reflection picks biggest size available.
	value int64
}

func NewIntegerValue(value int64) *IntegerValue {
	return &IntegerValue{value: value}
}

func (integerValue *IntegerValue) Type() DataType {
	return INTEGER
}

func (integerValue *IntegerValue) Value() int64 {
	return integerValue.value
}

func (integerValue *IntegerValue) String() string {
	return strconv.FormatInt(integerValue.value, 10)
}

func (integerValue *IntegerValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(integerValue.value)
}
