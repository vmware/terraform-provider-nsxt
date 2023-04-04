/* Copyright © 2019, 2021 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package common

import (
	"github.com/google/uuid"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/lib"
)

func NewDefaultApplicationContext() *core.ApplicationContext {
	appContext := core.NewApplicationContext(nil)
	InsertOperationId(appContext)
	return appContext
}

func NewUUID() string {
	return uuid.NewString()
}

func InsertOperationId(appContext *core.ApplicationContext) {
	opId := NewUUID()
	appContext.SetProperty(lib.OPID, &opId)
}
