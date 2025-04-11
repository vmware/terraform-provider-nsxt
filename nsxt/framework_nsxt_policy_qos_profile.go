// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/customtypes"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_                                        resource.Resource              = &PolicyQOSProfileResource{}
	_                                        resource.ResourceWithConfigure = &PolicyQOSProfileResource{}
	rateFrameworkShaperScales                                               = []string{"mbps", "kbps", "mbps"}
	ingressFrameworkRateShaperIndex                                         = 0
	ingressFrameworkBroadcastRateShaperIndex                                = 1
	egressFrameworkRateShaperIndex                                          = 2
	frameworkRateShaperSchemaKeys                                           = []string{"ingress_rate_shaper", "ingress_broadcast_rate_shaper", "egress_rate_shaper"}
	frameworkRateShaperScales                                               = []string{"mbps", "kbps", "mbps"}
	frameworkRateShaperResourceTypes                                        = []string{"IngressRateShaper", "IngressBroadcastRateShaper", "EgressRateShaper"}
	frameworkRateLimiterResourceTypes                                       = []string{
		model.QosBaseRateLimiter_RESOURCE_TYPE_INGRESSRATELIMITER,
		model.QosBaseRateLimiter_RESOURCE_TYPE_INGRESSBROADCASTRATELIMITER,
		model.QosBaseRateLimiter_RESOURCE_TYPE_EGRESSRATELIMITER,
	}
)

var frameworkRateShaperSchemaNames []basetypes.ListValue
var frameworkRateShaperSchemaNamesPtr []*basetypes.ListValue

type policyQOSProfileModel struct {
	NsxID                      types.String `tfsdk:"nsx_id"`
	Path                       types.String `tfsdk:"path"`
	DisplayName                types.String `tfsdk:"display_name"`
	Description                types.String `tfsdk:"description"`
	Revision                   types.Int64  `tfsdk:"revision"`
	Tags                       types.List   `tfsdk:"tag"`
	Context                    types.Object `tfsdk:"context"`
	ClassOfService             types.Int64  `tfsdk:"class_of_service"`
	DscpTrusted                types.Bool   `tfsdk:"dscp_trusted"`
	DscpPriority               types.Int64  `tfsdk:"dscp_priority"`
	IngressRateShaper          types.List   `tfsdk:"ingress_rate_shaper"`
	IngressBroadcastRateShaper types.List   `tfsdk:"ingress_broadcast_rate_shaper"`
	EgressRateShaper           types.List   `tfsdk:"egress_rate_shaper"`
}

func NewPolicyQOSProfileResource() resource.Resource {
	return &PolicyQOSProfileResource{}
}

type PolicyQOSProfileResource struct {
	client interface{}
}

// Configure adds the provider configured client to the resource.
func (r *PolicyQOSProfileResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData
}

// Metadata returns the resource type name.
func (r *PolicyQOSProfileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy_qos_profile"
}

