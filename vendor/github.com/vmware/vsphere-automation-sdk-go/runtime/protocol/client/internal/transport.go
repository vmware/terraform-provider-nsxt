/* Copyright © 2019, 2021 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package internal

import (
	"encoding/asn1"
	"fmt"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/l10n"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/lib"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/log"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client/metadata"
	vapiHttp "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/http"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/security"
	"net"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"syscall"
)

// StreamingTransport is a contract for managing streaming related options
type StreamingTransport interface {
	SetStreamingProtocol(string)
	SetClientFrameDeserializer(deserializer vapiHttp.ClientFrameDeserializer)
}

// HttpClientTransport provides contract for access to http specific protocols
type HttpClientTransport interface {
	SetHttpClient(client *http.Client)
	WithRequestProcessors(...core.RequestProcessor)
	WithResponseAcceptors(...core.ResponseAcceptor)
}

// HttpTransport type provides functionality for making http calls to a server.
type HttpTransport struct {
	httpClient        *http.Client
	requestProcessors []core.RequestProcessor
	responseAcceptors []core.ResponseAcceptor
}

// NewHttpTransport instantiates instance of HttpTransport type
func NewHttpTransport() *HttpTransport {
	return &HttpTransport{
		httpClient: &http.Client{},
	}
}

// SetHttpClient sets different from default http.Client to be used for client to server communication.
func (h *HttpTransport) SetHttpClient(client *http.Client) {
	h.httpClient = client
}

var _ HttpClientTransport = &HttpTransport{}

// WithRequestProcessors sets request processors to be called before making a request to the server.
func (h *HttpTransport) WithRequestProcessors(processors ...core.RequestProcessor) {
	h.requestProcessors = processors
}

// RequestProcessors gets list of all request processors defined in connector
func (h *HttpTransport) RequestProcessors() []core.RequestProcessor {
	return h.requestProcessors
}

// WithResponseAcceptors specifies response acceptors to be called after a response is returned from the server.
func (h *HttpTransport) WithResponseAcceptors(acceptors ...core.ResponseAcceptor) {
	h.responseAcceptors = acceptors
}

// ResponseAcceptors gets list of all response acceptors defined in connector
func (h *HttpTransport) ResponseAcceptors() []core.ResponseAcceptor {
	return h.responseAcceptors
}

// executeRequest makes http call to the server
func (h *HttpTransport) executeRequest(ctx *core.ExecutionContext, request *http.Request) (*http.Response, error) {
	if ctx == nil {
		panic("execution context must not be nil")
	}

	// Allow client to access the req object before executing it
	allRequestProcessors := append(ctx.RuntimeData().GetRequestProcessors(), h.requestProcessors...)
	for _, preProcessor := range allRequestProcessors {
		err := preProcessor(request)
		if err != nil {
			return nil, err
		}
	}

	if ctx.Context() != nil {
		request = request.WithContext(ctx.Context())
	}

	response, err := h.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	acceptableResponseType := core.AcceptableResponseType(ctx.Context())

	contentType := response.Header.Get(lib.HTTP_CONTENT_TYPE_HEADER)

	if acceptableResponseType == core.OnlyStreamResponse && contentType != lib.VAPI_STREAMING_CONTENT_TYPE && contentType != lib.VAPI_STREAMING_CLEAN_JSON_CONTENT_TYPE {
		return nil, core.UnacceptableContent
	}

	// Provides access to raw response
	allResponseAcceptors := append(ctx.RuntimeData().GetResponseAcceptors(), h.responseAcceptors...)
	for _, responseAcceptor := range allResponseAcceptors {
		responseAcceptor(response)
	}

	return response, nil
}

// GetRuntimeUserAgentHeader returns User-Agent header for go runtime
func GetRuntimeUserAgentHeader() string {
	return fmt.Sprintf("vAPI/%s Go/%s (%s; %s)", metadata.RuntimeVersion, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}

// copyContextsToHeaders sets request headers using execution context properties
func copyContextsToHeaders(ctx *core.ExecutionContext, header http.Header) {
	appCtx := ctx.ApplicationContext()
	secCtx := ctx.SecurityContext()

	if appCtx != nil {
		for key, value := range appCtx.GetAllProperties() {
			keyLowerCase := strings.ToLower(key)
			switch keyLowerCase {
			case lib.HTTP_USER_AGENT_HEADER:
				// Prepend application user agent to runtime user agent
				vapiUserAgent := header.Get(lib.HTTP_USER_AGENT_HEADER)
				userAgent := fmt.Sprintf("%s %s", *value, vapiUserAgent)
				header.Set(lib.HTTP_USER_AGENT_HEADER, userAgent)
			case lib.HTTP_ACCEPT_LANGUAGE:
				header.Set(lib.HTTP_ACCEPT_LANGUAGE, *value)
			default:
				header.Set(lib.VAPI_HEADER_PREFIX+keyLowerCase, *value)
			}
		}
	}

	if secCtx != nil {
		if secCtx.Property(security.AUTHENTICATION_SCHEME_ID) == security.SESSION_SCHEME_ID {
			if sessionId, ok := secCtx.Property(security.SESSION_ID).(string); ok {
				header.Set(lib.VAPI_SESSION_HEADER, sessionId)
			} else {
				log.Errorf("Invalid session ID in security context. Skipping setting request header. Expected string but was %s",
					reflect.TypeOf(secCtx.Property(security.SESSION_ID)))
			}
		}
	}
}

type localizationError = l10n.Error
type responseError struct {
	*localizationError
	ErrorDefinition data.ErrorDefinition
}

// newResponseError instantiates instance of responseError type
func newResponseError(definition data.ErrorDefinition, message string, arguments map[string]string) *responseError {
	return &responseError{
		localizationError: l10n.NewRuntimeError(message, arguments),
		ErrorDefinition:   definition,
	}
}

// LocalizationError gets localization error
func (r *responseError) LocalizationError() *l10n.Error {
	return r.localizationError
}

// Error gets string representation of responseError
func (r *responseError) Error() string {
	return r.localizationError.Error()
}

// getVAPIError converts go error to appropriate vapi error
func getVAPIError(err error) *responseError {
	if netError, ok := err.(net.Error); ok && netError.Timeout() {
		return newResponseError(
			bindings.TIMEDOUT_ERROR_DEF,
			"vapi.server.timedout",
			map[string]string{"errMsg": err.Error()})
	}

	switch t := err.(type) {
	case *net.OpError:
		return newResponseError(
			bindings.SERVICE_UNAVAILABLE_ERROR_DEF,
			"vapi.server.unavailable",
			map[string]string{"errMsg": err.Error()})
	case syscall.Errno:
		if t == syscall.ECONNREFUSED {
			return newResponseError(
				bindings.SERVICE_UNAVAILABLE_ERROR_DEF,
				"vapi.server.unavailable",
				map[string]string{"errMsg": err.Error()})
		}
	case asn1.SyntaxError:
		return newResponseError(
			bindings.SERVICE_UNAVAILABLE_ERROR_DEF,
			"vapi.security.authentication.certificate.invalid",
			map[string]string{"errMsg": err.Error()})
	case asn1.StructuralError:
		return newResponseError(
			bindings.SERVICE_UNAVAILABLE_ERROR_DEF,
			"vapi.security.authentication.certificate.invalid",
			map[string]string{"errMsg": err.Error()})
	}

	return newResponseError(
		bindings.SERVICE_UNAVAILABLE_ERROR_DEF,
		"vapi.protocol.client.request.error",
		map[string]string{"errMsg": err.Error()})
}

// getRequestError from error message creates vapi error related to building a request
func getRequestError(errMsg string) *responseError {
	responseError := newResponseError(
		bindings.INVALID_REQUEST_ERROR_DEF,
		"vapi.protocol.client.request.error",
		map[string]string{"errMsg": errMsg})
	log.Error(responseError)
	return responseError
}

// getResponseError from error message creates vapi error related to handling a response
func getResponseError(errMsg string) *responseError {
	responseError := newResponseError(
		bindings.INTERNAL_SERVER_ERROR_DEF,
		"vapi.protocol.client.response.error",
		map[string]string{"errMsg": errMsg})
	log.Error(responseError)
	return responseError
}

// getErrorMethodResult from vapi error creates core.MethodResult object
func getErrorMethodResult(responseError *responseError) core.MethodResult {
	errVal := toErrorValue(responseError)
	return core.NewErrorResult(errVal)
}

func getErrorMonoResult(responseError *responseError) core.MonoResult {
	errVal := toErrorValue(responseError)
	return core.NewMonoResult(nil, errVal)
}

func toErrorValue(responseError *responseError) *data.ErrorValue {
	log.Error(responseError)
	errVal := bindings.CreateErrorValueFromMessages(
		responseError.ErrorDefinition,
		[]error{responseError.LocalizationError()})
	return errVal
}
