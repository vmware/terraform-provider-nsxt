/* Copyright © 2021-2022 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package client

import (
	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data/serializers/rest"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client/internal"
	vapiHttp "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/http"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/security"
	"net/http"
)

type ConnectorOption func(*connector)

// UsingRest instruments connector to use REST protocol
// setting nil for restOptions would assume default REST client options
func UsingRest(restOptions protocol.RestClientOptions) ConnectorOption {
	return func(connector *connector) {
		if restOptions == nil {
			restOptions = NewRestClientOptions(
				true,
				map[string]protocol.SecurityContextSerializer{
					security.USER_PASSWORD_SCHEME_ID: rest.NewUserPwdSecContextSerializer(),
					security.SESSION_SCHEME_ID:       rest.NewSessionSecContextSerializer(),
					security.OAUTH_SCHEME_ID:         rest.NewOauthSecContextSerializer(),
				})
		}
		connector.protocol = internal.NewRESTHttpProtocol(connector, restOptions)
	}
}

// WithHttpClient sets specific http.Client on connector instance
func WithHttpClient(client *http.Client) ConnectorOption {
	return func(connector *connector) {
		if httpTransport, ok := connector.protocol.(internal.HttpClientTransport); ok {
			httpTransport.SetHttpClient(client)
		} else {
			panic("underlying protocol does not implement HttpClientTransport interface")
		}
	}
}

// WithStreamingProtocol sets streaming wire protocol
// Default option is lib.VAPI_STREAMING_CONTENT_TYPE
// lib.VAPI_STREAMING_CLEAN_JSON_CONTENT_TYPE can be used as well for shorter wire format.
func WithStreamingProtocol(acceptHeader string) ConnectorOption {
	return func(connector *connector) {
		if streamingTransport, ok := connector.protocol.(internal.StreamingTransport); ok {
			streamingTransport.SetStreamingProtocol(acceptHeader)
		} else {
			panic("underlying protocol does not implement StreamingTransport")
		}
	}
}

// WithClientFrameDeserializer sets streaming frame deserializer.
// Frame deserializer is used to handle raw []byte frame data
// and deserialize it into runtime specific channel of core.MethodResult type.
func WithClientFrameDeserializer(deserializer vapiHttp.ClientFrameDeserializer) ConnectorOption {
	return func(connector *connector) {
		if streamingTransport, ok := connector.protocol.(internal.StreamingTransport); ok {
			streamingTransport.SetClientFrameDeserializer(deserializer)
		} else {
			panic("underlying protocol does not implement StreamingTransport")
		}
	}
}

// WithSecurityContext sets specified security context on connector instance
func WithSecurityContext(secCtx core.SecurityContext) ConnectorOption {
	return func(connector *connector) {
		connector.securityContext = secCtx
	}
}

// WithApplicationContext sets specified application context on connector instance
func WithApplicationContext(appCtx *core.ApplicationContext) ConnectorOption {
	return func(connector *connector) {
		connector.appContext = appCtx
	}
}

// WithDecorators sets decorating APIProvider.
// APIProvider decorators wrap execution of Invoke method to extend runtime with additional functionality.
func WithDecorators(decorators ...core.APIProviderDecorator) ConnectorOption {
	return func(connector *connector) {
		connector.decorators = append(connector.decorators, decorators...)
	}
}

// WithRequestProcessors defines request processors for connector
// Used to read or extend http request object
func WithRequestProcessors(requestProcessors ...core.RequestProcessor) ConnectorOption {
	return func(connector *connector) {
		if httpTransport, ok := connector.protocol.(internal.HttpClientTransport); ok {
			httpTransport.WithRequestProcessors(requestProcessors...)
		} else {
			panic("Underlying protocol does not implement HttpClientTransport interface")
		}
	}
}

// WithResponseAcceptors defines response acceptors for connector
// Used to read specific information from http response object
func WithResponseAcceptors(responseAcceptors ...core.ResponseAcceptor) ConnectorOption {
	return func(connector *connector) {
		if httpTransport, ok := connector.protocol.(internal.HttpClientTransport); ok {
			httpTransport.WithResponseAcceptors(responseAcceptors...)
		} else {
			panic("Underlying protocol does not implement HttpClientTransport interface")
		}
	}
}

// WithAPIProvider sets custom API provider to be used for invoking of operations
// Useful if e.g. server and client bindings reside in the same host.
// You can then directly connect client to server's implementation using WithAPIProvider
//
// // Server implementation
// localProvider := NewLocalProvider()
// ...
//
// // client connector
// connector:= NewConnector("dummyAddress", WithAPIProvider(localProvider))
// // instantiate bindings client
// client := NewVMClient(connector) // from generated bindings
// client.CreateVM() // from generated bindings
func WithAPIProvider(apiProvider core.APIProvider) ConnectorOption {
	return func(connector *connector) {
		connector.provider = apiProvider
	}
}
