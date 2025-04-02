/* Copyright © 2024 Broadcom, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package customtypes

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type policyTagModel struct {
	Scope types.String `tfsdk:"scope"`
	Tag   types.String `tfsdk:"tag"`
}

func (m policyTagModel) objectType() types.ObjectType {
	return types.ObjectType{AttrTypes: m.objectAttributeTypes()}
}

func (m policyTagModel) objectAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"scope": types.StringType,
		"tag":   types.StringType,
	}
}

func GetTagsSchema(required bool, forceNew bool) schema.ListNestedBlock {
	var planModifiers []planmodifier.String
	if forceNew {
		planModifiers = []planmodifier.String{stringplanmodifier.RequiresReplace()}
	}
	return schema.ListNestedBlock{
		Description: "Description for this resource",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"scope": schema.StringAttribute{
					Optional:      true,
					PlanModifiers: planModifiers,
				},
				"tag": schema.StringAttribute{
					Optional:      true,
					PlanModifiers: planModifiers,
				},
			},
		},
	}
}

func GetPolicyTagsFromSchema(ctx context.Context, tags types.List, diag diag.Diagnostics) []model.Tag {
	elems := make([]policyTagModel, len(tags.Elements()))
	sdkTags := make([]model.Tag, len(tags.Elements()))
	diags := tags.ElementsAs(ctx, &elems, false)
	diag.Append(diags...)
	if !diag.HasError() {
		for i, tag := range elems {
			sdkTags[i] = model.Tag{
				Scope: tag.Scope.ValueStringPointer(),
				Tag:   tag.Tag.ValueStringPointer(),
			}
		}
	}
	return sdkTags
}

func SetPolicyTagsInSchema(sdkTags []model.Tag, diag diag.Diagnostics) types.List {
	values := make([]attr.Value, len(sdkTags))
	for i, tag := range sdkTags {
		t, d := types.ObjectValue(policyTagModel{}.objectAttributeTypes(),
			map[string]attr.Value{
				"scope": types.StringPointerValue(tag.Scope),
				"tag":   types.StringPointerValue(tag.Tag),
			})
		if d.HasError() {
			diag.Append(d...)
			return types.List{}
		}
		values[i] = t
	}
	l, d := types.ListValue(policyTagModel{}.objectType(), values)
	diag.Append(d...)
	return l
}
