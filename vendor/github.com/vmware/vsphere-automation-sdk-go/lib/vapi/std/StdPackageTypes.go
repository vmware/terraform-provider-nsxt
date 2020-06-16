/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Data type definitions file for package: com.vmware.vapi.std.
 * Includes binding types of a top level structures and enumerations.
 * Shared by client-side stubs and server-side skeletons to ensure type
 * compatibility.
 */

package std

import (
	"reflect"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/log"
	"time"
)


// The AuthenticationScheme class defines constants for authentication scheme identifiers for authentication mechanisms present in the vAPI infrastructure shipped by VMware. 
//
//  A third party extension can define and implements it's own authentication mechanism and define a constant in a different IDL file.
type AuthenticationScheme struct {
}
// Indicates that the request doesn't need any authentication.
const AuthenticationScheme_NO_AUTHENTICATION = "com.vmware.vapi.std.security.no_authentication"
// Indicates that the security context in a request is using a SAML bearer token based authentication scheme. 
//
//  In this scheme, the following pieces of information has to be passed in the SecurityContext structure in the execution context of the request: 
//
// * The scheme identifier: com.vmware.vapi.std.security.saml_bearer_token
// * The token itself
//
//  
//
//  Sample security context in JSON format that matches the specification: ``{
// 'schemeId': 'com.vmware.vapi.std.security.saml_bearer_token',
// 'token': 'the token itself'
// }`` vAPI runtime provide convenient factory methods that take SAML bearer token and to create the security context that conforms to the above mentioned format.
const AuthenticationScheme_SAML_BEARER_TOKEN = "com.vmware.vapi.std.security.saml_bearer_token"
// Indicates that the security context in a request is using a SAML holder-of-key token based authentication scheme. 
//
//  In this scheme, the following pieces of information has to be passed in the SecurityContext structure in the execution context of the request: 
//
// * The scheme identifier: com.vmware.vapi.std.security.saml_hok_token
// * Signature of the request: This includes - algorithm used for signing the request, SAML holder of key token and signature digest
// * Request timestamp: This includes the ``created`` and ``expires`` timestamp of the request. The timestamp should match the following format - YYYY-MM-DDThh:mm:ss.sssZ (e.g. 1878-03-03T19:20:30.451Z).
//
//  
//
//  Sample security context in JSON format that matches the specification: ``{
// 'schemeId': 'com.vmware.vapi.std.security.saml_hok_token',
// 'signature': {
// 'alg': 'RS256',
// 'samlToken': ...,
// 'value': ...,``, 'timestamp': { 'created': '2012-10-26T12:24:18.941Z', 'expires': '2012-10-26T12:44:18.941Z', } } } vAPI runtime provide convenient factory methods that take SAML holder of key token and private key to create the security context that conforms to the above mentioned format.
const AuthenticationScheme_SAML_HOK_TOKEN = "com.vmware.vapi.std.security.saml_hok_token"
// Indicates that the security context in a request is using a session identifier based authentication scheme. 
//
//  In this scheme, the following pieces of information has to be passed in the SecurityContext structure in the execution context of the request: 
//
// * The scheme identifier - com.vmware.vapi.std.security.session_id
// * Valid session identifier - This is usually returned by a login method of a session manager interface for a particular vAPI service of this authentication scheme
//
//  Sample security context in JSON format that matches the specification: ``{
// 'schemeId': 'com.vmware.vapi.std.security.session_id',
// 'sessionId': ....,
// }`` vAPI runtime provides convenient factory methods that take session identifier as input parameter and create a security context that conforms to the above format.
const AuthenticationScheme_SESSION_ID = "com.vmware.vapi.std.security.session_id"
// Indicates that the security context in a request is using username/password based authentication scheme. 
//
//  In this scheme, the following pieces of information has to be passed in the SecurityContext structure in the execution context of the request: 
//
// * The scheme identifier - com.vmware.vapi.std.security.user_pass
// * Username
// * Password
//
//  
//
//  Sample security context in JSON format that matches the specification: ``{
// 'schemeId': 'com.vmware.vapi.std.security.user_pass',
// 'userName': ....,
// 'password': ...
// }`` 
//  vAPI runtime provides convenient factory methods that take username and password as input parameters and create a security context that conforms to the above format.
const AuthenticationScheme_USER_PASSWORD = "com.vmware.vapi.std.security.user_pass"
// Indicates that the security context in a request is using OAuth2 based authentication scheme. 
//
//  In this scheme, the following pieces of information has to be passed in the SecurityContext structure in the execution context of the request: 
//
// * The scheme identifier - com.vmware.vapi.std.security.oauth
// * Valid OAuth2 access token - This is usually acquired by OAuth2 Authorization Server after successful authentication of the end user.
//
//  
//
//  Sample security context in JSON format that matches the specification: ``{
// 'schemeId': 'com.vmware.vapi.std.security.oauth',
// 'accesstoken': ....
// }`` 
//  vAPI runtime provides convenient factory methods that takes OAuth2 access token as input parameter and creates a security context that conforms to the above format.
const AuthenticationScheme_OAUTH_ACCESS_TOKEN = "com.vmware.vapi.std.security.oauth"

