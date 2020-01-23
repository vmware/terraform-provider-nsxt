/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package bindings

import (
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
)

type StructValueExtractor struct {
	structType    StructType
	structValue   *data.StructValue
	typeConverter *TypeConverter
}

func NewStructValueExtractor(structType StructType, structValue *data.StructValue, typeConverter *TypeConverter) *StructValueExtractor {
	return &StructValueExtractor{structType: structType, structValue: structValue, typeConverter: typeConverter}
}

func (s *StructValueExtractor) ExtractValue(fieldName string) (interface{}, []error) {
	fieldValue, _ := s.structValue.Field(fieldName)
	fieldType := s.structType.Field(fieldName)
	return s.typeConverter.ConvertToGolang(fieldValue, fieldType)
}
