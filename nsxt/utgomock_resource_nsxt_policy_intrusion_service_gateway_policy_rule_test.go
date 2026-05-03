//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	intrusionservicegatewaypolicies "github.com/vmware/terraform-provider-nsxt/api/infra/domains/intrusion_service_gateway_policies"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	idsgwrulemocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/domains/intrusion_service_gateway_policies"
)

var (
	idsGwRuleID          = "ids-gw-rule-id"
	idsGwRuleDisplayName = "Test IDS Gateway Policy Rule"
	idsGwRuleDescription = "Test Intrusion Service Gateway Policy Rule"
	idsGwRuleDomain      = "default"
	idsGwRuleAction      = "DETECT"
	idsGwRuleDirection   = "IN_OUT"
	idsGwRuleIPVersion   = "IPV4_IPV6"
	idsGwRuleLogged      = false
	idsGwRuleDisabled    = false
	idsGwRuleSeqNum      = int64(0)
	idsGwRuleRuleID      = int64(2001)
	idsGwRuleScope       = []string{"/infra/tier-1s/test-t1-gw"}
	idsGwRuleProfilePath = []string{"/infra/settings/firewall/security/intrusion-services/profiles/DefaultIDSProfile"}
	idsGwRulePolicyID    = "ids-gw-policy-id"
	idsGwRulePolicyPath  = "/infra/domains/default/intrusion-service-gateway-policies/ids-gw-policy-id"

	// Updated values for testing modifications
	idsGwRuleUpdatedDirection   = "IN"
	idsGwRuleUpdatedIPVersion   = "IPV4"
	idsGwRuleUpdatedLogged      = true
	idsGwRuleUpdatedDisabled    = true
	idsGwRuleUpdatedSeqNum      = int64(1)
	idsGwRuleUpdatedAction      = "DETECT_PREVENT"
	idsGwRuleUpdatedDisplayName = "Test IDS Gateway Policy Rule Update"
	idsGwRuleUpdatedDescription = "Test Intrusion Service Gateway Policy Rule Update"
)

// idsGwRuleAPIResponse creates a mock API response for IDS Gateway Rule with default values.
// This represents what the NSX-T API would return for a successfully created/retrieved rule.
func idsGwRuleAPIResponse() nsxModel.IdsRule {
	resourceType := "IdsRule"
	srcEx := false
	dstEx := false
	return nsxModel.IdsRule{
		Id:                   &idsGwRuleID,
		DisplayName:          &idsGwRuleDisplayName,
		Description:          &idsGwRuleDescription,
		ResourceType:         &resourceType,
		Action:               &idsGwRuleAction,
		Direction:            &idsGwRuleDirection,
		IpProtocol:           &idsGwRuleIPVersion,
		RuleId:               &idsGwRuleRuleID,
		Scope:                idsGwRuleScope,
		IdsProfiles:          idsGwRuleProfilePath,
		SequenceNumber:       &idsGwRuleSeqNum,
		Logged:               &idsGwRuleLogged,
		Disabled:             &idsGwRuleDisabled,
		SourcesExcluded:      &srcEx,
		DestinationsExcluded: &dstEx,
	}
}

// idsGwRuleAPIResponseUpdated creates a mock API response for IDS Gateway Rule with updated values.
// This represents what the NSX-T API would return after a successful update operation.
func idsGwRuleAPIResponseUpdated() nsxModel.IdsRule {
	resourceType := "IdsRule"
	srcEx := false
	dstEx := false
	return nsxModel.IdsRule{
		Id:                   &idsGwRuleID,
		DisplayName:          &idsGwRuleUpdatedDisplayName,
		Description:          &idsGwRuleUpdatedDescription,
		ResourceType:         &resourceType,
		Action:               &idsGwRuleUpdatedAction,
		Direction:            &idsGwRuleUpdatedDirection,
		IpProtocol:           &idsGwRuleUpdatedIPVersion,
		RuleId:               &idsGwRuleRuleID,
		Scope:                idsGwRuleScope,
		IdsProfiles:          idsGwRuleProfilePath,
		SequenceNumber:       &idsGwRuleUpdatedSeqNum,
		Logged:               &idsGwRuleUpdatedLogged,
		Disabled:             &idsGwRuleUpdatedDisabled,
		SourcesExcluded:      &srcEx,
		DestinationsExcluded: &dstEx,
	}
}

func minimalIdsGwRuleData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":    idsGwRuleDisplayName,
		"description":     idsGwRuleDescription,
		"policy_path":     idsGwRulePolicyPath,
		"action":          idsGwRuleAction,
		"direction":       idsGwRuleDirection,
		"ip_version":      idsGwRuleIPVersion,
		"sequence_number": 1,
		"logged":          false,
		"disabled":        false,
	}
}

