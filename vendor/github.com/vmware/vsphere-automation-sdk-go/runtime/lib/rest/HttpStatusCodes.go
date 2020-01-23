/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package rest


const (
	HTTP_200_OK                              = 200
	HTTP_201_CREATED                         = 201
	HTTP_202_ACCEPTED                        = 202
	HTTP_204_NO_CONTENT                      = 204
	HTTP_304_NOT_MODIFIED                    = 304
	HTTP_400_BAD_REQUEST                     = 400
	HTTP_401_UNAUTHORIZED                    = 401
	HTTP_402_PAYMENT_REQUIRED                = 402
	HTTP_403_FORBIDDEN                       = 403
	HTTP_404_NOT_FOUND                       = 404
	HTTP_405_METHOD_NOT_ALLOWED              = 405
	HTTP_406_NOT_ACCEPTABLE                  = 406
	HTTP_407_PROXY_AUTHENTICATION_REQUIRED   = 407
	HTTP_408_REQUEST_TIMEOUT                 = 408
	HTTP_409_CONFLICT                        = 409
	HTTP_410_GONE                            = 410
	HTTP_411_LENGTH_REQUIRED                 = 411
	HTTP_412_PRECONDITION_FAILED             = 412
	HTTP_413_REQUEST_ENTITY_TOO_LARGE        = 413
	HTTP_414_REQUEST_URI_TOO_LARGE           = 414
	HTTP_415_UNSUPPORTED_MEDIA_TYPE          = 415
	HTTP_416_REQUEST_RANGE_NOT_SATISFIABLE   = 416
	HTTP_417_EXPECTATION_FAILED              = 417
	HTTP_418_IM_A_TEAPOT                     = 418
	HTTP_422_UNPROCESSABLE_ENTITY            = 422
	HTTP_423_LOCKED                          = 423
	HTTP_424_FAILED_DEPENDENCY               = 424
	HTTP_426_UPGRADE_REQUIRED                = 426
	HTTP_428_PRECONDITION_REQUIRED           = 428
	HTTP_429_TOO_MANY_REQUESTS               = 429
	HTTP_431_REQUEST_HEADER_FIELDS_TOO_LARGE = 431
	HTTP_500_INTERNAL_SERVER_ERROR           = 500
	HTTP_501_NOT_IMPLEMENTED                 = 501
	HTTP_503_SERVICE_UNAVAILABLE             = 503
	HTTP_504_GATEWAY_TIMEOUT                 = 504
	HTTP_505_HTTP_VERSION_NOT_SUPPORTED      = 505
	HTTP_506_VARIANT_ALSO_NEGOTIATES         = 506
	HTTP_507_INSUFFICIENT_STORAGE            = 507
	HTTP_508_LOOP_DETECTED                   = 508
	HTTP_509_BANDWIDTH_LIMIT_EXCEEDED        = 509
	HTTP_510_NOT_EXTENDED                    = 510
	HTTP_511_NETWORK_AUTHENTICATION_REQUIRED = 511
)

var VAPI_TO_HTTP_ERROR_MAP = map[string]int{
	"com.vmware.vapi.std.errors.already_exists":                HTTP_400_BAD_REQUEST,
	"com.vmware.vapi.std.errors.already_in_desired_state":      HTTP_400_BAD_REQUEST,
	"com.vmware.vapi.std.errors.feature_in_use":                HTTP_400_BAD_REQUEST,
	"com.vmware.vapi.std.errors.internal_server_error":         HTTP_500_INTERNAL_SERVER_ERROR,
	"com.vmware.vapi.std.errors.invalid_argument":              HTTP_400_BAD_REQUEST,
	"com.vmware.vapi.std.errors.invalid_element_configuration": HTTP_400_BAD_REQUEST,
	"com.vmware.vapi.std.errors.invalid_element_type":          HTTP_400_BAD_REQUEST,
	"com.vmware.vapi.std.errors.invalid_request":               HTTP_400_BAD_REQUEST,
	"com.vmware.vapi.std.errors.not_found":                     HTTP_404_NOT_FOUND,
	"com.vmware.vapi.std.errors.not_allowed_in_current_state":  HTTP_400_BAD_REQUEST,
	"com.vmware.vapi.std.errors.operation_not_found":           HTTP_404_NOT_FOUND,
	"com.vmware.vapi.std.errors.resource_busy":                 HTTP_400_BAD_REQUEST,
	"com.vmware.vapi.std.errors.resource_in_use":               HTTP_400_BAD_REQUEST,
	"com.vmware.vapi.std.errors.resource_inaccessible":         HTTP_400_BAD_REQUEST,
	"com.vmware.vapi.std.errors.concurrent_change":             HTTP_409_CONFLICT,
	"com.vmware.vapi.std.errors.service_unavailable":           HTTP_503_SERVICE_UNAVAILABLE,
	"com.vmware.vapi.std.errors.timed_out":                     HTTP_504_GATEWAY_TIMEOUT,
	"com.vmware.vapi.std.errors.unable_to_allocate_resource":   HTTP_400_BAD_REQUEST,
	"com.vmware.vapi.std.errors.unauthenticated":               HTTP_401_UNAUTHORIZED,
	"com.vmware.vapi.std.errors.unauthorized":                  HTTP_403_FORBIDDEN,
	"com.vmware.vapi.std.errors.unsupported":                   HTTP_400_BAD_REQUEST,
}

