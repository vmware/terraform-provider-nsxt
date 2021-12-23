/* Copyright © 2019-2021 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package core

import (
	"github.com/vmware/vsphere-automation-sdk-go/runtime/l10n"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/lib"
	"golang.org/x/net/context"
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

// Set request context
func (e *ExecutionContext) WithContext(ctx context.Context) {
	e.ctx = ctx
}

// Get request context
func (e *ExecutionContext) Context() context.Context {
	return e.ctx
}

func (e *ExecutionContext) SecurityContext() SecurityContext {
	return e.securityContext
}

func (e *ExecutionContext) ApplicationContext() *ApplicationContext {
	return e.applicationContext
}

func (e *ExecutionContext) RuntimeData() *RuntimeData {
	return e.runtimeData
}

func (e *ExecutionContext) SetSecurityContext(secContext SecurityContext) {
	e.securityContext = secContext
}

// Construct a message formatter from the localization headers set in application
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
