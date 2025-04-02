/* Copyright © 2024 Broadcom, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func HandleCreateError(resourceType string, resource string, err error, diag diag.Diagnostics) {
	diag.AddError(fmt.Sprintf("Failed to create %s %s", resourceType, resource), err.Error())
}

func HandleUpdateError(resourceType string, resource string, err error, diag diag.Diagnostics) {
	diag.AddError(fmt.Sprintf("Failed to update %s %s", resourceType, resource), err.Error())
}

func HandleReadError(resourceType string, resource string, err error, diag diag.Diagnostics) {
	msg := fmt.Sprintf("Failed to read %s %s", resourceType, resource)
	if IsNotFoundError(err) {
		log.Print(msg)
		return
	}
	diag.AddError(msg, err.Error())
}

func HandleDeleteError(resourceType string, resourceID string, err error, diag diag.Diagnostics) {
	if isNotFoundError(err) {
		diag.AddWarning("Resource not found on delete", fmt.Sprintf("[WARNING] %s %s not found on backend", resourceType, resourceID))
		return
	}
	msg := fmt.Sprintf("Failed to delete %s %s", resourceType, resourceID)
	diag.AddError(msg, err.Error())
}
