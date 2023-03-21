/* Copyright © 2021-2022 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Package retry provides Decorator which can be setup for a retry of
// operation calls for various reasons.
// Decorator is set up in NewConnector, when instantiating a Connector instance,
// as optional connector option.
// To set it up provide connector option through WithDecorators function and
// inside call NewRetryDecorator function:
//
//		connector := client.NewConnector(
//			url,
//			client.WithDecorators(
//				retry.NewRetryDecorator(2, retryFunc)))
//
// Setting up a retry decorator requires maximum number of retry attempts (
// this does not include first attempt), and a retry boolean function.
// If the retry function returns True additional retry would be executed
// until the limits of maxRetries are reached.
// This function could be used to set up a waiting time between each retry
// request as well
// Below is example of simple use of retry decorator which
// retries on returned response status of 503 and sleeps for 10 seconds
// between each attempt:
//
//	retryFunc := func(retryContext retry.RetryContext) bool {
//		if retryContext.Response.StatusCode != 503 {
//			return false
//		}
//
//		// sleeping for some time before next retry
//		time.Sleep(10 * time.Second)
//
//		return true
//	}
//
//  // retries each request maximum two times in case of 503 response from server
//	connector := client.NewConnector(
//		url,
//		client.UsingRest(nil),
//		client.WithDecorators(
//			retry.NewRetryDecorator(
//				2,
//				retryFunc)))
//
package retry

import (
	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/log"
	"net/http"
)

// RetryFunc declares the syntax of a retry function
type RetryFunc func(retryCtx RetryContext) bool

// Decorator provides mechanism for retrying requests based on a retry
// function set up by the user of the decorator.
type Decorator struct {
	next       core.APIProvider
	maxRetries uint
	retryFunc  RetryFunc
}

// RetryContext holds information of an operation's result, response,
// attempt, and various other options which could be used in retry decorator
// 's retry function.
type RetryContext struct {
	Result      core.MethodResult
	Response    *http.Response
	Attempt     uint
	ServiceId   string
	OperationId string
	Input       data.DataValue
}

// NewRetryDecorator builds a new retry decorator based on provided maximum
// number of call attempts,
// and a retry function which determines whether and when a retry should be
// attempted.
func NewRetryDecorator(max uint, retryFunc RetryFunc) core.APIProviderDecorator {
	return func(next core.APIProvider) core.APIProvider {
		return Decorator{
			next:       next,
			maxRetries: max,
			retryFunc:  retryFunc,
		}
	}
}

// Invoke holds the logic for retrying requests.
// Inside it calls next provider decorated by retry decorator.
func (d Decorator) Invoke(serviceID string, operationID string,
	input data.DataValue, ctx *core.ExecutionContext) core.MethodResult {
	var result core.MethodResult
	var response *http.Response
	getResponse := func(resp *http.Response) {
		response = resp
	}

	extendedExecutionContext := ctx.WithResponseAcceptor(getResponse)

	// first attempt is actual invocation so counting retries till maxRetries + 1
	for attempt := uint(0); attempt < d.maxRetries+1; attempt++ {
		if attempt > 0 {
			// we only want to have retry logs after first attempt
			log.Infof("Retrying operation '%v' in service '%v'; attempt"+
				": %v'", serviceID, operationID, attempt)
		}

		result = d.next.Invoke(serviceID, operationID, input, extendedExecutionContext)

		retryContext := RetryContext{
			Result:      result,
			Response:    response,
			Attempt:     attempt,
			ServiceId:   serviceID,
			OperationId: operationID,
			Input:       input,
		}

		if d.retryFunc(retryContext) {
			continue
		}

		return result
	}

	return result
}