func (s AuthenticationScheme) GetType__() bindings.BindingType {
	return AuthenticationSchemeBindingType()
}

func (s AuthenticationScheme) GetDataValue__() (data.DataValue, []error) {
	typeConverter := bindings.NewTypeConverter()
	typeConverter.SetMode(bindings.JSONRPC)
	dataVal, err := typeConverter.ConvertToVapi(s, s.GetType__())
	if err != nil {
		log.Errorf("Error in ConvertToVapi for AuthenticationScheme._GetDataValue method - %s",
			bindings.VAPIerrorsToError(err).Error())
		return nil, err
	}
	return dataVal, nil
}


// The ``DynamicID`` class represents an identifier for a resource of an arbitrary type.
type DynamicID struct {
    // The type of resource being identified (for example ``com.acme.Person``). 
    //
    //  Interfaces that contain methods for creating and deleting resources typically contain a constant field specifying the resource type for the resources being created and deleted. The API metamodel metadata interfaces include a interface that allows retrieving all the known resource types.
	Type_ string
    // The identifier for a resource whose type is specified by DynamicID#type.
	Id string
}

func (s DynamicID) GetType__() bindings.BindingType {
	return DynamicIDBindingType()
}

func (s DynamicID) GetDataValue__() (data.DataValue, []error) {
	typeConverter := bindings.NewTypeConverter()
	typeConverter.SetMode(bindings.JSONRPC)
	dataVal, err := typeConverter.ConvertToVapi(s, s.GetType__())
	if err != nil {
		log.Errorf("Error in ConvertToVapi for DynamicID._GetDataValue method - %s",
			bindings.VAPIerrorsToError(err).Error())
		return nil, err
	}
	return dataVal, nil
}


// The ``LocalizableMessage`` class represents localizable string and message template. Interfaces include one or more localizable message templates in the exceptions they report so that clients can display diagnostic messages in the native language of the user. Interfaces can include localizable strings in the data returned from methods to allow clients to display localized status information in the native language of the user.
type LocalizableMessage struct {
    // Unique identifier of the localizable string or message template. 
    //
    //  This identifier is typically used to retrieve a locale-specific string or message template from a message catalog.
	Id string
    // The value of this localizable string or message template in the ``en_US`` (English) locale. If LocalizableMessage#id refers to a message template, the default message will contain the substituted arguments. This value can be used by clients that do not need to display strings and messages in the native language of the user. It could also be used as a fallback if a client is unable to access the appropriate message catalog.
	DefaultMessage string
    // Positional arguments to be substituted into the message template. This list will be empty if the message uses named arguments or has no arguments.
	Args []string
    // Named arguments to be substituted into the message template.
	Params map[string]LocalizationParam
    // Localized string value as per request requirements.
	Localized *string
}

func (s LocalizableMessage) GetType__() bindings.BindingType {
	return LocalizableMessageBindingType()
}

func (s LocalizableMessage) GetDataValue__() (data.DataValue, []error) {
	typeConverter := bindings.NewTypeConverter()
	typeConverter.SetMode(bindings.JSONRPC)
	dataVal, err := typeConverter.ConvertToVapi(s, s.GetType__())
	if err != nil {
		log.Errorf("Error in ConvertToVapi for LocalizableMessage._GetDataValue method - %s",
			bindings.VAPIerrorsToError(err).Error())
		return nil, err
	}
	return dataVal, nil
}


// This class holds a single message parameter and formatting settings for it. The class has fields for string, int64, float64, date time and nested messages. Only one will be used depending on the type of data sent. For date, float64 and int64 it is possible to set additional formatting details.
type LocalizationParam struct {
    // String value associated with the parameter.
	S *string
    // Date and time value associated with the parameter. Use the ``format`` property to specify date and time display style.
	Dt *time.Time
    // int64 value associated with the parameter.
	I *int64
    // The float64 value associated with the parameter. The number of displayed fractional digits is changed via ``precision`` property.
	D *float64
    // Nested localizable value associated with the parameter. This is useful construct to convert to human readable localized form enumeration class and bool values. It can also be used for proper handling of pluralization and gender forms in localization. Recursive ``NestedLocalizableMessage`` instances can be used for localizing short lists of items.
	L *NestedLocalizableMessage
    // Format associated with the date and time value in ``dt`` property. The enumeration constant ``SHORT_DATETIME`` will be used as default.
	Format *LocalizationParamDateTimeFormat
    // Number of fractional digits to include in formatted float64 value.
	Precision *int64
}

