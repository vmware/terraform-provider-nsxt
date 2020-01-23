/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package bindings

import (
	"fmt"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/l10n"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/lib"
	"net/url"
	"reflect"
	"time"
)

type GolangToVapiDataValueVisitor struct {
	inValue  interface{}
	outValue data.DataValue
}

func NewGolangToVapiDataValueVisitor(inValue interface{}) *GolangToVapiDataValueVisitor {
	return &GolangToVapiDataValueVisitor{inValue: inValue}
}

func (g *GolangToVapiDataValueVisitor) OutputValue() data.DataValue {
	return g.outValue
}

func (g *GolangToVapiDataValueVisitor) visit(bindingType BindingType) []error {
	switch reflect.TypeOf(bindingType) {
	case IntegerBindingType:
		return g.visitIntegerType()
	case DoubleBindingType:
		return g.visitDoubleType()
	case StringBindingType:
		return g.visitStringType()
	case BooleanBindingType:
		return g.visitBooleanType()
	case OpaqueBindingType:
		return g.visitOpaqueType()
	case IdBindingType:
		return g.visitIdType()
	case EnumBindingType:
		return g.visitEnumType()
	case VoidBindingType:
		g.visitVoidType()
		return nil
	case BlobBindingType:
		return g.visitBlobType()
	case SecretBindingType:
		return g.visitSecretType()
	case DateTimeBindingType:
		return g.visitDateTimeType()
	case UriBindingType:
		return g.visitURIType()
	case OptionalBindingType:
		return g.visitOptionalType(bindingType.(OptionalType))
	case StructBindingType:
		return g.visitStructType(bindingType.(StructType))
	case DynamicStructBindingType:
		return g.visitDynamicStructType(bindingType.(DynamicStructType))
	case ErrorBindingType:
		return g.visitErrorType(bindingType.(ErrorType))
	case ListBindingType:
		return g.visitListType(bindingType.(ListType))
	case MapBindingType:
		return g.visitMapType(bindingType.(MapType))
	case SetBindingType:
		return g.visitSetType(bindingType.(SetType))
	case AnyErrorBindingType:
		return g.visitAnyErrorType(bindingType.(AnyErrorType))
	case ReferenceBindingType:
		return g.visitReferenceType(bindingType.(ReferenceType))
	default:
		return []error{l10n.NewRuntimeError("vapi.bindings.typeconverter.invalid.type",
			map[string]string{"bindingType": reflect.TypeOf(bindingType).String()})}
	}
}

func (g *GolangToVapiDataValueVisitor) visitReferenceType(bindingType ReferenceType) []error {
	return g.visitStructType(bindingType.Resolve().(StructType))
}

func (g *GolangToVapiDataValueVisitor) visitVoidType() {
	g.outValue = data.NewVoidValue()
}

func (g *GolangToVapiDataValueVisitor) visitIntegerType() []error {
	if intVal, ok := g.inValue.(int64); ok {
		g.outValue = data.NewIntegerValue(intVal)
		return nil
	}
	return g.unexpectedValueError("int64")
}

func (g *GolangToVapiDataValueVisitor) visitDoubleType() []error {
	if floatVal, ok := g.inValue.(float64); ok {
		g.outValue = data.NewDoubleValue(floatVal)
		return nil
	}
	return g.unexpectedValueError("float64")
}

func (g *GolangToVapiDataValueVisitor) visitStringType() []error {
	if stringVal, ok := g.inValue.(string); ok {
		g.outValue = data.NewStringValue(stringVal)
		return nil
	}
	return g.unexpectedValueError("string")
}

func (g *GolangToVapiDataValueVisitor) visitIdType() []error {
	return g.visitStringType()
}

func (g *GolangToVapiDataValueVisitor) visitBooleanType() []error {
	if boolVal, ok := g.inValue.(bool); ok {
		g.outValue = data.NewBooleanValue(boolVal)
		return nil
	}
	return g.unexpectedValueError("bool")
}

func (g *GolangToVapiDataValueVisitor) visitSecretType() []error {
	if stringVal, ok := g.inValue.(string); ok {
		g.outValue = data.NewSecretValue(stringVal)
		return nil
	}
	return g.unexpectedValueError("string")
}

func (g *GolangToVapiDataValueVisitor) visitOpaqueType() []error {
	if dataVal, ok := g.inValue.(data.DataValue); ok {
		g.outValue = dataVal
		return nil
	}
	return g.unexpectedValueError("DataValue")

}

func (g *GolangToVapiDataValueVisitor) visitBlobType() []error {
	inValue := g.inValue
	if byteVal, ok := inValue.([]byte); ok {
		g.outValue = data.NewBlobValue(byteVal)
		return nil
	}
	return g.unexpectedValueError("[]byte")

}

func (g *GolangToVapiDataValueVisitor) visitDateTimeType() []error {
	if t1, ok := g.inValue.(time.Time); ok {
		stringValue := data.NewStringValue(t1.Format(VAPI_DATETIME_LAYOUT))
		g.outValue = stringValue
		return nil
	}
	return g.unexpectedValueError("time.Time")
}

