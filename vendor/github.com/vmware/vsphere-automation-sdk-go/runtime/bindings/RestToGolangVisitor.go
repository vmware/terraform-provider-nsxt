/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package bindings

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"time"

	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/l10n"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/lib"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/log"
)

// Visitor to convert from RestNative DataValue to Golang native value
type RestToGolangVisitor struct {
	// Value which will be converted to a golang object.
	inValue data.DataValue

	//output of the visitor
	outValue interface{}

	typeConverter *TypeConverter
}

func NewRestToGolangVisitor(input data.DataValue, typeConverter *TypeConverter) *RestToGolangVisitor {
	return &RestToGolangVisitor{inValue: input, typeConverter: typeConverter}
}

func (v *RestToGolangVisitor) OutputValue() interface{} {
	return v.outValue
}

func (v *RestToGolangVisitor) visit(bindingType BindingType) []error {
	return v.visitInternal(bindingType, false)

}

func (v *RestToGolangVisitor) setOutValue(optional bool, x reflect.Value) {
	if optional {
		v.outValue = x.Interface()
	} else {
		v.outValue = x.Elem().Interface()
	}
}
func (v *RestToGolangVisitor) visitInteger(optional bool) []error {
	if v.inValue == nil {
		var x *int64 = nil
		v.outValue = x
		return nil
	}
	x := reflect.New(reflect.TypeOf(int64(0)))
	err := v.setIntegerType(x)
	if err != nil {
		return err
	}
	v.setOutValue(optional, x)
	return nil
}

func (v *RestToGolangVisitor) visitDouble(optional bool) []error {
	if v.inValue == nil {
		var x *float64 = nil
		v.outValue = x
		return nil
	}
	x := reflect.New(reflect.TypeOf(float64(0)))
	err := v.setDoubleType(x)
	if err != nil {
		return err
	}
	v.setOutValue(optional, x)
	return nil
}

func (v *RestToGolangVisitor) visitString(optional bool) []error {
	if v.inValue == nil {
		var x *string = nil
		v.outValue = x
		return nil
	}
	x := reflect.New(reflect.TypeOf(""))
	err := v.setStringType(x)
	if err != nil {
		return err
	}
	v.setOutValue(optional, x)
	return nil
}

func (v *RestToGolangVisitor) visitBoolean(optional bool) []error {
	if v.inValue == nil {
		var x *bool = nil
		v.outValue = x
		return nil
	}
	x := reflect.New(reflect.TypeOf(true))
	err := v.setBooleanType(x)
	if err != nil {
		return err
	}
	v.setOutValue(optional, x)
	return nil
}

func (v *RestToGolangVisitor) visitOptional(bindingType BindingType) []error {
	if optionalValue, ok := v.inValue.(*data.OptionalValue); ok {
		v.inValue = optionalValue.Value()
	} else if v.typeConverter.permissive {
		log.Debugf("Expected OptionalValue but found %s", reflect.TypeOf(v.inValue))
		log.Debugf("Tolerating absence of optional value")
	} else {
		return v.unexpectedValueError("OptionalValue")
	}
	return v.visitInternal(bindingType, true)
}

func (v *RestToGolangVisitor) visitBlob() []error {
	if v.inValue == nil {
		var nilSlice []byte = nil
		v.outValue = nilSlice
		return nil
	}
	//slice is not nil. initialize an empty slice.
	// nil slice, empty slice are different.
	slice := reflect.MakeSlice(reflect.TypeOf([]uint8{}), 0, 0)
	x := reflect.New(slice.Type())
	err := v.setBlobType(x)
	if err != nil {
		return err
	}
	v.outValue = x.Elem().Interface()
	return nil
}
func (v *RestToGolangVisitor) visitStruct(structType StructType, optional bool) []error {
	if v.inValue == nil {
		v.outValue = zeroPtr(structType.BindingStruct())
		return nil
	}

	var x = reflect.New(structType.BindingStruct())
	err := v.setStructType(structType, optional, x)
	if err != nil {
		return err
	}
	v.setOutValue(optional, x)
	return nil
}

