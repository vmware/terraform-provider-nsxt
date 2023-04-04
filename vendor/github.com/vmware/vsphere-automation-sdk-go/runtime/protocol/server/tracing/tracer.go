/* Copyright © 2021-2023 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package tracing

import (
	"context"
	"net/http"

	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
)

// Defines interfaces needed for trace integration of transports.

type WireProtocol string

const (
	Rest2018  = WireProtocol("rest-2018")
	JsonRpc   = WireProtocol("jsonrpc")
	JsonRpc11 = WireProtocol("jsonrpc1.1")
)

// StartSpan is called by server handlers to start a trace span. It reads the trace info from the request and
// initializes a server span instance
type StartSpan func(serviceId string, operationId string, protocol WireProtocol, r *http.Request) (context.Context, ServerSpan)

// StartTaskSpan is used by task infrastructure to start a trace span for a task
// (i.e. a long-running operation or "lro" for short).
//
// The given context must contain the context values returned by StartSpan, so
// that the task span can be linked as a child of the span created by StartSpan
// for the API operation which triggered the task.
type StartTaskSpan func(ctx context.Context, serviceId string, operationId string, taskId string) (context.Context, ServerSpan)

// ServerSpan describes the interface for VAPI server handlers to the tracing system span objects.
type ServerSpan interface {
	// Completes the ServerSpan
	Finish()
	// LogError writes error in the ServerSpan and cancels the related context from StartSpan
	LogError(error)
	// LogVapiError writes VAPI error to the ServerSpan.
	LogVapiError(value *data.ErrorValue)
}

// No-op implementation. Does not depend on any 3rd party library.
var NoopServerSpan = &noopServerSpan{}

// No-op implementation. Does not depend on any 3rd party library.
var NoopTracer StartSpan = func(serviceId string, operationId string, protocol WireProtocol, r *http.Request) (context.Context, ServerSpan) {
	return r.Context(), &noopServerSpan{}
}

// No-op implementation. Does not depend on any 3rd party library.
var NoopTaskTracer StartTaskSpan = func(ctx context.Context, serviceId string, operationId string, taskId string) (context.Context, ServerSpan) {
	return ctx, NoopServerSpan
}

type noopServerSpan struct {
}

func (x *noopServerSpan) Finish() {
}

func (x *noopServerSpan) LogError(error) {
}

func (x *noopServerSpan) LogVapiError(value *data.ErrorValue) {
}
