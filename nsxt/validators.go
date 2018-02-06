/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
)

func validateIntegerInRange(min, max int) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(int)
		if value < min {
			errors = append(errors, fmt.Errorf(
				"%q cannot be lower than %d: %d", k, min, value))
		}
		if value > max {
			errors = append(errors, fmt.Errorf(
				"%q cannot be higher than %d: %d", k, max, value))
		}
		return
	}
}

func validateValueInList(v interface{}, k string, legal_values []string) (ws []string, errors []error) {
	value := v.(string)
	for _, option := range legal_values {
		if value == option {
			return
		}
	}
	errors = append(errors, fmt.Errorf("%q must be one of %s", k, legal_values))
	return
}
