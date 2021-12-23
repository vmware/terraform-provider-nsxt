/* Copyright © 2019-2020 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package data

import (
	"errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/l10n"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/lib"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/log"
	"strconv"
)

type DataDefinition interface {
	/**
	* Returns the {@link DataType} for this type.
	*
	* @return  the {@link DataType} for this type
	 */
	Type() DataType

	/**
	 * Validates that the specified {@link DataValue} is an instance of this
	 * data definition.
	 *
	 * <p>Validates that supplied <code>value</code> is not <code>nil</code>
	 * and it's type matches the type of this definition.
	 *
	 * @param value  the <code>DataValue</code> to validate
	 * @return       a list of messages describing validation problems; empty
	 *               list indicates that validation was successful
	 */
	Validate(value DataValue) []error

	String() string
}

type AnyErrorDefinition struct {
}

func NewAnyErrorDefinition() AnyErrorDefinition {
	return AnyErrorDefinition{}
}

func (a AnyErrorDefinition) String() string {
	return a.Type().String()
}

func (a AnyErrorDefinition) Type() DataType {
	return ANY_ERROR
}

func (a AnyErrorDefinition) Validate(value DataValue) []error {
	if value != nil && value.Type() == ERROR {
		return nil
	}
	return a.Type().Validate(value)
}

type BlobDefinition struct{}

func NewBlobDefinition() BlobDefinition {
	var instance = BlobDefinition{}
	return instance
}

func (b BlobDefinition) Type() DataType {
	return BLOB
}

func (b BlobDefinition) Validate(value DataValue) []error {
	return b.Type().Validate(value)
}

func (b BlobDefinition) String() string {
	return BLOB.String()
}

func (b BlobDefinition) NewValue(value []byte) *BlobValue {
	return NewBlobValue(value)
}

type BooleanDefinition struct{}

func NewBooleanDefinition() BooleanDefinition {
	return BooleanDefinition{}
}

func (b BooleanDefinition) Type() DataType {
	return BOOLEAN
}

func (b BooleanDefinition) Validate(value DataValue) []error {
	return b.Type().Validate(value)
}

func (b BooleanDefinition) String() string {
	return BOOLEAN.String()
}

func (b BooleanDefinition) Equals(other DataDefinition) bool {
	return other.Type() == BOOLEAN
}

type DoubleDefinition struct{}

func NewDoubleDefinition() DoubleDefinition {
	return DoubleDefinition{}
}
func (d DoubleDefinition) Type() DataType {
	return DOUBLE
}

func (d DoubleDefinition) Validate(value DataValue) []error {
	return d.Type().Validate(value)
}

func (d DoubleDefinition) String() string {
	return DOUBLE.String()
}

type DynamicStructDefinition struct {
	validDataTypes []DataType
}

func NewDynamicStructDefinition() DynamicStructDefinition {
	var validDataTypes = []DataType{STRUCTURE}
	return DynamicStructDefinition{validDataTypes: validDataTypes}
}

func (dynamicStructure DynamicStructDefinition) Type() DataType {
	return DYNAMIC_STRUCTURE
}

func (dynamicStructure DynamicStructDefinition) Validate(value DataValue) []error {
	if value != nil && value.Type() == STRUCTURE {
		return nil
	}
	return dynamicStructure.Type().Validate(value)
}

func (dynamicStructure DynamicStructDefinition) String() string {
	return dynamicStructure.Type().String()
}

type ErrorDefinition struct {
	name   string
	fields map[string]DataDefinition
}

func NewErrorDefinition(name string, fields map[string]DataDefinition) ErrorDefinition {
	if len(name) == 0 {
		log.Error("ErrorDefinition name missing")
	}
	if fields == nil {
		log.Error("Missing fields for ErrorDefinition")
	}
	for key, value := range fields {
		if value == nil {
			log.Error("Missing value for field " + key)
		}
	}
	return ErrorDefinition{name: name, fields: fields}
}

func (errorDefinition ErrorDefinition) Type() DataType {
	return ERROR
}

func (errorDefinition ErrorDefinition) Name() string {
	return errorDefinition.name
}

func (errorDefinition ErrorDefinition) FieldNames() []string {
	var keys = make([]string, len(errorDefinition.fields))
	var i = 0
	for k := range errorDefinition.fields {
		keys[i] = k
		i++
	}
	return keys
}

