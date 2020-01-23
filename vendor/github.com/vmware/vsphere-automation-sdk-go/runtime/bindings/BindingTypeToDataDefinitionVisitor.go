/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package bindings

import (
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/l10n"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/lib"
	"reflect"
)

var stringDefinition = data.NewStringDefinition()

/**
 * Builds DataDefinition by visiting a BindingType
 */
type BindingTypeToDataDefinitionVisitor struct {
	ctx            *data.ReferenceResolver
	seenStructures map[string]bool
	inValue        BindingType
	outValue       data.DataDefinition
}

func NewBindingTypeToDataDefinitionVisitor(inValue BindingType, ctx *data.ReferenceResolver, seenStructures map[string]bool) *BindingTypeToDataDefinitionVisitor {
	return &BindingTypeToDataDefinitionVisitor{inValue: inValue, ctx: ctx, seenStructures: seenStructures}
}
func (b *BindingTypeToDataDefinitionVisitor) OutputValue() data.DataDefinition {
	return b.outValue
}

func (b *BindingTypeToDataDefinitionVisitor) visit(typ BindingType) []error {
	switch reflect.TypeOf(typ) {
	case VoidBindingType:
		b.visitVoidType(typ.(VoidType))
	case DoubleBindingType:
		b.outValue = data.NewDoubleDefinition()
	case BlobBindingType:
		b.outValue = data.NewBlobDefinition()
	case IntegerBindingType:
		b.visitIntegerType(typ.(IntegerType))
	case StringBindingType:
		b.outValue = stringDefinition
	case SecretBindingType:
		b.outValue = data.NewSecretDefinition()
	case OptionalBindingType:
		return b.visitOptionalType(typ.(OptionalType))
	case BooleanBindingType:
		b.visitBooleanType(typ.(BooleanType))
	case OpaqueBindingType:
		b.visitOpaqueType(typ.(OpaqueType))
	case StructBindingType:
		b.visitStructType(typ.(StructType))
	case ListBindingType:
		return b.visitListType(typ.(ListType))
	case MapBindingType:
		return b.visitMapType(typ.(MapType))
	case IdBindingType:
		b.visitIdType(typ.(IdType))
	case EnumBindingType:
		b.visitEnumType(typ.(EnumType))
	case SetBindingType:
		return b.visitSetType(typ.(SetType))
	case DynamicStructBindingType:
		b.visitDynamicStructType(typ.(DynamicStructType))
	case DateTimeBindingType:
		b.visitDateTimeType(typ.(DateTimeType))
	case UriBindingType:
		b.visitUriType(typ.(UriType))
	case ReferenceBindingType:
		b.visitReferenceType(typ.(ReferenceType))
	case ErrorBindingType:
		b.visitErrorType(typ.(ErrorType))
	case AnyErrorBindingType:
		b.visitAnyErrorType(typ.(AnyErrorType))
	default:
		return []error{l10n.NewRuntimeError(
			"vapi.bindings.typeconverter.invalid.type",
			map[string]string{"bindingType": reflect.TypeOf(typ).String()})}
	}
	return nil
}

func (visitor *BindingTypeToDataDefinitionVisitor) visitReferenceType(i ReferenceType) {
	visitor.visitStructType(i.Resolve().(StructType))
}
func (visitor *BindingTypeToDataDefinitionVisitor) visitIntegerType(i IntegerType) {
	visitor.outValue = data.NewIntegerDefinition()
}
func (visitor *BindingTypeToDataDefinitionVisitor) visitSetType(i SetType) []error {
	err := i.ElementType().Accept(visitor)
	if err != nil {
		return err
	}
	elementDef := visitor.outValue
	visitor.outValue = data.NewListDefinition(elementDef)
	return nil
}

func (visitor *BindingTypeToDataDefinitionVisitor) visitEnumType(i EnumType) {
	visitor.outValue = stringDefinition
}

func (visitor *BindingTypeToDataDefinitionVisitor) visitIdType(i IdType) {
	visitor.outValue = stringDefinition
}

