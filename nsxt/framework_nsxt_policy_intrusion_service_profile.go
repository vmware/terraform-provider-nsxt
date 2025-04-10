// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	services "github.com/vmware/terraform-provider-nsxt/api/infra/settings/firewall/security/intrusion_services"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/customtypes"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &PolicyIntrusionServiceProfileResource{}
	_ resource.ResourceWithConfigure = &PolicyIntrusionServiceProfileResource{}
	// _ resource.ResourceWithConfigValidators = &PolicyIntrusionServiceProfileResource{}
	// TODO: revert this
	idsFrameworkProfileSeverityValues = []string{
		model.IdsProfile_PROFILE_SEVERITY_MEDIUM,
		model.IdsProfile_PROFILE_SEVERITY_HIGH,
		model.IdsProfile_PROFILE_SEVERITY_LOW,
		model.IdsProfile_PROFILE_SEVERITY_CRITICAL,
		model.IdsProfile_PROFILE_SEVERITY_SUSPICIOUS,
	}
	// idsFrameworkProfileSeverityValues = []string{"CRITICAL", "HIGH", "MEDIUM", "LOW", "SUSPICIOUS"}
)

type policyIntrusionServiceProfileModel struct {
	NsxID               types.String `tfsdk:"nsx_id"`
	Path                types.String `tfsdk:"path"`
	DisplayName         types.String `tfsdk:"display_name"`
	Description         types.String `tfsdk:"description"`
	Revision            types.Int64  `tfsdk:"revision"`
	Tags                types.List   `tfsdk:"tag"`
	Context             types.Object `tfsdk:"context"`
	Criteria            types.List   `tfsdk:"criteria"`
	OverriddenSignature types.Set    `tfsdk:"overridden_signature"`
	Severities          types.Set    `tfsdk:"severities"`
}

func NewPolicyIntrusionServiceProfileResource() resource.Resource {
	return &PolicyIntrusionServiceProfileResource{}
}

type PolicyIntrusionServiceProfileResource struct {
	client interface{}
}

// Configure adds the provider configured client to the resource.
func (r *PolicyIntrusionServiceProfileResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData
}

// Metadata returns the resource type name.
func (r *PolicyIntrusionServiceProfileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy_intrusion_service_profile"
}

// Schema defines the schema for the resource.
func (r *PolicyIntrusionServiceProfileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"nsx_id":       customtypes.GetNsxIDSchema(),
			"path":         customtypes.GetPathSchema(),
			"display_name": customtypes.GetDisplayNameSchema(),
			"description":  customtypes.GetDescriptionSchema(),
			"revision":     customtypes.GetRevisionSchema(),

			"severities": schema.SetAttribute{
				Required:    true,
				Description: "Severities of signatures which are part of this profile",
				ElementType: types.StringType,
				Validators: []validator.Set{setvalidator.ValueStringsAre(
					stringvalidator.OneOf(idsFrameworkProfileSeverityValues...),
				)},
			},
		},
		Blocks: map[string]schema.Block{
			"context":              customtypes.GetContextSchema(),
			"tag":                  customtypes.GetTagsSchema(false, false),
			"criteria":             getFrameworkIdsProfileCriteriaSchema(),
			"overridden_signature": getFrameworkIdsProfileSignatureSchema(),
		},
	}
}

// func (r *PolicyIntrusionServiceProfileResource) ConfigValidators() resource.ConfigValidator {
// }
func getFrameworkIdsProfileSignatureSchema() schema.SetNestedBlock {
	return schema.SetNestedBlock{
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"signature_id": schema.StringAttribute{
					Description: "List of attack type criteria",
					Required:    true,
				},
				"enabled": schema.BoolAttribute{
					Default:  booldefault.StaticBool(true),
					Optional: true,
					Computed: true,
				},
				"action": schema.StringAttribute{
					Description: "This will take precedence over IDS signature action",
					Optional:    true,
				},
			},
		},
	}
}