func (v *RestToGolangVisitor) visitListType(listType ListType, optional bool) []error {
	if v.inValue == nil {
		v.outValue = zeroPtr(listType.bindingStruct)
		return nil
	}
	s := listType.bindingStruct
	if listValue, ok := v.inValue.(*data.ListValue); ok {
		slice := reflect.MakeSlice(s, len(listValue.List()), len(listValue.List()))
		x := reflect.New(slice.Type())
		x.Elem().Set(slice)
		err := v.setListType(listType, optional, x)
		if err != nil {
			return err
		}
		v.outValue = slice.Interface()
		return nil
	}
	return v.unexpectedValueError("ListValue")
}

func (v *RestToGolangVisitor) visitMapType(mapType MapType, optional bool) []error {
	if v.inValue == nil {
		v.outValue = zeroPtr(mapType.bindingStruct)
		return nil
	}
	if structValue, ok := v.inValue.(*data.StructValue); ok {
		s := mapType.bindingStruct
		x := reflect.MakeMapWithSize(s, len(structValue.Fields()))
		err := v.setMapType(mapType, optional, x)
		if err != nil {
			return err
		}
		v.outValue = x.Interface()
		return nil
	}
	return v.unexpectedValueError("ListValue")
}

func (v *RestToGolangVisitor) visitEnumType(enumType EnumType, optional bool) []error {
	e := enumType.bindingStruct
	if v.inValue == nil {
		v.outValue = zeroPtr(e)
		return nil
	}
	x := reflect.New(e)
	err := v.setEnumType(x)
	if err != nil {
		return err
	}
	v.setOutValue(optional, x)
	return nil
}

func (v *RestToGolangVisitor) visitSetType(setType SetType, optional bool) []error {
	inValue := v.inValue
	if inValue == nil {
		v.outValue = zeroPtr(setType.bindingStruct)
		return nil
	}
	if listValue, ok := v.inValue.(*data.ListValue); ok {
		var structVal = data.NewStructValue(lib.MAP_ENTRY, nil)
		for _, element := range listValue.List() {
			setBindingType := reflect.TypeOf(setType.ElementType())
			if setBindingType == IntegerBindingType {
				intVal := element.(*data.IntegerValue).Value()
				structVal.SetField(strconv.FormatInt(intVal, 10), data.NewBooleanValue(true))
			} else {
				structVal.SetField(element.(*data.StringValue).Value(), data.NewBooleanValue(true))
			}
		}
		s := setType.bindingStruct
		x := reflect.MakeMapWithSize(s, len(listValue.List()))
		bType := NewMapType(setType.ElementType(), NewBooleanType(), s)
		v.inValue = structVal
		err := v.setMapType(bType, optional, x)
		if err != nil {
			return err
		}
		v.outValue = x.Interface()
		return nil
	}
	return v.unexpectedValueError("ListValue")
}

func (v *RestToGolangVisitor) visitErrorType(errorType ErrorType, optional bool) []error {
	if v.inValue == nil {
		v.outValue = zeroPtr(errorType.BindingStruct())
		return nil
	}
	var x = reflect.New(errorType.bindingStruct)
	err := v.setErrorType(errorType, optional, x)
	if err != nil {
		return err
	}
	v.setOutValue(optional, x)
	return nil
}

func (v *RestToGolangVisitor) visitSecretType(optional bool) []error {
	if v.inValue == nil {
		var x *string = nil
		v.outValue = x
		return nil
	}
	var x = reflect.New(reflect.TypeOf(""))
	err := v.setSecretType(x)
	if err != nil {
		return err
	}
	v.setOutValue(optional, x)
	return nil
}

func (v *RestToGolangVisitor) visitDateTimeType(optional bool) []error {
	if v.inValue == nil {
		var x *time.Time = nil
		v.outValue = x
		return nil
	}
	var x = reflect.New(reflect.TypeOf(time.Time{}))
	err := v.setDateTimeType(x)
	if err != nil {
		return err
	}
	v.setOutValue(optional, x)
	return nil
}

