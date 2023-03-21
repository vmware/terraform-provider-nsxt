/* Copyright © 2021 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package common

import "strings"

// StringSliceContains checks if specified string is contained in the elements of the specified slice.
// TODO replace strings.Contains with equals once a header parser is implemented
func StringSliceContains(arr []string, value string) bool {
	for _, item := range arr {
		if strings.Contains(item, value) {
			return true
		}
	}
	return false
}

// GetErrorsString constructs a single string from provided slice of errors
func GetErrorsString(errs []error) string {
	errStrings := make([]string, len(errs))
	for _, err := range errs {
		errStrings = append(errStrings, err.Error())
	}

	return strings.Join(errStrings, "; ")
}