func getFrameworkIdsProfileCriteriaSchema() schema.ListNestedBlock {
	return schema.ListNestedBlock{
		// Optional: true,
		// Validators: []validator.List{listvalidator.AtLeastOneOf(
		// 	path.MatchRoot("criteria").AtName("attack_types"),
		// 	path.MatchRoot("criteria").AtMapKey("attack_targets"),
		// 	path.MatchRoot("criteria").AtMapKey("cvss"),
		// 	path.MatchRoot("criteria").AtMapKey("products_affected"),
		// )},
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"attack_types": schema.SetAttribute{
					Description: "List of attack type criteria",
					Optional:    true,
					ElementType: types.StringType,
				},
				"attack_targets": schema.SetAttribute{
					Description: "List of attack target criteria",
					Optional:    true,
					ElementType: types.StringType,
				},
				"cvss": schema.SetAttribute{
					Description: "Common Vulnerability Scoring System Ranges",
					Optional:    true,
					ElementType: types.StringType,
				},
				"products_affected": schema.SetAttribute{
					Description: "List of products affected",
					Optional:    true,
					ElementType: types.StringType,
				},
			},
		},
	}
}

func getFrameworkIdsProfileSignaturesFromSchema(ctx context.Context, config policyIntrusionServiceProfileModel) ([]model.IdsProfileLocalSignature, error) {
	var result []model.IdsProfileLocalSignature

	// signatures := d.Get("overridden_signature").(*schema.Set).List()
	signatures := config.OverriddenSignature.Elements()
	// if len(signatures) == 0 {
	// 	return result, fmt.Errorf("No signatures provided")
	// }
	if len(signatures) == 0 {
		return result, nil
	}

	for _, item := range signatures {
		// sig, err := item.ToTerraformValue(ctx)
		// sig.Type()
		var dataMap (map[string]interface{})
		err := serializeShaperData(&dataMap, item.String())
		if err != nil {

		}
		action := dataMap["action"].(string)
		enabled := dataMap["enabled"].(bool)
		id := dataMap["signature_id"].(string)
		signature := model.IdsProfileLocalSignature{
			Action:      &action,
			Enable:      &enabled,
			SignatureId: &id,
		}

		result = append(result, signature)
	}

	return result, nil
}

func buildFrameworkIdsProfileCriteriaFilter(name string, values []string) (*data.StructValue, error) {
	converter := bindings.NewTypeConverter()
	item := model.IdsProfileFilterCriteria{
		FilterName:   &name,
		FilterValue:  values,
		ResourceType: model.IdsProfileCriteria_RESOURCE_TYPE_IDSPROFILEFILTERCRITERIA,
	}

	dataValue, errs := converter.ConvertToVapi(item, model.IdsProfileFilterCriteriaBindingType())
	if errs != nil {
		return nil, errs[0]
	}
	return dataValue.(*data.StructValue), nil
}

func buildFrameworkIdsProfileCriteriaOperator() (*data.StructValue, error) {
	op := "AND"
	operator := model.IdsProfileConjunctionOperator{
		Operator:     &op,
		ResourceType: model.IdsProfileCriteria_RESOURCE_TYPE_IDSPROFILECONJUNCTIONOPERATOR,
	}

	converter := bindings.NewTypeConverter()

	dataValue, errs := converter.ConvertToVapi(operator, model.IdsProfileConjunctionOperatorBindingType())
	if errs != nil {
		return nil, errs[0]
	}

	return dataValue.(*data.StructValue), nil
}

func serializeCriteriaData(criteriaData map[string][]string, shaperStr string) error {
	data := make(map[string]interface{})
	err := json.Unmarshal([]byte(shaperStr), &data)
	if err != nil {
		return fmt.Errorf("Unable to parse shaperData :%v", err)
	}

	for k, v := range data {
		l := []string{}
		for _, i := range v.([]interface{}) {
			l = append(l, i.(string))
		}
		if len(l) > 0 {
			criteriaData[k] = l
		}
	}
	return nil
}

