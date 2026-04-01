//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/infra/SegmentSecurityProfilesClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt/infra/SegmentSecurityProfilesClient.go SegmentSecurityProfilesClient

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
	sspmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	segmentSecProfileID          = "segment-sec-1"
	segmentSecProfileDisplayName = "segment-sec-fooname"
	segmentSecProfileDescription = "segment security profile mock"
	segmentSecProfilePath        = "/infra/segment-security-profiles/segment-sec-1"
	segmentSecProfileRevision    = int64(1)
)

func TestMockResourceNsxtPolicySegmentSecurityProfileCreate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSSPSDK := sspmocks.NewMockSegmentSecurityProfilesClient(ctrl)
	mockWrapper := &cliinfra.SegmentSecurityProfileClientContext{
		Client:     mockSSPSDK,
		ClientType: utl.Local,
	}

	originalCli := cliSegmentSecurityProfilesClient
	defer func() { cliSegmentSecurityProfilesClient = originalCli }()
	cliSegmentSecurityProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.SegmentSecurityProfileClientContext {
		return mockWrapper
	}

	t.Run("Create success", func(t *testing.T) {
		mockSSPSDK.EXPECT().Patch(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		mockSSPSDK.EXPECT().Get(gomock.Any()).Return(model.SegmentSecurityProfile{
			Id:          &segmentSecProfileID,
			DisplayName: &segmentSecProfileDisplayName,
			Description: &segmentSecProfileDescription,
			Path:        &segmentSecProfilePath,
			Revision:    &segmentSecProfileRevision,
		}, nil)

		res := resourceNsxtPolicySegmentSecurityProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSegmentSecurityProfileData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicySegmentSecurityProfileCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, d.Id(), d.Get("nsx_id"))
	})

	t.Run("Create fails when resource already exists", func(t *testing.T) {
		mockSSPSDK.EXPECT().Get("existing-id").Return(model.SegmentSecurityProfile{Id: &segmentSecProfileID}, nil)

		res := resourceNsxtPolicySegmentSecurityProfile()
		data := minimalSegmentSecurityProfileData()
		data["nsx_id"] = "existing-id"
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicySegmentSecurityProfileCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicySegmentSecurityProfileRead(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSSPSDK := sspmocks.NewMockSegmentSecurityProfilesClient(ctrl)
	mockWrapper := &cliinfra.SegmentSecurityProfileClientContext{
		Client:     mockSSPSDK,
		ClientType: utl.Local,
	}

	originalCli := cliSegmentSecurityProfilesClient
	defer func() { cliSegmentSecurityProfilesClient = originalCli }()
	cliSegmentSecurityProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.SegmentSecurityProfileClientContext {
		return mockWrapper
	}

	t.Run("Read success", func(t *testing.T) {
		mockSSPSDK.EXPECT().Get(segmentSecProfileID).Return(model.SegmentSecurityProfile{
			Id:          &segmentSecProfileID,
			DisplayName: &segmentSecProfileDisplayName,
			Description: &segmentSecProfileDescription,
			Path:        &segmentSecProfilePath,
			Revision:    &segmentSecProfileRevision,
		}, nil)

		res := resourceNsxtPolicySegmentSecurityProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(segmentSecProfileID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicySegmentSecurityProfileRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, segmentSecProfileDisplayName, d.Get("display_name"))
		assert.Equal(t, segmentSecProfileDescription, d.Get("description"))
		assert.Equal(t, segmentSecProfilePath, d.Get("path"))
		assert.Equal(t, int(segmentSecProfileRevision), d.Get("revision"))
		assert.Equal(t, segmentSecProfileID, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicySegmentSecurityProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicySegmentSecurityProfileRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining SegmentSecurityProfile ID")
	})
}

func TestMockResourceNsxtPolicySegmentSecurityProfileUpdate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSSPSDK := sspmocks.NewMockSegmentSecurityProfilesClient(ctrl)
	mockWrapper := &cliinfra.SegmentSecurityProfileClientContext{
		Client:     mockSSPSDK,
		ClientType: utl.Local,
	}

	originalCli := cliSegmentSecurityProfilesClient
	defer func() { cliSegmentSecurityProfilesClient = originalCli }()
	cliSegmentSecurityProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.SegmentSecurityProfileClientContext {
		return mockWrapper
	}

	t.Run("Update success", func(t *testing.T) {
		mockSSPSDK.EXPECT().Patch(segmentSecProfileID, gomock.Any(), gomock.Any()).Return(nil)
		mockSSPSDK.EXPECT().Get(segmentSecProfileID).Return(model.SegmentSecurityProfile{
			Id:          &segmentSecProfileID,
			DisplayName: &segmentSecProfileDisplayName,
			Description: &segmentSecProfileDescription,
			Path:        &segmentSecProfilePath,
			Revision:    &segmentSecProfileRevision,
		}, nil)

		res := resourceNsxtPolicySegmentSecurityProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSegmentSecurityProfileData())
		d.SetId(segmentSecProfileID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicySegmentSecurityProfileUpdate(d, m)
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicySegmentSecurityProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSegmentSecurityProfileData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicySegmentSecurityProfileUpdate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining SegmentSecurityProfile ID")
	})
}

func TestMockResourceNsxtPolicySegmentSecurityProfileDelete(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSSPSDK := sspmocks.NewMockSegmentSecurityProfilesClient(ctrl)
	mockWrapper := &cliinfra.SegmentSecurityProfileClientContext{
		Client:     mockSSPSDK,
		ClientType: utl.Local,
	}

	originalCli := cliSegmentSecurityProfilesClient
	defer func() { cliSegmentSecurityProfilesClient = originalCli }()
	cliSegmentSecurityProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.SegmentSecurityProfileClientContext {
		return mockWrapper
	}

	t.Run("Delete success", func(t *testing.T) {
		mockSSPSDK.EXPECT().Delete(segmentSecProfileID, gomock.Any()).Return(nil)

		res := resourceNsxtPolicySegmentSecurityProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(segmentSecProfileID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicySegmentSecurityProfileDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicySegmentSecurityProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicySegmentSecurityProfileDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining SegmentSecurityProfile ID")
	})

	t.Run("Delete fails when API returns error", func(t *testing.T) {
		mockSSPSDK.EXPECT().Delete(segmentSecProfileID, gomock.Any()).Return(errors.New("API error"))

		res := resourceNsxtPolicySegmentSecurityProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(segmentSecProfileID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicySegmentSecurityProfileDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})
}

func minimalSegmentSecurityProfileData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": segmentSecProfileDisplayName,
		"description":  segmentSecProfileDescription,
	}
}
