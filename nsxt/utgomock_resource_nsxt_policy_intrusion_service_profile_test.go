//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	intrusionservices "github.com/vmware/terraform-provider-nsxt/api/infra/settings/firewall/security/intrusion_services"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	isprofmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/settings/firewall/security/intrusion_services"
)

var (
	isprofID          = "ids-profile-1"
	isprofDisplayName = "Test IDS Profile"
	isprofDescription = "test intrusion service profile"
	isprofRevision    = int64(1)
	isprofPath        = "/infra/settings/firewall/security/intrusion-services/profiles/ids-profile-1"
	isprofSeverity    = nsxModel.IdsProfile_PROFILE_SEVERITY_HIGH
)

func isprofAPIResponse() nsxModel.IdsProfile {
	return nsxModel.IdsProfile{
		Id:              &isprofID,
		DisplayName:     &isprofDisplayName,
		Description:     &isprofDescription,
		Revision:        &isprofRevision,
		Path:            &isprofPath,
		ProfileSeverity: []string{isprofSeverity},
	}
}

func minimalIsprofData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": isprofDisplayName,
		"description":  isprofDescription,
		"nsx_id":       isprofID,
		"severities":   []interface{}{isprofSeverity},
	}
}

func setupIsprofMock(t *testing.T, ctrl *gomock.Controller) (*isprofmocks.MockProfilesClient, func()) {
	mockSDK := isprofmocks.NewMockProfilesClient(ctrl)
	mockWrapper := &intrusionservices.IdsProfileClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	original := cliIntrusionServiceProfilesClient
	cliIntrusionServiceProfilesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *intrusionservices.IdsProfileClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliIntrusionServiceProfilesClient = original }
}

func TestMockResourceNsxtPolicyIntrusionServiceProfileCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIsprofMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(isprofID).Return(nsxModel.IdsProfile{}, notFoundErr),
			mockSDK.EXPECT().Patch(isprofID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(isprofID).Return(isprofAPIResponse(), nil),
		)

		res := resourceNsxtPolicyIntrusionServiceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIsprofData())

		err := resourceNsxtPolicyIntrusionServiceProfileCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, isprofID, d.Id())
		assert.Equal(t, isprofDisplayName, d.Get("display_name"))
	})

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(isprofID).Return(isprofAPIResponse(), nil)

		res := resourceNsxtPolicyIntrusionServiceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIsprofData())

		err := resourceNsxtPolicyIntrusionServiceProfileCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIntrusionServiceProfileRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIsprofMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(isprofID).Return(isprofAPIResponse(), nil)

		res := resourceNsxtPolicyIntrusionServiceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIsprofData())
		d.SetId(isprofID)

		err := resourceNsxtPolicyIntrusionServiceProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, isprofDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(isprofID).Return(nsxModel.IdsProfile{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyIntrusionServiceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIsprofData())
		d.SetId(isprofID)

		err := resourceNsxtPolicyIntrusionServiceProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIntrusionServiceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIsprofData())

		err := resourceNsxtPolicyIntrusionServiceProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIntrusionServiceProfileUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIsprofMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Patch(isprofID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(isprofID).Return(isprofAPIResponse(), nil),
		)

		res := resourceNsxtPolicyIntrusionServiceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIsprofData())
		d.SetId(isprofID)

		err := resourceNsxtPolicyIntrusionServiceProfileUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIntrusionServiceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIsprofData())

		err := resourceNsxtPolicyIntrusionServiceProfileUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIntrusionServiceProfileDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIsprofMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(isprofID).Return(nil)

		res := resourceNsxtPolicyIntrusionServiceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIsprofData())
		d.SetId(isprofID)

		err := resourceNsxtPolicyIntrusionServiceProfileDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIntrusionServiceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIsprofData())

		err := resourceNsxtPolicyIntrusionServiceProfileDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestUnitNsxt_buildIdsProfileCriteriaFilter(t *testing.T) {
	sv, err := buildIdsProfileCriteriaFilter(nsxModel.IdsProfileFilterCriteria_FILTER_NAME_ATTACK_TYPE, []string{"a", "b"})
	require.NoError(t, err)
	require.NotNil(t, sv)

	converter := bindings.NewTypeConverter()
	golang, errs := converter.ConvertToGolang(sv, nsxModel.IdsProfileFilterCriteriaBindingType())
	require.Empty(t, errs)
	criteria := golang.(nsxModel.IdsProfileFilterCriteria)
	assert.Equal(t, nsxModel.IdsProfileFilterCriteria_FILTER_NAME_ATTACK_TYPE, *criteria.FilterName)
	assert.Equal(t, []string{"a", "b"}, criteria.FilterValue)
}

