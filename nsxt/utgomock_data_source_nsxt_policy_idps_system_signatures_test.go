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
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/settings/firewall/security/intrusion_services/signature_versions"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	signaturesmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/settings/firewall/security/intrusion_services/signature_versions"
)

// dataSourceNsxtPolicyIdpsSystemSignaturesRead reads through the cliIdsSignatureVersionsClient
// and cliIdsSignaturesClient package-level vars (see
// utgomock_resource_nsxt_policy_idps_signature_version_test.go for the version-client mock
// helper), so both the guard condition and the mocked lookups are covered here.

func setupIdpsSignaturesMock(ctrl *gomock.Controller) (*signaturesmocks.MockSignaturesClient, func()) {
	mockClient := signaturesmocks.NewMockSignaturesClient(ctrl)
	orig := cliIdsSignaturesClient
	cliIdsSignaturesClient = func(_ client.Connector) signature_versions.SignaturesClient {
		return mockClient
	}
	return mockClient, func() { cliIdsSignaturesClient = orig }
}

func TestMockDataSourceNsxtPolicyIdpsSystemSignaturesReadGuard(t *testing.T) {
	t.Run("Read fails for global manager", func(t *testing.T) {
		ds := dataSourceNsxtPolicyIdpsSystemSignatures()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		err := dataSourceNsxtPolicyIdpsSystemSignaturesRead(d, newGoMockGlobalProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Global Manager")
	})
}

func TestMockDataSourceNsxtPolicyIdpsSystemSignaturesReadExplicitVersion(t *testing.T) {
	t.Run("skips the version lookup when version_id is set and applies filters", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSigs, restoreSigs := setupIdpsSignaturesMock(ctrl)
		defer restoreSigs()

		matchID, otherID := "sig-1", "sig-2"
		matchName, otherName := "Match Name", "Other Name"
		high, low := "HIGH", "LOW"
		count := int64(2)
		mockSigs.EXPECT().List("v1", (*string)(nil), nil, nil, nil, nil, nil, nil).Return(
			model.IdsSignatureListResult{
				Results: []model.IdsSignature{
					{Id: &matchID, DisplayName: &matchName, Severity: &high},
					{Id: &otherID, DisplayName: &otherName, Severity: &low},
				},
				ResultCount: &count,
			}, nil,
		)

		ds := dataSourceNsxtPolicyIdpsSystemSignatures()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"version_id":   "v1",
			"severity":     "HIGH",
			"display_name": "match",
		})

		err := dataSourceNsxtPolicyIdpsSystemSignaturesRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		sigs := d.Get("signatures").([]interface{})
		require.Len(t, sigs, 1)
		assert.Equal(t, "sig-1", sigs[0].(map[string]interface{})["id"])
	})

	t.Run("signatures List error propagates", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSigs, restoreSigs := setupIdpsSignaturesMock(ctrl)
		defer restoreSigs()
		mockSigs.EXPECT().List("v1", (*string)(nil), nil, nil, nil, nil, nil, nil).Return(
			model.IdsSignatureListResult{}, vapiErrors.InternalServerError{},
		)

		ds := dataSourceNsxtPolicyIdpsSystemSignatures()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"version_id": "v1"})

		err := dataSourceNsxtPolicyIdpsSystemSignaturesRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("paginates across cursors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSigs, restoreSigs := setupIdpsSignaturesMock(ctrl)
		defer restoreSigs()

		id1, id2 := "sig-1", "sig-2"
		count := int64(2)
		cursor := "next-page"
		mockSigs.EXPECT().List("v1", (*string)(nil), nil, nil, nil, nil, nil, nil).Return(
			model.IdsSignatureListResult{Results: []model.IdsSignature{{Id: &id1}}, ResultCount: &count, Cursor: &cursor}, nil,
		)
		mockSigs.EXPECT().List("v1", &cursor, nil, nil, nil, nil, nil, nil).Return(
			model.IdsSignatureListResult{Results: []model.IdsSignature{{Id: &id2}}, ResultCount: &count}, nil,
		)

		ds := dataSourceNsxtPolicyIdpsSystemSignatures()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"version_id": "v1"})

		err := dataSourceNsxtPolicyIdpsSystemSignaturesRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Len(t, d.Get("signatures").([]interface{}), 2)
	})
}