func (v *RestToGolangVisitor) visitUriType(optional bool) []error {
	if v.inValue == nil {
		var x *url.URL = nil
		v.outValue = x
		return nil
	}
	var x = reflect.New(reflect.TypeOf(url.URL{}))
	err := v.setURIType(x)
	if err != nil {
		return err
	}
	v.setOutValue(optional, x)
	return nil
}

func (v *RestToGolangVisitor) visitAnyErrorType(optional bool) []error {
	if v.inValue == nil {
		var x *data.ErrorValue
		v.outValue = x
		return nil
	}
	if errorValue, ok := v.inValue.(*data.ErrorValue); ok {
		v.outValue = errorValue
		return nil
	}
	return v.unexpectedValueError("ErrorValue")
}

func (v *RestToGolangVisitor) visitDynamicStructure(structType DynamicStructType, optional bool) []error {
	if v.inValue == nil {
		var x *data.StructValue
		v.outValue = x
		return nil
	}
	if structVal, ok := v.inValue.(*data.StructValue); ok {
		v.outValue = structVal
		msgs := structType.Validate(structVal)
		if msgs != nil {
			return msgs
		}
		return nil
	}
	return v.unexpectedValueError("StructValue")
}
func (v *RestToGolangVisitor) visitInternal(bindingType BindingType, optional bool) []error {
	switch reflect.TypeOf(bindingType) {
	case BlobBindingType:
		return v.visitBlob()
	case IntegerBindingType:
		return v.visitInteger(optional)
	case DoubleBindingType:
		return v.visitDouble(optional)
	case StringBindingType:
		return v.visitString(optional)
	case BooleanBindingType:
		return v.visitBoolean(optional)
	case OptionalBindingType:
		return v.visitOptional(bindingType.(OptionalType).ElementType())
	case OpaqueBindingType:
		v.outValue = v.inValue
		return nil
	case StructBindingType:
		return v.visitStruct(bindingType.(StructType), optional)
	case ListBindingType:
		return v.visitListType(bindingType.(ListType), optional)
	case MapBindingType:
		return v.visitMapType(bindingType.(MapType), optional)
	case IdBindingType:
		return v.visitString(optional)
	case EnumBindingType:
		return v.visitEnumType(bindingType.(EnumType), optional)
	case SetBindingType:
		return v.visitSetType(bindingType.(SetType), optional)
	case ErrorBindingType:
		return v.visitErrorType(bindingType.(ErrorType), optional)
	case SecretBindingType:
		return v.visitSecretType(optional)
	case DateTimeBindingType:
		return v.visitDateTimeType(optional)
	case UriBindingType:
		return v.visitUriType(optional)
	case VoidBindingType:
		v.outValue = nil
		return nil
	case AnyErrorBindingType:
		return v.visitAnyErrorType(optional)
	case DynamicStructBindingType:
		return v.visitDynamicStructure(bindingType.(DynamicStructType), optional)
	case ReferenceBindingType:
		return v.visitInternal(bindingType.(ReferenceType).Resolve().(StructType), optional)
	default:
		return []error{l10n.NewRuntimeError("vapi.bindings.typeconverter.invalid.type",
			map[string]string{"bindingType": reflect.TypeOf(bindingType).String()})}
	}
}

func (v *RestToGolangVisitor) setIntegerType(value reflect.Value) []error {
	if v.setZeroValue(value) {
		return nil
	}

	if integerValue, ok := v.inValue.(*data.IntegerValue); ok {
		x := integerValue.Value()
		value.Elem().Set(reflect.ValueOf(x))
		return nil
	} else if strValue, ok := v.inValue.(*data.StringValue); ok {
		// If string value passed check if it can be converted to integer
		i, err := strconv.ParseInt(strValue.Value(), 10, 64)
		if err == nil {
			value.Elem().Set(reflect.ValueOf(i))
			return nil
		}
	}
	return v.unexpectedValueError("IntegerValue")
}

