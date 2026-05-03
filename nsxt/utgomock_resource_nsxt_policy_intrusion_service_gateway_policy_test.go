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

	apipkg "github.com/vmware/terraform-provider-nsxt/api"
	apidomains "github.com/vmware/terraform-provider-nsxt/api/infra/domains"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	inframocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	domainmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/domains"
)

var (
	idsGwPolicyID          = "ids-gw-policy-id"
	idsGwPolicyDisplayName = "Test IDS Gateway Policy"
	idsGwPolicyDescription = "Test Intrusion Service Gateway Policy"
	idsGwPolicyDomain      = "default"
	idsGwPolicyPath        = "/infra/domains/default/intrusion-service-gateway-policies/ids-gw-policy-id"
	idsGwPolicyCategory    = "LocalGatewayRules"
	idsGwPolicySeqNum      = int64(0)
	idsGwPolicyStateful    = true
	idsGwPolicyLocked      = false
	idsGwPolicyComments    = ""

	// Updated values
	idsGwPolicyUpdatedSeqNum      = int64(1)
	idsGwPolicyUpdatedLocked      = true
	idsGwPolicyUpdatedComments    = "Updated Policy"
	idsGwPolicyUpdatedDisplayName = "Test IDS Gateway Policy Update"
	idsGwPolicyUpdatedDescription = "Test Intrusion Service Gateway Policy Update"
)

var (
	// Embedded rule variables
	embeddedRuleID           = "ids-gw-policy-rule-id"
	embeddedRuleDisplayName  = "Test IDS Gateway Policy Rule"
	embeddedRuleDescription  = "Test Intrusion Service Gateway Policy Rule"
	embeddedRuleResourceType = "IdsRule"
	embeddedRuleAction       = "DETECT"
	embeddedRuleDirection    = "IN_OUT"
	embeddedRuleIPVersion    = "IPV4_IPV6"
	embeddedRuleRuleID       = int64(1001)
	embeddedRuleSeqNum       = int64(1)
	embeddedRuleLogged       = false
	embeddedRuleDisabled     = false
	embeddedRuleSrcEx        = false
	embeddedRuleDstEx        = false
	embeddedRuleProfiles     = []string{"/infra/settings/firewall/security/intrusion-services/profiles/DefaultIDSProfile"}

	// Updated embedded rule variables
	embeddedRuleUpdatedDisplayName = "Test IDS Gateway Policy Rule Update"
	embeddedRuleUpdatedDescription = "Test Intrusion Service Gateway Policy Rule Update"
	embeddedRuleUpdatedAction      = "DETECT_PREVENT"
	embeddedRuleUpdatedDirection   = "IN"
	embeddedRuleUpdatedIPVersion   = "IPV4"
	embeddedRuleUpdatedSeqNum      = int64(2)
	embeddedRuleUpdatedLogged      = true
	embeddedRuleUpdatedDisabled    = true
	embeddedRuleUpdatedSrcEx       = true
	embeddedRuleUpdatedDstEx       = true
)

func createTestIdsRule() nsxModel.IdsRule {
	return nsxModel.IdsRule{
		Id:                   &embeddedRuleID,
		DisplayName:          &embeddedRuleDisplayName,
		Description:          &embeddedRuleDescription,
		ResourceType:         &embeddedRuleResourceType,
		Action:               &embeddedRuleAction,
		Direction:            &embeddedRuleDirection,
		IpProtocol:           &embeddedRuleIPVersion,
		RuleId:               &embeddedRuleRuleID,
		IdsProfiles:          embeddedRuleProfiles,
		SequenceNumber:       &embeddedRuleSeqNum,
		Logged:               &embeddedRuleLogged,
		Disabled:             &embeddedRuleDisabled,
		SourcesExcluded:      &embeddedRuleSrcEx,
		DestinationsExcluded: &embeddedRuleDstEx,
	}
}

