/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package bindings

import (
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"reflect"
	"time"
)

// the layout numbers used in parse function have specific meanings.
// they are not random. refer here for more.
// https://golang.org/src/time/format.go   ## 2006-01-02T15:04:05.999999999Z07:00
const RFC3339Nano_DATETIME_LAYOUT = time.RFC3339Nano

const VAPI_DATETIME_LAYOUT = "2006-01-02T15:04:05.000Z"

var IntegerBindingType = reflect.TypeOf(IntegerType{})
var BooleanBindingType = reflect.TypeOf(BooleanType{})
var OptionalBindingType = reflect.TypeOf(OptionalType{})
var StringBindingType = reflect.TypeOf(StringType{})
var OpaqueBindingType = reflect.TypeOf(OpaqueType{})
var StructBindingType = reflect.TypeOf(StructType{})
var MapBindingType = reflect.TypeOf(MapType{})
var ListBindingType = reflect.TypeOf(ListType{})
var IdBindingType = reflect.TypeOf(IdType{})
var EnumBindingType = reflect.TypeOf(EnumType{})
var SetBindingType = reflect.TypeOf(SetType{})
var VoidBindingType = reflect.TypeOf(VoidType{})
var DynamicStructBindingType = reflect.TypeOf(DynamicStructType{})
var ErrorBindingType = reflect.TypeOf(ErrorType{})
var DoubleBindingType = reflect.TypeOf(DoubleType{})
var BlobBindingType = reflect.TypeOf(BlobType{})
var DateTimeBindingType = reflect.TypeOf(DateTimeType{})
var SecretBindingType = reflect.TypeOf(SecretType{})
var UriBindingType = reflect.TypeOf(UriType{})
var AnyErrorBindingType = reflect.TypeOf(AnyErrorType{})
var ReferenceBindingType = reflect.TypeOf(ReferenceType{})

var ALREADY_EXISTS_ERROR_DEF = CreateStdErrorDefinition("com.vmware.vapi.std.errors.already_exists")
var ALREADY_IN_DESIRED_STATE_ERROR_DEF = CreateStdErrorDefinition("com.vmware.vapi.std.errors.already_in_desired_state")
var CANCELED_ERROR_DEF = CreateStdErrorDefinition("com.vmware.vapi.std.errors.canceled")
var CONCURRENT_CHANGE_ERROR_DEF = CreateStdErrorDefinition("com.vmware.vapi.std.errors.concurrent_change")
var STANDARD_ERROR_DEF = CreateStdErrorDefinition("com.vmware.vapi.std.errors.error")
var FEATURE_IN_USE_ERROR_DEF = CreateStdErrorDefinition("com.vmware.vapi.std.errors.feature_in_use")
var INTERNAL_SERVER_ERROR_DEF = CreateStdErrorDefinition("com.vmware.vapi.std.errors.internal_server_error")
var INVALID_ARGUMENT_ERROR_DEF = CreateStdErrorDefinition("com.vmware.vapi.std.errors.invalid_argument")
var INVALID_ELEMENT_CONFIGURATION_ERROR_DEF = CreateStdErrorDefinition("com.vmware.vapi.std.errors.invalid_element_configuration")
var INVALID_ELEMENT_TYPE_ERROR_DEF = CreateStdErrorDefinition("com.vmware.vapi.std.errors.invalid_element_type")
var INVALID_REQUEST_ERROR_DEF = CreateStdErrorDefinition("com.vmware.vapi.std.errors.invalid_request")
var NOT_ALLOWED_IN_CURRENT_STATE_ERROR_DEF = CreateStdErrorDefinition("com.vmware.vapi.std.errors.not_allowed_in_current_state")
var NOT_FOUND_ERROR_DEF = CreateStdErrorDefinition("com.vmware.vapi.std.errors.not_found")
var OP_NOT_FOUND_ERROR_DEF = CreateStdErrorDefinition("com.vmware.vapi.std.errors.operation_not_found")
var RESOURCE_BUSY_ERROR_DEF = CreateStdErrorDefinition("com.vmware.vapi.std.errors.resource_busy")
var RESOURCE_IN_USE_ERROR_DEF = CreateStdErrorDefinition("com.vmware.vapi.std.errors.resource_in_use")
var RESOURCE_INACCESSIBLE_ERROR_DEF = CreateStdErrorDefinition("com.vmware.vapi.std.errors.resource_inaccessible")
var SERVICE_UNAVAILABLE_ERROR_DEF = CreateStdErrorDefinition("com.vmware.vapi.std.errors.service_unavailable")
var TIMEDOUT_ERROR_DEF = CreateStdErrorDefinition("com.vmware.vapi.std.errors.timed_out")
var UNABLE_TO_ALLOCATE_RESOURCE_ERROR_DEF = CreateStdErrorDefinition("com.vmware.vapi.std.errors.unable_to_allocate_resource")
var UNAUTHENTICATED_ERROR_DEF = CreateStdErrorDefinition("com.vmware.vapi.std.errors.unauthenticated")
var UNAUTHORIZED_ERROR_DEF = CreateStdErrorDefinition("com.vmware.vapi.std.errors.unauthorized")
var UNEXPECTED_INPUT_ERROR_DEF = CreateStdErrorDefinition("com.vmware.vapi.std.errors.unexpected_input")
var UNSUPPORTED_ERROR_DEF = CreateStdErrorDefinition("com.vmware.vapi.std.errors.unsupported")
var UNVERIFIED_PEER_ERROR_DEF = CreateStdErrorDefinition("com.vmware.vapi.std.errors.unverified_peer")

