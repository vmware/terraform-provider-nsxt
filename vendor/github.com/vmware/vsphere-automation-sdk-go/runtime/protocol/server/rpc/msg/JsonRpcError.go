/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package msg

var JSONRPC_INVALID_REQUEST = int64(-32600)
var JSONRPC_METHOD_NOT_FOUND = int64(-32601)
var JSONRPC_INVALID_PARAMS = int64(-32602)
var JSONRPC_INTERNAL_ERROR = int64(-32603)
var JSONRPC_PARSE_ERROR = int64(-32700)
var JSONRPC_OPERATION_ID_MISMATCH = int64(-31001)

/* TRANSPORT_ERROR is defined in xmlrpc error code
 *  http://xmlrpc-epi.sourceforge.net/specs/rfc.fault_codes.php
 * but not JSON RPC 2.0
  Need this for server connection error */
var JSONRPC_TRANSPORT_ERROR = int64(-32300)
var SERVER_ERROR_RANGE_MIN = int64(-32768)
var SERVER_ERROR_RANGE_MAX = int64(-32000)

type JsonRpcRequestError struct {
	jsonRpc20Error   *JsonRpc20Error
	jsonRpc20Request *JsonRpc20Request
}

func NewJsonRpcRequestError(jsonrpc20Error *JsonRpc20Error, request *JsonRpc20Request) *JsonRpcRequestError {
	return &JsonRpcRequestError{jsonRpc20Error: jsonrpc20Error, jsonRpc20Request: request}
}

func (j *JsonRpcRequestError) Error() *JsonRpc20Error {
	return j.jsonRpc20Error
}

func (j *JsonRpcRequestError) Request() *JsonRpc20Request {
	return j.jsonRpc20Request
}

type JsonRpc20Error struct {
	//A Number that indicates the error type that occurred. This MUST be an integer.
	code int64

	/**
	 * A Primitive or Structured value that contains additional information about the error.
	 * This may be omitted.
	 */
	data interface{}

	/**
	 *  A String providing a short description of the error.
	 *	The message SHOULD be limited to a concise single sentence.
	 */
	message string
}

func getMessage(errorcode int64) string {
	switch errorcode {
	case JSONRPC_INVALID_REQUEST:
		return "Invalid Request"
	case JSONRPC_METHOD_NOT_FOUND:
		return "Method not found"
	case JSONRPC_INVALID_PARAMS:
		return "Invalid params"
	case JSONRPC_INTERNAL_ERROR:
		return "Internal Error"
	case JSONRPC_PARSE_ERROR:
		return "Parse Error"
	case JSONRPC_TRANSPORT_ERROR:
		return "Transport Error"
	case JSONRPC_OPERATION_ID_MISMATCH:
		return "Mismatching operation identifier in HTTP header and payload"
	default:
		return "Server Error"
	}
}

func NewJsonRpc20Error(code int64, data interface{}) *JsonRpc20Error {
	var message string
	if code >= SERVER_ERROR_RANGE_MIN || code <= SERVER_ERROR_RANGE_MAX {
		message = getMessage(code)
	}
	return &JsonRpc20Error{code: code, data: data, message: message}
}

func NewJsonRpcErrorParseError(data interface{}) *JsonRpc20Error {
	return NewJsonRpc20Error(JSONRPC_PARSE_ERROR, data)
}

func NewJsonRpcErrorInvalidRequest(data interface{}) *JsonRpc20Error {
	return NewJsonRpc20Error(JSONRPC_INVALID_REQUEST, data)
}

func NewJsonRpcErrorMethodNotFound(data interface{}) *JsonRpc20Error {
	return NewJsonRpc20Error(JSONRPC_METHOD_NOT_FOUND, data)
}

func NewJsonRpcErrorInvalidParams(data interface{}) *JsonRpc20Error {
	return NewJsonRpc20Error(JSONRPC_INVALID_PARAMS, data)
}

func NewJsonRpcErrorInternalError(data interface{}) *JsonRpc20Error {
	return NewJsonRpc20Error(JSONRPC_INTERNAL_ERROR, data)
}

func NewJsonRpcErrorTransportError(data interface{}) *JsonRpc20Error {
	return NewJsonRpc20Error(JSONRPC_TRANSPORT_ERROR, data)
}

func NewJsonRpcMismatchOperationIdError(data interface{}) *JsonRpc20Error {
	return NewJsonRpc20Error(JSONRPC_OPERATION_ID_MISMATCH, data)
}

func (j *JsonRpc20Error) Code() int64 {
	return j.code
}

func (j *JsonRpc20Error) Message() string {
	return j.message
}
