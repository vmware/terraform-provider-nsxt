/* Copyright © 2019-2020 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package bindings

import (
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/l10n"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/lib"
	"reflect"
)

var stringDefinition = data.NewStringDefinition()

// BindingTypeToDataDefinitionConverter creates data.DataDefinition from BindingType
type BindingTypeToDataDefinitionConverter struct {
	ctx            *data.ReferenceResolver
	seenStructures map[string]bool
	inValue        BindingType
	outValue       data.DataDefinition
}

func NewBindingTypeToDataDefinitionConverter(inValue BindingType, ctx *data.ReferenceResolver, seenStructures map[string]bool) *BindingTypeToDataDefinitionConverter {
	return &BindingTypeToDataDefinitionConverter{inValue: inValue, ctx: ctx, seenStructures: seenStructures}
}
func (converter *BindingTypeToDataDefinitionConverter) OutputValue() data.DataDefinition {
	return converter.outValue
}

func (converter *BindingTypeToDataDefinitionConverter) convert(typ BindingType) []error {
	switch reflect.TypeOf(typ) {
	case VoidBindingType:
		converter.visitVoidType(typ.(VoidType))
	case DoubleBindingType:
		converter.outValue = data.NewDoubleDefinition()
	case BlobBindingType:
		converter.outValue = data.NewBlobDefinition()
	case IntegerBindingType:
		converter.visitIntegerType(typ.(IntegerType))
	case StringBindingType:
		converter.outValue = stringDefinition
	case SecretBindingType:
		converter.outValue = data.NewSecretDefinition()
	case OptionalBindingType:
		return converter.visitOptionalType(typ.(OptionalType))
	case BooleanBindingType:
		converter.visitBooleanType(typ.(BooleanType))
	case OpaqueBindingType:
		converter.visitOpaqueType(typ.(OpaqueType))
	case StructBindingType:
		converter.visitStructType(typ.(StructType))
	case ListBindingType:
		return converter.visitListType(typ.(ListType))
	case MapBindingType:
		return converter.visitMapType(typ.(MapType))
	case IdBindingType:
		converter.visitIdType(typ.(IdType))
	case EnumBindingType:
		converter.visitEnumType(typ.(EnumType))
	case SetBindingType:
		return converter.visitSetType(typ.(SetType))
	case DynamicStructBindingType:
		converter.visitDynamicStructType(typ.(DynamicStructType))
	case DateTimeBindingType:
		converter.visitDateTimeType(typ.(DateTimeType))
	case UriBindingType:
		converter.visitUriType(typ.(UriType))
	case ReferenceBindingType:
		converter.visitReferenceType(typ.(ReferenceType))
	case ErrorBindingType:
		converter.visitErrorType(typ.(ErrorType))
	case AnyErrorBindingType:
		converter.visitAnyErrorType(typ.(AnyErrorType))
	default:
		return []error{l10n.NewRuntimeError(
			"vapi.bindings.typeconverter.invalid.type",
			map[string]string{"bindingType": reflect.TypeOf(typ).String()})}
	}
	return nil
}

func (converter *BindingTypeToDataDefinitionConverter) visitReferenceType(i ReferenceType) {
	converter.visitStructType(i.Resolve().(StructType))
}
func (converter *BindingTypeToDataDefinitionConverter) visitIntegerType(i IntegerType) {
	converter.outValue = data.NewIntegerDefinition()
}
func (converter *BindingTypeToDataDefinitionConverter) visitSetType(i SetType) []error {
	err := converter.convert(i.ElementType())
	if err != nil {
		return err
	}
	elementDef := converter.outValue
	converter.outValue = data.NewListDefinition(elementDef)
	return nil
}

func (converter *BindingTypeToDataDefinitionConverter) visitEnumType(i EnumType) {
	converter.outValue = stringDefinition
}

func (converter *BindingTypeToDataDefinitionConverter) visitIdType(i IdType) {
	converter.outValue = stringDefinition
}

func (converter *BindingTypeToDataDefinitionConverter) visitMapType(i MapType) []error {
	fieldDefs := make(map[string]data.DataDefinition)
	err := converter.convert(i.KeyType)
	if err != nil {
		return err
	}
	keyDef := converter.outValue
	fieldDefs[lib.MAP_KEY_FIELD] = keyDef
	err = converter.convert(i.ValueType)
	if err != nil {
		return err
	}
	valueDef := converter.outValue
	fieldDefs[lib.MAP_VALUE_FIELD] = valueDef
	elementDef := data.NewStructDefinition(lib.MAP_ENTRY, fieldDefs)
	converter.outValue = data.NewListDefinition(elementDef)
	return nil
}
func (converter *BindingTypeToDataDefinitionConverter) visitListType(i ListType) []error {
	err := converter.convert(i.ElementType())
	if err != nil {
		return err
	}
	elementDef := converter.outValue
	converter.outValue = data.NewListDefinition(elementDef)
	return nil
}

func (converter *BindingTypeToDataDefinitionConverter) visitStructType(i StructType) {
	structName := i.Name()
	if converter.seenStructures[structName] {
		if converter.ctx.IsDefined(structName) {
			converter.outValue = converter.ctx.Definition(structName)
		} else {
			// For recursive/self-referencing structures the resolution
			// of first/top-level is not yet completed when we convert the
			// second level of structure reference.
			// So, to break the cycle in building definitions we create
			// an unresolved reference (StructRefDefinition) to end it.
			structRef := data.NewStructRefDefinition(structName, nil)
			converter.ctx.AddReference(structRef)
			converter.outValue = structRef
		}
	} else {
		converter.seenStructures[structName] = true
		fieldDefMap := make(map[string]data.DataDefinition)
		for key, field := range i.fields {
			converter.convert(field)
			fieldDefMap[key] = converter.outValue
		}
		structDef := data.NewStructDefinition(i.name, fieldDefMap)
		converter.ctx.AddDefinition(structDef)
		converter.outValue = structDef
	}

}

func (converter *BindingTypeToDataDefinitionConverter) visitOpaqueType(i OpaqueType) {
	converter.outValue = data.NewOpaqueDefinition()
}

func (converter *BindingTypeToDataDefinitionConverter) visitBooleanType(i BooleanType) {
	converter.outValue = data.NewBooleanDefinition()
}

func (converter *BindingTypeToDataDefinitionConverter) visitOptionalType(i OptionalType) []error {
	err := converter.convert(i.ElementType())
	if err != nil {
		return err
	}
	elementDef := converter.outValue
	converter.outValue = data.NewOptionalDefinition(elementDef)
	return nil
}
func (converter *BindingTypeToDataDefinitionConverter) visitVoidType(i VoidType) {
	converter.outValue = data.NewVoidDefinition()
}
func (converter *BindingTypeToDataDefinitionConverter) visitDynamicStructType(i DynamicStructType) {
	converter.outValue = data.NewDynamicStructDefinition()
}
func (converter *BindingTypeToDataDefinitionConverter) visitDateTimeType(timeType DateTimeType) {
	converter.outValue = stringDefinition
}
func (converter *BindingTypeToDataDefinitionConverter) visitUriType(i UriType) {
	converter.outValue = stringDefinition
}
func (converter *BindingTypeToDataDefinitionConverter) visitErrorType(i ErrorType) {
	fieldDefMap := make(map[string]data.DataDefinition)
	for key, field := range i.fields {
		converter.convert(field)
		fieldDefMap[key] = converter.outValue
	}
	errorDef := data.NewErrorDefinition(i.name, fieldDefMap)
	converter.outValue = errorDef
}
func (converter *BindingTypeToDataDefinitionConverter) visitAnyErrorType(i AnyErrorType) {
	converter.outValue = data.NewAnyErrorDefinition()
}