func idsGwPolicyAPIResponse() nsxModel.IdsGatewayPolicy {
	resourceType := "IdsGatewayPolicy"
	rules := []nsxModel.IdsRule{createTestIdsRule()}
	return nsxModel.IdsGatewayPolicy{
		Id:             &idsGwPolicyID,
		DisplayName:    &idsGwPolicyDisplayName,
		Description:    &idsGwPolicyDescription,
		Path:           &idsGwPolicyPath,
		Category:       &idsGwPolicyCategory,
		ResourceType:   &resourceType,
		SequenceNumber: &idsGwPolicySeqNum,
		Stateful:       &idsGwPolicyStateful,
		Locked:         &idsGwPolicyLocked,
		Comments:       &idsGwPolicyComments,
		Rules:          rules,
	}
}

func createTestIdsRuleUpdate() nsxModel.IdsRule {
	return nsxModel.IdsRule{
		Id:                   &embeddedRuleID,
		DisplayName:          &embeddedRuleUpdatedDisplayName,
		Description:          &embeddedRuleUpdatedDescription,
		ResourceType:         &embeddedRuleResourceType,
		Action:               &embeddedRuleUpdatedAction,
		Direction:            &embeddedRuleUpdatedDirection,
		IpProtocol:           &embeddedRuleUpdatedIPVersion,
		RuleId:               &embeddedRuleRuleID,
		IdsProfiles:          embeddedRuleProfiles,
		SequenceNumber:       &embeddedRuleUpdatedSeqNum,
		Logged:               &embeddedRuleUpdatedLogged,
		Disabled:             &embeddedRuleUpdatedDisabled,
		SourcesExcluded:      &embeddedRuleUpdatedSrcEx,
		DestinationsExcluded: &embeddedRuleUpdatedDstEx,
	}
}

func idsGwPolicyAPIResponseUpdated() nsxModel.IdsGatewayPolicy {
	resourceType := "IdsGatewayPolicy"
	rules := []nsxModel.IdsRule{createTestIdsRuleUpdate()}
	return nsxModel.IdsGatewayPolicy{
		Id:             &idsGwPolicyID,
		DisplayName:    &idsGwPolicyUpdatedDisplayName,
		Description:    &idsGwPolicyUpdatedDescription,
		Path:           &idsGwPolicyPath,
		Category:       &idsGwPolicyCategory,
		ResourceType:   &resourceType,
		SequenceNumber: &idsGwPolicyUpdatedSeqNum,
		Stateful:       &idsGwPolicyStateful,
		Locked:         &idsGwPolicyUpdatedLocked,
		Comments:       &idsGwPolicyUpdatedComments,
		Rules:          rules,
	}
}

func minimalEmbeddedIdsGwRuleData() map[string]interface{} {
	// Convert []string to []interface{}
	profiles := make([]interface{}, len(embeddedRuleProfiles))
	for i, profile := range embeddedRuleProfiles {
		profiles[i] = profile
	}
	return map[string]interface{}{
		"display_name":    embeddedRuleDisplayName,
		"description":     embeddedRuleDescription,
		"action":          embeddedRuleAction,
		"direction":       embeddedRuleDirection,
		"ip_version":      embeddedRuleIPVersion,
		"sequence_number": int(embeddedRuleSeqNum),
		"logged":          embeddedRuleLogged,
		"disabled":        embeddedRuleDisabled,
		"ids_profiles":    profiles,
	}
}

func minimalIdsGwPolicyData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":    idsGwPolicyDisplayName,
		"description":     idsGwPolicyDescription,
		"nsx_id":          idsGwPolicyID,
		"domain":          idsGwPolicyDomain,
		"category":        idsGwPolicyCategory,
		"sequence_number": 0,
		"stateful":        true,
		"locked":          false,
		"rule": []interface{}{
			minimalEmbeddedIdsGwRuleData(),
		},
	}
}

func minimalIdsGwPolicyDataUpdate() map[string]interface{} {
	return map[string]interface{}{
		"display_name":    idsGwPolicyDisplayName,
		"description":     idsGwPolicyDescription,
		"nsx_id":          idsGwPolicyID,
		"domain":          idsGwPolicyDomain,
		"sequence_number": 0,
		"locked":          false,
		"rule": []interface{}{
			minimalEmbeddedIdsGwRuleData(),
		},
	}
}