func serializeCriteriaData1(val attr.Value) (map[string][]string, error) {
	// Check if value is null or unknown
	if val.IsNull() || val.IsUnknown() {
		return nil, nil
	}

	// Assert value is Object
	objVal, ok := val.(types.Object)
	if !ok {
		return nil, fmt.Errorf("expected Object type")
	}

	result := make(map[string][]string)

	for key, attrVal := range objVal.Attributes() {
		if attrVal.IsNull() || attrVal.IsUnknown() {
			continue
		}

		// Expect each field to be a Set of strings
		setVal, ok := attrVal.(types.Set)
		if !ok {
			return nil, fmt.Errorf("expected Set type for key %q", key)
		}

		var items []string
		for _, elem := range setVal.Elements() {
			strVal, ok := elem.(types.String)
			if !ok {
				return nil, fmt.Errorf("expected string in set for key %q", key)
			}
			items = append(items, strVal.ValueString())
		}

		result[key] = items
	}

	return result, nil
}

func getFrameworkIdsProfileCriteriaFromSchema(ctx context.Context, config policyIntrusionServiceProfileModel, diag diag.Diagnostics) ([]*data.StructValue, diag.Diagnostics) {
	var result []*data.StructValue

	criteria := config.Criteria.Elements()
	if len(criteria) == 0 || criteria[0] == nil {
		return result, nil
	}

	criteriaData, err := serializeCriteriaData1(criteria[0])
	if err != nil {
		diag.AddError("getFrameworkIdsProfileCriteriaFromSchema", fmt.Sprintf("Failed while serializing the criteria : %v--%v", err, criteriaData))
		return []*data.StructValue{}, diag
	}
	var attackTypes []string
	if criteriaData["attack_types"] != nil {
		attackTypes = criteriaData["attack_types"]
	}
	var attackTargets []string
	if criteriaData["attack_targets"] != nil {
		attackTargets = criteriaData["attack_targets"]
	}
	var cvss []string
	if len(criteriaData["cvss"]) > 0 {
		cvss = criteriaData["cvss"]
	}
	var productsAffected []string
	if criteriaData["products_affected"] != nil {
		productsAffected = criteriaData["products_affected"]
	}

	if len(attackTypes) > 0 {
		item, err := buildFrameworkIdsProfileCriteriaFilter(model.IdsProfileFilterCriteria_FILTER_NAME_ATTACK_TYPE, attackTypes)
		if err != nil {
			diag.AddError("getFrameworkIdsProfileCriteriaFromSchema", fmt.Sprintf("Failed in buildFrameworkIdsProfileCriteriaFilter : %v", err.Error()))
			return result, diag
		}

		result = append(result, item)

		operator, err1 := buildFrameworkIdsProfileCriteriaOperator()
		if err1 != nil {
			diag.AddError("getFrameworkIdsProfileCriteriaFromSchema", fmt.Sprintf("Failed in buildFrameworkIdsProfileCriteriaOperator : %v", err.Error()))
			return result, diag
		}
		result = append(result, operator)
	}

	if len(attackTargets) > 0 {
		item, err := buildFrameworkIdsProfileCriteriaFilter(model.IdsProfileFilterCriteria_FILTER_NAME_ATTACK_TARGET, attackTargets)
		if err != nil {
			diag.AddError("getFrameworkIdsProfileCriteriaFromSchema", fmt.Sprintf("Failed in buildFrameworkIdsProfileCriteriaFilter : %v", err.Error()))
			return result, diag
		}

		result = append(result, item)

		operator, err1 := buildFrameworkIdsProfileCriteriaOperator()
		if err1 != nil {
			diag.AddError("getFrameworkIdsProfileCriteriaFromSchema", fmt.Sprintf("Failed in buildFrameworkIdsProfileCriteriaOperator : %v", err.Error()))
			return result, diag
		}
		result = append(result, operator)
	}

	if len(cvss) > 0 {
		item, err := buildFrameworkIdsProfileCriteriaFilter(model.IdsProfileFilterCriteria_FILTER_NAME_CVSS, cvss)
		if err != nil {
			diag.AddError("getFrameworkIdsProfileCriteriaFromSchema", fmt.Sprintf("Failed in buildFrameworkIdsProfileCriteriaFilter : %v", err.Error()))
			return result, diag
		}

		result = append(result, item)

		operator, err1 := buildFrameworkIdsProfileCriteriaOperator()
		if err1 != nil {
			diag.AddError("getFrameworkIdsProfileCriteriaFromSchema", fmt.Sprintf("Failed in buildFrameworkIdsProfileCriteriaOperator : %v", err.Error()))
			return result, diag
		}
		result = append(result, operator)
	}

	if len(productsAffected) > 0 {
		item, err := buildFrameworkIdsProfileCriteriaFilter(model.IdsProfileFilterCriteria_FILTER_NAME_PRODUCT_AFFECTED, productsAffected)
		if err != nil {
			diag.AddError("getFrameworkIdsProfileCriteriaFromSchema", fmt.Sprintf("Failed in buildFrameworkIdsProfileCriteriaFilter : %v", err.Error()))
			return result, diag
		}

		result = append(result, item)

		operator, err1 := buildFrameworkIdsProfileCriteriaOperator()
		if err1 != nil {
			diag.AddError("getFrameworkIdsProfileCriteriaFromSchema", fmt.Sprintf("Failed in buildFrameworkIdsProfileCriteriaOperator : %v", err.Error()))
			return result, diag
		}
		result = append(result, operator)
	}
	return result[:len(result)-1], nil
}

