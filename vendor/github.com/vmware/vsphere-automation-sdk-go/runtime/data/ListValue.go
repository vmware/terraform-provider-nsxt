/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package data

type ListValue struct {
	list []DataValue
}

func NewListValue() *ListValue {
	var list = make([]DataValue, 0)
	return &ListValue{list: list}
}

func (listValue *ListValue) Type() DataType {
	return LIST
}

func (listValue *ListValue) Add(value DataValue) {
	listValue.list = append(listValue.list, value)
}

func (listValue *ListValue) Get(index int) DataValue {
	return listValue.list[index]
}

func (listValue *ListValue) IsEmpty() bool {
	return len(listValue.list) == 0
}

func (listValue *ListValue) List() []DataValue {
	return listValue.list
}
