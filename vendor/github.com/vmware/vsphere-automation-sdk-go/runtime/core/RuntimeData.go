/* Copyright © 2021 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package core

import "net/http"

// RequestProcessor defines contract function for accessing and modifying a
// request object
type RequestProcessor func(*http.Request) error

// ResponseAcceptor defines contract function for accessing and getting
// information from response object
type ResponseAcceptor func(*http.Response)

// RuntimeData holds custom runtime information
type RuntimeData struct {
	requestProcessors []RequestProcessor
	responseAcceptors []ResponseAcceptor
}

// NewRuntimeData creates instance of RuntimeData object
func NewRuntimeData(requestProcessors []RequestProcessor,
	responseAcceptors []ResponseAcceptor) *RuntimeData {
	return &RuntimeData{
		requestProcessors: requestProcessors,
		responseAcceptors: responseAcceptors}
}

// GetRequestProcessors returns slice of request processing functions
// executed right before making request to the server
func (r *RuntimeData) GetRequestProcessors() []RequestProcessor {
	return r.requestProcessors
}

// GetResponseAcceptors returns slice of response accepting functions
// executed right after a response from the server
func (r *RuntimeData) GetResponseAcceptors() []ResponseAcceptor {
	return r.responseAcceptors
}