// Schema defines the schema for the resource.
func (r *PolicyQOSProfileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"nsx_id":       customtypes.GetNsxIDSchema(),
			"path":         customtypes.GetPathSchema(),
			"display_name": customtypes.GetDisplayNameSchema(),
			"description":  customtypes.GetDescriptionSchema(),
			"revision":     customtypes.GetRevisionSchema(),
			"class_of_service": schema.Int64Attribute{
				Description: "Class of service",
				Optional:    true,
				Validators:  []validator.Int64{int64validator.Between(0, 7)},
			},
			"dscp_trusted": schema.BoolAttribute{
				Description: "Trust mode for DSCP",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},

			"dscp_priority": schema.Int64Attribute{
				Description: "DSCP Priority",
				Optional:    true,
				Validators:  []validator.Int64{int64validator.Between(0, 63)},
			},
		},
		Blocks: map[string]schema.Block{
			"context":                       customtypes.GetContextSchema(),
			"tag":                           customtypes.GetTagsSchema(false, false),
			"ingress_rate_shaper":           getFrameworkQosRateShaperSchema(ingressFrameworkRateShaperIndex),
			"ingress_broadcast_rate_shaper": getFrameworkQosRateShaperSchema(ingressFrameworkBroadcastRateShaperIndex),
			"egress_rate_shaper":            getFrameworkQosRateShaperSchema(egressFrameworkRateShaperIndex),
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
// Create a new resource.
func (r *PolicyQOSProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var config policyQOSProfileModel
	diags := req.Plan.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
	connector := GetPolicyConnector(r.client)
	id := GetOrGenerateID2(ctx, r.client, config.NsxID, config.Context, resp.Diagnostics, resourceNsxtPolicyQOSProfileExistsForFramework)
	if resp.Diagnostics.HasError() {
		return
	}
	displayName := config.DisplayName.ValueString()
	description := config.Description.ValueString()
	tags := customtypes.GetPolicyTagsFromSchema(ctx, config.Tags, resp.Diagnostics)
	classOfService := config.ClassOfService.ValueInt64()
	dscpTrusted := "UNTRUSTED"
	if config.DscpTrusted.ValueBool() {
		dscpTrusted = "TRUSTED"
	}
	dscpPriority := config.DscpPriority.ValueInt64()
	var shapers []*data.StructValue

	frameworkRateShaperSchemaNames = []basetypes.ListValue{config.IngressRateShaper, config.IngressBroadcastRateShaper, config.EgressRateShaper}
	for index := ingressFrameworkRateShaperIndex; index <= egressFrameworkRateShaperIndex; index++ {

		shaper, err := getFrameworkPolicyQosRateShaperFromSchema(ctx, config, index, diags)
		if shaper != nil {
			shapers = append(shapers, shaper)
		}
		if err != nil {
			resp.Diagnostics.AddError(
				"getFrameworkPolicyQosRateShaperFromSchema error",
				fmt.Sprintf("Failed in getFrameworkPolicyQosRateShaperFromSchema : %v", err),
			)
			return
		}
	}
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"getFrameworkPolicyQosRateShaperFromSchema",
			"Failed in getFrameworkPolicyQosRateShaperFromSchema",
		)
		return
	}
	obj := model.QosProfile{
		DisplayName:    &displayName,
		Description:    &description,
		Tags:           tags,
		ClassOfService: &classOfService,
		Dscp: &model.QosDscp{
			Mode:     &dscpTrusted,
			Priority: &dscpPriority,
		},
		ShaperConfigurations: shapers,
	}
	// Create the resource using PATCH
	boolFalse := false
	client := infra.NewQosProfilesClient(GetSessionContext(ctx, r.client, config.Context, resp.Diagnostics), connector)
	if client == nil {
		resp.Diagnostics.AddError(
			"Error creating QosProfile client",
			"This NSX policy resource is not supported with given provider settings",
		)
		return

	}
	err := client.Patch(id, obj, &boolFalse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating QosProfile",
			fmt.Sprintf("Failed to create QosProfile %s: %s", id, err.Error()),
		)
		return

	}
	config.Revision = types.Int64PointerValue(obj.Revision)
	config.Path = types.StringPointerValue(obj.Path)
	config.NsxID = types.StringValue(id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)

}

