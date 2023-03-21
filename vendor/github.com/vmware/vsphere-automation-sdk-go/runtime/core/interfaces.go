/* Copyright © 2022 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package core

type JSONRPCRequestPreProcessor interface {
	Process(requestBody *map[string]interface{}) error
}
