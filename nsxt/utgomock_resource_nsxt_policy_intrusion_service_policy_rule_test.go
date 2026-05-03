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

	idsruleapi "github.com/vmware/terraform-provider-nsxt/api/infra/domains/intrusion_service_policies"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	idsrulemocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/domains/intrusion_service_policies"
)

var (
	idsRuleID               = "ids-rule-id"
	idsRuleDisplayName      = "Test IDS Policy Rule"
	idsRuleDescription      = "Test Intrusion Service Policy Rule"
	idsRuleDomain           = "default"
	idsRuleAction           = "DETECT"
	idsRuleDirection        = "IN_OUT"
	idsRuleIPVersion        = "IPV4_IPV6"
	idsRuleLogged           = false
	idsRuleDisabled         = false
	idsRuleSeqNum           = int64(0)
	idsRuleRuleID           = int64(1001)
	idsRuleScope            = []string{"/infra/tier-1s/test-t1-gw"}
	idsRuleProfilePath      = []string{"/infra/settings/firewall/security/intrusion-services/profiles/DefaultIDSProfile"}
	idsRulePolicyID         = "ids-policy-id"
	idsRulePolicyPath       = "/infra/domains/default/intrusion-service-policies/ids-policy-id"
	idsRuleOversubscription = "INHERIT_GLOBAL"
	// Updated values
	idsRuleUpdatedDirection        = "IN"
	idsRuleUpdatedIPVersion        = "IPV4"
	idsRuleUpdatedLogged           = true
	idsRuleUpdatedDisabled         = true
	idsRuleUpdatedSeqNum           = int64(1)
	idsRuleUpdatedAction           = "DETECT_PREVENT"
	idsRuleUpdatedOversubscription = "DROPPED"
	idsRuleUpdatedDisplayName      = "Test IDS Policy Rule Update"
	idsRuleUpdatedDescription      = "Test Intrusion Service Policy Rule Update"
)

func idsRuleAPIResponse() nsxModel.IdsRule {
	resourceType := "IdsRule"
	srcEx := false
	dstEx := false
	return nsxModel.IdsRule{
		Id:                   &idsRuleID,
		DisplayName:          &idsRuleDisplayName,
		Description:          &idsRuleDescription,
		ResourceType:         &resourceType,
		Action:               &idsRuleAction,
		Direction:            &idsRuleDirection,
		IpProtocol:           &idsRuleIPVersion,
		RuleId:               &idsRuleRuleID,
		Scope:                idsRuleScope,
		IdsProfiles:          idsRuleProfilePath,
		SequenceNumber:       &idsRuleSeqNum,
		Logged:               &idsRuleLogged,
		Disabled:             &idsRuleDisabled,
		SourcesExcluded:      &srcEx,
		DestinationsExcluded: &dstEx,
		Oversubscription:     &idsRuleOversubscription,
	}
}

func idsRuleAPIResponseUpdated() nsxModel.IdsRule {
	resourceType := "IdsRule"
	srcEx := false
	dstEx := false
	return nsxModel.IdsRule{
		Id:                   &idsRuleID,
		DisplayName:          &idsRuleUpdatedDisplayName,
		Description:          &idsRuleUpdatedDescription,
		ResourceType:         &resourceType,
		Action:               &idsRuleUpdatedAction,
		Direction:            &idsRuleUpdatedDirection,
		IpProtocol:           &idsRuleUpdatedIPVersion,
		RuleId:               &idsRuleRuleID,
		Scope:                idsRuleScope,
		IdsProfiles:          idsRuleProfilePath,
		SequenceNumber:       &idsRuleUpdatedSeqNum,
		Logged:               &idsRuleUpdatedLogged,
		Disabled:             &idsRuleUpdatedDisabled,
		SourcesExcluded:      &srcEx,
		DestinationsExcluded: &dstEx,
		Oversubscription:     &idsRuleUpdatedOversubscription,
	}
}

func minimalIdsRuleData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":     idsRuleDisplayName,
		"description":      idsRuleDescription,
		"policy_path":      idsRulePolicyPath,
		"action":           idsRuleAction,
		"direction":        idsRuleDirection,
		"ip_version":       idsRuleIPVersion,
		"sequence_number":  int(idsRuleSeqNum),
		"logged":           idsRuleLogged,
		"disabled":         idsRuleDisabled,
		"oversubscription": idsRuleOversubscription,
	}
}