// setIdsGwRuleProfiles sets the ids_profiles TypeSet on d (must be called after TestResourceDataRaw).
// This helper function properly handles the TypeSet creation for IDS profiles, which cannot be done
// directly in the minimalIdsGwRuleData() map due to TypeSet requirements.
func setIdsGwRuleProfiles(d *schema.ResourceData, t *testing.T) {
	profilesSet := schema.NewSet(schema.HashString, []interface{}{
		"/infra/settings/firewall/security/intrusion-services/profiles/DefaultIDSProfile",
	})
	if err := d.Set("ids_profiles", profilesSet); err != nil {
		t.Fatal(err)
	}
}

// setIdsGwRuleScope sets the scope TypeSet on d (must be called after TestResourceDataRaw).
// This helper function properly handles the TypeSet creation for Gateway scope paths, which is
// required for Gateway rules but optional for DFW rules.
func setIdsGwRuleScope(d *schema.ResourceData, t *testing.T) {
	scopeSet := schema.NewSet(schema.HashString, []interface{}{
		"/infra/tier-1s/test-t1-gw",
	})
	if err := d.Set("scope", scopeSet); err != nil {
		t.Fatal(err)
	}
}

func setupIdsGwRuleMock(t *testing.T, ctrl *gomock.Controller) (*idsgwrulemocks.MockRulesClient, func()) {
	t.Helper()
	mockSDK := idsgwrulemocks.NewMockRulesClient(ctrl)
	mockWrapper := &intrusionservicegatewaypolicies.IntrusionServiceGatewayRuleClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	original := cliIntrusionServiceGatewayPolicyRulesClient
	cliIntrusionServiceGatewayPolicyRulesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *intrusionservicegatewaypolicies.IntrusionServiceGatewayRuleClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliIntrusionServiceGatewayPolicyRulesClient = original }
}

