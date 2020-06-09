/* Copyright © 2019-2020 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package bindings

import (
	"reflect"
	"time"

	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/l10n"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/log"
)

//localizable message fields
const ID_FIELD_NAME = "id"
const DEFAULT_MESSAGE_FIELD_NAME = "default_message"
const ARGS_FIELD_NAME = "args"
const PARAMS_FIELD_NAME = "params"
const LOCALIZED_FIELD_NAME = "localized"

const MESSAGES_FIELD_NAME = "messages"
const DATA_FIELD_NAME = "data"
const ERROR_TYPE_FIELD_NAME = "error_type"

func getErrorDefFields() map[string]data.DataDefinition {
	var messageListDef = data.NewListDefinition(getLocalizableMessageDefinition())
	var dataOptionalDynamicStructureDef = data.NewOptionalDefinition(data.NewDynamicStructDefinition())

	var ERROR_DEF_FIELDS = make(map[string]data.DataDefinition, 2)
	ERROR_DEF_FIELDS[MESSAGES_FIELD_NAME] = messageListDef
	ERROR_DEF_FIELDS[DATA_FIELD_NAME] = dataOptionalDynamicStructureDef
	ERROR_DEF_FIELDS[ERROR_TYPE_FIELD_NAME] = data.NewOptionalDefinition(data.NewStringDefinition())
	return ERROR_DEF_FIELDS
}
func LocalizableMessageType() BindingType {
	fields := make(map[string]BindingType)
	fieldNameMap := make(map[string]string)
	fields["id"] = NewStringType()
	fieldNameMap["id"] = "Id"

	fields["default_message"] = NewStringType()
	fieldNameMap["default_message"] = "DefaultMessage"

	fields["args"] = NewListType(NewStringType(), reflect.TypeOf(make([]string, 1)[0]))
	fieldNameMap["args"] = "Args"

	fields["params"] = NewOptionalType(NewMapType(NewStringType(), NewReferenceType(LocalizationParamType), reflect.TypeOf(make(map[string]LocalizationParam))))
	fieldNameMap["params"] = "Params"

	fields["localized"] = NewOptionalType(NewStringType())
	fieldNameMap["localized"] = "Localized"

	return NewStructType("com.vmware.vapi.std.localizable_message", fields,
		reflect.TypeOf(LocalizableMessage{}), fieldNameMap, nil)
}

func LocalizationParamType() BindingType {
	fields := make(map[string]BindingType)
	fieldNameMap := make(map[string]string)
	fields["s"] = NewOptionalType(NewStringType())
	fieldNameMap["s"] = "S"

	fields["dt"] = NewOptionalType(NewDateTimeType())
	fieldNameMap["dt"] = "Dt"

	fields["i"] = NewOptionalType(NewIntegerType())
	fieldNameMap["i"] = "I"

	fields["d"] = NewOptionalType(NewDoubleType())
	fieldNameMap["d"] = "D"

	fields["l"] = NewOptionalType(NewReferenceType(NestedLocalizableMessageType))
	fieldNameMap["l"] = "L"

	fields["format"] = NewOptionalType(NewEnumType("com.vmware.vapi.std.localization_param.date_time_format", reflect.TypeOf(LocalizationParam_DateTimeFormat(""))))
	fieldNameMap["format"] = "Format"

	fields["precision"] = NewOptionalType(NewIntegerType())
	fieldNameMap["precision"] = "Precision"

	return NewStructType("com.vmware.vapi.std.localization_param", fields,
		reflect.TypeOf(LocalizationParam{}), fieldNameMap, nil)
}

func NestedLocalizableMessageType() BindingType {
	fields := make(map[string]BindingType)
	fieldNameMap := make(map[string]string)
	fields["id"] = NewStringType()
	fieldNameMap["id"] = "Id"

	fields["params"] = NewOptionalType(NewMapType(NewStringType(), NewReferenceType(LocalizationParamType), reflect.TypeOf(make(map[string]LocalizationParam))))
	fieldNameMap["params"] = "Params"

	return NewStructType("com.vmware.vapi.std.nested_localizable_message", fields,
		reflect.TypeOf(NestedLocalizableMessage{}), fieldNameMap, nil)
}

type NestedLocalizableMessage struct {
	Id     string
	Params map[string]LocalizationParam
}

type LocalizationParam struct {
	S         *string
	Dt        *time.Time
	I         *int64
	D         *float64
	L         *NestedLocalizableMessage
	Format    *LocalizationParam_DateTimeFormat
	Precision *int64
}

type LocalizationParam_DateTimeFormat string

type LocalizableMessage struct {
	Id             string
	DefaultMessage string
	Args           []string
	Params         map[string]LocalizationParam
	Localized      *string
}

func getLocalizableMessageDefinition() data.StructDefinition {
	typeConverter := NewTypeConverter()
	locMsgDefinition, err := typeConverter.ConvertToDataDefinition(LocalizableMessageType())
	if err != nil {
		log.Error("Error creating definition for LocalizableMessage")
		return data.StructDefinition{}
	}
	return locMsgDefinition.(data.StructDefinition)
}

/**
 * Internal function to create a "standard" ErrorDefinition for use only by vAPI runtime.
 */
func CreateStdErrorDefinition(name string) data.ErrorDefinition {
	var errorDef = data.NewErrorDefinition(name, getErrorDefFields())
	return errorDef
}

