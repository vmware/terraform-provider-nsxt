/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package info


type StructureInfo struct {
	identifier      string
	isError         bool
	documentation   string
	privileges      []string
	fieldinfo       map[string]*FieldInfo
	enumerationInfo map[string]*EnumerationInfo
	unionInfo       map[string]*UnionInfo
}

func NewStructureInfo() *StructureInfo {
	fieldinfo := make(map[string]*FieldInfo)
	enumerationInfo := make(map[string]*EnumerationInfo)
	unionInfo := make(map[string]*UnionInfo)
	return &StructureInfo{fieldinfo: fieldinfo, enumerationInfo: enumerationInfo, unionInfo: unionInfo}
}

// Identifier
func (s *StructureInfo) Identifier() string {
	return s.identifier
}

func (s *StructureInfo) SetIdentifier(identifier string) {
	s.identifier = identifier
}

// IsError
func (s *StructureInfo) IsError() bool {
	return s.isError
}

func (s *StructureInfo) SetIsError(isError bool) {
	s.isError = isError
}

// Fieldinfo
func (s *StructureInfo) FieldInfo(field string) *FieldInfo {
	return s.fieldinfo[field]
}

func (s *StructureInfo) FieldInfoMap() map[string]*FieldInfo {
	return s.fieldinfo
}

func (s *StructureInfo) SetFieldInfo(field string, fieldInfo *FieldInfo) {
	s.fieldinfo[field] = fieldInfo
}

func (s *StructureInfo) ListFieldInfo() []string {
	keys := []string{}
	for k, _ := range s.fieldinfo {
		keys = append(keys, k)
	}
	return keys
}
