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
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	enforcementpoints "github.com/vmware/terraform-provider-nsxt/api/infra/sites/enforcement_points"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	epmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/sites/enforcement_points"
)

func setupEdgeHAProfileDataSourceMock(t *testing.T, ctrl *gomock.Controller) (*epmocks.MockEdgeClusterHighAvailabilityProfilesClient, func()) {
	t.Helper()
	mockSDK := epmocks.NewMockEdgeClusterHighAvailabilityProfilesClient(ctrl)
	wrapper := &enforcementpoints.PolicyEdgeHighAvailabilityProfileClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	orig := cliEdgeClusterHighAvailabilityProfilesClient
	cliEdgeClusterHighAvailabilityProfilesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *enforcementpoints.PolicyEdgeHighAvailabilityProfileClientContext {
		return wrapper
	}
	return mockSDK, func() { cliEdgeClusterHighAvailabilityProfilesClient = orig }
}

func dsEdgeHAProfileModel() model.PolicyEdgeHighAvailabilityProfile {
	uid := "unique-ha-1"
	return model.PolicyEdgeHighAvailabilityProfile{
		Id:          &policyHAProfileID,
		DisplayName: &policyHAProfileName,
		Path:        &policyHAProfilePath,
		Revision:    &policyHAProfileRevision,
		UniqueId:    &uid,
	}
}

func TestMockDataSourceNsxtPolicyEdgeHighAvailabilityProfileRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupEdgeHAProfileDataSourceMock(t, ctrl)
	defer restore()
	m := newGoMockProviderClient()
	ep := getPolicyEnforcementPoint(m)

	t.Run("by id success", func(t *testing.T) {
		mockSDK.EXPECT().Get(policyHAProfileSiteID, ep, policyHAProfileID).Return(dsEdgeHAProfileModel(), nil)

		ds := dataSourceNsxtPolicyEdgeHighAvailabilityProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": policyHAProfileID,
		})

		err := dataSourceNsxtPolicyEdgeHighAvailabilityProfileRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, policyHAProfileID, d.Id())
		assert.Equal(t, policyHAProfileName, d.Get("display_name"))
		assert.Equal(t, "unique-ha-1", d.Get("unique_id"))
	})

	t.Run("by id not found", func(t *testing.T) {
		mockSDK.EXPECT().Get(policyHAProfileSiteID, ep, policyHAProfileID).Return(model.PolicyEdgeHighAvailabilityProfile{}, vapiErrors.NotFound{})

		ds := dataSourceNsxtPolicyEdgeHighAvailabilityProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": policyHAProfileID,
		})

		err := dataSourceNsxtPolicyEdgeHighAvailabilityProfileRead(d, m)
		require.Error(t, err)
	})

	t.Run("by display_name via list", func(t *testing.T) {
		mockSDK.EXPECT().List(policyHAProfileSiteID, ep, nil, nil, nil, nil, nil, nil, nil).Return(model.EdgeClusterHighAvailabilityProfileListResult{
			Results: []model.PolicyEdgeHighAvailabilityProfile{dsEdgeHAProfileModel()},
		}, nil)

		ds := dataSourceNsxtPolicyEdgeHighAvailabilityProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": policyHAProfileName,
		})

		err := dataSourceNsxtPolicyEdgeHighAvailabilityProfileRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, policyHAProfileID, d.Id())
	})

	t.Run("missing id and name", func(t *testing.T) {
		ds := dataSourceNsxtPolicyEdgeHighAvailabilityProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		err := dataSourceNsxtPolicyEdgeHighAvailabilityProfileRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "ID or name")
	})
}