func (errorDefinition ErrorDefinition) Field(field string) DataDefinition {
	//TODO:
	// zero value of pointer is nil
	// what is the zero value of a structure
	var value, ok = errorDefinition.fields[field]
	if ok {
		return value
	} else {
		return nil
	}

}

func (errorDefinition ErrorDefinition) HasField(field string) bool {
	var _, ok = errorDefinition.fields[field]
	return ok
}

func (errorDefinition ErrorDefinition) Validate(value DataValue) []error {
	var result = errorDefinition.Type().Validate(value)
	if result != nil {
		return result
	}
	var errorValue = value.(*ErrorValue)
	if errorDefinition.Name() != errorValue.Name() {
		var args = map[string]string{
			"actualName":   errorValue.Name(),
			"expectedName": errorDefinition.Name()}
		return []error{l10n.NewRuntimeError("vapi.data.structure.name.mismatch", args)}
	}
	// Make sure no fields are missing and those that are present are valid
	for _, fieldName := range errorDefinition.FieldNames() {
		var fieldDef = errorDefinition.Field(fieldName)
		var fieldVal, err = errorValue.Field(fieldName)
		if err != nil {
			var args = map[string]string{
				"structName": errorDefinition.Name(),
				"fieldName":  fieldName}
			return []error{l10n.NewRuntimeError("vapi.data.structure.field.missing", args)}
		}

		var errs = fieldDef.Validate(fieldVal)
		if errs != nil {
			var args = map[string]string{
				"structName": errorDefinition.Name(),
				"fieldName":  fieldName}
			var msg = l10n.NewRuntimeError("vapi.data.structure.field.invalid", args)
			errs = append(errs, msg)
			return errs
		}
	}
	return nil
}

func (errorDefinition ErrorDefinition) String() string {
	return errorDefinition.Type().String()
}

func (errorDefinition ErrorDefinition) NewValue() *ErrorValue {
	return NewErrorValue(errorDefinition.name, nil)
}

type IntegerDefinition struct{}

func NewIntegerDefinition() IntegerDefinition {
	return IntegerDefinition{}
}

func (integerDefinition IntegerDefinition) Type() DataType {
	return INTEGER
}

func (integerDefinition IntegerDefinition) Validate(value DataValue) []error {
	return integerDefinition.Type().Validate(value)
}

func (integerDefinition IntegerDefinition) String() string {
	return INTEGER.String()
}

type ListDefinition struct {
	/**
	 * The Definition of the type that is contained in a list.
	 */
	elementType DataDefinition
}

func NewListDefinition(elementType DataDefinition) ListDefinition {
	return ListDefinition{elementType: elementType}
}

func (listDefinition ListDefinition) Type() DataType {
	return LIST
}

func (listDefinition ListDefinition) ElementType() DataDefinition {
	return listDefinition.elementType
}

func (listDefinition ListDefinition) Validate(value DataValue) []error {
	// relax specific validation scenarios to support rest
	if _, ok := value.(*StructValue); ok {
		return nil
	}
	var result = listDefinition.Type().Validate(value)
	if result != nil {
		return result
	}
	var listValue = (value).(*ListValue)
	elementType := listDefinition.elementType
	for index, listElement := range listValue.List() {
		var subErrors = elementType.Validate(listElement)
		if len(subErrors) != 0 {
			args := map[string]string{
				"value": listElement.Type().String(),
				"index": strconv.Itoa(index)}
			var msg = l10n.NewRuntimeError("vapi.data.list.invalid.entry", args)
			result = append(result, subErrors...)
			result = append(result, msg)
		}
	}
	return result
}

func (listDefinition ListDefinition) String() string {
	return LIST.String()
}

func (listDefinition ListDefinition) NewValue() *ListValue {
	return NewListValue()
}

type SecretDefinition struct{}

func (s SecretDefinition) Type() DataType {
	return SECRET
}

func NewSecretDefinition() SecretDefinition {
	var instance = SecretDefinition{}
	return instance
}

func (s SecretDefinition) Validate(value DataValue) []error {
	return s.Type().Validate(value)
}

func (s SecretDefinition) String() string {
	return SECRET.String()
}

func (s SecretDefinition) NewValue(value string) *SecretValue {
	return NewSecretValue(value)
}

type StringDefinition struct{}

func NewStringDefinition() StringDefinition {
	var instance = StringDefinition{}
	return instance
}

func (s StringDefinition) Type() DataType {
	return STRING
}

