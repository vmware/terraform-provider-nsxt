/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package core

import "github.com/vmware/vsphere-automation-sdk-go/runtime/data"

/**
 * The <code>ApiInterface</code> interface provides introspection
 * APIs for a vAPI interface; it is implemented by API providers.
 */
type ApiInterface interface {
	Identifier() InterfaceIdentifier

	Definition() InterfaceDefinition

	MethodDefinition(MethodIdentifier) *MethodDefinition

	Invoke(ctx *ExecutionContext, methodId MethodIdentifier, input data.DataValue) MethodResult
}