func (g *GolangToVapiDataValueVisitor) visitURIType() []error {
	if uriValue, ok := g.inValue.(url.URL); ok {
		g.outValue = data.NewStringValue((&uriValue).String())
		return nil
	}
	return g.unexpectedValueError("url.URL")
}

func (g *GolangToVapiDataValueVisitor) visitEnumType() []error {
	g.outValue = data.NewStringValue(reflect.ValueOf(g.inValue).String())
	return nil
}

func (g *GolangToVapiDataValueVisitor) visitStructType(structType StructType) []error {
	var inValue = g.inValue
	var e reflect.Value
	if isPointer(inValue) {
		e = reflect.ValueOf(inValue).Elem()
	} else {
		e = reflect.ValueOf(inValue)
	}
	var result = data.NewStructValue(structType.Name(), nil)
	var fieldNames = structType.FieldNames()
	for _, fieldName := range fieldNames {
		var x = e.FieldByName(structType.canonicalFieldMap[fieldName])
		//todo
		//this is a possible place to check required fields are not nil.
		//useful error can be returned from here.
		if x.Kind() == reflect.Ptr {
			if x.IsNil() {
				g.inValue = nil
			} else {
				bindingType := structType.Field(fieldName)
				if _, ok := bindingType.(DynamicStructType); ok {
					g.inValue = x.Interface()
				} else if optionalType, optType := bindingType.(OptionalType); optType {
					if _, optDS := optionalType.ElementType().(DynamicStructType); optDS {
						//if optional dynamic structure
						g.inValue = x.Interface()
					} else if _, optAE := optionalType.ElementType().(AnyErrorType); optAE {
						//if optional AnyErrorType
						g.inValue = x.Interface()
					} else {
						g.inValue = x.Elem().Interface()
					}
				} else if _, ok := bindingType.(AnyErrorType); ok {
					g.inValue = x.Interface()
				}
			}
		} else {
			if x.IsValid() {
				g.inValue = x.Interface()
			} else {
				g.inValue = nil
			}
		}
		//Where is missing field error?
		err := g.visit(structType.Field(fieldName))
		if err != nil {
			g.inValue = inValue
			g.outValue = result
			err = append(err, l10n.NewRuntimeError("vapi.bindings.typeconverter.struct.field.invalid",
				map[string]string{"fieldName": fieldName, "structName": structType.name}))
			return err
		}
		result.SetField(fieldName, g.outValue)
	}
	g.inValue = inValue
	g.outValue = result
	err := structType.Validate(result)
	if err != nil {
		return err
	}
	return nil
}

func (g *GolangToVapiDataValueVisitor) visitErrorType(errorType ErrorType) []error {
	var inValue = g.inValue
	var e reflect.Value
	if isPointer(inValue) {
		e = reflect.ValueOf(inValue).Elem()
	} else {
		e = reflect.ValueOf(inValue)
	}
	var result = data.NewErrorValue(errorType.Name(), nil)
	var fieldNames = errorType.FieldNames()
	for _, fieldName := range fieldNames {
		var x = e.FieldByName(errorType.canonicalFieldMap[fieldName])
		if x.Kind() == reflect.Ptr {
			if x.IsNil() {
				g.inValue = nil
			} else {
				g.inValue = x.Elem().Interface()
			}

		} else {
			g.inValue = x.Interface()
		}
		err := g.visit(errorType.Field(fieldName))
		if err != nil {
			g.inValue = inValue
			g.outValue = result
			err = append(err, l10n.NewRuntimeError("vapi.bindings.typeconverter.struct.field.invalid",
				map[string]string{"fieldName": fieldName, "structName": errorType.name}))
			return err
		}
		result.SetField(fieldName, g.outValue)
	}
	g.inValue = inValue
	g.outValue = result
	return nil
}

func (g *GolangToVapiDataValueVisitor) visitDynamicStructType(structType DynamicStructType) []error {

	if structVal, ok := g.inValue.(*data.StructValue); ok {
		g.outValue = structVal
		msgs := structType.Validate(structVal)
		if msgs != nil {
			return msgs
		}
		return nil
	}
	return g.unexpectedValueError("*data.StructValue")
}

func (g *GolangToVapiDataValueVisitor) visitAnyErrorType(i AnyErrorType) []error {
	if errorVal, ok := g.inValue.(*data.ErrorValue); ok {
		g.outValue = errorVal
		return nil
	}
	return g.unexpectedValueError("ErrorValue")
}

