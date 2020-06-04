/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package cleanjson

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/l10n"
)

// Serializes DataValue to clean json.
type DataValueToJsonEncoder struct {
}

func NewDataValueToJsonEncoder() *DataValueToJsonEncoder {
	return &DataValueToJsonEncoder{}
}

func (d *DataValueToJsonEncoder) Encode(val interface{}) (string, error) {
	marshaller, err := getSerializer(val)
	if err != nil {
		return "", err
	}
	jsonBytes, err := jsonMarshallDisableEscapeHTML(marshaller)
	if err != nil {
		marshallError := l10n.NewRuntimeError("vapi.data.serializers.json.marshall.error",
			map[string]string{"errorMessage": err.Error()})
		return "", marshallError
	}
	return string(jsonBytes), nil
}

func getSerializer(val interface{}) (json.Marshaler, error) {
	switch reflect.TypeOf(val) {
	case data.StructValuePtr:
		return NewStructValueSerializer(val.(*data.StructValue)), nil
	case data.StringValuePtr:
		return NewStringValueSerializer(val.(*data.StringValue)), nil
	case data.IntegerValuePtr:
		return NewIntegerValueSerializer(val.(*data.IntegerValue)), nil
	case data.DoubleValuePtr:
		return NewDoubleValueSerializer(val.(*data.DoubleValue)), nil
	case data.ListValuePtr:
		return NewListValueSerializer(val.(*data.ListValue)), nil
	case data.OptionalValuePtr:
		return NewOptionalValueSerializer(val.(*data.OptionalValue)), nil
	case data.ErrorValuePtr:
		return NewErrorValueSerializer(val.(*data.ErrorValue)), nil
	case data.VoidValuePtr:
		return NewVoidValueSerializer(), nil
	case data.BoolValuePtr:
		return NewBooleanValueSerializer(val.(*data.BooleanValue)), nil
	case data.BlobValuePtr:
		return NewBlobValueSerializer(val.(*data.BlobValue)), nil
	case data.SecretValuePtr:
		return NewSecretValueSerializer(val.(*data.SecretValue)), nil
	default:
		var serializerNotFound = l10n.NewRuntimeError("vapi.data.serializers.rest.datavalue.error",
			map[string]string{"type": fmt.Sprintf("%s", reflect.TypeOf(val))})
		return nil, serializerNotFound
	}
}

type SecretValueSerializer struct {
	secretValue *data.SecretValue
}

func (svs *SecretValueSerializer) MarshalJSON() ([]byte, error) {
	return json.Marshal(svs.secretValue.Value())
}

func NewSecretValueSerializer(value *data.SecretValue) *SecretValueSerializer {
	return &SecretValueSerializer{secretValue: value}
}

type BlobValueSerializer struct {
	blobValue *data.BlobValue
}

func (bvs *BlobValueSerializer) MarshalJSON() ([]byte, error) {
	base64EncodedString := base64.StdEncoding.EncodeToString(bvs.blobValue.Value())
	return json.Marshal(base64EncodedString)
}

func NewBlobValueSerializer(value *data.BlobValue) *BlobValueSerializer {
	return &BlobValueSerializer{blobValue: value}
}

type StructValueSerializer struct {
	structValue *data.StructValue
}

func NewStructValueSerializer(value *data.StructValue) *StructValueSerializer {
	return &StructValueSerializer{structValue: value}
}

func (svs *StructValueSerializer) MarshalJSON() ([]byte, error) {
	var items = make(map[string]interface{})
	for key, val := range svs.structValue.Fields() {
		// Do not serialize empty optional value
		if isEmptyOptionalValue(val) {
			continue
		}

		var err error
		items[key], err = getSerializer(val)
		if err != nil {
			return nil, err
		}
	}
	return json.Marshal(items)
}

type StringValueSerializer struct {
	stringValue *data.StringValue
}

func NewStringValueSerializer(value *data.StringValue) *StringValueSerializer {
	return &StringValueSerializer{stringValue: value}
}

func (svs *StringValueSerializer) MarshalJSON() ([]byte, error) {
	return jsonMarshallDisableEscapeHTML(svs.stringValue.Value())
}

