/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package info


type Info interface {
	Identifier() string
	SetFieldInfo(name string, fieldinfo *FieldInfo)
}
