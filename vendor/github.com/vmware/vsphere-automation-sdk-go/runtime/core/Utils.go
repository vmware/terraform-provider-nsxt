/* Copyright © 2020-2021 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package core

import (
	"context"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/l10n"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/lib"
	"strings"
)

// The SupportedByRuntimeVersion variables are referenced from generated bindings
// to ensure compatibility with go runtime version used.  The latest
// support package version is 2.
//
// Older versions are kept for compatibility. They may be removed if
// compatibility cannot be maintained.
//
// These constants should not be referenced from any other code.
const (
	SupportedByRuntimeVersion2 = true
)

// GetErrorChan return a channel with the error given as argument.
// Usually useful when an error occurs and it needs to be send over through an error channel
func GetErrorChan(inputError []error) chan error {
	errChan := make(chan error, 1)
	errChan <- bindings.VAPIerrorsToError(inputError)
	close(errChan)
	return errChan
}

// ClientAcceptsStream checks if a stream response is acceptable.
func ClientAcceptsStream(ctx context.Context) bool {
	acceptTypes := AcceptableResponseType(ctx)
	return acceptTypes.AcceptsStreamResponse()
}

// ClientAcceptsMono checks if a mono response is acceptable.
func ClientAcceptsMono(ctx context.Context) bool {
	acceptTypes := AcceptableResponseType(ctx)
	return acceptTypes.AcceptsMonoResponse()
}

// BindingAcceptOptions holds bindings information to validate against incoming client request
type BindingAcceptOptions struct {
	MethodName string
	IsStream   bool
	IsTaskOnly bool
}

// ValidateServerBindings validates whether bindings can handle initiated request
// todo: this method does not belong to core package, rather to the bindings one
// However due to cyclic imports it is currently not possible to have it there.
// We need to reorganize the code in a way that bindings import core not the vice versa.
func ValidateServerBindings(ctx *ExecutionContext, bindingAcceptOptions *BindingAcceptOptions) (bool, MethodResult) {
	if bindingAcceptOptions.IsStream && !ClientAcceptsStream(ctx.Context()) {
		err := l10n.NewRuntimeError("vapi.protocol.server.response.stream_type_not_acceptable", nil)
		errorValue := bindings.CreateErrorValueFromMessages(bindings.INVALID_REQUEST_ERROR_DEF, []error{err})
		return false, NewStreamErrorMethodResult(ctx, errorValue)
	}

	if !bindingAcceptOptions.IsStream && !ClientAcceptsMono(ctx.Context()) {
		err := l10n.NewRuntimeError("vapi.protocol.server.response.mono_type_not_acceptable", nil)
		errorValue := bindings.CreateErrorValueFromMessages(bindings.INVALID_REQUEST_ERROR_DEF, []error{err})
		return false, NewErrorResult(errorValue)
	}

	if bindingAcceptOptions.IsTaskOnly && !strings.HasSuffix(bindingAcceptOptions.MethodName, lib.TaskInvocationString) {
		err := l10n.NewRuntimeError("vapi.protocol.server.response.non_task_not_acceptable", nil)
		errorValue := bindings.CreateErrorValueFromMessages(bindings.INVALID_REQUEST_ERROR_DEF, []error{err})
		return false, NewErrorResult(errorValue)
	}

	return true, nil
}
