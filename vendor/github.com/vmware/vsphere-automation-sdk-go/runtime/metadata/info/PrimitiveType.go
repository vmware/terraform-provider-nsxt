/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package info


import "github.com/vmware/vsphere-automation-sdk-go/runtime/data"

const ID PrimitiveType = 16

type PrimitiveType data.DataType

func PrimitiveTypeFromString(str string) PrimitiveType {
	switch str {
	case "VOID":
		return PrimitiveType(data.VOID)
	case "INTEGER":
		return PrimitiveType(data.INTEGER)
	case "DOUBLE":
		return PrimitiveType(data.DOUBLE)
	case "BOOLEAN":
		return PrimitiveType(data.BOOLEAN)
	case "BLOB":
		return PrimitiveType(data.BLOB)
	case "STRING":
		return PrimitiveType(data.STRING)
	case "ERROR":
		return PrimitiveType(data.ERROR)
	case "OPTIONAL":
		return PrimitiveType(data.OPTIONAL)
	case "LIST":
		return PrimitiveType(data.LIST)
	case "STRUCTURE":
		return PrimitiveType(data.STRUCTURE)
	case "OPAQUE":
		return PrimitiveType(data.OPAQUE)
	case "SECRET":
		return PrimitiveType(data.SECRET)
	case "STRUCTURE_REF":
		return PrimitiveType(data.STRUCTURE_REF)
	case "DYNAMIC_STRUCTURE":
		return PrimitiveType(data.DYNAMIC_STRUCTURE)
	case "ANY_ERROR":
		return PrimitiveType(data.ANY_ERROR)
	case "ID":
		return ID
	}
	return 0
}
