// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"regexp"

	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

// validateBMSExternalID validates the format of BMS external IDs
func validateBMSExternalID(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if value == "" {
		errors = append(errors, fmt.Errorf("%s cannot be empty", k))
		return
	}

	// Check for UUID format (basic validation)
	uuidRegex := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	if !uuidRegex.MatchString(value) {
		errors = append(errors, fmt.Errorf("%s must be a valid UUID format (e.g., 71be0142-2ed1-1d53-9c60-5564cf4b7e2e)", k))
		return
	}

	return warnings, errors
}

// validateBMSVersionRequirement checks NSX version compatibility
func validateBMSVersionRequirement() error {
	if !util.NsxVersionHigherOrEqual("9.0.0") {
		return fmt.Errorf("Bare Metal Server features require NSX-T version 9.0.0 or higher. Current version does not support BMS management. Please upgrade NSX-T to version 9.0.0 or later")
	}
	return nil
}
