/* Copyright © 2019, 2021 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package msg

import (
	"encoding/json"
	"errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/lib"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/log"
	"reflect"
)

type JsonRpc20Response struct {
	//A String specifying the version of the JSON-RPC protocol. MUST be exactly "2.0"
	version string

	/**
	 *  This member is REQUIRED.
	 *	It MUST be the same as the value of the id member in the Request Object.
	 *	If there was an error in detecting the id in the Request object (e.g. Parse error/Invalid Request), it MUST be Null.
	 */
	id interface{}

	/**
	 *	This member is REQUIRED on success.
	 *	This member MUST NOT exist if there was an error invoking the method.
	 *	The value of this member is determined by the method invoked on the Server.
	 */
	result interface{}

	/**
	 *	This member is REQUIRED on error.
	 *	This member MUST NOT exist if there was no error triggered during invocation.
	 *	The value for this member MUST be an Object as defined in section 5.1.
	 */
	error map[string]interface{}
}

func NewJsonRpc20Response(version string, id interface{}, result interface{}, error map[string]interface{}) JsonRpc20Response {
	return JsonRpc20Response{version: version, id: id, result: result, error: error}
}

func (j JsonRpc20Response) Id() interface{} {
	return j.id
}

func (j JsonRpc20Response) Version() string {
	return j.version
}

func (j JsonRpc20Response) Result() interface{} {
	return j.result
}

func (j JsonRpc20Response) Error() map[string]interface{} {
	return j.error
}

// getJsonRPCIdValue gets id value from request or response
// id can be string, int or nil.
// it is important to find out the type of id and return it.
// when id is a number, its type is json.Number but its default representation is string.
// In this case, type has to be represented as int otherwise request and response id will be of different type.
func getJsonRPCIdValue(id interface{}) (interface{}, error) {
	var idString, isString = id.(string)
	if !isString {
		//check if id is an int.
		var jsonNumber, numError = id.(json.Number)
		if !numError {
			//id is not string or json.Number. check if its nil
			if reflect.TypeOf(id) != nil {
				log.Errorf("JSON RPC request id must be string or int or nil")
				return nil, errors.New("JSON RPC request id must be string or int or nil")
			} else {
				return nil, nil
			}
		} else {
			//check if number is int64
			var idInt, intErr = jsonNumber.Int64()
			if intErr != nil {
				log.Errorf("JSON RPC request id must be string or int or nil")
				return nil, errors.New("JSON RPC request id must be string or int or nil")
			} else {
				return idInt, nil
			}
		}

	} else {
		return idString, nil
	}
}

//DeSerializeResponse gets JsonRpc20Response object from provided string or byte array json
func DeSerializeResponse(response interface{}) (JsonRpc20Response, *JsonRpc20Error) {
	responseObj, err := DeSerializeJson(response)
	if err != nil {
		return JsonRpc20Response{}, NewJsonRpcErrorParseError(err)
	}
	return GetJsonRpc20Response(responseObj)
}

//GetJsonRpc20Response gets JsonRpc20Response object from unmarshalled json map
func GetJsonRpc20Response(response map[string]interface{}) (JsonRpc20Response, *JsonRpc20Error) {
	var version string
	if versionValue, ok := response[lib.JSONRPC]; ok {
		version = versionValue.(string)
	} else {
		return JsonRpc20Response{}, NewJsonRpcErrorInvalidParams("jsonrpc version not present")
	}
	var id interface{}
	if responseId, ok := response[lib.JSONRPC_ID]; ok {
		idVal, err := getJsonRPCIdValue(responseId)
		if err != nil {
			return JsonRpc20Response{}, NewJsonRpcErrorInvalidParams("Invalid Id")
		}
		id = idVal
	}
	if result, ok := response[lib.METHOD_RESULT]; ok {
		return NewJsonRpc20Response(version, id, result, nil), nil
	} else if err, ok := response[lib.METHOD_RESULT_ERROR]; ok {
		errMap, ok := err.(map[string]interface{})
		if !ok {
			return JsonRpc20Response{}, NewJsonRpcErrorInvalidParams("Invalid error type")
		}
		return NewJsonRpc20Response(version, id, nil, errMap), nil
	} else {
		return JsonRpc20Response{}, NewJsonRpcErrorInvalidParams("Invalid json rpc parameters")
	}
}