func (s LocalizationParam) GetType__() bindings.BindingType {
	return LocalizationParamBindingType()
}

func (s LocalizationParam) GetDataValue__() (data.DataValue, []error) {
	typeConverter := bindings.NewTypeConverter()
	typeConverter.SetMode(bindings.JSONRPC)
	dataVal, err := typeConverter.ConvertToVapi(s, s.GetType__())
	if err != nil {
		log.Errorf("Error in ConvertToVapi for LocalizationParam._GetDataValue method - %s",
			bindings.VAPIerrorsToError(err).Error())
		return nil, err
	}
	return dataVal, nil
}


// The ``DateTimeFormat`` enumeration class lists possible date and time formatting options. It combines the Unicode CLDR format types - full, long, medium and short with 3 different presentations - date only, time only and combined date and time presentation.
//
// <p> See {@link com.vmware.vapi.bindings.ApiEnumeration enumerated types description}.
type LocalizationParamDateTimeFormat string

const (
    // The date and time value will be formatted as short date, for example *2019-01-28*
	LocalizationParamDateTimeFormat_SHORT_DATE LocalizationParamDateTimeFormat = "SHORT_DATE"
    // The date and time value will be formatted as medium date, for example *2019 Jan 28*
	LocalizationParamDateTimeFormat_MED_DATE LocalizationParamDateTimeFormat = "MED_DATE"
    // The date and time value will be formatted as long date, for example *2019 Jan 28*
	LocalizationParamDateTimeFormat_LONG_DATE LocalizationParamDateTimeFormat = "LONG_DATE"
    // The date and time value will be formatted as full date, for example *2019 Jan 28, Mon*
	LocalizationParamDateTimeFormat_FULL_DATE LocalizationParamDateTimeFormat = "FULL_DATE"
    // The date and time value will be formatted as short time, for example *12:59*
	LocalizationParamDateTimeFormat_SHORT_TIME LocalizationParamDateTimeFormat = "SHORT_TIME"
    // The date and time value will be formatted as medium time, for example *12:59:01*
	LocalizationParamDateTimeFormat_MED_TIME LocalizationParamDateTimeFormat = "MED_TIME"
    // The date and time value will be formatted as long time, for example *12:59:01 Z*
	LocalizationParamDateTimeFormat_LONG_TIME LocalizationParamDateTimeFormat = "LONG_TIME"
    // The date and time value will be formatted as full time, for example *12:59:01 Z*
	LocalizationParamDateTimeFormat_FULL_TIME LocalizationParamDateTimeFormat = "FULL_TIME"
    // The date and time value will be formatted as short date and time, for example *2019-01-28 12:59*
	LocalizationParamDateTimeFormat_SHORT_DATE_TIME LocalizationParamDateTimeFormat = "SHORT_DATE_TIME"
    // The date and time value will be formatted as medium date and time, for example *2019 Jan 28 12:59:01*
	LocalizationParamDateTimeFormat_MED_DATE_TIME LocalizationParamDateTimeFormat = "MED_DATE_TIME"
    // The date and time value will be formatted as long date and time, for example *2019 Jan 28 12:59:01 Z*
	LocalizationParamDateTimeFormat_LONG_DATE_TIME LocalizationParamDateTimeFormat = "LONG_DATE_TIME"
    // The date and time value will be formatted as full date and time, for example *2019 Jan 28, Mon 12:59:01 Z*
	LocalizationParamDateTimeFormat_FULL_DATE_TIME LocalizationParamDateTimeFormat = "FULL_DATE_TIME"
)

func (d LocalizationParamDateTimeFormat) LocalizationParamDateTimeFormat() bool {
	switch d {
	case LocalizationParamDateTimeFormat_SHORT_DATE:
		return true
	case LocalizationParamDateTimeFormat_MED_DATE:
		return true
	case LocalizationParamDateTimeFormat_LONG_DATE:
		return true
	case LocalizationParamDateTimeFormat_FULL_DATE:
		return true
	case LocalizationParamDateTimeFormat_SHORT_TIME:
		return true
	case LocalizationParamDateTimeFormat_MED_TIME:
		return true
	case LocalizationParamDateTimeFormat_LONG_TIME:
		return true
	case LocalizationParamDateTimeFormat_FULL_TIME:
		return true
	case LocalizationParamDateTimeFormat_SHORT_DATE_TIME:
		return true
	case LocalizationParamDateTimeFormat_MED_DATE_TIME:
		return true
	case LocalizationParamDateTimeFormat_LONG_DATE_TIME:
		return true
	case LocalizationParamDateTimeFormat_FULL_DATE_TIME:
		return true
	default:
		return false
	}
}