func (g *GolangToVapiDataValueVisitor) visitOptionalType(typ OptionalType) []error {
	if g.inValue == nil {
		g.outValue = data.NewOptionalValue(nil)
	} else if reflect.TypeOf(typ.ElementType()) == OpaqueBindingType {
		g.outValue = data.NewOptionalValue(g.inValue.(data.DataValue))
	} else if reflect.TypeOf(g.inValue).Kind() == reflect.Map && reflect.ValueOf(g.inValue).IsNil() {
		g.outValue = data.NewOptionalValue(nil)
	} else if reflect.TypeOf(g.inValue).Kind() == reflect.Slice && reflect.ValueOf(g.inValue).IsNil() {
		g.outValue = data.NewOptionalValue(nil)
	} else if reflect.TypeOf(g.inValue).Kind() == reflect.Ptr && reflect.ValueOf(g.inValue).IsNil() {
		g.outValue = data.NewOptionalValue(nil)
	} else if reflect.TypeOf(typ.ElementType()) == DynamicStructBindingType || reflect.TypeOf(typ.ElementType()) == AnyErrorBindingType {
		// in case of optional dynamicstructure and anyerrortype, do not de-reference the pointer
		// since optional and non optional cases are represented in the same way.
		err := g.visit(typ.ElementType())
		if err != nil {
			return err
		}
		g.outValue = data.NewOptionalValue(g.outValue)
	} else {
		x := reflect.TypeOf(g.inValue)
		if x.Kind() == reflect.Ptr {
			g.inValue = reflect.ValueOf(g.inValue).Elem().Interface()
		}
		err := g.visit(typ.ElementType())
		if err != nil {
			return err
		}
		g.outValue = data.NewOptionalValue(g.outValue)
	}
	return nil
}

func (g *GolangToVapiDataValueVisitor) visitListType(listType ListType) []error {
	var inValue = g.inValue
	if g.inValue == nil || reflect.ValueOf(g.inValue).IsNil() {
		g.outValue = nil
		return []error{l10n.NewRuntimeErrorNoParam("vapi.bindings.typeconverter.value.nil")}
	}
	result := data.NewListValue()
	s := reflect.ValueOf(inValue)
	for i := 0; i < s.Len(); i++ {
		g.inValue = s.Index(i).Interface()
		err := g.visit(listType.ElementType())
		if err != nil {
			err = append(err, l10n.NewRuntimeError("vapi.bindings.typeconverter.list.entry.invalid",
				map[string]string{"index": fmt.Sprintf("%d", i)}))
			return err
		}
		result.Add(g.outValue)
	}
	g.inValue = inValue
	g.outValue = result
	return nil
}

func (g *GolangToVapiDataValueVisitor) visitMapType(mapType MapType) []error {
	inValue := g.inValue
	if g.inValue == nil || reflect.ValueOf(g.inValue).IsNil() {
		g.outValue = nil
		return []error{l10n.NewRuntimeErrorNoParam("vapi.bindings.typeconverter.value.nil")}
	}
	result := data.NewListValue()
	m := reflect.ValueOf(inValue)
	keys := m.MapKeys()
	for _, key := range keys {
		structVal := data.NewStructValue(lib.MAP_ENTRY, nil)
		g.inValue = key.Interface()
		err := g.visit(mapType.KeyType)
		if err != nil {
			err = append(err, l10n.NewRuntimeErrorNoParam("vapi.bindings.typeconverter.dict.key.invalid"))
			return err
		}
		structVal.SetField(lib.MAP_KEY_FIELD, g.outValue)
		g.inValue = m.MapIndex(key).Interface()
		err = g.visit(mapType.ValueType)
		if err != nil {
			err = append(err, l10n.NewRuntimeError("vapi.bindings.typeconverter.dict.value.invalid",
				map[string]string{"key": key.String()}))
			return err
		}
		structVal.SetField(lib.MAP_VALUE_FIELD, g.outValue)
		result.Add(structVal)
	}
	g.inValue = inValue
	g.outValue = result
	return nil
}

func (g *GolangToVapiDataValueVisitor) visitSetType(setType SetType) []error {
	var inValue = g.inValue
	if g.inValue == nil || reflect.ValueOf(inValue).IsNil() {
		g.outValue = nil
		return []error{l10n.NewRuntimeErrorNoParam("vapi.bindings.typeconverter.value.nil")}
	}
	result := data.NewListValue()
	keys := reflect.ValueOf(inValue).MapKeys()
	for i := 0; i < len(keys); i++ {
		g.inValue = keys[i].Interface()
		err := g.visit(setType.ElementType())
		if err != nil {
			err = append(err, l10n.NewRuntimeErrorNoParam("vapi.bindings.typeconverter.set.invalid"))
			return err
		}
		result.Add(g.outValue)
	}
	g.inValue = inValue
	g.outValue = result
	return nil
}

func (g *GolangToVapiDataValueVisitor) unexpectedValueError(expectedType string) []error {
	var actualType string
	if g.inValue == nil {
		actualType = "nil"
	} else{
		actualType = reflect.TypeOf(g.inValue).String()
	}
	var args = map[string]string{
		"expectedType": expectedType,
		"actualType": actualType  }
	return []error{l10n.NewRuntimeError("vapi.bindings.typeconverter.unexpected.runtime.value", args)}
}
