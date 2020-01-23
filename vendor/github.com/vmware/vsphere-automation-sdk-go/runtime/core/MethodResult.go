/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package core

import (
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
)

type MethodResult struct {
	output data.DataValue
	error  *data.ErrorValue
}

func NewMethodResult(output data.DataValue, error *data.ErrorValue) MethodResult {
	return MethodResult{output: output, error: error}
}

func (methodResult MethodResult) Output() data.DataValue {
	return methodResult.output
}
func (methodResult MethodResult) Error() *data.ErrorValue {
	return methodResult.error
}
func (methodResult MethodResult) IsSuccess() bool {
	return methodResult.error == (*data.ErrorValue)(nil)
}
