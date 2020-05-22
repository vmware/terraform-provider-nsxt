/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package lib

const EXECUTION_CONTEXT = "ctx"

const APPLICATION_CONTEXT = "appCtx"
const OPID = "opId"

const SECURITY_CONTEXT = "securityCtx"

// Structure name for the StructValues that represent
// map entries in the runtime
const MAP_STRUCT = "map-struct"
const MAP_ENTRY = "map-entry"
const MAP_KEY_FIELD = "key"
const MAP_VALUE_FIELD = "value"

const STRUCT_FIELDS = "fields"

// Structure name for the StructValues that represent
// operation input in the runtime
const OPERATION_INPUT = "operation-input"

// Constants for JSONRPC presentation Layer
const JSONRPC = "jsonrpc"
const JSONRPC_VERSION = "2.0"
const JSONRPC_INVOKE = "invoke"
const JSONRPC_ID = "id"
const JSONRPC_METHOD = "method"
const JSONRPC_PARAMS = "params"
const REQUEST_OPERATION_ID = "operationId"
const REQUEST_SERVICE_ID = "serviceId"
const REQUEST_INPUT = "input"
const METHOD_RESULT = "result"
const METHOD_RESULT_OUTPUT = "output"
const METHOD_RESULT_ERROR = "error"
const ERROR_CODE = "code"
const ERROR_MESSAGE = "message"
const ERROR_DATA = "data"

// Constants for REST presentation Layer
const REST_OP_ID_HEADER = "X-Request-ID"
const LOCALE = "locale"
const AUTH_HEADER = "Authorization"

//HTTP headers
const FORM_URL_CONTENT_TYPE = "application/x-www-form-urlencoded"
const JSON_CONTENT_TYPE = "application/json"
const TEXT_PLAIN_CONTENT_TYPE = "text/plain"
const HTTP_USER_AGENT_HEADER = "user-agent"
const HTTP_CONTENT_TYPE_HEADER = "Content-Type"
const HTTP_ACCEPT_LANGUAGE = "accept-language"
const HTTP_ACCEPT = "Accept"

// VAPI headers
const VAPI_SERVICE_HEADER = "vapi-service"
const VAPI_OPERATION_HEADER = "vapi-operation"
const VAPI_ERROR = "vapi-error"
const VAPI_HEADER_PREFIX = "vapi-ctx-"
const VAPI_SESSION_HEADER = "vmware-api-session-id"
const VAPI_L10N_FORMAT_LOCALE = "format-locale"
const VAPI_L10N_TIMEZONE = "timezone"
const VAPI_STREAMING_HEADER_VALUE = "application/vnd.vmware.vapi.stream.json,application/json"
const VAPI_STREAMING_CONTENT_TYPE = "application/vnd.vmware.vapi.stream.json"

const REST_METADATA = "rest-metadata"

var CRLFBytes = []byte("\r\n")