func TestMockResourceNsxtPolicyIntrusionServiceGatewayPolicyRuleCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIdsGwRuleMock(t, ctrl)
	defer restore()

	t.Run("Create success with auto-generated ID", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Patch(idsGwRuleDomain, idsGwRulePolicyID, gomock.Any(), gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(idsGwRuleDomain, idsGwRulePolicyID, gomock.Any()).Return(idsGwRuleAPIResponse(), nil),
		)

		res := resourceNsxtPolicyIntrusionServiceGatewayPolicyRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsGwRuleData())
		setIdsGwRuleProfiles(d, t)

		err := resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
	})

	t.Run("Create fails when patch fails", func(t *testing.T) {
		mockSDK.EXPECT().Patch(idsGwRuleDomain, idsGwRulePolicyID, gomock.Any(), gomock.Any()).Return(vapiErrors.NewInvalidRequest())

		res := resourceNsxtPolicyIntrusionServiceGatewayPolicyRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsGwRuleData())
		setIdsGwRuleProfiles(d, t)

		err := resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("Create fails when get after patch fails", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Patch(idsGwRuleDomain, idsGwRulePolicyID, gomock.Any(), gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(idsGwRuleDomain, idsGwRulePolicyID, gomock.Any()).Return(nsxModel.IdsRule{}, vapiErrors.NewInternalServerError()),
		)

		res := resourceNsxtPolicyIntrusionServiceGatewayPolicyRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsGwRuleData())
		setIdsGwRuleProfiles(d, t)

		err := resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIntrusionServiceGatewayPolicyRuleRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIdsGwRuleMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(idsGwRuleDomain, idsGwRulePolicyID, idsGwRuleID).Return(idsGwRuleAPIResponse(), nil)

		res := resourceNsxtPolicyIntrusionServiceGatewayPolicyRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsGwRuleData())
		setIdsGwRuleProfiles(d, t)
		d.SetId(idsGwRuleID)

		err := resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, idsGwRuleLogged, d.Get("logged"))
		assert.Equal(t, idsGwRuleDisabled, d.Get("disabled"))
		assert.Equal(t, int(idsGwRuleSeqNum), d.Get("sequence_number"))
		assert.Equal(t, idsGwRuleDisplayName, d.Get("display_name"))
		assert.Equal(t, idsGwRuleDescription, d.Get("description"))
		assert.Equal(t, idsGwRuleAction, d.Get("action"))
		assert.Equal(t, idsGwRuleDirection, d.Get("direction"))
		assert.Equal(t, idsGwRuleIPVersion, d.Get("ip_version"))

		// Verify scope is properly handled
		scopeSet := d.Get("scope").(*schema.Set)
		assert.Equal(t, 1, scopeSet.Len())
		assert.True(t, scopeSet.Contains("/infra/tier-1s/test-t1-gw"))

		// Verify ids profile is properly handled
		idsProfileSet := d.Get("ids_profiles").(*schema.Set)
		assert.Equal(t, 1, idsProfileSet.Len())
		assert.True(t, idsProfileSet.Contains("/infra/settings/firewall/security/intrusion-services/profiles/DefaultIDSProfile"))
	})

	t.Run("Read not found", func(t *testing.T) {
		mockSDK.EXPECT().Get(idsGwRuleDomain, idsGwRulePolicyID, idsGwRuleID).Return(nsxModel.IdsRule{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyIntrusionServiceGatewayPolicyRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsGwRuleData())
		setIdsGwRuleProfiles(d, t)
		d.SetId(idsGwRuleID)

		err := resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIntrusionServiceGatewayPolicyRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsGwRuleData())
		setIdsGwRuleProfiles(d, t)
		// Don't set ID

		err := resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "ID")
		assert.Contains(t, err.Error(), "Intrusion Service Gateway Policy Rule")
	})

	t.Run("Read fails on API error", func(t *testing.T) {
		mockSDK.EXPECT().Get(idsGwRuleDomain, idsGwRulePolicyID, idsGwRuleID).Return(nsxModel.IdsRule{}, vapiErrors.NewInternalServerError())

		res := resourceNsxtPolicyIntrusionServiceGatewayPolicyRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsGwRuleData())
		setIdsGwRuleProfiles(d, t)
		d.SetId(idsGwRuleID)

		err := resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIntrusionServiceGatewayPolicyRuleUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIdsGwRuleMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Patch(idsGwRuleDomain, idsGwRulePolicyID, idsGwRuleID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(idsGwRuleDomain, idsGwRulePolicyID, idsGwRuleID).Return(idsGwRuleAPIResponseUpdated(), nil),
		)

		updatedData := minimalIdsGwRuleData()
		updatedData["sequence_number"] = int(idsGwRuleUpdatedSeqNum)
		updatedData["logged"] = idsGwRuleUpdatedLogged
		updatedData["disabled"] = idsGwRuleUpdatedDisabled
		updatedData["action"] = idsGwRuleUpdatedAction
		updatedData["display_name"] = idsGwRuleUpdatedDisplayName
		updatedData["description"] = idsGwRuleUpdatedDescription
		updatedData["direction"] = idsGwRuleUpdatedDirection
		updatedData["ip_version"] = idsGwRuleUpdatedIPVersion

		res := resourceNsxtPolicyIntrusionServiceGatewayPolicyRule()
		d := schema.TestResourceDataRaw(t, res.Schema, updatedData)
		setIdsGwRuleProfiles(d, t)
		d.SetId(idsGwRuleID)

		err := resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, int(idsGwRuleUpdatedSeqNum), d.Get("sequence_number"))
		assert.Equal(t, idsGwRuleUpdatedLogged, d.Get("logged"))
		assert.Equal(t, idsGwRuleUpdatedDisabled, d.Get("disabled"))
		assert.Equal(t, idsGwRuleUpdatedAction, d.Get("action"))
		assert.Equal(t, idsGwRuleUpdatedDisplayName, d.Get("display_name"))
		assert.Equal(t, idsGwRuleUpdatedDescription, d.Get("description"))
		assert.Equal(t, idsGwRuleUpdatedDirection, d.Get("direction"))
		assert.Equal(t, idsGwRuleUpdatedIPVersion, d.Get("ip_version"))
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIntrusionServiceGatewayPolicyRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsGwRuleData())
		setIdsGwRuleProfiles(d, t)
		// Don't set ID

		err := resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "ID")
	})

	t.Run("Update fails on patch error", func(t *testing.T) {
		mockSDK.EXPECT().Patch(idsGwRuleDomain, idsGwRulePolicyID, idsGwRuleID, gomock.Any()).Return(vapiErrors.NewInvalidRequest())

		res := resourceNsxtPolicyIntrusionServiceGatewayPolicyRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsGwRuleData())
		setIdsGwRuleProfiles(d, t)
		d.SetId(idsGwRuleID)

		err := resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Update fails when get after patch fails", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Patch(idsGwRuleDomain, idsGwRulePolicyID, idsGwRuleID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(idsGwRuleDomain, idsGwRulePolicyID, idsGwRuleID).Return(nsxModel.IdsRule{}, vapiErrors.NewInternalServerError()),
		)

		res := resourceNsxtPolicyIntrusionServiceGatewayPolicyRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsGwRuleData())
		setIdsGwRuleProfiles(d, t)
		d.SetId(idsGwRuleID)

		err := resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIntrusionServiceGatewayPolicyRuleDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIdsGwRuleMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(idsGwRuleDomain, idsGwRulePolicyID, idsGwRuleID).Return(nil)

		res := resourceNsxtPolicyIntrusionServiceGatewayPolicyRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsGwRuleData())
		setIdsGwRuleProfiles(d, t)
		d.SetId(idsGwRuleID)
		d.Set("nsx_id", idsGwRuleID)

		err := resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIntrusionServiceGatewayPolicyRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsGwRuleData())
		setIdsGwRuleProfiles(d, t)
		// Don't set ID

		err := resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleDelete(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "ID")
	})

	t.Run("Delete fails on API error", func(t *testing.T) {
		mockSDK.EXPECT().Delete(idsGwRuleDomain, idsGwRulePolicyID, idsGwRuleID).Return(vapiErrors.NewInternalServerError())

		res := resourceNsxtPolicyIntrusionServiceGatewayPolicyRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsGwRuleData())
		setIdsGwRuleProfiles(d, t)
		d.SetId(idsGwRuleID)

		err := resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