func ConvertMessageToStructValue(message *l10n.Error) *data.StructValue {
	var localizableMessage = getLocalizableMessageDefinition()

	var result = localizableMessage.NewValue()
	result.SetField(ID_FIELD_NAME, data.NewStringValue(message.ID()))
	result.SetField(DEFAULT_MESSAGE_FIELD_NAME, data.NewStringValue(message.Error()))
	var listValue = data.NewListValue()
	for _, arg := range message.Args() {
		listValue.Add(data.NewStringValue(arg))
	}
	result.SetField(ARGS_FIELD_NAME, listValue)
	//TODO: JIRA 1630
	// construct localized message from localization parameters in execution context
	result.SetField(LOCALIZED_FIELD_NAME, data.NewOptionalValue(nil))
	result.SetField(PARAMS_FIELD_NAME, data.NewOptionalValue(nil))
	return result
}

func GetNestedLocalizableMessageStructValue() *data.StructValue {
	structDef := data.NewStructDefinition("name", map[string]data.DataDefinition{})
	structVal := structDef.NewValue()
	structVal.SetField("id", data.NewStringValue(""))
	structVal.SetField("params", data.NewOptionalValue(nil))
	return structVal
}

func GetLocalizationParamStructValue() *data.StructValue {
	structDef := data.NewStructDefinition("name", map[string]data.DataDefinition{})
	structVal := structDef.NewValue()
	structVal.SetField("s", data.NewOptionalValue(data.NewStringValue("")))
	structVal.SetField("dt", data.NewOptionalValue(data.NewStringValue("")))
	structVal.SetField("i", data.NewOptionalValue(data.NewIntegerValue(int64(0))))
	structVal.SetField("d", data.NewOptionalValue(data.NewDoubleValue(float64(0))))
	structVal.SetField("format", data.NewOptionalValue(data.NewStringValue("")))
	structVal.SetField("precision", data.NewOptionalValue(data.NewIntegerValue(int64(0))))
	structVal.SetField("l", GetNestedLocalizableMessageStructValue())
	return structVal
}

func CreateErrorValueFromMessagesAndData(errorDef data.ErrorDefinition, messageList []error, data *data.StructValue) *data.ErrorValue {
	errorValue := CreateErrorValueFromMessages(errorDef, messageList)
	errorValue.SetField(DATA_FIELD_NAME, data)
	return errorValue
}

func CreateErrorValueFromMessages(errorDef data.ErrorDefinition, messageList []error) *data.ErrorValue {
	var messages = data.NewListValue()
	for _, msg := range messageList {
		if l10nErr, ok := msg.(*l10n.Error); ok {
			messages.Add(ConvertMessageToStructValue(l10nErr))
		} else {
			log.Errorf("Expected msg of type *l10n.Error but found %s", reflect.TypeOf(msg))
		}

	}
	var dataOptionalDynamicStructureDef = data.NewOptionalDefinition(data.NewDynamicStructDefinition())

	var dataD = dataOptionalDynamicStructureDef.NewValue(nil)
	var errorValue = errorDef.NewValue()
	errorValue.SetField(MESSAGES_FIELD_NAME, messages)
	errorValue.SetField(DATA_FIELD_NAME, dataD)
	errorValue.SetField(ERROR_TYPE_FIELD_NAME, data.NewOptionalValue(data.NewStringValue(ERROR_TYPE_MAP[errorDef.Name()])))
	return errorValue
}

func CreateErrorValueFromMessageId(errorDef data.ErrorDefinition, msgId string, args map[string]string) *data.ErrorValue {
	var msg = l10n.NewRuntimeError(msgId, args)
	var messages = data.NewListValue()
	messages.Add(ConvertMessageToStructValue(msg))
	var dataOptionalDynamicStructureDef = data.NewOptionalDefinition(data.NewDynamicStructDefinition())

	var dataO = dataOptionalDynamicStructureDef.NewValue(nil)
	var errorValue = errorDef.NewValue()
	errorValue.SetField(MESSAGES_FIELD_NAME, messages)
	errorValue.SetField(DATA_FIELD_NAME, dataO)
	errorValue.SetField(ERROR_TYPE_FIELD_NAME, data.NewOptionalValue(data.NewStringValue(ERROR_TYPE_MAP[errorDef.Name()])))
	return errorValue
}

func CreateErrorValueFromErrorValueAndMessages(errorDef data.ErrorDefinition, cause *data.ErrorValue, messages []error) *data.ErrorValue {
	var messageList = data.NewListValue()
	for _, msg := range messages {
		if l10nErr, ok := msg.(*l10n.Error); ok {
			messageList.Add(ConvertMessageToStructValue(l10nErr))
		} else {
			log.Errorf("Expected msg of type *l10n.Error but found %s", reflect.TypeOf(msg))
		}

	}

	var causeMsgList, err = cause.Field(MESSAGES_FIELD_NAME)
	if err == nil && causeMsgList.Type() == data.LIST {
		var list = causeMsgList.(*data.ListValue).List()
		for _, causeMsg := range list {
			messageList.Add(causeMsg)
		}
	}
	var dataOptionalDynamicStructureDef = data.NewOptionalDefinition(data.NewDynamicStructDefinition())

	var errorValue = errorDef.NewValue()
	errorValue.SetField(MESSAGES_FIELD_NAME, messageList)
	errorValue.SetField(DATA_FIELD_NAME, dataOptionalDynamicStructureDef.NewValue(nil))
	errorValue.SetField(ERROR_TYPE_FIELD_NAME, data.NewOptionalValue(data.NewStringValue(ERROR_TYPE_MAP[errorDef.Name()])))
	return errorValue

}
