/* Copyright © 2024 Broadcom, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package customtypes

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func GetNsxIDSchema() schema.StringAttribute {
	return schema.StringAttribute{
		Optional:    true,
		Computed:    true,
		Description: "NSX ID for this resource",
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	}
}

func GetPathSchema() schema.StringAttribute {
	return schema.StringAttribute{
		Computed:    true,
		Description: "Policy path for this resource",
	}
}

func GetDisplayNameSchema() schema.StringAttribute {
	return schema.StringAttribute{
		Required:    true,
		Description: "Display name for this resource",
	}
}

func GetDescriptionSchema() schema.StringAttribute {
	return schema.StringAttribute{
		Optional:    true,
		Description: "Description for this resource",
	}
}

func GetRevisionSchema() schema.Int64Attribute {
	return schema.Int64Attribute{
		Computed:    true,
		Description: "The _revision property describes the current revision of the resource. To prevent clients from overwriting each other's changes, PUT operations must include the current _revision of the resource, which clients should obtain by issuing a GET operation. If the _revision provided in a PUT request is missing or stale, the operation will be rejected",
	}
}

func GetContextSchema() schema.SingleNestedBlock {
	return schema.SingleNestedBlock{

		Attributes: map[string]schema.Attribute{
			"project_id": schema.StringAttribute{
				Description:   "Id of the project which the resource belongs to.",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Validators:    []validator.String{stringvalidator.LengthAtLeast(1)},
			},
		},
	}
}
