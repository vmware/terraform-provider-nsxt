/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package core

type SecurityContext interface {

	Property(key string) interface{}
	GetAllProperties() map[string]interface{}
	SetProperty(key string, value interface{})

}