var HTTP_TO_VAPI_ERROR_MAP = map[int]string{
	HTTP_400_BAD_REQUEST:                     "com.vmware.vapi.std.errors.invalid_request",
	HTTP_401_UNAUTHORIZED:                    "com.vmware.vapi.std.errors.unauthenticated",
	HTTP_402_PAYMENT_REQUIRED:                "com.vmware.vapi.std.errors.unauthorized",
	HTTP_403_FORBIDDEN:                       "com.vmware.vapi.std.errors.unauthorized",
	HTTP_404_NOT_FOUND:                       "com.vmware.vapi.std.errors.not_found",
	HTTP_405_METHOD_NOT_ALLOWED:              "com.vmware.vapi.std.errors.invalid_request",
	HTTP_406_NOT_ACCEPTABLE:                  "com.vmware.vapi.std.errors.invalid_request",
	HTTP_407_PROXY_AUTHENTICATION_REQUIRED:   "com.vmware.vapi.std.errors.unauthenticated",
	HTTP_408_REQUEST_TIMEOUT:                 "com.vmware.vapi.std.errors.timed_out",
	HTTP_409_CONFLICT:                        "com.vmware.vapi.std.errors.concurrent_change",
	HTTP_410_GONE:                            "com.vmware.vapi.std.errors.not_found",
	HTTP_411_LENGTH_REQUIRED:                 "com.vmware.vapi.std.errors.invalid_request",
	HTTP_412_PRECONDITION_FAILED:             "com.vmware.vapi.std.errors.invalid_request",
	HTTP_413_REQUEST_ENTITY_TOO_LARGE:        "com.vmware.vapi.std.errors.invalid_request",
	HTTP_414_REQUEST_URI_TOO_LARGE:           "com.vmware.vapi.std.errors.invalid_request",
	HTTP_415_UNSUPPORTED_MEDIA_TYPE:          "com.vmware.vapi.std.errors.invalid_request",
	HTTP_416_REQUEST_RANGE_NOT_SATISFIABLE:   "com.vmware.vapi.std.errors.resource_inaccessible",
	HTTP_417_EXPECTATION_FAILED:              "com.vmware.vapi.std.errors.invalid_request",
	HTTP_418_IM_A_TEAPOT:                     "com.vmware.vapi.std.errors.error",
	HTTP_422_UNPROCESSABLE_ENTITY:            "com.vmware.vapi.std.errors.invalid_request",
	HTTP_423_LOCKED:                          "com.vmware.vapi.std.errors.resource_busy",
	HTTP_424_FAILED_DEPENDENCY:               "com.vmware.vapi.std.errors.invalid_request",
	HTTP_426_UPGRADE_REQUIRED:                "com.vmware.vapi.std.errors.invalid_request",
	HTTP_428_PRECONDITION_REQUIRED:           "com.vmware.vapi.std.errors.concurrent_change",
	HTTP_429_TOO_MANY_REQUESTS:               "com.vmware.vapi.std.errors.service_unavailable",
	HTTP_431_REQUEST_HEADER_FIELDS_TOO_LARGE: "com.vmware.vapi.std.errors.invalid_request",
	HTTP_500_INTERNAL_SERVER_ERROR:           "com.vmware.vapi.std.errors.internal_server_error",
	HTTP_501_NOT_IMPLEMENTED:                 "com.vmware.vapi.std.errors.error",
	HTTP_503_SERVICE_UNAVAILABLE:             "com.vmware.vapi.std.errors.service_unavailable",
	HTTP_504_GATEWAY_TIMEOUT:                 "com.vmware.vapi.std.errors.timed_out",
	HTTP_505_HTTP_VERSION_NOT_SUPPORTED:      "com.vmware.vapi.std.errors.invalid_request",
	HTTP_506_VARIANT_ALSO_NEGOTIATES:         "com.vmware.vapi.std.errors.internal_server_error",
	HTTP_507_INSUFFICIENT_STORAGE:            "com.vmware.vapi.std.errors.unable_to_allocate_resource",
	HTTP_508_LOOP_DETECTED:                   "com.vmware.vapi.std.errors.internal_server_error",
	HTTP_509_BANDWIDTH_LIMIT_EXCEEDED:        "com.vmware.vapi.std.errors.unable_to_allocate_resource",
	HTTP_510_NOT_EXTENDED:                    "com.vmware.vapi.std.errors.invalid_request",
	HTTP_511_NETWORK_AUTHENTICATION_REQUIRED: "com.vmware.vapi.std.errors.unauthenticated",
}

var SUCCESSFUL_STATUS_CODES = []int{
	HTTP_200_OK,
	HTTP_201_CREATED,
	HTTP_202_ACCEPTED,
	HTTP_204_NO_CONTENT,
	HTTP_304_NOT_MODIFIED,
}
