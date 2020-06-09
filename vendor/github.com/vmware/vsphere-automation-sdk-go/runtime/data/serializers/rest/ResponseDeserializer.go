/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package rest

import (
	"encoding/json"
	"strings"

	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data/serializers/cleanjson"
	httpStatus "github.com/vmware/vsphere-automation-sdk-go/runtime/lib/rest"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/log"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
)

var jsonToDataValueDecoder *cleanjson.JsonToDataValueDecoder = cleanjson.NewJsonToDataValueDecoder()

func deserializeResponseToDataValue(status int, headers map[string][]string, response string, restmetadata *protocol.OperationRestMetadata) (data.DataValue, error) {

	if response == "" && len(headers) == 0 {
		return data.NewVoidValue(), nil
	}

	var responseValue data.DataValue
	var err error

	if response != "" {
		d := json.NewDecoder(strings.NewReader(response))
		d.UseNumber()
		var jsondata interface{}
		if err := d.Decode(&jsondata); err != nil {
			return nil, err
		}
		if bodyName := restmetadata.ResponseBodyName(); bodyName != "" {
			responseValue = data.NewStructValue("", nil)
			responseBody, err := jsonToDataValueDecoder.Decode(jsondata)
			if err != nil {
				return nil, err
			}
			responseValue.(*data.StructValue).SetField(bodyName, responseBody)
		} else {
			responseValue, err = jsonToDataValueDecoder.Decode(jsondata)
			if err != nil {
				return nil, err
			}
		}
	} else {
		responseValue = data.NewStructValue("", nil)
	}

	// DeSerialize response headers into method result output
	// we deserialize errors only for struct values
	if _, ok := responseValue.(*data.StructValue); ok && isStatusSuccess(status) {
		fillDataValueFromHeaders(headers, responseValue.(*data.StructValue), restmetadata)
	}

	return responseValue, nil
}

// DeserializeResponse Deserializes returned server response into data-value-like method result object
func DeserializeResponse(status int, headers map[string][]string, response string, restmetadata *protocol.OperationRestMetadata) (core.MethodResult, error) {
	dataVal, err := deserializeResponseToDataValue(status, headers, response, restmetadata)
	if err != nil {
		return core.MethodResult{}, err
	}

	if isStatusSuccess(status) {
		return core.NewMethodResult(dataVal, nil), nil
	}

	//create errorval
	var errorVal *data.ErrorValue
	if restErrorName, ok := getErrorNameUsingStatus(status, restmetadata.ErrorCodeMap()); ok {
		isVapiStdError := false
		if _, ok := bindings.ERROR_MAP[restErrorName]; ok {
			isVapiStdError = true
		}
		errorVal = bindings.CreateErrorValueFromMessages(bindings.CreateStdErrorDefinition(restErrorName), []error{})
		if structVal, ok := dataVal.(*data.StructValue); !ok {
			panic("Struct value expected") // todo: localise this
		} else if isVapiStdError {
			errorVal.SetField("messages", data.NewListValue())
			var dataStruct *data.StructValue = nil
			if len(structVal.Fields()) > 0 {
				dataStruct = data.NewStructValue("", nil)
				for fieldName, fieldDataVal := range structVal.Fields() {
					dataStruct.SetField(fieldName, fieldDataVal)
				}
			}
			errorVal.SetField("data", data.NewOptionalValue(dataStruct))
		} else {
			for fieldName, fieldDataVal := range structVal.Fields() {
				errorVal.SetField(fieldName, fieldDataVal)
			}
		}
		fillDataValueFromHeaders(headers, errorVal, restmetadata)
		return core.NewMethodResult(nil, errorVal), nil
	}

	vapiStdErrorName, ok := httpStatus.HTTP_TO_VAPI_ERROR_MAP[status]
	if !ok {
		// Use "com.vmware.vapi.std.errors.internal_server_error" if response error code is not present
		vapiStdErrorName = httpStatus.HTTP_TO_VAPI_ERROR_MAP[500]
	}
	errorVal = bindings.CreateErrorValueFromMessages(bindings.ERROR_MAP[vapiStdErrorName], []error{})
	errorVal.SetField("messages", data.NewListValue())
	errorVal.SetField("data", dataVal)
	fillDataValueFromHeaders(headers, errorVal, restmetadata)
	return core.NewMethodResult(nil, errorVal), nil
}

func isStatusSuccess(status int) bool {
	for _, successfulStatus := range httpStatus.SUCCESSFUL_STATUS_CODES {
		if status == successfulStatus {
			return true
		}
	}
	return false
}

func mapKeysTolowerCase(data map[string][]string) map[string][]string {
	res := make(map[string][]string)
	for k, v := range data {
		if headerVal, ok := res[strings.ToLower(k)]; ok {
			v = append(headerVal, v...)
		}
		res[strings.ToLower(k)] = v
	}
	return res
}

func getErrorNameUsingStatus(status int, errorcodeMap map[string]int) (string, bool) {
	for errName, errStatus := range errorcodeMap {
		if errStatus == status {
			return errName, true
		}
	}
	return "", false
}

func fillDataValueFromHeaders(headers map[string][]string, dataValue data.DataValue, restmetadata *protocol.OperationRestMetadata) {
	headersLower := mapKeysTolowerCase(headers)
	var dvHeaders map[string]string
	var structValue *data.StructValue

	switch dataValue.Type() {
	case data.ERROR:
		structValue = &dataValue.(*data.ErrorValue).StructValue
		if dvHeadersFound, exist := restmetadata.ErrorHeadersNameMap()[structValue.Name()]; exist {
			dvHeaders = dvHeadersFound
		} else {
			return
		}
	case data.STRUCTURE:
		structValue = dataValue.(*data.StructValue)
		dvHeaders = restmetadata.ResultHeadersNameMap()
	}

	for dataValName, wireName := range dvHeaders {
		val, exists := headersLower[strings.ToLower(wireName)]
		if exists {
			headerDataVal, err := jsonToDataValueDecoder.Decode(strings.Join(val, ","))
			if err != nil {
				log.Error(err)
				panic(err)
			}
			structValue.SetField(dataValName, headerDataVal)
		}
	}
}