func TestMockDataSourceNsxtPolicyIdpsSystemSignaturesReadVersionLookup(t *testing.T) {
	t.Run("uses the ACTIVE version when version_id is unset", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockVersions, restoreVersions := setupIdpsSignatureVersionMock(ctrl)
		defer restoreVersions()
		mockSigs, restoreSigs := setupIdpsSignaturesMock(ctrl)
		defer restoreSigs()

		activeID, notActiveID := "active-ver", "old-ver"
		activeState, notActiveState := "ACTIVE", "NOTACTIVE"
		mockVersions.EXPECT().List(nil, nil, nil, nil, nil, nil, nil).Return(
			model.IdsSignatureVersionListResult{Results: []model.IdsSignatureVersion{
				{Id: &notActiveID, State: &notActiveState},
				{Id: &activeID, State: &activeState},
			}}, nil,
		)
		count := int64(0)
		mockSigs.EXPECT().List("active-ver", (*string)(nil), nil, nil, nil, nil, nil, nil).Return(
			model.IdsSignatureListResult{ResultCount: &count}, nil,
		)

		ds := dataSourceNsxtPolicyIdpsSystemSignatures()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		err := dataSourceNsxtPolicyIdpsSystemSignaturesRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "active-ver", d.Get("version_id"))
	})

	t.Run("falls back to the first version when none is ACTIVE", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockVersions, restoreVersions := setupIdpsSignatureVersionMock(ctrl)
		defer restoreVersions()
		mockSigs, restoreSigs := setupIdpsSignaturesMock(ctrl)
		defer restoreSigs()

		firstID := "first-ver"
		notActiveState := "NOTACTIVE"
		mockVersions.EXPECT().List(nil, nil, nil, nil, nil, nil, nil).Return(
			model.IdsSignatureVersionListResult{Results: []model.IdsSignatureVersion{{Id: &firstID, State: &notActiveState}}}, nil,
		)
		count := int64(0)
		mockSigs.EXPECT().List("first-ver", (*string)(nil), nil, nil, nil, nil, nil, nil).Return(
			model.IdsSignatureListResult{ResultCount: &count}, nil,
		)

		ds := dataSourceNsxtPolicyIdpsSystemSignatures()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		err := dataSourceNsxtPolicyIdpsSystemSignaturesRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "first-ver", d.Get("version_id"))
	})

	t.Run("falls back to DEFAULT when the version list is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockVersions, restoreVersions := setupIdpsSignatureVersionMock(ctrl)
		defer restoreVersions()
		mockSigs, restoreSigs := setupIdpsSignaturesMock(ctrl)
		defer restoreSigs()

		mockVersions.EXPECT().List(nil, nil, nil, nil, nil, nil, nil).Return(model.IdsSignatureVersionListResult{}, nil)
		count := int64(0)
		mockSigs.EXPECT().List("DEFAULT", (*string)(nil), nil, nil, nil, nil, nil, nil).Return(
			model.IdsSignatureListResult{ResultCount: &count}, nil,
		)

		ds := dataSourceNsxtPolicyIdpsSystemSignatures()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		err := dataSourceNsxtPolicyIdpsSystemSignaturesRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "DEFAULT", d.Get("version_id"))
	})

	t.Run("version list error propagates", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockVersions, restoreVersions := setupIdpsSignatureVersionMock(ctrl)
		defer restoreVersions()

		mockVersions.EXPECT().List(nil, nil, nil, nil, nil, nil, nil).Return(
			model.IdsSignatureVersionListResult{}, vapiErrors.InternalServerError{},
		)

		ds := dataSourceNsxtPolicyIdpsSystemSignatures()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		err := dataSourceNsxtPolicyIdpsSystemSignaturesRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
