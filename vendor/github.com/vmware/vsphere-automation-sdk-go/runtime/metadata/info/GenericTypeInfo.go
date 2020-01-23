/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package info


type GenericTypeInfo struct {
	name         string
	isList       bool
	isOptional   bool
	isSet        bool
	isMap        bool
	elementType  TypeInfo
	mapKeyType   TypeInfo
	mapValueType TypeInfo
}

func NewGenericTypeInfo(name string) *GenericTypeInfo {
	genericTypeInfo := GenericTypeInfo{name: name, isList: false, isOptional: false, isSet: false, isMap: false}
	return &genericTypeInfo
}

// Name
func (gt *GenericTypeInfo) Name() string {
	return gt.name
}

func (gt *GenericTypeInfo) SetName(name string) {
	gt.name = name
}

// Is list
func (gt *GenericTypeInfo) IsList() bool {
	return gt.isList
}

func (gt *GenericTypeInfo) SetIsList(isList bool) {
	gt.isList = isList
}

// Is Optional
func (gt *GenericTypeInfo) IsOptional() bool {
	return gt.isOptional
}

func (gt *GenericTypeInfo) SetIsOptional(isOptional bool) {
	gt.isOptional = isOptional
}

// Is Set
func (gt *GenericTypeInfo) IsSet() bool {
	return gt.isSet
}

func (gt *GenericTypeInfo) SetIsSet(isSet bool) {
	gt.isSet = isSet
}

// Is Map
func (gt *GenericTypeInfo) IsMap() bool {
	return gt.isMap
}

func (gt *GenericTypeInfo) SetIsMap(isMap bool) {
	gt.isMap = isMap
}

// Element Type
func (gt *GenericTypeInfo) SetElementType(elementType TypeInfo) {
	gt.elementType = elementType
}

func (gt *GenericTypeInfo) ElementType() TypeInfo {
	return gt.elementType
}

// Map key type
func (gt *GenericTypeInfo) SetMapKeyType(mapKeyType TypeInfo) {
	gt.mapKeyType = mapKeyType
}

// Map Value type
func (gt *GenericTypeInfo) MapValueType() TypeInfo {
	return gt.mapValueType
}

func (gt *GenericTypeInfo) SetMapValueType(mapValueType TypeInfo) {
	gt.mapValueType = mapValueType
}