// setIdsRuleProfiles sets the ids_profiles TypeSet on d (must be called after TestResourceDataRaw).
func setIdsRuleProfiles(d *schema.ResourceData, t *testing.T) {
	profilesSet := schema.NewSet(schema.HashString, []interface{}{
		"/infra/settings/firewall/security/intrusion-services/profiles/DefaultIDSProfile",
	})
	if err := d.Set("ids_profiles", profilesSet); err != nil {
		t.Fatal(err)
	}
}

// setIdsRuleScope sets the scope TypeSet on d (must be called after TestResourceDataRaw).
func setIdsRuleScope(d *schema.ResourceData, t *testing.T) {
	scopeSet := schema.NewSet(schema.HashString, []interface{}{
		"/infra/tier-1s/test-t1-gw",
	})
	if err := d.Set("scope", scopeSet); err != nil {
		t.Fatal(err)
	}
}

func setupIdsRuleMock(t *testing.T, ctrl *gomock.Controller) (*idsrulemocks.MockRulesClient, func()) {
	mockSDK := idsrulemocks.NewMockRulesClient(ctrl)
	mockWrapper := &idsruleapi.IntrusionServicePolicyRuleClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	original := cliIntrusionServicePolicyRulesClient
	cliIntrusionServicePolicyRulesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *idsruleapi.IntrusionServicePolicyRuleClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliIntrusionServicePolicyRulesClient = original }
}

