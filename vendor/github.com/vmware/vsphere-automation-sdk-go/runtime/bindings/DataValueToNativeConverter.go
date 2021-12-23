/* Copyright © 2019-2020 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package bindings

import (
	"encoding/base64"
	"encoding/json"
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

// DataValueToNativeConverter converts DataValue to golang native
type DataValueToNativeConverter struct {
	// Value which will be converted to a golang object.
	inValue data.DataValue

	//output of the visitor
	outValue interface{}

	typeConverter *TypeConverter
}

func NewDataValueToNativeConverter(input data.DataValue, typeConverter *TypeConverter) *DataValueToNativeConverter {
	return &DataValueToNativeConverter{inValue: input, typeConverter: typeConverter}
}

func (v *DataValueToNativeConverter) OutputValue() interface{} {
	return v.outValue
}

func (v *DataValueToNativeConverter) visit(bindingType BindingType) []error {
	return v.visitInternal(bindingType, false)

}

func (v *DataValueToNativeConverter) setOutValue(optional bool, x reflect.Value) {
	if optional {
		v.outValue = x.Interface()
	} else {
		v.outValue = x.Elem().Interface()
	}
}
func (v *DataValueToNativeConverter) visitInteger(optional bool) []error {
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

func (v *DataValueToNativeConverter) visitDouble(optional bool) []error {
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

func (v *DataValueToNativeConverter) visitString(optional bool) []error {
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

func (v *DataValueToNativeConverter) visitBoolean(optional bool) []error {
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

func (v *DataValueToNativeConverter) visitOptional(bindingType BindingType) []error {
	if optionalValue, ok := v.inValue.(*data.OptionalValue); ok {
		v.inValue = optionalValue.Value()
	} else {
		log.Debugf("Tolerating absence of optional value")
	}
	return v.visitInternal(bindingType, true)
}

func (v *DataValueToNativeConverter) visitBlob() []error {
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
func (v *DataValueToNativeConverter) visitStruct(structType StructType, optional bool) []error {
	if v.inValue == nil {
		v.outValue = zeroPtr(structType.BindingStruct())
		return nil
	}

	var x = reflect.New(structType.BindingStruct())
	err := v.setStructType(structType, x)
	if err != nil {
		return err
	}
	v.setOutValue(optional, x)
	return nil
}

func (v *DataValueToNativeConverter) visitListType(listType ListType) []error {
	if v.inValue == nil {
		v.outValue = zeroPtr(listType.bindingStruct)
		return nil
	}
	s := listType.bindingStruct
	if listValue, ok := v.inValue.(*data.ListValue); ok {
		slice := reflect.MakeSlice(s, len(listValue.List()), len(listValue.List()))
		x := reflect.New(slice.Type())
		x.Elem().Set(slice)
		err := v.setListType(listType, x)
		if err != nil {
			return err
		}
		v.outValue = slice.Interface()
		return nil
	}
	return v.unexpectedTypeError(listType.Type().String())
}

func (v *DataValueToNativeConverter) visitMapType(mapType MapType) []error {
	if v.inValue == nil {
		v.outValue = zeroPtr(mapType.bindingStruct)
		return nil
	}
	if structValue, ok := v.inValue.(*data.StructValue); ok {
		s := mapType.bindingStruct
		x := reflect.MakeMapWithSize(s, len(structValue.Fields()))
		err := v.setMapStructType(mapType, x)
		if err != nil {
			return err
		}
		v.outValue = x.Interface()
		return nil
	}
	if listValue, ok := v.inValue.(*data.ListValue); ok {
		s := mapType.bindingStruct
		x := reflect.MakeMapWithSize(s, len(listValue.List()))
		err := v.setMapListType(mapType, x)
		if err != nil {
			return err
		}
		v.outValue = x.Interface()
		return nil
	}
	return v.unexpectedTypeError(mapType.Type().String())
}

func (v *DataValueToNativeConverter) visitEnumType(enumType EnumType, optional bool) []error {
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

func (v *DataValueToNativeConverter) visitSetType(setType SetType) []error {
	inValue := v.inValue
	if inValue == nil {
		v.outValue = zeroPtr(setType.bindingStruct)
		return nil
	}
	if listValue, ok := v.inValue.(*data.ListValue); ok {
		var modifiedListVal = data.NewListValue()
		for _, element := range listValue.List() {
			var field = make(map[string]data.DataValue)
			field[lib.MAP_KEY_FIELD] = element
			field[lib.MAP_VALUE_FIELD] = data.NewBooleanValue(true)
			var structValue = data.NewStructValue(lib.MAP_ENTRY, field)
			modifiedListVal.Add(structValue)
		}
		s := setType.bindingStruct
		x := reflect.MakeMapWithSize(s, len(listValue.List()))
		bType := NewMapType(setType.ElementType(), NewBooleanType(), s)
		v.inValue = modifiedListVal
		err := v.setMapListType(bType, x)
		if err != nil {
			return err
		}
		v.outValue = x.Interface()
		return nil
	}
	return v.unexpectedTypeError(setType.Type().String())
}

func (v *DataValueToNativeConverter) visitErrorType(errorType ErrorType, optional bool) []error {
	if v.inValue == nil {
		if optional {
			v.outValue = zeroPtr(errorType.BindingStruct())
			return nil
		} else {
			return v.unexpectedTypeError(errorType.Type().String())
		}
	}
	var x = reflect.New(errorType.bindingStruct)
	err := v.setErrorType(errorType, x)
	if err != nil {
		return err
	}
	v.setOutValue(optional, x)
	return nil
}

func (v *DataValueToNativeConverter) visitSecretType(optional bool) []error {
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

func (v *DataValueToNativeConverter) visitDateTimeType(optional bool) []error {
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

func (v *DataValueToNativeConverter) visitUriType(optional bool) []error {
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

func (v *DataValueToNativeConverter) visitAnyErrorType() []error {
	if v.inValue == nil {
		var x *data.ErrorValue
		v.outValue = x
		return nil
	}
	if errorValue, ok := v.inValue.(*data.ErrorValue); ok {
		v.outValue = errorValue
		return nil
	}
	return v.unexpectedTypeError(data.ERROR.String())
}

func (v *DataValueToNativeConverter) visitDynamicStructure(structType DynamicStructType) []error {
	if v.inValue == nil {
		var x *data.StructValue
		v.outValue = x
		return nil
	}
	if structVal, ok := v.inValue.(*data.StructValue); ok {
		v.outValue = structVal
		messages := structType.Validate(structVal)
		if messages != nil {
			return messages
		}
		return nil
	}
	return v.unexpectedTypeError(structType.Type().String())
}
func (v *DataValueToNativeConverter) visitInternal(bindingType BindingType, optional bool) []error {
	concreteBindingType := reflect.TypeOf(bindingType)

	// throw error if not optional and inValue is nil
	if v.inValue == nil &&
		!optional &&
		concreteBindingType != OpaqueBindingType &&
		concreteBindingType != VoidBindingType &&
		concreteBindingType != OptionalBindingType {
		return v.unexpectedTypeError(bindingType.Type().String())
	}

	switch concreteBindingType {
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
		return v.visitListType(bindingType.(ListType))
	case MapBindingType:
		return v.visitMapType(bindingType.(MapType))
	case IdBindingType:
		return v.visitString(optional)
	case EnumBindingType:
		return v.visitEnumType(bindingType.(EnumType), optional)
	case SetBindingType:
		return v.visitSetType(bindingType.(SetType))
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
		return v.visitAnyErrorType()
	case DynamicStructBindingType:
		return v.visitDynamicStructure(bindingType.(DynamicStructType))
	case ReferenceBindingType:
		return v.visitInternal(bindingType.(ReferenceType).Resolve().(StructType), optional)
	default:
		return []error{l10n.NewRuntimeError("vapi.bindings.typeconverter.invalid.type",
			map[string]string{"bindingType": concreteBindingType.String()})}
	}
}

func (v *DataValueToNativeConverter) setIntegerType(value reflect.Value) []error {
	if v.setZeroValue(value) {
		return nil
	}

	if integerValue, ok := v.inValue.(*data.IntegerValue); ok {
		x := integerValue.Value()
		value.Elem().Set(reflect.ValueOf(x))
		return nil
	} else if strValue, ok := v.inValue.(*data.StringValue); ok {
		// If string value passed check if it can be converted to integer
		var val int64
		err := json.Unmarshal([]byte(strValue.Value()), &val)
		if err == nil {
			value.Elem().Set(reflect.ValueOf(val))
			return nil
		}
	}
	return v.unexpectedTypeError(data.INTEGER.String())
}

func (v *DataValueToNativeConverter) setDoubleType(value reflect.Value) []error {
	if v.setZeroValue(value) {
		return nil
	}
	val := value.Elem()
	if doubleValue, ok := v.inValue.(*data.DoubleValue); ok {
		x := doubleValue.Value()
		val.Set(reflect.ValueOf(x))
		return nil
	} else if intValue, ok := v.inValue.(*data.IntegerValue); ok {
		val.Set(reflect.ValueOf(float64(intValue.Value())))
		return nil
	} else if stringValue, ok := v.inValue.(*data.StringValue); ok {
		var val float64
		err := json.Unmarshal([]byte(stringValue.Value()), &val)
		if err == nil {
			value.Elem().Set(reflect.ValueOf(val))
			return nil
		}
	}
	return v.unexpectedTypeError(data.DOUBLE.String())

}

func (v *DataValueToNativeConverter) setBooleanType(value reflect.Value) []error {
	if v.setZeroValue(value) {
		return nil
	}
	val := value.Elem()
	if booleanValue, ok := v.inValue.(*data.BooleanValue); ok {
		x := booleanValue.Value()
		val.Set(reflect.ValueOf(x))
		return nil
	} else if stringValue, ok := v.inValue.(*data.StringValue); ok {
		var val bool
		err := json.Unmarshal([]byte(stringValue.Value()), &val)
		if err == nil {
			value.Elem().Set(reflect.ValueOf(val))
			return nil
		}
	}
	return v.unexpectedTypeError(data.BOOLEAN.String())

}

func (v *DataValueToNativeConverter) setStringType(value reflect.Value) []error {
	if v.setZeroValue(value) {
		return nil
	}
	val := value.Elem()
	if stringValue, ok := v.inValue.(*data.StringValue); ok {
		x := stringValue.Value()
		val.Set(reflect.ValueOf(x))
		return nil
	}
	return v.unexpectedTypeError(data.STRING.String())
}

func (v *DataValueToNativeConverter) setSecretType(value reflect.Value) []error {
	if v.setZeroValue(value) {
		return nil
	}
	val := value.Elem()
	if secretValue, ok := v.inValue.(*data.SecretValue); ok {
		x := secretValue.Value()
		val.Set(reflect.ValueOf(x))
		return nil
	}
	log.Debug("Expected SecretValue. Checking for StringValue")
	if stringValue, ok := v.inValue.(*data.StringValue); ok {
		x := stringValue.Value()
		val.Set(reflect.ValueOf(x))
		log.Debug("Expected SecretValue. Found StringValue")
		return nil
	}
	return v.unexpectedTypeError(data.SECRET.String())
}
func (v *DataValueToNativeConverter) setBlobType(value reflect.Value) []error {
	if blobValue, ok := v.inValue.(*data.BlobValue); ok {
		value.Elem().Set(reflect.ValueOf(blobValue.Value()))
		return nil
	}
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
	return v.unexpectedTypeError(data.BLOB.String())
}

func (v *DataValueToNativeConverter) setDateTimeType(value reflect.Value) []error {
	if v.setZeroValue(value) {
		return nil
	}
	val := value.Elem()
	if stringValue, ok := v.inValue.(*data.StringValue); ok {
		x := stringValue.Value()
		datetimeLayout := RFC3339Nano_DATETIME_LAYOUT
		datetime, err := time.Parse(datetimeLayout, x)
		if err != nil {
			return []error{l10n.NewRuntimeError("vapi.bindings.typeconverter.datetime.invalid",
				map[string]string{"dateTime": x, "vapiFormat": datetimeLayout, "errorMessage": err.Error()})}
		}
		restDatetimeStr := datetime.UTC().Format(VAPI_DATETIME_LAYOUT)

		// In above line, we are converting Format of time which is in RFC3339Nano to VAPI layout but it gives back a string.
		// so we need to parse it from resulted string restDatetimeStr, Now we are already aware that restDatetimeStr is in
		// format compliant with VAPI_DATETIME_LAYOUT hence error handling is not required.
		t, _ := time.Parse(VAPI_DATETIME_LAYOUT, restDatetimeStr)
		val.Set(reflect.ValueOf(t))
		return nil
	}
	return v.unexpectedTypeError(data.STRING.String())
}

func (v *DataValueToNativeConverter) setURIType(value reflect.Value) []error {
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
	return v.unexpectedTypeError(data.STRING.String())
}

func (v *DataValueToNativeConverter) setEnumType(value reflect.Value) []error {
	value = value.Elem()
	if stringValue, ok := v.inValue.(*data.StringValue); ok {
		value.SetString(stringValue.Value())
		return nil
	}
	return v.unexpectedTypeError(data.STRING.String())
}

func (v *DataValueToNativeConverter) setStructType(typ StructType, outputPtr reflect.Value) []error {
	if v.setZeroValue(outputPtr) {
		return nil
	}
	output := outputPtr.Elem()
	var structValue *data.StructValue = nil
	if structVal, ok := v.inValue.(*data.StructValue); !ok {
		return v.unexpectedTypeError(data.STRUCTURE.String())
	} else {
		structValue = structVal
	}
	for _, fieldName := range typ.FieldNames() {
		var field = output.FieldByName(typ.canonicalFieldMap[fieldName])
		var bindingType = typ.Field(fieldName)
		var err error
		v.inValue, err = structValue.Field(fieldName)
		if err != nil {
			// if the field is optional, absence of the field in DataValue should be tolerated.
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
	messages := typ.Validate(structValue)
	if messages != nil {
		return messages
	}
	v.inValue = structValue
	return nil
}

func (v *DataValueToNativeConverter) setErrorType(typ ErrorType, outputPtr reflect.Value) []error {
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
		return v.unexpectedTypeError(typ.Type().String())
	}
	for _, fieldName := range typ.FieldNames() {
		var field = output.FieldByName(typ.canonicalFieldMap[fieldName])
		var bindingType = typ.Field(fieldName)
		var err error
		v.inValue, err = errorValue.Field(fieldName)
		if err != nil {
			// if the field is optional, absence of the field in DataValue should be tolerated.
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

func (v *DataValueToNativeConverter) setListType(listType ListType, slice reflect.Value) []error {
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
		return v.unexpectedTypeError(listType.Type().String())
	}
	v.inValue = inValue
	v.outValue = slice
	return nil
}

func (v *DataValueToNativeConverter) setMapStructType(mapType MapType, result reflect.Value) []error {

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
	return v.unexpectedTypeError(data.STRUCTURE.String())
}

func (v *DataValueToNativeConverter) setMapListType(mapType MapType, result reflect.Value) []error {

	if listValue, ok := v.inValue.(*data.ListValue); ok {
		inValue := v.inValue
		for _, listValElem := range listValue.List() {
			if structVal, ok := listValElem.(*data.StructValue); ok {
				//process key
				keyVal, _ := structVal.Field(lib.MAP_KEY_FIELD)
				v.inValue = keyVal
				err := v.visit(mapType.KeyType)
				if err != nil {
					err = append(err, l10n.NewRuntimeErrorNoParam("vapi.bindings.typeconverter.dict.key.invalid"))
					return err
				}
				mKey := reflect.ValueOf(v.outValue)

				//process value
				mapVal, _ := structVal.Field(lib.MAP_VALUE_FIELD)
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
		}
		v.inValue = inValue
		return nil
	}
	return v.unexpectedTypeError(data.LIST.String())
}

func (v *DataValueToNativeConverter) setZeroValue(value reflect.Value) bool {
	if v.inValue == nil {
		value.Set(reflect.Zero(reflect.TypeOf(value)))
		return true
	}
	return false
}

func (v *DataValueToNativeConverter) unexpectedTypeError(expectedType string) []error {
	var actualType string
	if v.inValue == nil {
		actualType = "nil"
	} else {
		actualType = v.inValue.Type().String()
	}

	var args = map[string]string{
		"expectedType": expectedType,
		"actualType":   actualType}
	return []error{l10n.NewRuntimeError("vapi.bindings.typeconverter.unexpected.runtime.value", args)}
}

// zeroPtr returns an interface whose value is nil (zero value of Ptr) and whose type is pointer to type specified by
// input. If the kind is Slice, Map or Array, returned interface value is nil and type is type specified by input.
func zeroPtr(input reflect.Type) interface{} {
	switch input.Kind() {
	case reflect.Slice, reflect.Map, reflect.Array, reflect.Ptr:
		x := reflect.New(input)
		return x.Elem().Interface()
	default:
		ptrReflectType := reflect.PtrTo(input)
		x := reflect.New(ptrReflectType)
		return x.Elem().Interface()
	}
}
