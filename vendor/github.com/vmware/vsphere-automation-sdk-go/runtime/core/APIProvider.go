/* Copyright © 2019, 2021 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package core

import "github.com/vmware/vsphere-automation-sdk-go/runtime/data"

// APIProvider defines vAPI provider contract
type APIProvider interface {
	// Invoke invokes the specified method using the input DataValue and
	// execution context.
	Invoke(serviceId string, operationId string, inputValue data.DataValue,
		ctx *ExecutionContext) MethodResult
}

// APIProviderDecorator defines a decorator wrapping function used to chain
// various APIProvider instances
type APIProviderDecorator func(next APIProvider) APIProvider