func getFrameworkIdsProfileSeverities(severities basetypes.SetValue) []string {
	result := make([]string, 0, len(severities.Elements()))
	for _, severity := range severities.Elements() {
		for _, i := range idsFrameworkProfileSeverityValues {
			if strings.ReplaceAll(severity.String(), `"`, ``) == string(i) {

				result = append(result, i)
			}
		}
	}
	return result
}

func resourceFrameworkNsxtPolicyIntrusionServiceProfileExists(sessionContext utl.SessionContext, id string, connector client.Connector, diag diag.Diagnostics) bool {
	var err error
	client := services.NewProfilesClient(sessionContext, connector)
	if client == nil {
		return false
	}
	_, err = client.Get(id)
	if err == nil {
		return true
	}

	if isNotFoundError(err) {
		return false
	}

	return false
}

// Create creates the resource and sets the initial Terraform state.
// Create a new resource.
func (r *PolicyIntrusionServiceProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var config policyIntrusionServiceProfileModel
	diags := req.Plan.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
	connector := GetPolicyConnector(r.client)
	id := GetOrGenerateID2(ctx, r.client, config.NsxID, config.Context, resp.Diagnostics, resourceFrameworkNsxtPolicyIntrusionServiceProfileExists)
	if resp.Diagnostics.HasError() {
		return
	}

	displayName := config.DisplayName.ValueString()
	description := config.Description.ValueString()
	tags := customtypes.GetPolicyTagsFromSchema(ctx, config.Tags, resp.Diagnostics)
	criteria, diags := getFrameworkIdsProfileCriteriaFromSchema(ctx, config, resp.Diagnostics)

	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"getFrameworkIdsProfileCriteriaFromSchema",
			fmt.Sprintf("Failed in getFrameworkIdsProfileCriteriaFromSchema"),
		)
		return
	}
	signatures, err := getFrameworkIdsProfileSignaturesFromSchema(ctx, config)
	if err != nil {
		resp.Diagnostics.AddError(
			"getFrameworkIdsProfileSignaturesFromSchema",
			fmt.Sprintf("Failed in getFrameworkIdsProfileSignaturesFromSchema : %v", err.Error()),
		)
		return
	}
	frameworkProfileSeverity := getFrameworkIdsProfileSeverities(config.Severities)

	obj := model.IdsProfile{
		DisplayName:          &displayName,
		Description:          &description,
		Tags:                 tags,
		Criteria:             criteria,
		OverriddenSignatures: signatures,
		ProfileSeverity:      frameworkProfileSeverity,
	}
	client := services.NewProfilesClient(GetSessionContext(ctx, r.client, config.Context, resp.Diagnostics), connector)
	if client == nil {
		resp.Diagnostics.AddError(
			"Failed to create client object",
			policyResourceNotSupportedError().Error(),
		)
		return
	}

	err = client.Patch(id, obj)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create resource",
			handleCreateError("Ids Profile", id, err).Error(),
		)
		return
	}

	config.Revision = types.Int64PointerValue(obj.Revision)
	config.Path = types.StringPointerValue(obj.Path)
	config.NsxID = types.StringValue(id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)

}

