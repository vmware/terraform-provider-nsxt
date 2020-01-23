/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package msg

import (
	"github.com/vmware/vsphere-automation-sdk-go/runtime/log"
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

/**
 *	 Validate a json rpc 2.0 response.
 *   Check for version / id mismatch with request
 */
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