var ERROR_MAP = map[string]data.ErrorDefinition{
	"com.vmware.vapi.std.errors.already_exists":                ALREADY_EXISTS_ERROR_DEF,
	"com.vmware.vapi.std.errors.already_in_desired_state":      ALREADY_IN_DESIRED_STATE_ERROR_DEF,
	"com.vmware.vapi.std.errors.canceled":                      CANCELED_ERROR_DEF,
	"com.vmware.vapi.std.errors.concurrent_change":             CONCURRENT_CHANGE_ERROR_DEF,
	"com.vmware.vapi.std.errors.error":                         STANDARD_ERROR_DEF,
	"com.vmware.vapi.std.errors.feature_in_use":                FEATURE_IN_USE_ERROR_DEF,
	"com.vmware.vapi.std.errors.internal_server_error":         INTERNAL_SERVER_ERROR_DEF,
	"com.vmware.vapi.std.errors.invalid_argument":              INVALID_ARGUMENT_ERROR_DEF,
	"com.vmware.vapi.std.errors.invalid_element_configuration": INVALID_ELEMENT_CONFIGURATION_ERROR_DEF,
	"com.vmware.vapi.std.errors.invalid_element_type":          INVALID_ELEMENT_TYPE_ERROR_DEF,
	"com.vmware.vapi.std.errors.invalid_request":               INVALID_REQUEST_ERROR_DEF,
	"com.vmware.vapi.std.errors.not_allowed_in_current_state":  NOT_ALLOWED_IN_CURRENT_STATE_ERROR_DEF,
	"com.vmware.vapi.std.errors.not_found":                     NOT_FOUND_ERROR_DEF,
	"com.vmware.vapi.std.errors.operation_not_found":           OP_NOT_FOUND_ERROR_DEF,
	"com.vmware.vapi.std.errors.resource_busy":                 RESOURCE_BUSY_ERROR_DEF,
	"com.vmware.vapi.std.errors.resource_in_use":               RESOURCE_IN_USE_ERROR_DEF,
	"com.vmware.vapi.std.errors.resource_inaccessible":         RESOURCE_INACCESSIBLE_ERROR_DEF,
	"com.vmware.vapi.std.errors.service_unavailable":           SERVICE_UNAVAILABLE_ERROR_DEF,
	"com.vmware.vapi.std.errors.timed_out":                     TIMEDOUT_ERROR_DEF,
	"com.vmware.vapi.std.errors.unable_to_allocate_resource":   UNABLE_TO_ALLOCATE_RESOURCE_ERROR_DEF,
	"com.vmware.vapi.std.errors.unauthenticated":               UNAUTHENTICATED_ERROR_DEF,
	"com.vmware.vapi.std.errors.unauthorized":                  UNAUTHORIZED_ERROR_DEF,
	"com.vmware.vapi.std.errors.unexpected_input":              UNEXPECTED_INPUT_ERROR_DEF,
	"com.vmware.vapi.std.errors.unsupported":                   UNSUPPORTED_ERROR_DEF,
	"com.vmware.vapi.std.errors.unverified_peer":               UNVERIFIED_PEER_ERROR_DEF,
}

var ERROR_TYPE_MAP = map[string]string{
	"com.vmware.vapi.std.errors.already_exists":                "ALREADY_EXISTS",
	"com.vmware.vapi.std.errors.already_in_desired_state":      "ALREADY_IN_DESIRED_STATE",
	"com.vmware.vapi.std.errors.canceled":                      "CANCELED",
	"com.vmware.vapi.std.errors.concurrent_change":             "CONCURRENT_CHANGE",
	"com.vmware.vapi.std.errors.error":                         "ERROR",
	"com.vmware.vapi.std.errors.feature_in_use":                "FEATURE_IN_USE",
	"com.vmware.vapi.std.errors.internal_server_error":         "INTERNAL_SERVER_ERROR",
	"com.vmware.vapi.std.errors.invalid_argument":              "INVALID_ARGUMENT",
	"com.vmware.vapi.std.errors.invalid_element_configuration": "INVALID_ELEMENT_CONFIGURATION",
	"com.vmware.vapi.std.errors.invalid_element_type":          "INVALID_ELEMENT_TYPE",
	"com.vmware.vapi.std.errors.invalid_request":               "INVALID_REQUEST",
	"com.vmware.vapi.std.errors.not_allowed_in_current_state":  "NOT_ALLOWED_IN_CURRENT_STATE",
	"com.vmware.vapi.std.errors.not_found":                     "NOT_FOUND",
	"com.vmware.vapi.std.errors.operation_not_found":           "OPERATION_NOT_FOUND",
	"com.vmware.vapi.std.errors.resource_busy":                 "RESOURCE_BUSY",
	"com.vmware.vapi.std.errors.resource_in_use":               "RESOURCE_IN_USE",
	"com.vmware.vapi.std.errors.resource_inaccessible":         "RESOURCE_INACCESSIBLE",
	"com.vmware.vapi.std.errors.service_unavailable":           "SERVICE_UNAVAILABLE",
	"com.vmware.vapi.std.errors.timed_out":                     "TIMED_OUT",
	"com.vmware.vapi.std.errors.unable_to_allocate_resource":   "UNABLE_TO_ALLOCATE_RESOURCE",
	"com.vmware.vapi.std.errors.unauthenticated":               "UNAUTHENTICATED",
	"com.vmware.vapi.std.errors.unauthorized":                  "UNAUTHORIZED",
	"com.vmware.vapi.std.errors.unexpected_input":              "UNEXPECTED_INPUT",
	"com.vmware.vapi.std.errors.unsupported":                   "UNSUPPORTED",
	"com.vmware.vapi.std.errors.unverified_peer":               "UNVERIFIED_PEER",
}