func (s StringDefinition) Validate(value DataValue) []error {
	return s.Type().Validate(value)
}

func (s StringDefinition) String() string {
	return STRING.String()
}

func (s StringDefinition) NewValue(value string) *StringValue {
	return NewStringValue(value)
}

type StructDefinition struct {
	fields map[string]DataDefinition
	name   string
}

/**
* Constructor.
*
* @param name name for the structure; must not be <code>null</code>
* @param fields field names and definitions for the structure;
*               must not be <code>null</code>
* @throws IllegalArgumentException if some of the arguments or the
*         values in <code>fields</code> is <code>null</code>
 */

func NewStructDefinition(name string, fields map[string]DataDefinition) StructDefinition {
	if len(name) == 0 {
		log.Error("StructDefinition name missing")
	}
	if fields == nil {
		log.Error("Missing fields for StructDefinition")
	}
	for key, value := range fields {
		if value == nil {
			log.Errorf("Missing value for field %s", key)
			return StructDefinition{}
		}
	}
	return StructDefinition{name: name, fields: fields}
}

func (structDefinition StructDefinition) Type() DataType {
	return STRUCTURE
}

func (structDefinition StructDefinition) Name() string {
	return structDefinition.name
}

func (structDefinition StructDefinition) FieldNames() []string {
	var keys = make([]string, len(structDefinition.fields))
	var i = 0
	for k := range structDefinition.fields {
		keys[i] = k
		i++
	}
	return keys
}

func (structDefinition StructDefinition) Field(field string) DataDefinition {
	//TODO:
	// zero value of pointer is nil
	// what is the zero value of a structure
	var value, ok = structDefinition.fields[field]
	if ok {
		return value
	} else {
		return nil
	}

}

func (structDefinition StructDefinition) HasField(field string) bool {
	var _, ok = structDefinition.fields[field]
	if ok {
		return true
	}
	return false
}

func (structDefinition StructDefinition) Validate(value DataValue) []error {
	var result = structDefinition.Type().Validate(value)
	if result != nil {
		return result
	}
	var structValue = value.(*StructValue)
	if structDefinition.Name() != structValue.Name() {
		if !(structDefinition.Name() == lib.MAP_STRUCT ||
			structDefinition.Name() == lib.MAP_ENTRY ||
			structDefinition.Name() == lib.OPERATION_INPUT) {
			var args = map[string]string{
				"actualName":   structValue.Name(),
				"expectedName": structDefinition.Name()}
			return []error{l10n.NewRuntimeError("vapi.data.structure.name.mismatch", args)}
		}
	}
	// Make sure no fields are missing and those that are present are valid
	for _, fieldName := range structDefinition.FieldNames() {
		var fieldDef = structDefinition.Field(fieldName)
		var fieldVal, err = structValue.Field(fieldName)
		if err != nil && fieldDef.Type() != OPTIONAL {
			var args = map[string]string{
				"structName": structDefinition.Name(),
				"fieldName":  fieldName}
			return []error{l10n.NewRuntimeError("vapi.data.structure.field.missing", args)}
		}

		var errs = fieldDef.Validate(fieldVal)
		if errs != nil {
			var args = map[string]string{
				"structName": structDefinition.Name(),
				"fieldName":  fieldName}
			var msg = l10n.NewRuntimeError("vapi.data.structure.field.invalid", args)
			errs = append(errs, msg)
			return errs
		}
	}
	return result
}

func (structDefinition StructDefinition) String() string {
	return STRUCTURE.String()
}

// TODO:
// lot of structdef tests depend on this method. remove it after correcting tests
func (structDefinition StructDefinition) Equals(other DataDefinition) bool {
	if other == nil {
		return false
	}
	if other.Type() != STRUCTURE {
		return false
	}
	var otherStruct = other.(StructDefinition)
	if (structDefinition.name) != (otherStruct.name) {
		return false
	}

	for key := range structDefinition.fields {
		if _, ok := otherStruct.fields[key]; !ok {
			return false
		}
	}
	return true

}

func (structDefinition StructDefinition) NewValue() *StructValue {
	return NewStructValue(structDefinition.name, nil)

}

/**
 * Reference to a {@link StructDefinition}. If the reference is resolved, it is
 * bound to a specific {@link StructDefinition} target. If the reference is
 * unresolved, its target is <code>nil</code>.
 */
/**
 * StructRef is different when compared to other Definition structs.
 * It returns a pointer instead of value because of the resolving task involved.
 * When StructRef gets resolved, all its references will also be resolved.
 */
