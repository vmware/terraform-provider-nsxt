/* Copyright © 2020-2021 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package rest

import "net/http"

// RequestProcessor Provides access to request object right before execution
// Deprecated: use core.RequestProcessor instead
type RequestProcessor interface {
	Process(*http.Request) error
}
