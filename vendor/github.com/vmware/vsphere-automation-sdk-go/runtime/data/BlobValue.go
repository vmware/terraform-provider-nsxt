/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package data

type BlobValue struct {
	value []byte
}

func NewBlobValue(value []byte) *BlobValue {
	return &BlobValue{value: value}
}

func (b *BlobValue) Type() DataType {
	return BLOB
}

func (b *BlobValue) Value() []byte {
	return (b.value)
}
