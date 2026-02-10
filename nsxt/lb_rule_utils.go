// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func getLbRuleInverseSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Whether to reverse match result of this condition",
		Optional:    true,
		Default:     false,
	}
}

func getLbRuleCaseSensitiveSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeBool,
		Description: "If true, case is significant in condition matching",
		Optional:    true,
		Default:     true,
	}
}

func getLbRuleMatchTypeSchema() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Description:  "Match type (STARTS_WITH, ENDS_WITH, EQUALS, CONTAINS, REGEX)",
		ValidateFunc: validation.StringInSlice([]string{"STARTS_WITH", "ENDS_WITH", "EQUALS", "CONTAINS", "REGEX"}, false),
		Required:     true,
	}
}
