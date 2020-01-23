/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package server
type RequestPreProcessor interface {
	Process(requestBody *map[string]interface{}) error
}
