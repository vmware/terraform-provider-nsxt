/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package bindings


import (
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/l10n"
	"reflect"
)

type ConverterMode string

var REST ConverterMode = "REST"
var JSONRPC ConverterMode = "JSONRPC"

// TypeConverter converts between Golang Native data model and API runtime data model.
type TypeConverter struct {

	//TODO
	//remove permissive mode. typeconverter's default mode is permissive.
	/**
	 * Activates permissive mode of converter which is disabled by default.
	 * In permissive mode, Optional wrappers need not be present for deserializing DynamicStructures.
	 */
	permissive bool

	mode ConverterMode
}

func NewTypeConverter() *TypeConverter {
	return &TypeConverter{permissive: true}
}

// ConvertToGolang converts vapiValue which is an API runtime representation to its equivalent golang native representation
// with the help of bindingType
func (t *TypeConverter) ConvertToGolang(vapiValue data.DataValue, bindingType BindingType) (interface{}, []error) {
	if bindingType == nil {
		return nil, []error{l10n.NewRuntimeErrorNoParam("vapi.bindings.typeconverter.nil.type")}
	}
	if t.mode == REST {
		var visitor = NewRestToGolangVisitor(vapiValue, t)
		err := bindingType.Accept(visitor)
		if err != nil {
			return nil, err
		}
		return visitor.OutputValue(), nil
	}
	var visitor = NewVapiJsonRpcDataValueToGolangVisitor(vapiValue, t)
	err := bindingType.Accept(visitor)
	if err != nil {
		return nil, err
	}
	return visitor.OutputValue(), nil
}

// ConvertToVapi converts golangValue which is native golang value to its equivalent api runtime representation
// with the help of bindingType.
func (t *TypeConverter) ConvertToVapi(golangValue interface{}, bindingType BindingType) (data.DataValue, []error) {
	if bindingType == nil {
		return nil, []error{l10n.NewRuntimeErrorNoParam("vapi.bindings.typeconverter.nil.type")}
	}
	if golangValue != nil {
		if reflect.TypeOf(golangValue).Kind() == reflect.Ptr && reflect.ValueOf(golangValue).IsNil() {
			golangValue = nil
		} else if reflect.TypeOf(golangValue).Kind() == reflect.Slice && reflect.ValueOf(golangValue).IsNil() {
			golangValue = nil
		} else if reflect.TypeOf(golangValue).Kind() == reflect.Map && reflect.ValueOf(golangValue).IsNil() {
			golangValue = nil
		}
	}

	if t.mode == REST {
		visitor := NewGolangToRestDataValueVisitor(golangValue)
		err := bindingType.Accept(visitor)
		if err != nil {
			return nil, err
		}
		return visitor.OutputValue(), nil
	} else {
		visitor := NewGolangToVapiDataValueVisitor(golangValue)
		err := bindingType.Accept(visitor)
		if err != nil {
			return nil, err
		}
		return visitor.OutputValue(), nil
	}
}

// ConvertToDataDefinition outputs DataDefinition representation of bindingType.
func (t *TypeConverter) ConvertToDataDefinition(bindingType BindingType) (data.DataDefinition, []error) {
	if bindingType == nil {
		return nil, []error{l10n.NewRuntimeErrorNoParam("vapi.bindings.typeconverter.nil.type")}
	}
	referenceResolver := data.NewReferenceResolver()
	seenStructures := map[string]bool{}
	var visitor = NewBindingTypeToDataDefinitionVisitor(bindingType, referenceResolver, seenStructures)
	err := bindingType.Accept(visitor)
	if err != nil {
		return nil, err
	}
	err = referenceResolver.ResolveReferences()
	if err != nil {
		return nil, err
	}
	return visitor.OutputValue(), nil
}

/**
 * Sets permissive mode for TypeConverter.
 * todo: remove this method. default mode is permissive.
 */
func (t *TypeConverter) SetPermissive(permissive bool) {
	t.permissive = permissive
}

func (t *TypeConverter) SetMode(mode ConverterMode) {
	t.mode = mode
}

func (t *TypeConverter) Mode() ConverterMode {
	return t.mode
}