// Read refreshes the Terraform state with the latest data.
// Read resource information.
func (r *PolicyQOSProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state policyQOSProfileModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	connector := GetPolicyConnector(r.client)
	var id string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("nsx_id"), &id)...)
	if id == "" {
		id = state.NsxID.ValueString()
	}
	if id == "" {
		resp.Diagnostics.AddError("Error obtaining QosProfile ID", "ID attribute not found")
		return
	}
	client := infra.NewQosProfilesClient(GetSessionContext(ctx, r.client, state.Context, resp.Diagnostics), connector)
	if client == nil {
		resp.Diagnostics.AddError(
			"Error creating QosProfile client",
			"This NSX policy resource is not supported with given provider settings",
		)
		return

	}
	obj, err := client.Get(id)
	if err != nil {
		HandleReadError("QosProfile", id, err, resp.Diagnostics)
		resp.State.RemoveResource(ctx)
	}

	state.DisplayName = types.StringPointerValue(obj.DisplayName)
	state.Description = types.StringPointerValue(obj.Description)
	state.Tags = customtypes.SetPolicyTagsInSchema(obj.Tags, resp.Diagnostics)
	state.Path = types.StringPointerValue(obj.Path)
	state.Revision = types.Int64PointerValue(obj.Revision)
	state.ClassOfService = types.Int64PointerValue(obj.ClassOfService)
	state.DscpTrusted = types.BoolValue(false)
	if *obj.Dscp.Mode == "TRUSTED" {
		state.DscpTrusted = types.BoolValue(true)
	}
	state.DscpPriority = types.Int64PointerValue(obj.Dscp.Priority)
	// TODO: Revisit again
	for _, shaperSchemaKeys := range obj.ShaperConfigurations {

		resourceType, temp, err := setFrameworkPolicyQosRateShaperInSchema(shaperSchemaKeys)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed in setFrameworkPolicyQosRateShaperInSchema",
				fmt.Sprintf("Failed in setFrameworkPolicyQosRateShaperInSchema : %v", err),
			)
			return
		}

		switch resourceType {
		case "IngressRateShaper":
			state.IngressRateShaper = temp
		case "IngressBroadcastRateShaper":
			state.IngressBroadcastRateShaper = temp
		case "EgressRateShaper":
			state.EgressRateShaper = temp
		}
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *PolicyQOSProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var state policyQOSProfileModel
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	connector := GetPolicyConnector(r.client)
	id := GetOrGenerateID2(ctx, r.client, state.NsxID, state.Context, resp.Diagnostics, resourceNsxtPolicyQOSProfileExistsForFramework)
	if resp.Diagnostics.HasError() {
		return
	}
	// resp.Diagnostics.AddError(
	// 	"xxxxxxxxxxxxxxxxxxxxxxxx",
	// 	fmt.Sprintf("yyyyyyyyyyyyyyyyyyyyyyyy  %v", "sca"),
	// )
	// return
	displayName := state.DisplayName.ValueString()
	description := state.Description.ValueString()
	tags := customtypes.GetPolicyTagsFromSchema(ctx, state.Tags, resp.Diagnostics)
	classOfService := state.ClassOfService.ValueInt64()
	dscpTrusted := "UNTRUSTED"
	if state.DscpTrusted.ValueBool() {
		dscpTrusted = "TRUSTED"
	}
	dscpPriority := state.DscpPriority.ValueInt64()
	var shapers []*data.StructValue
	frameworkRateShaperSchemaNames = []basetypes.ListValue{state.IngressRateShaper, state.IngressBroadcastRateShaper, state.EgressRateShaper}
	for index := ingressFrameworkRateShaperIndex; index <= egressFrameworkRateShaperIndex; index++ {

		shaper, _ := getFrameworkPolicyQosRateShaperFromSchema(ctx, state, index, diags)
		if shaper != nil {
			shapers = append(shapers, shaper)
		}
	}
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"getFrameworkPolicyQosRateShaperFromSchema",
			"Failed in getFrameworkPolicyQosRateShaperFromSchema",
		)
		return
	}
	obj := model.QosProfile{
		DisplayName:    &displayName,
		Description:    &description,
		Tags:           tags,
		ClassOfService: &classOfService,
		Dscp: &model.QosDscp{
			Mode:     &dscpTrusted,
			Priority: &dscpPriority,
		},
		ShaperConfigurations: shapers,
	}

	// Create the resource using PATCH
	boolFalse := false
	client := infra.NewQosProfilesClient(GetSessionContext(ctx, r.client, state.Context, resp.Diagnostics), connector)
	if client == nil {
		resp.Diagnostics.AddError(
			"Error creating QosProfile client",
			"This NSX policy resource is not supported with given provider settings",
		)
		return
		// return policyResourceNotSupportedError()
	}
	err := client.Patch(id, obj, &boolFalse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating QosProfile",
			fmt.Sprintf("Failed to create QosProfile %s: %s", id, err.Error()),
		)
		return
		// return handleCreateError("QosProfile", id, err)
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *PolicyQOSProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state policyQOSProfileModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	connector := GetPolicyConnector(r.client)
	boolFalse := false
	var id string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("nsx_id"), &id)...)
	if id == "" {
		id = state.NsxID.ValueString()
	}
	if id == "" {
		resp.Diagnostics.AddError("Error obtaining QosProfile ID", "ID attribute not found")
		return
	}
	client := infra.NewQosProfilesClient(GetSessionContext(ctx, r.client, state.Context, resp.Diagnostics), connector)
	if client == nil {
		resp.Diagnostics.AddError(
			"Error creating QosProfile client",
			"This NSX policy resource is not supported with given provider settings",
		)
		return
		// return policyResourceNotSupportedError()
	}
	err := client.Delete(id, &boolFalse)

	if err != nil {
		HandleDeleteError("QosProfile", id, err, resp.Diagnostics)
		return
	}
}

func getFrameworkQosRateShaperSchema(index int) schema.ListNestedBlock {
	scale := rateFrameworkShaperScales[index]
	return schema.ListNestedBlock{
		Validators: []validator.List{listvalidator.SizeAtMost(1)},
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"enabled": schema.BoolAttribute{
					Description: "Whether this rate shaper is enabled",
					Optional:    true,
					Computed:    true,
					Default:     booldefault.StaticBool(true),
				},
				fmt.Sprintf("average_bw_%s", scale): schema.Int64Attribute{
					Description: fmt.Sprintf("Average Bandwidth in %s", scale),
					Optional:    true,
				},
				"burst_size": schema.Int64Attribute{
					Description: "Burst size in bytes",
					Optional:    true,
				},
				fmt.Sprintf("peak_bw_%s", scale): schema.Int64Attribute{
					Description: fmt.Sprintf("Average Bandwidth in %s", scale),
					Optional:    true,
				},
			},
		},
	}
}

