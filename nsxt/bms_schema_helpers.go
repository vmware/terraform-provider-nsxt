// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// getBMSExternalIDSchema returns the external_id schema with validation for resources
func getBMSExternalIDSchema(description string) *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Description:  description,
		Required:     true,
		ForceNew:     true,
		ValidateFunc: validateBMSExternalID,
	}
}

// getBMSExternalIDDataSourceSchema returns the external_id schema with validation for data sources
func getBMSExternalIDDataSourceSchema(description string, required bool) *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Description:  description,
		Required:     required,
		Optional:     !required,
		ValidateFunc: validateBMSExternalID,
	}
}
