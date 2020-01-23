/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package bindings

import (
	"errors"
	"reflect"
)

func isPointer(x interface{}) bool {
	return reflect.ValueOf(x).Kind() == reflect.Ptr
}

func VAPIerrorsToError(errs []error) error {
	message := ""
	for _, e := range errs {
		message = message + " " + e.Error()
	}
	return errors.New(message)
}
