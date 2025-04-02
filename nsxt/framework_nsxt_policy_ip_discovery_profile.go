/* Copyright © 2024 Broadcom, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/customtypes"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

type policyIPDiscoveryProfileModel struct {
	NsxID                       types.String `tfsdk:"nsx_id"`
	Path                        types.String `tfsdk:"path"`
	DisplayName                 types.String `tfsdk:"display_name"`
	Description                 types.String `tfsdk:"description"`
	Revision                    types.Int64  `tfsdk:"revision"`
	Tags                        types.List   `tfsdk:"tag"`
	Context                     types.Object `tfsdk:"context"`
	ArpNdBindingTimeout         types.Int64  `tfsdk:"arp_nd_binding_timeout"`
	DuplicateIpDetectionEnabled types.Bool   `tfsdk:"duplicate_ip_detection_enabled"`
	ArpBindingLimit             types.Int64  `tfsdk:"arp_binding_limit"`
	ArpSnoopingEnabled          types.Bool   `tfsdk:"arp_snooping_enabled"`
	DhcpSnoopingEnabled         types.Bool   `tfsdk:"dhcp_snooping_enabled"`
	VmtoolsEnabled              types.Bool   `tfsdk:"vmtools_enabled"`
	DhcpSnoopingV6Enabled       types.Bool   `tfsdk:"dhcp_snooping_v6_enabled"`
	NdSnoopingEnabled           types.Bool   `tfsdk:"nd_snooping_enabled"`
	NdSnoopingLimit             types.Int64  `tfsdk:"nd_snooping_limit"`
	VmtoolsV6Enabled            types.Bool   `tfsdk:"vmtools_v6_enabled"`
	TofuEnabled                 types.Bool   `tfsdk:"tofu_enabled"`
}

type PolicyIPDiscoveryProfileResource struct {
	client interface{}
}

func NewPolicyIPDiscoveryProfileResource() resource.Resource {
	return &PolicyIPDiscoveryProfileResource{}
}

func (r *PolicyIPDiscoveryProfileResource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"nsx_id":       customtypes.GetNsxIDSchema(),
			"path":         customtypes.GetPathSchema(),
			"display_name": customtypes.GetDisplayNameSchema(),
			"description":  customtypes.GetDescriptionSchema(),
			"revision":     customtypes.GetRevisionSchema(),
			"arp_nd_binding_timeout": schema.Int64Attribute{
				Description: "ARP and ND cache timeout (in minutes)",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(10),
				Validators:  []validator.Int64{int64validator.Between(5, 120)},
			},
			"duplicate_ip_detection_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Duplicate IP detection",
			},
			"arp_binding_limit": schema.Int64Attribute{
				Description: "Maximum number of ARP bindings",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(1),
				Validators:  []validator.Int64{int64validator.Between(1, 256)},
			},
			"arp_snooping_enabled": schema.BoolAttribute{
				Description: "Is ARP snooping enabled or not",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(true),
			},
			"dhcp_snooping_enabled": schema.BoolAttribute{
				Description: "Is DHCP snooping enabled or not",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(true),
			},
			"vmtools_enabled": schema.BoolAttribute{
				Description: "Is VM tools enabled or not",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(true),
			},
			"dhcp_snooping_v6_enabled": schema.BoolAttribute{
				Description: "Is DHCP snoping v6 enabled or not",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"nd_snooping_enabled": schema.BoolAttribute{
				Description: "Is ND snooping enabled or not",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"nd_snooping_limit": schema.Int64Attribute{
				Description: "Maximum number of ND (Neighbor Discovery Protocol) bindings",
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(3),
				Validators:  []validator.Int64{int64validator.Between(2, 15)},
			},
			"vmtools_v6_enabled": schema.BoolAttribute{
				Description: "Is VM tools enabled or not",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"tofu_enabled": schema.BoolAttribute{
				Description: "Is TOFU enabled or not",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
		},
		Blocks: map[string]schema.Block{
			"context": customtypes.GetContextSchema(),
			"tag":     customtypes.GetTagsSchema(false, false),
		},
	}
}

func (r *PolicyIPDiscoveryProfileResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	r.client = request.ProviderData
}

func (r *PolicyIPDiscoveryProfileResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_policy_ip_discovery_profile"
}

func ipDiscoveryProfileObjFromFrameworkSchema(ctx context.Context, plan policyIPDiscoveryProfileModel, diag diag.Diagnostics) model.IPDiscoveryProfile {
	tags := customtypes.GetPolicyTagsFromSchema(ctx, plan.Tags, diag)

	return model.IPDiscoveryProfile{
		DisplayName:         plan.DisplayName.ValueStringPointer(),
		Description:         plan.Description.ValueStringPointer(),
		Tags:                tags,
		ArpNdBindingTimeout: plan.ArpNdBindingTimeout.ValueInt64Pointer(),
		DuplicateIpDetection: &model.DuplicateIPDetectionOptions{
			DuplicateIpDetectionEnabled: plan.DuplicateIpDetectionEnabled.ValueBoolPointer(),
		},
		IpV4DiscoveryOptions: &model.IPv4DiscoveryOptions{
			ArpSnoopingConfig: &model.ArpSnoopingConfig{
				ArpBindingLimit:    plan.ArpBindingLimit.ValueInt64Pointer(),
				ArpSnoopingEnabled: plan.ArpSnoopingEnabled.ValueBoolPointer(),
			},
			DhcpSnoopingEnabled: plan.DhcpSnoopingV6Enabled.ValueBoolPointer(),
			VmtoolsEnabled:      plan.VmtoolsEnabled.ValueBoolPointer(),
		},
		IpV6DiscoveryOptions: &model.IPv6DiscoveryOptions{
			DhcpSnoopingV6Enabled: plan.DhcpSnoopingV6Enabled.ValueBoolPointer(),
			NdSnoopingConfig: &model.NdSnoopingConfig{
				NdSnoopingEnabled: plan.NdSnoopingEnabled.ValueBoolPointer(),
				NdSnoopingLimit:   plan.NdSnoopingLimit.ValueInt64Pointer(),
			},
			VmtoolsV6Enabled: plan.VmtoolsV6Enabled.ValueBoolPointer(),
		},
		TofuEnabled: plan.TofuEnabled.ValueBoolPointer(),
	}
}

func (r *PolicyIPDiscoveryProfileResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan policyIPDiscoveryProfileModel

	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	connector := GetPolicyConnector(r.client)

	// Initialize resource Id and verify this ID is not yet used
	id := GetOrGenerateID2(ctx, r.client, plan.NsxID, plan.Context, response.Diagnostics, resourceNsxtPolicyIPDiscoveryProfileExistsForFramework)
	if response.Diagnostics.HasError() {
		return
	}

	obj := ipDiscoveryProfileObjFromFrameworkSchema(ctx, plan, response.Diagnostics)
	if response.Diagnostics.HasError() {
		return
	}

	boolFalse := false
	client := infra.NewIpDiscoveryProfilesClient(GetSessionContext(ctx, r.client, plan.Context, response.Diagnostics), connector)
	err := client.Patch(id, obj, &boolFalse)
	if err != nil {
		HandleCreateError("IPDiscoveryProfile", id, err, response.Diagnostics)
		return
	}
	obj, err = client.Get(id)
	if err != nil {
		HandleCreateError("IPDiscoveryProfile", id, err, response.Diagnostics)
		return
	}

	plan.Revision = types.Int64PointerValue(obj.Revision)
	plan.Path = types.StringPointerValue(obj.Path)
	plan.NsxID = types.StringValue(id)
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	// response.Diagnostics.AddError(
	// 	"Error creating the ip discovery policy -- From the plugin framework",
	// 	"Could not create ip discovery policy, unexpected error: From the plugin framework",
	// )
	// return
}

func resourceNsxtPolicyIPDiscoveryProfileExistsForFramework(sessionContext utl.SessionContext, id string, connector client.Connector, diag diag.Diagnostics) bool {
	client := infra.NewIpDiscoveryProfilesClient(sessionContext, connector)
	_, err := client.Get(id)
	if err == nil {
		return true
	}

	if IsNotFoundError(err) {
		return false
	}

	diag.AddError("Error retrieving resource", err.Error())
	return false
}

func (r *PolicyIPDiscoveryProfileResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state policyIPDiscoveryProfileModel
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	connector := GetPolicyConnector(r.client)

	var id string
	response.Diagnostics.Append(request.State.GetAttribute(ctx, path.Root("nsx_id"), &id)...)
	if id == "" {
		id = state.NsxID.ValueString()
	}
	if id == "" {
		response.Diagnostics.AddError("Error obtaining IPDiscoveryProfile ID", "ID attribute not found")
		return
	}
	//var context policyContextModel
	//response.Diagnostics.Append(request.State.GetAttribute(ctx, path.Root("context"), &context)...)
	//log.Printf("=====> proj_id=%v", context)

	client := infra.NewIpDiscoveryProfilesClient(GetSessionContext(ctx, r.client, state.Context, response.Diagnostics), connector)
	obj, err := client.Get(id)
	if err != nil {
		HandleReadError("IPDiscoveryProfile", id, err, response.Diagnostics)
		response.State.RemoveResource(ctx)

	}

	state.DisplayName = types.StringPointerValue(obj.DisplayName)
	state.Description = types.StringPointerValue(obj.Description)
	state.Path = types.StringPointerValue(obj.Path)
	state.Revision = types.Int64PointerValue(obj.Revision)
	state.Tags = customtypes.SetPolicyTagsInSchema(obj.Tags, response.Diagnostics)
	state.ArpNdBindingTimeout = types.Int64PointerValue(obj.ArpNdBindingTimeout)
	state.DuplicateIpDetectionEnabled = types.BoolPointerValue(obj.DuplicateIpDetection.DuplicateIpDetectionEnabled)
	state.ArpBindingLimit = types.Int64PointerValue(obj.IpV4DiscoveryOptions.ArpSnoopingConfig.ArpBindingLimit)
	state.ArpSnoopingEnabled = types.BoolPointerValue(obj.IpV4DiscoveryOptions.ArpSnoopingConfig.ArpSnoopingEnabled)
	state.DhcpSnoopingEnabled = types.BoolPointerValue(obj.IpV4DiscoveryOptions.DhcpSnoopingEnabled)
	state.VmtoolsEnabled = types.BoolPointerValue(obj.IpV4DiscoveryOptions.VmtoolsEnabled)
	state.DhcpSnoopingV6Enabled = types.BoolPointerValue(obj.IpV6DiscoveryOptions.DhcpSnoopingV6Enabled)
	state.NdSnoopingEnabled = types.BoolPointerValue(obj.IpV6DiscoveryOptions.NdSnoopingConfig.NdSnoopingEnabled)
	state.NdSnoopingLimit = types.Int64PointerValue(obj.IpV6DiscoveryOptions.NdSnoopingConfig.NdSnoopingLimit)
	state.VmtoolsV6Enabled = types.BoolPointerValue(obj.IpV6DiscoveryOptions.VmtoolsV6Enabled)
	state.TofuEnabled = types.BoolPointerValue(obj.TofuEnabled)
}

func (r *PolicyIPDiscoveryProfileResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var state policyIPDiscoveryProfileModel

	response.Diagnostics.Append(request.Plan.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	//Read the rest of the configured parameters
	obj := ipDiscoveryProfileObjFromFrameworkSchema(ctx, state, response.Diagnostics)
	if response.Diagnostics.HasError() {
		return
	}
	id := state.NsxID.ValueString()
	if id == "" {
		response.Diagnostics.AddError("Error obtaining IPDiscoveryProfile ID", "Resource IPDiscoveryProfile ID field is not set")
		return
	}

	connector := GetPolicyConnector(r.client)

	client := infra.NewIpDiscoveryProfilesClient(GetSessionContext(ctx, r.client, state.Context, response.Diagnostics), connector)
	if client == nil {
		response.Diagnostics.AddError(fmt.Sprintf("Error Updating IPDiscoveryProfile %s", id), policyResourceNotSupportedError().Error())
	}

	// Create the resource using PATCH
	log.Printf("[INFO] Updating IPDiscoveryProfile with ID %s", id)
	boolFalse := false
	err := client.Patch(id, obj, &boolFalse)
	if err != nil {
		HandleUpdateError("IPDiscoveryProfile", id, err, response.Diagnostics)
		return
	}
}

func (r *PolicyIPDiscoveryProfileResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state policyIPDiscoveryProfileModel
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	id := state.NsxID.ValueString()
	if id == "" {
		response.Diagnostics.AddError("Error obtaining IPDiscoveryProfile ID", "Resource IPDiscoveryProfile ID field is not set in state")
		return
	}

	connector := GetPolicyConnector(r.client)
	boolFalse := false

	client := infra.NewIpDiscoveryProfilesClient(GetSessionContext(ctx, r.client, state.Context, response.Diagnostics), connector)
	err := client.Delete(id, &boolFalse)

	if err != nil {
		HandleDeleteError("IPDiscoveryProfile", id, err, response.Diagnostics)
		return
	}
}

func (r *PolicyIPDiscoveryProfileResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	NsxtPolicyPathResourceImporter(ctx, request, response)
}