// The ``NestedLocalizableMessage`` class represents a nested within a parameter localizable string or message template. This class is useful for modeling composite messages. Such messages are necessary to do correct pluralization of phrases, represent lists of several items etc.
type NestedLocalizableMessage struct {
    // Unique identifier of the localizable string or message template. 
    //
    //  This identifier is typically used to retrieve a locale-specific string or message template from a message catalog.
	Id string
    // Named Arguments to be substituted into the message template.
	Params map[string]LocalizationParam
}

func (s NestedLocalizableMessage) GetType__() bindings.BindingType {
	return NestedLocalizableMessageBindingType()
}

func (s NestedLocalizableMessage) GetDataValue__() (data.DataValue, []error) {
	typeConverter := bindings.NewTypeConverter()
	typeConverter.SetMode(bindings.JSONRPC)
	dataVal, err := typeConverter.ConvertToVapi(s, s.GetType__())
	if err != nil {
		log.Errorf("Error in ConvertToVapi for NestedLocalizableMessage._GetDataValue method - %s",
			bindings.VAPIerrorsToError(err).Error())
		return nil, err
	}
	return dataVal, nil
}





func AuthenticationSchemeBindingType() bindings.BindingType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	var validators = []bindings.Validator{}
	return bindings.NewStructType("com.vmware.vapi.std.authentication_scheme", fields, reflect.TypeOf(AuthenticationScheme{}), fieldNameMap, validators)
}

func DynamicIDBindingType() bindings.BindingType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["type"] = bindings.NewStringType()
	fieldNameMap["type"] = "Type_"
	fields["id"] = bindings.NewIdType(nil, "type")
	fieldNameMap["id"] = "Id"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("com.vmware.vapi.std.dynamic_ID", fields, reflect.TypeOf(DynamicID{}), fieldNameMap, validators)
}

func LocalizableMessageBindingType() bindings.BindingType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["id"] = bindings.NewStringType()
	fieldNameMap["id"] = "Id"
	fields["default_message"] = bindings.NewStringType()
	fieldNameMap["default_message"] = "DefaultMessage"
	fields["args"] = bindings.NewListType(bindings.NewStringType(), reflect.TypeOf([]string{}))
	fieldNameMap["args"] = "Args"
	fields["params"] = bindings.NewOptionalType(bindings.NewMapType(bindings.NewStringType(), bindings.NewReferenceType(LocalizationParamBindingType),reflect.TypeOf(map[string]LocalizationParam{})))
	fieldNameMap["params"] = "Params"
	fields["localized"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["localized"] = "Localized"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("com.vmware.vapi.std.localizable_message", fields, reflect.TypeOf(LocalizableMessage{}), fieldNameMap, validators)
}

func LocalizationParamBindingType() bindings.BindingType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["s"] = bindings.NewOptionalType(bindings.NewStringType())
	fieldNameMap["s"] = "S"
	fields["dt"] = bindings.NewOptionalType(bindings.NewDateTimeType())
	fieldNameMap["dt"] = "Dt"
	fields["i"] = bindings.NewOptionalType(bindings.NewIntegerType())
	fieldNameMap["i"] = "I"
	fields["d"] = bindings.NewOptionalType(bindings.NewDoubleType())
	fieldNameMap["d"] = "D"
	fields["l"] = bindings.NewOptionalType(bindings.NewReferenceType(NestedLocalizableMessageBindingType))
	fieldNameMap["l"] = "L"
	fields["format"] = bindings.NewOptionalType(bindings.NewEnumType("com.vmware.vapi.std.localization_param.date_time_format", reflect.TypeOf(LocalizationParamDateTimeFormat(LocalizationParamDateTimeFormat_SHORT_DATE))))
	fieldNameMap["format"] = "Format"
	fields["precision"] = bindings.NewOptionalType(bindings.NewIntegerType())
	fieldNameMap["precision"] = "Precision"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("com.vmware.vapi.std.localization_param", fields, reflect.TypeOf(LocalizationParam{}), fieldNameMap, validators)
}

func NestedLocalizableMessageBindingType() bindings.BindingType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["id"] = bindings.NewStringType()
	fieldNameMap["id"] = "Id"
	fields["params"] = bindings.NewOptionalType(bindings.NewMapType(bindings.NewStringType(), bindings.NewReferenceType(LocalizationParamBindingType),reflect.TypeOf(map[string]LocalizationParam{})))
	fieldNameMap["params"] = "Params"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("com.vmware.vapi.std.nested_localizable_message", fields, reflect.TypeOf(NestedLocalizableMessage{}), fieldNameMap, validators)
}