type StructRefDefinition struct {
	name             string
	structDefinition *StructDefinition
}

func NewStructRefDefinition(name string, structDefinition *StructDefinition) *StructRefDefinition {
	if name == "" {
		log.Error("StructRef name missing")
	}
	return &StructRefDefinition{name: name, structDefinition: structDefinition}
}

func (structRefDefinition StructRefDefinition) Name() string {
	return structRefDefinition.name
}

func (structRefDefinition StructRefDefinition) Target() *StructDefinition {
	return structRefDefinition.structDefinition
}

func (structRefDefinition *StructRefDefinition) SetTarget(structDefinition *StructDefinition) error {
	if structDefinition == nil {
		return errors.New("StructDefinition may not be nil")
	}
	if structRefDefinition.structDefinition != nil {
		return errors.New("structref already resolved")
	}
	if structRefDefinition.name != structDefinition.name {
		return errors.New("type mismatch")
	}

	structRefDefinition.structDefinition = structDefinition
	return nil
}

func (structRefDefinition StructRefDefinition) Type() DataType {
	return STRUCTURE_REF
}

func (structRefDefinition StructRefDefinition) Validate(value DataValue) []error {
	var errorList = structRefDefinition.CheckResolved()
	if errorList != nil {
		return errorList
	}
	return structRefDefinition.structDefinition.Validate(value)

}

func (structRefDefinition StructRefDefinition) String() string {
	return STRUCTURE_REF.String()
}

func (structRefDefinition StructRefDefinition) CheckResolved() []error {
	if structRefDefinition.structDefinition == nil {
		return []error{l10n.NewRuntimeError("vapi.data.structref.not.resolved",
			map[string]string{"referenceType": structRefDefinition.name})}
	}
	return nil
}

type OpaqueDefinition struct{}

func NewOpaqueDefinition() OpaqueDefinition {
	return OpaqueDefinition{}
}

func (opaqueDefinition OpaqueDefinition) Type() DataType {
	return OPAQUE
}

func (opaqueDefinition OpaqueDefinition) String() string {
	return OPAQUE.String()
}

/**
 * Validates that the specified DataValue is an instance of this data
 * definition.
 *
 * <p>Only validates that supplied value is not <code>nil</code>.
 */
func (opaqueDefinition OpaqueDefinition) Validate(value DataValue) []error {
	if value == nil {
		var msg = l10n.NewRuntimeError("vapi.data.opaque.definition.null.value", map[string]string{})
		return []error{msg}
	}
	return nil
}

type OptionalDefinition struct {
	//cannot be nil
	optionalElementType DataDefinition
}

func NewOptionalDefinition(optionalElementType DataDefinition) OptionalDefinition {
	return OptionalDefinition{optionalElementType: optionalElementType}
}

func (optionalDefinition OptionalDefinition) Type() DataType {
	return OPTIONAL
}

func (optionalDefinition OptionalDefinition) ElementType() DataDefinition {
	return optionalDefinition.optionalElementType
}

func (optionalDefinition OptionalDefinition) Validate(value DataValue) []error {
	var result = optionalDefinition.Type().Validate(value)
	if result != nil {
		// Error is returned for nil value. But optional value can be nil in permissive mode?
		return result
	}
	var optionalValue = value.(*OptionalValue)

	if optionalValue.IsSet() {
		var subErrors = optionalDefinition.ElementType().Validate(optionalValue.Value())
		if subErrors != nil {
			result = append(result, subErrors...)
		}
	} else {
		// nulls are valid for optional so this is a
		// no-op for a good reason
	}
	return result
}

func (optionalDefinition OptionalDefinition) String() string {
	//TODO:sreeshas
	// implement this
	return ""
}

func (optionalDefinition OptionalDefinition) NewValue(value DataValue) *OptionalValue {
	return NewOptionalValue(value)
}

type VoidDefinition struct{}

func NewVoidDefinition() VoidDefinition {
	return VoidDefinition{}
}

func (voidDefinition VoidDefinition) Type() DataType {
	return VOID
}

func (voidDefinition VoidDefinition) Validate(value DataValue) []error {
	return voidDefinition.Type().Validate(value)
}

func (voidDefinition VoidDefinition) String() string {
	return VOID.String()
}

func (voidDefinition VoidDefinition) NewValue() *VoidValue {
	return NewVoidValue()
}