// Read refreshes the Terraform state with the latest data.

func (r *PolicyIntrusionServiceProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state

	var state policyIntrusionServiceProfileModel
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
		resp.Diagnostics.AddError("Error obtaining Ids Profile ID", "ID attribute not found")
		return
	}
	client := services.NewProfilesClient(GetSessionContext(ctx, r.client, state.Context, resp.Diagnostics), connector)
	if client == nil {
		resp.Diagnostics.AddError(
			"Error creating Ids Profile client",
			"This NSX policy resource is not supported with given provider settings",
		)
		return

	}
	obj, err := client.Get(id)
	if err != nil {
		HandleReadError("IdsProfile", id, err, resp.Diagnostics)
		resp.State.RemoveResource(ctx)
	}

	state.DisplayName = types.StringPointerValue(obj.DisplayName)
	state.Description = types.StringPointerValue(obj.Description)
	state.Tags = customtypes.SetPolicyTagsInSchema(obj.Tags, resp.Diagnostics)
	state.Path = types.StringPointerValue(obj.Path)
	state.Revision = types.Int64PointerValue(obj.Revision)

	state.Criteria, diags = setFrameworkIdsProfileCriteriaInSchema(state, obj.Criteria, resp.Diagnostics)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Severities, diags = setFrameworkIdsProfileSeverities(state, obj.ProfileSeverity)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func setFrameworkIdsProfileSeverities(state policyIntrusionServiceProfileModel, severities []string) (basetypes.SetValue, diag.Diagnostics) {
	var elements []attr.Value
	for _, severity := range severities {
		elements = append(elements, basetypes.NewStringValue(severity))
	}
	return basetypes.NewSetValue(types.StringType, elements)
}

func convertMapListToListValue(input []map[string][]attr.Value) (types.List, diag.Diagnostics) {
	// Define the object type that matches your schema
	criteriaObjectType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"attack_types":      types.SetType{ElemType: types.StringType},
			"attack_targets":    types.SetType{ElemType: types.StringType},
			"cvss":              types.SetType{ElemType: types.StringType},
			"products_affected": types.SetType{ElemType: types.StringType},
		},
	}

	var objs []attr.Value
	converted := make(map[string]attr.Value)
	var diag diag.Diagnostics

	for _, item := range input {

		for key, valSlice := range item {

			setVal, d := basetypes.NewSetValue(types.StringType, valSlice)

			diag.Append(d...)
			if diag.HasError() {
				return types.List{}, diag
			}
			converted[key] = setVal

			diag.Append(d...)
			if diag.HasError() {
				return types.List{}, diag
			}
		}
	}
	for key, _ := range criteriaObjectType.AttrTypes {
		if _, ok := converted[key]; !ok {
			converted[key] = basetypes.NewSetNull(types.StringType)
		}
	}
	objVal, d := types.ObjectValue(criteriaObjectType.AttrTypes, converted)

	diag.Append(d...)
	if diag.HasError() {
		return types.List{}, diag
	}

	objs = append(objs, objVal)
	listValue, d := types.ListValue(criteriaObjectType, objs)
	diag.Append(d...)
	return listValue, diag
}