func (visitor *BindingTypeToDataDefinitionVisitor) visitMapType(i MapType) []error {
	fieldDefs := make(map[string]data.DataDefinition)
	err := i.KeyType.Accept(visitor)
	if err != nil {
		return err
	}
	keyDef := visitor.outValue
	fieldDefs[lib.MAP_KEY_FIELD] = keyDef
	err = i.ValueType.Accept(visitor)
	if err != nil {
		return err
	}
	valueDef := visitor.outValue
	fieldDefs[lib.MAP_VALUE_FIELD] = valueDef
	elementDef := data.NewStructDefinition(lib.MAP_ENTRY, fieldDefs)
	visitor.outValue = data.NewListDefinition(elementDef)
	return nil
}
func (visitor *BindingTypeToDataDefinitionVisitor) visitListType(i ListType) []error {
	err := i.ElementType().Accept(visitor)
	if err != nil {
		return err
	}
	elementDef := visitor.outValue
	visitor.outValue = data.NewListDefinition(elementDef)
	return nil
}

func (visitor *BindingTypeToDataDefinitionVisitor) visitStructType(i StructType) {
	structName := i.Name()
	if visitor.seenStructures[structName] {
		if visitor.ctx.IsDefined(structName) {
			visitor.outValue = visitor.ctx.Definition(structName)
		} else {
			// For recursive/self-referencing structures the resolution
			// of first/top-level is not yet completed when we visit the
			// second level of structure reference.
			// So, to break the cycle in building definitions we create
			// an unresolved reference (StructRefDefinition) to end it.
			structRef := data.NewStructRefDefinition(structName, nil)
			visitor.ctx.AddReference(structRef)
			visitor.outValue = structRef
		}
	} else {
		visitor.seenStructures[structName] = true
		fieldDefMap := make(map[string]data.DataDefinition)
		for key, field := range i.fields {
			field.Accept(visitor)
			fieldDefMap[key] = visitor.outValue
		}
		structDef := data.NewStructDefinition(i.name, fieldDefMap)
		visitor.ctx.AddDefinition(structDef)
		visitor.outValue = structDef
	}

}

func (visitor *BindingTypeToDataDefinitionVisitor) visitOpaqueType(i OpaqueType) {
	visitor.outValue = data.NewOpaqueDefinition()
}

func (visitor *BindingTypeToDataDefinitionVisitor) visitBooleanType(i BooleanType) {
	visitor.outValue = data.NewBooleanDefinition()
}

func (visitor *BindingTypeToDataDefinitionVisitor) visitOptionalType(i OptionalType) []error {
	err := i.ElementType().Accept(visitor)
	if err != nil {
		return err
	}
	elementDef := visitor.outValue
	visitor.outValue = data.NewOptionalDefinition(elementDef)
	return nil
}
func (visitor *BindingTypeToDataDefinitionVisitor) visitVoidType(i VoidType) {
	visitor.outValue = data.NewVoidDefinition()
}
func (visitor *BindingTypeToDataDefinitionVisitor) visitDynamicStructType(i DynamicStructType) {
	visitor.outValue = data.NewDynamicStructDefinition()
}
func (visitor *BindingTypeToDataDefinitionVisitor) visitDateTimeType(timeType DateTimeType) {
	visitor.outValue = stringDefinition
}
func (visitor *BindingTypeToDataDefinitionVisitor) visitUriType(i UriType) {
	visitor.outValue = stringDefinition
}
func (visitor *BindingTypeToDataDefinitionVisitor) visitErrorType(i ErrorType) {
	fieldDefMap := make(map[string]data.DataDefinition)
	for key, field := range i.fields {
		field.Accept(visitor)
		fieldDefMap[key] = visitor.outValue
	}
	errorDef := data.NewErrorDefinition(i.name, fieldDefMap)
	visitor.outValue = errorDef
}
func (visitor *BindingTypeToDataDefinitionVisitor) visitAnyErrorType(i AnyErrorType) {
	visitor.outValue = data.NewAnyErrorDefinition()
}
