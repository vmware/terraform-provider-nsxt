//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/infra/GatewayQosProfilesClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt/infra/GatewayQosProfilesClient.go GatewayQosProfilesClient

package nsxt

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	cliinfra "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	qosmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	gwQosProfileID           = "gw-qos-profile-1"
	gwQosProfileDisplayName  = "gw-qos-fooname"
	gwQosProfileDescription  = "gateway qos profile mock"
	gwQosProfilePath         = "/infra/gateway-qos-profiles/gw-qos-profile-1"
	gwQosProfileRevision     = int64(1)
	gwQosProfileBurstSize    = int64(1)
	gwQosProfileCommittedBw  = int64(1)
	gwQosProfileExcessAction = model.GatewayQosProfile_EXCESS_ACTION_DROP
)

func TestMockResourceNsxtPolicyGatewayQosProfileCreate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQosSDK := qosmocks.NewMockGatewayQosProfilesClient(ctrl)
	mockWrapper := &cliinfra.GatewayQosProfileClientContext{
		Client:     mockQosSDK,
		ClientType: utl.Local,
	}

	originalCli := cliGatewayQosProfilesClient
	defer func() { cliGatewayQosProfilesClient = originalCli }()
	cliGatewayQosProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.GatewayQosProfileClientContext {
		return mockWrapper
	}

	t.Run("Create success", func(t *testing.T) {
		mockQosSDK.EXPECT().Patch(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		mockQosSDK.EXPECT().Get(gomock.Any()).Return(model.GatewayQosProfile{
			Id:                 &gwQosProfileID,
			DisplayName:        &gwQosProfileDisplayName,
			Description:        &gwQosProfileDescription,
			Path:               &gwQosProfilePath,
			Revision:           &gwQosProfileRevision,
			BurstSize:          &gwQosProfileBurstSize,
			CommittedBandwidth: &gwQosProfileCommittedBw,
			ExcessAction:       &gwQosProfileExcessAction,
		}, nil)

		res := resourceNsxtPolicyGatewayQosProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGatewayQosProfileData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyGatewayQosProfileCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, d.Id(), d.Get("nsx_id"))
	})

	t.Run("Create fails when resource already exists", func(t *testing.T) {
		mockQosSDK.EXPECT().Get("existing-id").Return(model.GatewayQosProfile{Id: &gwQosProfileID}, nil)

		res := resourceNsxtPolicyGatewayQosProfile()
		data := minimalGatewayQosProfileData()
		data["nsx_id"] = "existing-id"
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyGatewayQosProfileCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicyGatewayQosProfileRead(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQosSDK := qosmocks.NewMockGatewayQosProfilesClient(ctrl)
	mockWrapper := &cliinfra.GatewayQosProfileClientContext{
		Client:     mockQosSDK,
		ClientType: utl.Local,
	}

	originalCli := cliGatewayQosProfilesClient
	defer func() { cliGatewayQosProfilesClient = originalCli }()
	cliGatewayQosProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.GatewayQosProfileClientContext {
		return mockWrapper
	}

	t.Run("Read success", func(t *testing.T) {
		mockQosSDK.EXPECT().Get(gwQosProfileID).Return(model.GatewayQosProfile{
			Id:                 &gwQosProfileID,
			DisplayName:        &gwQosProfileDisplayName,
			Description:        &gwQosProfileDescription,
			Path:               &gwQosProfilePath,
			Revision:           &gwQosProfileRevision,
			BurstSize:          &gwQosProfileBurstSize,
			CommittedBandwidth: &gwQosProfileCommittedBw,
			ExcessAction:       &gwQosProfileExcessAction,
		}, nil)

		res := resourceNsxtPolicyGatewayQosProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(gwQosProfileID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyGatewayQosProfileRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, gwQosProfileDisplayName, d.Get("display_name"))
		assert.Equal(t, gwQosProfileDescription, d.Get("description"))
		assert.Equal(t, gwQosProfilePath, d.Get("path"))
		assert.Equal(t, int(gwQosProfileRevision), d.Get("revision"))
		assert.Equal(t, gwQosProfileID, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyGatewayQosProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyGatewayQosProfileRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining GatewayQosProfile ID")
	})
}

func TestMockResourceNsxtPolicyGatewayQosProfileUpdate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQosSDK := qosmocks.NewMockGatewayQosProfilesClient(ctrl)
	mockWrapper := &cliinfra.GatewayQosProfileClientContext{
		Client:     mockQosSDK,
		ClientType: utl.Local,
	}

	originalCli := cliGatewayQosProfilesClient
	defer func() { cliGatewayQosProfilesClient = originalCli }()
	cliGatewayQosProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.GatewayQosProfileClientContext {
		return mockWrapper
	}

	t.Run("Update success", func(t *testing.T) {
		mockQosSDK.EXPECT().Patch(gwQosProfileID, gomock.Any(), gomock.Any()).Return(nil)
		mockQosSDK.EXPECT().Get(gwQosProfileID).Return(model.GatewayQosProfile{
			Id:                 &gwQosProfileID,
			DisplayName:        &gwQosProfileDisplayName,
			Description:        &gwQosProfileDescription,
			Path:               &gwQosProfilePath,
			Revision:           &gwQosProfileRevision,
			BurstSize:          &gwQosProfileBurstSize,
			CommittedBandwidth: &gwQosProfileCommittedBw,
			ExcessAction:       &gwQosProfileExcessAction,
		}, nil)

		res := resourceNsxtPolicyGatewayQosProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGatewayQosProfileData())
		d.SetId(gwQosProfileID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyGatewayQosProfileUpdate(d, m)
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyGatewayQosProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGatewayQosProfileData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyGatewayQosProfileUpdate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining GatewayQosProfile ID")
	})
}

func TestMockResourceNsxtPolicyGatewayQosProfileDelete(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQosSDK := qosmocks.NewMockGatewayQosProfilesClient(ctrl)
	mockWrapper := &cliinfra.GatewayQosProfileClientContext{
		Client:     mockQosSDK,
		ClientType: utl.Local,
	}

	originalCli := cliGatewayQosProfilesClient
	defer func() { cliGatewayQosProfilesClient = originalCli }()
	cliGatewayQosProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.GatewayQosProfileClientContext {
		return mockWrapper
	}

	t.Run("Delete success", func(t *testing.T) {
		mockQosSDK.EXPECT().Delete(gwQosProfileID, nil).Return(nil)

		res := resourceNsxtPolicyGatewayQosProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(gwQosProfileID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyGatewayQosProfileDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyGatewayQosProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyGatewayQosProfileDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining GatewayQosProfile ID")
	})

	t.Run("Delete fails when API returns error", func(t *testing.T) {
		mockQosSDK.EXPECT().Delete(gwQosProfileID, nil).Return(errors.New("API error"))

		res := resourceNsxtPolicyGatewayQosProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(gwQosProfileID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyGatewayQosProfileDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})
}

func minimalGatewayQosProfileData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":        gwQosProfileDisplayName,
		"description":         gwQosProfileDescription,
		"burst_size":          int(gwQosProfileBurstSize),
		"committed_bandwidth": int(gwQosProfileCommittedBw),
		"excess_action":       gwQosProfileExcessAction,
	}
}
