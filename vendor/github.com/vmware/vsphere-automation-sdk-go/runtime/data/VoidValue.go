/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package data

type VoidValue struct {
}

func NewVoidValue() *VoidValue {
	return &VoidValue{}
}

func (voidValue *VoidValue) Value() DataValue {
	return nil
}

func (voidValue *VoidValue) Type() DataType {
	return VOID
}
