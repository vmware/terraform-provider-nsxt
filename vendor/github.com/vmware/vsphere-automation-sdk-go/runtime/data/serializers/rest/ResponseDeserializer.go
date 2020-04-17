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

func isStatusSuccess(status int) bool {
	for _, val := range httpStatus.SUCCESSFUL_STATUS_CODES {
		if status == val {
			return true
		}
	}
	return false
}

func DeSerializeResponseToDataValue(headers map[string][]string, response string, restmetadata *protocol.OperationRestMetadata) (data.DataValue, error) {

	if response == "" && len(headers) == 0 {
		return data.NewVoidValue(), nil
	}

	var responseValue data.DataValue
	var err error

	jsonToDataValueDecoder := cleanjson.NewJsonToDataValueDecoder()
	resultHeader := restmetadata.ResultHeadersNameMap()
	if response != "" {
		d := json.NewDecoder(strings.NewReader(response))
		d.UseNumber()
		var jsondata interface{}
		if err := d.Decode(&jsondata); err != nil {
			return nil, err
		}
		if bodyName := restmetadata.ResponseBodyName(); bodyName != "" && len(resultHeader) > 0 {
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
	headersLower := mapKeysTolowerCase(headers)
	for dataValName, wireName := range resultHeader {
		val, exists := headersLower[strings.ToLower(wireName)]
		respStruct, isStruct := responseValue.(*data.StructValue)
		if exists && isStruct {
			headerDataVal, err := jsonToDataValueDecoder.Decode(strings.Join(val, ","))
			if err != nil {
				log.Error(err)
				return nil, err
			}
			respStruct.SetField(dataValName, headerDataVal)
		}
	}

	return responseValue, nil
}

func mapKeysTolowerCase(data map[string][]string) map[string][]string {
	res := make(map[string][]string)
	for k, v := range data {
		res[strings.ToLower(k)] = v
	}
	return res
}

func DeserializeResponse(status int, headers map[string][]string, response string, restmetadata *protocol.OperationRestMetadata) (core.MethodResult, error) {
	dataVal, err := DeSerializeResponseToDataValue(headers, response, restmetadata)
	if err != nil {
		return core.MethodResult{}, err
	}
	if isStatusSuccess(status) {
		return core.NewMethodResult(dataVal, nil), nil
	}
	//create errorval
	var errorVal *data.ErrorValue
	if vapiStdErrorName, ok := getErrorNameUsingStatus(status, restmetadata.ErrorCodeMap()); ok {
		errorVal = bindings.CreateErrorValueFromMessages(bindings.CreateStdErrorDefinition(vapiStdErrorName), []error{})
		errorVal.SetField("messages", data.NewListValue())
		if structVal, ok := dataVal.(*data.StructValue); ok {
			for fieldName, fieldDataVal := range structVal.Fields() {
				errorVal.SetField(fieldName, fieldDataVal)
			}
		} else if errVal, ok := dataVal.(*data.ErrorValue); ok {
			for fieldName, fieldDataVal := range errVal.Fields() {
				errorVal.SetField(fieldName, fieldDataVal)
			}
		}
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
	return core.NewMethodResult(nil, errorVal), nil
}

func getErrorNameUsingStatus(status int, errorcodeMap map[string]int) (string, bool) {
	for errName, errStatus := range errorcodeMap {
		if errStatus == status {
			return errName, true
		}
	}
	return "", false
}