func (v *RestToGolangVisitor) setDoubleType(value reflect.Value) []error {
	if v.setZeroValue(value) {
		return nil
	}
	val := value.Elem()
	if doubleValue, ok := v.inValue.(*data.DoubleValue); ok {
		x := doubleValue.Value()
		val.Set(reflect.ValueOf(x))
		return nil
	}
	return v.unexpectedValueError("DoubleValue")

}

func (v *RestToGolangVisitor) setBooleanType(value reflect.Value) []error {
	if v.setZeroValue(value) {
		return nil
	}
	val := value.Elem()
	if booleanValue, ok := v.inValue.(*data.BooleanValue); ok {
		x := booleanValue.Value()
		val.Set(reflect.ValueOf(x))
		return nil
	}
	return v.unexpectedValueError("BooleanValue")

}

func (v *RestToGolangVisitor) setStringType(value reflect.Value) []error {
	if v.setZeroValue(value) {
		return nil
	}
	val := value.Elem()
	if stringValue, ok := v.inValue.(*data.StringValue); ok {
		x := stringValue.Value()
		val.Set(reflect.ValueOf(x))
		return nil
	}
	return v.unexpectedValueError("StringValue")
}

func (v *RestToGolangVisitor) setSecretType(value reflect.Value) []error {
	if v.setZeroValue(value) {
		return nil
	}
	val := value.Elem()
	if secretValue, ok := v.inValue.(*data.SecretValue); ok {
		x := secretValue.Value()
		val.Set(reflect.ValueOf(x))
		return nil
	}
	if v.typeConverter.permissive {
		log.Debug("Expected SecretValue. Checking for StringValue in permissive mode")
		if stringValue, ok := v.inValue.(*data.StringValue); ok {
			x := stringValue.Value()
			val.Set(reflect.ValueOf(x))
			log.Debug("Expected SecretValue. Found StringValue in permissive mode")
			return nil
		}
		return v.unexpectedValueError("SecretValue or StringValue")
	}
	return v.unexpectedValueError("SecretValue")
}
func (v *RestToGolangVisitor) setBlobType(value reflect.Value) []error {
	if blobValue, ok := v.inValue.(*data.BlobValue); ok {
		value.Elem().Set(reflect.ValueOf(blobValue.Value()))
		return nil
	}
	if v.typeConverter.permissive {
		if stringValue, ok := v.inValue.(*data.StringValue); ok {
			log.Debug("Expected BlobValue but found StringValue instead.")
			//decode base 64 encoded string.
			decodedString, decodeErr := base64.StdEncoding.DecodeString(stringValue.Value())
			if decodeErr != nil {
				var args = map[string]string{
					"errMsg": decodeErr.Error()}
				return []error{l10n.NewRuntimeError("vapi.bindings.typeconverter.blob.base64.decode.error", args)}
			}
			value.Elem().Set(reflect.ValueOf(decodedString))
			return nil
		}
		return v.unexpectedValueError("BlobValue or StringValue")
	}
	return v.unexpectedValueError("BlobValue")
}

func (v *RestToGolangVisitor) setDateTimeType(value reflect.Value) []error {
	if v.setZeroValue(value) {
		return nil
	}
	val := value.Elem()
	if stringValue, ok := v.inValue.(*data.StringValue); ok {
		x := stringValue.Value()
		datetime_layout := RFC3339Nano_DATETIME_LAYOUT
		datetime, err := time.Parse(datetime_layout, x)
		if err != nil {
			return []error{l10n.NewRuntimeError("vapi.bindings.typeconverter.datetime.invalid",
				map[string]string{"dateTime": x, "vapiFormat": datetime_layout, "errorMessage": err.Error()})}
		}
		restDatetimeStr := datetime.UTC().Format(VAPI_DATETIME_LAYOUT)

		// In above line, we are converting Format of time which is in RFC3339Nano to VAPI layout but it gives back a string.
		// so we need to parse it from resulted string restDatetimeStr, Now we are already aware that restDatetimeStr is in
		// format compliant with VAPI_DATEIME_LAYOUT hence error handling is not required.
		t, _ := time.Parse(VAPI_DATETIME_LAYOUT, restDatetimeStr)
		val.Set(reflect.ValueOf(t))
		return nil
	}
	return v.unexpectedValueError("StringValue")
}

