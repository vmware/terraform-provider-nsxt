/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package core

import "github.com/vmware/vsphere-automation-sdk-go/runtime/data"

type APIProvider interface {
	/**
	 * Invokes the specified method using the input DataValue and execution context.
	 */
	Invoke(serviceId string, operationId string, inputValue data.DataValue, ctx *ExecutionContext) MethodResult
}
