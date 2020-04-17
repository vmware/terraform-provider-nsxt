/* Copyright © 2019-2020 VMware, Inc. All Rights Reserved.
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
	ctx                context.Context
}

func NewExecutionContext(applicationContext *ApplicationContext, securityContext SecurityContext) *ExecutionContext {
	if applicationContext == nil {
		applicationContext = NewApplicationContext(nil)
	}

	return &ExecutionContext{applicationContext: applicationContext, securityContext: securityContext}
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

func (e *ExecutionContext) SetSecurityContext(secContext SecurityContext) {
	e.securityContext = secContext
}

// Construct a message formatter from the localization headers set in application
// context in the execution context of a request
func (ctx *ExecutionContext) GetMessageFormatter(m l10n.LocalizableMessageFactory) (l10n.MessageFormatter, error) {
	if ctx == nil || ctx.ApplicationContext() == nil {
		return *m.GetDefaultFormatter(), nil
	}
	applicationCtx := ctx.ApplicationContext()

	formatter, _ := m.GetFormatterForLocalizationParams(
		applicationCtx.GetProperty(lib.HTTP_ACCEPT_LANGUAGE),
		applicationCtx.GetProperty(lib.VAPI_L10N_FORMAT_LOCALE),
		applicationCtx.GetProperty(lib.VAPI_L10N_TIMEZONE))

	return formatter, nil
}
