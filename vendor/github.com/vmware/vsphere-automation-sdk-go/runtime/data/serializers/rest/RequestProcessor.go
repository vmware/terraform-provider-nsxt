/* Copyright © 2020 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package rest

import "net/http"

type RequestProcessor interface {
	Process(*http.Request) error
}
