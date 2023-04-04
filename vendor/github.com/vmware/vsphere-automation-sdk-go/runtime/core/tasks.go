/* Copyright © 2021 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package core

// TaskResultProvider provides tasks result outcome
// Used in bindings
type TaskResultProvider interface {
	GetTaskResult(taskId string) (MonoResult, error)
}
