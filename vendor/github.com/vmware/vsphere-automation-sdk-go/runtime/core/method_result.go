/* Copyright © 2019, 2021 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package core

import (
	"context"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
)

type state int

const (
	initialState state = iota
	streamResult
	monoResult
)
const timeoutDuration = 20

// MonoResult stores a single piece of MethodResult data
// It consists of an output value and an error value.
type MonoResult interface {
	Output() data.DataValue
	Error() *data.ErrorValue
	IsSuccess() bool
}

// MethodResult contains data value result of an API operation.
type MethodResult interface {
	MonoResult

	ResponseStream() chan MonoResult
	IsResponseStream() bool
}

type defaultMonoResult struct {
	output data.DataValue
	error  *data.ErrorValue
}

var _ MonoResult = &defaultMonoResult{}

func NewMonoResult(output data.DataValue, error *data.ErrorValue) MonoResult {
	return &defaultMonoResult{output, error}
}

func (monoResult *defaultMonoResult) Output() data.DataValue {
	return monoResult.output
}

func (monoResult *defaultMonoResult) Error() *data.ErrorValue {
	return monoResult.error
}

func (monoResult *defaultMonoResult) IsSuccess() bool {
	return monoResult.error == (*data.ErrorValue)(nil)
}

type defaultMethodResult struct {
	initAsStream   bool
	resultState    state
	result         MonoResult
	responseStream chan MonoResult
	closer         context.CancelFunc
}

var _ MethodResult = &defaultMethodResult{}

func createMonoStream(output data.DataValue, error *data.ErrorValue) chan MonoResult {
	responseStream := make(chan MonoResult, 1)
	responseStream <- NewMonoResult(output, error)
	close(responseStream)
	return responseStream
}

//Deprecated: use NewErrorResult or NewDataResult
func NewMethodResult(output data.DataValue, error *data.ErrorValue) MethodResult {
	if output != nil && error != nil {
		panic("Can not set both output and error in MethodResult")
	}
	responseStream := createMonoStream(output, error)
	return &defaultMethodResult{initAsStream: false, resultState: initialState, responseStream: responseStream, closer: func() {}}
}

// NewDataResult creates a mono MethodResult object with an output and nil error.
func NewDataResult(output data.DataValue) MethodResult {
	responseStream := createMonoStream(output, nil)
	return &defaultMethodResult{initAsStream: false, resultState: initialState, responseStream: responseStream, closer: func() {}}
}

// NewErrorResult creates a mono MethodResult object with an error and a nil output.
func NewErrorResult(error *data.ErrorValue) MethodResult {
	responseStream := createMonoStream(nil, error)
	return &defaultMethodResult{initAsStream: false, resultState: initialState, responseStream: responseStream, closer: func() {}}
}

// NewStreamMethodResult creates a streaming MethodResult object.
// stream - a channel which will store all of the frames
// closer - a CancelFunc to close the stream channel if the object is used as a mono MethodResult.
func NewStreamMethodResult(stream chan MonoResult, closer context.CancelFunc) MethodResult {
	if stream == nil || closer == nil {
		panic("core.NewStreamMethodResult requires result channel and closer function")
	}
	return &defaultMethodResult{initAsStream: true, resultState: initialState, responseStream: stream, closer: closer}
}

// NewStreamErrorMethodResult returns a streaming MethodResult object
// with only 1 error frame.
func NewStreamErrorMethodResult(ctx *ExecutionContext, err *data.ErrorValue) MethodResult {
	cancelCtx, cancelFunc := context.WithCancel(ctx.Context())
	ctx.WithContext(cancelCtx)

	methodResultChan := make(chan MonoResult, 1)
	methodResultChan <- NewErrorResult(err)
	close(methodResultChan)
	return NewStreamMethodResult(methodResultChan, cancelFunc)
}

// Output returns the output field if methodResult is mono.
// If methodResult is a stream it returns the output field of the first frame.
// After the first use of Output, Error ot IsSuccess this instance of methodResult can no longer be used
// as a stream methodResult, it's state is set to monoResult.
func (methodResult *defaultMethodResult) Output() data.DataValue {
	methodResult.readResult()
	return methodResult.result.Output()
}

// Error returns the error field if methodResult is mono.
// If methodResult is a stream it returns the error field of the first frame.
// After the first use of Output, Error ot IsSuccess this instance of methodResult can no longer be used
// as a stream methodResult, it's state is set to monoResult.
func (methodResult *defaultMethodResult) Error() *data.ErrorValue {
	methodResult.readResult()
	return methodResult.result.Error()
}

// IsSuccess returns true if the error field for a mono methodResult is nil.
// If methodResult is a stream it returns true if the error field of the first frame is nil.
// After the first use of Output, Error ot IsSuccess this instance of methodResult can no longer be used
// as a stream methodResult, it's state is set to monoResult.
func (methodResult *defaultMethodResult) IsSuccess() bool {
	methodResult.readResult()
	return methodResult.result.Error() == (*data.ErrorValue)(nil)
}

// ResponseStream returns the stream channel of the methodResult.
// After the use of ResponseStream this instance of methodResult can't be used
// as a mono methodResult, it's state is set to streamResult.
func (methodResult *defaultMethodResult) ResponseStream() chan MonoResult {
	if methodResult.resultState == monoResult {
		panic("MethodResult is already consumed as mono result and cannot be converted to stream")
	}
	methodResult.resultState = streamResult
	return methodResult.responseStream
}

// IsResponseStream returns true if methodResult has been created with NewStreamMethodResult.
func (methodResult *defaultMethodResult) IsResponseStream() bool {
	return methodResult.initAsStream
}

func (methodResult *defaultMethodResult) readResult() {
	if methodResult.resultState == initialState {
		methodResult.result = <-methodResult.responseStream
		methodResult.resultState = monoResult
		methodResult.closer()
	} else if methodResult.resultState == streamResult {
		panic("MethodResult is already consumed as stream and cannot be converted to mono result")
	}
}