func (v *RestToGolangVisitor) setURIType(value reflect.Value) []error {
	if v.setZeroValue(value) {
		return nil
	}
	if stringValue, ok := v.inValue.(*data.StringValue); ok {
		x := stringValue.Value()
		u, err := url.Parse(x)
		if err != nil {
			return []error{l10n.NewRuntimeError("vapi.bindings.typeconverter.uri.invalid",
				map[string]string{"uriValue": x, "errorMessage": err.Error()})}
		}
		value.Elem().Set(reflect.ValueOf(*u))
		return nil
	}
	return v.unexpectedValueError("StringValue")
}

func (v *RestToGolangVisitor) setEnumType(value reflect.Value) []error {
	value = value.Elem()
	if stringValue, ok := v.inValue.(*data.StringValue); ok {
		value.SetString(stringValue.Value())
		return nil
	}
	return v.unexpectedValueError("StringValue")
}

func (v *RestToGolangVisitor) setStructType(typ StructType, optional bool, outputPtr reflect.Value) []error {
	if v.setZeroValue(outputPtr) {
		return nil
	}
	output := outputPtr.Elem()
	var structValue *data.StructValue = nil
	if structVal, ok := v.inValue.(*data.StructValue); !ok {
		return v.unexpectedValueError("StructValue")
	} else {
		structValue = structVal
	}
	for _, fieldName := range typ.FieldNames() {
		var field = output.FieldByName(typ.canonicalFieldMap[fieldName])
		var bindingType = typ.Field(fieldName)
		var err error
		v.inValue, err = structValue.Field(fieldName)
		if err != nil {
			// if the field is optional, absence of the field in datavalue should be tolerated.
			if _, ok := bindingType.(OptionalType); !ok {
				// return error only if field is not optional
				return []error{l10n.NewRuntimeError("vapi.bindings.typeconverter.struct.field.missing",
					map[string]string{"fieldName": fieldName, "structName": typ.name})}
			}
		}
		if field.IsValid() {
			if field.CanSet() {
				err := v.visit(bindingType)
				if err != nil {
					err = append(err, l10n.NewRuntimeError("vapi.bindings.typeconverter.struct.field.invalid",
						map[string]string{"fieldName": fieldName, "structName": typ.name}))
					return err
				}
				if v.outValue != nil {
					fieldVal := reflect.ValueOf(v.outValue)
					field.Set(fieldVal)
				}
			}
			//error if cannot set?
		} else {
			// is this error right?
			return []error{l10n.NewRuntimeError("vapi.bindings.typeconverter.struct.field.invalid",
				map[string]string{"fieldName": fieldName, "structName": typ.name})}
		}
	}
	msgs := typ.Validate(structValue)
	if msgs != nil {
		return msgs
	}
	v.inValue = structValue
	return nil
}

func (v *RestToGolangVisitor) setErrorType(typ ErrorType, optional bool, outputPtr reflect.Value) []error {
	if v.setZeroValue(outputPtr) {
		return nil
	}
	output := outputPtr.Elem()
	var errorValue *data.ErrorValue = nil

	if errorDataValue, ok := v.inValue.(*data.ErrorValue); ok {
		errorValue = errorDataValue
	} else if structValue, ok := v.inValue.(*data.StructValue); ok {
		errorValue = data.NewErrorValue(structValue.Name(), structValue.Fields())
	} else {
		return v.unexpectedValueError("ErrorValue")
	}
	for _, fieldName := range typ.FieldNames() {
		var field = output.FieldByName(typ.canonicalFieldMap[fieldName])
		var bindingType = typ.Field(fieldName)
		var err error
		v.inValue, err = errorValue.Field(fieldName)
		if err != nil {
			// if the field is optional, absence of the field in datavalue should be tolerated.
			if _, ok := bindingType.(OptionalType); !ok {
				// return error only if field is not optional
				return []error{l10n.NewRuntimeError("vapi.bindings.typeconverter.struct.field.missing",
					map[string]string{"fieldName": fieldName, "structName": typ.name})}
			}
		}
		if field.IsValid() {
			if field.CanSet() {
				err := v.visit(bindingType)
				if err != nil {
					err = append(err, l10n.NewRuntimeError("vapi.bindings.typeconverter.struct.field.invalid",
						map[string]string{"fieldName": fieldName, "structName": typ.name}))
					return err
				}
				fieldVal := reflect.ValueOf(v.outValue)
				field.Set(fieldVal)
			}
			//error if cannot set?
		} else {
			// is this error right?
			return []error{l10n.NewRuntimeError("vapi.bindings.typeconverter.struct.field.invalid",
				map[string]string{"fieldName": fieldName, "structName": typ.name})}
		}
	}
	v.inValue = errorValue
	return nil
}

