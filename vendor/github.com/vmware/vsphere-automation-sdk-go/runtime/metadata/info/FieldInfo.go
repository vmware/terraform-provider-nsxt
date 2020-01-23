/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package info


type FieldInfo struct {
	fieldName               string
	documentation           string
	isIdentifierPolymorphic bool
	identifierType          string
	identifierTypes         []string
	identifierTypeHolder    string
	typeInfo                TypeInfo
	hasFieldsOfTypeNames    []string
	resource                string
	isOneOf                 []string
}

func NewFieldInfo(fieldName string) *FieldInfo {
	fieldInfo := FieldInfo{fieldName: fieldName}
	return &fieldInfo
}

// Field Name
func (field *FieldInfo) FieldName() string {
	return field.fieldName
}

func (field *FieldInfo) SetFieldName(fieldName string) {
	field.fieldName = fieldName
}

// Documentation
func (field *FieldInfo) Documentation() string {
	return field.documentation
}

func (field *FieldInfo) SetDocumentation(documentation string) {
	field.documentation = documentation
}

// Is Identifier Polymorphic
func (field *FieldInfo) IsIdentifierPolymorphic() bool {
	return field.isIdentifierPolymorphic
}

func (field *FieldInfo) SetIsIdentifierPolymorphic(isIdentifierPolymorphic bool) {
	field.isIdentifierPolymorphic = isIdentifierPolymorphic
}

// Identifier Type
func (field *FieldInfo) IdentifierType() string {
	return field.identifierType
}

func (field *FieldInfo) SetIdentifierType(identifierType string) {
	field.identifierType = identifierType
}

// Identifier Types
func (field *FieldInfo) IdentifierTypes() []string {
	return field.identifierTypes
}

func (field *FieldInfo) SetIdentifierTypes(identifierTypes []string) {
	field.identifierTypes = identifierTypes
}

// Identifier Type Holder
func (field *FieldInfo) IdentifierTypeHolder() string {
	return field.identifierTypeHolder
}

func (field *FieldInfo) SetIdentifierTypeHolder(identifierTypeHolder string) {
	field.identifierTypeHolder = identifierTypeHolder
}

// type info
func (field *FieldInfo) TypeInfo() TypeInfo {
	return field.typeInfo
}

func (field *FieldInfo) SetTypeInfo(typeInfo TypeInfo) {
	field.typeInfo = typeInfo
}

// Has Fields Of Type Names
func (field *FieldInfo) HasFieldsOfTypeNames() []string {
	return field.hasFieldsOfTypeNames
}

func (field *FieldInfo) SetHasFieldsOfTypeNames(hasFieldsOfTypeNames []string) {
	field.hasFieldsOfTypeNames = hasFieldsOfTypeNames
}

func (field *FieldInfo) AddHasFieldsOfTypeName(hasFieldsOfTypeName string) {
	if field.hasFieldsOfTypeNames == nil {
		field.hasFieldsOfTypeNames = []string{hasFieldsOfTypeName}
	} else {
		field.hasFieldsOfTypeNames = append(field.hasFieldsOfTypeNames, hasFieldsOfTypeName)
	}
}

// Resource
func (field *FieldInfo) Resource() string {
	return field.resource
}

func (field *FieldInfo) SetResource(resource string) {
	field.resource = resource
}

// Is One Of
func (field *FieldInfo) IsOneOf() []string {
	return field.isOneOf
}

func (field *FieldInfo) SetIsOneOf(isOneOf []string) {
	field.isOneOf = isOneOf
}

func (field *FieldInfo) AddIsOneOf(element string) {
	if field.isOneOf == nil {
		field.isOneOf = []string{element}
	} else {
		field.isOneOf = append(field.isOneOf, element)
	}
}
