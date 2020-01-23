/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package info


type UserDefinedTypeInfo struct {
	name          string
	structureInfo *StructureInfo
}

func NewUserDefinedTypeInfo(name string, structureInfo *StructureInfo) *UserDefinedTypeInfo {
	return &UserDefinedTypeInfo{name: name, structureInfo: structureInfo}
}

// Name
func (udt *UserDefinedTypeInfo) Name() string {
	return udt.name
}

// StructureInfo
func (udt *UserDefinedTypeInfo) StructureInfo() *StructureInfo {
	return udt.structureInfo
}