func (v *RestToGolangVisitor) setListType(listType ListType, optional bool, slice reflect.Value) []error {
	//https://play.golang.org/p/0aB01KBniI
	//https://stackoverflow.com/questions/25384640/why-golang-reflect-makeslice-returns-un-addressable-value
	slice = slice.Elem()
	inValue := v.inValue
	if listValue, ok := inValue.(*data.ListValue); ok {
		for index, element := range listValue.List() {
			v.inValue = element
			err := v.visit(listType.ElementType())
			if err != nil {
				err = append(err, l10n.NewRuntimeError("vapi.bindings.typeconverter.list.entry.invalid",
					map[string]string{"index": fmt.Sprintf("%d", index)}))
				return err
			}
			x := reflect.ValueOf(v.outValue)
			if x.IsValid() {
				slice.Index(index).Set(x)
			}
			// invalid is error?
		}
	} else {
		return v.unexpectedValueError("ListValue")
	}
	v.inValue = inValue
	v.outValue = slice
	return nil
}

func (v *RestToGolangVisitor) setMapType(mapType MapType, optional bool, result reflect.Value) []error {

	if structValue, ok := v.inValue.(*data.StructValue); ok {
		//s := mapType.bindingStruct
		inValue := v.inValue
		for _, fieldName := range structValue.FieldNames() {
			keyBindingType := reflect.TypeOf(mapType.KeyType)
			if keyBindingType == IntegerBindingType {
				n, err := strconv.ParseInt(fieldName, 10, 64)
				if err != nil {
					log.Errorf("Error converting string to int64 %s", err)
					return []error{l10n.NewRuntimeError("vapi.bindings.typeconverter.dict.value.invalid",
						map[string]string{"key": fieldName})}
				}
				v.inValue = data.NewIntegerValue(n)
			} else { //StringBindingType, IdBindingType, URIBindingType, EnumBindingType
				v.inValue = data.NewStringValue(fieldName)
			}
			err := v.visit(mapType.KeyType)
			if err != nil {
				err = append(err, l10n.NewRuntimeErrorNoParam("vapi.bindings.typeconverter.dict.key.invalid"))
				return err
			}
			mKey := reflect.ValueOf(v.outValue)

			//process value
			mapVal, _ := structValue.Field(fieldName)
			v.inValue = mapVal
			err = v.visit(mapType.ValueType)
			if err != nil {
				err = append(err, l10n.NewRuntimeError("vapi.bindings.typeconverter.dict.value.invalid",
					map[string]string{"key": mKey.String()}))
				return err
			}
			mVal := reflect.ValueOf(v.outValue)
			result.SetMapIndex(mKey, mVal)

		}
		v.inValue = inValue
		return nil
	}
	return v.unexpectedValueError("StructValue")
}

func (v *RestToGolangVisitor) setZeroValue(value reflect.Value) bool {
	if v.inValue == nil {
		value.Set(reflect.Zero(reflect.TypeOf(value)))
		return true
	}
	return false
}

func (v *RestToGolangVisitor) unexpectedValueError(expectedType string) []error {
	var args = map[string]string{
		"expectedType": expectedType,
		"actualType":   reflect.TypeOf(v.inValue).String()}
	return []error{l10n.NewRuntimeError("vapi.bindings.typeconverter.unexpected.runtime.value", args)}
}