func setupIdsGwPolicyMock(t *testing.T, ctrl *gomock.Controller) (*domainmocks.MockIntrusionServiceGatewayPoliciesClient, *inframocks.MockInfraClient, func()) {
	t.Helper()
	mockSDK := domainmocks.NewMockIntrusionServiceGatewayPoliciesClient(ctrl)
	mockWrapper := &apidomains.IntrusionServiceGatewayPolicyClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	originalGWP := cliIntrusionServiceGatewayPoliciesClient
	cliIntrusionServiceGatewayPoliciesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apidomains.IntrusionServiceGatewayPolicyClientContext {
		return mockWrapper
	}

	mockInfraSDK := inframocks.NewMockInfraClient(ctrl)
	originalInfra := cliInfraClient
	cliInfraClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apipkg.InfraClientContext {
		return &apipkg.InfraClientContext{Client: mockInfraSDK, ClientType: utl.Local}
	}

	return mockSDK, mockInfraSDK, func() {
		cliIntrusionServiceGatewayPoliciesClient = originalGWP
		cliInfraClient = originalInfra
	}
}

func TestMockResourceNsxtPolicyIntrusionServiceGatewayPolicyCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, mockInfra, restore := setupIdsGwPolicyMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(idsGwPolicyDomain, idsGwPolicyID).Return(nsxModel.IdsGatewayPolicy{}, notFoundErr),
			mockInfra.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(idsGwPolicyDomain, idsGwPolicyID).Return(idsGwPolicyAPIResponse(), nil),
		)

		res := resourceNsxtPolicyIntrusionServiceGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsGwPolicyData())

		err := resourceNsxtPolicyIntrusionServiceGatewayPolicyCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, idsGwPolicyID, d.Id())
		assert.Equal(t, idsGwPolicyDisplayName, d.Get("display_name"))
	})

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(idsGwPolicyDomain, idsGwPolicyID).Return(idsGwPolicyAPIResponse(), nil)

		res := resourceNsxtPolicyIntrusionServiceGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsGwPolicyData())

		err := resourceNsxtPolicyIntrusionServiceGatewayPolicyCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicyIntrusionServiceGatewayPolicyRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, _, restore := setupIdsGwPolicyMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(idsGwPolicyDomain, idsGwPolicyID).Return(idsGwPolicyAPIResponse(), nil)

		res := resourceNsxtPolicyIntrusionServiceGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsGwPolicyData())
		d.SetId(idsGwPolicyID)

		err := resourceNsxtPolicyIntrusionServiceGatewayPolicyRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, idsGwPolicyDomain, d.Get("domain"))
		assert.Equal(t, idsGwPolicyCategory, d.Get("category"))
		assert.Equal(t, idsGwPolicyDisplayName, d.Get("display_name"))
		assert.Equal(t, idsGwPolicyDescription, d.Get("description"))
		assert.Equal(t, int(idsGwPolicySeqNum), d.Get("sequence_number"))
		assert.Equal(t, idsGwPolicyStateful, d.Get("stateful"))
		assert.Equal(t, idsGwPolicyLocked, d.Get("locked"))
		// Verify rule count matches actual embedded rules
		// rule_count removed for consistency with other policies
		rules := d.Get("rule").([]interface{})
		// rule_count removed for consistency with other policies

		// Verify the embedded rule
		assert.Equal(t, 1, len(rules))
		rule := rules[0].(map[string]interface{})
		assert.Equal(t, embeddedRuleDisplayName, rule["display_name"])
		assert.Equal(t, embeddedRuleDescription, rule["description"])
		assert.Equal(t, embeddedRuleAction, rule["action"])
		assert.Equal(t, int(embeddedRuleSeqNum), rule["sequence_number"])
		assert.Equal(t, embeddedRuleLogged, rule["logged"])
		assert.Equal(t, embeddedRuleDisabled, rule["disabled"])
		assert.Equal(t, embeddedRuleDirection, rule["direction"])
		assert.Equal(t, embeddedRuleIPVersion, rule["ip_version"])
		// Convert Set to slice for comparison
		idsProfilesSet := rule["ids_profiles"].(*schema.Set)
		idsProfilesList := make([]string, idsProfilesSet.Len())
		for i, v := range idsProfilesSet.List() {
			idsProfilesList[i] = v.(string)
		}
		assert.Equal(t, embeddedRuleProfiles, idsProfilesList)
	})

	t.Run("Read not found", func(t *testing.T) {
		mockSDK.EXPECT().Get(idsGwPolicyDomain, idsGwPolicyID).Return(nsxModel.IdsGatewayPolicy{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyIntrusionServiceGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsGwPolicyData())
		d.SetId(idsGwPolicyID)

		err := resourceNsxtPolicyIntrusionServiceGatewayPolicyRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIntrusionServiceGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsGwPolicyData())

		err := resourceNsxtPolicyIntrusionServiceGatewayPolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIntrusionServiceGatewayPolicyUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, mockInfra, restore := setupIdsGwPolicyMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockInfra.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(idsGwPolicyDomain, idsGwPolicyID).Return(idsGwPolicyAPIResponseUpdated(), nil),
		)

		updateData := minimalIdsGwPolicyDataUpdate()
		updateData["sequence_number"] = int(idsGwPolicyUpdatedSeqNum)
		updateData["locked"] = idsGwPolicyUpdatedLocked
		updateData["comments"] = idsGwPolicyUpdatedComments
		updateData["display_name"] = idsGwPolicyUpdatedDisplayName
		updateData["description"] = idsGwPolicyUpdatedDescription

		// Update the embedded rule
		updateRule := minimalEmbeddedIdsGwRuleData()
		updateRule["display_name"] = "Test IDS Gateway Policy Rule Updated"
		updateRule["description"] = "Test Intrusion Service Gateway Policy Rule Updated"
		updateRule["action"] = "DETECT_PREVENT"
		updateRule["sequence_number"] = 2
		updateRule["logged"] = true
		updateRule["disabled"] = true
		updateRule["direction"] = "IN"
		updateRule["ip_version"] = "IPV4"
		updateData["rule"] = []interface{}{updateRule}

		res := resourceNsxtPolicyIntrusionServiceGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, updateData)
		d.SetId(idsGwPolicyID)

		err := resourceNsxtPolicyIntrusionServiceGatewayPolicyUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, int(idsGwPolicyUpdatedSeqNum), d.Get("sequence_number"))
		assert.Equal(t, idsGwPolicyStateful, d.Get("stateful"))
		assert.Equal(t, idsGwPolicyUpdatedLocked, d.Get("locked"))
		assert.Equal(t, idsGwPolicyUpdatedComments, d.Get("comments"))
		assert.Equal(t, idsGwPolicyCategory, d.Get("category"))
		assert.Equal(t, idsGwPolicyUpdatedDisplayName, d.Get("display_name"))
		assert.Equal(t, idsGwPolicyUpdatedDescription, d.Get("description"))
		// rule_count removed for consistency with other policies

		// Verify the update embedded rule
		rules := d.Get("rule").([]interface{})
		assert.Equal(t, 1, len(rules))
		rule := rules[0].(map[string]interface{})
		assert.Equal(t, embeddedRuleUpdatedDisplayName, rule["display_name"])
		assert.Equal(t, embeddedRuleUpdatedDescription, rule["description"])
		assert.Equal(t, embeddedRuleUpdatedAction, rule["action"])
		assert.Equal(t, int(embeddedRuleUpdatedSeqNum), rule["sequence_number"])
		assert.Equal(t, embeddedRuleUpdatedLogged, rule["logged"])
		assert.Equal(t, embeddedRuleUpdatedDirection, rule["direction"])
		assert.Equal(t, embeddedRuleUpdatedIPVersion, rule["ip_version"])
		// tcp_strict not available for IDS rules
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIntrusionServiceGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsGwPolicyData())

		err := resourceNsxtPolicyIntrusionServiceGatewayPolicyUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIntrusionServiceGatewayPolicyDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, _, restore := setupIdsGwPolicyMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(idsGwPolicyDomain, idsGwPolicyID).Return(nil)

		res := resourceNsxtPolicyIntrusionServiceGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsGwPolicyData())
		d.SetId(idsGwPolicyID)

		err := resourceNsxtPolicyIntrusionServiceGatewayPolicyDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIntrusionServiceGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsGwPolicyData())

		err := resourceNsxtPolicyIntrusionServiceGatewayPolicyDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