func resourceNsxtPolicyQOSProfileExistsForFramework(sessionContext utl.SessionContext, id string, connector client.Connector, diag diag.Diagnostics) bool {
	client := infra.NewQosProfilesClient(sessionContext, connector)
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

func getFrameworkPolicyQosRateShaperFromSchema(ctx context.Context, config policyQOSProfileModel, index int, diag diag.Diagnostics) (*data.StructValue, error) {
	scale := frameworkRateShaperScales[index]
	listValue := frameworkRateShaperSchemaNames[index]
	resourceType := frameworkRateLimiterResourceTypes[index]

	elemType := listValue.Elements()

	for _, shaperData := range elemType {
		// 	// only 1 is allowed for each type
		var serShaperData map[string]interface{}
		err := serializeShaperData(&serShaperData, shaperData.String())
		if err != nil {

		}

		enabled := serShaperData["enabled"].(bool)
		averageBW := int64(serShaperData[fmt.Sprintf("average_bw_%s", scale)].(float64))

		burstSize := int64(serShaperData["burst_size"].(float64))
		peakBW := int64(serShaperData[fmt.Sprintf("peak_bw_%s", scale)].(float64))

		shaper := model.IngressRateLimiter{
			ResourceType:     resourceType,
			Enabled:          &enabled,
			BurstSize:        &burstSize,
			AverageBandwidth: &averageBW,
			PeakBandwidth:    &peakBW,
		}
		converter := bindings.NewTypeConverter()
		dataValue, errList := converter.ConvertToVapi(shaper, model.IngressRateLimiterBindingType())
		if len(errList) > 0 {
			return nil, fmt.Errorf("Error while converting IngressRateLimiter: %v", errList[0])
		}

		return dataValue.(*data.StructValue), nil
	}

	return nil, nil
}

func serializeShaperData(data *map[string]interface{}, shaperStr string) error {

	err := json.Unmarshal([]byte(shaperStr), &data)
	if err != nil {
		return fmt.Errorf("Unable to parse shaperData :%v", err)
	}
	return nil
}

type IngressRateLimiter struct {
	AverageBandwidth *int64
	BurstSize        *int64
	PeakBandwidth    *int64
	Enabled          *bool
	ResourceType     string
}

func (m IngressRateLimiter) objectType(resourceType string) types.ObjectType {
	return types.ObjectType{AttrTypes: m.objectAttributeTypes(resourceType)}
}

func (m IngressRateLimiter) objectAttributeTypes(resourceType string) map[string]attr.Type {
	scale := "mbps"
	if resourceType == "IngressBroadcastRateShaper" {
		scale = "kbps"
	}
	return map[string]attr.Type{
		"enabled":                           types.BoolType,
		fmt.Sprintf("average_bw_%s", scale): types.Int64Type,
		"burst_size":                        types.Int64Type,
		fmt.Sprintf("peak_bw_%s", scale):    types.Int64Type,
	}
}

// TODO: Revisit again
func setFrameworkPolicyQosRateShaperInSchema(shaperConfs *data.StructValue) (string, basetypes.ListValue, error) {
	converter := bindings.NewTypeConverter()

	dataValue, _ := converter.ConvertToGolang(shaperConfs, model.IngressRateLimiterBindingType())
	shaper := dataValue.(model.IngressRateLimiter)
	scale := "mbps"
	if shaper.ResourceType == "IngressBroadcastRateShaper" {
		scale = "kbps"
	}
	t, _ := types.ObjectValue(IngressRateLimiter{}.objectAttributeTypes(shaper.ResourceType),
		map[string]attr.Value{
			"enabled":                           types.BoolPointerValue(shaper.Enabled),
			fmt.Sprintf("average_bw_%s", scale): types.Int64PointerValue(shaper.AverageBandwidth),
			"burst_size":                        types.Int64PointerValue(shaper.BurstSize),
			fmt.Sprintf("peak_bw_%s", scale):    types.Int64PointerValue(shaper.PeakBandwidth),
		})
	l, _ := types.ListValue(IngressRateLimiter{}.objectType(shaper.ResourceType), []attr.Value{t})

	return shaper.ResourceType, l, nil
}
