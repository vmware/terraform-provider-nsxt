/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package data


type BlobDefinition struct{}

func NewBlobDefinition() BlobDefinition {
	var instance = BlobDefinition{}
	return instance
}

func (b BlobDefinition) Type() DataType {
	return BLOB
}

func (b BlobDefinition) ValidInstanceOf(value DataValue) bool {
	var errors = b.Validate(value)
	return len(errors) == 0
}

func (b BlobDefinition) Validate(value DataValue) []error {
	return b.Type().Validate(value)
}

func (b BlobDefinition) CompleteValue(value DataValue) {

}

func (b BlobDefinition) String() string {
	return BLOB.String()
}

func (b BlobDefinition) NewValue(value []byte) *BlobValue {
	return NewBlobValue(value)
}
