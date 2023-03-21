/* Copyright © 2019-2022 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/l10n"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/lib"
	"log"
)

// ConnectionMetadataKey are keys used to store connection metadata in
// the ExecutionContext's context.
type ConnectionMetadataKey string

const RESTMetadataKey ConnectionMetadataKey = "RESTMetadata"
const ResponseTypeKey ConnectionMetadataKey = "ResponseType"

// ResponseType defines the acceptable for the VAPI consumer response types.
type ResponseType struct {
	monoResponse   bool
	streamResponse bool
}

func NewResponseType(acceptsMono bool, acceptsStream bool) *ResponseType {
	return &ResponseType{monoResponse: acceptsMono, streamResponse: acceptsStream}
}

// AcceptsMonoResponse returns true is the consumer can accept mono response.
func (resType *ResponseType) AcceptsMonoResponse() bool {
	return resType.monoResponse
}

// AcceptsStreamResponse returns true is the consumer can accept stream response.
func (resType *ResponseType) AcceptsStreamResponse() bool {
	return resType.streamResponse
}

var (
	OnlyMonoResponse     = &ResponseType{monoResponse: true, streamResponse: false}
	OnlyStreamResponse   = &ResponseType{monoResponse: false, streamResponse: true}
	MonoAndStreamRespose = &ResponseType{monoResponse: true, streamResponse: true}
)

type ExecutionContext struct {
	securityContext    SecurityContext
	applicationContext *ApplicationContext
	runtimeData        *RuntimeData
	ctx                context.Context
}

func NewExecutionContext(applicationContext *ApplicationContext, securityContext SecurityContext) *ExecutionContext {
	if applicationContext == nil {
		applicationContext = NewApplicationContext(nil)
	}

	return &ExecutionContext{
		applicationContext: applicationContext,
		securityContext:    securityContext,
		runtimeData:        NewRuntimeData(nil, nil),

		// TODO: Go context should live only on the stack
		// see here for more info: https://golang.org/pkg/context/
		// It would probably make more sense to rewrite runtime so that
		// ExecutionContext is inside go's context (as a value) not vice versa
		ctx: context.Background(),
	}
}

func (e *ExecutionContext) SetConnectionMetadata(key ConnectionMetadataKey, metadata interface{}) {
	e.ctx = context.WithValue(e.ctx, key, metadata)
}

func (e *ExecutionContext) ConnectionMetadata(key ConnectionMetadataKey) (interface{}, error) {
	metadata := e.ctx.Value(key)
	if metadata == nil {
		return nil, errors.New("cannot fetch metadata, it has not been set")
	}
	return metadata, nil
}

// Set request context
func (e *ExecutionContext) WithContext(ctx context.Context) {
	e.ctx = ctx
}

// Get request context
func (e *ExecutionContext) Context() context.Context {
	return e.ctx
}

// AcceptableResponseType returns the value of the key core.ResponseTypeKey
// Extracts the acceptable response types for the consumer VAPI.
// Defaults to core.OnlyMonoResponse.
func AcceptableResponseType(ctx context.Context) *ResponseType {
	responseType, ok := ctx.Value(ResponseTypeKey).(*ResponseType)
	if !ok {
		return OnlyMonoResponse
	}
	return responseType
}

func (e *ExecutionContext) SecurityContext() SecurityContext {
	return e.securityContext
}

func (e *ExecutionContext) ApplicationContext() *ApplicationContext {
	return e.applicationContext
}

func (e *ExecutionContext) RuntimeData() *RuntimeData {
	if e == nil {
		return nil
	}
	return e.runtimeData
}

func (e *ExecutionContext) SetSecurityContext(secContext SecurityContext) {
	e.securityContext = secContext
}

func (e *ExecutionContext) JSON() (map[string]interface{}, error) {
	var result map[string]interface{}
	tmp, err := json.Marshal(e)
	if err != nil {
		log.Fatalf("Error occurred trying to marshal execution context: %s", err)
		return nil, err
	}

	err = json.Unmarshal(tmp, &result)
	if err != nil {
		log.Fatalf("Error occurred trying to unmarshal execution context: %s", err)
		return nil, err
	}

	return result, nil
}

func (e *ExecutionContext) MarshalJSON() ([]byte, error) {
	ecJSON := struct {
		SecurityCtx SecurityContext     `json:"securityCtx,omitempty"`
		AppCtx      *ApplicationContext `json:"appCtx"`
	}{
		SecurityCtx: e.securityContext,
		AppCtx:      e.applicationContext,
	}
	return json.Marshal(ecJSON)
}

// GetMessageFormatter constructs a message formatter from the localization headers set in application
// context in the execution context of a request
func (e *ExecutionContext) GetMessageFormatter(m l10n.LocalizableMessageFactory) (l10n.MessageFormatter, error) {
	if e == nil || e.ApplicationContext() == nil {
		return *m.GetDefaultFormatter(), nil
	}
	applicationCtx := e.ApplicationContext()

	formatter, _ := m.GetFormatterForLocalizationParams(
		applicationCtx.GetProperty(lib.HTTP_ACCEPT_LANGUAGE),
		applicationCtx.GetProperty(lib.VAPI_L10N_FORMAT_LOCALE),
		applicationCtx.GetProperty(lib.VAPI_L10N_TIMEZONE))

	return formatter, nil
}

// WithResponseAcceptor returns shallow copy of ExecutionContext object and updates RuntimeData with given
// response acceptor
func (e *ExecutionContext) WithResponseAcceptor(acceptor ResponseAcceptor) *ExecutionContext {
	responseAcceptors := append(e.runtimeData.GetResponseAcceptors(), acceptor)
	runtimeData := NewRuntimeData(e.runtimeData.GetRequestProcessors(), responseAcceptors)
	return &ExecutionContext{
		applicationContext: e.applicationContext,
		securityContext:    e.securityContext,
		runtimeData:        runtimeData,
		ctx:                e.ctx,
	}
}

type TaskInvocationKey string

var taskInvocationKey TaskInvocationKey = "task invocation context"

type TaskInvocationContext struct {
	ServiceId         string
	OperationId       string
	OutputBindingType bindings.BindingType
}

func WithTaskInvocationContext(parent context.Context, taskCtx TaskInvocationContext) context.Context {
	return context.WithValue(parent, taskInvocationKey, taskCtx)
}

func GetTaskInvocationContext(ctx context.Context) (TaskInvocationContext, error) {
	val := ctx.Value(taskInvocationKey)
	if val == nil {
		return TaskInvocationContext{}, errors.New(fmt.Sprintf("context with key %s not found", taskInvocationKey))
	}

	return val.(TaskInvocationContext), nil
}