// TODO: revist
func setFrameworkIdsProfileCriteriaInSchema(state policyIntrusionServiceProfileModel, criteriaList []*data.StructValue, diag diag.Diagnostics) (basetypes.ListValue, diag.Diagnostics) {
	var schemaList []map[string][]attr.Value
	converter := bindings.NewTypeConverter()
	criteriaMap := make(map[string][]attr.Value)
	var dataValues []attr.Value
	for i, item := range criteriaList {
		// Odd elements are AND operators - ignoring them
		if i%2 == 1 {
			continue
		}

		dataValue, errs := converter.ConvertToGolang(item, model.IdsProfileFilterCriteriaBindingType())
		if errs != nil {
			diag.AddError("Failed while converting *data.StructValue to interface{}", fmt.Sprintf("Error: %v", errs))
		}

		criteria := dataValue.(model.IdsProfileFilterCriteria)
		for _, at := range criteria.FilterValue {
			dataValues = append(dataValues, basetypes.NewStringValue(at))
		}
		if *criteria.FilterName == model.IdsProfileFilterCriteria_FILTER_NAME_ATTACK_TYPE {

			criteriaMap["attack_types"] = dataValues

		} else if *criteria.FilterName == model.IdsProfileFilterCriteria_FILTER_NAME_ATTACK_TARGET {

			criteriaMap["attack_targets"] = dataValues

		} else if *criteria.FilterName == model.IdsProfileFilterCriteria_FILTER_NAME_CVSS {

			criteriaMap["cvss"] = dataValues

		} else if *criteria.FilterName == model.IdsProfileFilterCriteria_FILTER_NAME_PRODUCT_AFFECTED {

			criteriaMap["products_affected"] = dataValues

		}
	}

	if len(criteriaMap) > 0 {
		schemaList = append(schemaList, criteriaMap)
	}
	return convertMapListToListValue(schemaList)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *PolicyIntrusionServiceProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var config policyIntrusionServiceProfileModel
	diags := req.Plan.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
	connector := GetPolicyConnector(r.client)
	id := GetOrGenerateID2(ctx, r.client, config.NsxID, config.Context, resp.Diagnostics, resourceFrameworkNsxtPolicyIntrusionServiceProfileExists)
	if resp.Diagnostics.HasError() {
		return
	}

	displayName := config.DisplayName.ValueString()
	description := config.Description.ValueString()
	tags := customtypes.GetPolicyTagsFromSchema(ctx, config.Tags, resp.Diagnostics)
	criteria, diags := getFrameworkIdsProfileCriteriaFromSchema(ctx, config, resp.Diagnostics)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"getFrameworkIdsProfileCriteriaFromSchema",
			fmt.Sprintf("Failed in getFrameworkIdsProfileCriteriaFromSchema"),
		)
		return
	}
	signatures, err := getFrameworkIdsProfileSignaturesFromSchema(ctx, config)
	if err != nil {
		resp.Diagnostics.AddError(
			"getFrameworkIdsProfileSignaturesFromSchema",
			fmt.Sprintf("Failed in getFrameworkIdsProfileSignaturesFromSchema : %v", err.Error()),
		)
		return
	}
	profileSeverity := getFrameworkIdsProfileSeverities(config.Severities)

	obj := model.IdsProfile{
		DisplayName:          &displayName,
		Description:          &description,
		Tags:                 tags,
		Criteria:             criteria,
		OverriddenSignatures: signatures,
		ProfileSeverity:      profileSeverity,
	}
	client := services.NewProfilesClient(GetSessionContext(ctx, r.client, config.Context, resp.Diagnostics), connector)
	if client == nil {
		resp.Diagnostics.AddError(
			"Failed to create client object",
			policyResourceNotSupportedError().Error(),
		)
		return
	}
	err = client.Patch(id, obj)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create resource",
			handleCreateError("Ids Profile", id, err).Error(),
		)
		return
	}

	config.Revision = types.Int64PointerValue(obj.Revision)
	config.Path = types.StringPointerValue(obj.Path)
	config.NsxID = types.StringValue(id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)

}

// Delete deletes the resource and removes the Terraform state on success.
func (r *PolicyIntrusionServiceProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state policyIntrusionServiceProfileModel
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
		resp.Diagnostics.AddError("Error obtaining Intrusion Service Profile ID", "ID attribute not found")
		return
	}
	client := services.NewProfilesClient(GetSessionContext(ctx, r.client, state.Context, resp.Diagnostics), connector)
	if client == nil {
		resp.Diagnostics.AddError(
			"Error creating Intrusion Service Profile client",
			"This NSX policy resource is not supported with given provider settings",
		)
		return

	}
	err := client.Delete(id)

	if err != nil {
		HandleDeleteError("IntrusionServiceProfile", id, err, resp.Diagnostics)
		return
	}
}
