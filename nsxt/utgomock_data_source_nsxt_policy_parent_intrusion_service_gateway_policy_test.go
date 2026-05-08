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

	apidomains "github.com/vmware/terraform-provider-nsxt/api/infra/domains"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	domainmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/domains"
)

var (
	idsParentGwPolDomain   = "default"
	idsParentGwPolID       = "ids-parent-gw-pol-1"
	idsParentGwPolName     = "ids-parent-gw-policy"
	idsParentGwPolPath     = "/infra/domains/default/intrusion-service-gateway-policies/ids-parent-gw-pol-1"
	idsParentGwPolCategory = "LocalGatewayRules"
	idsParentGwPolDesc     = "parent gateway policy desc"
)

func idsParentGwPolicyModel() nsxModel.IdsGatewayPolicy {
	stateful := true
	locked := false
	sequenceNumber := int64(15)
	return nsxModel.IdsGatewayPolicy{
		Id:             &idsParentGwPolID,
		DisplayName:    &idsParentGwPolName,
		Description:    &idsParentGwPolDesc,
		Path:           &idsParentGwPolPath,
		Category:       &idsParentGwPolCategory,
		Stateful:       &stateful,
		Locked:         &locked,
		SequenceNumber: &sequenceNumber,
	}
}

func setupParentIntrusionServiceGwPolicyDataSourceMock(t *testing.T, ctrl *gomock.Controller) (*domainmocks.MockIntrusionServiceGatewayPoliciesClient, func()) {
	t.Helper()
	mockSDK := domainmocks.NewMockIntrusionServiceGatewayPoliciesClient(ctrl)
	wrapper := &apidomains.IntrusionServiceGatewayPolicyClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	orig := cliIntrusionServiceGatewayPoliciesClient
	cliIntrusionServiceGatewayPoliciesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apidomains.IntrusionServiceGatewayPolicyClientContext {
		return wrapper
	}
	return mockSDK, func() { cliIntrusionServiceGatewayPoliciesClient = orig }
}

func TestMockDataSourceNsxtPolicyParentIntrusionServiceGatewayPolicyRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupParentIntrusionServiceGwPolicyDataSourceMock(t, ctrl)
	defer restore()

	boolFalse := false

	t.Run("by id success", func(t *testing.T) {
		mockSDK.EXPECT().Get(idsParentGwPolDomain, idsParentGwPolID).Return(idsParentGwPolicyModel(), nil)

		ds := dataSourceNsxtPolicyParentIntrusionServiceGatewayPolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":     idsParentGwPolID,
			"domain": idsParentGwPolDomain,
		})

		err := dataSourceNsxtPolicyParentIntrusionServiceGatewayPolicyRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, idsParentGwPolID, d.Id())
		assert.Equal(t, idsParentGwPolName, d.Get("display_name"))
		assert.Equal(t, idsParentGwPolCategory, d.Get("category"))
		assert.Equal(t, true, d.Get("stateful"))
		assert.Equal(t, false, d.Get("locked"))
		assert.Equal(t, 15, d.Get("sequence_number"))
	})

	t.Run("by id not found", func(t *testing.T) {
		mockSDK.EXPECT().Get(idsParentGwPolDomain, idsParentGwPolID).Return(nsxModel.IdsGatewayPolicy{}, vapiErrors.NotFound{})

		ds := dataSourceNsxtPolicyParentIntrusionServiceGatewayPolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":     idsParentGwPolID,
			"domain": idsParentGwPolDomain,
		})

		err := dataSourceNsxtPolicyParentIntrusionServiceGatewayPolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})

	t.Run("by id API error", func(t *testing.T) {
		mockSDK.EXPECT().Get(idsParentGwPolDomain, idsParentGwPolID).Return(nsxModel.IdsGatewayPolicy{}, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtPolicyParentIntrusionServiceGatewayPolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":     idsParentGwPolID,
			"domain": idsParentGwPolDomain,
		})

		err := dataSourceNsxtPolicyParentIntrusionServiceGatewayPolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error while reading")
	})

	t.Run("missing id name and category", func(t *testing.T) {
		ds := dataSourceNsxtPolicyParentIntrusionServiceGatewayPolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"domain": idsParentGwPolDomain,
		})

		err := dataSourceNsxtPolicyParentIntrusionServiceGatewayPolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must be specified")
	})

	t.Run("by display_name via list", func(t *testing.T) {
		rc := int64(1)
		mockSDK.EXPECT().List(idsParentGwPolDomain, nil, &boolFalse, nil, nil, nil, nil, nil).Return(nsxModel.IdsGatewayPolicyListResult{
			Results:     []nsxModel.IdsGatewayPolicy{idsParentGwPolicyModel()},
			ResultCount: &rc,
		}, nil)

		ds := dataSourceNsxtPolicyParentIntrusionServiceGatewayPolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"domain":       idsParentGwPolDomain,
			"display_name": idsParentGwPolName,
		})

		err := dataSourceNsxtPolicyParentIntrusionServiceGatewayPolicyRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, idsParentGwPolID, d.Id())
		assert.Equal(t, idsParentGwPolName, d.Get("display_name"))
		assert.Equal(t, idsParentGwPolCategory, d.Get("category"))
	})

	t.Run("by category via list", func(t *testing.T) {
		rc := int64(1)
		mockSDK.EXPECT().List(idsParentGwPolDomain, nil, &boolFalse, nil, nil, nil, nil, nil).Return(nsxModel.IdsGatewayPolicyListResult{
			Results:     []nsxModel.IdsGatewayPolicy{idsParentGwPolicyModel()},
			ResultCount: &rc,
		}, nil)

		ds := dataSourceNsxtPolicyParentIntrusionServiceGatewayPolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"domain":   idsParentGwPolDomain,
			"category": idsParentGwPolCategory,
		})

		err := dataSourceNsxtPolicyParentIntrusionServiceGatewayPolicyRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, idsParentGwPolID, d.Id())
		assert.Equal(t, idsParentGwPolName, d.Get("display_name"))
		assert.Equal(t, idsParentGwPolCategory, d.Get("category"))
	})

	t.Run("list error", func(t *testing.T) {
		mockSDK.EXPECT().List(idsParentGwPolDomain, nil, &boolFalse, nil, nil, nil, nil, nil).Return(nsxModel.IdsGatewayPolicyListResult{}, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtPolicyParentIntrusionServiceGatewayPolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"domain":       idsParentGwPolDomain,
			"display_name": idsParentGwPolName,
		})

		err := dataSourceNsxtPolicyParentIntrusionServiceGatewayPolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("multiple matches error", func(t *testing.T) {
		rc := int64(2)
		policy2 := idsParentGwPolicyModel()
		id2 := "ids-parent-gw-pol-2"
		policy2.Id = &id2

		mockSDK.EXPECT().List(idsParentGwPolDomain, nil, &boolFalse, nil, nil, nil, nil, nil).Return(nsxModel.IdsGatewayPolicyListResult{
			Results:     []nsxModel.IdsGatewayPolicy{idsParentGwPolicyModel(), policy2},
			ResultCount: &rc,
		}, nil)

		ds := dataSourceNsxtPolicyParentIntrusionServiceGatewayPolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"domain":       idsParentGwPolDomain,
			"display_name": idsParentGwPolName,
		})

		err := dataSourceNsxtPolicyParentIntrusionServiceGatewayPolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Found multiple")
	})
}