type IntegerValueSerializer struct {
	integerValue *data.IntegerValue
}

func (svs *IntegerValueSerializer) MarshalJSON() ([]byte, error) {
	return json.Marshal(svs.integerValue.Value())
}

func NewIntegerValueSerializer(value *data.IntegerValue) *IntegerValueSerializer {

	return &IntegerValueSerializer{integerValue: value}
}

type DoubleValueSerializer struct {
	doubleValue *data.DoubleValue
}

func (svs *DoubleValueSerializer) MarshalJSON() ([]byte, error) {
	var buf []byte
	splitResult := strings.Split(svs.doubleValue.String(), ".")
	var prec int
	if len(splitResult) > 1 {
		prec = len(splitResult[1])
	} else {
		//setting minimum precision to 1
		prec = 1
	}
	buf = strconv.AppendFloat(buf, svs.doubleValue.Value(), 'f', prec, 64)
	return buf, nil
}

func NewDoubleValueSerializer(value *data.DoubleValue) *DoubleValueSerializer {
	return &DoubleValueSerializer{doubleValue: value}
}

type ListValueSerializer struct {
	listValue *data.ListValue
}

func (lvs *ListValueSerializer) MarshalJSON() ([]byte, error) {
	result := make([]interface{}, 0)
	for _, element := range lvs.listValue.List() {

		var serializer, err = getSerializer(element)
		if err != nil {
			return nil, err
		}
		result = append(result, serializer)
	}
	return json.Marshal(result)
}
func NewListValueSerializer(value *data.ListValue) *ListValueSerializer {
	return &ListValueSerializer{listValue: value}
}

type OptionalValueSerializer struct {
	optionalValue *data.OptionalValue
}

func (ovs *OptionalValueSerializer) MarshalJSON() ([]byte, error) {
	if ovs.optionalValue.IsSet() {
		var optSerializer, err = getSerializer(ovs.optionalValue.Value())
		if err != nil {
			return nil, err
		}
		return json.Marshal(optSerializer)
	} else {
		return json.Marshal(nil)
	}
}

func NewOptionalValueSerializer(value *data.OptionalValue) *OptionalValueSerializer {
	return &OptionalValueSerializer{optionalValue: value}
}

type ErrorValueSerializer struct {
	errorValue *data.ErrorValue
}

func (evs *ErrorValueSerializer) MarshalJSON() ([]byte, error) {
	var items = make(map[string]interface{})
	for key, val := range evs.errorValue.Fields() {
		// Do not serialize empty optional value
		if isEmptyOptionalValue(val) {
			continue
		}
		var err error
		items[key], err = getSerializer(val)
		if err != nil {
			return nil, err
		}
	}
	return json.Marshal(items)
}

func NewErrorValueSerializer(value *data.ErrorValue) *ErrorValueSerializer {
	return &ErrorValueSerializer{errorValue: value}
}

type VoidValueSerializer struct {
}

func (vvs *VoidValueSerializer) MarshalJSON() ([]byte, error) {
	return json.Marshal(nil)
}
func NewVoidValueSerializer() *VoidValueSerializer {
	return &VoidValueSerializer{}
}

type BooleanValueSerializer struct {
	booleanValue *data.BooleanValue
}

func (bvs *BooleanValueSerializer) MarshalJSON() ([]byte, error) {
	return json.Marshal(bvs.booleanValue.Value())
}

func NewBooleanValueSerializer(value *data.BooleanValue) *BooleanValueSerializer {
	return &BooleanValueSerializer{booleanValue: value}
}

func isEmptyOptionalValue(val data.DataValue) bool {
	if optVal, ok := val.(*data.OptionalValue); ok {
		if !optVal.IsSet() {
			return true
		}
	}
	return false
}

// custom Json Marshal to escape HTML values
func jsonMarshallDisableEscapeHTML(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	if err != nil {
		return []byte{}, err
	}

	// explicitly trim trailing '\n', this is required because setEscapeHtml is set to false
	s := string(buffer.Bytes())
	res := strings.TrimSuffix(s, "\n")

	return []byte(res), err
}
