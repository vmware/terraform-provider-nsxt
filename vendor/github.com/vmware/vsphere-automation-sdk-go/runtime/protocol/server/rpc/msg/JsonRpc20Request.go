/* Copyright © 2019, 2021-2022 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package msg

import (
	"github.com/vmware/vsphere-automation-sdk-go/runtime/lib"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/log"
	"reflect"
)

type JsonRpc20Request struct {

	//A String specifying the version of the JSON-RPC protocol. MUST be exactly "2.0".
	version string
	//A String containing the name of the method to be invoked.
	method string
	//A Structured value that holds the parameter values to be used during the invocation of the method.
	params map[string]interface{}

	id interface{}

	notification bool
}

func NewJsonRpc20Request(version string, method string, params map[string]interface{}, id interface{}, notification bool) JsonRpc20Request {
	//TODO:sreeshas
	//handle errors
	return JsonRpc20Request{version: version, method: method, params: params, id: id, notification: notification}
}

func (j JsonRpc20Request) JSON() map[string]interface{} {
	var result = make(map[string]interface{})
	result[lib.JSONRPC] = j.version
	result[lib.JSONRPC_METHOD] = j.method
	result[lib.JSONRPC_PARAMS] = j.params
	if !j.notification {
		result[lib.JSONRPC_ID] = j.id
	}
	return result
}

// ValidateResponse Validate a json rpc 2.0 response. Check for version / id mismatch with request.
func (j JsonRpc20Request) ValidateResponse(response JsonRpc20Response) *JsonRpc20Error {
	if j.notification {
		log.Error("JSON RPC notification does not have response")
		return NewJsonRpcErrorInvalidParams(nil)

	}
	//Check jsonrpc version
	if j.version != response.version {
		log.Errorf("JSON RPC incompatible version: %s. Expecting %s", response.version, j.version)
		return NewJsonRpcErrorInvalidParams(nil)
	}
	var id_ interface{}
	if response.error != nil {
		var responseCode = response.error["code"]
		if responseCode == JSONRPC_PARSE_ERROR || responseCode == JSONRPC_INVALID_REQUEST {
			id_ = nil
		}
	} else {
		id_ = j.id
	}
	if id_ != response.id {
		log.Errorf("Json RPC response id mismatch: %s. Expecting %s", response.id, id_)
		return NewJsonRpcErrorInvalidParams(nil)
	}
	return nil
}

func (j JsonRpc20Request) Notification() bool {
	return j.notification
}

func (j JsonRpc20Request) Version() string {
	return j.version
}

func (j JsonRpc20Request) Id() interface{} {
	return j.id
}

func (j JsonRpc20Request) Params() map[string]interface{} {
	return j.params
}

// DeSerializeRequest gets JsonRpc20Request object from provided string or byte array json
func DeSerializeRequest(request interface{}) (JsonRpc20Request, *JsonRpc20Error) {
	requestObject, err := DeSerializeJson(request)
	if err != nil {
		return JsonRpc20Request{}, NewJsonRpcErrorParseError(err)
	}
	return getJsonRpc20Request(requestObject)
}

// getJsonRpc20Request gets JsonRpc20Request from unmarshalled json map
func getJsonRpc20Request(request map[string]interface{}) (JsonRpc20Request, *JsonRpc20Error) {
	var version = request[lib.JSONRPC].(string)
	var notification = true
	var id interface{}
	var err error
	if idValue, ok := request[lib.JSONRPC_ID]; ok {
		notification = false
		id, err = getJsonRPCIdValue(idValue)
		if err != nil {
			return JsonRpc20Request{}, NewJsonRpcErrorInvalidRequest(err.Error())
		}
	}
	var method, ok = request[lib.JSONRPC_METHOD].(string)
	if !ok {
		log.Errorf("JSON RPC request method must be string ")
		return JsonRpc20Request{}, NewJsonRpcErrorInvalidRequest(nil)
	}

	if paramValue, ok := request[lib.JSONRPC_PARAMS]; ok {
		paramMap, isObject := paramValue.(map[string]interface{})
		if !isObject && reflect.TypeOf(paramValue) != nil {
			return JsonRpc20Request{}, NewJsonRpcErrorInvalidRequest(nil)
		}
		return NewJsonRpc20Request(version, method, paramMap, id, notification), nil
	} else {
		return NewJsonRpc20Request(version, method, nil, id, notification), nil
	}
}
