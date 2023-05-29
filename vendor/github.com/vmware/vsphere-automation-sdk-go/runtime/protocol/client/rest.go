/* Copyright © 2019-2021, 2023 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package client

import (
	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data/serializers/rest"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client/internal"
	"net/http"
)

// DefaultRestClientOptions contains rest connector specific options
type DefaultRestClientOptions struct {
	// securityContextSerializerMap provides mapping between core.SecurityContext schemeId
	// and protocol.SecurityContextSerializer. protocol.SecurityContextSerializer is used to serialize
	// core.SecurityContext object into http headers
	securityContextSerializerMap map[string]protocol.SecurityContextSerializer
	// enableDefaultContentType when set to true overrides Content-Type header with value 'application/json'
	enableDefaultContentType bool
}

// NewRestClientOptions instantiates instance of RestClientOptions
func NewRestClientOptions(
	enableDefaultContentType bool,
	secCtxSerMap map[string]protocol.SecurityContextSerializer) *DefaultRestClientOptions {
	return &DefaultRestClientOptions{
		securityContextSerializerMap: secCtxSerMap,
		enableDefaultContentType:     enableDefaultContentType,
	}
}

// SecurityContextSerializers get map between security context schemeId and specific protocol.SecurityContextSerializer
func (d *DefaultRestClientOptions) SecurityContextSerializers() map[string]protocol.SecurityContextSerializer {
	return d.securityContextSerializerMap
}

// EnableDefaultContentType if true overrides Content-Type header with 'application/json' value
func (d *DefaultRestClientOptions) EnableDefaultContentType() bool {
	return d.enableDefaultContentType
}

// RestConnector extends connector type to provide REST protocol specific logic
// Deprecated: use NewConnector(url, UsingRest())
type RestConnector struct {
	*internal.RESTHttpProtocol
	connector  *connector
	statusCode int
}

// NewRestConnector instantiates instance of RestConnector
// Deprecated: use NewConnector(url, UsingRest(nil), WithHttpClient(client)) instead
func NewRestConnector(url string, client http.Client, options ...ConnectorOption) *RestConnector {
	options = append(options,
		UsingRest(nil),
		WithHttpClient(&client))

	connector := NewConnector(url, options...)
	restProtocol := connector.protocol.(*internal.RESTHttpProtocol)
	restConnector := &RestConnector{
		RESTHttpProtocol: restProtocol,
		connector:        connector,
	}

	// using ResponseAcceptor to satisfy RestConnector's (deprecated) StatusCode function
	statusCodeGetter := func(r *http.Response) {
		restConnector.statusCode = r.StatusCode
	}

	restProtocol.WithResponseAcceptors(statusCodeGetter)

	return restConnector
}

// AddRequestProcessor adds request processor to connector.
// Request processors are executed right before request is made to the server
// Deprecated: use WithRequestProcessors instead, e.g.:
// myProcessor := func(r *http.Request) error {
// 		// your request processing logic goes here
// }
// NewConnector(address, UsingRest(nil), WithRequestProcessors(myProcessor)) instead
func (r *RestConnector) AddRequestProcessor(processor rest.RequestProcessor) {
	deprecatedProcessorWrapper := func(r *http.Request) error {
		return processor.Process(r)
	}
	r.WithRequestProcessors(append(r.RequestProcessors(), deprecatedProcessorWrapper)...)
}

// StatusCode is used to get response status code
// Deprecated: use connector's ResponseAcceptors instead to read response status code value
func (r *RestConnector) StatusCode() int {
	return r.statusCode
}

// SetEnableDefaultContentType specifies whether header parameters which maps to Content-Type to be overwritten
// to value 'application/json'
// Deprecated: enabling default content type should only be set when initializing connector:
// rOptions = &RestClientOptions{
//     enableDefaultContentType: false,
// }
// NewConnector(address, UsingRest(rOptions))
func (r *RestConnector) SetEnableDefaultContentType(enableDefaultContentType bool) {
	r.Options.(*DefaultRestClientOptions).enableDefaultContentType = enableDefaultContentType
}

// SetSecCtxSerializer sets a serializer to be used for specified authentication schemeID
// Deprecated: security context serializers should only be set when initializing connector:
// rOptions = &RestClientOptions{
//     securityContextSerializerMap: securityContextSerializerMap,
// }
// NewConnector(address, UsingRest(rOptions))
func (r *RestConnector) SetSecCtxSerializer(schemeID string, serializer protocol.SecurityContextSerializer) {
	r.Options.(*DefaultRestClientOptions).securityContextSerializerMap[schemeID] = serializer
}

// SetApplicationContext specifies security context to be used by Connector instance
// Deprecated: security context should be specified when instantiating Connector instance
// Use WithApplicationContext ConnectorOption helper method
func (r *RestConnector) SetApplicationContext(context *core.ApplicationContext) {
	r.connector.appContext = context
}

// SetSecurityContext specifies application context to be used by Connector instance
// Deprecated: application context should be specified when instantiating Connector instance
// Use WithSecurityContext ConnectorOption helper method
func (r *RestConnector) SetSecurityContext(context core.SecurityContext) {
	r.connector.SetSecurityContext(context)
}