func TestUnitNsxt_buildIdsProfileCriteriaOperator(t *testing.T) {
	sv, err := buildIdsProfileCriteriaOperator()
	require.NoError(t, err)
	require.NotNil(t, sv)

	converter := bindings.NewTypeConverter()
	golang, errs := converter.ConvertToGolang(sv, nsxModel.IdsProfileConjunctionOperatorBindingType())
	require.Empty(t, errs)
	operator := golang.(nsxModel.IdsProfileConjunctionOperator)
	assert.Equal(t, "AND", *operator.Operator)
}

func idsProfileTestResource() *schema.Resource {
	return resourceNsxtPolicyIntrusionServiceProfile()
}

func TestUnitNsxt_getIdsProfileCriteriaFromSchema(t *testing.T) {
	res := idsProfileTestResource()
	converter := bindings.NewTypeConverter()

	t.Run("no criteria returns empty list", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		result, err := getIdsProfileCriteriaFromSchema(d)
		require.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("single criteria field produces one filter, no trailing operator", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"criteria": []interface{}{
				map[string]interface{}{
					"attack_types": []interface{}{"malware"},
				},
			},
		})
		result, err := getIdsProfileCriteriaFromSchema(d)
		require.NoError(t, err)
		require.Len(t, result, 1)

		golang, errs := converter.ConvertToGolang(result[0], nsxModel.IdsProfileFilterCriteriaBindingType())
		require.Empty(t, errs)
		criteria := golang.(nsxModel.IdsProfileFilterCriteria)
		assert.Equal(t, nsxModel.IdsProfileFilterCriteria_FILTER_NAME_ATTACK_TYPE, *criteria.FilterName)
	})

	t.Run("multiple criteria fields are joined with AND operators", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"criteria": []interface{}{
				map[string]interface{}{
					"attack_types":   []interface{}{"malware"},
					"attack_targets": []interface{}{"server"},
				},
			},
		})
		result, err := getIdsProfileCriteriaFromSchema(d)
		require.NoError(t, err)
		// filter, AND, filter -> the trailing operator after the last filter is stripped
		require.Len(t, result, 3)

		golang0, errs := converter.ConvertToGolang(result[0], nsxModel.IdsProfileFilterCriteriaBindingType())
		require.Empty(t, errs)
		assert.Equal(t, nsxModel.IdsProfileFilterCriteria_FILTER_NAME_ATTACK_TYPE, *golang0.(nsxModel.IdsProfileFilterCriteria).FilterName)

		golangOp, errs := converter.ConvertToGolang(result[1], nsxModel.IdsProfileConjunctionOperatorBindingType())
		require.Empty(t, errs)
		assert.Equal(t, "AND", *golangOp.(nsxModel.IdsProfileConjunctionOperator).Operator)

		golang2, errs := converter.ConvertToGolang(result[2], nsxModel.IdsProfileFilterCriteriaBindingType())
		require.Empty(t, errs)
		assert.Equal(t, nsxModel.IdsProfileFilterCriteria_FILTER_NAME_ATTACK_TARGET, *golang2.(nsxModel.IdsProfileFilterCriteria).FilterName)
	})
}