func TestMockResourceNsxtPolicyIntrusionServicePolicyRuleCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIdsRuleMock(t, ctrl)
	defer restore()

	t.Run("Create success with auto-generated ID", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Patch(idsRuleDomain, idsRulePolicyID, gomock.Any(), gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(idsRuleDomain, idsRulePolicyID, gomock.Any()).Return(idsRuleAPIResponse(), nil),
		)

		res := resourceNsxtPolicyIntrusionServicePolicyRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsRuleData())
		setIdsRuleProfiles(d, t)
		setIdsRuleScope(d, t)

		err := resourceNsxtPolicyIntrusionServicePolicyRuleCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
	})

	t.Run("Create fails when patch fails", func(t *testing.T) {
		mockSDK.EXPECT().Patch(idsRuleDomain, idsRulePolicyID, gomock.Any(), gomock.Any()).Return(vapiErrors.NewInvalidRequest())

		res := resourceNsxtPolicyIntrusionServicePolicyRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsRuleData())
		setIdsRuleProfiles(d, t)
		setIdsRuleScope(d, t)

		err := resourceNsxtPolicyIntrusionServicePolicyRuleCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("Create fails when get after patch fails", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Patch(idsRuleDomain, idsRulePolicyID, gomock.Any(), gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(idsRuleDomain, idsRulePolicyID, gomock.Any()).Return(nsxModel.IdsRule{}, vapiErrors.NewInternalServerError()),
		)

		res := resourceNsxtPolicyIntrusionServicePolicyRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsRuleData())
		setIdsRuleProfiles(d, t)
		setIdsRuleScope(d, t)

		err := resourceNsxtPolicyIntrusionServicePolicyRuleCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIntrusionServicePolicyRuleRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIdsRuleMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(idsRuleDomain, idsRulePolicyID, idsRuleID).Return(idsRuleAPIResponse(), nil)

		res := resourceNsxtPolicyIntrusionServicePolicyRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsRuleData())
		setIdsRuleProfiles(d, t)
		setIdsRuleScope(d, t)
		d.SetId(idsRuleID)

		err := resourceNsxtPolicyIntrusionServicePolicyRuleRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, idsRuleDisplayName, d.Get("display_name"))
		assert.Equal(t, idsRuleDescription, d.Get("description"))
		assert.Equal(t, idsRuleAction, d.Get("action"))
		assert.Equal(t, idsRuleDirection, d.Get("direction"))
		assert.Equal(t, idsRuleIPVersion, d.Get("ip_version"))
		assert.Equal(t, int(idsRuleSeqNum), d.Get("sequence_number"))
		assert.Equal(t, idsRuleLogged, d.Get("logged"))
		assert.Equal(t, idsRuleDisabled, d.Get("disabled"))
		assert.Equal(t, idsRuleOversubscription, d.Get("oversubscription"))
	})

	t.Run("Read not found", func(t *testing.T) {
		mockSDK.EXPECT().Get(idsRuleDomain, idsRulePolicyID, idsRuleID).Return(nsxModel.IdsRule{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyIntrusionServicePolicyRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsRuleData())
		setIdsRuleProfiles(d, t)
		setIdsRuleScope(d, t)
		d.SetId(idsRuleID)

		err := resourceNsxtPolicyIntrusionServicePolicyRuleRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIntrusionServicePolicyRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsRuleData())
		setIdsRuleProfiles(d, t)
		setIdsRuleScope(d, t)

		err := resourceNsxtPolicyIntrusionServicePolicyRuleRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read fails on API error", func(t *testing.T) {
		mockSDK.EXPECT().Get(idsRuleDomain, idsRulePolicyID, idsRuleID).Return(nsxModel.IdsRule{}, vapiErrors.NewInternalServerError())

		res := resourceNsxtPolicyIntrusionServicePolicyRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsRuleData())
		setIdsRuleProfiles(d, t)
		setIdsRuleScope(d, t)
		d.SetId(idsRuleID)

		err := resourceNsxtPolicyIntrusionServicePolicyRuleRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIntrusionServicePolicyRuleUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIdsRuleMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Patch(idsRuleDomain, idsRulePolicyID, idsRuleID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(idsRuleDomain, idsRulePolicyID, idsRuleID).Return(idsRuleAPIResponseUpdated(), nil),
		)

		updatedData := minimalIdsRuleData()
		updatedData["sequence_number"] = int(idsRuleUpdatedSeqNum)
		updatedData["logged"] = idsRuleUpdatedLogged
		updatedData["disabled"] = idsRuleUpdatedDisabled
		updatedData["action"] = idsRuleUpdatedAction
		updatedData["display_name"] = idsRuleUpdatedDisplayName
		updatedData["description"] = idsRuleUpdatedDescription
		updatedData["direction"] = idsRuleUpdatedDirection
		updatedData["ip_version"] = idsRuleUpdatedIPVersion
		updatedData["oversubscription"] = idsRuleUpdatedOversubscription

		res := resourceNsxtPolicyIntrusionServicePolicyRule()
		d := schema.TestResourceDataRaw(t, res.Schema, updatedData)
		setIdsRuleProfiles(d, t)
		setIdsRuleScope(d, t)
		d.SetId(idsRuleID)

		err := resourceNsxtPolicyIntrusionServicePolicyRuleUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, int(idsRuleUpdatedSeqNum), d.Get("sequence_number"))
		assert.Equal(t, idsRuleUpdatedLogged, d.Get("logged"))
		assert.Equal(t, idsRuleUpdatedDisabled, d.Get("disabled"))
		assert.Equal(t, idsRuleUpdatedAction, d.Get("action"))
		assert.Equal(t, idsRuleUpdatedDisplayName, d.Get("display_name"))
		assert.Equal(t, idsRuleUpdatedDescription, d.Get("description"))
		assert.Equal(t, idsRuleUpdatedDirection, d.Get("direction"))
		assert.Equal(t, idsRuleUpdatedIPVersion, d.Get("ip_version"))
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIntrusionServicePolicyRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsRuleData())
		setIdsRuleProfiles(d, t)
		setIdsRuleScope(d, t)
		// Don't set ID

		err := resourceNsxtPolicyIntrusionServicePolicyRuleUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Update fails on patch error", func(t *testing.T) {
		mockSDK.EXPECT().Patch(idsRuleDomain, idsRulePolicyID, idsRuleID, gomock.Any()).Return(vapiErrors.NewInternalServerError())

		res := resourceNsxtPolicyIntrusionServicePolicyRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsRuleData())
		setIdsRuleProfiles(d, t)
		setIdsRuleScope(d, t)
		d.SetId(idsRuleID)

		err := resourceNsxtPolicyIntrusionServicePolicyRuleUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIntrusionServicePolicyRuleDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIdsRuleMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(idsRuleDomain, idsRulePolicyID, idsRuleID).Return(nil)

		res := resourceNsxtPolicyIntrusionServicePolicyRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsRuleData())
		setIdsRuleProfiles(d, t)
		setIdsRuleScope(d, t)
		d.SetId(idsRuleID)

		err := resourceNsxtPolicyIntrusionServicePolicyRuleDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIntrusionServicePolicyRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsRuleData())
		setIdsRuleProfiles(d, t)
		setIdsRuleScope(d, t)
		// Don't set ID

		err := resourceNsxtPolicyIntrusionServicePolicyRuleDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Delete fails on API error", func(t *testing.T) {
		mockSDK.EXPECT().Delete(idsRuleDomain, idsRulePolicyID, idsRuleID).Return(vapiErrors.NewInternalServerError())

		res := resourceNsxtPolicyIntrusionServicePolicyRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsRuleData())
		setIdsRuleProfiles(d, t)
		setIdsRuleScope(d, t)
		d.SetId(idsRuleID)

		err := resourceNsxtPolicyIntrusionServicePolicyRuleDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
