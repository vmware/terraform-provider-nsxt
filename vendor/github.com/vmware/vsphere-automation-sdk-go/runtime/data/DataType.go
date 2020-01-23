/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package data

import "github.com/vmware/vsphere-automation-sdk-go/runtime/l10n"

type DataType int

const (
	VOID DataType = 1 + iota
	INTEGER
	DOUBLE
	BOOLEAN
	BLOB
	STRING
	ERROR
	OPTIONAL
	LIST
	STRUCTURE
	OPAQUE
	SECRET
	STRUCTURE_REF
	DYNAMIC_STRUCTURE
	ANY_ERROR
)

func (e DataType) String() string {
	switch e {
	case VOID:
		return "VOID"
	case INTEGER:
		return "INTEGER"
	case DOUBLE:
		return "DOUBLE"
	case BOOLEAN:
		return "BOOLEAN"
	case BLOB:
		return "BLOB"
	case STRING:
		return "STRING"
	case ERROR:
		return "ERROR"
	case OPTIONAL:
		return "OPTIONAL"
	case LIST:
		return "LIST"
	case STRUCTURE:
		return "STRUCTURE"
	case OPAQUE:
		return "OPAQUE"
	case SECRET:
		return "SECRET"
	case STRUCTURE_REF:
		return "STRUCTURE_REF"
	case DYNAMIC_STRUCTURE:
		return "DYNAMIC_STRUCTURE"
	case ANY_ERROR:
		return "ANY_ERROR"
	}
	return ""
}

func (expectedType DataType) Validate(value DataValue) []error {
	actualType := "nil"
	if value != nil {
		actualType = value.Type().String()
	}
	expectedTypeStr := expectedType.String()
	if actualType == expectedTypeStr {
		return nil
	}
	var args = map[string]string{
		"actualType":   actualType,
		"expectedType": expectedType.String()}
	return []error{l10n.NewRuntimeError("vapi.data.validate.mismatch", args)}
}
