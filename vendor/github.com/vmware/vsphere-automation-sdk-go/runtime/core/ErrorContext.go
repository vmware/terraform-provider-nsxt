/* Copyright © 2021 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package core

import (
	"context"
)

// ErrorContext provides a context wrapper with custom error.
// It can be used for example to cancel the response stream when encountering an error
// from a lower layer.
type ErrorContext struct {
	context.Context
	cancel context.CancelFunc
	err    error
}

// ErrorFunc can be used to signal error e.g. to wrap ErrorContext.Cancel
type ErrorFunc func(err error)

// Err gives access to an error a context was canceled with.
// Err is usually retrieved on successful read from context's Done channel
func (ctx *ErrorContext) Err() error {
	if ctx.err != nil {
		return ctx.err
	}
	// if no error set, use parent's error. The parent could be a cancel context.
	return ctx.Context.Err()
}

// Cancel calls the cancel function of the embedded context and has the option to set
// specific error to propagate to context listeners
func (ctx *ErrorContext) Cancel(err error) {
	ctx.err = err
	ctx.cancel()
}

// WithErrorContext returns a new copy of ErrorContext with the parent Context
// as input parameter.
// The returned context's Done channel is closed when the returned cancel function is
// called or when the parent context's Done channel is closed, whichever happens first.
func WithErrorContext(parentCtx context.Context) *ErrorContext {
	cancelCtx, cancelFun := context.WithCancel(parentCtx)
	myCtx := &ErrorContext{Context: cancelCtx, cancel: cancelFun}
	return myCtx
}