func TestUnitNsxt_setIdsProfileCriteriaInSchema(t *testing.T) {
	res := idsProfileTestResource()

	t.Run("filter criteria without operators populate the schema", func(t *testing.T) {
		d := res.TestResourceData()
		filter, err := buildIdsProfileCriteriaFilter(nsxModel.IdsProfileFilterCriteria_FILTER_NAME_CVSS, []string{"HIGH"})
		require.NoError(t, err)

		require.NoError(t, setIdsProfileCriteriaInSchema([]*data.StructValue{filter}, d))

		criteria := d.Get("criteria").([]interface{})
		require.Len(t, criteria, 1)
		elem := criteria[0].(map[string]interface{})
		assert.ElementsMatch(t, []string{"HIGH"}, elem["cvss"].(*schema.Set).List())
	})

	t.Run("AND operators at odd indices are skipped", func(t *testing.T) {
		d := res.TestResourceData()
		filter1, err := buildIdsProfileCriteriaFilter(nsxModel.IdsProfileFilterCriteria_FILTER_NAME_ATTACK_TYPE, []string{"malware"})
		require.NoError(t, err)
		op, err := buildIdsProfileCriteriaOperator()
		require.NoError(t, err)
		filter2, err := buildIdsProfileCriteriaFilter(nsxModel.IdsProfileFilterCriteria_FILTER_NAME_ATTACK_TARGET, []string{"server"})
		require.NoError(t, err)

		require.NoError(t, setIdsProfileCriteriaInSchema([]*data.StructValue{filter1, op, filter2}, d))

		criteria := d.Get("criteria").([]interface{})
		require.Len(t, criteria, 1)
		elem := criteria[0].(map[string]interface{})
		assert.ElementsMatch(t, []string{"malware"}, elem["attack_types"].(*schema.Set).List())
		assert.ElementsMatch(t, []string{"server"}, elem["attack_targets"].(*schema.Set).List())
	})

	t.Run("empty list clears the schema", func(t *testing.T) {
		d := res.TestResourceData()
		require.NoError(t, setIdsProfileCriteriaInSchema(nil, d))
		assert.Empty(t, d.Get("criteria").([]interface{}))
	})
}

func TestUnitNsxt_getIdsProfileSignaturesFromSchema(t *testing.T) {
	res := idsProfileTestResource()

	t.Run("no signatures returns empty list", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		result := getIdsProfileSignaturesFromSchema(d)
		assert.Empty(t, result)
	})

	t.Run("converts overridden signatures", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"overridden_signature": []interface{}{
				map[string]interface{}{
					"signature_id": "sig-1",
					"enabled":      true,
					"action":       nsxModel.IdsProfileLocalSignature_ACTION_DROP,
				},
			},
		})
		result := getIdsProfileSignaturesFromSchema(d)
		require.Len(t, result, 1)
		assert.Equal(t, "sig-1", *result[0].SignatureId)
		assert.True(t, *result[0].Enable)
		assert.Equal(t, nsxModel.IdsProfileLocalSignature_ACTION_DROP, *result[0].Action)
	})
}

func TestUnitNsxt_setIdsProfileSignaturesInSchema(t *testing.T) {
	res := idsProfileTestResource()

	t.Run("populates schema from signature models", func(t *testing.T) {
		d := res.TestResourceData()
		require.NoError(t, setIdsProfileSignaturesInSchema([]nsxModel.IdsProfileLocalSignature{
			{SignatureId: strPtr("sig-1"), Enable: boolPtr(true), Action: strPtr(nsxModel.IdsProfileLocalSignature_ACTION_REJECT)},
		}, d))

		signatures := d.Get("overridden_signature").(*schema.Set).List()
		require.Len(t, signatures, 1)
		elem := signatures[0].(map[string]interface{})
		assert.Equal(t, "sig-1", elem["signature_id"])
		assert.Equal(t, true, elem["enabled"])
		assert.Equal(t, nsxModel.IdsProfileLocalSignature_ACTION_REJECT, elem["action"])
	})

	t.Run("empty list clears the schema", func(t *testing.T) {
		d := res.TestResourceData()
		require.NoError(t, setIdsProfileSignaturesInSchema(nil, d))
		assert.Empty(t, d.Get("overridden_signature").(*schema.Set).List())
	})
}
