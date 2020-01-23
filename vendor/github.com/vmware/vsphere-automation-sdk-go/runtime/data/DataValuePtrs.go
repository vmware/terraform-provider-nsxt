/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package data


import (
	"reflect"
)

var StructValuePtr = reflect.TypeOf((*StructValue)(nil))
var StringValuePtr = reflect.TypeOf((*StringValue)(nil))
var IntegerValuePtr = reflect.TypeOf((*IntegerValue)(nil))
var DoubleValuePtr = reflect.TypeOf((*DoubleValue)(nil))
var ListValuePtr = reflect.TypeOf((*ListValue)(nil))
var OptionalValuePtr = reflect.TypeOf((*OptionalValue)(nil))
var ErrorValuePtr = reflect.TypeOf((*ErrorValue)(nil))
var VoidValuePtr = reflect.TypeOf((*VoidValue)(nil))
var BoolValuePtr = reflect.TypeOf((*BooleanValue)(nil))
var BlobValuePtr = reflect.TypeOf(((*BlobValue)(nil)))
var SecretValuePtr = reflect.TypeOf((*SecretValue)(nil))